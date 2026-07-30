package server

import (
	"fmt"
	"go-redis/aof"
	"go-redis/cache"
	"net"
)

func Run(cache *cache.Cache, aofFile *aof.AOF) error {
	listener, err := net.Listen("tcp", ":6379")
	if err != nil {
		return err
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
		go HandleClient(conn, cache, aofFile)
	}
}
