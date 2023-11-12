package main

import "fmt"

func main() {
	K := NewBoundedKnapsackRemovable(100, -1)
	K.Add(1, 5)
	K.Add(2, 5)
	fmt.Println(K.Query(10))
	K.Add(3, 1)
	fmt.Println(K.Query(10))
	// K.Remove(1, 5)
	fmt.Println(K.Query(10))
	K.Remove(3, 1)
	fmt.Println(K.Query(10))
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

// !注意需要特判重量为0.
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
