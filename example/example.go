package example

import (
	"fmt"

	"github.com/lwch/bencode"
)

func test() {
	type t struct {
		T string `bencode:"t"`
		Y string `bencode:"y"`
		Q string `bencode:"q"`
		A struct {
			ID [20]byte `bencode:"id"`
		} `bencode:"a"`
	}
	var obj t
	obj.T = "aa"
	obj.Y = "q"
	obj.Q = "ping"
	copy(obj.A.ID[:], []byte("abcdefghij0123456789"))
	data, err := bencode.Encode(obj)
	if err != nil {
		fmt.Printf("FATAL: encode: %v\n", err)
		return
	}
	var dec t
	err = bencode.Decode(data, &dec)
	if err != nil {
		fmt.Printf("FATAL: decode: %v\n", err)
		return
	}
}
