package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	fmt.Println(GeometricSequenceSumK(2, 4, 3, 1e9+7))
}

// https://atcoder.jp/contests/abc293/editorial/5955
func abc293e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var x, n, mod int
	fmt.Fscan(in, &x, &n, &mod)
	_, sum := GeometicSequenceSum(x, n, mod)
	fmt.Fprintln(out, sum)
}

// GeometicSequenceSum(等比数列求和)
// O(logn) 求 (x^n, x^0 + ... + x^(n - 1)) 模 mod
// x,n,mod >=1
// eg:
// x = 2, n = 4, mod = 1000000007
// res = (16, 15)
// 其中
// res[0] = 2^4 = 16
// res[1] = 2^0 + 2^1 + 2^2 + 2^3 = 15
func GeometicSequenceSum(x, n int, mod int) (int, int) {
	if mod == 1 {
		return 0, 0
	}
	sum, p := 1, x // res = x^0 + ... + x^(len - 1), p = x^len
	start := bits.Len(uint(n)) - 2
	for d := start; d >= 0; d-- {
		sum = (sum * (p + 1)) % mod
		p = p * p % mod
		if ((n >> d) & 1) == 1 {
			sum = (sum + p) % mod
			p = p * x % mod
		}
	}
	return p % mod, sum % mod
}

// res[j] = sum(i^k * x^i, i = 0..n-1), j = 0..k
// eg:
// x = 2, n = 4, k = 3
// res = [15 34 90 250]
// 其中
// res[0] = 1 * 2^0 + 1 * 2^1 + 1 * 2^2 + 1 * 2^3 = 15
// res[1] = 0 * 2^0 + 1 * 2^1 + 2 * 2^2 + 3 * 2^3 = 34
// res[2] = 0 * 2^0 + 1 * 2^1 + 4 * 2^2 + 9 * 2^3 = 90
// res[3] = 0 * 2^0 + 1 * 2^1 + 8 * 2^2 + 27 * 2^3 = 250
func GeometricSequenceSumK(x, n, k int, mod int) []int {
	comb := make([][]int, k+1)
	for i := 0; i <= k; i++ {
		comb[i] = make([]int, k+1)
	}
	comb[0][0] = 1
	for i := 0; i < k; i++ {
		for j := 0; j <= i; j++ {
			comb[i+1][j] = (comb[i+1][j] + comb[i][j]) % mod
			comb[i+1][j+1] = (comb[i+1][j+1] + comb[i][j]) % mod
		}
	}

	// (n ,x^n, sum(i^k * x^i, i = 0..n-1))
	type triple struct {
		v1, v2 int
		v3     []int // 长 k+1
	}
	mul := func(left, right *triple) *triple {
		n1, r1, a1 := left.v1, left.v2, left.v3
		n2, r2, a2 := right.v1, right.v2, right.v3
		a3 := make([]int, k+1)
		for i := 0; i <= k; i++ {
			a3[i] = a1[i]
			tmp := r1
			for j := 0; j < i+1; j++ {
				a3[i] = (a3[i] + comb[i][j]*tmp*a2[i-j]) % mod
				tmp = tmp * n1 % mod
			}
		}
		return &triple{(n1 + n2) % mod, (r1 * r2) % mod, a3}
	}

	res := &triple{0, 1, make([]int, k+1)}
	base := &triple{1, x, make([]int, k+1)}
	base.v3[0] = 1
	for n > 0 {
		if n&1 == 1 {
			res = mul(res, base)
		}
		base = mul(base, base)
		n >>= 1
	}
	return res.v3
}
