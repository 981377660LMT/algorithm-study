package main

import (
	"fmt"
	"math"
	"sort"
)

func main() {
	EnumeratePrefixPrimeFactors(100, func(num, primeFactor int) {
		fmt.Println(num, primeFactor)
	})

	nums := []int{0, 1, 1, 1, 1, 1, 1, 1, 1, 1, 1}
	// DivisorZeta(nums)
	// DivisorMobius(nums)
	MultiplierZeta(nums)
	// MultiplierMobius(nums)
	fmt.Println(nums)
}

// SoS DP (Sum over Subsets) 数论版本.
func DivisorZeta(nums []int) {
	if nums[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(nums) - 1
	EnumeratePrefixPrimeFactors(n, func(num, primeFactor int) {
		nums[num] += nums[primeFactor] // add
	})
}

func DivisorMobius(nums []int) {
	if nums[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(nums) - 1
	EnumeratePrefixPrimeFactors(n, func(num, primeFactor int) {
		nums[num] -= nums[primeFactor] // sub
	})
}

// SoS DP (Sum over Subsets) 数论版本.
func MultiplierZeta(nums []int) {
	if nums[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(nums) - 1
	EnumeratePrefixPrimeFactors(n, func(num, primeFactor int) {
		nums[primeFactor] += nums[num] // add
	})
}

func MultiplierMobius(nums []int) {
	if nums[0] != 0 {
		panic("nums[0] != 0")
	}
	n := len(nums) - 1
	EnumeratePrefixPrimeFactors(n, func(num, primeFactor int) {
		nums[primeFactor] -= nums[num] // sub
	})
}

// 倒序遍历[2, n]内所有数的所有素因子.
func EnumeratePrefixPrimeFactors(n int, f func(num, primeFactor int)) {
	primes := P.GetPrimes(n)
	for _, p := range primes {
		for x := n / p; x >= 1; x-- {
			f(x*p, p)
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
