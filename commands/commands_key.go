package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
)

func executeEXISTS(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("Error: EXISTS command requires a key")
	}
	count := cache.Exists(tokens[1:])
	return protocol.NewInteger(count), nil
}
func executeDEL(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("Error: DELETE command requires a key")
	}
	count := cache.Delete(tokens[1:])
	return protocol.NewInteger(count), nil
}
func executeKEYS(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'KEYS' command")
	}
	if tokens[1] != "*" {
		return protocol.Response{}, fmt.Errorf("ERR only KEYS * is currently supported")
	}

	keys := cache.Keys()
	responses := make([]protocol.Response, 0, len(keys))
	for _, key := range keys {
		responses = append(responses, protocol.NewBulkString(key))
	}
	return protocol.NewArray(responses), nil
}
func executeDBSIZE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 1 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'DBSIZE' command")
	}
	size := cache.Size()
	return protocol.NewInteger(size), nil
}
func executeFLUSHDB(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 1 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'FLUSHDB' command")
	}
	cache.Flush()
	return protocol.NewSimpleString("OK"), nil
}
