// BinaryGcd
// 二进制gcd

package main

import (
	"fmt"
	"math/bits"
	"math/rand"
	"time"
)

func main() {
	fmt.Println(binaryGcd(10, 5)) // 0
	fmt.Println(naiveGcd(10, 5))  // 0

	// test perf
	nums := [1e7]int{}
	denos := [1e7]int{}
	for i := 0; i < 1e7; i++ {
		nums[i] = rand.Intn(1e9) - 8e5
		denos[i] = rand.Intn(1e9) + 1
	}

	time1 := time.Now()
	for i := 0; i < 1e7; i++ {
		naiveGcd(nums[i], denos[i])
	}
	fmt.Println(time.Since(time1)) // 2.069395s

	time1 = time.Now()
	for i := 0; i < 1e7; i++ {
		binaryGcd(nums[i], denos[i])
	}
	fmt.Println(time.Since(time1)) // 707.0429ms
}

// binary binaryGcd
func binaryGcd(a, b int) int {
	// 取绝对值
	x, y := a, b
	if x < 0 {
		x = -x
	}
	if y < 0 {
		y = -y
	}

	if x == 0 || y == 0 {
		return x + y
	}
	n := bits.TrailingZeros(uint(x))
	m := bits.TrailingZeros(uint(y))
	x >>= n
	y >>= m
	for x != y {
		d := bits.TrailingZeros(uint(x - y))
		f := x > y
		var c int
		if f {
			c = x
		} else {
			c = y
		}
		if !f {
			y = x
		}
		x = (c - y) >> d
	}

	return x << min(n, m)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func naiveGcd(a, b int) int {
	if a < 0 {
		a = -a
	}
	if b < 0 {
		b = -b
	}

	for a != 0 {
		a, b = b%a, a
	}
	return b
}
