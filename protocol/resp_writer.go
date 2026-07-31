package protocol

import (
	"fmt"
)

func writeSimpleString(message string) []byte {
	return []byte(fmt.Sprintf("+%s\r\n", message))
}
func writeError(message string) []byte {
	return []byte(fmt.Sprintf("-ERR %s\r\n", message))
}
func writeInteger(value int) []byte {
	return []byte(fmt.Sprintf("+%d\r\n", value))
}
func writeBulkString(value string) []byte {
	return []byte(fmt.Sprintf("$%d\r\n%s\r\n", len(value), value))
}
func writeNullBulkString() []byte {
	return []byte("$-1\r\n")
}

func writeArray(elements []Response) ([]byte, error) {
	data := []byte(fmt.Sprintf("*%d\r\n", len(elements)))

	for _, e := range elements {
		b, err := EncodeResponse(e)
		if err != nil {
			return nil, err
		}
		data = append(data, b...)
	}

	return data, nil
}
func EncodeResponse(response Response) ([]byte, error) {
	switch response.Type {

	case SimpleString:
		return writeSimpleString(response.Value), nil

	case BulkString:
		return writeBulkString(response.Value), nil

	case Integer:
		return writeInteger(response.Integer), nil

	case NullBulkString:
		return writeNullBulkString(), nil

	case Array:
		return writeArray(response.Elements)

	case Error:
		return writeError(response.Value), nil

	default:
		return nil, fmt.Errorf("unknown response type")
	}

}
func EncodeError(message string) []byte {
	return []byte(fmt.Sprintf("-ERR %s\r\n", message))
}
