package encoding

import (
	"encoding/binary"
	"io"
	"reflect"
)

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		input: r,
	}
}

type Decoder struct {
	input   io.Reader
	scratch [4]byte
}

func (d *Decoder) Decode(out interface{}) error {
	if _, err := io.ReadFull(d.input, d.scratch[:]); err != nil {
		return err
	}
	ln := uint32(d.scratch[0]) | uint32(d.scratch[1])<<8 | uint32(d.scratch[2])<<16 | uint32(d.scratch[3])<<24
	decode := decodePool.Get().(*decode)
	decode.free()
	if cap(decode.block) < int(ln) {
		decode.block = make([]byte, 0, ln)
	}
	decode.block = decode.block[:ln]
	if _, err := io.ReadFull(d.input, decode.block); err != nil {
		return err
	}
	decodeStruct(decode, reflect.ValueOf(out).Elem())
	decodePool.Put(decode)
	return nil
}

type decode struct {
	block   []byte
	offset  int
	columns columns
}

func (decode *decode) free() {
	decode.block = decode.block[0:0]
	decode.offset = 0
	decode.columns = decode.columns[0:0]
}

func (decode *decode) uvarint() (uint64, error) {
	return binary.ReadUvarint(decode)
}

func (decode *decode) uint8() (uint8, error) {
	byte, err := decode.ReadByte()
	if err != nil {
		return 0, err
	}
	return uint8(byte), nil
}

func (decode *decode) uint32() (uint32, error) {
	b, err := decode.readFixed(4)
	if err != nil {
		return 0, err
	}
	return uint32(b[0]) | uint32(b[1])<<8 | uint32(b[2])<<16 | uint32(b[3])<<24, nil
}

func (decode *decode) uint64() (uint64, error) {
	b, err := decode.readFixed(6)
	if err != nil {
		return 0, err
	}
	return uint64(b[0]) | uint64(b[1])<<8 | uint64(b[2])<<16 | uint64(b[3])<<24 | uint64(b[4])<<48 | uint64(b[5])<<56, nil
}

func (decode *decode) string() (string, error) {
	strlen, err := decode.uvarint()
	if err != nil {
		return "", err
	}
	bytes, err := decode.readFixed(int(strlen))
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func (decode *decode) readFixed(ln int) ([]byte, error) {
	idx := decode.offset
	decode.offset = idx + ln
	return decode.block[idx : idx+ln], nil
}

func (decode *decode) ReadByte() (byte, error) {
	idx := decode.offset
	decode.offset++
	return decode.block[idx], nil
}
