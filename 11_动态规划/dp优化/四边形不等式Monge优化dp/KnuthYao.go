// https://blog.csdn.net/weixin_43914593/article/details/105150937
// 有一些常见的DP问题，通常是区间DP问题
// !dp[i][j]=min{dp[i][k]+dp[k+1][j]+w[i][j]} (i<=k<j)
// !dp[i][j]表示从i状态到j状态的最小花费,k在i和j之间滑动,k有一个最优值使得dp[i][j]最小
// !w[i][j]的性质非常重要,如果它满足四边形不等式和单调性，那么用DP计算dp的时候，就能进行四边形不等式优化
// !w[i1][i4]+w[i2][i3]>=w[i1][i3]+w[i2][i4] (i1<=i2,i3<=i4)
// 拿到题目后，先判断w是否单调、是否满足四边形不等式，再使用四边形不等式优化DP。

// https://beet-aizu.github.io/library/algorithm/knuthyao.cpp
// 四边形不等式优化 将区间dp从 O(n^3) 优化到 O(n^2)
// cost(i, k, j) 表示区间[i, j]的最优划分点为k时的代价
// eg:
// !1. 石子合并
// cost(i, k, j) = func(i, k, j int) int {
// 	return preSum[j+1] - preSum[i]
// }
// !2. Tree Construction 代价与分割点有关 (https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2488)
// cost(i, k, j) = func(i, k, j int) int {
//  return (xs[k+1]-xs[i])+(ys[k]-ys[j])
// }
// !注意有的时候不是dp[i][k]+dp[k+1][j],而是dp[i][k-1]+dp[k+1][j] ,需要稍微修改

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

// !dp[i][j]=min{dp[i][k]+dp[k+1][j]+cost(i,j,k)} (0<=i<=k<j<n)
func KnuthYao1(n int, cost func(i, k, j int) int) int {
	dp := make([][]int, n)
	ar := make([][]int, n)
	for i := 0; i < n; i++ {
		dp[i] = make([]int, n)
		ar[i] = make([]int, n)
		dp[i][i] = 0
		ar[i][i] = i
	}

	for w := 1; w < n; w++ {
		for i := 0; i+w < n; i++ {
			j := i + w
			p := ar[i][j-1]
			q := ar[i+1][j]
			dp[i][j] = INF
			ar[i][j] = p
			for k := p; k <= q && k+1 <= j; k++ {
				res := dp[i][k] + dp[k+1][j] + cost(i, k, j)
				if res < dp[i][j] {
					dp[i][j] = res
					ar[i][j] = k
				}
			}
		}
	}

	return dp[0][n-1]
}

// !dp[i][j]=min(dp[i][k-1]+dp[k+1][j]+cost(i,k,j)) (0<=i<=k<j<=n)
func KnuthYao2(n int, cost func(i, k, j int) int) int {
	dp := make([][]int, n+2)
	pos := make([][]int, n+2)
	for i := 0; i < n+2; i++ {
		dp[i] = make([]int, n+2)
		pos[i] = make([]int, n+2)
	}

	for i := 1; i <= n; i++ {
		for j := i; j <= n; j++ {
			dp[i][j] = INF
		}
	}

	for i := 1; i <= n; i++ {
		dp[i][i] = 0
		pos[i][i] = i
	}

	for len := 2; len <= n; len++ {
		for i := 1; i+len-1 <= n; i++ {
			j := i + len - 1
			for k := pos[i][j-1]; k <= pos[i+1][j]; k++ {
				res := dp[i][k-1] + dp[k+1][j] + cost(i-1, k-1, j-1)
				if res < dp[i][j] {
					dp[i][j] = res
					pos[i][j] = k
				}
			}
		}
	}

	return dp[1][n]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	xs, ys := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &xs[i], &ys[i])
	}

	res := KnuthYao1(n, func(i, k, j int) int {
		return (xs[k+1] - xs[i]) + (ys[k] - ys[j])
	})

	fmt.Fprintln(out, res)
}
