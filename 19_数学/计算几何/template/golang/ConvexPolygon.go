// 凸多边形.

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	CF87E()
}

// https://www.luogu.com.cn/problem/CF87E
func CF87E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	get := func() *ConvexPolygon {
		var n int
		fmt.Fscan(in, &n)
		points := make([][2]int, n)
		for i := 0; i < n; i++ {
			fmt.Fscan(in, &points[i][0], &points[i][1])
		}
		return NewConvexPolygon(points)
	}

	P1, P2, P3 := get(), get(), get()
	P1 = MinkowskiSum(P1, P2)
	P1 = MinkowskiSum(P1, P3)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var p [2]int
		fmt.Fscan(in, &p[0], &p[1])
		p[0] *= 3
		p[1] *= 3
		res := P1.Side(p)
		if res >= 0 {
			fmt.Fprintln(out, "YES")
		} else {
			fmt.Fprintln(out, "NO")
		}
	}
}

// 闵可夫斯基和.
func MinkowskiSum(polygon1, polygon2 *ConvexPolygon) *ConvexPolygon {
	points := [][2]int{}
	p := [2]int{0, 0}
	for i := 0; i < 2; i++ {
		polygon1, polygon2 = polygon2, polygon1
		ps := polygon1.points
		n := int32(len(ps))
		for i := int32(0); i < n; i++ {
			j := (i + 1) % n
			points = append(points, [2]int{ps[j][0] - ps[i][0], ps[j][1] - ps[i][1]})
		}
		minP := minPoint(ps)
		p[0] += minP[0]
		p[1] += minP[1]
	}
	indices := angleArgSort(points)
	n := int32(len(indices))
	points = reArrage(points, indices)
	newPoints := make([][2]int, 0, n)
	newPoints = append(newPoints, [2]int{0, 0})
	for i := int32(0); i < n-1; i++ {
		newPoints = append(newPoints, [2]int{newPoints[i][0] + points[i][0], newPoints[i][1] + points[i][1]})
	}
	minP := minPoint(newPoints)
	add := [2]int{p[0] - minP[0], p[1] - minP[1]}
	for i := range newPoints {
		newPoints[i][0] += add[0]
		newPoints[i][1] += add[1]
	}
	order := ConvexHull(newPoints, Full, false)
	newPoints = reArrage(newPoints, order)
	return NewConvexPolygon(newPoints)
}

func minPoint(points [][2]int) [2]int {
	res := points[0]
	for _, p := range points {
		if p[0] < res[0] || (p[0] == res[0] && p[1] < res[1]) {
			res = p
		}
	}
	return res
}

// 极角排序，返回值为点的下标
func angleArgSort(points [][2]int) []int32 {
	lower, origin, upper := []int32{}, []int32{}, []int32{}
	O := [2]int{0, 0}
	for i := int32(0); i < int32(len(points)); i++ {
		p := points[i]
		if p == O {
			origin = append(origin, i)
		} else if p[1] < 0 || (p[1] == 0 && p[0] > 0) {
			lower = append(lower, i)
		} else {
			upper = append(upper, i)
		}
	}

	sort.Slice(lower, func(i, j int) bool {
		oi, oj := lower[i], lower[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})
	sort.Slice(upper, func(i, j int) bool {
		oi, oj := upper[i], upper[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})

	res := lower
	res = append(res, origin...)
	res = append(res, upper...)
	return res
}

type ConvexPolygon struct {
	n      int32
	points [][2]int
}

func NewConvexPolygon(points [][2]int) *ConvexPolygon {
	n := int32(len(points))
	if n < 3 {
		panic("ConvexPolygon: n < 3")
	}
	points = append(points[:0:0], points...)
	res := &ConvexPolygon{n: n, points: points}
	for i := int32(0); i < n; i++ {
		j := res.nextIndex(i)
		k := res.nextIndex(j)
		dx1 := res.points[j][0] - res.points[i][0]
		dy1 := res.points[j][1] - res.points[i][1]
		dx2 := res.points[k][0] - res.points[i][0]
		dy2 := res.points[k][1] - res.points[i][1]
		if dx1*dy2 <= dx2*dy1 {
			panic("ConvexPolygon: not convex")
		}
	}
	return res

}

// 点在多边内/点在多边形边上/点在多边形外.
// 内部: 1, 边界: 0, 外部: -1.
func (cp *ConvexPolygon) Side(p [2]int) int8 {
	l, r := int32(1), cp.n-1
	ps := cp.points
	a := det(ps[l][0]-ps[0][0], ps[l][1]-ps[0][1], p[0]-ps[0][0], p[1]-ps[0][1])
	b := det(ps[r][0]-ps[0][0], ps[r][1]-ps[0][1], p[0]-ps[0][0], p[1]-ps[0][1])
	if a < 0 || b > 0 {
		return -1
	}
	for r-l >= 2 {
		m := (l + r) / 2
		c := det(ps[m][0]-ps[0][0], ps[m][1]-ps[0][1], p[0]-ps[0][0], p[1]-ps[0][1])
		if c < 0 {
			r = m
			b = c
		} else {
			l = m
			a = c
		}
	}
	c := det(ps[r][0]-ps[l][0], ps[r][1]-ps[l][1], p[0]-ps[l][0], p[1]-ps[l][1])
	x := min(min(a, -b), c)
	if x < 0 {
		return -1
	}
	if x > 0 {
		return 1
	}
	if p == ps[0] {
		return 0
	}
	if c != 0 && a == 0 && l != 1 {
		return 1
	}
	if c != 0 && b == 0 && r != cp.n-1 {
		return 1
	}
	return 0
}

func (c *ConvexPolygon) MinDot(p [2]int) (int, int32) {
	ps := c.points
	index := c.periodicMinComp(func(i, j int32) bool {
		d1 := ps[i][0]*p[0] + ps[i][1]*p[1]
		d2 := ps[j][0]*p[0] + ps[j][1]*p[1]
		return d1 < d2
	})
	d := ps[index][0]*p[0] + ps[index][1]*p[1]
	return d, index
}

func (c *ConvexPolygon) MaxDot(p [2]int) (int, int32) {
	ps := c.points
	index := c.periodicMinComp(func(i, j int32) bool {
		d1 := ps[i][0]*p[0] + ps[i][1]*p[1]
		d2 := ps[j][0]*p[0] + ps[j][1]*p[1]
		return d1 > d2
	})
	d := ps[index][0]*p[0] + ps[index][1]*p[1]
	return d, index
}

func (c *ConvexPolygon) VisibleRange(p [2]int) (int32, int32) {
	ps := c.points
	a := c.periodicMinComp(func(i, j int32) bool {
		x1, y1 := ps[i][0]-p[0], ps[i][1]-p[1]
		x2, y2 := ps[j][0]-p[0], ps[j][1]-p[1]
		return det(x1, y1, x2, y2) < 0
	})
	b := c.periodicMinComp(func(i, j int32) bool {
		x1, y1 := ps[i][0]-p[0], ps[i][1]-p[1]
		x2, y2 := ps[j][0]-p[0], ps[j][1]-p[1]
		return det(x1, y1, x2, y2) > 0
	})
	{
		prev := c.prevIndex(a)
		x1, y1 := p[0]-ps[a][0], p[1]-ps[a][1]
		x2, y2 := p[0]-ps[prev][0], p[1]-ps[prev][1]
		if det(x1, y1, x2, y2) == 0 {
			a = prev
		}
	}
	{
		next := c.nextIndex(b)
		x1, y1 := p[0]-ps[b][0], p[1]-ps[b][1]
		x2, y2 := p[0]-ps[next][0], p[1]-ps[next][1]
		if det(x1, y1, x2, y2) == 0 {
			b = next
		}
	}
	return a, b
}

// 线段是否与多边形`内部`相交.
func (c *ConvexPolygon) IsCross(a, b [2]int) bool {
	for i := 0; i < 2; i++ {
		a, b = b, a
		l, r := c.VisibleRange(a)
		if det(c.points[l][0]-a[0], c.points[l][1]-a[1], b[0]-a[0], b[1]-a[1]) >= 0 {
			return false
		}
		if det(c.points[r][0]-a[0], c.points[r][1]-a[1], b[0]-a[0], b[1]-a[1]) <= 0 {
			return false
		}
	}
	return true
}

func (c *ConvexPolygon) periodicMinComp(f func(i, j int32) bool) int32 {
	n := c.n
	l, m, r := int32(0), n, n+n
	for {
		if r-l == 2 {
			break
		}
		l1, r1 := (l+m)/2, (m+r+1)/2
		if f(l1%n, m%n) {
			r = m
			m = l1
		} else if f(r1%n, m%n) {
			l = m
			m = r1
		} else {
			l = l1
			r = r1
		}
	}
	return m % n
}

func (c *ConvexPolygon) nextIndex(i int32) int32 {
	if i+1 == c.n {
		return 0
	}
	return i + 1
}

func (c *ConvexPolygon) prevIndex(i int32) int32 {
	if i == 0 {
		return c.n - 1
	}
	return i - 1
}

type Mode uint8

const (
	Full Mode = iota
	Lower
	Upper
)

const INF int = 4e18

// (凸包/上凸包/下凸包).
func ConvexHull(points [][2]int, mode Mode, isPointsSorted bool) []int32 {
	n := len(points)
	if n == 1 {
		return []int32{0}
	}

	compare := func(i, j int32) int8 {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[j][0], points[j][1]
		if x1 < x2 || (x1 == x2 && y1 < y2) {
			return -1
		}
		if x1 == x2 && y1 == y2 {
			return 0
		}
		return 1
	}

	if n == 2 {
		res := compare(0, 1)
		if res == -1 {
			return []int32{0, 1}
		}
		if res == 1 {
			return []int32{1, 0}
		}
		return []int32{0}
	}

	order := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		order[i] = i
	}
	if !isPointsSorted {
		sort.Slice(order, func(i, j int) bool { return compare(order[i], order[j]) == -1 })
	}

	check := func(i, j, k int32) bool {
		x1, y1 := points[j][0]-points[i][0], points[j][1]-points[i][1]
		x2, y2 := points[k][0]-points[i][0], points[k][1]-points[i][1]
		return x1*y2 > x2*y1
	}

	calc := func() []int32 {
		var p []int32
		for _, k := range order {
			for len(p) > 1 {
				i, j := p[len(p)-2], p[len(p)-1]
				if check(i, j, k) {
					break
				}
				p = p[:len(p)-1]
			}
			p = append(p, k)
		}
		return p
	}

	var p []int32
	if mode == Full || mode == Lower {
		p = append(p, calc()...)
	}
	if mode == Full || mode == Upper {
		if len(p) > 0 {
			p = p[:len(p)-1]
		}
		reverse(order)
		p = append(p, calc()...)
	}
	if mode == Upper {
		reverse(p)
	}
	for len(p) >= 2 && points[p[0]] == points[p[len(p)-1]] {
		p = p[:len(p)-1]
	}
	return p
}

func reverse[T any](nums []T) {
	for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
}

func reArrage[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}

func det(x1, y1, x2, y2 int) int {
	return x1*y2 - x2*y1
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
