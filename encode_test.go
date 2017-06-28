package encoding

import (
	"bytes"
	"testing"
)

func Test_Encode(t *testing.T) {
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
	{
		NewEncoder(&buff).Encode(v)
	}
	var z T
	NewDecoder(&buff).Decode(&z)
	t.Logf("%#v", z)
	t.Log(z.UInt32, z.Uint64)
	/*


		ln := buff.Len()
		buff.Reset()
		{
			gob.NewEncoder(&buff).Encode(v)
		}
		t.Logf("[gob ] reduce=%.4f%%", 100-(float64(ln)/float64(buff.Len())*100))
		buff.Reset()
		{
			json.NewEncoder(&buff).Encode(v)
		}
		t.Logf("[json] reduce=%.4f%%", 100-(float64(ln)/float64(buff.Len())*100))
	*/
}

/*
func Test_UInt8(t *testing.T) {
	var prev uint8
	for i := 1; i < 2; i++ {
		var (
			buf    = newBuffer(10)
			encode = encode{
				buf: buf,
			}
		)
		if err := encode.uint8(uint8(i)); assert.NoError(t, err) {
			decode := decode{
				in:      bytes.NewBuffer(buf.bytes()),
				scratch: make([]byte, binary.MaxVarintLen64),
			}
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
	var prev string
	for i := 0; i < 255; i++ {
		var (
			buf    = newBuffer(10)
			encode = encode{
				buf: buf,
			}
		)
		if err := encode.string(fmt.Sprintf("str_%d", i)); assert.NoError(t, err) {
			decode := decode{
				in:      bytes.NewBuffer(buf.bytes()),
				scratch: make([]byte, binary.MaxVarintLen64),
			}
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
*/
