// Point: 点
// Line: 直线
// Segment: 线段
// Circle: 圆
// Polygon: 多边形

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	yosupoAngleSort()
}

// https://judge.yosupo.jp/problem/sort_points_by_argument
func yosupoAngleSort() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	points := make([][2]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	order := AngleArgSort2(points)
	newPoints := ReArrange(func(i int32) [2]int { return points[i] }, order)
	for _, p := range order {
		fmt.Fprintln(out, newPoints[p][0], newPoints[p][1])
	}
}

type Integer interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64
}

type Float interface {
	float32 | float64
}

type Number interface {
	Integer | Float
}

type Point[T Number] struct{ x, y T }

func NewPoint[T Number](x, y T) Point[T] { return Point[T]{x, y} }

func (p Point[T]) Add(q Point[T]) Point[T] { return Point[T]{p.x + q.x, p.y + q.y} }
func (p Point[T]) Sub(q Point[T]) Point[T] { return Point[T]{p.x - q.x, p.y - q.y} }

// 叉积.
func (p Point[T]) Det(q Point[T]) T { return p.x*q.y - p.y*q.x }

// 点积.
func (p Point[T]) Dot(q Point[T]) T { return p.x*q.x + p.y*q.y }

func (p Point[T]) Norm() float64  { return math.Sqrt(float64(p.x*p.x + p.y*p.y)) }
func (p Point[T]) Angle() float64 { return math.Atan2(float64(p.y), float64(p.x)) }

func Rotate(p Point[float64], theta float64) Point[float64] {
	c, s := math.Cos(theta), math.Sin(theta)
	return Point[float64]{c*p.x - s*p.y, s*p.x + c*p.y}
}

// CCW: counter clockwise.
// a->b->c 逆时针为1，顺时针为-1，共线为0.
func CCW[T Number](a, b, c Point[T]) int32 {
	x := (b.Sub(a)).Det(c.Sub(a))
	if x > 0 {
		return 1
	}
	if x < 0 {
		return -1
	}
	return 0
}

func Dist[T Number](a, b Point[T]) float64 {
	c := a.Sub(b)
	return math.Sqrt(float64(c.x*c.x + c.y*c.y))
}

// 直线的一般式方程：ax+by+c=0
type Line[T Number] struct{ a, b, c T }

func NewLine[T Number](a, b, c T) Line[T] { return Line[T]{a, b, c} }
func NewLineFromPoints[T Number](a, b Point[T]) Line[T] {
	return Line[T]{a.y - b.y, b.x - a.x, a.x*b.y - a.y*b.x}
}
func NewLineFromCoords[T Number](x1, y1, x2, y2 T) Line[T] {
	return NewLineFromPoints(NewPoint(x1, y1), NewPoint(x2, y2))
}

func Normalize[T Integer](l *Line[T]) {
	g := Gcd(Gcd(Abs(l.a), Abs(l.b)), Abs(l.c))
	l.a /= g
	l.b /= g
	l.c /= g
	if l.b < 0 {
		l.a, l.b, l.c = -l.a, -l.b, -l.c
	}
	if l.b == 0 && l.a < 0 {
		l.a, l.b, l.c = -l.a, -l.b, -l.c
	}
}

// 点(x,y)代入直线方程ax+by+c求值.
func (l Line[T]) Eval(p Point[T]) T   { return l.a*p.x + l.b*p.y + l.c }
func (l Line[T]) EvalCoords(x, y T) T { return l.a*x + l.b*y + l.c }

// 直线平行/直线垂直.
func (l Line[T]) IsParallel(other Line[T]) bool   { return l.a*other.b-l.b*other.a == 0 }
func (l Line[T]) IsOrthogonal(other Line[T]) bool { return l.a*other.a+l.b*other.b == 0 }

type Segment[T Number] struct{ p1, p2 Point[T] }

func NewSegment[T Number](p1, p2 Point[T]) Segment[T] { return Segment[T]{p1, p2} }
func NewSegmentFromCoords[T Number](x1, y1, x2, y2 T) Segment[T] {
	return NewSegment(NewPoint(x1, y1), NewPoint(x2, y2))
}

func (s Segment[T]) Contains(p Point[T]) bool {
	det := (p.Sub(s.p1)).Det(s.p2.Sub(s.p1))
	if det != 0 {
		return false
	}
	return (p.Sub(s.p1)).Dot(s.p2.Sub(s.p1)) >= 0 && (p.Sub(s.p2)).Dot(s.p1.Sub(s.p2)) >= 0
}

func (s Segment[T]) ToLine() Line[T] { return NewLineFromPoints(s.p1, s.p2) }

type Circle[T Number] struct {
	center Point[T]
	radius T
}

func NewCircle[T Number](center Point[T], radius T) Circle[T] {
	return Circle[T]{center, radius}
}
func NewCircleFromCoords[T Number](x, y, r T) Circle[T] {
	return NewCircle(NewPoint(x, y), r)
}

// 点(x,y)是否在圆内/圆上.
func (c Circle[T]) Contains(p Point[T]) bool {
	dx := p.x - c.center.x
	dy := p.y - c.center.y
	return dx*dx+dy*dy <= c.radius*c.radius
}

type Polygon[T Number] struct {
	points []Point[T]
	a      T
}

func NewPolygon[T Number](points []Point[T]) Polygon[T] {
	p := Polygon[T]{points: points}
	p.build()
	return p
}

func NewPolygonFromPairs[T Number](pairs [][2]T) Polygon[T] {
	points := make([]Point[T], len(pairs))
	for i, p := range pairs {
		points[i] = NewPoint(p[0], p[1])
	}
	return NewPolygon(points)
}

func (p Polygon[T]) Size() int { return len(p.points) }

func (p Polygon[T]) Area() T  { return p.a / 2 }
func (p Polygon[T]) Area2() T { return p.a }
func (p Polygon[T]) IsConvex() bool {
	for j := 0; j < len(p.points); j++ {
		i := j - 1
		if i < 0 {
			i = len(p.points) - 1
		}
		k := j + 1
		if k == len(p.points) {
			k = 0
		}
		if (p.points[j].Sub(p.points[i])).Det(p.points[k].Sub(p.points[j])) < 0 {
			return false
		}
	}
	return true
}

func (p *Polygon[T]) build() {
	p.a = 0
	for i := 0; i < len(p.points); i++ {
		j := i + 1
		if j == len(p.points) {
			j = 0
		}
		p.a += p.points[i].Det(p.points[j])
	}
	if p.a < 0 {
		p.a = -p.a
		for i, j := 0, len(p.points)-1; i < j; i, j = i+1, j-1 {
			p.points[i], p.points[j] = p.points[j], p.points[i]
		}
	}
}

// 极角排序，返回值为点的下标.
func AngleArgSort[T Number](points []Point[T]) []int32 {
	var origin, lower, upper []int32
	O := Point[T]{0, 0}
	for i := int32(0); i < int32(len(points)); i++ {
		p := points[i]
		if p == O {
			origin = append(origin, i)
		} else if p.y < 0 || (p.y == 0 && p.x > 0) {
			lower = append(lower, i)
		} else {
			upper = append(upper, i)
		}
	}

	sort.Slice(lower, func(i, j int) bool {
		oi, oj := lower[i], lower[j]
		return points[oi].Det(points[oj]) > 0
	})
	sort.Slice(upper, func(i, j int) bool {
		oi, oj := upper[i], upper[j]
		return points[oi].Det(points[oj]) > 0
	})

	res := lower
	res = append(res, origin...)
	res = append(res, upper...)
	return res
}

func AngleArgSort2[T Number](points [][2]T) []int32 {
	ps := make([]Point[T], len(points))
	for i := 0; i < len(points); i++ {
		ps[i] = Point[T]{points[i][0], points[i][1]}
	}
	return AngleArgSort(ps)
}

func ReArrange[T any](f func(i int32) T, order []int32) []T {
	res := make([]T, len(order))
	for _, v := range order {
		res[v] = f(v)
	}
	return res
}

func Gcd[T Integer](a, b T) T {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func Abs[T Number](a T) T {
	if a < 0 {
		return -a
	}
	return a
}
