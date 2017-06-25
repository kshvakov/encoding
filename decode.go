package encoding

import (
	"encoding/binary"
	"io"
)

type decode struct {
	in      io.Reader
	scratch []byte
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

func (decode *decode) readFixed(l int) ([]byte, error) {
	if len(decode.scratch) < l {
		decode.scratch = make([]byte, 0, l)
	}
	if _, err := io.ReadFull(decode.in, decode.scratch[:l]); err != nil {
		return nil, err
	}
	return decode.scratch[:l], nil
}

func (decode *decode) ReadByte() (byte, error) {
	if _, err := io.ReadFull(decode.in, decode.scratch[:1]); err != nil {
		return 0x0, err
	}
	return decode.scratch[0], nil
}
