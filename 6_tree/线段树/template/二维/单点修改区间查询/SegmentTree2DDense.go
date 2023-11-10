// 如果点的范围很大,需要对xs和ys进行离散化,
// 查询坐标x为变为 `sort.SearchInts(sortedXs,x)`

// API:
// 1. NewSegmentTree2DDense(row,col int) *SegmentTree2DDense
// 3. (st *SegmentTree2DDense) Build(f func(r,c int) E)
// 4. (st *SegmentTree2DDense) Get(row,col int) E
// 5. (st *SegmentTree2DDense) Set(row,col int,value E)
// 6. (st *SegmentTree2DDense) Query(row1,row2,col1,col2 int) E

package main

import (
	"math/bits"
	"sort"
)

type NumMatrix struct {
	seg *SegmentTree2DDense
}

func Constructor(matrix [][]int) NumMatrix {
	seg := NewSegmentTree2DDense(len(matrix), len(matrix[0]))
	seg.Build(func(r, c int) E { return matrix[r][c] })
	return NumMatrix{seg}
}

func (this *NumMatrix) Update(row int, col int, val int) {
	this.seg.Set(row, col, val)
}

func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	return this.seg.Query(row1, row2+1, col1, col2+1)
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

func (*SegmentTree2DDense) e() E        { return 0 }
func (*SegmentTree2DDense) op(a, b E) E { return a + b }

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

type SegmentTree2DDense struct {
	Row, Col             int
	rowOffset, colOffset int
	tree                 []E
	unit                 E
}

func NewSegmentTree2DDense(row, col int) *SegmentTree2DDense {
	res := &SegmentTree2DDense{}
	res.unit = res.e()
	rowOffset, colOffset := 1<<uint(bits.Len(uint(row-1))), 1<<uint(bits.Len(uint(col-1)))
	tree := make([]E, (rowOffset*colOffset)<<2)
	for i := range tree {
		tree[i] = res.unit
	}
	res.Row, res.Col = row, col
	res.rowOffset, res.colOffset, res.tree = rowOffset, colOffset, tree
	return res
}

func (st *SegmentTree2DDense) Build(f func(r, c int) E) {
	ROW, COL := st.Row, st.Col
	rowOffset, colOffset := st.rowOffset, st.colOffset
	for r := 0; r < ROW; r++ {
		for c := 0; c < COL; c++ {
			st.tree[st._id(r+rowOffset, c+colOffset)] = f(r, c)
		}
	}

	for c := colOffset; c < colOffset<<1; c++ {
		for r := rowOffset - 1; r > 0; r-- {
			st.tree[st._id(r, c)] = st.op(st.tree[st._id(r<<1, c)], st.tree[st._id((r<<1)|1, c)])
		}
	}
	for r := 0; r < rowOffset<<1; r++ {
		for c := colOffset - 1; c > 0; c-- {
			st.tree[st._id(r, c)] = st.op(st.tree[st._id(r, c<<1)], st.tree[st._id(r, (c<<1)|1)])
		}
	}
}

func (st *SegmentTree2DDense) Get(row, col int) E {
	return st.tree[st._id(row+st.rowOffset, col+st.colOffset)]
}

func (st *SegmentTree2DDense) Set(row, col int, value E) {
	r, c := row+st.rowOffset, col+st.colOffset
	st.tree[st._id(r, c)] = value
	for i := r >> 1; i > 0; i >>= 1 {
		st.tree[st._id(i, c)] = st.op(st.tree[st._id(i<<1, c)], st.tree[st._id((i<<1)|1, c)])
	}
	for ; r > 0; r >>= 1 {
		for j := c >> 1; j > 0; j >>= 1 {
			st.tree[st._id(r, j)] = st.op(st.tree[st._id(r, j<<1)], st.tree[st._id(r, (j<<1)|1)])
		}
	}
}

// `[row1,row2) x [col1,col2)`
//
//	0<=row1<row2<=row, 0<=col1<col2<=col
func (st *SegmentTree2DDense) Query(row1, row2, col1, col2 int) E {
	if row1 >= row2 || col1 >= col2 {
		return st.unit
	}
	res := st.unit
	row1 += st.rowOffset
	row2 += st.rowOffset
	col1 += st.colOffset
	col2 += st.colOffset
	for row1 < row2 {
		if row1&1 == 1 {
			res = st.op(res, st._query(row1, col1, col2))
			row1++
		}
		if row2&1 == 1 {
			row2--
			res = st.op(res, st._query(row2, col1, col2))
		}
		row1 >>= 1
		row2 >>= 1
	}
	return res
}

func (st *SegmentTree2DDense) _query(r, c1, c2 int) E {
	res := st.unit
	for c1 < c2 {
		if c1&1 == 1 {
			res = st.op(res, st.tree[st._id(r, c1)])
			c1++
		}
		if c2&1 == 1 {
			c2--
			res = st.op(res, st.tree[st._id(r, c2)])
		}
		c1 >>= 1
		c2 >>= 1
	}
	return res
}

func (st *SegmentTree2DDense) _id(r, c int) int { return ((r * st.colOffset) << 1) + c }
