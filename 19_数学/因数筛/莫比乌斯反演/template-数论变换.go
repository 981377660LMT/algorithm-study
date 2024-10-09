// https://hitonanode.github.io/cplib-cpp/number/zeta_moebius_transform.hpp
// !约数/倍数 zeta/μ变换   (约数/倍数的前缀和)
// 整除関係に基づくゼータ変換・メビウス変換．
// またこれらを使用した添字 GCD・添字 LCM 畳み込み．計算量は O(N loglog N)．
// !なお，引数の vector<> について第 0 要素（f[0]）の値は無視される．

package main

import (
	"sort"
)

// 3312. 查询排序后的最大公约数
// https://leetcode.cn/problems/sorted-gcd-pair-queries/description/
func gcdValues(nums []int, queries []int64) []int {
	upper := maxs(nums...) + 1
	c := make([]int, upper)
	for _, v := range nums {
		c[v]++
	}
	MultipleZeta(c)
	for i := 1; i < upper; i++ {
		c[i] = c[i] * (c[i] - 1) / 2
	}
	MultipleMoebius(c)

	presum := make([]int, len(c))
	presum[0] = c[0]
	for i := 1; i < len(c); i++ {
		presum[i] = presum[i-1] + c[i]
	}
	res := make([]int, len(queries))
	for i, kth := range queries {
		res[i] = sort.SearchInts(presum, int(kth)+1)
	}
	return res
}

// 倍数的前缀和变换 O(N loglog N)
//
//	c2[v] = c1[v] + c1[2v] + c1[3v] + ... + c1[dv]
//
// 用于gcd相关问题.
func MultipleZeta(c1 []int) {
	n := len(c1) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := n / p; j > 0; j-- {
				c1[j] += c1[j*p] // op
			}
		}
	}
}

// 倍数的前缀和逆变换 O(N loglog N)
// 用于gcd相关问题.
func MultipleMoebius(c2 []int) {
	n := len(c2) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := 1; j*p <= n; j++ {
				c2[j] -= c2[j*p] // op
			}
		}
	}
}

// 約数的前缀和变换 O(N loglog N)
//
//	c2[v] = c1[f1] + c1[f2] + c1[f3] + ... + c1[fm]
//
// 用于lcm相关问题.
func DivisorZeta(c1 []int) {
	n := len(c1) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := 1; j*p <= n; j++ {
				c1[j*p] += c1[j] // op
			}
		}
	}
}

// 約数的前缀和逆变换 O(N loglog N).
// 用于lcm相关问题.
func DivisorMoebius(c2 []int) {
	n := len(c2) - 1
	isPrime := make([]bool, n+1)
	for i := range isPrime {
		isPrime[i] = true
	}
	for p := 2; p <= n; p++ {
		if isPrime[p] {
			for q := p * 2; q <= n; q += p {
				isPrime[q] = false
			}
			for j := n / p; j > 0; j-- {
				c2[j*p] -= c2[j] // op
			}
		}
	}
}

// gcd卷积
// ret[k] = sum(f[i] * g[j] | gcd(i, j) = k)
func GcdConv(f, g []int) []int {
	n := len(f)
	if n != len(g) {
		panic("len(f) != len(g)")
	}
	MultipleZeta(f)
	MultipleZeta(g)
	for i := range f {
		f[i] *= g[i]
	}
	MultipleMoebius(f)
	return f
}

// lcm卷积
// ret[k] = sum(f[i] * g[j] | lcm(i, j) = k)
func LcmConv(f, g []int) []int {
	n := len(f)
	if n != len(g) {
		panic("len(f) != len(g)")
	}
	DivisorZeta(f)
	DivisorZeta(g)
	for i := range f {
		f[i] *= g[i]
	}
	DivisorMoebius(f)
	return f
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
