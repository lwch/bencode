package bencode

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"testing"
)

func TestEncode(t *testing.T) {
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
	copy(obj.A.ID[:], []byte("abcd"))
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(obj)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Println(hex.Dump(buf.Bytes()))
}
