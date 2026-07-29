package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
	"strconv"
)

func executeLPUSH(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'LPUSH' command")
	}

	length, err := cache.LPush(tokens[1], tokens[2:])
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(length), nil
}
func executeRPUSH(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'RPUSH' command")
	}

	length, err := cache.RPush(tokens[1], tokens[2:])
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewInteger(length), nil
}
func executeLLEN(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'LLEN' command")
	}

	length, err := cache.LLen(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(length), nil
}
func executeLRANGE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'LRANGE' command")
	}

	start, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("ERR value is not an integer or out of range")
	}

	stop, err := strconv.Atoi(tokens[3])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("ERR value is not an integer or out of range")
	}

	list, err := cache.LRange(tokens[1], start, stop)
	if err != nil {
		return protocol.Response{}, err
	}

	responses := make([]protocol.Response, len(list))
	for i, value := range list {
		responses[i] = protocol.NewBulkString(value)
	}

	return protocol.NewArray(responses), nil
}
func executeLPOP(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'LPOP' command")
	}
	value, exists, err := cache.LPop(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}

	if !exists {
		return protocol.NewNullBulkString(), nil
	}

	return protocol.NewBulkString(value), nil
}
func executeRPOP(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'RPOP' command")
	}

	value, exists, err := cache.RPop(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}

	if !exists {
		return protocol.NewNullBulkString(), nil
	}

	return protocol.NewBulkString(value), nil
}
