// EnumerateRangePrimeFactors/EnumerateIntervalPrimeFactors/EnumeratePrimeFactorsInterval
//
// 区间筛/区间质因数分解

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	abc227g()
}

func demo() {
	table := NewPrimeTable(1e2)
	table.EnumerateRangePrimeFactors(100, 110, func(num, primeFactor int) {
		fmt.Println(num, primeFactor)
	})
}

// G - Divisors of Binomial Coefficient
// https://atcoder.jp/contests/abc227/tasks/abc227_g
// 二项式系数的约数个数模998244353.
// n<=1e12,k<=min(n,1e6).
// C(n,k)=n!/(k!(n-k)!).
// 即为 n*(n-1)*...*(n-k+1) / k!.
func abc227g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var n, k int
	fmt.Fscan(in, &n, &k)

	mp := make(map[int]int)
	table := NewPrimeTable(1e6 + 10)
	table.EnumerateRangePrimeFactors(n-k+1, n+1, func(num, primeFactor int) {
		mp[primeFactor]++
	})
	table.EnumerateRangePrimeFactors(1, k+1, func(num, primeFactor int) {
		mp[primeFactor]--
	})

	res := 1
	for _, v := range mp {
		res = res * (v + 1) % MOD
	}
	fmt.Fprintln(out, res)
}

const S int = 32768

type PrimeTable struct {
	done   int
	primes []int
	sieve  []bool
}

func NewPrimeTable(limit int) *PrimeTable {
	res := &PrimeTable{
		done:   2,
		primes: []int{2},
		sieve:  make([]bool, S+1),
	}
	res.expand(limit + 1)
	return res
}

// 返回小于等于limit的所有素数.
func (table *PrimeTable) GetPrimes(limit int) []int {
	limit++
	table.expand(limit)
	k := sort.Search(len(table.primes), func(i int) bool { return table.primes[i] >= limit })
	return table.primes[:k]
}

// 区间质因数分解(区间筛).
// 遍历区间[start, end)内`每个数的所有素因子`.
// f(n, factor)会被调用多次, 其中n是[start, end)内的数, factor是n的一个素因子.
// end<=1e14.
// O((end-start)*loglog(end)).
func (table *PrimeTable) EnumerateRangePrimeFactors(start, end int, f func(num, primeFactor int)) {
	n := end - start
	primes := table.GetPrimes(int(math.Sqrt(float64(end))))
	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = start + i
	}
	for _, p := range primes {
		pp := 1
		for {
			if pp > end/p {
				break
			}
			pp *= p
			s := ((start + pp - 1) / pp) * pp
			for s < end {
				f(s, p)
				res[s-start] /= p
				s += pp
			}
		}
	}
	for i, v := range res {
		if v > 1 {
			f(start+i, v)
		}
	}
}

func (table *PrimeTable) expand(limit int) {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
