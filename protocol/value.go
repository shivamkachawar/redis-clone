package protocol

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

type HashValue struct {
	Fields map[string]string
}

func (*HashValue) isRedisValue() {}
