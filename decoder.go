package bencode

import (
	"bufio"
	"bytes"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
)

// Decoder bencode decoder
type Decoder struct {
	r io.Reader
}

// NewDecoder create decoder from io.Reader
func NewDecoder(r io.Reader) Decoder {
	return Decoder{r: r}
}

// Decode decode data
func (dec Decoder) Decode(data interface{}) error {
	return decode(bufio.NewReader(dec.r), reflect.ValueOf(data))
}

// Decode decode data in raw
func Decode(data []byte, value interface{}) error {
	return NewDecoder(bytes.NewReader(data)).Decode(value)
}

func decode(r *bufio.Reader, v reflect.Value) error {
	if v.Kind() != reflect.Ptr {
		return errors.New("input value is not pointer")
	}
	ch, err := r.ReadByte()
	if err != nil {
		return err
	}
	switch ch {
	case 'i':
		n, err := decodeNumber(r, v.Elem())
		if err != nil {
			return err
		}
		if n.isUnsigned {
			v.SetUint(n.unsigned)
		} else {
			v.SetInt(n.signed)
		}
		return nil
	case 'd':
		return decodeDict(r, v.Elem())
	case 'l':
		return decodeList(r, v.Elem())
	default:
		return decodeString(r, v.Elem(), ch)
	}
}

type number struct {
	isUnsigned bool
	signed     int64
	unsigned   uint64
}

func decodeNumber(r *bufio.Reader, v reflect.Value) (number, error) {
	var ret number
	var str []byte
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return ret, fmt.Errorf("decode number: %v", err)
		}
		if ch == 'e' {
			switch v.Kind() {
			case reflect.Int,
				reflect.Int8, reflect.Int16,
				reflect.Int32, reflect.Int64:
				n, err := strconv.ParseInt(string(str), 10, v.Type().Bits())
				if err != nil {
					return ret, fmt.Errorf("can not parse %s to %s value", string(str), v.Kind().String())
				}
				return number{
					isUnsigned: false,
					signed:     n,
				}, nil
			case reflect.Uint,
				reflect.Uint8, reflect.Uint16,
				reflect.Uint32, reflect.Uint64:
				n, err := strconv.ParseUint(string(str), 10, v.Type().Bits())
				if err != nil {
					return ret, fmt.Errorf("can not parse %s to %s value", string(str), v.Kind().String())
				}
				return number{
					isUnsigned: true,
					unsigned:   n,
				}, nil
			default:
				return ret, fmt.Errorf("can not set number value to variable of type %s", v.Kind().String())
			}
		}
		str = append(str, ch)
	}
}

func decodeDict(r *bufio.Reader, v reflect.Value) error {
	if v.Kind() == reflect.Map && v.IsNil() {
		v.Set(reflect.MakeMap(v.Type()))
	}
	key := reflect.New(reflect.TypeOf(""))
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("decode dict: %v", err)
		}
		if ch == 'e' {
			return nil
		}
		err = decodeString(r, key.Elem(), ch)
		if err != nil {
			return err
		}
		ch, err = r.ReadByte()
		switch ch {
		case 'i':
			err = setDictNumber(r, key.Elem().String(), v)
		case 'd':
			err = setDictDict(r, key.Elem().String(), v)
		default:
			err = setDictString(r, key.Elem().String(), v, ch)
		}
		if err != nil {
			return err
		}
	}
}

func decodeList(r *bufio.Reader, v reflect.Value) error {
	if v.Kind() != reflect.Slice && v.Kind() != reflect.Array {
		return fmt.Errorf("can not set list value to variable of type %s", v.Kind().String())
	}
	i := 0
	for {
		fmt.Println(decode(r, v.Index(i)))
		i++
	}
}

func decodeString(r *bufio.Reader, v reflect.Value, ch byte) error {
	var len []byte
	len = append(len, ch)
	for {
		ch, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("decode string: %v", err)
		}
		if ch == ':' {
			size, err := strconv.ParseUint(string(len), 10, 64)
			if err != nil {
				return fmt.Errorf("can not parse string size: %s", string(len))
			}
			data := make([]byte, size)
			for i := 0; uint64(i) < size; i++ {
				data[i], err = r.ReadByte()
				if err != nil {
					return fmt.Errorf("decode string value: %v", err)
				}
			}
			switch v.Kind() {
			case reflect.String:
				v.SetString(string(data))
			case reflect.Slice, reflect.Array:
				if v.Type().ConvertibleTo(bytesType) {
					v.SetBytes([]byte(string(data)))
					return nil
				}
				min := size
				if uint64(v.Len()) < min {
					min = uint64(v.Len())
				}
				bt := []byte(string(data))
				for i := 0; i < int(min); i++ {
					if v.Index(i).Kind() != reflect.Uint8 {
						return fmt.Errorf("can not set string value to variable of type %s", v.Kind().String())
					}
					v.Index(i).SetUint(uint64(bt[i]))
				}
			default:
				return fmt.Errorf("can not set string value to variable of type %s", v.Kind().String())
			}
			return nil
		}
		len = append(len, ch)
	}
}

func setDictDict(r *bufio.Reader, key string, v reflect.Value) error {
	run := func(v reflect.Value) (error, bool) {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if kField.Tag.Get("bencode") == key {
				return decodeDict(r, v.Field(i)), true
			}
		}
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if strings.ToLower(kField.Name) == key {
				return decodeDict(r, v.Field(i)), true
			}
		}
		return nil, false
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		kField := t.Field(i)
		if kField.Anonymous {
			vField := v.Field(i)
			err, ok := run(vField)
			if err != nil {
				return err
			}
			if ok {
				return err
			}
		}
	}
	err, ok := run(v)
	if err != nil {
		return err
	}
	if ok {
		return err
	}
	return decodeDict(r, reflect.New(reflect.StructOf(nil)).Elem())
}

func numberByType(n number, t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(n.signed))
	case reflect.Int8:
		return reflect.ValueOf(int8(n.signed))
	case reflect.Int16:
		return reflect.ValueOf(int16(n.signed))
	case reflect.Int32:
		return reflect.ValueOf(int32(n.signed))
	case reflect.Int64:
		return reflect.ValueOf(int64(n.signed))
	case reflect.Uint:
		return reflect.ValueOf(uint(n.unsigned))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(n.unsigned))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(n.unsigned))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(n.unsigned))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(n.unsigned))
	default:
		return reflect.ValueOf(nil)
	}
}

func newNumberValue(t reflect.Type) reflect.Value {
	switch t.Kind() {
	case reflect.Int:
		return reflect.ValueOf(int(0))
	case reflect.Int8:
		return reflect.ValueOf(int8(0))
	case reflect.Int16:
		return reflect.ValueOf(int16(0))
	case reflect.Int32:
		return reflect.ValueOf(int32(0))
	case reflect.Int64:
		return reflect.ValueOf(int64(0))
	case reflect.Uint:
		return reflect.ValueOf(uint(0))
	case reflect.Uint8:
		return reflect.ValueOf(uint8(0))
	case reflect.Uint16:
		return reflect.ValueOf(uint16(0))
	case reflect.Uint32:
		return reflect.ValueOf(uint32(0))
	case reflect.Uint64:
		return reflect.ValueOf(uint64(0))
	default:
		return reflect.ValueOf(nil)
	}
}

func setDictNumber(r *bufio.Reader, key string, v reflect.Value) error {
	if v.Kind() == reflect.Map {
		kvalue := reflect.ValueOf(key)
		vn := v.MapIndex(kvalue)
		if !vn.IsValid() {
			if v.Type().Elem().Kind() == reflect.Interface {
				vn = newNumberValue(reflect.TypeOf(0))
			} else {
				vn = newNumberValue(v.Type().Elem())
			}
		}
		n, err := decodeNumber(r, vn)
		if err != nil {
			return err
		}
		v.SetMapIndex(kvalue, numberByType(n, vn.Type()))
		return nil
	}
	run := func(v reflect.Value) (error, bool) {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if kField.Tag.Get("bencode") == key {
				n, err := decodeNumber(r, v.Field(i))
				if err != nil {
					return err, true
				}
				if n.isUnsigned {
					v.SetUint(n.unsigned)
				} else {
					v.SetInt(n.signed)
				}
				return nil, true
			}
		}
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if strings.ToLower(kField.Name) == key {
				n, err := decodeNumber(r, v.Field(i))
				if err != nil {
					return err, true
				}
				if n.isUnsigned {
					v.SetUint(n.unsigned)
				} else {
					v.SetInt(n.signed)
				}
				return nil, true
			}
		}
		return nil, false
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		kField := t.Field(i)
		if kField.Anonymous {
			vField := v.Field(i)
			err, ok := run(vField)
			if err != nil {
				return err
			}
			if ok {
				return err
			}
		}
	}
	err, ok := run(v)
	if err != nil {
		return err
	}
	if ok {
		return err
	}
	_, err = decodeNumber(r, reflect.New(reflect.TypeOf(0)).Elem())
	if err != nil {
		return err
	}
	return nil
}

func setDictString(r *bufio.Reader, key string, v reflect.Value, ch byte) error {
	if v.Kind() == reflect.Map {
		kvalue := reflect.ValueOf(key)
		vn := v.MapIndex(kvalue)
		if !vn.IsValid() {
			t := reflect.TypeOf("")
			if v.Type().Elem().Kind() != reflect.Interface {
				t = v.Type().Elem()
			}
			vn = reflect.New(t).Elem()
		}
		err := decodeString(r, vn, ch)
		if err != nil {
			return err
		}
		v.SetMapIndex(kvalue, vn)
		return nil
	}
	run := func(v reflect.Value) (error, bool) {
		t := v.Type()
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if kField.Tag.Get("bencode") == key {
				return decodeString(r, v.Field(i), ch), true
			}
		}
		for i := 0; i < t.NumField(); i++ {
			kField := t.Field(i)
			if strings.ToLower(kField.Name) == key {
				return decodeString(r, v.Field(i), ch), true
			}
		}
		return nil, false
	}
	t := v.Type()
	for i := 0; i < t.NumField(); i++ {
		kField := t.Field(i)
		if kField.Anonymous {
			vField := v.Field(i)
			err, ok := run(vField)
			if err != nil {
				return err
			}
			if ok {
				return err
			}
		}
	}
	err, ok := run(v)
	if err != nil {
		return err
	}
	if ok {
		return err
	}
	return decodeString(r, reflect.New(reflect.TypeOf("")).Elem(), ch)
}
