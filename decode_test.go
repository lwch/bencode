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
	var intf interface{}
	err = Decode(str, &intf)
	if err != nil {
		t.Fatalf("FATAL: decode number to interface: %v", err)
	}
	if intf != 1 {
		t.Fatalf("unexpected number value of interface: %d", i)
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
	var intf interface{}
	err = Decode(str, &intf)
	if err != nil {
		t.Fatalf("FATAL: decode bytes to interface: %v", err)
	}
	if intf.(string) != s {
		t.Fatalf("unexpected string value to interface: %s", string(bs[:]))
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
	var intf interface{}
	err = Decode(str, &intf)
	if err != nil {
		t.Fatalf("FATAL: decode map to interface: %v", err)
	}
	m2 := intf.(map[string]interface{})
	if len(m2) != 2 {
		t.Fatalf("unexpected map size: %d", len(m2))
	}
	if m2["a"] != 1 {
		t.Fatalf("unexpected key a of map to interface: %d", m2["a"])
	}
	if m2["b"] != "abc" {
		t.Fatalf("unexpected key b of map to interface: %s", m2["b"])
	}
}

func TestDecodeListString(t *testing.T) {
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

	var sbs [][]byte
	err = Decode(str, &sbs)
	if err != nil {
		t.Fatalf("FATAL: decode bytes slice: %v", err)
	}
	if len(sbs) != 3 {
		t.Fatalf("unexpected slice bytes size: %d", len(sbs))
	}
	if len(sbs[0]) != 1 || sbs[0][0] != 'a' {
		t.Fatalf("unexpected slice bytes value of 0: %s", string(sbs[0]))
	}
	if len(sbs[1]) != 1 || sbs[1][0] != 'b' {
		t.Fatalf("unexpected slice bytes value of 1: %s", string(sbs[1]))
	}
	if len(sbs[2]) != 1 || sbs[2][0] != 'c' {
		t.Fatalf("unexpected slice bytes value of 2: %s", string(sbs[2]))
	}
	var sba [2][1]byte
	err = Decode(str, &sba)
	if err != nil {
		t.Fatalf("FATAL: decode bytes array: %v", err)
	}
	if sba[0][0] != 'a' {
		t.Fatalf("unexpected array byte array value of 0: %s", string(sba[0][:]))
	}
	if sba[1][0] != 'b' {
		t.Fatalf("unexpected array byte array value of 1: %s", string(sbs[0]))
	}
}

func TestDecodeListInt(t *testing.T) {
	str := []byte("li1ei2ei3ee")
	var si []int
	err := Decode(str, &si)
	if err != nil {
		t.Fatalf("FATAL: decode int slice: %v", err)
	}
	if len(si) != 3 {
		t.Fatalf("unexpected number slice size: %d", len(si))
	}
	for i := 0; i < len(si); i++ {
		if si[i] != i+1 {
			t.Fatalf("unexpected number slice value of %d: %d", i, si[i])
		}
	}
	var ai [2]int
	err = Decode(str, &ai)
	if err != nil {
		t.Fatalf("FATAL: decode int array: %v", err)
	}
	for i := 0; i < len(ai); i++ {
		if ai[i] != i+1 {
			t.Fatalf("unexpected number array value of %d: %d", i, ai[i])
		}
	}

	var sintf []interface{}
	err = Decode(str, &sintf)
	if err != nil {
		t.Fatalf("FATAL: decode int interface slice: %v", err)
	}
	if len(sintf) != 3 {
		t.Fatalf("unexpected interface slice size: %d", len(sintf))
	}
	for i := 0; i < len(sintf); i++ {
		if sintf[i] != i+1 {
			t.Fatalf("unexpected int interface slice value of %d: %d", i, ai[i])
		}
	}
	var intf interface{}
	err = Decode(str, &intf)
	if err != nil {
		t.Fatalf("FATAL: decode int interface slice: %v", err)
	}
	sintf = intf.([]interface{})
	if len(sintf) != 3 {
		t.Fatalf("unexpected interface slice size: %d", len(sintf))
	}
	for i := 0; i < len(sintf); i++ {
		if sintf[i] != i+1 {
			t.Fatalf("unexpected interface value of %d: %d", i, ai[i])
		}
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

func TestDecodeAnnounce(t *testing.T) {
	data := []byte{
		0x64, 0x31, 0x3a, 0x61, 0x64, 0x32, 0x3a, 0x69, 0x64, 0x32, 0x30, 0x3a, 0xf5, 0xe1, 0x44, 0x56,
		0x65, 0x97, 0x4e, 0xf2, 0xa0, 0xfb, 0x28, 0x95, 0x66, 0xd8, 0x41, 0x7f, 0x48, 0x90, 0x40, 0x4b,
		0x39, 0x3a, 0x69, 0x6e, 0x66, 0x6f, 0x5f, 0x68, 0x61, 0x73, 0x68, 0x32, 0x30, 0x3a, 0x2d, 0x5b,
		0xb1, 0x38, 0x02, 0xc4, 0xe4, 0xcd, 0xd2, 0xba, 0x4d, 0x0b, 0x9e, 0x08, 0x61, 0x4b, 0xef, 0xed,
		0xf6, 0x2c, 0x34, 0x3a, 0x70, 0x6f, 0x72, 0x74, 0x69, 0x38, 0x30, 0x38, 0x37, 0x65, 0x35, 0x3a,
		0x74, 0x6f, 0x6b, 0x65, 0x6e, 0x30, 0x3a, 0x65, 0x31, 0x3a, 0x71, 0x31, 0x33, 0x3a, 0x61, 0x6e,
		0x6e, 0x6f, 0x75, 0x6e, 0x63, 0x65, 0x5f, 0x70, 0x65, 0x65, 0x72, 0x31, 0x3a, 0x74, 0x38, 0x3a,
		0x5b, 0x33, 0x4a, 0xfa, 0xee, 0x1d, 0x7a, 0x47, 0x31, 0x3a, 0x79, 0x31, 0x3a, 0x71, 0x65,
	}
	type hdr struct {
		Transaction string `bencode:"t"`
		Type        string `bencode:"y"`
	}
	type announceRequest struct {
		hdr
		Response struct {
			ID      [20]byte `bencode:"id"`
			Hash    [20]byte `bencode:"info_hash"`
			Implied int      `bencode:"implied_port"`
			Port    uint16   `bencode:"port"`
			Token   string   `bencode:"token"`
		} `bencode:"a"`
	}
	var req announceRequest
	err := Decode(data, &req)
	if err != nil {
		t.Fatalf("FATAL: decode announce request: %v", err)
	}
}
