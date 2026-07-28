package main

import (
	"fmt"
	"strconv"
)

func executeEXPIRE(tokens []string, cache *Cache) (Response, error) {
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
}
func executeTTL(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("Error: TTL command requires a key")
	}
	key := tokens[1]
	ttl := cache.TTL(key)
	return NewInteger(ttl), nil
}
func executePERSIST(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("Error: PERSIST command requires a key")
	}
	key := tokens[1]
	if cache.Persist(key) {
		return NewInteger(1), nil
	}
	return NewInteger(0), nil
}
