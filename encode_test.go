package encoding

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_UInt8(t *testing.T) {
	var (
		buf    bytes.Buffer
		encode = encode{
			buf: &buf,
		}
		decode = decode{
			in:      &buf,
			scratch: make([]byte, binary.MaxVarintLen64),
		}
	)
	var prev uint8
	for i := 1; i < 255; i++ {
		if err := encode.uint8(uint8(i)); assert.NoError(t, err) {
			if v, err := decode.uint8(); assert.NoError(t, err) {
				if assert.Equal(t, uint8(i), v) {
					if assert.NotEqual(t, prev, v) {
						prev = v
					}
				}
			}
		}
	}
}

func Test_String(t *testing.T) {
	var (
		buf    bytes.Buffer
		encode = encode{
			buf: &buf,
		}
		decode = decode{
			in:      &buf,
			scratch: make([]byte, binary.MaxVarintLen64),
		}
	)
	var prev string
	for i := 0; i < 255; i++ {
		if err := encode.string(fmt.Sprintf("str_%d", i)); assert.NoError(t, err) {
			if v, err := decode.string(); assert.NoError(t, err) {
				if assert.Equal(t, fmt.Sprintf("str_%d", i), v) {
					if assert.NotEqual(t, prev, v) {
						prev = v
					}
				}
			}
		}
	}
}
