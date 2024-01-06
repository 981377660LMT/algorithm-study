package main

import (
	"math"
	"sort"
)

// template <typename T>
// void divisor_zeta(vc<T>& A) {
//   assert(A[0] == 0);
//   int N = len(A) - 1;
//   auto P = primetable(N);
//   for (auto&& p: P) { FOR3(x, 1, N / p + 1) A[p * x] += A[x]; }
// }

// template <typename T>
// void divisor_mobius(vc<T>& A) {
//   assert(A[0] == 0);
//   int N = len(A) - 1;
//   auto P = primetable(N);
//   for (auto&& p: P) { FOR3_R(x, 1, N / p + 1) A[p * x] -= A[x]; }
// }

// template <typename T>
// void multiplier_zeta(vc<T>& A) {
//   assert(A[0] == 0);
//   int N = len(A) - 1;
//   auto P = primetable(N);
//   for (auto&& p: P) { FOR3_R(x, 1, N / p + 1) A[x] += A[p * x]; }
// }

// template <typename T>
// void multiplier_mobius(vc<T>& A) {
//   assert(A[0] == 0);
//   int N = len(A) - 1;
//   auto P = primetable(N);
//   for (auto&& p: P) { FOR3(x, 1, N / p + 1) A[x] -= A[p * x]; }
// }

const S int = 32768

var T = NewPT()

type PT struct {
	done   int
	primes []int
	sieve  []int
}

func NewPT() *PT {
	return &PT{
		done:   2,
		primes: []int{2},
		sieve:  make([]int, S+1),
	}
}

// 返回小于等于limit的所有素数.
func (table *PT) GetPrimes(limit int) []int {
	limit++
	if table.done < limit {
		table.done = limit
		R := limit / 2
		for i := range table.sieve {
			table.sieve[i] = 0
		}
		table.primes = make([]int, 0, int((float64(limit)/math.Log(float64(limit)))*1.1))
		table.primes = append(table.primes, 2)
		cp := [][2]int{}
		for i := int(3); i <= S; i += 2 {
			if table.sieve[i] == 0 {
				cp = append(cp, [2]int{i, i * i / 2})
				for j := i * i; j <= S; j += 2 * i {
					table.sieve[j] = 1
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

	k := sort.Search(len(table.primes), func(i int) bool { return table.primes[i] >= limit })
	return table.primes[:k]
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
