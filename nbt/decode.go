package nbt

import (
	"fmt"
	"reflect"
)

func Unmarshal(data []byte, v interface{}) error {
	if err := checkValid(data); err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Ptr {
		return fmt.Errorf("nbt: interface must be a pointer")
	}

	// 2 bytes for string length, 1 accounts for starting 0 index
	length := int(uint16(data[1])<<8|uint16(data[2])) + 2 + 1 //FIXME

	t := (&parser{buf: data, at: length}).parseTag(data[0])

	return compound(rv.Elem(), t)
}

func compound(v reflect.Value, t tag) error {
	if v.Kind() == reflect.Map && t.typ == tagCompound {
		v.Set(reflect.ValueOf(t.c.Map()))
		return nil
	}

	if v.Kind() == reflect.Struct && t.typ == tagCompound {

		rt := v.Type()
		for i := 0; i < v.NumField(); i++ {

			f := rt.Field(i)

			if f.IsExported() == false {
				continue
			}

			fieldKey := f.Name
			if k, ok := f.Tag.Lookup("nbt"); ok {
				fieldKey = k
			}

			nbtTagType := typeToTag(f.Type)
			if nbtTagType != 0 {
				//fmt.Printf("nbtTagType %v, field key %v\n", nbtTagType, fieldKey)
				//fmt.Printf("named tags %+v\n", t.c.namedTags)
				if nbtValue, ok := t.c.getTag(nbtTagType, fieldKey); ok {

					if nbtTagType == tagCompound {
						if err := compound(v.Field(i), nbtValue); err != nil {
							return err
						}
					} else if nbtTagType == tagList {
						if err := list(v.Field(i), nbtValue); err != nil {
							return err
						}
					} else {
						v.Field(i).Set(nbtValue.rv)
					}
				}
			}
		}
	}

	return nil
}

func list(v reflect.Value, t tag) error {
	if v.Type().Elem().Kind() == reflect.Struct {
		s := reflect.MakeSlice(v.Type(), len(t.l), len(t.l))

		for i := 0; i < len(t.l); i++ {
			if err := compound(s.Index(i), t.l[i]); err != nil {
				return err
			}
		}

		v.Set(s)
		return nil
	}

	s := reflect.MakeSlice(v.Type(), len(t.l), len(t.l))

	for i := 0; i < len(t.l); i++ {
		s.Index(i).Set(t.l[i].rv)
	}

	v.Set(s)

	return nil
}

func typeToTag(t reflect.Type) byte {
	switch t.Kind() {

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
		if t.Elem().Kind() == reflect.Int32 {
			return tagIntArray
		} else if t.Elem().Kind() == reflect.Int64 {
			return tagLongArray
		}

		return tagList
	}

	return 0
}
