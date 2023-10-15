package nbt

import (
	"fmt"
	"io"
	"reflect"
	"unsafe"
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

	case reflect.Int8:
		e.w.Write([]byte{byte(v.Int())})

	case reflect.Int16:
		e.write16(int(v.Int()))

	case reflect.Int32:
		e.write32(int(v.Int()))

	case reflect.Int64:
		e.write64(v.Int())

	case reflect.String:
		e.writeString(v.String())

	case reflect.Slice:
		switch v.Type().Elem().Kind() {
		default:
			_, _ = e.w.Write([]byte{nbtId(v.Type().Elem())})
			ln := v.Len()
			e.write32(ln)

			for i := 0; i < ln; i++ {
				if err := e.encode(v.Index(i)); err != nil {
					return fmt.Errorf("%v encoding %v element in list tag", err, i)
				}
			}

		case reflect.Slice:
			_, _ = e.w.Write([]byte{nbtId(v.Type().Elem())})
			ln := v.Len()
			e.write32(ln)

			for i := 0; i < ln; i++ {
				if err := e.encode(v.Index(i)); err != nil {
					return fmt.Errorf("%v encoding %v element in list tag", err, i)
				}
			}

		case reflect.Int8:
			ln := v.Len()
			e.write32(ln)

			x := v.Interface().([]int8)
			e.w.Write(*(*[]byte)(unsafe.Pointer(&x)))

		case reflect.Int64:
			ln := v.Len()
			e.write32(ln)

			x := v.Interface().([]int64)
			for i := 0; i < ln; i++ {
				e.write64(x[i])
			}

		case reflect.Int32:
			ln := v.Len()
			e.write32(ln)

			x := v.Interface().([]int32)
			for i := 0; i < ln; i++ {
				e.write32(int(x[i]))
			}
		}

	case reflect.Struct:
		t := v.Type()
		for i := 0; i < v.NumField(); i++ {
			fieldType := t.Field(i)
			if !fieldType.IsExported() {
				continue
			}

			name, ok := fieldType.Tag.Lookup("nbt")
			if !ok {
				name = fieldType.Name
			}

			fv := v.Field(i)
			if fv.Kind() == reflect.Interface {
				//extracts the value hiding behind the interface
				fv = v.Elem()
			}
			if fv.Kind() == reflect.Map {
				if fv.IsNil() {
					continue
				}
			}

			if err := e.encodeNameTag(fv, name); err != nil {
				return fmt.Errorf("%v encoding %v in compound tag", err, name)
			}
		}

		e.w.Write([]byte{tagEnd})

	case reflect.Map:
		typ := v.Type().Key()
		if typ.Kind() != reflect.String {
			return fmt.Errorf("cant encode map with %v keys", v.Kind())
		}

		iter := v.MapRange()
		for iter.Next() {
			k := iter.Key()
			val := iter.Value()

			if val.Kind() == reflect.Interface {
				val = val.Elem()
			}

			if val.Kind() == reflect.Map {
				if val.IsNil() {
					continue
				}
			}

			if err := e.encodeNameTag(val, k.String()); err != nil {
				return fmt.Errorf("%v encoding %v in compound tag", err, k.String())
			}
		}
		e.w.Write([]byte{tagEnd})
	}

	return nil
}

func (e *Encoder) encodeNameTag(v reflect.Value, name string) error {
	_, _ = e.w.Write([]byte{nbtId(v.Type())})
	e.writeString(name)

	return e.encode(v)
}

func (e *Encoder) writeString(x string) {
	e.write16(len(x))
	_, _ = e.w.Write(*(*[]byte)(unsafe.Pointer(&x)))
}

func (e *Encoder) write16(x int) {
	e.w.Write([]byte{byte(x >> 8), byte(x)})
}

func (e *Encoder) write32(x int) {
	e.w.Write([]byte{byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
}

func (e *Encoder) write64(x int64) {
	e.w.Write([]byte{
		byte(x >> 56), byte(x >> 48), byte(x >> 40), byte(x >> 32),
		byte(x >> 24), byte(x >> 16), byte(x >> 8), byte(x)})
}

func nbtId(v reflect.Type) byte {
	switch v.Kind() {

	case reflect.Int8:
		return tagByte

	case reflect.Int16:
		return tagShort

	case reflect.Int32:
		return tagInt

	case reflect.Int64:
		return tagLong

	case reflect.Float32:
		return tagFloat

	case reflect.Float64:
		return tagDouble

	case reflect.String:
		return tagString

	case reflect.Struct, reflect.Map:
		return tagCompound

	case reflect.Slice:
		switch v.Elem().Kind() {

		default:
			return tagList

		case reflect.Int8:
			return tagByteArray

		case reflect.Int32:
			return tagIntArray

		case reflect.Int64:
			return tagLongArray
		}
	}

	panic("shouldnt happen")
}
