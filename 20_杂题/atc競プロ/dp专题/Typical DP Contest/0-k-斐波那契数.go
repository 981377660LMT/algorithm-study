// https://atcoder.jp/contests/abc401/editorial/12689
// 给定正整数 **N** 和 **K**，定义一个长度为 **N+1** 的数列 **A = (A₀, A₁, ..., Aₙ)**，其各元素的值按照以下方式定义：
//
// - 当 **0 ≤ i < K** 时，**Aᵢ = 1**
// - 当 **K ≤ i** 时，**Aᵢ = Aᵢ₋ₖ + Aᵢ₋ₖ₊₁ + ... + Aᵢ₋₁**
//
// 求 **Aₙ** 除以 **10⁹** 的余数。
//
// ### 约束条件
// - **1 ≤ N, K ≤ 10⁶**
// - 输入的数值均为整数

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9

	var n, k int
	fmt.Fscan(in, &n, &k)

	fmt.Fprintln(out, KFibonacci(n, k, MOD))
}

func KFibonacci(n, k, mod int) int {
	if n < k {
		return 1 % mod
	}

	pre := make([]int, k)
	for i := range pre {
		pre[i] = 1
	}
	cur := k % mod
	for index := 0; index < n-k+1; index++ {
		pos := index % k
		tmp := pre[pos]
		pre[pos] = cur
		cur = (2*cur - tmp) % mod
	}
	res := pre[n%k]
	if res < 0 {
		res += mod
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
