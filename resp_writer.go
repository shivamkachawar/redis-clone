package main

import (
	"fmt"
	"net"
)

func writeSimpleString(conn net.Conn, message string) error {

	_, err := fmt.Fprintf(conn, "+%s\r\n", message)

	return err
}
func writeError(conn net.Conn, message string) error {
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
