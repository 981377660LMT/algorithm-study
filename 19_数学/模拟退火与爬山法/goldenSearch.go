package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

// Ex - Disk and Segments(三分套三分)
// https://atcoder.jp/contests/abc314/tasks/abc314_h
// 给定n条不相交的线段.
// 需要在平面上放置一个圆，使得所有线段与圆有交点或者在圆内.
// 求圆的半径的最小值.相对或者绝对误差小于1e-5.
// n<=100.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	type Point = [2]float64
	type Segment = [2]Point
	dot := func(a, b Point) float64 { return a[0]*b[0] + a[1]*b[1] }
	distP2P := func(a, b Point) float64 {
		dx, dy := a[0]-b[0], a[1]-b[1]
		return math.Sqrt(dx*dx + dy*dy)
	}
	// 点到线段的距离/点到直线的距离
	distS2P := func(seg Segment, p Point) float64 {
		p1, p2 := seg[0], seg[1]
		b1 := dot([2]float64{p2[0] - p1[0], p2[1] - p1[1]}, [2]float64{p[0] - p1[0], p[1] - p1[1]}) >= 0
		b2 := dot([2]float64{p1[0] - p2[0], p1[1] - p2[1]}, [2]float64{p[0] - p2[0], p[1] - p2[1]}) >= 0
		if b1 && !b2 {
			return distP2P(p2, p)
		}
		if !b1 && b2 {
			return distP2P(p1, p)
		}
		a, b, c := p1[1]-p2[1], p2[0]-p1[0], p1[0]*p2[1]-p1[1]*p2[0]
		return abs64(a*p[0]+b*p[1]+c) / math.Sqrt(a*a+b*b)
	}

	var n int32
	fmt.Fscan(in, &n)
	A, B := make([]Point, n), make([]Point, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &A[i][0], &A[i][1], &B[i][0], &B[i][1])
	}
	seg := make([]Segment, n)
	for i := int32(0); i < n; i++ {
		seg[i] = Segment{A[i], B[i]}
	}

	maxDist := func(x, y float64) float64 {
		p := Point{x, y}
		res := 0.0
		for i := int32(0); i < n; i++ {
			res = max64(res, distS2P(seg[i], p))
		}
		return res
	}

	INF := 1e6
	f := func(y float64) float64 {
		fx, _ := GoldenSearch(func(x float64) float64 { return maxDist(x, y) }, -INF, INF, -1)
		return fx
	}

	res, _ := GoldenSearch(func(y float64) float64 { return f(y) }, -INF, INF, -1)
	fmt.Fprintln(out, res)
}

func abs64(a float64) float64 {
	if a < 0 {
		return -a
	}
	return a
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

// GoldenSearch : 二分探索の一種で、関数の最小値を求めるアルゴリズム.
// f が評価される回数：2 + iter
// 幅が 1/phi^{iter} 倍になる. iter = 44: 1e-9.
func GoldenSearch(fun func(x float64) float64, min, max float64, iter int32) (fx, x float64) {
	if iter < 0 {
		iter = 50
	}
	invPhi := (math.Sqrt(5) - 1.0) * 0.5
	invPhi2 := invPhi * invPhi
	x1, x4 := min, max
	x2 := x1 + invPhi2*(x4-x1)
	x3 := x1 + invPhi*(x4-x1)
	y2, y3 := fun(x2), fun(x3)
	for i := int32(0); i < iter; i++ {
		if y2 < y3 {
			x4, x3, y3 = x3, x2, y2
			x2 = x1 + invPhi2*(x4-x1)
			y2 = fun(x2)
		} else {
			x1, x2, y2 = x2, x3, y3
			x3 = x1 + invPhi*(x4-x1)
			y3 = fun(x3)
		}
	}

	if y2 < y3 {
		return y2, x2
	} else {
		return y3, x3
	}
}
