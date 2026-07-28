package main

import (
	"fmt"
	"net"
)

func writeSimpleString(conn net.Conn, message string) error {

	_, err := fmt.Fprintf(conn, "+%s\r\n", message)

	return err
}
