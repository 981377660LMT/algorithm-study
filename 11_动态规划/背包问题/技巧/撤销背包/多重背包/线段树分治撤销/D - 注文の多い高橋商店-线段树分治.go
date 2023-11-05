package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 1e9 + 7

// https://atcoder.jp/contests/arc028/tasks/arc028_4
// D - 注文の多い高橋商店
// 有n种商品，第i种商品有counts[i]个。给定购买的商品数need。
// 给定q次查询，每次查询给定i和x，求第i种商品恰好拿走x个的情况下，在n种物品中一共拿走need个物品的方案数。
// n,need,counts[i]<=2e3
// q<=5e5
// O(n*need*logn)
// !由于背包都可撤销(具有 remove 接口)，所以这种方法并不是最优的
func ARC028D(counts []int, need int, queries [][2]int) []int {
	type entry struct{ take, qid int }
	n := len(counts)
	q := len(queries)
	queryGroups := make([][]entry, n)
	for i, query := range queries {
		pos, take := query[0], query[1]
		if take > need || take > counts[pos] {
			continue
		}
		queryGroups[pos] = append(queryGroups[pos], entry{take, i})
	}

	res := make([]int, q)
	initState := NewBoundedKnapsack(need, MOD)
	MutateWithoutOneCopy(
		initState,
		0, n,
		func(state *S) *S {
			return state.Copy()
		},
		func(state *S, index int) {
			state.Add(1, counts[index])
		},
		func(state *S, index int) {
			// 回答 index 处的查询
			group := queryGroups[index]
			for _, queryWithId := range group {
				take, qid := queryWithId.take, queryWithId.qid
				res[qid] = state.Query(need - take)
			}
		},
	)

	return res
}

type S = BoundedKnapsack

// 线段树分治的特殊情形.
func MutateWithoutOneCopy(
	initState *S,
	start, end int,
	copy func(state *S) *S,
	mutate func(state *S, index int),
	query func(state *S, index int),
) {
	var dfs func(state *S, curStart, curEnd int)
	dfs = func(state *S, curStart, curEnd int) {
		if curEnd == curStart+1 {
			query(state, curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		leftCopy := copy(state)
		for i := curStart; i < mid; i++ {
			mutate(leftCopy, i)
		}
		dfs(leftCopy, mid, curEnd)

		rightCopy := copy(state)
		for i := mid; i < curEnd; i++ {
			mutate(rightCopy, i)
		}
		dfs(rightCopy, curStart, mid)
	}

	dfs(initState, start, end)
}

// 多重背包求方案数.
type BoundedKnapsack struct {
	dp       []int
	mod      int
	maxValue int
	maxJ     int
}

// maxWeight: 需要的价值上限.
// mod: 模数，传入-1表示不需要取模.
func NewBoundedKnapsack(maxValue int, mod int) *BoundedKnapsack {
	dp := make([]int, maxValue+1)
	dp[0] = 1
	return &BoundedKnapsack{
		dp:       dp,
		mod:      mod,
		maxValue: maxValue,
	}
}

// 加入一个价值为value(value>0)的物品，数量为count.
// O(maxValue).
func (ks *BoundedKnapsack) Add(value, count int) {
	if value <= 0 {
		panic(fmt.Sprintf("value must be positive, but got %d", value))
	}
	ks.maxJ = min(ks.maxJ+count*value, ks.maxValue)
	if ks.mod == -1 {
		for j := value; j <= ks.maxJ; j++ {
			ks.dp[j] += ks.dp[j-value]
		}
		for j := ks.maxJ; j >= value*(count+1); j-- {
			ks.dp[j] -= ks.dp[j-value*(count+1)]
		}
	} else {
		for j := value; j <= ks.maxJ; j++ {
			ks.dp[j] = (ks.dp[j] + ks.dp[j-value]) % ks.mod
		}
		for j := ks.maxJ; j >= value*(count+1); j-- {
			ks.dp[j] = (ks.dp[j] - ks.dp[j-value*(count+1)]) % ks.mod
		}
	}
}

func (ks *BoundedKnapsack) Query(value int) int {
	if value < 0 || value > ks.maxValue {
		return 0
	}
	if ks.mod == -1 {
		return ks.dp[value]
	}
	if ks.dp[value] < 0 {
		ks.dp[value] += ks.mod
	}
	return ks.dp[value]
}

func (ks *BoundedKnapsack) Copy() *BoundedKnapsack {
	dp := append(ks.dp[:0:0], ks.dp...)
	return &BoundedKnapsack{
		dp:       dp,
		mod:      ks.mod,
		maxValue: ks.maxValue,
		maxJ:     ks.maxJ,
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, need, q int
	fmt.Fscan(in, &n, &need, &q)
	counts := make([]int, n)
	for i := range counts {
		fmt.Fscan(in, &counts[i])
	}
	queries := make([][2]int, q)
	for i := range queries {
		var qi, qv int
		fmt.Fscan(in, &qi, &qv)
		qi--
		queries[i] = [2]int{qi, qv}
	}

	res := ARC028D(counts, need, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}
