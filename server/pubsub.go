package server

import (
	"fmt"
	"go-redis/protocol"
	"syscall"
)

type PubSub struct {
	Channels map[string]map[int]*Client
}

func executeSUBSCRIBE(client *Client, pubsub *PubSub, tokens []string) (protocol.Response, error) {
	if len(tokens) != 2 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'SUBSCRIBE' command")
	}

	channel := tokens[1]

	// Already subscribed
	if _, ok := client.Subscriptions[channel]; ok {
		if _, ok := client.Subscriptions[channel]; ok {
			return protocol.NewArray([]protocol.Response{
				protocol.NewBulkString("subscribe"),
				protocol.NewBulkString(channel),
				protocol.NewInteger(len(client.Subscriptions)),
			}), nil
		}
	}

	// Add channel to the client
	client.Subscriptions[channel] = struct{}{}

	// Create channel if it doesn't exist
	if _, ok := pubsub.Channels[channel]; !ok {
		pubsub.Channels[channel] = make(map[int]*Client)
	}

	// Add client to the channel
	pubsub.Channels[channel][client.FD] = client

	return protocol.NewArray([]protocol.Response{
		protocol.NewBulkString("subscribe"),
		protocol.NewBulkString(channel),
		protocol.NewInteger(len(client.Subscriptions)),
	}), nil
}
func executePUBLISH(pubSub *PubSub, tokens []string) (protocol.Response, error) {

	if len(tokens) != 3 {
		return protocol.Response{}, fmt.Errorf("ERR wrong number of arguments for 'PUBLISH' command")
	}
	channel := tokens[1]
	message := tokens[2]

	subscribers, ok := pubSub.Channels[channel]

	if !ok {
		return protocol.NewInteger(0), nil
	}
	count := 0
	for _, subscriber := range subscribers {
		response := protocol.NewArray([]protocol.Response{
			protocol.NewBulkString("message"),
			protocol.NewBulkString(channel),
			protocol.NewBulkString(message),
		})

		resp, err := protocol.EncodeResponse(response)
		if err != nil {
			return protocol.Response{}, err
		}

		_, err = syscall.Write(subscriber.FD, resp)
		if err != nil {
			continue
		}

		count++
	}
	return protocol.NewInteger(count), nil
}

func removeClientFromPubSub(client *Client, pubSub *PubSub) {
	for channel := range client.Subscriptions {

		delete(pubSub.Channels[channel], client.FD)

		if len(pubSub.Channels[channel]) == 0 {
			delete(pubSub.Channels, channel)
		}
	}
}
