package main

import (
	"bufio"
	"fmt"
	"go-redis/cache"
	"go-redis/commands"
	"go-redis/protocol"
	"log"
	"net"
	"strings"
)

const (
	replayAOF  = true
	rewriteAOF = false
)

func main() {
	cache := cache.NewCache()

	aof, err := NewAOF("appendonly.aof")
	if err != nil {
		log.Fatal(err)
	}
	defer aof.file.Close()

	if replayAOF {
		err = aof.Replay(cache)
		if err != nil {
			log.Fatal(err)
		}
	}

	if rewriteAOF {
		err = aof.Rewrite(cache)
		if err != nil {
			log.Fatal(err)
		}
	}
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listener.Close()

	fmt.Println("Server listening on Port 6379")

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println(err)
			continue
		}

		fmt.Println("Client Connected:", conn.RemoteAddr())
		go handleClient(conn, cache, aof)
	}
}
func handleClient(conn net.Conn, cache *cache.Cache, aof *AOF) {
	reader := bufio.NewReader(conn)

	for {
		tokens, err := protocol.ParseRESP(reader)
		if err != nil {
			fmt.Println("Parse error:", err)
			break
		}

		response, err := commands.ExecuteCommand(tokens, cache)
		if err == nil && isWriteCommand(tokens[0]) {
			if err := aof.Write(tokens); err != nil {
				fmt.Println("AOF write error:", err)
				break
			}
		}
		if err != nil {
			protocol.WriteError(conn, err.Error())
			continue
		}
		err = protocol.WriteResponse(conn, response)

		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
func isWriteCommand(command string) bool {
	switch strings.ToUpper(command) {
	case "SET", "MSET", "DEL", "APPEND", "GETSET", "INCR", "DECR", "INCRBY", "DECRBY", "EXPIRE", "PERSIST", "FLUSHDB", "LPUSH", "HSET", "HINCRBY", "HSETNX":
		return true
	default:
		return false
	}
}
