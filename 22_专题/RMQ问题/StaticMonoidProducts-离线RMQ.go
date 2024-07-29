package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

// https://judge.yosupo.jp/problem/staticrmq
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][2]int32, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}
	res := StaticMonoidProducts(nums, func() S { return INF }, func(s1, s2 S) S { return min(s1, s2) }, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type S = int

// 离线RMQ.
func StaticMonoidProducts(arr []S, e func() S, op func(S, S) S, query [][2]int32) []S {
	n, q := int32(len(arr)), int32(len(query))
	res := make([]S, q)
	ids := make([][]int32, n)
	for qi := int32(0); qi < q; qi++ {
		start, end := query[qi][0], query[qi][1]
		if start >= end {
			res[qi] = e()
		} else if end <= start+32 {
			res[qi] = arr[start]
			for i := start + 1; i < end; i++ {
				res[qi] = op(res[qi], arr[i])
			}
		} else {
			end--
			k := bits.Len32(uint32(start^end)) - 1
			m := end >> k << k
			ids[m] = append(ids[m], qi)
		}
	}

	dp := make([]S, n+1)
	for m := int32(0); m < n; m++ {
		pos := ids[m]
		if len(pos) == 0 {
			continue
		}
		minA, maxB := m, m
		for _, qi := range pos {
			a, b := query[qi][0], query[qi][1]
			if a < minA {
				minA = a
			}
			if b > maxB {
				maxB = b
			}
		}
		dp[m] = e()
		for i := m; i > minA; i-- {
			dp[i-1] = op(arr[i-1], dp[i])
		}
		for i := m; i < maxB; i++ {
			dp[i+1] = op(dp[i], arr[i])
		}
		for _, qi := range pos {
			a, b := query[qi][0], query[qi][1]
			res[qi] = op(dp[a], dp[b])
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
