// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/math_fwt.go

package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

// https://www.luogu.com.cn/problem/P4717
// 给定长度为2^n的两个序列(n<=17)
// !记 Ck = ∑Ai*Bj (其中 op(i,j)=k,op可以是or,and,xor)
// 求出C0,C1,...,C2^n-1 模998244353的值
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n uint8
	fmt.Fscan(in, &n)
	nums1, nums2 := make([]int, 1<<n), make([]int, 1<<n)
	for i := range nums1 {
		fmt.Fscan(in, &nums1[i])
	}
	for i := range nums2 {
		fmt.Fscan(in, &nums2[i])
	}

	printRes := func(res []int) {
		for _, num := range res {
			fmt.Fprint(out, num, " ")
		}
		fmt.Fprintln(out)
	}

	solve := func(fwtFunc func([]int, bool)) {
		a, b := append([]int(nil), nums1...), append([]int(nil), nums2...)
		fwtFunc(a, false)
		fwtFunc(b, false)
		for i, v := range b {
			a[i] = a[i] * v % MOD // !是否需要MOD
		}
		fwtFunc(a, true)
		printRes(a)
	}

	solve(FwtOr)
	solve(FwtAnd)
	solve(FwtXor)
}

func FwtXor(a []int, inv bool) {
	if !inv {
		n := len(a)
		if n&(n-1) != 0 {
			panic("n is not a power of 2")
		}
		for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
			for i := 0; i < n; i += l {
				for j := 0; j < k; j++ {
					a[i+j], a[i+j+k] = add(a[i+j], a[i+j+k]), sub(a[i+j], a[i+j+k])
				}
			}
		}
	} else {
		n := len(a)
		if n&(n-1) != 0 {
			panic("n is not a power of 2")
		}
		for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
			for i := 0; i < n; i += l {
				for j := 0; j < k; j++ {
					a[i+j], a[i+j+k] = div2(add(a[i+j], a[i+j+k])), div2(sub(a[i+j], a[i+j+k]))
				}
			}
		}
	}
}

func FwtOr(a []int, inv bool) {
	if !inv {
		subsetZetaTransform(a)
	} else {
		subsetMoebiusTransform(a)
	}
}

func FwtAnd(a []int, inv bool) {
	if !inv {
		supersetZetaTransform(a)
	} else {
		supersetMoebiusTransform(a)
	}
}

func supersetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j] = add(f[j], f[j|i])
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
				f[j] = sub(f[j], f[j|i])
			}
		}
	}
}

func subsetZetaTransform(f []int) {
	n := len(f)
	if n&(n-1) != 0 {
		panic("n is not a power of 2")
	}
	for i := 1; i < n; i <<= 1 {
		for j := 0; j < n; j++ {
			if (j & i) == 0 {
				f[j|i] = add(f[j|i], f[j])
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
				f[j|i] = sub(f[j|i], f[j])
			}
		}
	}
}

func add(a, b int) int {
	a += b
	if a >= MOD {
		a -= MOD
	}
	return a
}

func sub(a, b int) int {
	a -= b
	if a < 0 {
		a += MOD
	}
	return a
}

func div2(a int) int {
	if a&1 > 0 {
		a += MOD
	}
	return a >> 1
}

// 将 (-mod,mod) 范围内的 a 变成 [0,mod) 范围内
// 原理是负数右移会不断补 1，所以最后二进制都是 1，因此返回值等价于 a+_mod
// 而对于非负数，右移后二进制全为 0，所以返回结果仍然是 a
func norm32(a int32, mod int32) int32 {
	return a + a>>31&mod
}

func norm64(a int64, mod int64) int64 {
	return a + a>>63&mod
}
