package main

import (
	"fmt"
	"strconv"
	"strings"
)

func executeCommand(tokens []string, cache *Cache) (Response, error) {
	command := tokens[0]
	switch command {
	case "PING":
		return NewSimpleString("PONG"), nil
	case "SET":
		if len(tokens) != 3 && len(tokens) != 5 {
			return Response{}, fmt.Errorf("Usage: SET <key> <value> [EX <seconds>]")
		}

		key := tokens[1]
		value := tokens[2]

		// Normal SET
		if len(tokens) == 3 {
			cache.Set(key, value, 0)
			return NewSimpleString("OK"), nil
		}

		// SET with EX
		if strings.ToUpper(tokens[3]) != "EX" {
			return Response{}, fmt.Errorf("Error: Expected EX")
		}

		seconds, err := strconv.Atoi(tokens[4])
		if err != nil {
			return Response{}, fmt.Errorf("Error: Invalid expiry time")
		}

		if seconds <= 0 {
			return Response{}, fmt.Errorf("Error: TTL must be greater than 0")
		}

		cache.Set(key, value, seconds)
		return NewSimpleString("OK"), nil
	case "GET":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: GET command requires a key")
		}
		key := tokens[1]
		value, exists := cache.Get(key)
		if !exists {
			return NewNullBulkString(), nil
		}
		return NewBulkString(value), nil
	case "DEL":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: DELETE command requires a key")
		}
		key := tokens[1]
		cache.Delete(key)
		return NewSimpleString("OK"), nil
	case "EXPIRE":
		if len(tokens) != 3 {
			return Response{}, fmt.Errorf("Error: EXPIRE command requires a key and a time in seconds")
		}
		key := tokens[1]
		seconds, err := strconv.Atoi(tokens[2])
		if err != nil {
			return Response{}, fmt.Errorf("Error: EXPIRE command requires a valid time in seconds")
		}
		if seconds <= 0 {
			return Response{}, fmt.Errorf("Error: Expiry time must be greater than 0")
		}
		if cache.Expire(key, seconds) {
			return NewInteger(1), nil
		}
		return NewInteger(0), nil
	case "TTL":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: TTL command requires a key")
		}
		key := tokens[1]
		ttl := cache.TTL(key)
		return NewInteger(ttl), nil
	case "PERSIST":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: PERSIST command requires a key")
		}
		key := tokens[1]
		if cache.Persist(key) {
			return NewInteger(1), nil
		}
		return NewInteger(0), nil
	case "EXISTS":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: EXISTS command requires a key")
		}
		key := tokens[1]
		if cache.Exists(key) {
			return NewInteger(1), nil
		}
		return NewInteger(0), nil
	case "INCR":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: INCR command requires a key")
		}

		key := tokens[1]

		newValue, err := cache.Increment(key)
		if err != nil {
			return Response{}, err
		}
		return NewInteger(newValue), nil
	case "DECR":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: DECR command requires a key")
		}

		key := tokens[1]

		newValue, err := cache.Decrement(key)
		if err != nil {
			return Response{}, err
		}
		return NewInteger(newValue), nil
	case "INCRBY":
		if len(tokens) != 3 {
			return Response{}, fmt.Errorf("Error: INCRBY command requires a key and a delta")
		}

		key := tokens[1]
		delta, err := strconv.Atoi(tokens[2])
		if err != nil {
			return Response{}, fmt.Errorf("Error: INCRBY command requires a valid integer delta")
		}

		newValue, err := cache.IncrementBy(key, delta)
		if err != nil {
			return Response{}, err
		}
		return NewInteger(newValue), nil
	case "DECRBY":
		if len(tokens) != 3 {
			return Response{}, fmt.Errorf("Error: DECRBY command requires a key and a delta")
		}

		key := tokens[1]
		delta, err := strconv.Atoi(tokens[2])
		if err != nil {
			return Response{}, fmt.Errorf("Error: DECRBY command requires a valid integer delta")
		}

		newValue, err := cache.DecrementBy(key, delta)
		if err != nil {
			return Response{}, err
		}
		return NewInteger(newValue), nil
	case "MGET":
		if len(tokens) < 2 {
			return Response{}, fmt.Errorf("Error: MGET command requires at least one key")
		}

		keys := tokens[1:]
		responses := cache.MGet(keys)

		return NewArray(responses), nil
	case "MSET":
		if len(tokens) < 3 {
			return Response{}, fmt.Errorf("Error: MSET command requires at least one key-value pair")
		}
		if len(tokens)%2 == 0 {
			return Response{}, fmt.Errorf("Error: MSET required key-vlue pairs, but got an odd number of arguments")
		}

		var keys []string
		var values []string

		for i := 1; i < len(tokens); i += 2 {
			keys = append(keys, tokens[i])
			values = append(values, tokens[i+1])
		}

		cache.MSet(keys, values)
		return NewSimpleString("OK"), nil
	case "APPEND":
		if len(tokens) != 3 {
			return Response{}, fmt.Errorf("Error: APPEND requires exactly two arguments")
		}
		length := cache.Append(tokens[1], tokens[2])

		return NewInteger(length), nil
	case "STRLEN":
		if len(tokens) != 2 {
			return Response{}, fmt.Errorf("Error: STRLEN requires exactly one argument")
		}
		length := cache.Strlen(tokens[1])

		return NewInteger(length), nil
	case "GETSET":
		if len(tokens) != 3 {
			return Response{}, fmt.Errorf("Error: GETSET requires exactly two arguments")
		}
		oldValue, exists := cache.GetSet(tokens[1], tokens[2])

		if !exists {
			return NewNullBulkString(), nil
		}
		return NewBulkString(oldValue), nil
	case "COMMAND":
		return NewSimpleString("OK"), nil
	}

	return Response{}, fmt.Errorf("Error: Unknown command")
}
