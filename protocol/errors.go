package protocol

import "errors"

var ErrWrongType = errors.New("WRONGTYPE Operation against a key holding the wrong kind of value")
var ErrHashValueNotInteger = errors.New("ERR hash value is not an integer")
var ErrWrongNumberOfArguments = errors.New("ERR command has wrong number of arguments")
var ErrIncomplete = errors.New("Incomplete RESP")
var ErrInvalidFloat = errors.New("Invalid float value")
