// https://atcoder.jp/contests/abc364/tasks/abc364_e
// !E - Maximum Glutton (abc364 E) dp的key和value互换
// n个食物，有咸度和甜度。安排一个吃的顺序，使得吃的食物尽可能多，
// 且一旦甜度之和>x或咸度之和>y就停下来不吃。
// n<=80
// x,y<=1e4
//
// !每个食物吃或不吃.
// !注意到x*y*n肯定不行，考虑技巧：优化dp状态定义(dp的key和value互换)
// !dp[i][j][k]表示前i个食物，选择j个食物，甜度和为k时的咸度和最小值
// 时间复杂度O(n^2*x)
//
// 由于答案需要查表，不建议写记忆化.

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e9

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, X, Y int
	fmt.Fscan(in, &N, &X, &Y)
	A, B := make([]int, N), make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &A[i], &B[i])
	}
	if X > Y {
		X, Y = Y, X
		A, B = B, A
	}

	makeDp := func() [][]int {
		res := make([][]int, N+1)
		for i := 0; i <= N; i++ {
			res[i] = make([]int, X+1)
			for j := 0; j <= X; j++ {
				res[i][j] = INF
			}
		}
		return res
	}

	dp, ndp := makeDp(), makeDp() // 滚动数组(双缓冲,double buffering)
	dp[0][0] = 0
	for i := 0; i < N; i++ {
		for j := 0; j <= i; j++ {
			for k := 0; k <= X; k++ {
				ndp[j][k] = min(ndp[j][k], dp[j][k])
				if k+A[i] <= X {
					ndp[j+1][k+A[i]] = min(ndp[j+1][k+A[i]], dp[j][k]+B[i])
				}
			}
		}
		dp, ndp = ndp, dp
	}

	for j := N; j >= 0; j-- {
		for k := 0; k <= X; k++ {
			if dp[j][k] <= Y {
				fmt.Fprintln(out, min(j+1, N))
				return
			}
		}
	}

	fmt.Fprintln(out, 1)
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
