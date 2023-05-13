package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://yukicoder.me/problems/no/2003
// 青蛙从(0,0)出发 可以跳到1<=dx+dy<=k的格子上
// 不能跳到'#'上
// !求跳到(n-1,m-1)的方案数
// 2<=n*m<=1e6
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var ROW, COL, k int
	fmt.Fscan(in, &ROW, &COL, &k)
	grid := make([]string, ROW)
	for i := 0; i < ROW; i++ {
		fmt.Fscan(in, &grid[i])
	}

	bit2 := NewFenwickTree2DDense(ROW, COL)
	dp := make([][]int, ROW)
	for i := 0; i < ROW; i++ {
		dp[i] = make([]int, COL)
	}
	dp[0][0] = 1
	bit2.Add(0, 0, 1)

	// x+y=s
	getRange := func(s int) (int, int) {
		lo := max(0, s-COL+1)
		hi := min(ROW-1, s)
		return lo, hi + 1
	}

	for sum := 1; sum < ROW+COL-1; sum++ {
		a, b := getRange(sum - k - 1)
		for x := a; x < b; x++ {
			y := sum - k - 1 - x
			if grid[x][y] == '#' {
				continue
			}
			bit2.Add(x, y, -dp[x][y])
		}

		a, b = getRange(sum)
		for x := a; x < b; x++ {
			y := sum - x
			if grid[x][y] == '#' {
				continue
			}
			dp[x][y] = bit2.QueryPrefix(x+1, y+1)
			bit2.Add(x, y, dp[x][y])
		}
	}

	fmt.Fprintln(out, dp[ROW-1][COL-1])
}

const MOD int = 998244353

type E = int

func e() E { return 0 }
func op(a, b E) E {
	res := (a + b) % MOD
	if res < 0 {
		res += MOD
	}
	return res
}
func inv(a E) E { return -a }

type FenwickTree2DDense struct {
	H, W int
	data []E
	unit E
}

// 指定二维树状数组的高度和宽度.
func NewFenwickTree2DDense(h, w int) *FenwickTree2DDense {
	res := &FenwickTree2DDense{H: h, W: w, data: make([]E, h*w), unit: e()}
	for i := 0; i < h*w; i++ {
		res.data[i] = res.unit
	}
	return res
}

// 点 (x,y) 的值加上 val.
func (ft *FenwickTree2DDense) Add(x, y int, val E) {
	x++
	for x <= ft.H {
		ft.addX(x, y, val)
		x += x & -x
	}
}

// [lx,rx) * [ly,ry)
func (ft *FenwickTree2DDense) Query(lx, rx, ly, ry int) E {
	pos, neg := ft.unit, ft.unit
	for lx < rx {
		pos = op(pos, ft.sumX(rx, ly, ry))
		rx -= rx & -rx
	}
	for rx < lx {
		neg = op(neg, ft.sumX(lx, ly, ry))
		lx -= lx & -lx
	}
	return op(pos, inv(neg))
}

// [0,rx) * [0,ry)
func (ft *FenwickTree2DDense) QueryPrefix(rx, ry int) E {
	pos := ft.unit
	for rx > 0 {
		pos = op(pos, ft.sumXPrefix(rx, ry))
		rx -= rx & -rx
	}
	return pos
}

func (ft *FenwickTree2DDense) idx(x, y int) int {
	return ft.W*(x-1) + (y - 1)
}

func (ft *FenwickTree2DDense) addX(x, y int, val E) {
	y++
	for y <= ft.W {
		ft.data[ft.idx(x, y)] = op(ft.data[ft.idx(x, y)], val)
		y += y & -y
	}
}

func (ft *FenwickTree2DDense) sumX(x, ly, ry int) E {
	pos, neg := ft.unit, ft.unit
	for ly < ry {
		pos = op(pos, ft.data[ft.idx(x, ry)])
		ry -= ry & -ry
	}
	for ry < ly {
		neg = op(neg, ft.data[ft.idx(x, ly)])
		ly -= ly & -ly
	}
	return op(pos, inv(neg))
}

func (ft *FenwickTree2DDense) sumXPrefix(x, ry int) E {
	pos := ft.unit
	for ry > 0 {
		pos = op(pos, ft.data[ft.idx(x, ry)])
		ry -= ry & -ry
	}
	return pos
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
