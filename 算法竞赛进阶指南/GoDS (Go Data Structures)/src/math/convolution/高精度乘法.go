package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// FFT求大数乘法(高精度乘法)
// github.com/EndlessCheng/codeforces-go
func 高精度乘法() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s1, s2 []byte
	fmt.Fscan(in, &s1, &s2)
	fmt.Println(string(BigMul(s1, s2)))
}

// 大数乘法
func BigMul(s1, s2 []byte) []rune {
	n, m := len(s1), len(s2)
	limit := 1 << uint(bits.Len(uint(n+m-1)))
	f := newNTT(limit)
	A := make([]int64, limit)
	for i, v := range s1 {
		A[n-1-i] = int64(v & 15)
	}
	B := make([]int64, limit)
	for i, v := range s2 {
		B[m-1-i] = int64(v & 15)
	}
	f.dft(A)
	f.dft(B)
	for i, v := range A {
		A[i] = v * B[i] % MOD
	}
	f.idft(A)

	res := make([]rune, n+m)
	for i := 0; i < n+m-1; i++ {
		res[i] += rune(A[i])
		res[i+1] = res[i] / 10
		res[i] %= 10
	}
	r := n + m
	for res[r-1] == 0 {
		r--
	}
	res = res[:r]
	for i := 0; i < r/2; i++ {
		res[i], res[r-1-i] = res[r-1-i]+'0', res[i]+'0'
	}
	if r&1 > 0 {
		res[r/2] += '0'
	}
	return res
}
