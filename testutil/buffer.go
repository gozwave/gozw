package testutil

import (
	"bytes"
	"io"
)

type TestBuffer struct {
	ReadableBytes io.Reader
	BytesWritten  *bytes.Buffer
}

func (t *TestBuffer) Read(p []byte) (int, error) {
	return t.ReadableBytes.Read(p)
}

func (t *TestBuffer) ReadByte() (byte, error) {
	var buf [1]byte
	_, err := t.ReadableBytes.Read(buf[:])
	return buf[0], err
}
func (t *TestBuffer) Write(buf []byte) (int, error) {
	return t.BytesWritten.Write(buf)
}
