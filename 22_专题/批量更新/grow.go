package main

import (
	"fmt"
	"slices"
)

// golang 扩容
func main() {
	var arr []int
	var caps []int
	for i := 0; i < int(1e5); i++ {
		arr = append(arr, i)
		caps = append(caps, cap(arr))
	}

	caps = slices.Compact(caps)
	fmt.Println(caps)
	//
	// [1 2 4 8 16 32 64 128 256 512 848 1280 1792 2560 3408 5120 7168 9216 12288 16384 21504 27648 34816 44032 55296 69632 88064 110592]
}
