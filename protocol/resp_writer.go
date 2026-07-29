package protocol

import (
	"fmt"
	"net"
)

func writeSimpleString(conn net.Conn, message string) error {

	_, err := fmt.Fprintf(conn, "+%s\r\n", message)

	return err
}
func WriteError(conn net.Conn, message string) error {
	_, err := fmt.Fprintf(conn, "-ERR %s\r\n", message)
	return err
}
func writeInteger(conn net.Conn, value int) error {
	_, err := fmt.Fprintf(conn, ":%d\r\n", value)
	return err
}
func writeBulkString(conn net.Conn, value string) error {
	_, err := fmt.Fprintf(conn, "$%d\r\n%s\r\n", len(value), value)
	return err
}
func writeNullBulkString(conn net.Conn) error {
	_, err := fmt.Fprint(conn, "$-1\r\n")
	return err
}
func writeArray(conn net.Conn, elements []Response) error {
	_, err := fmt.Fprintf(conn, "*%d\r\n", len(elements))
	if err != nil {
		return err
	}

	for _, element := range elements {
		err := WriteResponse(conn, element)
		if err != nil {
			return err
		}
	}

	return nil
}
func WriteResponse(conn net.Conn, response Response) error {
	switch response.Type {

	case SimpleString:
		return writeSimpleString(conn, response.Value)

	case BulkString:
		return writeBulkString(conn, response.Value)

	case Integer:
		return writeInteger(conn, response.Integer)

	case NullBulkString:
		return writeNullBulkString(conn)

	case Array:
		return writeArray(conn, response.Elements)

	case Error:
		return WriteError(conn, response.Value)

	default:
		return fmt.Errorf("unknown response type")
	}

}
