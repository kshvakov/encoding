package encoding

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"
	"time"
)

func Benchmark_Test(b *testing.B) {
	benchmarkEncoderFn(b, NewEncoder(ioutil.Discard))
}

func Benchmark_TestJson(b *testing.B) {
	benchmarkEncoderFn(b, json.NewEncoder(ioutil.Discard))
}

func Benchmark_TestGob(b *testing.B) {
	benchmarkEncoderFn(b, gob.NewEncoder(ioutil.Discard))
}

type benchmarkEncoder interface {
	Encode(interface{}) error
}

func benchmarkEncoderFn(b *testing.B, encoder benchmarkEncoder) {
	type (
		In struct {
			V string
		}
		T struct {
			Fieldname  string
			Fieldname2 string
			UInt32     uint32
			Uint64     uint64
			In         In
		}
	)
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		encoder.Encode(T{
			Fieldname:  "A",
			Fieldname2: "B",
			UInt32:     256,
			Uint64:     542,
			In: In{
				V: "AAAAAAAAAAa",
			},
		})
	}
}

func Benchmark_EncodeBool(b *testing.B) {
	enc := encode{
		buf: newBuffer(1),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.bool((i % 10) == 0); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeUInt8(b *testing.B) {
	enc := encode{
		buf: newBuffer(1),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.uint8(uint8(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeUInt16(b *testing.B) {
	enc := encode{
		buf: newBuffer(2),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.uint16(uint16(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeUInt32(b *testing.B) {
	enc := encode{
		buf: newBuffer(4),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.uint32(uint32(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeUInt64(b *testing.B) {
	enc := encode{
		buf: newBuffer(8),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.uint64(uint64(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeInt64(b *testing.B) {
	enc := encode{
		buf: newBuffer(8),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.int64(int64(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeUvarint(b *testing.B) {
	enc := encode{
		buf: newBuffer(10),
	}
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.uvarint(uint64(i % 255)); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}

func Benchmark_EncodeString(b *testing.B) {
	var (
		enc = encode{
			buf: newBuffer(100),
		}
		str = fmt.Sprintf("abc_%d", time.Now().Unix())
	)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if err := enc.string(str); err != nil {
			b.Fatal(err)
		}
		enc.buf.free()
	}
}
