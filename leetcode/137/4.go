package main

import (
	"fmt"
	"math/bits"
)

// board = [[1,2,3],[4,5,6],[7,8,9]]
// [[-6,63,97],[43,-46,-19],[51,-7,52]]
// 158
// [[-53,-86,-80],[-28,16,-42],[-88,38,-66]]
func main() {
	board := [][]int{{1, 2, 3}, {4, 5, 6}, {7, 8, 9}}
	board = [][]int{{-6, 63, 97}, {43, -46, -19}, {51, -7, 52}}
	board = [][]int{{-53, -86, -80}, {-28, 16, -42}, {-88, 38, -66}}
	fmt.Println(maximumValueSum(board))
}

const INF int = 1e18

func maximumValueSum(board [][]int) int64 {
	type E = struct {
		max      int
		row, col int32
	}
	e := func() E { return E{-INF, -1, -1} }
	op := func(a, b E) E {
		if a.max > b.max {
			return a
		}
		return b
	}
	newGrid := make([][]E, len(board))
	for i := range newGrid {
		newGrid[i] = make([]E, len(board[i]))
		for j := range newGrid[i] {
			newGrid[i][j] = E{board[i][j], int32(i), int32(j)}
		}
	}
	st := NewSparseTableOnSegTreeFrom[E](newGrid, e, op)
	res := -INF
	ROW, COL := int32(len(board)), int32(len(board[0]))
	for r := int32(0); r < ROW; r++ {
		for c := int32(0); c < COL; c++ {
			// 分成了四个部分
			max1 := st.Query(0, r, 0, c)
			max2 := st.Query(0, r, c+1, COL)
			max3 := st.Query(r+1, ROW, 0, c)
			max4 := st.Query(r+1, ROW, c+1, COL)
			max2Sum := -INF
			max2Sum = max(max2Sum, max1.max+max4.max)
			max2Sum = max(max2Sum, max2.max+max3.max)

			// 1和2，但是不同行
			{
				if max1.row != max2.row {
					max2Sum = max(max2Sum, max1.max+max2.max)
				} else {
					// max1和2的次大值
					row2 := max2.row
					if row2 != -1 {
						cand1 := st.Query(0, row2, c+1, COL)
						cand2 := st.Query(row2+1, r, c+1, COL)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max1.max+candMax)
					}
					// max2和1的次大值
					row1 := max1.row
					if row1 != -1 {
						cand1 := st.Query(0, row1, 0, c)
						cand2 := st.Query(row1+1, r, 0, c)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max2.max+candMax)
					}
				}
			}

			// 1和3，但是不同列
			{
				if max1.col != max3.col {
					max2Sum = max(max2Sum, max1.max+max3.max)
				} else {
					// max1和3的次大值
					col2 := max3.col
					if col2 != -1 {
						cand1 := st.Query(r+1, ROW, 0, col2)
						cand2 := st.Query(r+1, ROW, col2+1, c)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max1.max+candMax)
					}
					// max3和1的次大值
					col1 := max1.col
					if col1 != -1 {
						cand1 := st.Query(0, r, 0, col1)
						cand2 := st.Query(0, r, col1+1, c)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max3.max+candMax)
					}
				}
			}

			// 2和4，但是不同列
			{
				if max2.col != max4.col {
					max2Sum = max(max2Sum, max2.max+max4.max)
				} else {
					// max2和4的次大值
					col4 := max4.col
					if col4 != -1 {
						cand1 := st.Query(r+1, ROW, c+1, col4)
						cand2 := st.Query(r+1, ROW, col4+1, COL)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max2.max+candMax)
					}
					// max4和2的次大值
					col2 := max2.col
					if col2 != -1 {
						cand1 := st.Query(0, r, c+1, col2)
						cand2 := st.Query(0, r, col2+1, COL)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max4.max+candMax)
					}
				}
			}

			// 3和4，但是不同行
			{
				if max3.row != max4.row {
					max2Sum = max(max2Sum, max3.max+max4.max)
				} else {
					// max3和4的次大值
					row4 := max4.row
					if row4 != -1 {
						cand1 := st.Query(r+1, row4, c+1, COL)
						cand2 := st.Query(row4+1, ROW, c+1, COL)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max3.max+candMax)
					}
					// max4和3的次大值
					row3 := max3.row
					if row3 != -1 {
						cand1 := st.Query(r+1, row3, 0, c)
						cand2 := st.Query(row3+1, ROW, 0, c)
						candMax := max(cand1.max, cand2.max)
						max2Sum = max(max2Sum, max4.max+candMax)
					}
				}
			}

			res = max(res, board[r][c]+max2Sum)
		}
	}
	return int64(res)
}

func max2Sum(a, b, c, d int) int {
	s1 := a + b
	s2 := a + c
	s3 := a + d
	s4 := b + c
	s5 := b + d
	s6 := c + d
	return max(max(max(max(max(s1, s2), s3), s4), s5), s6)
}

// 更快的 SparseTable2DFast.
type SparseTableOnSegTree[E any] struct {
	row, col int32
	e        func() E
	op       func(E, E) E
	data     []*SparseTable[E]
}

func NewSparseTableOnSegTreeFrom[E any](grid [][]E, e func() E, op func(E, E) E) *SparseTableOnSegTree[E] {
	row := int32(len(grid))
	col := int32(0)
	if row > 0 {
		col = int32(len(grid[0]))
	}
	data := make([]*SparseTable[E], 2*row)
	for i := int32(0); i < row; i++ {
		data[row+i] = NewSparseTableFrom(grid[i], e, op)
	}
	for i := row - 1; i > 0; i-- {
		data[i] = NewSparseTable(
			col,
			func(j int32) E {
				x := data[2*i].Query(j, j+1)
				y := data[2*i+1].Query(j, j+1)
				return op(x, y)
			},
			e, op,
		)
	}
	return &SparseTableOnSegTree[E]{row: row, col: col, e: e, op: op, data: data}
}

func (st *SparseTableOnSegTree[E]) Query(rowStart, rowEnd, colStart, colEnd int32) E {
	if !(0 <= rowStart && rowStart <= rowEnd && rowEnd <= st.row) {
		return st.e()
	}
	if !(0 <= colStart && colStart <= colEnd && colEnd <= st.col) {
		return st.e()
	}
	res := st.e()
	rowStart += st.row
	rowEnd += st.row
	for rowStart < rowEnd {
		if rowStart&1 != 0 {
			res = st.op(res, st.data[rowStart].Query(colStart, colEnd))
			rowStart++
		}
		if rowEnd&1 != 0 {
			rowEnd--
			res = st.op(res, st.data[rowEnd].Query(colStart, colEnd))
		}
		rowStart >>= 1
		rowEnd >>= 1
	}
	return res
}

type SparseTable[E any] struct {
	st [][]E
	e  func() E
	op func(E, E) E
	n  int32
}

func NewSparseTable[E any](n int32, f func(int32) E, e func() E, op func(E, E) E) *SparseTable[E] {
	res := &SparseTable[E]{}

	b := int32(bits.Len32(uint32(n)))
	st := make([][]E, b)
	for i := range st {
		st[i] = make([]E, n)
	}
	for i := int32(0); i < n; i++ {
		st[0][i] = f(i)
	}
	for i := int32(1); i < b; i++ {
		for j := int32(0); j+(1<<i) <= n; j++ {
			st[i][j] = op(st[i-1][j], st[i-1][j+(1<<(i-1))])
		}
	}
	res.st = st
	res.e = e
	res.op = op
	res.n = n
	return res
}

func NewSparseTableFrom[E any](leaves []E, e func() E, op func(E, E) E) *SparseTable[E] {
	return NewSparseTable(int32(len(leaves)), func(i int32) E { return leaves[i] }, e, op)
}

// 查询区间 [start, end) 的贡献值.
func (st *SparseTable[E]) Query(start, end int32) E {
	if start >= end {
		return st.e()
	}
	b := int32(bits.Len32(uint32(end-start))) - 1
	return st.op(st.st[b][start], st.st[b][end-(1<<b)])
}

func max32(a, b int32) int32 {
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
