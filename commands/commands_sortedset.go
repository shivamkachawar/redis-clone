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
func executeZRANGE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	start, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidInteger
	}

	stop, err := strconv.Atoi(tokens[3])
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidInteger
	}

	members, err := cache.ZRange(tokens[1], start, stop)
	if err != nil {
		return protocol.Response{}, err
	}

	responses := make([]protocol.Response, 0, len(members))

	for _, member := range members {
		responses = append(responses, protocol.NewBulkString(member))
	}

	return protocol.NewArray(responses), nil
}
func executeZRANK(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	rank, found, err := cache.ZRank(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}

	if !found {
		return protocol.NewNullBulkString(), nil
	}

	return protocol.NewInteger(rank), nil
}
func executeZREVRANK(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	rank, found, err := cache.ZRevRank(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}

	if !found {
		return protocol.NewNullBulkString(), nil
	}

	return protocol.NewInteger(rank), nil
}
func executeZREVRANGE(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	start, err := strconv.Atoi(tokens[2])
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidInteger
	}

	stop, err := strconv.Atoi(tokens[3])
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidInteger
	}

	members, err := cache.ZRevRange(tokens[1], start, stop)
	if err != nil {
		return protocol.Response{}, err
	}

	responses := make([]protocol.Response, 0, len(members))

	for _, member := range members {
		responses = append(responses, protocol.NewBulkString(member))
	}

	return protocol.NewArray(responses), nil
}
func executeZCOUNT(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 4 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	min, err := strconv.ParseFloat(tokens[2], 64)
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidFloat
	}

	max, err := strconv.ParseFloat(tokens[3], 64)
	if err != nil {
		return protocol.Response{}, protocol.ErrInvalidFloat
	}

	count, err := cache.ZCount(tokens[1], min, max)
	if err != nil {
		return protocol.Response{}, err
	}

	return protocol.NewInteger(count), nil
}
