// 内存对齐
// https://geektutu.com/post/hpg-struct-alignment.html

package main

import (
	"fmt"
	"unsafe"
)

type Args struct {
	num1 int
	num2 int
}

type Flag struct {
	num1 int16
	num2 int32
}

type Flag2 struct {
	num2 int32
	num1 int16
}

type Event1 struct {
	kind bool
	num2 int32
	num1 int16
}

type Event2 struct {
	kind bool
	num1 int32
	num2 int32
}

type Event3 struct {
	kind1 int32
	kind2 int32
	kind3 int32
	kind4 int32
	time  int
}

func main() {
	fmt.Println(unsafe.Sizeof(Args{}))
	fmt.Println(unsafe.Sizeof(Flag{}))
	fmt.Println(unsafe.Sizeof(Flag2{}))
	fmt.Println(unsafe.Sizeof(Event1{}))
	fmt.Println(unsafe.Sizeof(Event2{}))
	fmt.Println(unsafe.Sizeof(Event3{}))
}
