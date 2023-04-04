// 遍历bits(非常快)

package main

import (
	"fmt"
	"math/bits"
	"time"
)

func main() {
	time1 := time.Now()
	for i := 0; i < 1e7; i++ {
		EnumerateBits(uint(i), func(bit int) {})
	}
	fmt.Println(time.Since(time1))
}

// 遍历每个为1的比特位
func EnumerateBits(s uint, f func(bit int)) {
	for s > 0 {
		i := bits.TrailingZeros(s)
		f(i)
		s ^= 1 << i
	}
}
