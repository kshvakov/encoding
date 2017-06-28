package encoding

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
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

func Benchmark_Decode(b *testing.B) {
	var buff bytes.Buffer
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
	v := T{
		Fieldname:  "Abc",
		Fieldname2: "Bcde",
		UInt32:     256,
		Uint64:     542,
		In: In{
			V: "AAAAAAAAAAa",
		},
	}
	NewEncoder(&buff).Encode(v)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var z T
		if err := NewDecoder(bytes.NewBuffer(buff.Bytes())).Decode(&z); err != nil {
			b.Fatal(err)
		}
		if z.UInt32 != 256 || z.In.V != "AAAAAAAAAAa" {
			b.Fatal("invalid value", z)
		}
	}
}

func Benchmark_DecodeGob(b *testing.B) {
	var buff bytes.Buffer
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
	v := T{
		Fieldname:  "Abc",
		Fieldname2: "Bcde",
		UInt32:     256,
		Uint64:     542,
		In: In{
			V: "AAAAAAAAAAa",
		},
	}
	gob.NewEncoder(&buff).Encode(v)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var z T
		if err := gob.NewDecoder(bytes.NewBuffer(buff.Bytes())).Decode(&z); err != nil {
			b.Fatal(err)
		}
		if z.UInt32 != 256 || z.In.V != "AAAAAAAAAAa" {
			b.Fatal("invalid value", z)
		}
	}
}

func Benchmark_DecodeJson(b *testing.B) {
	var buff bytes.Buffer
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
	v := T{
		Fieldname:  "Abc",
		Fieldname2: "Bcde",
		UInt32:     256,
		Uint64:     542,
		In: In{
			V: "AAAAAAAAAAa",
		},
	}
	json.NewEncoder(&buff).Encode(v)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		var z T
		if err := json.NewDecoder(bytes.NewBuffer(buff.Bytes())).Decode(&z); err != nil {
			b.Fatal(err)
		}
		if z.UInt32 != 256 || z.In.V != "AAAAAAAAAAa" {
			b.Fatal("invalid value", z)
		}
	}
}
