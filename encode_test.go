package bencode

import (
	"bytes"
	"fmt"
	"testing"
)

func TestEncodeNumber(t *testing.T) {
	run := func(i interface{}, name string) {
		data, err := Encode(i)
		if err != nil {
			t.Fatalf("FATAL: encode %s: %v", name, err)
		}
		if bytes.Compare(data, []byte(fmt.Sprintf("i%de", i))) != 0 {
			t.Fatalf("unexpected value: %s", string(data))
		}
	}
	i := 1
	run(i, "int")
	i8 := int8(1)
	run(i8, "int8")
	i16 := int16(1)
	run(i16, "int16")
	i32 := int32(1)
	run(i32, "int32")
	i64 := int64(1)
	run(i64, "int64")
	u := uint(1)
	run(u, "uint")
	u8 := uint8(1)
	run(u8, "uint8")
	u16 := uint16(1)
	run(u16, "uint16")
	u32 := uint32(1)
	run(u32, "uint32")
	u64 := uint64(1)
	run(u64, "uint64")
	uneg := -1
	run(uneg, "negative")
}

func TestEncodeString(t *testing.T) {
	str := "abc"
	data, err := Encode(str)
	if err != nil {
		t.Fatalf("FATAL: encode string: %v", err)
	}
	if bytes.Compare(data, []byte(fmt.Sprintf("%d:%s", len(str), str))) != 0 {
		t.Fatalf("unexpected value: %s", string(data))
	}

	bs := []byte(str)
	data, err = Encode(str)
	if err != nil {
		t.Fatalf("FATAL: encode bytes: %v", err)
	}
	if bytes.Compare(data, []byte(fmt.Sprintf("%d:%s", len(bs), string(bs)))) != 0 {
		t.Fatalf("unexpected value: %s", string(data))
	}

	var ba [3]byte
	copy(ba[:], str)
	data, err = Encode(ba)
	if err != nil {
		t.Fatalf("FATAL: encode byte array: %v", err)
	}
	if bytes.Compare(data, []byte(fmt.Sprintf("%d:%s", len(ba), string(ba[:])))) != 0 {
		t.Fatalf("unexpected value: %s", string(data))
	}
}

func TestEncodeMap(t *testing.T) {
	mp := map[string]interface{}{
		"a": 1,
		"b": "abc",
	}
	data, err := Encode(mp)
	if err != nil {
		t.Fatalf("FATAL: encode map: %v", err)
	}
	if data[0] != 'd' {
		t.Fatalf("unexpected first char: %c", data[0])
	}
	if data[len(data)-1] != 'e' {
		t.Fatalf("unexpected end char: %c", data[len(data)-1])
	}
	if !bytes.Contains(data, []byte("1:ai1e")) {
		t.Fatal("value of key a not found")
	}
	if !bytes.Contains(data, []byte("1:b3:abc")) {
		t.Fatal("value of key b not found")
	}
}

func TestEncodeStruct(t *testing.T) {
	// example of http://www.bittorrent.org/beps/bep_0005.html
	var obj struct {
		T string `bencode:"t"`
		Y string `bencode:"y"`
		Q string `bencode:"q"`
		A struct {
			ID [20]byte `bencode:"id"`
		} `bencode:"a"`
	}
	obj.T = "aa"
	obj.Y = "q"
	obj.Q = "ping"
	copy(obj.A.ID[:], []byte("abcdefghij0123456789"))
	data, err := Encode(obj)
	if err != nil {
		t.Fatalf("FATAL: encode struct: %v", err)
	}
	if data[0] != 'd' {
		t.Fatalf("unexpected first char: %c", data[0])
	}
	if data[len(data)-1] != 'e' {
		t.Fatalf("unexpected end char: %c", data[len(data)-1])
	}
	if !bytes.Contains(data, []byte("1:ad2:id20:abcdefghij0123456789e")) {
		t.Fatal("value of key a not found")
	}
	if !bytes.Contains(data, []byte("1:q4:ping")) {
		t.Fatal("value of key q not found")
	}
	if !bytes.Contains(data, []byte("1:t2:aa")) {
		t.Fatal("value of key t not found")
	}
	if !bytes.Contains(data, []byte("1:y1:q")) {
		t.Fatal("value of key y not found")
	}
}

func TestEncodeInterface(t *testing.T) {
	var intf interface{}
	intf = 1
	data, err := Encode(intf)
	if err != nil {
		t.Fatalf("FATAL: encode interface int: %v", err)
	}
	if bytes.Compare(data, []byte("i1e")) != 0 {
		t.Fatalf("unexpected value: %s", string(data))
	}

	type common struct {
		Transaction string `bencode:"t"`
		Type        string `bencode:"y"`
	}
	type query struct {
		Action string      `bencode:"q"`
		Data   interface{} `bencode:"a"`
	}
	type pingRequest struct {
		common
		query
	}
	var s pingRequest
	s.common = common{
		Transaction: "aa",
		Type:        "q",
	}
	var id [20]byte
	copy(id[:], "abcdefghij0123456789")
	s.query = query{
		Action: "ping",
		Data: map[string][20]byte{
			"id": id,
		},
	}
	data, err = Encode(s)
	if err != nil {
		t.Fatalf("FATAL: encode interface map: %v", err)
	}
	if data[0] != 'd' {
		t.Fatalf("unexpected first char: %c", data[0])
	}
	if data[len(data)-1] != 'e' {
		t.Fatalf("unexpected end char: %c", data[len(data)-1])
	}
	if !bytes.Contains(data, []byte("1:ad2:id20:abcdefghij0123456789e")) {
		t.Fatal("value of key a not found")
	}
	if !bytes.Contains(data, []byte("1:q4:ping")) {
		t.Fatal("value of key q not found")
	}
	if !bytes.Contains(data, []byte("1:t2:aa")) {
		t.Fatal("value of key t not found")
	}
	if !bytes.Contains(data, []byte("1:y1:q")) {
		t.Fatal("value of key y not found")
	}
}
