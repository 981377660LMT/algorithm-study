package main

import (
	"fmt"
	"math"
	"math/big"
)

func main() {
	ps := NewPrimeSum(2e8)
	ps.CalSum()
	fmt.Println(ps.Get(2e8))
}

type PrimeSum struct {
	n, sqN       int
	sumLo, sumHi []int
}

// 给定n和函数f, 计算前缀和 sum_{p <= x} f(p),
// 其中x必须要形如floor(n/i)的形式。
//
//	O(n^(3/4)/logn) time, O(n^(1/2)) space.
func NewPrimeSum(n int) *PrimeSum {
	return &PrimeSum{n: n, sqN: int(math.Sqrt(float64(n)))}
}

func (ps *PrimeSum) Cal(f func(int) int) {
	primes := newEratosthenesSieve(ps.sqN).GetPrimes()
	ps.sumLo = make([]int, ps.sqN+1)
	ps.sumHi = make([]int, ps.sqN+1)
	for i := 1; i <= ps.sqN; i++ {
		ps.sumLo[i] = f(i) - 1
		ps.sumHi[i] = f(ps.n/i) - 1
	}
	for _, p := range primes {
		pp := p * p
		if pp > ps.n {
			break
		}
		R := min(ps.sqN, ps.n/pp)
		M := ps.sqN / p
		x := ps.sumLo[p-1]
		fp := ps.sumLo[p] - ps.sumLo[p-1]
		for i := 1; i <= M; i++ {
			ps.sumHi[i] -= fp * (ps.sumHi[i*p] - x)
		}
		for i := M + 1; i <= R; i++ {
			ps.sumHi[i] -= fp * (ps.sumLo[ps.n/(i*p)] - x)
		}
		for i := ps.sqN; i >= pp; i-- {
			ps.sumLo[i] -= fp * (ps.sumLo[i/p] - x)
		}
	}
}

func (ps *PrimeSum) CalCount() {
	ps.Cal(func(x int) int { return x })
}

func (ps *PrimeSum) CalSum() {
	ps.Cal(func(x int) int { return (x * (x + 1)) >> 1 })
}

// x 必须要是 n/i 的形式
func (ps *PrimeSum) Get(x int) int {
	if x <= ps.sqN {
		return ps.sumLo[x]
	}
	return ps.sumHi[ps.n/x]
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 埃氏筛
type eratosthenesSieve struct {
	minPrime []int
}

func newEratosthenesSieve(maxN int) *eratosthenesSieve {
	minPrime := make([]int, maxN+1)
	for i := range minPrime {
		minPrime[i] = i
	}
	upper := int(math.Sqrt(float64(maxN))) + 1
	for i := 2; i < upper; i++ {
		if minPrime[i] < i {
			continue
		}
		for j := i * i; j <= maxN; j += i {
			if minPrime[j] == j {
				minPrime[j] = i
			}
		}
	}
	return &eratosthenesSieve{minPrime}
}

func (es *eratosthenesSieve) IsPrime(n int) bool {
	if n < 2 {
		return false
	}
	return es.minPrime[n] == n
}

func (es *eratosthenesSieve) GetPrimeFactors(n int) map[int]int {
	res := make(map[int]int)
	for n > 1 {
		m := es.minPrime[n]
		res[m]++
		n /= m
	}
	return res
}

func (es *eratosthenesSieve) GetPrimes() []int {
	res := []int{}
	for i, x := range es.minPrime {
		if i >= 2 && i == x {
			res = append(res, x)
		}
	}
	return res
}
