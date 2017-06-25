package encoding

import (
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
	"time"
)

func Test_Fields(t *testing.T) {
	var (
		val struct {
			Int     int
			Float32 float32
			String  *string
			Time    time.Time `encoder:"time"`
		}
		assets = []struct {
			name string
		}{
			{
				name: "Int",
			},
			{
				name: "Float32",
			},
			{
				name: "String",
			},
			{
				name: "time",
			},
		}
	)
	if fields := fields(reflect.TypeOf(val)); assert.Len(t, fields, 4) {
		for i, asset := range assets {
			assert.Equal(t, asset.name, fields[i].name)
		}
	}
}

func Benchmark_Fields(b *testing.B) {
	var v struct {
		A int
		B string
		C string `encoder:"ccc"`
	}
	r := reflect.TypeOf(v)
	b.ResetTimer()
	b.ReportAllocs()
	for i := 0; i < b.N; i++ {
		if len(fields(r)) != 3 {
			b.Fatal("wrong result")
		}
	}
}
