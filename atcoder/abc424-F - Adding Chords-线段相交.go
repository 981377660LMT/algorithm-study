// https://atcoder.jp/contests/abc424/tasks/abc424_f
// F - Adding Chords (添加弦)
//
// 问题描述
//
// 在一个圆周上，有 N 个等距分布的点，按顺时针方向编号为 1, 2, ..., N。
//
// 你需要按顺序处理 Q 个查询，每个查询的形式如下：
//
// 尝试画一条连接点 A_i 和点 B_i 的弦。但是，如果这条弦与任何已经画好的弦相交，则不画。
// 这里保证所有 2Q 个端点 A_1, ..., A_Q, B_1, ..., B_Q 都是互不相同的。
//
// 对于每条弦，请回答它是否被画上了。
// !圆上弦相交问题经典结论：两条弦相交当且仅当它们的四个端点在圆周上是交错排列的——比如按照顺时针方向，出现ACBD或者CADB这样的交错顺序。
// !故原问题等价于有Q个区间 ，每次询问是否与已有区间相交。
// https://atcoder.jp/contests/abc424/editorial/13900

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

	var n, q int
	fmt.Fscan(in, &n, &q)
	A, B := make([]int, q), make([]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &A[i], &B[i])
		A[i]--
		B[i]--
	}

	segMin := NewSegmentTreeGeneric(n, func(int) int { return INF }, func() int { return INF }, min)
	segMax := NewSegmentTreeGeneric(n, func(int) int { return -INF }, func() int { return -INF }, max)
	for i := 0; i < q; i++ {
		start, end := A[i], B[i]
		if start > end {
			start, end = end, start
		}

		min_, max_ := segMin.Query(start+1, end), segMax.Query(start+1, end)
		if min_ < start || max_ > end {
			fmt.Fprintln(out, "No")
		} else {
			fmt.Fprintln(out, "Yes")
			segMin.Set(end, start)
			segMax.Set(start, end)
		}
	}
}

const INF int = 1e18

// SegmentTreeGeneric
type SegmentTreeGeneric[E any] struct {
	n, size int
	seg     []E
	e       func() E
	op      func(a, b E) E
}

func NewSegmentTreeGeneric[E any](
	n int, f func(int) E,
	e func() E, op func(a, b E) E,
) *SegmentTreeGeneric[E] {
	res := &SegmentTreeGeneric[E]{e: e, op: op}
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
func (st *SegmentTreeGeneric[E]) Get(index int) E {
	if index < 0 || index >= st.n {
		return st.e()
	}
	return st.seg[index+st.size]
}
func (st *SegmentTreeGeneric[E]) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = st.op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *SegmentTreeGeneric[E]) Update(index int, value E) {
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
func (st *SegmentTreeGeneric[E]) Query(start, end int) E {
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
func (st *SegmentTreeGeneric[E]) QueryAll() E { return st.seg[1] }
func (st *SegmentTreeGeneric[E]) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
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
