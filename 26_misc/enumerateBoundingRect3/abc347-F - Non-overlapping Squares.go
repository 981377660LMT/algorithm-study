// F - Non-overlapping Squares
// https://atcoder.jp/contests/abc347/tasks/abc347_f
// 给定一个n×n的网格，求三个不重叠的m×m的网格覆盖，和的最大值。
// n<=1000,m<=n//2

package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"strconv"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, m := int32(io.NextInt()), int32(io.NextInt())
	grid := make([][]int, n)
	for i := int32(0); i < n; i++ {
		grid[i] = make([]int, n)
		for j := int32(0); j < n; j++ {
			grid[i][j] = io.NextInt()
		}
	}

	preSum := NewPreSum2DFrom(grid)
	squareSum := make([][]int, n-m+1)
	for r := int32(0); r < n-m+1; r++ {
		squareSum[r] = make([]int, n-m+1)
		for c := int32(0); c < n-m+1; c++ {
			squareSum[r][c] = preSum.QueryRange(r, c, r+m-1, c+m-1)
		}
	}

	st := NewSparseTableOnSegTreeFrom(squareSum, func() int { return 0 }, max)

	// 求矩形范围内最大m*m正方形的和.
	calc := func(b BoundingRect) int {
		top, bottom, left, right := b[0], b[1], b[2], b[3]
		bottom, right = bottom-m+1, right-m+1
		return st.Query(top, bottom+1, left, right+1)
	}

	checkSize := func(b BoundingRect) bool {
		top, bottom, left, right := b[0], b[1], b[2], b[3]
		return bottom-top+1 >= m && right-left+1 >= m
	}

	res := 0
	EnumerateBoundingRect3(n, n, func(b1, b2, b3 BoundingRect) {
		if checkSize(b1) && checkSize(b2) && checkSize(b3) {
			res = max(res, calc(b1)+calc(b2)+calc(b3))
		}
	})
	io.Println(res)
}

type PreSum2D struct {
	preSum [][]int
}

func NewPreSum2D(row, col int32, f func(int32, int32) int) *PreSum2D {
	preSum := make([][]int, row+1)
	for i := range preSum {
		preSum[i] = make([]int, col+1)
	}
	for r := int32(0); r < row; r++ {
		for c := int32(0); c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2D{preSum}
}

func NewPreSum2DFrom(mat [][]int) *PreSum2D {
	return NewPreSum2D(int32(len(mat)), int32(len(mat[0])), func(r, c int32) int { return mat[r][c] })
}

// 查询sum(A[r1:r2+1, c1:c2+1])的值.
// 0 <= r1 <= r2 < row, 0 <= c1 <= c2 < col.
func (ps *PreSum2D) QueryRange(row1, col1, row2, col2 int32) int {
	return ps.preSum[row2+1][col2+1] - ps.preSum[row2+1][col1] - ps.preSum[row1][col2+1] + ps.preSum[row1][col1]
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

type BoundingRect = [4]int32 // (top,bottom,left,right)

// 给定一个row*col的矩阵,分割成3个不重合的矩形,返回所有可能的分割方法.
func EnumerateBoundingRect3(row, col int32, consumer func(b1, b2, b3 BoundingRect)) {
	// 三横
	for r1 := int32(0); r1 < row-2; r1++ {
		for r2 := r1 + 1; r2 < row-1; r2++ {
			consumer([4]int32{0, r1, 0, col - 1}, [4]int32{r1 + 1, r2, 0, col - 1}, [4]int32{r2 + 1, row - 1, 0, col - 1})
		}
	}

	// 三竖
	for c1 := int32(0); c1 < col-2; c1++ {
		for c2 := c1 + 1; c2 < col-1; c2++ {
			consumer([4]int32{0, row - 1, 0, c1}, [4]int32{0, row - 1, c1 + 1, c2}, [4]int32{0, row - 1, c2 + 1, col - 1})
		}
	}

	// 先一横 然后再切一竖
	for r := int32(0); r < row-1; r++ {
		for c := int32(0); c < col-1; c++ {
			consumer([4]int32{0, r, 0, c}, [4]int32{0, r, c + 1, col - 1}, [4]int32{r + 1, row - 1, 0, col - 1})
		}
		for c := int32(0); c < col-1; c++ {
			consumer([4]int32{0, r, 0, col - 1}, [4]int32{r + 1, row - 1, c + 1, col - 1}, [4]int32{r + 1, row - 1, 0, c})
		}
	}

	// 先一竖 再切一横
	for c := int32(0); c < col-1; c++ {
		for r := int32(0); r < row-1; r++ {
			consumer([4]int32{0, r, 0, c}, [4]int32{r + 1, row - 1, 0, c}, [4]int32{0, row - 1, c + 1, col - 1})
		}
		for r := int32(0); r < row-1; r++ {
			consumer([4]int32{0, row - 1, 0, c}, [4]int32{0, r, c + 1, col - 1}, [4]int32{r + 1, row - 1, c + 1, col - 1})
		}
	}
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
