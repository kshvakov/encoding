package encoding

import (
	"encoding/binary"
	"io"
	"reflect"
	"unsafe"
)

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{
		out: w,
	}
}

type Encoder struct {
	out io.Writer
}

func (e *Encoder) Encode(v interface{}) error {
	encode := &encode{}
	if err := encode.serialize(v); err != nil {
		return err
	}
	if _, err := encode.buf.WriteTo(e.out); err != nil {
		return err
	}
	return nil
}

type encode struct {
	buf interface {
		io.Writer
		io.WriterTo
	}
	scratch [binary.MaxVarintLen64]byte
}

func (enc *encode) serialize(v interface{}) error {
	return nil
}

func (enc *encode) bool(v bool) error {
	if v {
		return enc.uint8(1)
	}
	return enc.uint8(0)
}

func (enc *encode) uint8(v uint8) error {
	enc.scratch[0] = v
	if _, err := enc.buf.Write(enc.scratch[:1]); err != nil {
		return err
	}
	return nil
}

func (enc *encode) uint16(v uint16) error {
	enc.scratch[0] = byte(v)
	enc.scratch[1] = byte(v >> 8)
	if _, err := enc.buf.Write(enc.scratch[:2]); err != nil {
		return err
	}
	return nil
}

func (enc *encode) uint32(v uint32) error {
	enc.scratch[0] = byte(v)
	enc.scratch[1] = byte(v >> 8)
	enc.scratch[2] = byte(v >> 16)
	enc.scratch[3] = byte(v >> 24)
	if _, err := enc.buf.Write(enc.scratch[:4]); err != nil {
		return err
	}
	return nil
}

func (enc *encode) uint64(v uint64) error {
	enc.scratch[0] = byte(v)
	enc.scratch[1] = byte(v >> 8)
	enc.scratch[2] = byte(v >> 16)
	enc.scratch[3] = byte(v >> 24)
	enc.scratch[4] = byte(v >> 48)
	enc.scratch[5] = byte(v >> 56)
	if _, err := enc.buf.Write(enc.scratch[:6]); err != nil {
		return err
	}
	return nil
}

func (enc *encode) int64(v int64) error {
	return enc.uint64(uint64(v))
}

func (enc *encode) uvarint(v uint64) error {
	len := binary.PutUvarint(enc.scratch[:binary.MaxVarintLen64], v)
	if _, err := enc.buf.Write(enc.scratch[0:len]); err != nil {
		return err
	}
	return nil
}

func (enc *encode) string(v string) error {
	str := str2bytes(v)
	if err := enc.uvarint(uint64(len(str))); err != nil {
		return err
	}
	if _, err := enc.buf.Write(str); err != nil {
		return err
	}
	return nil
}

func str2bytes(str string) []byte {
	header := (*reflect.SliceHeader)(unsafe.Pointer(&str))
	header.Len = len(str)
	header.Cap = header.Len
	return *(*[]byte)(unsafe.Pointer(header))
}
