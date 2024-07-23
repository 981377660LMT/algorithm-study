// https://geektutu.com/post/hpg-struct-alignment.html

// CPU 只从对齐的地址开始加载数据
// CPU 读取块的大小是固定的，通常为 B 的 2 的整数幂次

package main

import (
	"fmt"
	"unsafe"
)

// 在对内存特别敏感的结构体的设计上，我们可以通过调整字段的顺序，减少内存的占用。
//
// 准则：从小到大非递减排列.
func main() {
	fmt.Println(unsafe.Sizeof(demo1{})) // 8
	fmt.Println(unsafe.Sizeof(demo2{})) // 12
	fmt.Println(unsafe.Sizeof(demo3{})) // 8
	fmt.Println(unsafe.Sizeof(demo4{})) // 8
	fmt.Println(unsafe.Sizeof(demo5{})) // 4
}

type demo1 struct {
	a int8
	b int16
	c int32
}

type demo2 struct {
	a int8
	c int32
	b int16
}

type demo3 struct {
	c int32
	a int8
	b int16
}
type demo4 struct {
	c int32
	a struct{}
}

type demo5 struct {
	a struct{}
	c int32
}
