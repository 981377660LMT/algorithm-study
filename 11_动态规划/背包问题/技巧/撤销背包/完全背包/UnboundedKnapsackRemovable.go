package main

// 可撤销完全背包，用于求解方案数/可行性.
type UnboundedKnapsackRemovable struct {
	dp        []int
	maxWeight int
	mod       int
}

// maxWeight: 背包最大容量.
// mod: 模数，传入-1表示不需要取模.
func NewUnboundedKnapsackRemovable(maxWeight int, mod int) *UnboundedKnapsackRemovable {
	dp := make([]int, maxWeight+1)
	dp[0] = 1
	return &UnboundedKnapsackRemovable{
		dp:        dp,
		maxWeight: maxWeight,
		mod:       mod,
	}
}

// 添加一个重量为weight的物品.
func (ks *UnboundedKnapsackRemovable) Add(weight int) {
	if ks.mod == -1 {
		for i := weight; i <= ks.maxWeight; i++ {
			ks.dp[i] += ks.dp[i-weight]
		}
	} else {
		for i := weight; i <= ks.maxWeight; i++ {
			ks.dp[i] = (ks.dp[i] + ks.dp[i-weight]) % ks.mod
		}
	}
}

// 移除一个重量为weight的物品.需要保证weight物品存在.
func (ks *UnboundedKnapsackRemovable) Remove(weight int) {
	if ks.mod == -1 {
		for i := ks.maxWeight; i >= weight; i-- {
			ks.dp[i] -= ks.dp[i-weight]
		}
	} else {
		for i := ks.maxWeight; i >= weight; i-- {
			ks.dp[i] = (ks.dp[i] - ks.dp[i-weight]) % ks.mod
		}
	}
}

// 查询组成重量为weight的物品有多少种方案.
// !注意需要特判重量为0.
func (ks *UnboundedKnapsackRemovable) Query(weight int) int {
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

func (ks *UnboundedKnapsackRemovable) Copy() *UnboundedKnapsackRemovable {
	dp := append(ks.dp[:0:0], ks.dp...)
	return &UnboundedKnapsackRemovable{
		dp:        dp,
		maxWeight: ks.maxWeight,
		mod:       ks.mod,
	}
}
