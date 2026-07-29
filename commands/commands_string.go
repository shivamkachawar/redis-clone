package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
	"strconv"
	"strings"
)

func executeSET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 && len(tokens) != 5 {
		return protocol.Response{}, fmt.Errorf("Usage: SET <key> <value> [EX <seconds>]")
	}

	key := tokens[1]
	value := tokens[2]

	// Normal SET
	if len(tokens) == 3 {
		cache.Set(key, value, 0)
		return protocol.NewSimpleString("OK"), nil
	}

	// SET with EX
	if strings.ToUpper(tokens[3]) != "EX" {
		return protocol.Response{}, fmt.Errorf("Error: Expected EX")
	}

	seconds, err := strconv.Atoi(tokens[4])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("Error: Invalid expiry time")
	}

	if seconds <= 0 {
		return protocol.Response{}, fmt.Errorf("Error: TTL must be greater than 0")
	}

	cache.Set(key, value, seconds)
	return protocol.NewSimpleString("OK"), nil
}
func executeGET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: GET command requires a key")
	}
	key := tokens[1]
	value, exists, err := cache.Get(key)
	if err != nil {
		return protocol.Response{}, err
	}
	if !exists {
		return protocol.NewNullBulkString(), nil
	}
	return protocol.NewBulkString(value), nil
}
func executeAPPEND(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("Error: APPEND requires exactly two arguments")
	}
	length := cache.Append(tokens[1], tokens[2])

	return protocol.NewInteger(length), nil
}
func executeSTRLEN(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: STRLEN requires exactly one argument")
	}
	length := cache.Strlen(tokens[1])

	return protocol.NewInteger(length), nil
}
func executeGETSET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("Error: GETSET requires exactly two arguments")
	}
	oldValue, exists := cache.GetSet(tokens[1], tokens[2])

	if !exists {
		return protocol.NewNullBulkString(), nil
	}
	return protocol.NewBulkString(oldValue), nil
}
func executeINCR(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: INCR command requires a key")
	}

	key := tokens[1]

	newValue, err := cache.Increment(key)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(newValue), nil
}
func executeDECR(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("Error: DECR command requires a key")
	}

	key := tokens[1]

	newValue, err := cache.Decrement(key)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(newValue), nil
}
func executeINCRBY(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("Error: INCRBY command requires a key and a delta")
	}

	key := tokens[1]
	delta, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("Error: INCRBY command requires a valid integer delta")
	}

	newValue, err := cache.IncrementBy(key, delta)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(newValue), nil
}
func executeDECRBY(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("Error: DECRBY command requires a key and a delta")
	}

	key := tokens[1]
	delta, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, fmt.Errorf("Error: DECRBY command requires a valid integer delta")
	}

	newValue, err := cache.DecrementBy(key, delta)
	if err != nil {
		return protocol.Response{}, err
	}
	return protocol.NewInteger(newValue), nil
}
func executeMSET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, fmt.Errorf("Error: MSET command requires at least one key-value pair")
	}
	if len(tokens)%2 == 0 {
		return protocol.Response{}, fmt.Errorf("Error: MSET required key-vlue pairs, but got an odd number of arguments")
	}

	var keys []string
	var values []string

	for i := 1; i < len(tokens); i += 2 {
		keys = append(keys, tokens[i])
		values = append(values, tokens[i+1])
	}

	cache.MSet(keys, values)
	return protocol.NewSimpleString("OK"), nil
}
func executeMGET(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 2 {
		return protocol.Response{}, fmt.Errorf("Error: MGET command requires at least one key")
	}

	keys := tokens[1:]
	responses, err := cache.MGet(keys)
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewArray(responses), nil
}
