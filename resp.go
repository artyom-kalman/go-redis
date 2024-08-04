package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

const (
	STRING  = '+'
	ERROR   = '-'
	INTEGER = ':'
	BULK    = '$'
	ARRAY   = '*'
)

// Resp represents a RESP reader
type Resp struct {
	reader *bufio.Reader
}

func NewResp(rd io.Reader) *Resp {
	return &Resp{
		reader: bufio.NewReader(rd),
	}
}

func (r *Resp) readLine() (line []byte, n int, err error) {
	for {
		b, err := r.reader.ReadByte()
		if err != nil {
			return nil, 0, err
		}
		n += 1

		line = append(line, b)

		if len(line) >= 2 && line[len(line)-2] == '\r' {
			break
		}
	}
	return line[:len(line)-2], n, nil
}

func (r *Resp) readInt() (x int, n int, err error) {
	line, n, err := r.readLine()
	if err != nil {
		return 0, 0, err
	}

	i64, err := strconv.ParseInt(string(line), 10, 64)
	if err != nil {
		return 0, 0, err
	}

	return int(i64), n, nil
}

func (r *Resp) readArray() (value Value, err error) {
	value.dataType = "array"

	arrayLen, _, err := r.readInt()
	if err != nil {
		return value, err
	}

	value.array = make([]Value, arrayLen)

	for i := 0; i < arrayLen; i++ {
		val, err := r.Read()
		if err != nil {
			return value, err
		}

		value.array[i] = val
	}
	return value, nil
}

func (r *Resp) readBulk() (value Value, err error) {
	value.dataType = "bulk"

	len, _, err := r.readInt()
	if err != nil {
		return value, err
	}

	bulk := make([]byte, len)
	r.reader.Read(bulk)

	value.bulk = string(bulk)

	r.readLine()

	return value, nil
}

// Read reads data from a buffer
func (r *Resp) Read() (Value, error) {
	_type, err := r.reader.ReadByte()
	if err != nil {
		return Value{}, err
	}

	switch _type {
	case ARRAY:
		return r.readArray()
	case BULK:
		return r.readBulk()
	default:
		fmt.Println("Unknown type: ", _type)
		return Value{}, errors.New("Invalid type")
	}
}
