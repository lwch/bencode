package bencode

import (
	"bytes"
	"errors"
	"fmt"
	"reflect"
)

// Encode encode data
func Encode(data interface{}) ([]byte, error) {
	var buf bytes.Buffer
	t := reflect.TypeOf(data)
	v := reflect.ValueOf(data)
	err := encode(&buf, t, v)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func encode(buf *bytes.Buffer, t reflect.Type, v reflect.Value) error {
	switch t.Kind() {
	case reflect.Bool:
		return errors.New("not supported bool value")
	case reflect.Int,
		reflect.Int8, reflect.Int16,
		reflect.Int32, reflect.Int64:
		_, err := buf.WriteString(fmt.Sprintf("i%de", v.Int()))
		return err
	case reflect.Uint,
		reflect.Uint8, reflect.Uint16,
		reflect.Uint32, reflect.Uint64:
		_, err := buf.WriteString(fmt.Sprintf("i%de", v.Uint()))
		return err
	}
	return nil
}
