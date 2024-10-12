// restoreSa/restoreSuffixArray.go

package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://atcoder.jp/contests/arc044/tasks/arc044_d
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	sa := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &sa[i])
		sa[i]--
	}

	res := RestoreSa(n, func(i int32) int32 { return sa[i] })
	if maxs32(res...) >= 26 {
		fmt.Fprintln(out, -1)
		return
	}

	sb := make([]byte, n)
	for i := int32(0); i < n; i++ {
		sb[i] = byte('A' + res[i])
	}
	fmt.Fprintln(out, string(sb))
}

// 给定一个表示后缀排名的[0, n)的排列，返回字典序最小的原数组.
func RestoreSa(n int32, f func(i int32) int32) []int32 {
	rank := make([]int32, n)
	for i := int32(0); i < n; i++ {
		rank[f(i)] = i
	}
	for i := int32(0); i < n; i++ {
		if rank[f(i)] != i || f(rank[i]) != i {
			panic("Invalid rank")
		}
	}
	res := make([]int32, n)
	for k := int32(1); k < n; k++ {
		i, j := f(k-1), f(k)
		res[j] = res[i]
		if i < n-1 && (j == n-1 || rank[i+1] > rank[j+1]) {
			res[j]++
		}
	}
	return res
}

func maxs32(nums ...int32) int32 {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
