package bencode

import (
	"bytes"
	"strings"
	"testing"
)

func TestEncode(t *testing.T) {
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
	var buf bytes.Buffer
	err := NewEncoder(&buf).Encode(obj)
	if err != nil {
		t.Fatalf("FATAL: encode %v", err)
	}
	data := buf.String()
	if data[0] != 'd' {
		t.Fatalf("unexpected first char: %c", data[0])
	}
	if data[len(data)-1] != 'e' {
		t.Fatalf("unexpected end char: %c", data[len(data)-1])
	}
	if !strings.Contains(data, "1:ad2:id20:abcdefghij0123456789e") {
		t.Fatal("value of key a not found")
	}
	if !strings.Contains(data, "1:q4:ping") {
		t.Fatal("value of key q not found")
	}
	if !strings.Contains(data, "1:t2:aa") {
		t.Fatal("value of key t not found")
	}
	if !strings.Contains(data, "1:y1:q") {
		t.Fatal("value of key y not found")
	}
}
