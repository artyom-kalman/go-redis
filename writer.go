package main

import (
	"io"
)

type Writer struct {
	writer io.Writer
}

func NewWriter(wt io.Writer) *Writer {
	return &Writer{
		writer: wt,
	}
}

func (w *Writer) Write(v Value) error {
	bytes := v.Marshal()

	_, err := w.writer.Write(bytes)
	if err != nil {
		return err
	}

	return nil
}
