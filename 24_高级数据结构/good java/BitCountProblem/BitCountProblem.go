// BitCountProblem

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fmt.Println(CountZeros(10))
	fmt.Println(CountOnes(10))
	fmt.Println(CountLength(10))

	fmt.Println(CountZerosRange(3, 10))
	fmt.Println(CountOnesRange(3, 10))
	fmt.Println(CountLengthRange(3, 10))
}

// [0, n] 之间的数中二进制表示中0的个数之和.
func CountZeros(n int) int {
	return CountLength(n) - CountOnes(n)
}

// [0, n] 之间的数中二进制表示中1的个数之和.
func CountOnes(n int) int {
	if n == 0 {
		return 0
	}
	res := 0
	len := floorLog64(uint64(n))
	for i := 0; i <= len; i++ {
		bit := 1 << i
		remain := n
		block := remain / (2 * bit)
		res += block * bit
		remain %= (2 * bit)
		res += max(0, remain-bit+1)
	}
	return res
}

// [0, n] 之间的数中二进制表示位数之和.
func CountLength(n int) int {
	if n == 0 {
		return 1
	}
	res := 0
	len := floorLog64(uint64(n))
	for i := 0; i <= len; i++ {
		l := 1 << i
		r := min(l*2-1, n)
		res += (i + 1) * (r - l + 1)
	}
	return res
}

// [left, right] 之间的数中二进制表示中0的个数之和.
func CountZerosRange(left, right int) int {
	res := CountZeros(right)
	if left > 0 {
		res -= CountZeros(left - 1)
	}
	return res
}

// [left, right] 之间的数中二进制表示中1的个数之和.
func CountOnesRange(left, right int) int {
	res := CountOnes(right)
	if left > 0 {
		res -= CountOnes(left - 1)
	}
	return res
}

// [left, right] 之间的数中二进制表示位数之和.
func CountLengthRange(left, right int) int {
	res := CountLength(right)
	if left > 0 {
		res -= CountLength(left - 1)
	}
	return res
}

func floorLog64(x uint64) int {
	if x <= 0 {
		panic("IllegalArgumentException")
	}
	return 63 - bits.LeadingZeros64(x)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
