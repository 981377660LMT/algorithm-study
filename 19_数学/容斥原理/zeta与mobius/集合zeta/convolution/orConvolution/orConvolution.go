// template <typename T>
// vector<T> and_convolution(vector<T> a, vector<T> b) {
//   superset_zeta_transform(a);
//   superset_zeta_transform(b);
//   for (int i = 0; i < (int)a.size(); i++) a[i] *= b[i];
//   superset_mobius_transform(a);
//   return a;
// }

package main

import (
	"fmt"
	"math/rand"
)

func main() {
	for i := 0; i < 100; i++ {
		// test and convolution with brute force
		n := rand.Intn(12) + 1
		A := make([]int, 1<<n)
		B := make([]int, 1<<n)
		for i := range A {
			A[i] = rand.Intn(10)
			B[i] = rand.Intn(10)
		}

		bruteForceOrConvolution := func(A, B []int) []int {
			n := len(A)
			C := make([]int, n)
			for i := 0; i < n; i++ {
				for j := 0; j < n; j++ {
					C[i|j] += A[i] * B[j]
				}
			}
			return C
		}

		C1 := OrConvolution(A, B)
		C2 := bruteForceOrConvolution(A, B)
		for i := range C1 {
			if C1[i] != C2[i] {
				panic("wrong answer")
			}
		}
	}

	fmt.Println("Correct")
}

const MOD int = 998244353

func OrConvolution(a, b []int) []int {
	a = append(a[:0:0], a...)
	b = append(b[:0:0], b...)
	subsetZetaTransform(a)
	subsetZetaTransform(b)
	for i := range a {
		a[i] = a[i] * b[i] % MOD
	}
	subsetMoebiusTransform(a)
	return a
}

func subsetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] = (f[j|i] + f[j]) % MOD
			}
		}
	}
}

func subsetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] = (f[j|i] - f[j] + MOD) % MOD
			}
		}
	}
}
