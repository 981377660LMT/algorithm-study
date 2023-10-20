package main

import "math"

func RangeSum(start, end int) int {
	if start >= end {
		return 0
	}
	return ((end - start) * (start + end - 1)) / 2
}

// 区间平方和.
func RangeSquareSum(start, end int) int {
	if start >= end {
		return 0
	}
	tmp1 := (end * (end - 1) * (2*end - 1)) / 6
	tmp2 := (start * (start - 1) * (2*start - 1)) / 6
	return tmp1 - tmp2
}

// 区间立方和.
func RangeCubeSum(start, end int) int {
	if start >= end {
		return 0
	}
	tmp1 := (end * (end - 1)) / 2
	tmp2 := (start * (start - 1)) / 2
	return tmp1*tmp1 - tmp2*tmp2
}

// 区间异或和.
func RangeXorSum(start, end int) int {
	if start >= end {
		return 0
	}
	preXor := func(upper int) int {
		mod := upper % 4
		if mod == 0 {
			return upper
		}
		if mod == 1 {
			return 1
		}
		if mod == 2 {
			return upper + 1
		}
		return 0
	}
	return preXor(end-1) ^ preXor(start-1)
}

// 区间以k为底的幂和.
func RangePowKSum(start, end, k int, mod int) int {
	if start >= end {
		return 0
	}
	if mod == 1 {
		return 0
	}

	cal := func(n int) int {
		sum_, p := 1, k
		start := int(math.Log2(float64(n))) - 1
		for d := start; d >= 0; d-- {
			sum_ *= p + 1
			p *= p
			if (n>>d)&1 == 1 {
				sum_ += p
				p *= k
			}
			sum_ %= mod
			p %= mod
		}
		return sum_
	}

	return cal(end) - cal(start)
}

func Pow(base int, exp int, mod int) int {
	base %= mod
	res := 1
	for exp != 0 {
		if exp&1 == 1 {
			res = res * base % mod
		}
		base = base * base % mod
		exp >>= 1
	}
	return res
}

// 区间以2为底的幂和.
func RangePow2Sum(start, end int, mod int) int {
	if start >= end {
		return 0
	}
	res := (Pow(2, end, mod) - Pow(2, start, mod)) % mod
	if res < 0 {
		res += mod
	}
	return res
}
