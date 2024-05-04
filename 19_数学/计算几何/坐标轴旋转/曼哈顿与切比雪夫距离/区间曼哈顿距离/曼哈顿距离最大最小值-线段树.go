// 区间曼哈顿距离最大最小值-线段树
// 区间最大曼哈顿距离/区间最小曼哈顿距离

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	cf1093g()
}

// https://leetcode.cn/problems/minimize-manhattan-distances/description/
// 3102. 最小化曼哈顿距离
// 给你一个下标从 0 开始的数组 points ，它表示二维平面上一些点的整数坐标，其中 points[i] = [xi, yi] 。
// 两点之间的距离定义为它们的曼哈顿距离。
// 请你恰好移除一个点，返回移除后任意两点之间的 最大 距离可能的 最小 值。
//
// !移除一个点，可以变成将这个点设为另一个点，消除了这个点的影响.
func minimumDistance(points [][]int) int {
	n := len(points)

	const K = 2
	type E = [1 << K]int
	fromElement := func(point []int) E {
		res := E{}
		for i := 0; i < 1<<K; i++ {
			for j := 0; j < K; j++ {
				if i&(1<<j) > 0 {
					res[i] += point[j]
				} else {
					res[i] -= point[j]
				}
			}
		}
		return res
	}
	e := func() E {
		res := E{}
		res[0] = INF
		return res
	}
	op := func(a, b E) E {
		if a[0] == INF {
			return b
		}
		if b[0] == INF {
			return a
		}
		for i := range a {
			a[i] = max(a[i], b[i])
		}
		return a
	}
	seg := NewSegmentTree(
		n, func(i int) E { return fromElement(points[i]) },
		e, op,
	)

	set := func(index int, point []int) {
		seg.Set(index, fromElement(point))
	}
	query := func(start, end int) int {
		tmp := seg.Query(start, end)
		res := 0
		for i := 0; i < 1<<K; i++ {
			res = max(res, tmp[i]+tmp[(1<<K)-1-i])
		}
		return res
	}
	_ = query
	queryAll := func() int {
		tmp := seg.QueryAll()
		res := 0
		for i := 0; i < 1<<K; i++ {
			res = max(res, tmp[i]+tmp[(1<<K)-1-i])
		}
		return res
	}

	res := INF
	x1, y1 := points[0][0], points[0][1]
	xn, yn := points[n-1][0], points[n-1][1]
	for i := 0; i < n; i++ {
		if i > 0 {
			set(i, []int{x1, y1})
		} else {
			set(i, []int{xn, yn})
		}
		res = min(res, queryAll())
		set(i, points[i])
	}
	return res
}

// CF1093G-Multidimensional Queries (曼哈顿距离最大值)
// https://www.luogu.com.cn/problem/CF1093G
// 给定n个k维点(k<=5).
// 1 i x1 x2 ... xk: 第i个点的坐标更新为(x1,x2,...,xk);
// 2 start end ：查询区间[start,end)内最大的两点间曼哈顿距离.
//
// !习惯性的把曼哈顿距离的绝对值拆出来，用二进制表示
// 31 的二进制表示是 11111，表示 5 维的一个点的坐标加入的正负情况都为正
// 即 x[0] - y[0] + x[1] - y[1] + x[2] - y[2] + x[3] - y[3] + x[4] - y[4]
// 29 的二进制表示是 11101，表示 5 维的一个点的坐标加入的正负情况为正、负、正、正、正
// 即 x[0] - y[0] + x[1] - y[1] + x[2] - y[2] - x[3] + y[3] + x[4] - y[4]
// 那么要求的就是 max(f[0]+f[31],f[1]+f[30],f[2]+f[29],...,f[15]+f[16])
// 用线段树维护最大值即可.
func cf1093g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	points := make([][]int, n)
	for i := 0; i < n; i++ {
		points[i] = make([]int, k)
		for j := 0; j < k; j++ {
			fmt.Fscan(in, &points[i][j])
		}
	}

	// PointSetRangeMax
	type E = [32]int // 长为1<<k
	fromElement := func(point []int) E {
		res := E{}
		for i := 0; i < 1<<k; i++ {
			for j := 0; j < k; j++ {
				if i&(1<<j) > 0 {
					res[i] += point[j]
				} else {
					res[i] -= point[j]
				}
			}
		}
		return res
	}
	e := func() E {
		res := E{}
		res[0] = INF
		return res
	}
	op := func(a, b E) E {
		if a[0] == INF {
			return b
		}
		if b[0] == INF {
			return a
		}
		for i := range a {
			a[i] = max(a[i], b[i])
		}
		return a
	}
	seg := NewSegmentTree(
		n, func(i int) E { return fromElement(points[i]) },
		e, op,
	)

	set := func(index int, point []int) {
		seg.Set(index, fromElement(point))
	}
	query := func(start, end int) int {
		tmp := seg.Query(start, end)
		res := 0
		for i := 0; i < 1<<k; i++ {
			res = max(res, tmp[i]+tmp[(1<<k)-1-i])
		}
		return res
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var pos int
			fmt.Fscan(in, &pos)
			pos--
			point := make([]int, k)
			for j := 0; j < k; j++ {
				fmt.Fscan(in, &point[j])
			}
			set(pos, point)
		} else {
			var start, end int
			fmt.Fscan(in, &start, &end)
			start--
			fmt.Fprintln(out, query(start, end))
		}
	}
}

const INF int = 1e18

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

const INF32 int32 = 1 << 30

type SegmentTree[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTree[E any](
	n int, f func(int) E,
	e func() E, op func(a, b E) E,
) *SegmentTree[E] {
	res := &SegmentTree[E]{e: e, op: op}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = res.e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = res.op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *SegmentTree[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTree[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTree[E]) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = st.op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *SegmentTree[E]) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return st.e()
	}
	leftRes, rightRes := st.e(), st.e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = st.op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = st.op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return st.op(leftRes, rightRes)
}
func (st *SegmentTree[E]) QueryAll() E { return st.seg[1] }
func (st *SegmentTree[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}
