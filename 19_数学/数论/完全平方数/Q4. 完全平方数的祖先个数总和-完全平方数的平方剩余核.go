// Q4. 完全平方数的祖先个数总和-完全平方数的平方剩余核

package main

import (
	"math"
)

// “平方剩余核”（square-free kernel）：把数中偶数次出现的质因子全部抵消后剩下的乘积.
// 若两个数 a、b 的平方剩余核相同 (sf[a]==sf[b])，则 a·b 为完全平方数.
func getSquareFreeKernel(e *eratosthenesSieve, n int) int {
	res := 1
	for p, c := range e.GetPrimeFactors(n) {
		if c&1 == 1 {
			res *= p
		}
	}
	return res
}

var e *eratosthenesSieve

func init() {
	e = newEratosthenesSieve(1e5 + 10)
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

func sumOfAncestors(n int, edges [][]int, nums []int) int64 {
	adjList := make([][]int, n)
	for _, e := range edges {
		u, v := e[0], e[1]
		adjList[u] = append(adjList[u], v)
		adjList[v] = append(adjList[v], u)
	}

	kernal := make([]int, n)
	for i, v := range nums {
		kernal[i] = getSquareFreeKernel(e, v)
	}

	res := 0
	counter := make(map[int]int)

	var dfs func(int, int)
	dfs = func(cur, parent int) {
		if parent != -1 {
			res += counter[kernal[cur]]
		}
		counter[kernal[cur]]++
		for _, v := range adjList[cur] {
			if v != parent {
				dfs(v, cur)
			}
		}
		counter[kernal[cur]]--
	}

	dfs(0, -1)
	return int64(res)
}
