// https://www.luogu.com.cn/problem/P6808
// 你需要修改序列中的一个数 P 为 Q，使得尽可能多的整数能够被表示出来。
// 如果有多种方案，则输出的 P 尽可能小。
// P 最小时如有多种方案，则输出的 Q 尽可能小。
// n <= 100, nums[i] <= 7000

// 求P:
// 记去掉每个数后,组成i的方案数为dp[i],可用撤销背包求出.
// !将这个数修改为一个很大的数后，可以表达数一定是2*dp[i]+1.那么只需要考虑最大的dp[i]对应的P.
// 求Q:
// !用一个 bitset 维护当前能组成的数的集合。加入一个数如果使其与原来的 bitset 没有交集，这个 Q 就是合法的

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int = 1e9 + 7

func Candies(nums []int) (p, q int) {
	sum := 0
	for _, v := range nums {
		sum += v
	}
	dp := NewKnapsack01Removable(sum, MOD)
	for _, v := range nums {
		dp.Add(v)
	}

	bestRes, bestI := -1, -1
	for i, v := range nums {
		tmp := dp.Copy()
		tmp.Remove(v)
		count := 0
		for j := 1; j <= sum; j++ {
			if tmp.Query(j) != 0 {
				count++
			}
		}
		if count > bestRes || (count == bestRes && v < nums[bestI]) {
			bestRes = count
			bestI = i
		}
	}
	p = nums[bestI]

	offset := (sum - p)
	bitset := _NewBS(offset * 2)
	bitset.Set(offset)

	for i, v := range nums {
		if i == bestI {
			continue
		}
		// bitset = bitset | (bitset << v) | (bitset >> v)
		tmp := bitset.Copy().Lsh(v)
		bitset.IOr(tmp)
		tmp.Rsh(2 * v)
		bitset.IOr(tmp)
	}

	for i := 1; i <= sum+1; i++ {
		if !bitset.Has(i + offset) {
			q = i
			break
		}
	}

	return
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

func (b _BS) And(c _BS) _BS {
	res := make(_BS, len(b))
	for i, v := range b {
		res[i] = v & c[i]
	}
	return res
}

func (b _BS) OnesCount() (c int) {
	for _, v := range b {
		c += bits.OnesCount(v)
	}
	return
}

type _BS []uint

func _NewBS(n int) _BS { return make(_BS, n>>6+1) } // (n+64-1)>>6

func (b _BS) Has(p int) bool { return b[p>>6]&(1<<(p&63)) != 0 } // get
func (b _BS) Flip(p int)     { b[p>>6] ^= 1 << (p & 63) }
func (b _BS) Set(p int)      { b[p>>6] |= 1 << (p & 63) }  // 置 1
func (b _BS) Reset(p int)    { b[p>>6] &^= 1 << (p & 63) } // 置 0

func (b _BS) Copy() _BS {
	res := make(_BS, len(b))
	copy(res, b)
	return res
}

func (b _BS) Lsh(k int) _BS {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	if offset == 0 {
		// Fast path
		copy(b[shift:], b)
	} else {
		for i := len(b) - 1; i > shift; i-- {
			b[i] = b[i-shift]<<offset | b[i-shift-1]>>(64-offset)
		}
		b[shift] = b[0] << offset
	}
	for i := 0; i < shift; i++ {
		b[i] = 0
	}
	return b
}

// 右移 k 位
func (b _BS) Rsh(k int) _BS {
	if k == 0 {
		return b
	}
	shift, offset := k>>6, k&63
	if shift >= len(b) {
		for i := range b {
			b[i] = 0
		}
		return b
	}
	lim := len(b) - 1 - shift
	if offset == 0 {
		// Fast path
		copy(b, b[shift:])
	} else {
		for i := 0; i < lim; i++ {
			b[i] = b[i+shift]>>offset | b[i+shift+1]<<(64-offset)
		}
		// 注意：若前后调用 lsh 和 rsh，需要注意超出 n 的范围的 1 对结果的影响（如果需要，可以把范围开大点）
		b[lim] = b[len(b)-1] >> offset
	}
	for i := lim + 1; i < len(b); i++ {
		b[i] = 0
	}
	return b
}

// 将 c 的元素合并进 b
func (b _BS) IOr(c _BS) {
	for i, v := range c {
		b[i] |= v
	}
}

func (b _BS) Or(c _BS) _BS {
	res := make(_BS, len(b))
	for i, v := range b {
		res[i] = v | c[i]
	}
	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	p, q := Candies(nums)
	fmt.Fprintln(out, p, q)
}
