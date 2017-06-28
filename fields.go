package encoding

import (
	"reflect"
	"sync"
)

var fieldsCache struct {
	mutex  sync.RWMutex
	fields map[reflect.Type][]field
}

func init() {
	fieldsCache.fields = make(map[reflect.Type][]field, 0)
}

type field struct {
	name   string
	encode encodeFunc
	decode decodeFunc
}

func fields(v reflect.Type) []field {
	fieldsCache.mutex.RLock()
	_, ok := fieldsCache.fields[v]
	fieldsCache.mutex.RUnlock()
	if !ok {
		var (
			numField = v.NumField()
			fields   = make([]field, 0, numField)
		)
		for i := 0; i < numField; i++ {
			f := v.Field(i)
			name := f.Name
			if n := f.Tag.Get("encoder"); len(n) != 0 {
				name = n
			}
			if f.Anonymous || f.PkgPath != "" || name == "-" {
				continue
			}
			fields = append(fields, field{
				name:   name,
				encode: getEncodeFunc(f.Type.Kind()),
				decode: getDecodeFunc(f.Type.Kind()),
			})
		}
		fieldsCache.mutex.Lock()
		fieldsCache.fields[v] = fields
		fieldsCache.mutex.Unlock()
	}
	return fieldsCache.fields[v]
}
