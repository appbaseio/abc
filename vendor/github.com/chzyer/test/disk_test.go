package test

import (
	"io"
	"testing"
)

func TestMemDisk(t *testing.T) {
	defer New(t)
	md := NewMemDisk()
	WriteAt(md, []byte("hello"), 4)
	h := make([]byte, 5)
	n, err := md.ReadAt(h, 3)
	Equals(n, 5, err, nil, h, []byte("\x00hell"))

	{
		n, err := md.ReadAt(h, 10)
		Equals(n, 0, err, io.EOF)
	}
}
