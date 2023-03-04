// 如果点的范围很大,需要对xs和ys进行离散化,
// 查询坐标x为变为 `sort.SearchInts(sortedXs,x)`

package main

import (
	"math/bits"
	"sort"
)

type NumMatrix struct {
	seg *SegmentTree2D
}

func Constructor(matrix [][]int) NumMatrix {
	seg := NewSegmentTree2D(len(matrix), len(matrix[0]))
	for i := range matrix {
		for j := range matrix[i] {
			seg.AddPoint(i, j, matrix[i][j])
		}
	}

	seg.Build()
	return NumMatrix{seg}
}

func (this *NumMatrix) Update(row int, col int, val int) {
	this.seg.Set(row, col, val)
}

func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	return this.seg.Query(row1, col1, row2, col2)
}

/**
 * Your NumMatrix object will be instantiated and called as such:
 * obj := Constructor(matrix);
 * obj.Update(row,col,val);
 * param_2 := obj.SumRegion(row1,col1,row2,col2);
 */

// PointSetRangeSum2D

const INF int = 1e18

type E = int

func (*SegmentTree2D) e() E        { return 0 }
func (*SegmentTree2D) op(a, b E) E { return a + b }

func sortedSet(nums []int) (getRank func(int) int) {
	set := make(map[int]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted := make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	getRank = func(x int) int { return sort.SearchInts(sorted, x) }
	return
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

type SegmentTree2D struct {
	row, col int
	tree     []E
}

func NewSegmentTree2D(row, col int) *SegmentTree2D {
	res := &SegmentTree2D{}
	row_, col_ := 1<<uint(bits.Len(uint(row-1))), 1<<uint(bits.Len(uint(col-1)))

	tree := make([]E, (row_*col_)<<2)
	for i := range tree {
		tree[i] = res.e()
	}
	res.row, res.col, res.tree = row_, col_, tree
	return res
}

// 在Build之前调用, 为每个点赋值
func (st *SegmentTree2D) AddPoint(row, col int, value E) {
	st.tree[st.id(row+st.row, col+st.col)] = value
}

func (st *SegmentTree2D) Build() {
	for c := st.col; c < st.col<<1; c++ {
		for r := st.row - 1; r > 0; r-- {
			st.tree[st.id(r, c)] = st.op(st.tree[st.id(r<<1, c)], st.tree[st.id((r<<1)|1, c)])
		}
	}
	for r := 0; r < st.row<<1; r++ {
		for c := st.col - 1; c > 0; c-- {
			st.tree[st.id(r, c)] = st.op(st.tree[st.id(r, c<<1)], st.tree[st.id(r, (c<<1)|1)])
		}
	}
}

func (st *SegmentTree2D) Get(row, col int) E {
	return st.tree[st.id(row+st.row, col+st.col)]
}

func (st *SegmentTree2D) Set(row, col int, value E) {
	r, c := row+st.row, col+st.col
	st.tree[st.id(r, c)] = value
	for i := r >> 1; i > 0; i >>= 1 {
		st.tree[st.id(i, c)] = st.op(st.tree[st.id(i<<1, c)], st.tree[st.id((i<<1)|1, c)])
	}
	for ; r > 0; r >>= 1 {
		for j := c >> 1; j > 0; j >>= 1 {
			st.tree[st.id(r, j)] = st.op(st.tree[st.id(r, j<<1)], st.tree[st.id(r, (j<<1)|1)])
		}
	}
}

// (row1,col1) 到 (row2,col2) 闭区间的值.
//  0<=row1<row2<row, 0<=col1<col2<col
func (st *SegmentTree2D) Query(row1, col1, row2, col2 int) E {
	row2++
	col2++
	if row1 >= row2 || col1 >= col2 {
		return st.e()
	}
	res := st.e()
	row1 += st.row
	row2 += st.row
	col1 += st.col
	col2 += st.col
	for row1 < row2 {
		if row1&1 == 1 {
			res = st.op(res, st.query(row1, col1, col2))
			row1++
		}
		if row2&1 == 1 {
			row2--
			res = st.op(res, st.query(row2, col1, col2))
		}
		row1 >>= 1
		row2 >>= 1
	}
	return res
}

func (st *SegmentTree2D) query(r, c1, c2 int) E {
	res := st.e()
	for c1 < c2 {
		if c1&1 == 1 {
			res = st.op(res, st.tree[st.id(r, c1)])
			c1++
		}
		if c2&1 == 1 {
			c2--
			res = st.op(res, st.tree[st.id(r, c2)])
		}
		c1 >>= 1
		c2 >>= 1
	}
	return res
}

func (st *SegmentTree2D) id(r, c int) int { return ((r * st.col) << 1) + c }
