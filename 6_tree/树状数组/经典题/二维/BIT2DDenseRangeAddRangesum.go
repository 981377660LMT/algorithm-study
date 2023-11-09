package main

import "fmt"

func main() {
	bit2d := NewBIT2DDenseRangeAddRangesum(10, 10)
	bit2d.AddRange(1, 2, 1, 2, 1)
	fmt.Println(bit2d)
}

type BIT2DDenseRangeAddRangesum struct {
	_row, _col                     int
	_tree1, _tree2, _tree3, _tree4 []int
}

func NewBIT2DDenseRangeAddRangesum(row, col int) *BIT2DDenseRangeAddRangesum {
	return &BIT2DDenseRangeAddRangesum{
		_row:   row,
		_col:   col,
		_tree1: make([]int, (row+1)*(col+1)),
		_tree2: make([]int, (row+1)*(col+1)),
		_tree3: make([]int, (row+1)*(col+1)),
		_tree4: make([]int, (row+1)*(col+1)),
	}
}

// [row1,row2) x [col1,col2) 的值加上 delta.
func (b *BIT2DDenseRangeAddRangesum) AddRange(row1, row2, col1, col2, delta int) {
	if row1 >= row2 || col1 >= col2 {
		return
	}
	b._add(row1, col1, delta)
	b._add(row2, col1, -delta)
	b._add(row1, col2, -delta)
	b._add(row2, col2, delta)
}

// [0,row) x [0,col) 的和.
func (b *BIT2DDenseRangeAddRangesum) QueryPrefix(row, col int) (res int) {
	if row > b._row {
		row = b._row
	}
	if col > b._col {
		col = b._col
	}
	for r := row; r > 0; r -= r & -r {
		for c := col; c > 0; c -= c & -c {
			id := b._id(r, c)
			res += row*col*b._tree1[id] -
				col*b._tree2[id] -
				row*b._tree3[id] +
				b._tree4[id]
		}
	}
	return res
}

// [row1,row2) x [col1,col2) 的和.
func (b *BIT2DDenseRangeAddRangesum) QueryRange(row1, row2, col1, col2 int) int {
	if row2 > b._row {
		row2 = b._row
	}
	if col2 > b._col {
		col2 = b._col
	}
	if row1 >= row2 || col1 >= col2 {
		return 0
	}
	return b.QueryPrefix(row2, col2) - b.QueryPrefix(row1, col2) - b.QueryPrefix(row2, col1) + b.QueryPrefix(row1, col1)
}

func (b *BIT2DDenseRangeAddRangesum) String() string {
	res := make([][]int, b._row)
	for r := 0; r < b._row; r++ {
		res[r] = make([]int, b._col)
		for c := 0; c < b._col; c++ {
			res[r][c] = b.QueryRange(r, r+1, c, c+1)
		}
	}
	return fmt.Sprintf("%v", res)
}

func (b *BIT2DDenseRangeAddRangesum) _add(row, col, delta int) {
	row++
	col++
	for r := row; r <= b._row; r += r & -r {
		for c := col; c <= b._col; c += c & -c {
			id := b._id(r, c)
			b._tree1[id] += delta
			b._tree2[id] += (row - 1) * delta
			b._tree3[id] += (col - 1) * delta
			b._tree4[id] += (row - 1) * (col - 1) * delta
		}
	}
}

func (b *BIT2DDenseRangeAddRangesum) _id(row, col int) int {
	return row*(b._col+1) + col
}

type NumMatrix struct {
	bit    *BIT2DDenseRangeAddRangesum
	matrix [][]int
}

func Constructor(matrix [][]int) NumMatrix {
	n := len(matrix)
	if n == 0 {
		return NumMatrix{}
	}
	m := len(matrix[0])
	bit := NewBIT2DDenseRangeAddRangesum(n, m)
	for i := 0; i < n; i++ {
		for j := 0; j < m; j++ {
			bit.AddRange(i, i+1, j, j+1, matrix[i][j])
		}
	}
	return NumMatrix{bit: bit, matrix: matrix}
}

func (this *NumMatrix) Update(row int, col int, val int) {
	this.bit.AddRange(row, row+1, col, col+1, val-this.matrix[row][col])
	this.matrix[row][col] = val
}

func (this *NumMatrix) SumRegion(row1 int, col1 int, row2 int, col2 int) int {
	return this.bit.QueryRange(row1, row2+1, col1, col2+1)
}

/**
 * Your NumMatrix object will be instantiated and called as such:
 * obj := Constructor(matrix);
 * obj.Update(row,col,val);
 * param_2 := obj.SumRegion(row1,col1,row2,col2);
 */
