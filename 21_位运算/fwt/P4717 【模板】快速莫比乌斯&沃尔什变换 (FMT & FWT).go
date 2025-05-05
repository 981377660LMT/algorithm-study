// deprecated.
//
// !FWT 参考 https://maspypy.github.io/library/setfunc/hadamard.hpp
// !https://github.dev/EndlessCheng/codeforces-go/blob/3dd70515200872705893d52dc5dad174f2c3b5f3/copypasta/math_fwt.go#L27
// !快速沃尔什变换 fast Walsh–Hadamard transform, FWT, FWHT
// 在算法竞赛中，FWT 是用于解决对下标进行【位运算卷积】问题的方法
// 一个常见的应用场景是对频率数组求 FWT
// 例如，求一个数组的三个元素的最大异或和，在值域不大的情况下，
// 可以先求出该数组的频率数组与频率数组的 FWT，即得到两个元素的所有异或和（及组成该异或和的元素对数），
// 然后枚举两元素异或和，在原数组的异或字典树上查询最大异或和
// 具体到名称，OR 上的 FWT 也叫 fast zeta transform，AND 上的 FWT 也叫 fast mobius transform
// fast Zeta transformation又名子集和dp(SOS DP, sum over subset)

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

	orRes := fwt(nums1, nums2, fwtOR, MOD-1)
	andRes := fwt(nums1, nums2, fwtAND, MOD-1)
	xorRes := fwt(nums1, nums2, fwtXOR, (MOD+1)/2) // /2
	printRes := func(res []int) {
		for _, num := range res {
			fmt.Fprint(out, num, " ")
		}
		fmt.Fprintln(out)
	}

	printRes(orRes)
	printRes(andRes)
	printRes(xorRes)
}

// 求 OR 和 AND 时 invOp = -1
// 求 XOR 时 invOp = inv(2)
func fwt(a, b []int, fwtFunc func([]int, int) []int, invOp int) []int {
	// 不修改原始数组
	a = fwtFunc(append([]int(nil), a...), 1)
	b = fwtFunc(append([]int(nil), b...), 1)
	for i, v := range b {
		a[i] = a[i] * v % MOD // !是否需要MOD
	}
	c := fwtFunc(a, invOp)
	return c
}

func fwtOR(a []int, op int) []int {
	n := len(a)
	for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
		for i := 0; i < n; i += l {
			for j := 0; j < k; j++ {
				a[i+j+k] = (a[i+j+k] + a[i+j]*op) % MOD
			}
		}
	}
	return a
}

func fwtAND(a []int, op int) []int {
	n := len(a)
	for l, k := 2, 1; l <= n; l, k = l<<1, k<<1 {
		for i := 0; i < n; i += l {
			for j := 0; j < k; j++ {
				a[i+j] = (a[i+j] + a[i+j+k]*op) % MOD
			}
		}
	}
	return a
}

func fwtXOR(a []int, op int) []int {
	n := len(a)
	for l_, k := 2, 1; l_ <= n; l_, k = l_<<1, k<<1 {
		for i := 0; i < n; i += l_ {
			for j := 0; j < k; j++ {
				a[i+j], a[i+j+k] = (a[i+j]+a[i+j+k])*op%MOD, (a[i+j]+MOD-a[i+j+k])*op%MOD
			}
		}
	}
	return a
}
