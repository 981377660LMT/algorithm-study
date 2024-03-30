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
