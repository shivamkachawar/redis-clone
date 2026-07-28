package main

import "fmt"

func executeEXISTS(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) < 2 {
		return Response{}, fmt.Errorf("Error: EXISTS command requires a key")
	}
	count := cache.Exists(tokens[1:])
	return NewInteger(count), nil
}
func executeDEL(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) < 2 {
		return Response{}, fmt.Errorf("Error: DELETE command requires a key")
	}
	count := cache.Delete(tokens[1:])
	return NewInteger(count), nil
}
func executeKEYS(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 2 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'KEYS' command")
	}
	if tokens[1] != "*" {
		return Response{}, fmt.Errorf("ERR only KEYS * is currently supported")
	}

	keys := cache.Keys()
	responses := make([]Response, 0, len(keys))
	for _, key := range keys {
		responses = append(responses, NewBulkString(key))
	}
	return NewArray(responses), nil
}
func executeDBSIZE(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 1 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'DBSIZE' command")
	}
	size := cache.Size()
	return NewInteger(size), nil
}
func executeFLUSHDB(tokens []string, cache *Cache) (Response, error) {
	if len(tokens) != 1 {
		return Response{}, fmt.Errorf("ERR wrong number of arguments for 'FLUSHDB' command")
	}
	cache.Flush()
	return NewSimpleString("OK"), nil
}
