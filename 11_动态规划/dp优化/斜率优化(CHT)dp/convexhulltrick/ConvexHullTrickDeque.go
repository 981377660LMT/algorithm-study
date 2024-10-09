// ConvexHullTrickLiDeque(addMonotone.)
// 单调队列维护凸包

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// abc2018()
	abc373_f()
}

// 3221. Maximum Array Hopping Score II (ConvexHullTrick优化dp/斜率优化dp)
// https://leetcode.cn/problems/maximum-array-hopping-score-ii/description/
// dp[j]=max(dp[j],dp[i]+(j-i)*nums[j])
// !dp[j]=max(dp[j],-i*nums[j]+dp[i]+j*nums[j])
// !dp过程中将直线(-i,dp[i])不断加入到CHT中，查询时查询x=nums[j]时的最大值即可
func maxScore(nums []int) int64 {
	n := len(nums)
	dp := make([]int, n)
	cht := NewConvexHullTrickDeque(false, n)
	cht.AddLineMonotone(0, 0, -1)
	for j, v := range nums {
		cur, _ := cht.QueryMonotoneInc(v)
		dp[j] = cur + v*j
		cht.AddLineMonotone(-j, dp[j], -1)
	}
	return int64(dp[n-1])
}

func abc2018() {
	// https://atcoder.jp/contests/colopl2018-final-open/tasks/colopl2018_final_c
	// 对每个i 求 f(i,j)=a[j]+(j-i)^2 的最小值
	// 化简得 f(i,j)=(a[j]+j^2-2ij)+i^2
	// 其中j变化的函数是关于i的一次函数(直线) -i*2j+j^2+a[j]
	// !将这n条直线加入到CHT中,然后对每个i求最小值即可

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}

	cht := NewConvexHullTrickDeque(true, n)
	for i := 0; i < n; i++ {
		cht.AddLineMonotone(-2*i, i*i+A[i], -1)
	}

	for i := 0; i < n; i++ {
		res, _ := cht.QueryMonotoneInc(i)
		// res, _ := cht.Query(i)
		fmt.Fprintln(out, res+i*i)
	}
}

// F - Knapsack with Diminishing Values (凸函数代价的多重背包问题/按模分组)
// https://atcoder.jp/contests/abc373/tasks/abc373_f
// n种物品，重量wi，价值vi，无限个。
// 背包容量W，现在选物品放入背包，不超背包容量，价值最大。
// 当第i种物品放k个时，其价值为 k*vi − k^2。
// N<=3000,W<=3000,wi<=W,vi<=1e9
//
// !dp[i][j]表示前i种物品，总重量为j时的最大价值
// dp[i][j]=max(dp[i-1][j-k*wi]+k*vi-k^2) 0<=k<=j/wi
//
// !将j按模wi分组转移，令 j' = div * wi + mod，则
// !f[i] = max(g[j] + (i-j)*vi - (i-j)^2) 0<=j<=i.
// 分离ij得到，f[i] = 2i*j + (g[j] - j*vi - j^2) + (-i^2 + i*vi)
func abc373_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, W int
	fmt.Fscan(in, &N, &W)
	ws, vs := make([]int, N), make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &ws[i], &vs[i])
	}

	dp := make([]int, W+1)
	tmp := []int{}
	for i := 0; i < N; i++ {
		for mod := 0; mod < ws[i]; mod++ {
			tmp = tmp[:0]
			for pos := mod; pos <= W; pos += ws[i] {
				tmp = append(tmp, dp[pos])
			}

			cht := NewConvexHullTrickDeque(false, 0)
			for j, pos := 0, mod; j < len(tmp); j, pos = j+1, pos+ws[i] {
				if j > 0 {
					preMax, _ := cht.QueryMonotoneInc(j)
					dp[pos] = max(dp[pos], preMax-j*j+j*vs[i])
				}
				cht.AddLineMonotone(2*j, tmp[j]-j*vs[i]-j*j, -1)
			}
		}
	}

	fmt.Fprintln(out, dp[W])
}

const INF int = 1e18

type Line struct{ k, b, id int }
type ConvexHullTrickDeque struct {
	isMin bool
	dq    *Deque
}

func NewConvexHullTrickDeque(isMin bool, capacity int) *ConvexHullTrickDeque {
	if capacity < 0 {
		capacity = 0
	}
	return &ConvexHullTrickDeque{
		isMin: isMin,
		dq:    NewDeque(capacity),
	}
}

// 追加一条直线.需要保证斜率k是单调递增或者是单调递减的.
func (cht *ConvexHullTrickDeque) AddLineMonotone(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}

	line := &Line{k, b, id}
	if cht.dq.Empty() {
		cht.dq.AppendLeft(line)
		return
	}

	if cht.dq.Front().k <= k {
		if cht.dq.Front().k == k {
			if cht.dq.Front().b <= b {
				return
			}
			cht.dq.PopLeft()
		}
		for cht.dq.Len() >= 2 && cht.check(line, cht.dq.Front(), cht.dq.At(1)) {
			cht.dq.PopLeft()
		}
		cht.dq.AppendLeft(line)
	} else {
		if cht.dq.Back().k == k {
			if cht.dq.Back().b <= b {
				return
			}
			cht.dq.Pop()
		}
		for cht.dq.Len() >= 2 && cht.check(cht.dq.At(cht.dq.Len()-2), cht.dq.Back(), line) {
			cht.dq.Pop()
		}
		cht.dq.Append(line)
	}
}

// O(logn) 查询 k*x + b 的最小(大)值以及对应的直线id.
// 如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) Query(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	left, right := -1, cht.dq.Len()-1
	for left+1 < right {
		mid := (left + right) >> 1
		a, _ := cht.getY(cht.dq.At(mid), x)
		b, _ := cht.getY(cht.dq.At(mid+1), x)
		if a >= b {
			left = mid
		} else {
			right = mid
		}
	}

	res, lineId = cht.getY(cht.dq.At(right), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递增的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneInc(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Front(), x)
		b, _ := cht.getY(cht.dq.At(1), x)
		if a < b {
			break
		}
		cht.dq.PopLeft()
	}

	res, lineId = cht.getY(cht.dq.Front(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

// O(1) 查询 k*x + b 的最小(大)值以及对应的直线id.
//
//	需要保证x是单调递减的.
//	如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickDeque) QueryMonotoneDec(x int) (res, lineId int) {
	if cht.dq.Empty() {
		res, lineId = INF, -1
		if !cht.isMin {
			res = -INF
		}
		return
	}

	for cht.dq.Len() >= 2 {
		a, _ := cht.getY(cht.dq.Back(), x)
		b, _ := cht.getY(cht.dq.At(cht.dq.Len()-2), x)
		if a < b {
			break
		}
		cht.dq.Pop()
	}

	res, lineId = cht.getY(cht.dq.Back(), x)
	if !cht.isMin {
		res = -res
	}
	return
}

func (cht *ConvexHullTrickDeque) check(a, b, c *Line) bool {
	if b.b == a.b || c.b == b.b {
		return sign(b.k-a.k)*sign(c.b-b.b) >= sign(c.k-b.k)*sign(b.b-a.b)
	}
	return (b.k-a.k)*sign(c.b-b.b)*abs(c.b-b.b) >= (c.k-b.k)*sign(b.b-a.b)*abs(b.b-a.b)
}

func (cht *ConvexHullTrickDeque) getY(line *Line, x int) (int, int) {
	return line.k*x + line.b, line.id
}

func sign(x int) int {
	if x == 0 {
		return 0
	} else if x > 0 {
		return 1
	} else {
		return -1
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

type E = *Line
type Deque struct{ l, r []E }

func NewDeque(capacity int) *Deque { return &Deque{make([]E, 0, capacity/2), make([]E, 0, capacity/2)} }
func (q Deque) Empty() bool        { return len(q.l) == 0 && len(q.r) == 0 }
func (q Deque) Len() int           { return len(q.l) + len(q.r) }
func (q *Deque) AppendLeft(v E)    { q.l = append(q.l, v) }
func (q *Deque) Append(v E)        { q.r = append(q.r, v) }
func (q *Deque) PopLeft() (v E) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v E) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() E {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() E {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) E {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
