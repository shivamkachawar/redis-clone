package commands

import (
	"go-redis/cache"
	"go-redis/protocol"
)

func executeSADD(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) < 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	key := tokens[1]
	added := 0

	for _, member := range tokens[2:] {
		count, err := cache.SAdd(key, member)
		if err != nil {
			return protocol.Response{}, err
		}
		added += count
	}

	return protocol.NewInteger(added), nil
}
func executeSISMEMBER(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	if len(tokens) != 3 {
		return protocol.Response{}, protocol.ErrWrongNumberOfArguments
	}

	exists, err := cache.SIsMember(tokens[1], tokens[2])
	if err != nil {
		return protocol.Response{}, err
	}

	if exists {
		return protocol.NewInteger(1), nil
	}

	return protocol.NewInteger(0), nil
}
