// https://atcoder.jp/contests/dp/tasks/dp_e
// 超大容量01背包 -> 维度转换

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

	var n, limit int
	fmt.Fscan(in, &n, &limit)
	values := make([]int, n)
	weights := make([]int, n)
	for i := range values {
		fmt.Fscan(in, &weights[i], &values[i])
	}

	fmt.Fprintln(out, Knapsack2(values, weights, limit))
}

const INF int = 1e18

func Knapsack2(values []int, weights []int, maxCapacity int) int {
	valueSum := 0
	for _, v := range values {
		valueSum += v
	}
	dp := make([]int, valueSum+1)
	for i := range dp {
		dp[i] = INF
	}
	dp[0] = 0
	for i := 0; i < len(values); i++ {
		v, w := values[i], weights[i]
		for j := valueSum; j >= v; j-- {
			dp[j] = min(dp[j], dp[j-v]+w)
		}
	}
	res := 0
	for i := valueSum; i >= 0; i-- {
		if dp[i] <= maxCapacity {
			res = i
			break
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
