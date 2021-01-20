package bencode

import (
	"bytes"
	"testing"
)

func TestDecodeNumber(t *testing.T) {
	str := []byte("i1e")
	var i int
	err := Decode(str, &i)
	if err != nil {
		t.Fatalf("FATAL: decode number: %v", err)
	}
	if i != 1 {
		t.Fatalf("unexpected number value: %d", i)
	}
}

func TestDecodeString(t *testing.T) {
	str := []byte("3:abc")
	var s string
	err := Decode(str, &s)
	if err != nil {
		t.Fatalf("FATAL: decode string: %v", err)
	}
	if s != "abc" {
		t.Fatalf("unexpected string value: %s", s)
	}
	var bs [3]byte
	err = Decode(str, &bs)
	if err != nil {
		t.Fatalf("FATAL: decode bytes: %v", err)
	}
	if !bytes.Equal(bs[:], []byte(s)) {
		t.Fatalf("unexpected bytes value: %s", string(bs[:]))
	}
}

func TestDecodeMap(t *testing.T) {
	str := []byte("d1:ai1e1:b3:abce")
	var m map[string]interface{}
	err := Decode(str, &m)
	if err != nil {
		t.Fatalf("FATAL: decode map: %v", err)
	}
	if len(m) != 2 {
		t.Fatalf("unexpected map size: %d", len(m))
	}
	if m["a"] != 1 {
		t.Fatalf("unexpected key a of map: %d", m["a"])
	}
	if m["b"] != "abc" {
		t.Fatalf("unexpected key b of map: %s", m["b"])
	}
}

func TestDecodeList(t *testing.T) {
	str := []byte("l1:a1:b1:ce")
	var ss []string
	err := Decode(str, &ss)
	if err != nil {
		t.Fatalf("FATAL: decode string slice: %v", err)
	}
	if len(ss) != 3 {
		t.Fatalf("unexpected slice size: %d", len(ss))
	}
	if ss[0] != "a" {
		t.Fatalf("unexpected slice value of 0: %s", ss[0])
	}
	if ss[1] != "b" {
		t.Fatalf("unexpected slice value of 1: %s", ss[1])
	}
	if ss[2] != "c" {
		t.Fatalf("unexpected slice value of 2: %s", ss[2])
	}
	var as [2]string
	err = Decode(str, &as)
	if err != nil {
		t.Fatalf("FATAL: decode string array: %v", err)
	}
	if as[0] != "a" {
		t.Fatalf("unexpected slice value of 0: %s", ss[0])
	}
	if as[1] != "b" {
		t.Fatalf("unexpected slice value of 1: %s", ss[1])
	}
}

func TestDecodeStruct(t *testing.T) {
	// example of http://www.bittorrent.org/beps/bep_0005.html
	str := []byte("d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t2:aa1:y1:qe")
	var obj struct {
		T string `bencode:"t"`
		Y string `bencode:"y"`
		Q string `bencode:"q"`
		A struct {
			ID [20]byte `bencode:"id"`
		} `bencode:"a"`
	}
	err := Decode(str, &obj)
	if err != nil {
		t.Fatalf("FATAL: decode struct: %v", err)
	}
	if obj.T != "aa" {
		t.Fatalf("unexpected value of t: %s", obj.T)
	}
	if obj.Y != "q" {
		t.Fatalf("unexpected value of y: %s", obj.Y)
	}
	if obj.Q != "ping" {
		t.Fatalf("unexpected value of q: %s", obj.Q)
	}
	if bytes.Compare(obj.A.ID[:], []byte("abcdefghij0123456789")) != 0 {
		t.Fatalf("unexpected value of a.id: %s", string(obj.A.ID[:]))
	}
}

func TestDecodeInherit(t *testing.T) {
	type hdr struct {
		Transaction string `bencode:"t"`
		Type        string `bencode:"y"`
	}
	type pingResponse struct {
		hdr
		Response struct {
			ID [20]byte `bencode:"id"`
		} `bencode:"r"`
	}
	var r pingResponse
	data := []byte{
		0x64, 0x32, 0x3a, 0x69, 0x70, 0x36, 0x3a, 0x74, 0x55, 0x1d, 0xa7, 0xd3, 0x73, 0x31, 0x3a, 0x72,
		0x64, 0x32, 0x3a, 0x69, 0x64, 0x32, 0x30, 0x3a, 0xf7, 0x45, 0xed, 0xe5, 0x01, 0x34, 0x7f, 0x9a,
		0x4b, 0x11, 0x3b, 0xe3, 0xc6, 0x36, 0xb5, 0x97, 0xd8, 0x49, 0xe7, 0x6c, 0x31, 0x3a, 0x70, 0x69,
		0x35, 0x34, 0x31, 0x33, 0x31, 0x65, 0x65, 0x31, 0x3a, 0x74, 0x33, 0x32, 0x3a, 0x32, 0x31, 0x30,
		0x31, 0x31, 0x38, 0x62, 0x63, 0x37, 0x31, 0x39, 0x32, 0x32, 0x36, 0x38, 0x33, 0x35, 0x64, 0x30,
		0x30, 0x66, 0x33, 0x31, 0x63, 0x62, 0x30, 0x35, 0x66, 0x63, 0x31, 0x31, 0x61, 0x31, 0x3a, 0x76,
		0x34, 0x3a, 0x4c, 0x54, 0x01, 0x01, 0x31, 0x3a, 0x79, 0x31, 0x3a, 0x72, 0x65,
	}
	err := Decode(data, &r)
	if err != nil {
		t.Fatalf("FATAL: decode inherit: %v", err)
	}
	if r.Transaction != "210118bc719226835d00f31cb05fc11a" {
		t.Fatalf("unexpected value of t: %s", r.Transaction)
	}
	if r.Type != "r" {
		t.Fatalf("unexpected value of y: %s", r.Type)
	}
	id := []byte{
		0xf7, 0x45, 0xed, 0xe5,
		0x01, 0x34, 0x7f, 0x9a,
		0x4b, 0x11, 0x3b, 0xe3,
		0xc6, 0x36, 0xb5, 0x97,
		0xd8, 0x49, 0xe7, 0x6c,
	}
	if !bytes.Equal(r.Response.ID[:], id) {
		t.Fatalf("unexpected value of r.id: %s", string(r.Response.ID[:]))
	}
}
