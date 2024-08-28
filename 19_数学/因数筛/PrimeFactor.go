// 获取数的质因数分解
// GetPrimeFactorsByLpf: O(logn), n<=10^7，适用于小数质因数分解，快
// GetPrimeFactorsBig: O(n^1/4), n<=10^18，适用于大数质因数分解，慢
// AllLcm

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"os"
	"sort"
	"time"
)

func main() {
	abc152e()
	// demo()
}

func demo() {
	time1 := time.Now()
	for i := int(1e17); i <= int(1e17+1000); i++ {
		GetPrimeFactorsBig(i)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))

	time1 = time.Now()
	lpf := Table.GetLpfTable(int(1e7))
	for i := 0; i < int(1e7); i++ {
		GetPrimeFactorsByLpf(i, lpf)
	}
	time2 = time.Now()
	fmt.Println(time2.Sub(time1))
}

// https://atcoder.jp/contests/abc152/tasks/abc152_e
// 给定数组A,要求构造一个数组B满足A[i]*B[i]=A[j]*B[j] (i!=j).
// 求B元素之和的最小值模1e9+7.
//
// 即 lcm * sum(1/A[i]) 的和.
func abc152e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}

	invSum := 0
	for i := 0; i < n; i++ {
		invSum += pow(A[i], MOD-2, MOD)
		invSum %= MOD
	}
	allLcm := AllLcm(A, MOD, true)
	fmt.Println(allLcm * invSum % MOD)
}

func GetPrimeFactorsByLpf(n int, lpf []int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		p := lpf[n]
		e := 0
		for n%p == 0 {
			n /= p
			e++
		}
		res[p] = e
	}
	return res
}

func GetPrimeFactorsBig(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		p := PollardRhoPrimeFactor(n)
		for n%p == 0 {
			n /= p
			res[p]++
		}
	}
	return res
}

// 所有数的最小公倍数模mod.
func AllLcm(nums []int, mod int, useLpf bool) int {
	if len(nums) == 0 {
		return 1
	}
	mp := map[int]int{}
	mx := 0
	for _, n := range nums {
		if n > mx {
			mx = n
		}
	}
	var lpf []int
	if useLpf {
		lpf = Table.GetLpfTable(mx)
	}
	if useLpf {
		for _, n := range nums {
			pfs := GetPrimeFactorsByLpf(n, lpf)
			for p, e := range pfs {
				if e > mp[p] {
					mp[p] = e
				}
			}
		}
	} else {
		for _, n := range nums {
			pfs := GetPrimeFactorsBig(n)
			for p, e := range pfs {
				if e > mp[p] {
					mp[p] = e
				}
			}
		}
	}

	pow := func(x, n int) int {
		res := 1
		for n > 0 {
			if n&1 == 1 {
				res = res * x % mod
			}
			x = x * x % mod
			n >>= 1
		}
		return res
	}

	x := 1
	for p, e := range mp {
		x = x * pow(p, e) % mod
	}
	return x
}

// Pollard-Rho 算法求出一个因子 O(n^1/4)
func PollardRhoFactor(n int) int {
	if n == 4 {
		return 2
	}
	if IsPrimeMillerRabin(n) {
		return n
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

// 判断质数+求最大质因子
// 先用 Pollard-Rho 算法求出一个因子，然后递归求最大质因子
// https://zhuanlan.zhihu.com/p/267884783
// https://www.luogu.com.cn/problem/P4718
func PollardRhoPrimeFactor(n int) int {
	if n == 4 {
		return 2
	}
	if IsPrimeMillerRabin(n) {
		return n
	}

	cache := map[int]int{}
	var getPrimeFactor func(int) int
	getPrimeFactor = func(x int) (res int) {
		if cache[x] > 0 {
			return cache[x]
		}
		p := PollardRhoFactor(x)
		if p == x {
			cache[x] = x
			return p
		}
		res = max(getPrimeFactor(p), getPrimeFactor(x/p))
		cache[x] = res
		return
	}

	return getPrimeFactor(n)
}

// IsPrime
func IsPrimeMillerRabin(n int) bool {
	return big.NewInt(int64(n)).ProbablyPrime(0)
}

const S int = 32768

var Table = NewPrimeTable(1e5 + 10)

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
		for i := int(3); i <= S; i += 2 {
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

func gcd(a, b int) int {
	for a != 0 {
		a, b = b%a, a
	}
	return b
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func pow(x, n, mod int) int {
	x %= mod
	res := 1
	for n > 0 {
		if n&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		n >>= 1
	}
	return res
}
