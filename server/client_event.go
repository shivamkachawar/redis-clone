package server

import (
	"errors"
	"fmt"
	"go-redis/aof"
	"go-redis/cache"
	"go-redis/commands"
	"go-redis/protocol"
	"syscall"
)

func HandleReadableClient(client *Client, cache *cache.Cache, aofFile *aof.AOF, pubsub *PubSub) bool {

	buffer := make([]byte, 4096)

	n, err := syscall.Read(client.FD, buffer)
	if err != nil {
		fmt.Println("Read Error:", err)
		return false
	}

	if n == 0 {
		fmt.Println("Client Disconnected")
		return false
	}
	client.InputBuffer = append(client.InputBuffer, buffer[:n]...)
	// Append newly received data to the client's persistent buffer.
	for {

		tokens, consumed, err := protocol.ParseRESP(client.InputBuffer)

		if len(tokens) > 0 && tokens[0] == "MULTI" {
			if client.InTransaction {
				resp := protocol.EncodeError("ERR MULTI calls can not be nested")
				syscall.Write(client.FD, resp)
			} else {
				client.InTransaction = true
				client.TransactionQueue = nil

				resp, _ := protocol.EncodeResponse(protocol.NewSimpleString("OK"))
				syscall.Write(client.FD, resp)
			}

			client.InputBuffer = client.InputBuffer[consumed:]
			continue
		}

		if err != nil {

			if errors.Is(err, protocol.ErrIncomplete) {
				// Wait for more bytes from the socket.
				return true
			}

			fmt.Println("Parse Error:", err)
			return true
		}
		if len(tokens) > 0 && tokens[0] == "DISCARD" {
			if !client.InTransaction {
				resp := protocol.EncodeError("ERR DISCARD without MULTI")
				syscall.Write(client.FD, resp)
			} else {
				client.InTransaction = false
				client.TransactionQueue = nil
				resp, _ := protocol.EncodeResponse(protocol.NewSimpleString("OK"))
				syscall.Write(client.FD, resp)
			}

			client.InputBuffer = client.InputBuffer[consumed:]
			continue
		}
		if len(tokens) > 0 && tokens[0] == "EXEC" {

			response, err := executeTransaction(client, cache, aofFile)

			client.InputBuffer = client.InputBuffer[consumed:]

			if err != nil {
				resp := protocol.EncodeError(err.Error())
				syscall.Write(client.FD, resp)
				continue
			}

			resp, err := protocol.EncodeResponse(response)
			if err != nil {
				fmt.Println("Encode Error:", err)
				return true
			}

			_, err = syscall.Write(client.FD, resp)
			if err != nil {
				fmt.Println("Write Error:", err)
				return true
			}

			continue
		}
		if len(tokens) > 0 && tokens[0] == "SUBSCRIBE" {

			response, err := executeSUBSCRIBE(client, pubsub, tokens)

			client.InputBuffer = client.InputBuffer[consumed:]

			if err != nil {
				resp := protocol.EncodeError(err.Error())
				syscall.Write(client.FD, resp)
				continue
			}

			resp, err := protocol.EncodeResponse(response)
			if err != nil {
				fmt.Println("Encode Error:", err)
				return true
			}

			_, err = syscall.Write(client.FD, resp)
			if err != nil {
				fmt.Println("Write Error:", err)
				return true
			}

			continue
		}

		if len(tokens) > 0 && tokens[0] == "PUBLISH" {

			response, err := executePUBLISH(pubsub, tokens)

			client.InputBuffer = client.InputBuffer[consumed:]

			if err != nil {
				resp := protocol.EncodeError(err.Error())
				syscall.Write(client.FD, resp)
				continue
			}

			resp, err := protocol.EncodeResponse(response)
			if err != nil {
				fmt.Println("Encode Error:", err)
				return true
			}

			_, err = syscall.Write(client.FD, resp)
			if err != nil {
				fmt.Println("Write Error:", err)
				return true
			}

			continue
		}
		if client.InTransaction {
			client.TransactionQueue = append(client.TransactionQueue, tokens)

			resp, _ := protocol.EncodeResponse(protocol.NewSimpleString("QUEUED"))
			syscall.Write(client.FD, resp)

			client.InputBuffer = client.InputBuffer[consumed:]
			continue
		}

		response, err := commands.ExecuteCommand(tokens, cache)

		if err == nil && len(tokens) > 0 && isWriteCommand(tokens[0]) {
			if err := aofFile.Write(tokens); err != nil {
				fmt.Println("AOF write error:", err)
				return true
			}
		}

		if err != nil {
			// Remove the command that was parsed successfully.
			client.InputBuffer = client.InputBuffer[consumed:]

			resp := protocol.EncodeError(err.Error())

			_, writeErr := syscall.Write(client.FD, resp)
			if writeErr != nil {
				fmt.Println("Write Error:", writeErr)
			}
			return true
		}

		resp, err := protocol.EncodeResponse(response)
		if err != nil {
			fmt.Println("Encode Error:", err)
			return true
		}

		_, err = syscall.Write(client.FD, resp)
		if err != nil {
			fmt.Println("Write Error:", err)
			return true
		}

		// Remove the processed command from the buffer.
		client.InputBuffer = client.InputBuffer[consumed:]
		if len(client.InputBuffer) == 0 {
			break
		}
	}
	return true
}
