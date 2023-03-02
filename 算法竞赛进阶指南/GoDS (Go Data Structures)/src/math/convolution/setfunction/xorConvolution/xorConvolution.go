package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {
	// https://judge.yosupo.jp/problem/bitwise_xor_convolution
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
	conv := XorConvolution(nums1, nums2)
	for i := 0; i < 1<<n; i++ {
		fmt.Fprint(out, conv[i], " ")
	}
}

func XorConvolution(a, b []int) []int {
	a = append(a[:0:0], a...)
	b = append(b[:0:0], b...)
	a = walshHadamardTransform(a, 1)
	b = walshHadamardTransform(b, 1)
	for i := range a {
		a[i] = a[i] * b[i] % MOD
	}
	res := walshHadamardTransform(a, (MOD+1)/2)
	return res
}

func walshHadamardTransform(f []int, op int) []int {
	n := len(f)
	for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
		for i := 0; i < n; i += l {
			for j := 0; j < k; j++ {
				f[i+j], f[i+j+k] = (f[i+j]+f[i+j+k])*op%MOD, (f[i+j]+MOD-f[i+j+k])*op%MOD
			}
		}
	}
	return f
}
