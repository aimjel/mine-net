package nbt

import (
	"fmt"
	"io"
	"reflect"
)

type Encoder struct {
	w io.Writer
}

func NewEncoder(w io.Writer) *Encoder {
	return &Encoder{w: w}
}

func (e *Encoder) Encode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Struct {
		//writes the nameless compound
		if _, err := e.w.Write([]byte{10, 0, 0}); err != nil {
			return err
		}

	}
	return e.encode(rv)
}

func (e *Encoder) encode(v reflect.Value) error {
	switch v.Kind() {

	case reflect.Struct:
		vt := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fieldType := vt.Field(i)

			if !fieldType.IsExported() {
				continue
			}

			fieldValue := v.Field(i)
			//writes tag id
			if id := nbtId(fieldValue); id == 0 {
				return fmt.Errorf("unknown reflect value")
			} else {
				_, _ = e.w.Write([]byte{id})
			}

			fieldName, ok := fieldType.Tag.Lookup("nbt")
			if !ok {
				fieldName = fieldType.Name
			}

			l := uint16(len(fieldName))
			//write name length
			_, _ = e.w.Write([]byte{byte(l >> 8), byte(l)})
			//write name
			_, _ = e.w.Write([]byte(fieldName))

			if err := e.encode(fieldValue); err != nil {
				return err
			}
		}

		//write end tag
		_, err := e.w.Write([]byte{0})
		return err

	case reflect.Slice:
		id := nbtId(v)

		switch id {

		case tagIntArray:
			l := v.Len()
			_, _ = e.w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})

			intArray := v.Interface().([]int32)
			for _, x := range intArray {
				_, _ = e.w.Write([]byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
			}

		case tagLongArray:
			l := v.Len()
			_, _ = e.w.Write([]byte{byte(l >> 24), byte(l >> 16), byte(l >> 8), byte(l)})

			longArray := v.Interface().([]int64)
			for _, x := range longArray {
				_, _ = e.w.Write([]byte{byte(x >> 56), byte(x >> 48), byte(x >> 40), byte(x >> 32), byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
			}
		}
	}

	return nil
}

func nbtId(v reflect.Value) byte {
	switch v.Kind() {

	case reflect.Struct:
		return tagCompound

	case reflect.Slice:

		switch v.Type().Elem().Kind() {

		case reflect.Int32:
			return tagIntArray

		case reflect.Int64:
			return tagLongArray
		}
	}

	return 0
}
