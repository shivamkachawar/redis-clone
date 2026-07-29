package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
	"strconv"
)

func executeEXPIRE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("Error: EXPIRE command requires a key and a time in seconds")
	}
	key := tokens[1]
	seconds, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("Error: EXPIRE command requires a valid time in seconds")
	}
	if seconds <= 0 {
		return protocol.Response{}, fmt.Errorf("Error: Expiry time must be greater than 0")
	}
	if cache.Expire(key, seconds) {
		return protocol.NewInteger(1), nil
	}
	return protocol.NewInteger(0), nil
}
func executeTTL(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: TTL command requires a key")
	}
	key := tokens[1]
	ttl := cache.TTL(key)
	return protocol.NewInteger(ttl), nil
}
func executePERSIST(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: PERSIST command requires a key")
	}
	key := tokens[1]
	if cache.Persist(key) {
		return protocol.NewInteger(1), nil
	}
	return protocol.NewInteger(0), nil
}
