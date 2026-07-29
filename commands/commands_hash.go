package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
	"strconv"
)

func executeHSET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 4 || (len(tokens)-2)%2 != 0 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HSET' command")
	}
	added := 0

	for i := 2; i < len(tokens); i += 2 {
		count, err := cache.HSet(tokens[1], tokens[i], tokens[i+1])
		if err != nil {
			return protocol.Response{}, err
		}

		added += count
	}

	return protocol.NewInteger(added), nil
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
func executeHGETALL(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HGETALL' command")
	}
	values, err := cache.HGetAll(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}
	elements := make([]protocol.Response, 0, len(values))

	for _, value := range values {
		elements = append(elements, protocol.NewBulkString(value))
	}

	return protocol.NewArray(elements), nil

}
func executeHKEYS(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HKEYS' command")
	}

	keys, err := cache.HKeys(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}

	elements := make([]protocol.Response, 0, len(keys))

	for _, key := range keys {
		elements = append(elements, protocol.NewBulkString(key))
	}

	return protocol.NewArray(elements), nil
}
func executeHVALS(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HVALS' command")
	}

	values, err := cache.HVals(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}

	elements := make([]protocol.Response, 0, len(values))

	for _, key := range values {
		elements = append(elements, protocol.NewBulkString(key))
	}

	return protocol.NewArray(elements), nil
}
func executeHINCRBY(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 4 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HINCRBY' command")
	}

	increment, err := strconv.Atoi(tokens[3])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("ERR value is not an integer")
	}

	newValue, err := cache.HIncrBy(tokens[1], tokens[2], increment)
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewInteger(newValue), nil
}
func executeHMGET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HMGET' command")
	}

	elements := make([]protocol.Response, 0, len(tokens)-2)

	for i := 2; i < len(tokens); i++ {
		value, found, err := cache.HGet(tokens[1], tokens[i])
		if err != nil {
			return protocol.Response{}, err
		}

		if !found {
			elements = append(elements, protocol.NewNullBulkString())
		} else {
			elements = append(elements, protocol.NewBulkString(value))
		}
	}

	return protocol.NewArray(elements), nil
}
func executeHSETNX(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'HSETNX' command")
	}

	count, err := cache.HSetNX(tokens[1], tokens[2], tokens[3])
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewInteger(count), nil
}
