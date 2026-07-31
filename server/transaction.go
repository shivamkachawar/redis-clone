package server

import (
	"errors"
	"go-redis/aof"
	"go-redis/cache"
	"go-redis/commands"
	"go-redis/protocol"
)

func executeTransaction(client *Client, cache *cache.Cache, aofFile *aof.AOF) (protocol.Response, error) {

	if !client.InTransaction {
		return protocol.Response{}, errors.New("ERR EXEC without MULTI")
	}

	responses := make([]protocol.Response, 0, len(client.TransactionQueue))

	for _, tokens := range client.TransactionQueue {
		response, err := commands.ExecuteCommand(tokens, cache)
		if err != nil {
			responses = append(responses, protocol.NewError(err.Error()))
			continue
		}
		if len(tokens) > 0 && isWriteCommand(tokens[0]) {
			if err := aofFile.Write(tokens); err != nil {
				return protocol.Response{}, err
			}
		}
		responses = append(responses, response)
	}
	client.InTransaction = false
	client.TransactionQueue = nil

	return protocol.NewArray(responses), nil
}
