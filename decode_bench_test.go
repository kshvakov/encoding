package encoding

import (
	"bytes"
	"encoding/binary"
	"testing"
)

type benchReader struct {
	buf []byte
}

func (b *benchReader) Read(p []byte) (int, error) {
	copy(p, b.buf[:len(p)])
	return len(p), nil
}

type benchStringReader struct {
	buf []byte
}

func (b *benchStringReader) Read(p []byte) (int, error) {
	switch {
	case len(p) == 1:
		copy(p, b.buf[:1])
	default:
		copy(p, b.buf[1:len(p)+1])
	}
	return len(p), nil
}

func Benchmark_DecodeUvarint(b *testing.B) {
	var buf bytes.Buffer
	encode := encode{
		buf: &buf,
	}
	encode.uvarint(42)
	decode := decode{
		in: &benchReader{
			buf: buf.Bytes(),
		},
		scratch: make([]byte, binary.MaxVarintLen64),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, err := decode.uvarint()
		switch {
		case err != nil:
			b.Fatal(err)
		case v != 42:
			b.Fatal("incorrect result")
		}
	}
}

func Benchmark_DecodeString(b *testing.B) {
	var (
		buf    bytes.Buffer
		encode = encode{
			buf: &buf,
		}
	)
	encode.string("string")
	decode := decode{
		in: &benchStringReader{
			buf: buf.Bytes(),
		},
		scratch: make([]byte, binary.MaxVarintLen64),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		v, err := decode.string()
		switch {
		case err != nil:
			b.Fatal(err)
		case v != "string":
			b.Fatalf("incorrect result: %v", v)
		}
	}
}

func Benchmark_DecodeUint8(b *testing.B) {
	decode := decode{
		in: &benchReader{
			buf: []byte{42},
		},
		scratch: make([]byte, binary.MaxVarintLen64),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if _, err := decode.uint8(); err != nil {
			b.Fatal(err)
		}
	}
}
