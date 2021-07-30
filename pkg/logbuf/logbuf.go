package logbuf

import "bytes"

type Buffer struct {
	buf []byte
	idx int
}

func New(size int) *Buffer {
	return &Buffer{
		buf: make([]byte, size),
		idx: 0,
	}
}

func (b *Buffer) Write(p []byte) (int, error) {
	for i := 0; i < len(p); i++ {
		b.buf[b.idx] = p[i]
		b.idx++
		if b.idx >= len(b.buf) {
			b.idx = 0
		}
	}

	return len(p), nil
}

func (b *Buffer) Bytes() []byte {
	var buf []byte
	buf = append(buf, b.buf[b.idx:]...)
	buf = append(buf, b.buf[0:b.idx]...)
	return bytes.Trim(buf, "\x00")
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}
