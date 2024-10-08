// ConvexHullTrickLichao
// 动态开点的李超线段树维护凸包

// https://ei1333.github.io/library/structure/convex-hull-trick/dynamic-li-chao-tree.hpp
// 追加直线/线段,查询k*x+b的最小值
// 如果要查询最大值,需要在插入直线时取反即(-k,-b),查询时返回-query(x)

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

}

// 3221. 最大数组跳跃得分 II
// https://leetcode.cn/problems/maximum-array-hopping-score-ii/solutions/2848972/python-convexhulltrickyou-hua-dpxie-lu-y-uy7t/
// dp[j]=max(dp[j],dp[i]+(j-i)*nums[j])
// !dp[j]=max(dp[j],-i*nums[j]+dp[i]+j*nums[j])
// !dp过程中将直线(-i,dp[i])不断加入到CHT中，查询时查询x=nums[j]时的最大值即可
func maxScore(nums []int) int64 {
	n := len(nums)
	dp := make([]int, n)
	cht := NewConvexHullTrickLichao(false, mins(nums...), maxs(nums...))
	cht.AddLine(0, 0, -1)
	for i, v := range nums {
		cur, _ := cht.Query(v)
		dp[i] = cur + v*i
		cht.AddLine(-i, dp[i], -1)
	}
	return int64(dp[n-1])
}

// 3282. 到达数组末尾的最大得分
// https://leetcode.cn/problems/reach-end-of-array-with-max-score/
// dp[j]=max(dp[j],dp[i]+(j-i)*nums[i])
// !dp[j]=max(dp[j],-i*nums[i]+dp[i]+j*nums[i])
// !dp过程中将直线(nums[i],dp[i]-i*nums[i])不断加入到CHT中，查询时查询x=j时的最大值即可
func findMaximumScore(nums []int) int64 {
	n := len(nums)
	dp := make([]int, n)
	cht := NewConvexHullTrickLichao(false, 0, n-1)
	cht.AddLine(0, 0, -1)
	for i, v := range nums {
		cur, _ := cht.Query(i)
		dp[i] = cur
		cht.AddLine(v, dp[i]-v*i, -1)
	}
	return int64(dp[n-1])
}

func abc2018() {
	// https://atcoder.jp/contests/colopl2018-final-open/tasks/colopl2018_final_c
	// 对每个i 求 f(i,j)=a[j]+(j-i)^2 的最小值
	// 化简得 f(i,j)=(a[j]+j^2-2ij)+i^2
	// 其中j变化的函数是关于i的一次函数(直线)
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

	cht := NewConvexHullTrickLichao(true, 0, n-1) // 查询自变量0<=x<n
	for i := 0; i < n; i++ {
		cht.AddLine(-2*i, i*i+A[i], i)
	}

	for i := 0; i < n; i++ {
		res, _ := cht.Query(i)
		fmt.Fprintln(out, res+i*i)
	}
}

const INF int = 1e18

type Line struct{ k, b, id int }
type LichaoNode struct {
	line *Line
	l, r *LichaoNode
}
type ConvexHullTrickLichao struct {
	isMin        bool
	lower, upper int
	root         *LichaoNode
}

// 根据待查询的自变量x的上下界[lower,upper]建立CHTLichao.
func NewConvexHullTrickLichao(isMin bool, lower, upper int) *ConvexHullTrickLichao {
	upper++
	return &ConvexHullTrickLichao{isMin: isMin, lower: lower, upper: upper}
}

// O(logN) 追加一条直线k*x+b, id为直线的编号.
func (cht *ConvexHullTrickLichao) AddLine(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := &Line{k, b, id}
	cht.root = cht.addLine(cht.root, line, cht.lower, cht.upper, cht.getY(line, cht.lower), cht.getY(line, cht.upper))
}

// O(logN^2) 追加一条左闭右开的线段[start,end)，所在直线k*x+b, id为线段的编号.
func (cht *ConvexHullTrickLichao) AddSegment(start, end, k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := &Line{k, b, id}
	cht.root = cht.addSegment(cht.root, line, start, end-1, cht.lower, cht.upper, cht.getY(line, cht.lower), cht.getY(line, cht.upper))
}

// O(logN) 查询k*x+b的最小/大值.如果不存在直线,返回的id为-1.
func (cht *ConvexHullTrickLichao) Query(x int) (res, id int) {
	res, id = cht.query(cht.root, cht.lower, cht.upper, x)
	if !cht.isMin {
		res = -res
	}
	return
}

func (cht *ConvexHullTrickLichao) addLine(t *LichaoNode, x *Line, l, r, xL, xR int) *LichaoNode {
	if t == nil {
		return &LichaoNode{line: x}
	}

	tL, tR := cht.getY(t.line, l), cht.getY(t.line, r)
	if tL <= xL && tR <= xR {
		return t
	} else if tL >= xL && tR >= xR {
		t.line = x
		return t
	} else {
		mid := (l + r) >> 1
		if mid == r {
			mid--
		}

		tM, xM := cht.getY(t.line, mid), cht.getY(x, mid)
		if tM > xM {
			t.line, x = x, t.line
			if xL >= tL {
				t.l = cht.addLine(t.l, x, l, mid, tL, tM)
			} else {
				t.r = cht.addLine(t.r, x, mid+1, r, tM+x.k, tR)
			}
		} else {
			if tL >= xL {
				t.l = cht.addLine(t.l, x, l, mid, xL, xM)
			} else {
				t.r = cht.addLine(t.r, x, mid+1, r, xM+x.k, xR)
			}
		}

		return t
	}
}

func (cht *ConvexHullTrickLichao) addSegment(t *LichaoNode, x *Line, a, b, l, r, xL, xR int) *LichaoNode {
	if r < a || b < l {
		return t
	}
	if a <= l && r <= b {
		y := &Line{x.k, x.b, x.id}
		return cht.addLine(t, y, l, r, xL, xR)
	}

	if t != nil {
		tL, tR := cht.getY(t.line, l), cht.getY(t.line, r)
		if tL <= xL && tR <= xR {
			return t
		}
	} else {
		t = &LichaoNode{line: &Line{0, INF, -1}}
	}

	mid := (l + r) >> 1
	if mid == r {
		mid--
	}
	xM := cht.getY(x, mid)
	t.l = cht.addSegment(t.l, x, a, b, l, mid, xL, xM)
	t.r = cht.addSegment(t.r, x, a, b, mid+1, r, xM+x.k, xR)
	return t
}

func (cht *ConvexHullTrickLichao) query(t *LichaoNode, l, r, x int) (res, id int) {
	if t == nil {
		res, id = INF, -1
		return
	}
	if l == r {
		res, id = cht.getY(t.line, x), t.line.id
		return
	}

	mid := (l + r) >> 1
	if mid == r {
		mid--
	}

	res, id = cht.getY(t.line, x), t.line.id
	if x <= mid {
		cand, candId := cht.query(t.l, l, mid, x)
		if cand < res {
			res, id = cand, candId
		}
	} else {
		cand, candId := cht.query(t.r, mid+1, r, x)
		if cand < res {
			res, id = cand, candId
		}
	}
	return
}

func (cht *ConvexHullTrickLichao) getY(line *Line, x int) int {
	return line.k*x + line.b
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
