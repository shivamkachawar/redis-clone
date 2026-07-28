package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
)

type AOF struct {
	file *os.File
}

func NewAOF(filename string) (*AOF, error) {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return nil, err
	}

	return &AOF{
		file: file,
	}, nil
}

func (a *AOF) Write(tokens []string) error {
	_, err := fmt.Fprintf(a.file, "*%d\r\n", len(tokens))
	if err != nil {
		return err
	}

	for _, token := range tokens {
		_, err := fmt.Fprintf(a.file, "$%d\r\n%s\r\n", len(token), token)
		if err != nil {
			return err
		}
	}
	err = a.file.Sync()
	if err != nil {
		return err
	}

	return nil
}
func (a *AOF) Replay(cache *Cache) error {
	_, err := a.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(a.file)
	for {
		tokens, err := parseRESP(reader)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// executeCommand(...)
		_, err = executeCommand(tokens, cache)
		if err != nil {
			return err
		}
	}
	return nil
}
