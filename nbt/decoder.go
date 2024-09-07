package nbt

import (
	"fmt"
	"io"
	"math"
	"reflect"
	"strings"
	"sync"
)

func Unmarshal(b []byte, v any) error {
	return (&Decoder{
		dec: newDecoderWithBytes(b),
	}).Decode(v)
}

type Decoder struct {
	dec *decoder
}

func NewDecoder(r io.Reader) *Decoder {
	return &Decoder{
		dec: newDecoder(r),
	}
}

func (d *Decoder) Decode(v any) error {
	rv := reflect.ValueOf(v)
	if rv.Kind() != reflect.Pointer {
		return fmt.Errorf("nbt: value passed must be a pointer")
	}

	id, err := d.dec.readTag()
	if err != nil {
		return err
	}

	return d.unmarshal(rv.Elem(), id)
}

func (d *Decoder) unmarshal(v reflect.Value, id byte) error {
	switch id {

	default:
		return fmt.Errorf("unknown tag id %v while unmarshalling", id)

	case tagByte:
		x, err := d.dec.readByte()
		if err != nil {
			return err
		}

		switch v.Kind() {

		default:
			return fmt.Errorf("nbt: cannot marshal byte tag into %v", v.Kind())

		case reflect.Int8:
			v.SetInt(int64(x))

		case reflect.Bool:
			v.SetBool(x == 1)

		case reflect.Interface:
			v.Set(reflect.ValueOf(x))
		}

	case tagShort:
		x, err := d.dec.readInt16()
		if err != nil {
			return err
		}
		switch v.Kind() {
		default:
			return fmt.Errorf("nbt: cannot marshal short tag into %v", v.Kind())

		case reflect.Int16:
			v.SetInt(int64(x))

		case reflect.Interface:
			v.Set(reflect.ValueOf(x))
		}

	case tagInt:
		x, err := d.dec.readInt32()
		if err != nil {
			return err
		}
		switch v.Kind() {
		case reflect.Int32, reflect.Int:
			v.SetInt(int64(x))
		case reflect.Interface:
			v.Set(reflect.ValueOf(x))
		}

	case tagLong:
		x, err := d.dec.readInt64()
		if err != nil {
			return err
		}
		switch v.Kind() {
		case reflect.Int64:
			v.SetInt(int64(x))
		case reflect.Interface:
			v.Set(reflect.ValueOf(x))
		}

	case tagFloat:
		x, err := d.dec.readInt32()
		if err != nil {
			return err
		}
		switch v.Kind() {
		default:
			return fmt.Errorf("nbt: cannot marshal float tag into %v", v.Kind())

		case reflect.Float32:
			v.SetFloat((float64)(math.Float32frombits((uint32)(x))))

		case reflect.Interface:
			v.Set(reflect.ValueOf((float64)(math.Float32frombits((uint32)(x)))))
		}

	case tagDouble:
		x, err := d.dec.readInt64()
		if err != nil {
			return err
		}
		switch v.Kind() {
		default:
			return fmt.Errorf("nbt: cannot marshal double tag into %v", v.Kind())

		case reflect.Float64:
			v.SetFloat(math.Float64frombits((uint64)(x)))

		case reflect.Interface:
			v.Set(reflect.ValueOf(math.Float64frombits((uint64)(x))))
		}

	case tagString:
		str, err := d.dec.readUnsafeString()
		if err != nil {
			return err
		}

		switch v.Kind() {
		default:
			return fmt.Errorf("nbt: cannot marshal string tag into %v", v.Kind())

		case reflect.String:
			v.SetString(strings.Clone(str))
		case reflect.Interface:
			v.Set(reflect.ValueOf(strings.Clone(str)))
		}

	case tagList, tagByteArray, tagIntArray, tagLongArray:
		if v.Kind() != reflect.Slice && v.Kind() != reflect.Interface {
			return fmt.Errorf("nbt: cannot marshal list tag into %v", v.Kind())
		}
		tagId := id
		var err error
		if id == tagList {
			tagId, err = d.dec.readByte() //get the type of list
			if err != nil {
				return err
			}
		}

		length, err := d.dec.readInt32()
		if err != nil {
			return err
		}

		return d.unmarshalList(v, tagId, length)

	case tagCompound:
		switch v.Kind() {
		case reflect.Interface:
			if v.IsNil() {
				v.Set(reflect.MakeMap(reflect.TypeOf(map[string]any{})))
				return d.unmarshalMap(v.Elem())
			}

		case reflect.Map:
			return d.unmarshalMap(v)

		case reflect.Struct:
			m := generateMap(v)
			defer func() {
				clear(m)
				valueMap.Put(m)
			}()

			for {
				tagId, er := d.dec.readByte()
				if er != nil {
					return er
				}

				if tagId == tagEnd {
					return nil
				}

				name, err := d.dec.readUnsafeString()
				if err != nil {
					return err
				}

				vv, ok := m[name]
				if !ok {
					if err = d.skip(tagId); err != nil {
						return err
					}
					continue
				}

				if err = d.unmarshal(vv, tagId); err != nil {
					return err
				}
			}

		default:
			return fmt.Errorf("nbt: cannot marshal compound tag into %v", v.Kind())
		}
	}

	return nil
}

var valueMap = sync.Pool{
	New: func() any {
		return make(map[string]reflect.Value)
	}}

func generateMap(v reflect.Value) map[string]reflect.Value {
	m := valueMap.Get().(map[string]reflect.Value)
	ty := v.Type()
	for i := 0; i < v.NumField(); i++ {
		ft := ty.Field(i)

		if !ft.IsExported() {
			continue
		}

		found := ft.Name
		if n, ok := ft.Tag.Lookup("nbt"); ok {
			found = n
		}
		if i := strings.Index(found, ",omitempty"); i != -1 {
			found = found[:i]
		}

		m[found] = v.Field(i)
	}

	return m
}

func (d *Decoder) unmarshalMap(v reflect.Value) error {
	if v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}

	for {
		id, err := d.dec.readByte()
		if err != nil {
			return err
		}
		if id == tagEnd {
			return nil
		}

		name, err := d.dec.readUnsafeString()
		if err != nil {
			return err
		}
		name = strings.Clone(name)

		switch id {

		case tagString:
			str, err := d.dec.readUnsafeString()
			if err != nil {
				return err
			}
			value := strings.Clone(str)
			if v.Type().Elem().Kind() == reflect.String {
				m := v.Interface().(map[string]string)
				m[name] = value
				continue
			}

			v.SetMapIndex(reflect.ValueOf(name), reflect.ValueOf(value))

		default:
			x := reflect.New(v.Type().Elem())
			if err := d.unmarshal(x.Elem(), id); err != nil {
				return fmt.Errorf("%v unmarshalling %v into map", err, name)
			}

			v.SetMapIndex(reflect.ValueOf(name), x.Elem())
		}
	}
}
func (d *Decoder) unmarshalList(v reflect.Value, id byte, length int) error {
	if v.Type().Kind() == reflect.Interface {
		var newV reflect.Value
		switch id {
		case tagByte, tagByteArray:
			newV = reflect.MakeSlice(reflect.TypeOf([]int8{}), length, length)
		case tagShort:
			newV = reflect.MakeSlice(reflect.TypeOf([]int16{}), length, length)
		case tagInt, tagIntArray:
			newV = reflect.MakeSlice(reflect.TypeOf([]int32{}), length, length)
		case tagLong, tagLongArray:
			newV = reflect.MakeSlice(reflect.TypeOf([]int64{}), length, length)
		case tagFloat:
			newV = reflect.MakeSlice(reflect.TypeOf([]float32{}), length, length)
		case tagDouble:
			newV = reflect.MakeSlice(reflect.TypeOf([]float64{}), length, length)
		case tagString:
			newV = reflect.MakeSlice(reflect.TypeOf([]string{}), length, length)
		case tagList, tagCompound:
			newV = reflect.MakeSlice(reflect.TypeOf([]any{}), length, length)
		}

		x := reflect.New(newV.Type()).Elem()
		x.Set(newV)
		err := d.unmarshalList(x, id, length)
		v.Set(x)
		return err
	}

	sliceType := v.Type().Elem().Kind()

	switch id {

	case tagByte, tagByteArray:
		if sliceType != reflect.Int8 {
			return fmt.Errorf("nbt: cannot marshal a byte array into %v", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := d.dec.readByte()
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
			b, err := d.dec.readInt16()
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
			b, err := d.dec.readInt32()
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
			b, err := d.dec.readInt64()
			if err != nil {
				return err
			}

			v.Index(i).SetInt(int64(b))
		}

	case tagFloat:
		if sliceType != reflect.Float32 {
			return fmt.Errorf("nbt: cannot marshal a float array into %v slice", sliceType)
		}

		v.Grow(length)
		v.SetLen(length)
		for i := 0; i < length; i++ {
			b, err := d.dec.readInt32()
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
			b, err := d.dec.readInt64()
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
			str, err := d.dec.readUnsafeString()
			if err != nil {
				return err
			}

			v.Index(i).SetString(strings.Clone(str))
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

func (d *Decoder) skip(id byte) error {
	switch id {
	default:
		return fmt.Errorf("unknown tag id %v while skipping", id)

	case tagByte, tagShort:
		return d.dec.skip(int(id))

	case tagInt, tagFloat:
		return d.dec.skip(4)

	case tagLong, tagDouble:
		return d.dec.skip(8)

	case tagString:
		_, err := d.dec.readUnsafeString()
		return err

	case tagList:
		tagId, err := d.dec.readByte()
		if err != nil {
			return err
		}

		ln, err := d.dec.readInt32()
		if err != nil {
			return err
		}

		if tagId == 0 || ln == 0 {
			return nil
		}

		for i := 0; i < ln; i++ {
			if err = d.skip(tagId); err != nil {
				return err
			}
		}

	case tagCompound:
		for {
			tagId, err := d.dec.readByte()
			if err != nil {
				return err
			}

			if tagId == tagEnd {
				return nil
			}

			_, err = d.dec.readUnsafeString()
			if err != nil {
				return err
			}

			if err = d.skip(tagId); err != nil {
				return err
			}
		}

	case tagByteArray, tagIntArray, tagLongArray:
		ln, err := d.dec.readInt32()
		if err != nil {
			return err
		}
		var size int
		switch id {
		case tagByteArray:
			size = 1
		case tagIntArray:
			size = 4
		case tagLongArray:
			size = 8
		}

		return d.dec.skip(size * ln)
	}

	return nil
}
