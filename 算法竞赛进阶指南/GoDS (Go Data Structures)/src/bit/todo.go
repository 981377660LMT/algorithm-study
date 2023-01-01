// TODO
// https://judge.yosupo.jp/problem/static_rectangle_add_rectangle_sum
package main

import (
	"bufio"
	"fmt"
	"os"
)

const MOD int = 998244353

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q, left, down, right, up, add int
	fmt.Fscan(in, &n, &q)
	bit := NewBIT2D(1e9, 1e9)

	// n times
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &left, &down, &right, &up, &add)
		bit.UpdateRange(up, left, down, right, add)
	}

	// q times
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &left, &down, &right, &up)
		fmt.Fprintln(out, bit.QueryRange(up, left, down, right))
	}
}

// !二维区间查询 区间修改
type BIT2D struct {
	row   int
	col   int
	tree1 map[int]map[int]int
	tree2 map[int]map[int]int
	tree3 map[int]map[int]int
	tree4 map[int]map[int]int
}

func NewBIT2D(row, col int) *BIT2D {
	return &BIT2D{
		row:   row,
		col:   col,
		tree1: make(map[int]map[int]int, 1<<4),
		tree2: make(map[int]map[int]int, 1<<4),
		tree3: make(map[int]map[int]int, 1<<4),
		tree4: make(map[int]map[int]int, 1<<4),
	}
}

func (b *BIT2D) UpdateRange(row1 int, col1 int, row2 int, col2 int, delta int) {
	b.update(row1, col1, delta)
	b.update(row2+1, col1, -delta)
	b.update(row1, col2+1, -delta)
	b.update(row2+1, col2+1, delta)
}

func (b *BIT2D) QueryRange(row1 int, col1 int, row2 int, col2 int) int {
	return (b.query(row2, col2) - b.query(row2, col1-1) - b.query(row1-1, col2) + b.query(row1-1, col1-1) + 2*MOD) % MOD
}

func (b *BIT2D) update(row int, col int, delta int) {
	row, col = row+1, col+1
	preRow, preCol := row, col

	for curRow := row; curRow <= b.row; curRow += curRow & -curRow {
		for curCol := col; curCol <= b.col; curCol += curCol & -curCol {
			setDeep(b.tree1, curRow, curCol, delta)
			setDeep(b.tree2, curRow, curCol, (preRow-1)*delta)
			setDeep(b.tree3, curRow, curCol, (preCol-1)*delta)
			setDeep(b.tree4, curRow, curCol, (preRow-1)*(preCol-1)*delta)
		}
	}
}

func (b *BIT2D) query(row, col int) (res int) {
	row, col = row+1, col+1
	if row > b.row {
		row = b.row
	}
	if col > b.col {
		col = b.col
	}

	preR, preC := row, col
	for curR := row; curR > 0; curR -= curR & -curR {
		for curC := col; curC > 0; curC -= curC & -curC {
			res += preR*preC*getDeep(b.tree1, curR, curC) - preC*getDeep(b.tree2, curR, curC) - preR*getDeep(b.tree3, curR, curC) + getDeep(b.tree4, curR, curC)
		}
	}

	return
}

func setDeep(mp map[int]map[int]int, key1, key2, delta int) {
	if _, ok := mp[key1]; !ok {
		mp[key1] = make(map[int]int, 1<<4)
	}
	mp[key1][key2] += delta
}

func getDeep(mp map[int]map[int]int, key1, key2 int) int {
	if _, ok := mp[key1]; !ok {
		return 0
	}
	return mp[key1][key2]
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
