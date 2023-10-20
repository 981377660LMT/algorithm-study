package main

import (
	"bufio"
	"fmt"
	"os"
)

// P3932 浮游大陆的68号岛
// https://www.luogu.com.cn/problem/P3932
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 19260817

	var n, q int
	fmt.Fscan(in, &n, &q)
	positions := make([]int, n)
	for i := 0; i < n-1; i++ {
		fmt.Fscan(in, &positions[i+1])
		positions[i+1] += positions[i]
	}
	weights := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &weights[i])
	}

	D := DistSumWeighted(positions, weights)
	for i := 0; i < q; i++ {
		var to, start, end int
		fmt.Fscan(in, &to, &start, &end)
		start--
		to--
		fmt.Fprintln(out, D(start, end, to)%MOD)
	}
}

// 数轴上按照顺序分布着n个点,每个点的位置为positions[i],权重为weights[i].
// 点i到点j的距离定义为 `weights[i]*abs(positions[i]-positions[j])`.
// 求区间[start,end)内的所有点到点to的距离之和.
func DistSumWeighted(positions, weights []int) func(start, end, to int) int {
	preSum := make([]int, len(weights)+1)
	preMul := make([]int, len(weights)+1)
	for i, w := range weights {
		preSum[i+1] = preSum[i] + w
		preMul[i+1] = preMul[i] + w*positions[i]
	}

	cal := func(start, end, to int, onLeft bool) int {
		if start >= end {
			return 0
		}
		res1 := (preSum[end] - preSum[start]) * positions[to]
		res2 := preMul[end] - preMul[start]
		if !onLeft {
			res1, res2 = res2, res1
		}
		return res1 - res2
	}

	return func(start, end, to int) int {
		res1 := cal(start, min(end, to), to, true)
		res2 := cal(max(start, to), end, to, false)
		return res1 + res2
	}
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
