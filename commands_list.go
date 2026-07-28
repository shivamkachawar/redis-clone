package main

import (
	"fmt"
	"strconv"
)

func executeLPUSH(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) < 3 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'LPUSH' command")
	}

	length, err := cache.LPush(tokens[1], tokens[2:])
	if err != nil {
		return Response{}, err
	}
	return NewInteger(length), nil
}
func executeRPUSH(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) < 3 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'RPUSH' command")
	}

	length, err := cache.RPush(tokens[1], tokens[2:])
	if err != nil {
		return Response{}, err
	}

	return NewInteger(length), nil
}
func executeLLEN(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'LLEN' command")
	}

	length, err := cache.LLen(tokens[1])
	if err != nil {
		return Response{}, err
	}
	return NewInteger(length), nil
}
func executeLRANGE(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 4 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'LRANGE' command")
	}

	start, err := strconv.Atoi(tokens[2])
	if err != nil {
		return Response{}, fmt.Errorf("ERR value is not an integer or out of range")
	}

	stop, err := strconv.Atoi(tokens[3])
	if err != nil {
		return Response{}, fmt.Errorf("ERR value is not an integer or out of range")
	}

	list, err := cache.LRange(tokens[1], start, stop)
	if err != nil {
		return Response{}, err
	}

	responses := make([]Response, len(list))
	for i, value := range list {
		responses[i] = NewBulkString(value)
	}

	return NewArray(responses), nil
}
func executeLPOP(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'LPOP' command")
	}
	value, exists, err := cache.LPop(tokens[1])
	if err != nil {
		return Response{}, err
	}

	if !exists {
		return NewNullBulkString(), nil
	}

	return NewBulkString(value), nil
}
func executeRPOP(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'RPOP' command")
	}

	value, exists, err := cache.RPop(tokens[1])
	if err != nil {
		return Response{}, err
	}

	if !exists {
		return NewNullBulkString(), nil
	}

	return NewBulkString(value), nil
}
