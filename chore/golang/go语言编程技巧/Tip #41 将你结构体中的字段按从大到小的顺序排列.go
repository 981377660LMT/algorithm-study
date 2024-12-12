package main

import (
	"fmt"
	"unsafe"
)

func main() {
	a := A{}
	fmt.Println("A size:", unsafe.Sizeof(a))
	aa := OptimizedA{}
	fmt.Println("OptimizedA size:", unsafe.Sizeof(aa))
}

// Sort your fields in your struct from largest to smallest.

// 32 bytes
type A struct {
	a byte  // 1-byte alignment
	b int32 // 4-byte alignment
	c byte  // 1-byte alignment
	d int64 // 8-byte alignment
	e byte  // 1-byte alignment
}

// 16 bytes
type OptimizedA struct {
	d int64 // 8-byte alignment
	b int32 // 4-byte alignment
	a byte  // 1-byte alignment
	c byte  // 1-byte alignment
	e byte  // 1-byte alignment
}
