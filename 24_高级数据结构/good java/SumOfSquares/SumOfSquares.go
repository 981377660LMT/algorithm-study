// 判断正整数n是否能拆分成不同的平方数之和
// 判断正整数n是否能拆分成1/2/3/4个平方数之和
// SquareSum

package main

import (
	"fmt"
	"math/big"
	"math/rand"
)

func main() {
	fmt.Println(IsSumOfDistinctSquares(5))
	fmt.Println(IsSumOfOneSquare(9))
	fmt.Println(IsSumOfTwoSquare(9))
	fmt.Println(IsSumOfThreeSquare(9))
	fmt.Println(IsSumOfFourSquare(9))
}

// https://leetcode.cn/problems/sum-of-square-numbers/description/
func judgeSquareSum(c int) bool {
	return IsSumOfTwoSquare(c)
}

var CAN_NOT_BE_REPRESENT = [...]uint8{
	2, 3, 6, 7, 8, 11, 12, 15, 18, 19, 22,
	23, 24, 27, 28, 31, 32, 33, 43, 44, 47,
	48, 60, 67, 72, 76, 92, 96, 108, 112, 128,
}

func IsSumOfDistinctSquares(n int) bool {
	if n < 0 {
		return false
	}
	if n <= 128 {
		n8 := uint8(n)
		for _, v := range CAN_NOT_BE_REPRESENT {
			if v == n8 {
				return false
			}
		}
	}
	return true
}

func IsSumOfOneSquare(n int) bool {
	if n == 0 {
		return true
	}
	isqrt := floorSqrt(n)
	return isqrt*isqrt == n
}

func IsSumOfTwoSquare(n int) bool {
	if n == 0 {
		return true
	}
	set := findAllFactors(n)
	for x := range set {
		pow := 0
		y := n
		for y%x == 0 {
			y /= x
			pow++
		}
		if x%4 == 3 && pow%2 == 1 {
			return false
		}
	}
	return true
}

func IsSumOfThreeSquare(n int) bool {
	if n == 0 {
		return true
	}
	for {
		if n%8 == 7 {
			return false
		}
		if n%4 != 0 {
			break
		}
		n /= 4
	}
	return true
}

func IsSumOfFourSquare(n int) bool {
	return true
}

func floorSqrt(x int) int {
	lo := 0
	hi := int(3e9)
	for lo < hi {
		mid := (lo + hi + 1) >> 1
		if mid*mid <= x {
			lo = mid
		} else {
			hi = mid - 1
		}
	}
	return lo
}

func findFactor(n int) int {
	if n == 1 {
		return 1
	}
	if n == 4 {
		return 2
	}
	if isPrimeMillerRabin(n) {
		return n
	}

	gcd := func(a, b int) int {
		for a != 0 {
			a, b = b%a, a
		}
		return b
	}

	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	mul := func(a, b int) (res int) {
		for ; b > 0; b >>= 1 {
			if b&1 == 1 {
				res = (res + a) % n
			}
			a = (a + a) % n
		}
		return
	}

	for {
		c := 1 + rand.Intn(n-1)
		f := func(x int) int { return (mul(x, x) + c) % n }
		for t, r := f(0), f(f(0)); t != r; t, r = f(t), f(f(r)) {
			if d := gcd(abs(t-r), n); d > 1 {
				return d
			}
		}
	}
}

func isPrimeMillerRabin(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

var SMALL_PRIMES = [...]int{2, 3, 5, 7, 11, 13, 17, 19}

func findAllFactors(n int) map[int]struct{} {
	if n == 0 {
		return map[int]struct{}{}
	}
	set := map[int]struct{}{}
	for _, p := range SMALL_PRIMES {
		if n%p == 0 {
			set[p] = struct{}{}
			for n%p == 0 {
				n /= p
			}
		}
	}
	findAllFactors2(set, n)
	return set
}

func findAllFactors2(set map[int]struct{}, n int) {
	if n == 1 {
		return
	}
	f := findFactor(n)
	if f == n {
		set[n] = struct{}{}
		return
	}
	other := n / f
	findAllFactors2(set, f)
	findAllFactors2(set, other)
}
