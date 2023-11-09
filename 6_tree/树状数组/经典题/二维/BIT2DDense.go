// Add : 单点加
// Query : 区间和
// QueryPrefix : 前缀和

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	tree := NewBIT2DDense(5, 5)
	tree.Update(1, 1, 1)
	tree.Update(1, 2, 10)
	tree.Build(func(x, y int) int { return x + y })
	fmt.Println(tree)
	fmt.Println(tree.QueryPrefix(2, 2))
}

func taiyaki() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=2842
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var H, W, T, Q int
	fmt.Fscan(in, &H, &W, &T, &Q)

	time := make([][]int, H)
	for i := 0; i < H; i++ {
		time[i] = make([]int, W)
		for j := 0; j < W; j++ {
			time[i][j] = 1e18
		}
	}
	A := NewBIT2DDense(H, W)
	B := NewBIT2DDense(H, W)

	end := NewDeque(Q)
	for i := 0; i < Q; i++ {
		var t, c, x, y int
		fmt.Fscan(in, &t, &c, &x, &y)
		x--
		y--
		for end.Size() > 0 && end.At(0)[2] <= t {
			e := end.PopLeft()
			x, y := e[0], e[1]
			A.Update(x, y, 1)
			B.Update(x, y, -1)
		}
		if c == 0 {
			B.Update(x, y, 1)
			end.Append([3]int{x, y, t + T})
		} else if c == 1 {
			if A.QueryRange(x, x+1, y, y+1) > 0 {
				A.Update(x, y, -1)
			}
		} else if c == 2 {
			var x2, y2 int
			fmt.Fscan(in, &x2, &y2)
			fmt.Fprintln(out, A.QueryRange(x, x2, y, y2), B.QueryRange(x, x2, y, y2))
		}
	}
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

// BIT2DDense
type BIT2DDense struct {
	H, W int
	data []E
	unit E
}

// 指定二维树状数组的高度和宽度.
func NewBIT2DDense(h, w int) *BIT2DDense {
	res := &BIT2DDense{H: h, W: w, data: make([]E, h*w), unit: e()}
	for i := 0; i < h*w; i++ {
		res.data[i] = res.unit
	}
	return res
}

func (ft *BIT2DDense) Build(f func(x, y int) int) {
	H, W := ft.H, ft.W
	for x := 0; x < H; x++ {
		for y := 0; y < W; y++ {
			ft.data[W*x+y] = f(x, y)
		}
	}
	for x := 1; x <= H; x++ {
		for y := 1; y <= W; y++ {
			ny := y + (y & -y)
			if ny <= W {
				ft.data[ft.idx(x, ny)] = op(ft.data[ft.idx(x, ny)], ft.data[ft.idx(x, y)])
			}
		}
	}
	for x := 1; x <= H; x++ {
		for y := 1; y <= W; y++ {
			nx := x + (x & -x)
			if nx <= H {
				ft.data[ft.idx(nx, y)] = op(ft.data[ft.idx(nx, y)], ft.data[ft.idx(x, y)])
			}
		}
	}
}

// 点 (x,y) 的值加上 val.
func (ft *BIT2DDense) Update(x, y int, val E) {
	x++
	for x <= ft.H {
		ft.addX(x, y, val)
		x += x & -x
	}
}

// [lx,rx) * [ly,ry)
func (ft *BIT2DDense) QueryRange(lx, rx, ly, ry int) E {
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
func (ft *BIT2DDense) QueryPrefix(rx, ry int) E {
	pos := ft.unit
	for rx > 0 {
		pos = op(pos, ft.sumXPrefix(rx, ry))
		rx -= rx & -rx
	}
	return pos
}

func (ft *BIT2DDense) String() string {
	res := make([][]string, ft.H)
	for i := 0; i < ft.H; i++ {
		res[i] = make([]string, ft.W)
		for j := 0; j < ft.W; j++ {
			res[i][j] = fmt.Sprintf("%v", ft.QueryRange(i, i+1, j, j+1))
		}

	}
	return fmt.Sprintf("%v", res)
}

func (ft *BIT2DDense) idx(x, y int) int {
	return ft.W*(x-1) + (y - 1)
}

func (ft *BIT2DDense) addX(x, y int, val E) {
	y++
	for y <= ft.W {
		ft.data[ft.idx(x, y)] = op(ft.data[ft.idx(x, y)], val)
		y += y & -y
	}
}

func (ft *BIT2DDense) sumX(x, ly, ry int) E {
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

func (ft *BIT2DDense) sumXPrefix(x, ry int) E {
	pos := ft.unit
	for ry > 0 {
		pos = op(pos, ft.data[ft.idx(x, ry)])
		ry -= ry & -ry
	}
	return pos
}

type D = [3]int
type Deque struct{ l, r []D }

func NewDeque(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
