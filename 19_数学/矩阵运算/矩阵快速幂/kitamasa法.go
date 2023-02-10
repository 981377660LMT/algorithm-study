// https://tjkendev.github.io/procon-library/python/series/kitamasa.html
// O(k^2logn) 求线性递推式的第n项 (比矩阵快速幂快一个k)
// 線形漸化式 dp[i+k] = c0*dp[i] + c1*dp[i+1] + ... + ci+k-1*dp[i+k-1] (i>=0) の第n項を求める
// C: 系数 c0,c1,...,ci+k-1
// A: dp[0]-dp[k-1] 初始值
// n: 第n项

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int = 1e9 + 7

func kitamasa(C, A []int, n int) int {
	if n == 0 {
		return A[0]
	}

	k := len(C)
	C0 := make([]int, k)
	C1 := make([]int, k)
	C0[1] = 1

	inc := func(k int, C0, C1 []int) {
		C1[0] = C0[k-1] * C[0] % MOD
		for i := 0; i < k-1; i++ {
			C1[i+1] = (C0[i] + C0[k-1]*C[i+1]) % MOD
		}
	}

	dbl := func(k int, C0, C1 []int) {
		D0 := make([]int, k)
		D1 := make([]int, k)
		copy(D0, C0)
		for j := 0; j < k; j++ {
			C1[j] = C0[0] * C0[j] % MOD
		}
		for i := 1; i < k; i++ {
			inc(k, D0, D1)
			for j := 0; j < k; j++ {
				C1[j] += C0[i] * D1[j] % MOD
			}
			D0, D1 = D1, D0
		}
		for i := 0; i < k; i++ {
			C1[i] %= MOD
		}
	}

	p := bits.Len(uint(n)) - 1

	for p > 0 {
		p--
		dbl(k, C0, C1)
		C0, C1 = C1, C0
		if n>>p&1 == 1 {
			inc(k, C0, C1)
			C0, C1 = C1, C0
		}
	}

	res := 0
	for i := 0; i < k; i++ {
		res = (res + C0[i]*A[i]) % MOD
	}
	return res
}

// 斐波那契数列第i项(i>=0)模1e9+7的值
// 0，1，1，2，3，5，8
func fib(n int) int {
	if n == 0 {
		return 0
	}
	C := []int{1, 1}
	A := []int{0, 1}
	return kitamasa(C, A, n-1)
}

func main() {
	// dp[0]=dp[1]=...=dp[k-1]=1
	// dp[i]=dp[i-1]+dp[i-2]+...+dp[i-k]
	// 求dp[n-1] 模 1e9+7 的值
	// k<=1000 n<=1e18

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var k, n int
	fmt.Fscan(in, &k, &n)
	C := make([]int, k)
	A := make([]int, k)
	for i := 0; i < k; i++ {
		C[i] = 1
		A[i] = 1
	}

	fmt.Fprintln(out, kitamasa(C, A, n-1))
}
