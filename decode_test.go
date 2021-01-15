package bencode

import (
	"fmt"
	"testing"
)

func TestDecode(t *testing.T) {
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
	fmt.Println(obj)
}
