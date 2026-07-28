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
		err = writeResponse(conn, response)

		if err != nil {
			fmt.Println("Write error:", err)
			break
		}

	}
}
