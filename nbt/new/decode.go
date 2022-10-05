package new

import (
	"fmt"
	"reflect"
	"sync"
)

type Decoder struct {
	p *parser
}

func Unmarshal(data []byte, v any) error {
	//todo: validate data

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("nbt: interface passed must be a pointer")
	}

	d := newDecoder()
	defer freeDecoder(d)
	d.p = &parser{buf: data, tags: make([]tag, 0, 16)}
	return d.unmarshal(rv.Elem())
}

func (d *Decoder) unmarshal(v reflect.Value) error {
	d.p.parse(tag{})
	fmt.Println(d.p.tags)
	switch v.Kind() {

	case reflect.Struct:
		t := d.p.getTag(tagCompound)
		fmt.Printf("%#v %#v\n", v.Type().Name(), t)

		rt := v.Type()
		for i := 0; i < v.NumField(); i++ {
			f := rt.Field(i)
			if !f.IsExported() {
				continue
			}

			if !(v.Field(i).Kind() == reflect.Struct) {
				continue
			}

			if f.Name != t.name {
				_, ok := f.Tag.Lookup("nbt")
				if !ok {
					continue
				}
			}

			d.decodeStruct(v.Field(i), t)
			fmt.Printf("%#v %#v\n", f.Name, t)
		}
	}
	return nil
}

func (d *Decoder) decodeStruct(v reflect.Value, t tag) {
	fmt.Println(1, d.p.tags)
	d.p.parse(t)
	fmt.Println(2, d.p.tags)

	rt := v.Type()
	for i := 0; i < v.NumField(); i++ {
		f := rt.Field(i)
		if !f.IsExported() {
			continue
		}

	}
}

func freeDecoder(d *Decoder) {
	decoders.Put(d)
}

func newDecoder() *Decoder {
	return decoders.Get().(*Decoder)
}

var decoders = sync.Pool{
	New: func() interface{} {
		return &Decoder{}
	},
}
