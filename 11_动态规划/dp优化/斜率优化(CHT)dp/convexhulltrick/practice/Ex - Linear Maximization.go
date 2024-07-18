// 二维CHT
// 给定q个查询,每个查询包含(ai,bi,xi,yi)
// 先添加直线ai*x+bi*y
// !再对于所有的点求出ai*x+bi*y的最大值
// q<=2e5

// TODO WA

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)
	queries := make([][4]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1], &queries[i][2], &queries[i][3])
	}
	res := LinearMaximization(queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

const INF int = 1e18

// (a,b,x,y)
//
//	(a,b) -> 添加直线ax+by
//	(x,y) -> 求出所有直线ax+by的最大值
func LinearMaximization(queries [][4]int) []int {
	res := make([]int, len(queries))
	chs := []*ConvexHull{}
	for i := 1; i <= len(queries); i++ {
		x, y, a, b := queries[i-1][0], queries[i-1][1], queries[i-1][2], queries[i-1][3]
		chs = append(chs, NewConvexHull([]*Point{{x, y}}))
		j := i
		for j&1 == 0 {
			j >>= 1
			p1 := chs[len(chs)-1].ps
			p2 := chs[len(chs)-2].ps
			chs = chs[:len(chs)-2]
			chs = append(chs, NewConvexHull(append(p1, p2...)))
		}
		cur := -INF
		for _, ch := range chs {
			cur = max(cur, ch.Get(a, b))
		}
		res[i-1] = cur
	}
	return res
}

type Point = [2]int

var _INV_PHI float64 = (math.Sqrt(5) - 1) / 2

func ccw(a, b, c *Point) int {
	cross := (b[0]-a[0])*(c[1]-b[1]) - (b[1]-a[1])*(c[0]-b[0])
	if cross > 0 { // a -> b -> c is counter clockwise
		return 1
	}
	if cross == 0 {
		return 0
	}
	return -1 // a -> b -> c is clockwise
}

type ConvexHull struct {
	ps           []*Point
	lower, upper []*Point
}

func NewConvexHull(ps []*Point) *ConvexHull {
	res := &ConvexHull{ps: ps}
	lower, upper := []*Point{}, []*Point{}
	sort.Slice(ps, func(i, j int) bool {
		if ps[i][0] == ps[j][0] {
			return ps[i][1] < ps[j][1]
		}
		return ps[i][0] < ps[j][0]
	})
	for i := 0; i < min(2, len(ps)); i++ {
		lower = append(lower, ps[i])
		upper = append(upper, ps[i])
	}
	for i := 2; i < len(ps); i++ {
		p := ps[i]
		lower = append(lower, p)
		upper = append(upper, p)
		for len(lower) >= 3 && ccw(lower[len(lower)-3], lower[len(lower)-2], lower[len(lower)-1]) <= 0 {
			lower = append(lower[:len(lower)-2], lower[len(lower)-1])
		}
		for len(upper) >= 3 && ccw(upper[len(upper)-3], upper[len(upper)-2], upper[len(upper)-1]) >= 0 {
			upper = append(upper[:len(upper)-2], upper[len(upper)-1])
		}
	}

	res.lower, res.upper = lower, upper
	return res
}

func (ch *ConvexHull) Get(a, b int) int {
	p := ch.lower
	if b > 0 {
		p = ch.upper
	}
	f := func(i int) int { return a*p[i][0] + b*p[i][1] }
	l := 0
	r := len(p) - 1
	r2 := int(math.Round(float64(r) * _INV_PHI))
	fr2 := f(r2)
	for abs(r-l) >= 6 {
		l2 := r + int(math.Round(float64(l-r)*_INV_PHI))
		fl2 := f(l2)
		if fl2 < fr2 {
			l, r = r, l2
		} else {
			r, r2, fr2 = r2, l2, fl2
		}
	}
	if l > r {
		l, r = r, l
	}
	res := f(l)
	for i := l + 1; i <= r; i++ {
		res = max(res, f(i))
	}
	return res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
