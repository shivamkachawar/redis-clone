package commands

import (
	"go-redis/cache"
	"go-redis/protocol"
	"strconv"
)

func executeZADD(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	score, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidFloat
	}

	if err := cache.ZAdd(tokens[1], score, tokens[3]); err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewSimpleString("OK"), nil
}
func executeZREM(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	deleted, err := cache.ZRem(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}

	if deleted {
		return protocol.NewInteger(1), nil
	}

	return protocol.NewInteger(0), nil
}
func executeZSCORE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	score, exists, err := cache.ZScore(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}

	if !exists {
		return protocol.NewNullBulkString(), nil
	}

	return protocol.NewBulkString(strconv.FormatFloat(score, 'f', -1, 64)), nil
}
func executeZCARD(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	count, err := cache.ZCard(tokens[1])
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewInteger(int(count)), nil
}
