// 単調最小値DP (aka. 分割統治DP) 优化 offlineDp
// https://ei1333.github.io/library/dp/divide-and-conquer-optimization.hpp
// !用于高速化 dp[k][j]=min(dp[k-1][i]+f(i,j)) (0<=i<j<=n) => !将区间[0,n)分成k组的最小代价
//  如果f满足决策单调性 那么对转移的每一行，可以采用 monotoneminima 寻找最值点
//  O(kn^2)优化到O(knlogn)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

// CF833B-The Bakery
// https://www.luogu.com.cn/problem/CF833B
// 将一个数组分为k段，使得总价值最大。
// 一段区间的价值表示为区间内不同数字的个数。
// n=3e4,k<=50
//
// dp[i][j]=max{dp[i-1][k]+cost(k+1,j) 1<=k<j
// dp[i][j]意为前j个数被分成i段时的最大总价值.
//
// 决策单调性+主席树求cost
// !O(nklog^2n),TLE
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int32, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	S := NewStaticRangeCountDistinct(nums)

	res := divideAndConquerOptimization(k, n, func(i, j, _ int) int {
		count := S.Query(i, j)
		return -count
	})

	fmt.Fprintln(out, -res[k][n])
}

// !f(i,j,step): 左闭右开区间[i,j)的代价(0<=i<j<=n)
//
//	可选:step表示当前在第几组(1<=step<=k)
func divideAndConquerOptimization(k, n int, f func(i, j, step int) int) [][]int {
	dp := make([][]int, k+1)
	for i := range dp {
		dp[i] = make([]int, n+1)
		for j := range dp[i] {
			dp[i][j] = INF
		}
	}
	dp[0][0] = 0

	for k_ := 1; k_ <= k; k_++ {
		getCost := func(y, x int) int {
			if x >= y {
				return INF
			}
			return dp[k_-1][x] + f(x, y, k_)
		}
		res := monotoneminima(n+1, n+1, getCost)
		for j := 0; j <= n; j++ {
			dp[k_][j] = res[j][1]
		}
	}

	return dp
}

// 对每个 0<=i<H 求出 f(i,j) 取得最小值的 (j, f(i,j)) (0<=j<W)
func monotoneminima(H, W int, f func(i, j int) int) [][2]int {
	dp := make([][2]int, H) // dp[i] 表示第i行取到`最小值`的(索引,值)

	var dfs func(top, bottom, left, right int)
	dfs = func(top, bottom, left, right int) {
		if top > bottom {
			return
		}

		mid := (top + bottom) / 2
		index := -1
		res := 0
		for i := left; i <= right; i++ {
			tmp := f(mid, i)
			if index == -1 || tmp < res { // !less if get min
				index = i
				res = tmp
			}
		}
		dp[mid] = [2]int{index, res}
		dfs(top, mid-1, left, index)
		dfs(mid+1, bottom, index, right)
	}

	dfs(0, H-1, 0, W-1)
	return dp
}

// 持久化权值线段树求区间颜色种类数.
// !对于询问区间[i,j]，直接用[i,j]减去重复颜色的数量。
// 思路是记录每个数上一次出现的位置pre.
type StaticRangeCountDistinctOnline struct {
	n                   int32
	ptr                 int32
	root                int32
	left, right, preSum []int32
	roots               []int32
}

func NewStaticRangeCountDistinct(nums []int32) *StaticRangeCountDistinctOnline {
	n := int32(len(nums))
	maxLog := int32(bits.Len(uint(n))) + 1
	size := 2 * n * maxLog
	res := &StaticRangeCountDistinctOnline{
		n:      n,
		ptr:    1,
		root:   1,
		left:   make([]int32, size),
		right:  make([]int32, size),
		preSum: make([]int32, size),
		roots:  make([]int32, n+1),
	}
	last := make(map[int32]int32)
	for i := int32(0); i < n; i++ {
		res.makeRoot()
		x := nums[i]
		res.add(last[x], -1, res.roots[res.root-1], 0, n+1)
		last[x] = i + 1
		res.add(i+1, 1, res.roots[res.root-1], 0, n+1)
	}
	return res
}

func (sr *StaticRangeCountDistinctOnline) Query(start, end int) int {
	a, b := int32(start+1), int32(end+1)
	return int(sr.get(a, b, sr.roots[b-1], 0, sr.n+1))
}

func (sr *StaticRangeCountDistinctOnline) add(pos, delta, root, left, right int32) {
	sr.preSum[root] += delta
	if right-left > 1 {
		mid := (left + right) >> 1
		if pos < mid {
			sr.left[root] = sr.copy(sr.left[root])
			sr.add(pos, delta, sr.left[root], left, mid)
		} else {
			sr.right[root] = sr.copy(sr.right[root])
			sr.add(pos, delta, sr.right[root], mid, right)
		}
	}
}

func (sr *StaticRangeCountDistinctOnline) get(a, b, root, left, right int32) int32 {
	if right <= a || b <= left {
		return 0
	} else if a <= left && right <= b {
		return sr.preSum[root]
	} else {
		mid := (left + right) >> 1
		return sr.get(a, b, sr.left[root], left, mid) + sr.get(a, b, sr.right[root], mid, right)
	}
}

func (sr *StaticRangeCountDistinctOnline) copy(v int32) int32 {
	sr.left[sr.ptr] = sr.left[v]
	sr.right[sr.ptr] = sr.right[v]
	sr.preSum[sr.ptr] = sr.preSum[v]
	sr.ptr++
	return sr.ptr - 1
}

func (sr *StaticRangeCountDistinctOnline) makeRoot() {
	sr.roots[sr.root] = sr.copy(sr.roots[sr.root-1])
	sr.root++
}
