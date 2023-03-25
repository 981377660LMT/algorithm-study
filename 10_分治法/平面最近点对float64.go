/*
URL:
https://onlinejudge.u-aizu.ac.jp/courses/library/4/CGL/5/CGL_5_A
*/

// 指定区间的平面最近点对
// https://maspypy.github.io/library/geo/range_closest_pair_query.hpp

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const (
	G_EPS         = 1e-10
	INF   float64 = 1 << 60
)

var (
	gabs  = math.Abs
	gmin  = math.Min
	gmax  = math.Max
	gsqrt = math.Sqrt
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([]Point, n)
	for i := 0; i < n; i++ {
		var x, y float64
		fmt.Fscan(in, &x, &y)
		points[i] = NewPoint(x, y)
	}

	res := ClosestPair(points)
	fmt.Fprintf(out, "%.4f", res)
}

type Point struct {
	x, y float64
}

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func fEq(v, w float64) bool {
	return gabs(v-w) < G_EPS
}

func pEq(p, q Point) bool {
	dx, dy := p.x-q.x, p.y-q.y
	return fEq(dx, 0.0) && fEq(dy, 0.0)
}

func pLess(a, b Point) bool {
	if !fEq(a.x, b.x) {
		return a.x < b.x
	}
	return a.y < b.y
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_5_A
// 最近点対
func ClosestPair(P []Point) float64 {
	var _rec func(P []Point, l, r int) float64
	_rec = func(P []Point, l, r int) float64 {
		if r-l <= 1 {
			return INF
		}

		mid := (l + r) / 2
		x := P[mid].x
		d := gmin(_rec(P, l, mid), _rec(P, mid, r))

		// merge by order of y Pointinate.
		L, R := []Point{}, []Point{}
		for i := l; i < r; i++ {
			if i < mid {
				L = append(L, P[i])
			} else {
				R = append(R, P[i])
			}
		}
		cur, j := l, 0
		for i := 0; i < len(L); i++ {
			for j < len(R) && L[i].y > R[j].y {
				P[cur] = R[j]
				cur, j = cur+1, j+1
			}

			P[cur] = L[i]
			cur++
		}
		for ; j < len(R); j++ {
			P[cur] = R[j]
			cur++
		}

		nearLine := []Point{}
		for i := l; i < r; i++ {
			if gabs(P[i].x-x) >= d {
				continue
			}

			sz := len(nearLine)
			for j := sz - 1; j >= 0; j-- {
				dx := P[i].x - nearLine[j].x
				dy := P[i].y - nearLine[j].y
				if dy >= d {
					break
				}
				d = gmin(d, gsqrt(dx*dx+dy*dy))
			}
			nearLine = append(nearLine, P[i])
		}

		return d
	}

	sort.Slice(P, func(i, j int) bool {
		return pLess(P[i], P[j])
	})

	return _rec(P, 0, len(P))
}

// abs is integer version of math.Abs
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
