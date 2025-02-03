package main

import (
	"fmt"
)

func main() {
	f1 := func() {
		s := make([]int, 10, 12)
		s1 := s[8:]
		fmt.Printf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
	}

	f2 := func() {
		s := make([]int, 10, 12)
		s1 := s[8:10]
		fmt.Printf("s1: %v, len of s1: %d, cap of s1: %d", s1, len(s1), cap(s1))
	}

	f3 := func() {
		s := []int{0, 1, 2, 3, 4}
		s = append(s[:2], s[3:]...)
		fmt.Printf("s: %v, len of s: %d, cap of s: %d", s, len(s), cap(s))
	}

	// 1. 扩容：512*1.25 + 256*0.75 = 832
	// 2. 内存对齐：8byte * 832 = 6656 byte，对拟申请空间的 object 进行大小补齐，最终 6656 byte 会被补齐到 6784 byte 的这一档次。
	// 3. 申请内存：扩容后实际的新容量为 cap = 6784/8 = 848
	f4 := func() {
		s := make([]int, 512)
		s = append(s, 1)
		fmt.Printf("len of s: %d, cap of s: %d", len(s), cap(s)) // len of s: 513, cap of s: 848
	}

	f1()
	fmt.Println()
	f2()
	fmt.Println()
	f3()
	fmt.Println()
	f4()
}
