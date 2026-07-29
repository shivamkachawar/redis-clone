package commands

import (
	"fmt"
	"go-redis/cache"
	"go-redis/protocol"
)

func ExecuteCommand(tokens []string, cache *cache.Cache) (protocol.Response, error) {
	command := tokens[0]
	switch command {
	case "PING":
		return protocol.NewSimpleString("PONG"), nil
	case "SET":
		return executeSET(tokens, cache)
	case "GET":
		return executeGET(tokens, cache)
	case "DEL":
		return executeDEL(tokens, cache)
	case "EXPIRE":
		return executeEXPIRE(tokens, cache)
	case "TTL":
		return executeTTL(tokens, cache)
	case "PERSIST":
		return executePERSIST(tokens, cache)
	case "EXISTS":
		return executeEXISTS(tokens, cache)
	case "INCR":
		return executeINCR(tokens, cache)
	case "DECR":
		return executeDECR(tokens, cache)
	case "INCRBY":
		return executeINCRBY(tokens, cache)
	case "DECRBY":
		return executeDECRBY(tokens, cache)
	case "MGET":
		return executeMGET(tokens, cache)
	case "MSET":
		return executeMSET(tokens, cache)
	case "APPEND":
		return executeAPPEND(tokens, cache)
	case "STRLEN":
		return executeSTRLEN(tokens, cache)
	case "GETSET":
		return executeGETSET(tokens, cache)
	case "COMMAND":
		return protocol.NewSimpleString("OK"), nil
	case "KEYS":
		return executeKEYS(tokens, cache)
	case "DBSIZE":
		return executeDBSIZE(tokens, cache)
	case "FLUSHDB":
		return executeFLUSHDB(tokens, cache)
	case "LPUSH":
		return executeLPUSH(tokens, cache)
	case "LRANGE":
		return executeLRANGE(tokens, cache)
	case "LLEN":
		return executeLLEN(tokens, cache)
	case "LPOP":
		return executeLPOP(tokens, cache)
	case "RPUSH":
		return executeRPUSH(tokens, cache)
	case "RPOP":
		return executeRPOP(tokens, cache)
	case "HSET":
		return executeHSET(tokens, cache)
	case "HGET":
		return executeHGET(tokens, cache)
	default:
		return protocol.Response{}, fmt.Errorf("Error: Unknown command")
	}
}
