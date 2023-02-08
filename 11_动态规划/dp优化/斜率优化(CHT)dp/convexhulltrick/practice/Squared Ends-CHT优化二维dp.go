// https://csacademy.com/contest/round-70/task/squared-ends/
// 给定一个数组,分成k个子数组,
// 子数组[A[l],...,A[r]]的代价为(A[l]-A[r])^2
// 最小化所有子数组的代价和
// !n<=1e4 k<=100 1<=A[i]<=1e9
// CHT或者分治优化dp都可以解决这个问题

// dp[k][j]表示将[0,j]区间分成k个子数组的最小代价和
// dp[k][j]=min(dp[k-1][i]+(A[i+1]-A[j])^2) (0<=i<j<n)
// 即 ndp[j] = min(dp[i]+(A[i+1]-A[j])^2) (0<=i<j<n)
// !变形得到 dp[j] = min(-2*A[i+1]*A[j] + (A[i+1]^2+dp[i]) + A[j]^2) (0<=i<j<n)
// 可以用CHT优化dp

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}

	dp := make([]int, n) // 1组
	for i := 0; i < n; i++ {
		dp[i] = (A[i] - A[0]) * (A[i] - A[0])
	}

	for k_ := 1; k_ < k; k_++ {
		ndp := make([]int, n)
		cht := NewConvexHullTrickLichao(true, 1, 1e9) // query自变量A[j]范围[1,1e9]
		for j := k_; j < n; j++ {
			cht.AddLine(-2*A[j], A[j]*A[j]+dp[j-1], j) // 注意这里是dp[j-1]
			best, _ := cht.Query(A[j])
			ndp[j] = best + A[j]*A[j]
		}
		dp = ndp
	}

	fmt.Fprintln(out, dp[n-1])
}

const INF int = 1e18

type Line struct{ k, b, id int }
type LichaoNode struct {
	line Line
	l, r *LichaoNode
}
type ConvexHullTrickLichao struct {
	isMin        bool
	lower, upper int
	root         *LichaoNode
}

// 根据待查询的自变量x的上下界[lower,upper]建立CHTLichao.
func NewConvexHullTrickLichao(isMin bool, lower, upper int) *ConvexHullTrickLichao {
	return &ConvexHullTrickLichao{isMin: isMin, lower: lower, upper: upper}
}

// O(logN) 追加一条直线k*x+b, id为直线的编号.
func (cht *ConvexHullTrickLichao) AddLine(k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := Line{k, b, id}
	cht.root = cht.addLine(cht.root, line, cht.lower, cht.upper, cht.getY(line, cht.lower), cht.getY(line, cht.upper))
}

// O(logN^2) 追加一条左闭右开的线段[start,end)，所在直线k*x+b, id为线段的编号.
func (cht *ConvexHullTrickLichao) AddSegment(start, end, k, b, id int) {
	if !cht.isMin {
		k, b = -k, -b
	}
	line := Line{k, b, id}
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

func (cht *ConvexHullTrickLichao) addLine(t *LichaoNode, x Line, l, r, xL, xR int) *LichaoNode {
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

func (cht *ConvexHullTrickLichao) addSegment(t *LichaoNode, x Line, a, b, l, r, xL, xR int) *LichaoNode {
	if r < a || b < l {
		return t
	}
	if a <= l && r <= b {
		y := Line{x.k, x.b, x.id}
		return cht.addLine(t, y, l, r, xL, xR)
	}

	if t != nil {
		tL, tR := cht.getY(t.line, l), cht.getY(t.line, r)
		if tL <= xL && tR <= xR {
			return t
		}
	} else {
		t = &LichaoNode{line: Line{0, INF, -1}}
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

func (cht *ConvexHullTrickLichao) getY(line Line, x int) int {
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
