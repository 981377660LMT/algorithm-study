// https://www.luogu.com.cn/problem/CF797E

package main

import (
	"bufio"
	"fmt"
	"os"
)

// 给定一个长为n的序列，q次询问.
// 每次询问给出 p,step,不断执行操作 p = p + a[p] + step, 问至少需要多少次才能使 p 大于 n。
//
// 类似弹飞绵羊.
// !step < sqrt => 预处理, step >= sqrt => 暴力.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	sqrt := 100               // int(math.Sqrt(float64(n)))
	dp := make([][]int, sqrt) // dp[step][pos]
	for i := range dp {
		dp[i] = make([]int, n)
	}
	for step := 0; step < sqrt; step++ {
		for pos := n - 1; pos >= 0; pos-- {
			if pos+nums[pos]+step >= n {
				dp[step][pos] = 1
			} else {
				dp[step][pos] = dp[step][pos+nums[pos]+step] + 1
			}
		}
	}

	query := func(pos, step int) int {
		if step < sqrt {
			return dp[step][pos]
		} else {
			res := 0
			for pos < n {
				res++
				pos += nums[pos] + step
			}
			return res
		}
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var pos, step int
		fmt.Fscan(in, &pos, &step)
		pos--
		fmt.Fprintln(out, query(pos, step))
	}
}
