package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	nums := []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	// DivisorZeta(nums)
	// DivisorMobius(nums)
	MultiplierZeta(nums)
	// MultiplierMobius(nums)
	fmt.Println(nums)
}

// !对于给定的整数数组 nums，快速回答“第 k 小的数对 (nums[i], nums[j]) 的最大公约数 (GCD) 是多少
func gcdValues(nums []int, queries []int64) []int {
	upper := maxs(nums...) + 1
	c := make([]int, upper)
	for _, v := range nums {
		c[v]++
	}
	MultiplierZeta(c)
	for i := 1; i < upper; i++ {
		c[i] = c[i] * (c[i] - 1) / 2
	}
	MultiplierMobius(c)

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

// SoS DP (Sum over Subsets) 数论版本.
func DivisorZeta(c1 []int) {
	if c1[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(c1) - 1
	primes := P.GetPrimes(n)
	for _, p := range primes {
		for x := 1; x < n/p+1; x++ {
			c1[x*p] += c1[x]
		}
	}
}

func DivisorMobius(c2 []int) {
	if c2[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(c2) - 1
	primes := P.GetPrimes(n)
	for _, p := range primes {
		for x := n / p; x > 0; x-- {
			c2[x*p] -= c2[x]
		}
	}
}

// SoS DP (Sum over Subsets) 数论版本.
func MultiplierZeta(c1 []int) {
	if c1[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(c1) - 1
	primes := P.GetPrimes(n)
	for _, p := range primes {
		for x := n / p; x > 0; x-- {
			c1[x] += c1[x*p]
		}
	}
}

func MultiplierMobius(c2 []int) {
	if c2[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(c2) - 1
	primes := P.GetPrimes(n)
	for _, p := range primes {
		for x := 1; x < n/p+1; x++ {
			c2[x] -= c2[x*p]
		}
	}
}

const S int = 32768

var P = NPT(1e5 + 10)

type PT struct {
	done   int
	primes []int
	sieve  []bool
}

func NPT(limit int) *PT {
	res := &PT{
		done:   2,
		primes: []int{2},
		sieve:  make([]bool, S+1),
	}
	res.expand(limit + 1)
	return res
}

// 返回小于等于limit的所有素数.
func (table *PT) GetPrimes(limit int) []int {
	limit++
	table.expand(limit)
	k := sort.Search(len(table.primes), func(i int) bool { return table.primes[i] >= limit })
	return table.primes[:k]
}

func (table *PT) expand(limit int) {
	if table.done < limit {
		table.done = limit
		R := limit / 2
		for i := range table.sieve {
			table.sieve[i] = false
		}
		table.primes = make([]int, 0, int((float64(limit)/math.Log(float64(limit)))*1.1))
		table.primes = append(table.primes, 2)
		cp := [][2]int{}
		for i := 3; i <= S; i += 2 {
			if !table.sieve[i] {
				cp = append(cp, [2]int{i, i * i / 2})
				for j := i * i; j <= S; j += 2 * i {
					table.sieve[j] = true
				}
			}
		}
		for L := 1; L <= R; L += S {
			block := [S]bool{}
			for i := range cp {
				pair := &cp[i]
				p, idx := pair[0], &pair[1]
				for j := *idx; j < S+L; {
					block[j-L] = true
					j += p
					*idx = j
				}
			}
			for i := 0; i < min(S, R-L); i++ {
				if !block[i] {
					table.primes = append(table.primes, (L+i)*2+1)
				}
			}
		}
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
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
