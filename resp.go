package main

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func parseRESP(reader *bufio.Reader) ([]string, error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	line = strings.TrimSuffix(line, "\r\n")

	if len(line) == 0 {
		return nil, fmt.Errorf("Empty request")
	}
	if line[0] != '*' {
		return nil, fmt.Errorf("Expected RESP array")
	}

	sizeStr := strings.TrimPrefix(line, "*")

	arraySize, err := strconv.Atoi(sizeStr)

	if err != nil {
		return nil, fmt.Errorf("Invalid array size")
	}
	tokens := make([]string, 0, arraySize)
	for i := 0; i < arraySize; i++ {
		bulkHeader, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		bulkHeader = strings.TrimSuffix(bulkHeader, "\r\n")
		if len(bulkHeader) == 0 || bulkHeader[0] != '$' {
			return nil, fmt.Errorf("Expected bulk string")
		}

		lengthStr := strings.TrimPrefix(bulkHeader, "$")
		length, err := strconv.Atoi(lengthStr)

		if err != nil {
			return nil, fmt.Errorf("Invalid bulk string length")
		}

		data := make([]byte, length)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}

		_, err = reader.ReadString('\n') // Read the trailing \r\n
		if err != nil {
			return nil, err
		}

		tokens = append(tokens, string(data))

	}
	fmt.Println(tokens)

	return tokens, nil
}
