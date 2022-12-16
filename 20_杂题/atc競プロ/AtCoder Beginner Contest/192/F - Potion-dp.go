//  F - Potion
//  魔法药水(背包dp)
//  有n个数，你可以在第0个时刻条任意多个数出来，并让计数器加上他们的和，
//  设你挑了k个数，那么之后每个时刻计数器都会加k，你要最小化计数器恰好到x的时间
//  !n<=100 Ai<=1e7 1e9<=x<=1e18
//  !枚举选的个数(固定k会比较方便) dp[index][remain][mod] O(n^4) =>
//  !选k个数,模k等于target,求选出的数之和的最大值

package main

import (
	"bufio"
	"fmt"
	"os"
)

func potion(magic []int, target int) int {
	const INF int = int(1e18)
	n := len(magic)

	cal := func(k int) int {
		memo := [110][110][110]int{}
		for i := 0; i < 110; i++ {
			for j := 0; j < 110; j++ {
				for k := 0; k < 110; k++ {
					memo[i][j][k] = INF
				}
			}
		}

		var dfs func(index, remain, mod int) int
		dfs = func(index, remain, mod int) int {
			if remain < 0 {
				return -INF
			}
			if index == n {
				if (remain == 0) && (mod == (target % k)) {
					return 0
				}
				return -INF
			}

			if cand := memo[index][remain][mod]; cand != INF {
				return cand
			}

			res := dfs(index+1, remain, mod) // jump
			if remain > 0 {                  // select
				res = max(res, magic[index]+dfs(index+1, remain-1, (mod+magic[index])%k))
			}
			memo[index][remain][mod] = res
			return res
		}

		return dfs(0, k, 0)
	}

	res := INF
	for k := 1; k <= n; k++ {
		maxSum := cal(k)
		if maxSum > 0 {
			res = min(res, (target-maxSum)/k)
		}
	}
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, target int
	fmt.Fscan(in, &n, &target)
	magic := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &magic[i])
	}

	fmt.Fprintln(out, potion(magic, target))
}
