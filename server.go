package main

import (
	"fmt"
	"go-redis/aof"
	"go-redis/cache"
	"go-redis/server"
	"log"
	"net"
)

const (
	replayAOF  = true
	rewriteAOF = false
)

func main() {
	cache := cache.NewCache()

	aof, err := aof.NewAOF("appendonly.aof")
	if err != nil {
		log.Fatal(err)
	}
	defer aof.Close()

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
		go server.HandleClient(conn, cache, aof)
	}
}
