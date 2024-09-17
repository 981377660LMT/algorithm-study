// 遍历bits(非常快)

package main

import (
	"fmt"
	"math/bits"
	"time"
)

func main() {
	time1 := time.Now()
	for i := 0; i < 1e8; i++ {
		EnumerateBits32(uint32(i), func(bit int) {})
	}
	fmt.Println(time.Since(time1))
}

// 遍历每个为1的比特位
func EnumerateBits32(s uint32, f func(bit int)) {
	for s > 0 {
		i := bits.TrailingZeros32(s)
		f(i)
		s ^= 1 << i
	}
}

func EnumerateBits64(s uint64, f func(bit int)) {
	for s > 0 {
		i := bits.TrailingZeros64(s)
		f(i)
		s ^= 1 << i
	}
}
