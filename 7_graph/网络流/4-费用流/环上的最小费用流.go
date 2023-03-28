// https://maspypy.github.io/library/flow/cyclic_bflow.hpp
// 环上的最小费用流(不需要费用流)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://atcoder.jp/contests/dwango2015-prelims/tasks/dwango2015_prelims_4

	// 给定一个长度为L的环,环上有n个城市,0为基准点。
	// !现在要求进入每个城市的人和出去的人的人数相同。
	// 每个城市距离基准点的距离为d[i],要出去的人数为out[i],可以进去的人数为in[i]。
	// 城市之间可以通过出租车来移动，1人跑距离为1的路程需要花费1元。
	// 创建一个程序，以确定移动总成本的最小可能值。
	// n<=1e5 L<=1e7
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, L int
	fmt.Fscan(in, &n, &L)
	pos := make([]int, n)
	supply := make([]int, n)
	for i := 0; i < n; i++ {
		var d, outFlow, inFlow int
		fmt.Fscan(in, &d, &outFlow, &inFlow)
		pos[i] = d
		supply[i] = outFlow - inFlow
	}
	cost := make([]int, n)
	for i := 0; i < n; i++ {
		pre := pos[i]
		var cur int
		if i+1 < n {
			cur = pos[i+1]
		} else {
			cur = L + pos[0]
		}
		cost[i] = cur - pre
	}
	fmt.Fprintln(out, CyclicFlow(supply, cost))
}

// CyclicFlow 环上的最小费用流
//  supply: 每个点的流量(正代表流出,负代表流入)
//  cost: i->i+1的费用
func CyclicFlow(supply []int, cost []int) int {
	n := len(supply)
	if n == 0 {
		return 0
	}
	cal := func(x int) int {
		res := abs(x) * cost[len(cost)-1]
		for i := 0; i < n-1; i++ {
			x += supply[i]
			res += abs(x) * cost[i]
		}
		return res
	}

	check := func(x int) bool { return cal(x) <= cal(x+1) }
	limit := 5
	for _, x := range supply {
		limit += max(x, 0)
	}

	ok, ng := limit, -limit
	for abs(ok-ng) > 1 {
		mid := (ok + ng) / 2
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return cal(ok)
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func max(x, y int) int {
	if x > y {
		return x
	}
	return y
}
