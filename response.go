package main

type ResponseType int

const (
	SimpleString ResponseType = iota
	BulkString
	Integer
	NullBulkString
	Error
)

type Response struct {
	Type    ResponseType
	Value   string
	Integer int
}

func NewSimpleString(value string) Response {
	return Response{
		Type:  SimpleString,
		Value: value,
	}
}

func NewBulkString(value string) Response {
	return Response{
		Type:  BulkString,
		Value: value,
	}
}

func NewInteger(value int) Response {
	return Response{
		Type:    Integer,
		Integer: value,
	}
}

func NewNullBulkString() Response {
	return Response{
		Type: NullBulkString,
	}
}
