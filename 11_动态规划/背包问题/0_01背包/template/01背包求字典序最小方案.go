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

	var n, v int
	fmt.Fscan(in, &n, &v)
	values := make([]int, n)
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i], &values[i])
	}

	res := ZeroOneKnapsackLexicographicallySmallest(values, weights, v)
	for _, i := range res {
		fmt.Fprint(out, i+1, " ")
	}
}

// 01背包求字典序最小的方案
// 倒序遍历物品，同时用 pre 数组记录转移来源，这样跑完 DP 后，从第一个物品开始即可得到字典序最小的方案
// https://www.acwing.com/problem/content/description/12/
func ZeroOneKnapsackLexicographicallySmallest(values, weights []int, maxW int) (res []int) {
	n := len(values)
	dp := make([]int, maxW+1) // fill
	dp[0] = 0
	pre := make([][]int, n)
	for i := n - 1; i >= 0; i-- {
		pre[i] = make([]int, maxW+1)
		for j := range pre[i] {
			pre[i][j] = j // 注意：<w 的转移来源也要标上！
		}
		v, w := values[i], weights[i]
		for j := maxW; j >= w; j-- {
			if dp[j-w]+v >= dp[j] { // !注意这里要取等号，从而保证尽可能地从字典序最小的方案转移过来
				dp[j] = dp[j-w] + v
				pre[i][j] = j - w
			}
		}
	}
	for i, j := 0, maxW; i < n; {
		if pre[i][j] == j { // &&  weights[i] > 0      考虑重量为 0 的情况，必须都选上
			i++
		} else {
			res = append(res, i) // 下标从 1 开始
			j = pre[i][j]
			i++ // 完全背包的情况，这行去掉
		}
	}
	return
}
