package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P4141()
}

// https://atcoder.jp/contests/abc321/tasks/abc321_f
func SubsetSumWithAddAndErase() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 998244353

	var q, maxWeight int
	fmt.Fscan(in, &q, &maxWeight)

	K := NewKnapsack01Removable(maxWeight, MOD)
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "+" {
			var w int
			fmt.Fscan(in, &w)
			K.Add(w)
		} else {
			var w int
			fmt.Fscan(in, &w)
			K.Remove(w)
		}
		fmt.Fprintln(out, K.Query(maxWeight))
	}
}

// P4141 消失之物
// https://www.luogu.com.cn/problem/P4141
// 对每个物品i，在不选择i的情况下输出容量为1-m的方案数.
// n,m<=2000
func P4141() {
	solve := func(weights []int, maxCapacity int) [][]int {
		const MOD int = 10
		dp := NewKnapsack01Removable(maxCapacity, MOD)
		for _, w := range weights {
			dp.Add(w)
		}
		res := make([][]int, len(weights))
		for i := range res {
			tmp := dp.Copy()
			tmp.Remove(weights[i])
			res[i] = make([]int, maxCapacity+1)
			for j := 1; j <= maxCapacity; j++ {
				res[i][j] = tmp.Query(j)
			}
		}
		return res
	}

	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, maxCapacity int
	fmt.Fscan(in, &n, &maxCapacity)
	weights := make([]int, n)
	for i := range weights {
		fmt.Fscan(in, &weights[i])
	}

	res := solve(weights, maxCapacity)
	for _, row := range res {
		for _, v := range row[1:] {
			fmt.Fprint(out, v, "")
		}
		fmt.Fprintln(out)
	}
}

// 可撤销01背包,用于求解方案数/可行性.
type Knapsack01Removable struct {
	dp        []int
	maxWeight int
	mod       int
}

// maxWeight: 背包最大容量.
// mod: 模数，传入-1表示不需要取模.
func NewKnapsack01Removable(maxWeight int, mod int) *Knapsack01Removable {
	dp := make([]int, maxWeight+1)
	dp[0] = 1
	return &Knapsack01Removable{
		dp:        dp,
		maxWeight: maxWeight,
		mod:       mod,
	}
}

// 添加一个重量为weight的物品.
func (ks *Knapsack01Removable) Add(weight int) {
	if ks.mod == -1 {
		for i := ks.maxWeight; i >= weight; i-- {
			ks.dp[i] += ks.dp[i-weight]
		}
	} else {
		for i := ks.maxWeight; i >= weight; i-- {
			ks.dp[i] = (ks.dp[i] + ks.dp[i-weight]) % ks.mod
		}
	}
}

// 移除一个重量为weight的物品.需要保证weight物品存在.
func (ks *Knapsack01Removable) Remove(weight int) {
	if ks.mod == -1 {
		for i := weight; i <= ks.maxWeight; i++ {
			ks.dp[i] -= ks.dp[i-weight]
		}
	} else {
		for i := weight; i <= ks.maxWeight; i++ {
			ks.dp[i] = (ks.dp[i] - ks.dp[i-weight]) % ks.mod
		}
	}
}

// 查询组成重量为weight的物品有多少种方案.
func (ks *Knapsack01Removable) Query(weight int) int {
	if weight < 0 || weight > ks.maxWeight {
		return 0
	}
	if ks.mod == -1 {
		return ks.dp[weight]
	}
	if ks.dp[weight] < 0 {
		ks.dp[weight] += ks.mod
	}
	return ks.dp[weight]
}

func (ks *Knapsack01Removable) Copy() *Knapsack01Removable {
	dp := append(ks.dp[:0:0], ks.dp...)
	return &Knapsack01Removable{
		dp:        dp,
		maxWeight: ks.maxWeight,
		mod:       ks.mod,
	}
}
