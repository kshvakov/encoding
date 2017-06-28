package encoding

import "reflect"

type encodeFunc func(enc *encode, v reflect.Value) error

var encodeFuncMap map[reflect.Kind]encodeFunc

func init() {
	encodeFuncMap = map[reflect.Kind]encodeFunc{
		reflect.Struct: encodeStruct,
		reflect.String: encodeString,
		reflect.Uint32: encodeUInt32,
		reflect.Uint64: encodeUInt64,
	}
}

func encodeStruct(enc *encode, v reflect.Value) error {
	fields := fields(v.Type())
	if err := enc.uvarint(uint64(len(fields))); err != nil {
		return err
	}
	for _, field := range fields {
		if err := enc.string(field.name); err != nil {
			return err
		}
	}
	offsets := enc.buf.alloc(4 * len(fields))
	for i, field := range fields {
		startOffset := enc.buf.len()
		if err := field.encode(enc, v.Field(i)); err != nil {
			return err
		}
		var (
			idx  = 4 * i
			bLen = int32(enc.buf.len() - startOffset)
		)
		{
			offsets[idx+0] = byte(bLen)
			offsets[idx+1] = byte(bLen >> 8)
			offsets[idx+2] = byte(bLen >> 16)
			offsets[idx+3] = byte(bLen >> 24)
		}
	}
	return nil
}

func encodeUInt32(enc *encode, v reflect.Value) error {
	return enc.uint32(uint32(v.Uint()))
}

func encodeUInt64(enc *encode, v reflect.Value) error {
	return enc.uint64(v.Uint())
}

func encodeString(enc *encode, v reflect.Value) error {
	return enc.string(v.String())
}

func getEncodeFunc(k reflect.Kind) encodeFunc {
	if fn, ok := encodeFuncMap[k]; ok {
		return fn
	}
	return func(enc *encode, v reflect.Value) error {
		return nil
	}
}
