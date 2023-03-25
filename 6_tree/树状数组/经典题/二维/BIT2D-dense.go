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
	tree := NewFenwickTree2DDense(5, 5)
	tree.Add(1, 1, 1)
	tree.Add(1, 2, 10)
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
	A := NewFenwickTree2DDense(H, W)
	B := NewFenwickTree2DDense(H, W)

	end := NewDeque(Q)
	for i := 0; i < Q; i++ {
		var t, c, x, y int
		fmt.Fscan(in, &t, &c, &x, &y)
		x--
		y--
		for end.Size() > 0 && end.At(0)[2] <= t {
			e := end.PopLeft()
			x, y := e[0], e[1]
			A.Add(x, y, 1)
			B.Add(x, y, -1)
		}
		if c == 0 {
			B.Add(x, y, 1)
			end.Append([3]int{x, y, t + T})
		} else if c == 1 {
			if A.Query(x, x+1, y, y+1) > 0 {
				A.Add(x, y, -1)
			}
		} else if c == 2 {
			var x2, y2 int
			fmt.Fscan(in, &x2, &y2)
			fmt.Fprintln(out, A.Query(x, x2, y, y2), B.Query(x, x2, y, y2))
		}
	}
}

type E = int

func e() E        { return 0 }
func op(a, b E) E { return a + b }
func inv(a E) E   { return -a }

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

//
//
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
