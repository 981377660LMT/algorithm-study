package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fmt.Println(CeilLog32(0))  // 0
	fmt.Println(CeilLog32(13)) // 4
	fmt.Println(CeilLog32(16)) // 4

	fmt.Println(FloorLog32(13)) // 3
	fmt.Println(FloorLog32(16)) // 4
}

func CeilLog32(x uint32) int {
	if x <= 0 {
		return 0
	}
	return 32 - bits.LeadingZeros32(x-1)
}

func FloorLog32(x uint32) int {
	if x <= 0 {
		panic("IllegalArgumentException")
	}
	return 31 - bits.LeadingZeros32(x)
}

func CeilLog64(x uint64) int {
	if x <= 0 {
		return 0
	}
	return 64 - bits.LeadingZeros64(x-1)
}

func FloorLog64(x uint64) int {
	if x <= 0 {
		panic("IllegalArgumentException")
	}
	return 63 - bits.LeadingZeros64(x)
}

// CachedLog2.
const BITS int32 = 16
const LIMIT int32 = 1 << BITS

var CACHE [LIMIT]int32

func init() {
	var b int32
	for i := int32(0); i < LIMIT; i++ {
		for 1<<(b+1) <= i {
			b++
		}
		CACHE[i] = b
	}
}

// CacheLog2 returns the floor of the base 2 logarithm of x.
func FloorLog(x int32) int32 {
	if x < LIMIT {
		return CACHE[x]
	}
	return BITS + CACHE[x>>BITS]
}
