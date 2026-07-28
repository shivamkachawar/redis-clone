package main

import (
	"bufio"
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
	reader := bufio.NewReader(conn)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Client Disconnected:", conn.RemoteAddr())
			break
		}

		cleanedMessage := strings.TrimSpace(line)
		tokens := strings.Fields(cleanedMessage)

		if len(tokens) == 0 {
			continue
		}

		response, err := executeCommand(tokens, cache)
		if err != nil {
			conn.Write([]byte(err.Error() + "\n"))
			continue
		}
		_, err = conn.Write([]byte(response + "\n"))
		if err != nil {
			fmt.Println("Write error:", err)
			break
		}
	}
}
