package testutil

import "bytes"

type TestBuffer struct {
	ReadableBytes *bytes.Buffer
	BytesWritten  *bytes.Buffer
}

func (t *TestBuffer) Read(p []byte) (int, error) {
	return t.ReadableBytes.Read(p)
}
func (t *TestBuffer) ReadByte() (byte, error) {
	return t.ReadableBytes.ReadByte()
}
func (t *TestBuffer) Write(buf []byte) (int, error) {
	return t.BytesWritten.Write(buf)
}
