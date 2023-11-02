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
// O(n*need)
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
	dp := NewBoundedKnapsackRemovable(need, MOD)
	for i := 0; i < n; i++ {
		dp.Add(1, counts[i])
	}
	for i := 0; i < n; i++ {
		tmp := dp.Copy()
		tmp.Remove(1, counts[i])
		group := queryGroups[i]
		for _, queryWithId := range group {
			take, qid := queryWithId.take, queryWithId.qid
			res[qid] = tmp.Query(need - take)
		}
	}

	return res
}

// 可撤销多重背包，用于求解方案数/可行性.
type BoundedKnapsackRemovable struct {
	dp       []int
	mod      int
	maxValue int
	countSum int
}

// maxWeight: 需要的价值上限.
// mod: 模数，传入-1表示不需要取模.
func NewBoundedKnapsackRemovable(maxValue int, mod int) *BoundedKnapsackRemovable {
	dp := make([]int, maxValue+1)
	dp[0] = 1
	return &BoundedKnapsackRemovable{
		dp:       dp,
		mod:      mod,
		maxValue: maxValue,
	}
}

// 加入一个价值为value(value>0)的物品，数量为count.
// O(maxValue).
func (ks *BoundedKnapsackRemovable) Add(value, count int) {
	if value <= 0 {
		panic(fmt.Sprintf("value must be positive, but got %d", value))
	}
	ks.countSum += count * value
	maxJ := min(ks.countSum, ks.maxValue)
	if ks.mod == -1 {
		for j := value; j <= maxJ; j++ {
			ks.dp[j] += ks.dp[j-value]
		}
		for j := maxJ; j >= value*(count+1); j-- {
			ks.dp[j] -= ks.dp[j-value*(count+1)]
		}
	} else {
		for j := value; j <= maxJ; j++ {
			ks.dp[j] = (ks.dp[j] + ks.dp[j-value]) % ks.mod
		}
		for j := maxJ; j >= value*(count+1); j-- {
			ks.dp[j] = (ks.dp[j] - ks.dp[j-value*(count+1)]) % ks.mod
		}
	}
}

// 移除一个价值为value(value>0)的物品，数量为count.需要保证value物品存在.
func (ks *BoundedKnapsackRemovable) Remove(value, count int) {
	maxJ := min(ks.countSum, ks.maxValue)
	if ks.mod == -1 {
		for i := (count + 1) * value; i <= maxJ; i++ {
			ks.dp[i] += ks.dp[i-(count+1)*value]
		}
		for i := maxJ; i >= value; i-- {
			ks.dp[i] -= ks.dp[i-value]
		}
	} else {
		for i := (count + 1) * value; i <= maxJ; i++ {
			ks.dp[i] = (ks.dp[i] + ks.dp[i-(count+1)*value]) % ks.mod
		}
		for i := maxJ; i >= value; i-- {
			ks.dp[i] = (ks.dp[i] - ks.dp[i-value]) % ks.mod
		}
	}

	ks.countSum -= count * value
}

func (ks *BoundedKnapsackRemovable) Query(value int) int {
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

func (ks *BoundedKnapsackRemovable) Copy() *BoundedKnapsackRemovable {
	dp := append(ks.dp[:0:0], ks.dp...)
	return &BoundedKnapsackRemovable{
		dp:       dp,
		mod:      ks.mod,
		maxValue: ks.maxValue,
		countSum: ks.countSum,
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
