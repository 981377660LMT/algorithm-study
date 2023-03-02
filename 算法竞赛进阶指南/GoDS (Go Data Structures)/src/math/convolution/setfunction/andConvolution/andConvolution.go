package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	// https://judge.yosupo.jp/problem/bitwise_and_convolution
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums1 := make([]int, 1<<n)
	nums2 := make([]int, 1<<n)
	for i := 0; i < 1<<n; i++ {
		fmt.Fscan(in, &nums1[i])
	}
	for i := 0; i < 1<<n; i++ {
		fmt.Fscan(in, &nums2[i])
	}
	conv := AndConvolution(nums1, nums2)
	for i := 0; i < 1<<n; i++ {
		fmt.Fprint(out, conv[i], " ")
	}
}

func AndConvolution(a, b []int) []int {
	a = append(a[:0:0], a...)
	b = append(b[:0:0], b...)
	supersetZetaTransform(a)
	supersetZetaTransform(b)
	for i := range a {
		a[i] = (a[i]*b[i]%MOD + MOD) % MOD
	}
	supersetMoebiusTransform(a)
	return a
}

func supersetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = (f[j] + f[j|i]) % MOD
			}
		}
	}
}

func supersetMoebiusTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = ((f[j]-f[j|i])%MOD + MOD) % MOD
			}
		}
	}
}
