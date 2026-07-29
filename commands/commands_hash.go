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
func executeHDEL(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HDEL' command")
	}
	count, err := cache.HDel(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(count), nil
}
func executeHEXISTS(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HEXISTS' command")
	}
	exists, err := cache.HExists(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}
	if !exists {
		return protocol.NewInteger(0), nil
	}
	return protocol.NewInteger(1), nil
}
func executeHLEN(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HLEN' command")
	}
	count, err := cache.HLen(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(count), nil
}
