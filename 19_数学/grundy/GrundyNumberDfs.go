// https://nyaannyaan.github.io/library/math/grundy-number.hpp

package main

import (
	"fmt"
	"math"
)

var E *eratosthenesSieve
var F []map[int]int

func init() {
	E = newEratosthenesSieve(1e5 + 10)
	F = make([]map[int]int, 1e5+10)
	for i := 0; i <= 1e5; i++ {
		F[i] = E.GetPrimeFactors(i)
	}
}

func main() {
	memo := make([]int, 1e5+10)
	for i := range memo {
		memo[i] = -1
	}
	var grundy func(state int) int
	grundy = func(state int) int {
		if memo[state] != -1 {
			return memo[state]
		}

		nextStates := make(map[int]struct{})
		for p := range F[state] {
			if state-p >= 2 { // 如果是拆分成两个子状态,则当前grundy数为子状态的异或和
				nextStates[grundy(state-p)] = struct{}{}
			}
		}

		mex := 0
		for {
			if _, ok := nextStates[mex]; !ok {
				break
			}
			mex++
		}
		memo[state] = mex
		return mex
	}
	n := 12
	fmt.Println(grundy(n))
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
