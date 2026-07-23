package main

import (
	"fmt"
	"net"
	"strings"
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
	for {
		buffer := make([]byte, 1024)
		n, err := conn.Read(buffer)

		if err != nil {
			fmt.Println("Client Disconnected : ", conn.RemoteAddr())
			break
		}

		message := string(buffer[:n])
		cleanedMessage := strings.TrimSpace(message)
		tokens := strings.Fields(cleanedMessage)

		if len(tokens) == 0 {
			continue
		}

		response, err := executeCommand(tokens, cache)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}
		conn.Write([]byte(response + "\n"))
	}
}
