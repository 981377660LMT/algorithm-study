// https://www.acwing.com/solution/content/48340/

// 设有 N 堆石子排成一排，其编号为 1，2，3，…，N。
// 每堆石子有一定的质量，可以用一个整数来描述，现在要将这 N 堆石子合并成为一堆。
// !每次只能合并`相邻`的两堆，合并的代价为这两堆石子的质量之和，
// 合并后与这两堆石子相邻的石子将和新堆相邻，合并时由于选择的顺序不同，合并的总代价也不相同。
// 例如有 4 堆石子分别为 1 3 5 2， 我们可以先合并 1、2 堆，代价为 4，得到 4 5 2，
// 又合并 1，2 堆，代价为 9，得到 9 2 ，再合并得到 11，总代价为 4+9+11=24；
// 如果第二步是先合并 2，3 堆，则代价为 7，得到 4 7，最后一次合并代价为 11，总代价为 4+7+11=22。
// !问题是：找出一种合理的方法，使总的代价最小，输出最小代价。
// n<=5000

//  dp(i, j)表示第i个石子到第j个石子的区间中所有石子合并的最小开销，最后求解答案是dp(0, n-1)
//  其中dp(i, j) = min { dp(i, k) + dp(k+1, j) + w(i, j) } (i <= k < j)
//  w(i, j) 是石子权值的区间和，满足四边形不等式，因此整个DP可以用四边形不等式的性质进行优化

package main

import (
	"bufio"
	"fmt"
	"os"
)

const INF int = 1e18

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

	preSum := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSum[i] = preSum[i-1] + nums[i-1]
	}

	res := KnuthYao1(n, func(i, k, j int) int {
		return preSum[j+1] - preSum[i]
	})
	fmt.Fprintln(out, res)
}

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
