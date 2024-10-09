// 送快递
// F - Shipping
// https://atcoder.jp/contests/abc374/tasks/abc374_f
// 给定N个货物，每个货物下单时间为Ti.
// 快递一次最多可以运送K个货物.
// 每次出发送快递后，需要X天才能再次发送快递.
// !记不满度为所有货物的送达时间减去下单时间之和.
// 求最小不满度.
// K<=N<=100
//
// !dfs(i,time) 表示当前送完了前i个货物，此时时间为time 状态下的最小不满度.
// O(N^2*K)

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, K, X int
	fmt.Fscan(in, &N, &K, &X)
	T := make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &T[i])
	}

	sort.Ints(T)
	preSum := make([]int, N+1)
	for i := 1; i <= N; i++ {
		preSum[i] = preSum[i-1] + T[i-1]
	}

	memo := make([]map[int]int, N)
	for i := range memo {
		memo[i] = make(map[int]int, 2*N)
	}

	var dfs func(index, time int) int
	dfs = func(index, time int) int {
		if index == N {
			return 0
		}
		if v, ok := memo[index][time]; ok {
			return v
		}

		res := INF
		for j := 0; j < K; j++ {
			if index+j >= N {
				break
			}
			if time >= T[index+j] {
				curCost := time*(j+1) - (preSum[index+j+1] - preSum[index])
				res = min(res, dfs(index+j+1, time+X)+curCost)
			} else {
				res = min(res, dfs(index, T[index+j]))
			}
		}
		memo[index][time] = res
		return res
	}

	res := dfs(0, 0)
	fmt.Fprintln(out, res)
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
