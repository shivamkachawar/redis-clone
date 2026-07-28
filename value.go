package main

type RedisValue interface {
	isRedisValue()
}
type StringValue struct {
	Value string
}

func (*StringValue) isRedisValue() {}

type ListValue struct {
	Values []string
}

func (*ListValue) isRedisValue() {}
