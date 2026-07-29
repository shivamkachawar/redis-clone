package main

import (
	"bufio"
	"fmt"
	"go-redis/cache"
	"go-redis/commands"
	"go-redis/protocol"
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

// writeRESP writes a RESP array to the given file. -> Helper Function for Write and Rewrite
func writeRESP(file *os.File, tokens []string) error {
	_, err := fmt.Fprintf(file, "*%d\r\n", len(tokens))
	if err != nil {
		return err
	}

	for _, token := range tokens {
		_, err := fmt.Fprintf(file, "$%d\r\n%s\r\n", len(token), token)
		if err != nil {
			return err
		}
	}

	return nil
}

func (a *AOF) Write(tokens []string) error {
	err := writeRESP(a.file, tokens)
	if err != nil {
		return err
	}

	return a.file.Sync()
}
func (a *AOF) Replay(cache *cache.Cache) error {
	_, err := a.file.Seek(0, io.SeekStart)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(a.file)
	for {
		tokens, err := protocol.ParseRESP(reader)

		if err == io.EOF {
			break
		}

		if err != nil {
			return err
		}

		// executeCommand(...)
		_, err = commands.ExecuteCommand(tokens, cache)
		if err != nil {
			return err
		}
	}
	return nil
}
func (a *AOF) Rewrite(cache *cache.Cache) error {
	file, err := os.Create("appendonly.new.aof")
	if err != nil {
		return err
	}

	entries := cache.Entries()

	for key, entry := range entries {
		str, err := cache.GetString(&entry)
		if err != nil {
			continue
		}
		tokens := []string{"SET", key, str.Value}

		err = writeRESP(file, tokens)
		if err != nil {
			return err
		}
	}
	err = file.Sync()
	if err != nil {
		return err
	}
	err = file.Close()
	if err != nil {
		return err
	}
	err = a.file.Close()
	if err != nil {
		return err
	}
	err = os.Rename("appendonly.new.aof", "appendonly.aof")
	if err != nil {
		return err
	}
	a.file, err = os.OpenFile("appendonly.aof", os.O_CREATE|os.O_RDWR|os.O_APPEND, 0644)
	if err != nil {
		return err
	}
	return nil
}
