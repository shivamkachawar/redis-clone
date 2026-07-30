package protocol

import (
	"fmt"
	"strconv"
	"strings"
)

func ParseRESP(data []byte) ([]string, int, error) {
	pos := 0

	line, err := readLine(data, &pos)
	if err != nil {
		return nil, 0, err
	}

	if len(line) == 0 {
		return nil, 0, fmt.Errorf("Empty request")
	}
	if line[0] != '*' {
		return nil, 0, fmt.Errorf("Expected RESP array")
	}

	sizeStr := strings.TrimPrefix(line, "*")

	arraySize, err := strconv.Atoi(sizeStr)

	if err != nil {
		return nil, 0, fmt.Errorf("Invalid array size")
	}
	tokens := make([]string, 0, arraySize)
	for i := 0; i < arraySize; i++ {
		bulkHeader, err := readLine(data, &pos)
		if err != nil {
			return nil, 0, err
		}

		bulkHeader = strings.TrimSuffix(bulkHeader, "\r\n")
		if len(bulkHeader) == 0 || bulkHeader[0] != '$' {
			return nil, 0, fmt.Errorf("Expected bulk string")
		}

		lengthStr := strings.TrimPrefix(bulkHeader, "$")
		length, err := strconv.Atoi(lengthStr)

		if err != nil {
			return nil, 0, fmt.Errorf("Invalid bulk string length")
		}

		value := make([]byte, length)
		if pos+length > len(data) {
			return nil, 0, ErrIncomplete
		}
		value = data[pos : pos+length]
		pos += length

		if err != nil {
			return nil, 0, err
		}

		if pos+2 > len(data) {
			return nil, 0, fmt.Errorf("incomplete bulk string")
		}
		if data[pos] != '\r' || data[pos+1] != '\n' {
			return nil, 0, fmt.Errorf("expected CRLF")
		}
		pos += 2

		tokens = append(tokens, string(value))

	}

	return tokens, pos, nil
}

func readLine(data []byte, pos *int) (string, error) {
	start := *pos
	for i := *pos; i < len(data)-1; i++ {
		if data[i] == '\r' && data[i+1] == '\n' {
			line := string(data[start:i])
			*pos = i + 2
			return line, nil
		}
	}
	return "", ErrIncomplete
}
