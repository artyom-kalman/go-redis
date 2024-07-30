package main

import (
	"bufio"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

type Aof struct {
	file  *os.File
	rd    *bufio.Reader
	mutex sync.Mutex
}

func NewAof(path string) (*Aof, error) {
	file, err := os.OpenFile(path, os.O_RDWR|os.O_CREATE, 0666)
	if err != nil {
		return nil, err
	}

	aof := &Aof{
		file: file,
		rd:   bufio.NewReader(file),
	}

	go func() {
		for {
			aof.mutex.Lock()
			aof.file.Sync()
			aof.mutex.Unlock()

			time.Sleep(time.Second)
		}
	}()

	return aof, nil
}

func (aof *Aof) Write(value Value) error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	_, err := aof.file.Write(value.Marshal())
	if err != nil {
		return err
	}

	return nil
}

func (aof *Aof) Read() error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	aof.file.Seek(0, io.SeekStart)

	resp := NewResp(aof.file)

	for {
		value, err := resp.Read()
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		command := strings.ToUpper(value.array[0].bulk)
		args := value.array[1:]

		handler := Handlers[command]
		handler(args)

	}

	return nil
}

func (aof *Aof) Close() error {
	aof.mutex.Lock()
	defer aof.mutex.Unlock()

	return aof.file.Close()
}
