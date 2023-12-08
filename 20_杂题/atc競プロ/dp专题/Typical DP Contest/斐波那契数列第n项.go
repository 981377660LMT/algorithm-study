package main

import (
	"fmt"
	"math/bits"
)

func main() {
	fmt.Println(KthFibonacci(0, 1000000007))  // 0
	fmt.Println(KthFibonacci(1, 1000000007))  // 1
	fmt.Println(KthFibonacci(2, 1000000007))  // 1
	fmt.Println(KthFibonacci(3, 1000000007))  // 2
	fmt.Println(KthFibonacci(10, 1000000007)) // 55
}

// 斐波那契数列第k(0-indexed)项.
func KthFibonacci(k int, mod int) int {
	if k <= 1 {
		return k
	}
	a, b := 0, 1
	i := 1 << (63 - bits.LeadingZeros64(uint64(k)) - 1)
	for ; i > 0; i >>= 1 {
		na := (a*a + b*b) % mod
		nb := (2*a + b) * b % mod
		a, b = na, nb
		if k&i != 0 {
			c := a + b
			if c >= mod {
				c -= mod
			}
			a, b = b, c
		}
	}
	return b
}
