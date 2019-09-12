package testutil

import "bytes"

// TestBuffer contains a test buffer.
type TestBuffer struct {
	ReadableBytes *bytes.Buffer
	BytesWritten  *bytes.Buffer
}

// Read implements io.Reader.
func (t *TestBuffer) Read(p []byte) (int, error) {
	return t.ReadableBytes.Read(p)
}

// ReadByte implements io.ByteReader.
func (t *TestBuffer) ReadByte() (byte, error) {
	return t.ReadableBytes.ReadByte()
}

// Write implements io.Writer.
func (t *TestBuffer) Write(buf []byte) (int, error) {
	return t.BytesWritten.Write(buf)
}
