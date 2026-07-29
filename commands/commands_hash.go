package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
)

func executeHSET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 4 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HSET' command")
	}
	count, err := cache.HSet(tokens[1], tokens[2], tokens[3])
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(count), nil
}
func executeHGET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HGET' command")
	}
	value, found, err := cache.HGet(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}
	if found == false {
		return protocol.NewNullBulkString(), nil
	}
	return protocol.NewBulkString(value), nil

}
