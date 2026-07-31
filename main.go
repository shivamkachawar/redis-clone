package main

import (
	"go-redis/aof"
	"go-redis/cache"
	"go-redis/server"
	"log"
)

const (
	replayAOF  = true
	rewriteAOF = false
)

func main() {
	cache := cache.NewCache(3)

	aofFile, err := aof.NewAOF("appendonly.aof")
	if err != nil {
		log.Fatal(err)
	}
	defer aofFile.Close()

	if replayAOF {
		err = aofFile.Replay(cache)
		if err != nil {
			log.Fatal(err)
		}
	}

	if rewriteAOF {
		err = aofFile.Rewrite(cache)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = server.Run(cache, aofFile)
	if err != nil {
		log.Fatal(err)
	}

}
