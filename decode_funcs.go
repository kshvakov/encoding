package encoding

import (
	"reflect"
	"sort"
	"sync"
)

type decodeFunc func(decode *decode, v reflect.Value) error

var decodeFuncMap map[reflect.Kind]decodeFunc

func init() {
	decodeFuncMap = map[reflect.Kind]decodeFunc{
		reflect.Struct: decodeStruct,
		reflect.String: decodeString,
		reflect.Uint32: decodeUInt32,
		reflect.Uint64: decodeUInt64,
	}
}

func getDecodeFunc(k reflect.Kind) decodeFunc {
	if fn, ok := decodeFuncMap[k]; ok {
		return fn
	}
	return func(decode *decode, v reflect.Value) error {
		return nil
	}
}

type column struct {
	name  string
	size  int
	block []byte
}

type columns []column

func (a columns) Len() int           { return len(a) }
func (a columns) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a columns) Less(i, j int) bool { return a[i].name < a[j].name }

var decodePool = sync.Pool{
	New: func() interface{} {
		return &decode{
			block:   make([]byte, 0, 100),
			columns: make(columns, 0, 25),
		}
	},
}

func decodeStruct(d *decode, v reflect.Value) error {
	fLen, err := d.uvarint()
	if err != nil {
		return err
	}
	if cap(d.columns) < int(fLen) {
		d.columns = make(columns, fLen)
	}
	columns := d.columns[:fLen]
	for i := 0; i < int(fLen); i++ {
		name, err := d.string()
		if err != nil {
			return err
		}
		columns[i].name = name
	}

	for i := 0; i < int(fLen); i++ {
		size, err := d.uint32()
		if err != nil {
			return err
		}
		columns[i].size = int(size)
	}

	for i, column := range columns {
		block, err := d.readFixed(column.size)
		if err != nil {
			return err
		}
		columns[i].block = block
	}

	sort.Sort(columns)

	for i, field := range fields(v.Type()) {
		f := sort.Search(len(columns), func(i int) bool { return columns[i].name >= field.name })
		if f < len(columns) && columns[f].name == field.name {
			decode := decodePool.Get().(*decode)
			decode.free()
			decode.block = append(decode.block, columns[f].block...)
			field.decode(decode, v.Field(i))
			decodePool.Put(decode)
		}
	}

	return nil
}

func decodeString(d *decode, v reflect.Value) error {
	str, err := d.string()
	if err != nil {
		return err
	}
	v.SetString(str)
	return nil
}

func decodeUInt32(d *decode, v reflect.Value) error {
	value, err := d.uint32()
	if err != nil {
		return err
	}
	v.SetUint(uint64(value))
	return nil
}

func decodeUInt64(d *decode, v reflect.Value) error {
	value, err := d.uint64()
	if err != nil {
		return err
	}
	v.SetUint(value)
	return nil
}
