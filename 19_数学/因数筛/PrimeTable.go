// enumeratePrimeFactors

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// demo()
	abc227g()
}

func demo() {
	// lpfTable := Table.GetLpfTable(2e6)
	fmt.Println(Table.GetLpfTable(100))
}

// https://atcoder.jp/contests/abc227/tasks/abc227_g
// 求C(n,k)的正约数个数, 模998244353.
// n<=1e12,k<=1e6.
func abc227g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD = 998244353

	var n, k int
	fmt.Fscan(in, &n, &k)

	counter := make(map[int]int)
	Table.EnumerateRangePrimeFactors(1, k+1, func(n, factor int) {
		counter[factor]--
	})
	Table.EnumerateRangePrimeFactors(n-k+1, n+1, func(n, factor int) {
		counter[factor]++
	})
	res := 1
	for _, v := range counter {
		res = (res * (v + 1)) % MOD
	}
	fmt.Fprintln(out, res)
}

const S int = 32768

var Table = NewPrimeTable(1e6 + 10)

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

// 遍历区间[start, end)内所有数的所有素因子.
// f(n, factor)会被调用多次, 其中n是[start, end)内的数, factor是n的一个素因子.
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

// 遍历[2, n]内所有数的所有素因子.
func (table *PrimeTable) EnumeratePrefixPrimeFactors(n int, f func(num, primeFactor int)) {
	primes := table.GetPrimes(n)
	for _, p := range primes {
		for x := n / p; x >= 1; x-- {
			f(x*p, p)
		}
	}
}

// 返回[0, limit]内所有数的最小素因子.0,1不是素数,返回-1.
// lpf: lowest prime factor(最小素因子)
func (table *PrimeTable) GetLpfTable(limit int) []int {
	primes := table.GetPrimes(limit)
	res := make([]int, limit+1)
	for i := len(primes) - 1; i >= 0; i-- {
		p := primes[i]
		upper := limit/p + 1
		for j := 1; j < upper; j++ {
			res[p*j] = p
		}
	}
	res[0] = -1
	if limit >= 1 {
		res[1] = -1
	}
	return res
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
