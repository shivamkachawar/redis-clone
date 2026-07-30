package server

import "strings"

func isWriteCommand(command string) bool {
	switch strings.ToUpper(command) {
	case "SET", "MSET", "DEL", "APPEND", "GETSET",
		"INCR", "DECR", "INCRBY", "DECRBY",
		"EXPIRE", "PERSIST", "FLUSHDB",
		"LPUSH",
		"HSET", "HINCRBY", "HSETNX",
		"SADD":
		return true
	default:
		return false
	}
}
