package server

import (
	"bufio"
	"fmt"
	"go-redis/aof"
	"go-redis/cache"
	"go-redis/commands"
	"go-redis/protocol"
	"net"
	"strings"
)

func HandleClient(conn net.Conn, cache *cache.Cache, aof *aof.AOF) {
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
	case "SET", "MSET", "DEL", "APPEND", "GETSET", "INCR", "DECR", "INCRBY", "DECRBY", "EXPIRE", "PERSIST", "FLUSHDB", "LPUSH", "HSET", "HINCRBY", "HSETNX", "SADD":
		return true
	default:
		return false
	}
}
