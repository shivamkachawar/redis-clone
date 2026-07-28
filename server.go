package main

import (
	"bufio"
	"fmt"
	"net"
)

func main() {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println("Server listening on Port 6379")
	cache := NewCache()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}
		fmt.Println("Client Connected : ", conn.RemoteAddr())
		go handleClient(conn, cache)
	}
}
func handleClient(conn net.Conn, cache *Cache) {
	reader := bufio.NewReader(conn)

	for {
		tokens, err := parseRESP(reader)
		if err != nil {
			fmt.Println("Parse error:", err)
			break
		}

		response, err := executeCommand(tokens, cache)
		if err != nil {
			writeError(conn, err.Error())
			continue
		}

		switch response.Type {
		case SimpleString:
			err = writeSimpleString(conn, response.Value)

		case BulkString:
			err = writeBulkString(conn, response.Value)

		case Integer:
			err = writeInteger(conn, response.Integer)

		case NullBulkString:
			err = writeNullBulkString(conn)

		case Error:
			err = writeError(conn, response.Value)

		default:
			err = writeError(conn, "Unknown response type")
		}

		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
