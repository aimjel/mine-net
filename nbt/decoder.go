package nbt

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"unsafe"
)

type Decoder struct {
	rd io.Reader

	//buf is a simple buffer to read values
	buf []byte
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{rd: r, buf: make([]byte, 8)}
}

func (dec *Decoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("nbt: value passed must be a pointer")
	}

	id, _, err := dec.readTag()
	if err != nil {
		return err
	}

	return dec.unmarshal(rv.Elem(), id)
}

func (dec *Decoder) unmarshal(v reflect.Value, id byte) error {
	switch id {

	case tagByte:
		switch v.Kind() {

		default:
			return fmt.Errorf("nbt: cannot marshal byte tag into %v", v.Kind())

		case reflect.Int8:
			x, err := dec.readByte()
			if err != nil {
				return err
			}

			v.SetInt(int64(x))

		case reflect.Bool:
			x, err := dec.readByte()
			if err != nil {
				return err
			}

			if x == 1 {
				v.SetBool(true)
			}
		}

	case tagShort:
		if v.Kind() != reflect.Int16 {
			return fmt.Errorf("nbt: cannot marshal short tag into %v", v.Kind())
		}

		x, err := dec.readShort()
		if err != nil {
			return err
		}
		v.SetInt(int64(x))

	case tagInt:
		if v.Kind() != reflect.Int32 && v.Kind() != reflect.Int {
			return fmt.Errorf("nbt: cannot marshal int tag into %v", v.Kind())
		}

		x, err := dec.readInt()
		if err != nil {
			return err
		}
		v.SetInt(int64(x))

	case tagLong:
		if v.Kind() != reflect.Int64 {
			return fmt.Errorf("nbt: cannot marshal long tag into %v", v.Kind())
		}
		x, err := dec.readLong()
		if err != nil {
			return err
		}
		v.SetInt(x)

	case tagFloat:
		if v.Kind() != reflect.Float32 {
			return fmt.Errorf("nbt: cannot marshal float tag into %v", v.Kind())
		}

		x, err := dec.readInt()
		if err != nil {
			return err
		}

		v.SetFloat((float64)(math.Float32frombits((uint32)(x))))

	case tagDouble:
		if v.Kind() != reflect.Float64 {
			return fmt.Errorf("nbt: cannot marshal double tag into %v", v.Kind())
		}
		x, err := dec.readLong()
		if err != nil {
			return err
		}

		v.SetFloat(math.Float64frombits((uint64)(x)))

	case tagString:
		if v.Kind() != reflect.String {
			return fmt.Errorf("nbt: cannot marshal string tag into %v", v.Kind())
		}

		str, err := dec.readString()
		if err != nil {
			return err
		}
		v.SetString(str)

	case tagList, tagByteArray, tagIntArray, tagLongArray:
		if v.Kind() != reflect.Slice {
			return fmt.Errorf("nbt: cannot marshal list tag into %v", v.Kind())
		}
		tagId := id
		var err error
		if id == tagList {
			tagId, err = dec.readByte() //get the type of list
			if err != nil {
				return err
			}
		}

		return dec.unmarshalList(v, tagId)

	case tagCompound:
		switch v.Kind() {

		case reflect.Map:
			return dec.unmarshalCompoundMap(v)

		case reflect.Struct:
			finder := newStructFields(v)
			for {
				tagId, er := dec.readByte()
				if er != nil {
					return er
				}

				if tagId == tagEnd {
					return nil
				}

				name, err := dec.readString()
				if err != nil {
					return err
				}

				//fmt.Println("trying to find", name, tagName(tagId))
				vv, ok := finder.find(name)
				if !ok {
					//fmt.Println("skipping", name, tagName(tagId))
					if err = dec.skip(tagId); err != nil {
						return err
					}
					//todo ability to read the tag without using unmarshal

					continue
				}

				//fmt.Println(vv.Kind(), tagName(tagId), name)
				if err = dec.unmarshal(vv, tagId); err != nil {
					return err
				}

			}

		default:
			return fmt.Errorf("nbt: cannot marshal compound tag into %v", v.Kind())
		}
	}

	return nil
}

func (dec *Decoder) skip(id byte) error {
	switch id {

	default:
		panic(id)

	case tagByte:
		if _, err := dec.rd.Read(dec.buf[:1]); err != nil {
			return err
		}

	case tagInt:
		if _, err := dec.rd.Read(dec.buf[:4]); err != nil {
			return err
		}

	case tagDouble:

	case tagString:
		if err := dec.skipString(); err != nil {
			return err
		}

	case tagList:
		tagId, err := dec.readByte()
		if err != nil {
			return err
		}

		ln, err := dec.readInt()
		if err != nil {
			return err
		}

		if tagId == 0 || ln == 0 {
			return nil
		}

		for i := 0; i < int(ln); i++ {
			if err = dec.skip(tagId); err != nil {
				return err
			}
		}

	case tagCompound:
		for {
			tagId, err := dec.readByte()
			if err != nil {
				return err
			}

			if tagId == tagEnd {
				return nil
			}

			_, err = dec.readString()
			if err != nil {
				return err
			}

			if err = dec.skip(tagId); err != nil {
				return err
			}
		}

	case tagByteArray:
		ln, err := dec.readInt()
		if err != nil {
			return err
		}

		for i := 0; i < int(ln); i++ {
			_, err = dec.rd.Read(dec.buf[:1])
			if err != nil {
				return err
			}
		}

	case tagIntArray:
		ln, err := dec.readInt()
		if err != nil {
			return err
		}

		for i := 0; i < int(ln); i++ {
			_, err = dec.rd.Read(dec.buf[:4])
			if err != nil {
				return err
			}
		}
	case tagLongArray:
		ln, err := dec.readInt()
		if err != nil {
			return err
		}

		for i := 0; i < int(ln); i++ {
			_, err = dec.rd.Read(dec.buf[:8])
			if err != nil {
				return err
			}
		}

	}

	return nil
}

func (dec *Decoder) unmarshalCompoundMap(v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	for {
		id, err := dec.readByte()
		if err != nil {
			return err
		}
		if id == tagEnd {
			return nil
		}

		name, err := dec.readString()
		if err != nil {
			return err
		}

		switch id {
		default:
			fmt.Println("cant unmarshal", tagName(id), "in map")

		case tagList:
			x := reflect.New(v.Type().Elem())

			if err := dec.unmarshal(x.Elem(), tagList); err != nil {
				return fmt.Errorf("%v unmarshalling %v in compound map", err, name)
			}

			v.SetMapIndex(reflect.ValueOf(name), x.Elem())

		case tagCompound:
			x := reflect.New(v.Type().Elem())

			if err := dec.unmarshal(x.Elem(), tagCompound); err != nil {
				return fmt.Errorf("%v unmarshalling %v in compound map", err, name)
			}

			v.SetMapIndex(reflect.ValueOf(name), x.Elem())

		case tagString:
			value, err := dec.readString()
			if err != nil {
				return err
			}
			if v.Type().Elem().Kind() == reflect.String {
				m := v.Interface().(map[string]string)
				m[name] = value
				continue
			}

			v.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))
		}
	}
}

func (dec *Decoder) unmarshalList(v reflect.Value, id byte) error {
	l, err := dec.readInt()
	if err != nil {
		return err
	}

	length := int(l)

	sliceType := v.Type().Elem().Kind()
	switch id {

	case tagByte, tagByteArray:
		if sliceType != reflect.Int8 {
			return fmt.Errorf("nbt: cannot marshal a byte array into %v", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readByte()
			if err != nil {
				return err
			}
			v.Index(i).SetInt((int64)(b))
		}

	case tagShort:
		if sliceType != reflect.Int16 {
			return fmt.Errorf("nbt: cannot marshal a short array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readShort()
			if err != nil {
				return err
			}
			v.Index(i).SetInt(int64(b))
		}

	case tagInt, tagIntArray:
		if sliceType != reflect.Int32 {
			return fmt.Errorf("nbt: cannot marshal a int array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readInt()
			if err != nil {
				return err
			}

			v.Index(i).SetInt(int64(b))
		}

	case tagLong, tagLongArray:
		if sliceType != reflect.Int64 {
			return fmt.Errorf("nbt: cannot marshal a long array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readLong()
			if err != nil {
				return err
			}

			v.Index(i).SetInt(b)
		}

	case tagFloat:
		if sliceType != reflect.Float32 {
			return fmt.Errorf("nbt: cannot marshal a float array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readInt()
			if err != nil {
				return err
			}

			v.Index(i).SetFloat(float64(math.Float32frombits(uint32(b))))
		}

	case tagDouble:
		if sliceType != reflect.Float64 {
			return fmt.Errorf("nbt: cannot marshal a double array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := dec.readLong()
			if err != nil {
				return err
			}

			v.Index(i).SetFloat(math.Float64frombits((uint64)(b)))
		}

	case tagString:
		if sliceType != reflect.String {
			return fmt.Errorf("nbt: cannot marshal a string array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			str, err := dec.readString()
			if err != nil {
				return err
			}

			v.Index(i).SetString(str)
		}

	case tagList, tagCompound:
		v.Grow(length)
		v.SetLen(length)

		for i := 0; i < length; i++ {
			if err := dec.unmarshal(v.Index(i), id); err != nil {
				return err
			}
		}
	}

	return nil
}

// structFields used to search fields by any string
type structFields struct {
	og reflect.Value

	vals []reflect.Value
}

func (s *structFields) find(name string) (reflect.Value, bool) {
	ty := s.og.Type()
	for i, val := range s.vals {
		ft := ty.Field(i)

		if !ft.IsExported() {
			continue
		}

		found := ft.Name
		if n, ok := ft.Tag.Lookup("nbt"); ok {
			found = n
		}

		if found == name {
			return val, true
		}
	}

	return reflect.Value{}, false
}

func newStructFields(v reflect.Value) (s structFields) {
	s.og = v

	s.vals = make([]reflect.Value, 0, v.NumField())
	for i := 0; i < v.NumField(); i++ {
		s.vals = append(s.vals, v.Field(i))
	}

	return
}

func (dec *Decoder) readTag() (id byte, key string, err error) {
	if id, err = dec.readByte(); err != nil {
		return
	}

	if id == tagCompound {
		key, err = dec.readString()
	}

	return
}

func (dec *Decoder) readByte() (byte, error) {
	_, err := dec.rd.Read(dec.buf[:1])

	return dec.buf[0], err
}

func (dec *Decoder) readShort() (int16, error) {
	_, err := dec.rd.Read(dec.buf[:2])

	v := int16(dec.buf[0])<<8 | int16(dec.buf[1])
	return v, err
}

func (dec *Decoder) readInt() (int32, error) {
	_, err := dec.rd.Read(dec.buf[:4])

	v := int32(dec.buf[0])<<24 | int32(dec.buf[1])<<16 | int32(dec.buf[2])<<8 | int32(dec.buf[3])
	return v, err
}

func (dec *Decoder) readLong() (int64, error) {
	_, err := dec.rd.Read(dec.buf)

	v := int64(dec.buf[0])<<56 | int64(dec.buf[1])<<48 | int64(dec.buf[2])<<40 | int64(dec.buf[3])<<32 |
		int64(dec.buf[4])<<24 | int64(dec.buf[5])<<16 | int64(dec.buf[6])<<8 | int64(dec.buf[7])
	return v, err
}

func (dec *Decoder) readString() (string, error) {
	ln, err := dec.readShort()
	if err != nil {
		return "", err
	}

	str := make([]byte, ln)
	_, err = dec.rd.Read(str)
	return *(*string)(unsafe.Pointer(&str)), err
}

func (dec *Decoder) skipString() error {
	ln, err := dec.readShort()
	if err != nil {
		return err
	}

	//how many bytes are left to read
	left := int(ln)
	for left != 0 {
		maxx := 8
		if left <= len(dec.buf) {
			maxx = left
		}

		n, err := dec.rd.Read(dec.buf[:maxx])
		if err != nil {
			return err
		}

		left -= n
	}

	return nil
}
