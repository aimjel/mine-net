package nbt

import (
	"fmt"
	"math"
	"reflect"
	"strings"
	"sync"
	"unsafe"
)

func Unmarshal(data []byte, v any) error {
	if err := checkValid(data); err != nil {
		return err
	}

	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("nbt: value passed must be a pointer")
	}

	dec := &decoder{buf: data}
	id, _ := dec.readTag()
	return dec.unmarshal(rv.Elem(), id)
}

func (d *decoder) unmarshal(v reflect.Value, id byte) error {
	switch id {

	case tagByte:
		if v.Kind() != reflect.Int8 {
			return fmt.Errorf("nbt: cannot marshal byte tag into %v", v.Kind())
		}
		v.SetInt(int64(d.readByte()))

	case tagShort:
		if v.Kind() != reflect.Int16 {
			return fmt.Errorf("nbt: cannot marshal short tag into %v", v.Kind())
		}
		v.SetInt(int64(d.readShort()))

	case tagInt:
		if v.Kind() != reflect.Int32 {
			return fmt.Errorf("nbt: cannot marshal int tag into %v", v.Kind())
		}
		v.SetInt(int64(d.readInt()))

	case tagLong:
		if v.Kind() != reflect.Int64 {
			return fmt.Errorf("nbt: cannot marshal long tag into %v", v.Kind())
		}
		v.SetInt(d.readLong())

	case tagFloat:
		if v.Kind() != reflect.Float32 {
			return fmt.Errorf("nbt: cannot marshal float tag into %v", v.Kind())
		}
		v.SetFloat((float64)(math.Float32frombits((uint32)(d.readInt()))))

	case tagDouble:
		if v.Kind() != reflect.Float64 {
			return fmt.Errorf("nbt: cannot marshal double tag into %v", v.Kind())
		}
		v.SetFloat(math.Float64frombits((uint64)(d.readLong())))

	case tagString:
		if v.Kind() != reflect.String {
			return fmt.Errorf("nbt: cannot marshal string tag into %v", v.Kind())
		}

		v.SetString(strings.Clone(d.readUnsafeString()))

	case tagList, tagByteArray, tagIntArray, tagLongArray:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("nbt: cannot marshal list tag into %v", v.Kind())
		}
		tagId := id
		if id == tagList {
			tagId = d.readByte() //get the type of list
		}

		return d.unmarshalList(v, tagId)

	case tagCompound:
		switch v.Kind() {

		case reflect.Map:
			d.unmarshalCompoundMap(v)

		case reflect.Struct:
			m := nameTags.Get().(map[string]nameTagMetaData)
			defer func() {
				for k := range m {
					delete(m, k)
				}
				nameTags.Put(m)
			}()
			endPos := d.fillMap(m)

			t := v.Type()
			for i := 0; i < v.NumField(); i++ {
				f := t.Field(i)
				if !f.IsExported() {
					continue
				}

				key, ok := f.Tag.Lookup("nbt")
				if !ok {
					key = f.Name
				}

				if tag, ok := m[key]; ok {
					d.at = tag.Index()
					if err := d.unmarshal(v.Field(i), tag.Type()); err != nil {
						return err
					}
				}
			}

			d.at = endPos

		default:
			return fmt.Errorf("nbt: cannot marshal compound tag into %v", v.Kind())
		}
	}

	return nil
}

func (d *decoder) unmarshalCompoundMap(v reflect.Value) {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	for {
		id := d.readByte()
		if id == tagEnd {
			return
		}

		name := strings.Clone(d.readUnsafeString())

		switch id {

		case tagString:
			value := strings.Clone(d.readUnsafeString())
			if v.Type().Elem().Kind() == reflect.String {
				m := v.Interface().(map[string]string)
				m[name] = value
				continue
			}

			v.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
		}
	}
}

func (d *decoder) fillMap(m map[string]nameTagMetaData) int {
	s := &scanner{buf: d.buf, at: d.at}
	for {
		id := s.buf[s.at]
		s.at++

		if id == tagEnd {
			return s.at
		}

		length := int(uint16(s.buf[s.at])<<8 | uint16(s.buf[s.at+1]))
		s.at += 2

		str := s.buf[s.at : s.at+length]
		key := *(*string)(unsafe.Pointer(&str))
		s.at += length
		m[key] = nameTagMetaData(int(id)<<59 | s.at)

		s.scan(id)
	}
}

func (d *decoder) unmarshalList(v reflect.Value, id byte) error {
	length := (int)(d.readInt())
	sliceType := v.Type().Elem().Kind()
	switch id {

	case tagByte, tagByteArray:
		if sliceType != reflect.Int8 {
			return fmt.Errorf("nbt: cannot marshal a byte array into %v", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetInt((int64)(d.readByte()))
		}

	case tagShort:
		if sliceType != reflect.Int16 {
			return fmt.Errorf("nbt: cannot marshal a short array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetInt(int64(d.readShort()))
		}

	case tagInt, tagIntArray:
		if sliceType != reflect.Int32 {
			return fmt.Errorf("nbt: cannot marshal a int array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetInt(int64(d.readInt()))
		}

	case tagLong, tagLongArray:
		if sliceType != reflect.Int64 {
			return fmt.Errorf("nbt: cannot marshal a long array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetInt(d.readLong())
		}

	case tagFloat:
		if sliceType != reflect.Float32 {
			return fmt.Errorf("nbt: cannot marshal a float array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetFloat(float64(math.Float32frombits((uint32)(d.readInt()))))
		}

	case tagDouble:
		if sliceType != reflect.Float64 {
			return fmt.Errorf("nbt: cannot marshal a double array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetFloat(math.Float64frombits((uint64)(d.readLong())))
		}

	case tagString:
		if sliceType != reflect.String {
			return fmt.Errorf("nbt: cannot marshal a string array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			v.Index(i).SetString(strings.Clone(d.readUnsafeString()))
		}

	case tagList, tagCompound:
		v.Grow(length)
		v.SetLen(length)

		for i := 0; i < length; i++ {
			if err := d.unmarshal(v.Index(i), id); err != nil {
				return err
			}
		}
	}

	return nil
}

type nameTagMetaData int

func (d nameTagMetaData) Index() int {
	return *(*int)(&d) & 0x7ffffffffffffff
}

func (d nameTagMetaData) Type() byte {
	return (byte)(d >> 59)
}

// nameTags returns a map where the nbt type and position of name tags are stored
var nameTags = sync.Pool{
	New: func() any {
		return make(map[string]nameTagMetaData)
	},
}

type decoder struct {
	buf []byte
	at  int
}

func (d *decoder) readTag() (id byte, key string) {
	if id = d.readByte(); id == tagCompound {
		key = d.readUnsafeString()
	}

	return
}

func (d *decoder) readByte() byte {
	v := d.buf[d.at]
	d.at++
	return v
}

func (d *decoder) readShort() int16 {
	v := int16(d.buf[d.at])<<8 | int16(d.buf[d.at+1])
	d.at += 2
	return v
}

func (d *decoder) readInt() int32 {
	v := int32(d.buf[d.at])<<24 | int32(d.buf[d.at+1])<<16 | int32(d.buf[d.at+2])<<8 | int32(d.buf[d.at+3])
	d.at += 4
	return v
}

func (d *decoder) readLong() int64 {
	v := int64(d.buf[d.at])<<56 | int64(d.buf[d.at+1])<<48 | int64(d.buf[d.at+2])<<40 | int64(d.buf[d.at+3])<<32 | int64(d.buf[d.at+4])<<24 | int64(d.buf[d.at+5])<<16 | int64(d.buf[d.at+6])<<8 | int64(d.buf[d.at+7])
	d.at += 8
	return v
}

func (d *decoder) readUnsafeString() string {
	v := int(d.buf[d.at])<<8 | int(d.buf[d.at+1])
	d.at += 2

	str := d.buf[d.at : d.at+v]
	d.at += v
	return *(*string)(unsafe.Pointer(&str))
}

const (
	tagEnd = iota
	tagByte
	tagShort
	tagInt
	tagLong
	tagFloat
	tagDouble
	tagByteArray
	tagString
	tagList
	tagCompound
	tagIntArray
	tagLongArray
)
