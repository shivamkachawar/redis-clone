package main

import (
	"fmt"
	"strconv"
	"strings"
)

func executeCommand(tokens []string, cache *Cache) (string, error) {
	command := tokens[0]
	switch command {
	case "PING":
		return "PONG", nil
	case "SET":
		if len(tokens) != 3 && len(tokens) != 5 {
			return "", fmt.Errorf("Usage: SET <key> <value> [EX <seconds>]")
		}

		key := tokens[1]
		value := tokens[2]

		// Normal SET
		if len(tokens) == 3 {
			cache.Set(key, value, 0)
			return "OK", nil
		}

		// SET with EX
		if strings.ToUpper(tokens[3]) != "EX" {
			return "", fmt.Errorf("Error: Expected EX")
		}

		seconds, err := strconv.Atoi(tokens[4])
		if err != nil {
			return "", fmt.Errorf("Error: Invalid expiry time")
		}

		if seconds <= 0 {
			return "", fmt.Errorf("Error: TTL must be greater than 0")
		}

		cache.Set(key, value, seconds)
		return "OK", nil
	case "GET":
		if len(tokens) != 2 {
			return "", fmt.Errorf("Error: GET command requires a key")
		}
		key := tokens[1]
		value, exists := cache.Get(key)
		if !exists {
			return "nil", nil
		}
		return value, nil
	case "DEL":
		if len(tokens) != 2 {
			return "", fmt.Errorf("Error: DELETE command requires a key")
		}
		key := tokens[1]
		cache.Delete(key)
		return "OK", nil
	case "EXPIRE":
		if len(tokens) != 3 {
			return "", fmt.Errorf("Error: EXPIRE command requires a key and a time in seconds")
		}
		key := tokens[1]
		seconds, err := strconv.Atoi(tokens[2])
		if err != nil {
			return "", fmt.Errorf("Error: EXPIRE command requires a valid time in seconds")
		}
		if seconds <= 0 {
			return "", fmt.Errorf("Error: Expiry time must be greater than 0")
		}
		if cache.Expire(key, seconds) {
			return "OK", nil
		}
		return "", fmt.Errorf("Error: Key not found")
	case "TTL":
		if len(tokens) != 2 {
			return "", fmt.Errorf("Error: TTL command requires a key")
		}
		key := tokens[1]
		ttl := cache.TTL(key)
		return strconv.Itoa(ttl), nil
	case "PERSIST":
		if len(tokens) != 2 {
			return "", fmt.Errorf("Error: PERSIST command requires a key")
		}
		key := tokens[1]
		if cache.Persist(key) {
			return "OK", nil
		}
		return "", fmt.Errorf("Error: Key not found")
	case "EXISTS":
		if len(tokens) != 2 {
			return "", fmt.Errorf("Error: EXISTS command requires a key")
		}
		key := tokens[1]
		if cache.Exists(key) {
			return "1", nil
		}
		return "0", nil
	}

	return "", fmt.Errorf("Error: Unknown command")
}
