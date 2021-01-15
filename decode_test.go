package bencode

import (
	"bytes"
	"testing"
)

func TestDecode(t *testing.T) {
	// example of http://www.bittorrent.org/beps/bep_0005.html
	str := "d1:ad2:id20:abcdefghij0123456789e1:q4:ping1:t2:aa1:y1:qe"
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
		t.Fatalf("FATAL: decode: %v", err)
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
