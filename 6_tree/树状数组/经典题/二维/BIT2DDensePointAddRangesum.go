package main

import "fmt"

func main() {
	bit2d := NewBIT2DDensePointAddRangeSum(3, 3)
	fmt.Println(bit2d)
	bit2d.Add(0, 0, 1)
	fmt.Println(bit2d)
	bit2d.Add(2, 1, 2)
	fmt.Println(bit2d)
	fmt.Println(bit2d.QueryRange(0, 3, 0, 3))
}

type BIT2DDensePointAddRangeSum struct {
	Row, Col int
	data     []int
}

func NewBIT2DDensePointAddRangeSum(row, col int) *BIT2DDensePointAddRangeSum {
	res := &BIT2DDensePointAddRangeSum{Row: row, Col: col, data: make([]int, row*col)}
	return res
}

func (ft *BIT2DDensePointAddRangeSum) Build(f func(x, y int) int) {
	Row, Col := ft.Row, ft.Col
	for x := 0; x < Row; x++ {
		for y := 0; y < Col; y++ {
			ft.data[Col*x+y] = f(x, y)
		}
	}
	for x := 1; x <= Row; x++ {
		for y := 1; y <= Col; y++ {
			ny := y + (y & -y)
			if ny <= Col {
				ft.data[ft.idx(x, ny)] += ft.data[ft.idx(x, y)]
			}
		}
	}
	for x := 1; x <= Row; x++ {
		for y := 1; y <= Col; y++ {
			nx := x + (x & -x)
			if nx <= Row {
				ft.data[ft.idx(nx, y)] += ft.data[ft.idx(x, y)]
			}
		}
	}
}

// 点 (x,y) 的值加上 delta.
func (ft *BIT2DDensePointAddRangeSum) Add(x, y int, delta int) {
	x++
	for x <= ft.Row {
		ft.addX(x, y, delta)
		x += x & -x
	}
}

// [lx,rx) * [ly,ry)
func (ft *BIT2DDensePointAddRangeSum) QueryRange(lx, rx, ly, ry int) int {
	if rx > ft.Row {
		rx = ft.Row
	}
	if ry > ft.Col {
		ry = ft.Col
	}
	if lx >= rx || ly >= ry {
		return 0
	}
	pos, neg := 0, 0
	for lx < rx {
		pos += ft.sumX(rx, ly, ry)
		rx -= rx & -rx
	}
	for rx < lx {
		neg += ft.sumX(lx, ly, ry)
		lx -= lx & -lx
	}
	return pos - neg
}

// [0,rx) * [0,ry)
func (ft *BIT2DDensePointAddRangeSum) QueryPrefix(rx, ry int) int {
	if rx > ft.Row {
		rx = ft.Row
	}
	if ry > ft.Col {
		ry = ft.Col
	}
	pos := 0
	for rx > 0 {
		pos += ft.sumXPrefix(rx, ry)
		rx -= rx & -rx
	}
	return pos
}

func (ft *BIT2DDensePointAddRangeSum) String() string {
	res := make([][]string, ft.Row)
	for i := 0; i < ft.Row; i++ {
		res[i] = make([]string, ft.Col)
		for j := 0; j < ft.Col; j++ {
			res[i][j] = fmt.Sprintf("%v", ft.QueryRange(i, i+1, j, j+1))
		}
	}
	return fmt.Sprintf("%v", res)
}

func (ft *BIT2DDensePointAddRangeSum) idx(x, y int) int {
	return ft.Col*(x-1) + (y - 1)
}

func (ft *BIT2DDensePointAddRangeSum) addX(x, y int, val int) {
	y++
	for y <= ft.Col {
		ft.data[ft.idx(x, y)] += val
		y += y & -y
	}
}

func (ft *BIT2DDensePointAddRangeSum) sumX(x, ly, ry int) int {
	pos, neg := 0, 0
	for ly < ry {
		pos += ft.data[ft.idx(x, ry)]
		ry -= ry & -ry
	}
	for ry < ly {
		neg += ft.data[ft.idx(x, ly)]
		ly -= ly & -ly
	}
	return pos - neg
}

func (ft *BIT2DDensePointAddRangeSum) sumXPrefix(x, ry int) int {
	pos := 0
	for ry > 0 {
		pos += ft.data[ft.idx(x, ry)]
		ry -= ry & -ry
	}
	return pos
}
