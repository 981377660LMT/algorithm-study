/*
URL:
https://onlinejudge.u-aizu.ac.jp/courses/library/4/CGL/5/CGL_5_A
*/

package main

import (
	"bufio"
	"errors"
	"fmt"
	"math"
	"math/cmplx"
	"os"
	"sort"
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
	fmt.Fprintf(out, "%.10f", res)

}

// originated from:
// https://ei1333.github.io/luzhiled/snippets/geometry/template.html

type Point struct {
	x, y float64
}

type Line struct {
	a, b Point
}

type Segment struct {
	a, b Point
}

type Circle struct {
	p Point
	r float64
}

func NewPoint(x, y float64) Point {
	return Point{x: x, y: y}
}

func (p Point) Add(q Point) Point {
	x, y := p.x+q.x, p.y+q.y
	return NewPoint(x, y)
}

func (p Point) Minus(q Point) Point {
	x, y := p.x-q.x, p.y-q.y
	return NewPoint(x, y)
}

func (p Point) Mul(a float64) Point {
	return NewPoint(p.x*a, p.y*a)
}

func (p Point) Dot(q Point) float64 {
	return p.x*q.x + p.y*q.y
}

func (p Point) Cross(q Point) float64 {
	return p.x*q.y - p.y*q.x
}

func Dot(p, q Point) float64 {
	return p.Dot(q)
}

func Cross(p, q Point) float64 {
	return p.Cross(q)
}

func (p Point) Norm2() float64 {
	return p.x*p.x + p.y*p.y
}

func (p Point) Norm() float64 {
	return gsqrt(p.Norm2())
}

func Arg(p Point) float64 {
	return gatan2(p.y, p.x)
}

func Conj(p Point) Point {
	return NewPoint(p.x, -p.y)
}

const (
	G_EPS = 1e-10
	G_PI  = math.Pi

	G_ONLINE_FRONT      = -2
	G_CLOCKWISE         = -1
	G_ON_SEGMENT        = 0
	G_COUNTER_CLOCKWISE = 1
	G_ONLINE_BACK       = 2

	G_OUT = 0
	G_ON  = 1
	G_IN  = 2
)

var (
	gcos   = math.Cos
	gsin   = math.Sin
	gacos  = math.Acos
	gatan2 = math.Atan2

	gabs  = math.Abs
	gmin  = math.Min
	gmax  = math.Max
	gsqrt = math.Sqrt
)

func fEq(v, w float64) bool {
	return gabs(v-w) < G_EPS
}

func pEq(p, q Point) bool {
	dx, dy := p.x-q.x, p.y-q.y
	return fEq(dx, 0.0) && fEq(dy, 0.0)
}

func RotateTheta(t float64, p Point) Point {
	x := gcos(t)*p.x - gsin(t)*p.y
	y := gsin(t)*p.x + gcos(t)*p.y
	return NewPoint(x, y)
}

func RadianToDegree(r float64) float64 {
	return (r * 180.0) / G_PI
}

func DegreeToRadian(d float64) float64 {
	return (d * G_PI) / 180.0
}

// 余弦定理でa-b-cの角度のうち小さい方を返す
func Angle(a, b, c Point) float64 {
	v, w, z := a.Minus(b), c.Minus(b), a.Minus(c)
	cosTheta := (v.Norm2() + w.Norm2() - z.Norm2()) / (2.0 * v.Norm() * w.Norm())
	theta := gacos(cosTheta)
	return theta
}

func pLess(a, b Point) bool {
	if !fEq(a.x, b.x) {
		return a.x < b.x
	}
	return a.y < b.y
}

func NewLine(a, b Point) Line {
	na, nb := NewPoint(a.x, a.y), NewPoint(b.x, b.y)
	return Line{a: na, b: nb}
}

func NewSegment(a, b Point) Segment {
	na, nb := NewPoint(a.x, a.y), NewPoint(b.x, b.y)
	return Segment{a: na, b: nb}
}

func NewCircle(p Point, r float64) Circle {
	np := NewPoint(p.x, p.y)
	return Circle{p: np, r: r}
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_1_C
// 点の回転方向
func Ccw(a, b, c Point) int {
	d, e := b.Minus(a), c.Minus(a)

	cross := Cross(d, e)
	if !fEq(cross, 0.0) {
		if cross > 0.0 {
			return G_COUNTER_CLOCKWISE
		}
		return G_CLOCKWISE
	}

	dot := Dot(d, e)
	if !fEq(dot, 0.0) && dot < 0.0 {
		return G_ONLINE_BACK
	}
	if d.Norm2() < e.Norm2() {
		return G_ONLINE_FRONT
	}

	return G_ON_SEGMENT
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_2_A
// 平行判定
func IsParallel(a, b Line) bool {
	AB := a.b.Minus(a.a)
	CD := b.b.Minus(b.a)
	cross := Cross(AB, CD)
	return fEq(cross, 0.0)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_2_A
// 垂直判定
func IsOrthogonal(a, b Line) bool {
	BA := a.a.Minus(a.b)
	DC := b.a.Minus(b.b)
	dot := Dot(BA, DC)
	return fEq(dot, 0.0)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_1_A
// 射影
// 直線 l に p から垂線を引いた交点を求める
func ProjectionToLine(l Line, p Point) Point {
	AP := p.Minus(l.a)
	BA := l.a.Minus(l.b)
	t := Dot(AP, BA) / BA.Norm2()
	return l.a.Add(BA.Mul(t))
}
func ProjectionToSegment(l Segment, p Point) Point {
	line := NewLine(l.a, l.b)
	return ProjectionToLine(line, p)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_1_B
// 反射
// 直線 l を対称軸として点 p  と線対称にある点を求める
func Reflection(l Line, p Point) Point {
	plus := ProjectionToLine(l, p).Minus(p).Mul(2.0)
	return p.Add(plus)
}

// 交差判定
func IsIntersectLinePoint(l Line, p Point) bool {
	ccw := Ccw(l.a, l.b, p)
	return ccw != G_CLOCKWISE && ccw != G_COUNTER_CLOCKWISE
}
func IsIntersectLineLine(l, m Line) bool {
	AB := l.b.Minus(l.a)
	CD := m.b.Minus(m.a)
	cross := Cross(AB, CD)
	return !fEq(cross, 0.0)
}
func IsIntersectSegmentPoint(s Segment, p Point) bool {
	return Ccw(s.a, s.b, p) == G_ON_SEGMENT
}
func IsIntersectLineSegment(l Line, s Segment) bool {
	AB := l.b.Minus(l.a)
	AC := s.a.Minus(l.a)
	AD := s.b.Minus(l.a)
	return Cross(AB, AC)*Cross(AB, AD) < G_EPS
}
func IsIntersectCircleLine(c Circle, l Line) bool {
	return DistanceLinePoint(l, c.p) <= c.r+G_EPS
}
func IsIntersectCirclePoint(c Circle, p Point) bool {
	CP := p.Minus(c.p)
	return gabs(CP.Norm()-c.r) < G_EPS
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_2_B
func IsIntersectSegmentSegment(s, t Segment) bool {
	lb := Ccw(s.a, s.b, t.a)*Ccw(s.a, s.b, t.b) <= 0
	rb := Ccw(t.a, t.b, s.a)*Ccw(t.a, t.b, s.b) <= 0
	return lb && rb
}

func IsIntersectCircleSegment(c Circle, s Segment) int {
	m := ProjectionToSegment(s, c.p)
	CM := m.Minus(c.p)

	if CM.Norm2()-c.r*c.r > G_EPS {
		return 0
	}

	AC := c.p.Minus(s.a)
	BC := c.p.Minus(s.b)
	d1, d2 := AC.Norm(), BC.Norm()
	if d1 < c.r+G_EPS && d2 < c.r+G_EPS {
		return 0
	}
	if (d1 < c.r-G_EPS && d2 > c.r+G_EPS) || (d1 > c.r+G_EPS && d2 < c.r-G_EPS) {
		return 1
	}

	MA := s.a.Minus(m)
	MB := s.b.Minus(m)
	if Dot(MA, MB) < 0.0 {
		return 2
	}

	return 0
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_A&lang=jp
func IsIntersectCircleCircle(c1, c2 Circle) int {
	if c1.r < c2.r {
		c1, c2 = NewCircle(c2.p, c2.r), NewCircle(c1.p, c1.r)
	}

	d := c1.p.Minus(c2.p).Norm()
	if c1.r+c2.r < d {
		return 4
	}
	if fEq(c1.r+c2.r, d) {
		return 3
	}
	if c1.r-c2.r < d {
		return 2
	}
	if fEq(c1.r-c2.r, d) {
		return 1
	}

	return 0
}

func DistancePointPoint(a, b Point) float64 {
	AB := a.Minus(b)
	return AB.Norm()
}
func DistanceLinePoint(l Line, p Point) float64 {
	q := ProjectionToLine(l, p)
	QP := p.Minus(q)
	return QP.Norm()
}
func DistanceLineLine(l, m Line) float64 {
	if IsIntersectLineLine(l, m) {
		return 0.0
	}
	return DistanceLinePoint(l, m.a)
}
func DistanceSegmentPoint(s Segment, p Point) float64 {
	r := ProjectionToSegment(s, p)

	if IsIntersectSegmentPoint(s, r) {
		RP := p.Minus(r)
		return RP.Norm()
	}

	PA := s.a.Minus(p)
	PB := s.b.Minus(p)
	return gmin(PA.Norm(), PB.Norm())
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_2_D
func DistanceSegmentSegment(a, b Segment) float64 {
	if IsIntersectSegmentSegment(a, b) {
		return 0.0
	}

	d1 := DistanceSegmentPoint(a, b.a)
	d2 := DistanceSegmentPoint(a, b.b)
	d3 := DistanceSegmentPoint(b, a.a)
	d4 := DistanceSegmentPoint(b, a.b)

	return gmin(d1, gmin(d2, gmin(d3, d4)))
}
func DistanceLineSegment(l Line, s Segment) float64 {
	if IsIntersectLineSegment(l, s) {
		return 0.0
	}

	d1 := DistanceLinePoint(l, s.a)
	d2 := DistanceLinePoint(l, s.b)
	return gmin(d1, d2)
}

func CrossPointLineLine(l, m Line) Point {
	AB := l.b.Minus(l.a)
	CD := m.b.Minus(m.a)
	CB := l.b.Minus(m.a)
	A := Cross(AB, CD)
	B := Cross(AB, CB)

	if fEq(gabs(A), 0.0) && fEq(gabs(B), 0.0) {
		return NewPoint(m.a.x, m.a.y)
	}

	tmp := CD.Mul(B).Mul(1.0 / A)
	return m.a.Add(tmp)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_2_C
func CrossPointSegmentSegment(l, m Segment) Point {
	a := NewLine(l.a, l.b)
	b := NewLine(m.a, m.b)
	return CrossPointLineLine(a, b)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_D
func CrossPointsCircleLine(c Circle, l Line) [2]Point {
	M := ProjectionToLine(l, c.p)
	AB := l.b.Minus(l.a)
	e := AB.Mul(1.0 / AB.Norm())

	if fEq(DistanceLinePoint(l, c.p), c.r) {
		return [2]Point{NewPoint(M.x, M.y), NewPoint(M.x, M.y)}
	}

	CM := M.Minus(c.p)
	base := gsqrt(c.r*c.r - CM.Norm2())
	eb := e.Mul(base)
	return [2]Point{M.Minus(eb), M.Add(eb)}
}

func CrossPointsCircleSegment(c Circle, l Segment) [2]Point {
	aa := NewLine(l.a, l.b)
	if IsIntersectCircleSegment(c, l) == 2 {
		return CrossPointsCircleLine(c, aa)
	}

	ret := CrossPointsCircleLine(c, aa)
	CA := l.a.Minus(ret[0])
	CB := l.b.Minus(ret[0])
	if Dot(CA, CB) < 0.0 {
		ret[1] = ret[0]
	} else {
		ret[0] = ret[1]
	}

	return ret
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_E
func CrossPointsCircleCircle(c1, c2 Circle) [2]Point {
	d := c1.p.Minus(c2.p).Norm()

	a := gacos((c1.r*c1.r + d*d - c2.r*c2.r) / (2 * c1.r * d))
	t := gatan2(c2.p.y-c1.p.y, c2.p.x-c1.p.x)

	tmp1 := NewPoint(gcos(t+a)*c1.r, gsin(t+a)*c1.r)
	tmp2 := NewPoint(gcos(t-a)*c1.r, gsin(t-a)*c1.r)
	p1 := c1.p.Add(tmp1)
	p2 := c1.p.Add(tmp2)

	return [2]Point{p1, p2}
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_F
// 点 p を通る円 c の接線
func TangentCirclePoint(c1 Circle, p2 Point) [2]Point {
	c2 := NewCircle(p2, gsqrt(c1.p.Minus(p2).Norm2()-c1.r*c1.r))
	return CrossPointsCircleCircle(c1, c2)
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_G
// 円 c1, c2 の共通接線
func TangentCircleCircle(c1, c2 Circle) []Line {
	ret := []Line{}

	if c1.r < c2.r {
		c1, c2 = NewCircle(c2.p, c2.r), NewCircle(c1.p, c1.r)
	}

	g := c1.p.Minus(c2.p).Norm2()
	if fEq(g, 0.0) {
		return ret
	}

	u := c2.p.Minus(c1.p).Mul(1.0 / gsqrt(g))
	v := RotateTheta(G_PI*0.5, u)
	for _, s := range []float64{-1.0, 1.0} {
		h := (c1.r + s*c2.r) / gsqrt(g)
		if fEq(1.0-h*h, 0.0) {
			tmp1 := u.Mul(c1.r)
			tmp2 := u.Add(v).Mul(c1.r)
			l := NewLine(c1.p.Add(tmp1), c1.p.Add(tmp2))
			ret = append(ret, l)
		} else {
			uu := u.Mul(h)
			vv := v.Mul(gsqrt(1.0 - h*h))
			tmp1 := uu.Add(vv).Mul(c1.r)
			tmp2 := uu.Add(vv).Mul(c2.r * s)
			tmp3 := uu.Minus(vv).Mul(c1.r)
			tmp4 := uu.Minus(vv).Mul(c2.r * s)
			l1 := NewLine(c1.p.Add(tmp1), c2.p.Minus(tmp2))
			l2 := NewLine(c1.p.Add(tmp3), c2.p.Minus(tmp4))
			ret = append(ret, l1, l2)
		}
	}

	return ret
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_3_B
// 凸性判定
func IsConvex(P []Point) bool {
	n := len(P)
	for i := 0; i < n; i++ {
		a, b, c := P[(i+n-1)%n], P[i], P[(i+1)%n]
		if Ccw(a, b, c) == G_CLOCKWISE {
			return false
		}
	}
	return true
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_4_A
// 凸包
func ConvexHull(P []Point) []Point {
	n := len(P)
	k := 0

	if n <= 2 {
		return P
	}

	sort.Slice(P, func(i, j int) bool {
		return pLess(P[i], P[j])
	})

	ch := make([]Point, 2*n)
	for i := 0; i < n; i++ {
		for k >= 2 {
			CB := ch[k-1].Minus(ch[k-2])
			BA := P[i].Minus(ch[k-1])
			if Cross(CB, BA) > 0.0 {
				break
			}
			k--
		}

		ch[k] = P[i]
		k++
	}
	for i, t := n-2, k; i >= 0; i-- {
		for k >= t {
			CB := ch[k-1].Minus(ch[k-2])
			BA := P[i].Minus(ch[k-1])
			if Cross(CB, BA) > 0.0 {
				break
			}
			k--
		}

		ch[k] = P[i]
		k++
	}

	return ch[:k-1]
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_3_C
// 多角形と点の包含判定
func Contains(P []Point, q Point) int {
	n := len(P)
	isIn := false

	for i := 0; i < n; i++ {
		a := P[i].Minus(q)
		b := P[(i+1)%n].Minus(q)

		if a.y > b.y {
			a, b = NewPoint(b.x, b.y), NewPoint(a.x, a.y)
		}

		cross := Cross(a, b)
		dot := Dot(a, b)
		if a.y <= 0.0 && 0.0 < b.y && !fEq(cross, 0.0) && cross < 0.0 {
			isIn = !isIn
		}
		if fEq(cross, 0.0) && dot <= 0.0 {
			return G_ON
		}
	}

	if isIn {
		return G_IN
	}
	return G_OUT
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1033
// 線分の重複除去
// func MergeSegments(segs []*Segment) {

// }

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1033
// 線分アレンジメント
// 任意の2線分の交点を頂点としたグラフを構築する
// func SegmentArrangement(segs []*Segment, ps []*Point) [][]int {

// }

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_4_C
// 凸多角形の切断
// 直線 l.a-l.b で切断しその左側にできる凸多角形を返す
func ConvexCut(U []Point, l Line) []Point {
	n := len(U)
	ret := []Point{}

	for i := 0; i < n; i++ {
		a, b := U[i], U[(i+1)%n]
		now := NewPoint(a.x, a.y)
		nxt := NewPoint(b.x, b.y)

		if Ccw(l.a, l.b, now) != G_CLOCKWISE {
			ret = append(ret, now)
		}
		if Ccw(l.a, l.b, now)*Ccw(l.a, l.b, nxt) < 0 {
			cp := CrossPointLineLine(NewLine(now, nxt), l)
			ret = append(ret, cp)
		}
	}

	return ret
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_3_A
// 多角形の面積
func AreaPolygon(P []Point) float64 {
	n := len(P)
	A := 0.0
	for i := 0; i < n; i++ {
		A += Cross(P[i], P[(i+1)%n])
	}
	return A * 0.5
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_7_H
// 円と多角形の共通部分の面積
func AreaPolygonCircle(P []Point, c Circle) float64 {
	n := len(P)
	if n < 3 {
		return 0.0
	}

	var _cross_area func(c Circle, a, b Point) float64
	_cross_area = func(c Circle, a, b Point) float64 {
		va := c.p.Minus(a)
		vb := c.p.Minus(b)
		f := Cross(va, vb)
		ret := 0.0

		if fEq(f, 0.0) {
			return ret
		}
		if gmax(va.Norm(), vb.Norm()) < c.r+G_EPS {
			return f
		}
		if DistanceSegmentPoint(NewSegment(a, b), c.p) > c.r-G_EPS {
			cvb := complex(vb.x, vb.y)
			cvaConj := complex(va.x, -va.y)
			arg := cmplx.Phase(cvb * cvaConj)

			return c.r * c.r * arg
		}

		u := CrossPointsCircleSegment(c, NewSegment(a, b))
		tot := []Point{a, u[0], u[1], b}
		for i := 0; i+1 < len(tot); i++ {
			ret += _cross_area(c, tot[i], tot[i+1])
		}

		return ret
	}

	A := 0.0
	for i := 0; i < n; i++ {
		A += _cross_area(c, P[i], P[(i+1)%n])
	}
	return A
}

// originated from:
// https://onlinejudge.u-aizu.ac.jp/solutions/problem/CGL_7_I/review/4554366/beet/C++14
func AreaCircleCircle(ac1, ac2 Circle) float64 {
	c1, c2 := NewCircle(ac1.p, ac1.r), NewCircle(ac2.p, ac2.r)

	d := c1.p.Minus(c2.p).Norm()
	if c1.r+c2.r <= d+G_EPS {
		return 0.0
	}
	if d <= gabs(c1.r-c2.r) {
		r := gmin(c1.r, c2.r)
		return G_PI * r * r
	}

	P := CrossPointsCircleCircle(c1, c2)

	res := 0.0
	for i := 0; i < 2; i++ {
		th := Angle(c2.p, c1.p, P[0]) * 2.0
		res += (th - gsin(th)) * c1.r * c1.r / 2.0

		c1, c2 = c2, c1
	}

	return res
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_4_B
// 凸多角形の直径(最遠頂点対間距離)
func ConvexDiameter(P []Point) (maxDist float64, mi, mj int) {
	n := len(P)

	is, js := 0, 0
	for i := 1; i < n; i++ {
		if P[i].y > P[is].y {
			is = i
		}
		if P[i].y < P[js].y {
			js = i
		}
	}
	maxdis := P[is].Minus(P[js]).Norm2()

	i, maxi := is, is
	j, maxj := js, js
	for {
		if Cross(P[(i+1)%n].Minus(P[i]), P[(j+1)%n].Minus(P[j])) >= 0.0 {
			j = (j + 1) % n
		} else {
			i = (i + 1) % n
		}

		if P[i].Minus(P[j]).Norm2() > maxdis {
			maxdis = P[i].Minus(P[j]).Norm2()
			maxi, maxj = i, j
		}

		if !(i != is || j != js) {
			break
		}
	}

	return gsqrt(maxdis), maxi, maxj
}

// http://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=CGL_5_A
// 最近点対
func ClosestPair(P []Point) float64 {
	var _rec func(P []Point, l, r int) float64
	_rec = func(P []Point, l, r int) float64 {
		if r-l <= 1 {
			return 1e60
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

/*******************************************************************/

/********** common constants **********/

const (
	MOD = 1000000000 + 7
	// MOD          = 998244353
	ALPH_N  = 26
	INF_I64 = math.MaxInt64
	INF_B60 = 1 << 60
	INF_I32 = math.MaxInt32
	INF_B30 = 1 << 30
	NIL     = -1
	EPS     = 1e-10
)

// mod can calculate a right residual whether value is positive or negative.
func mod(val, m int) int {
	res := val % m
	if res < 0 {
		res += m
	}
	return res
}

// min returns the min integer among input set.
// This function needs at least 1 argument (no argument causes panic).
func min(integers ...int) int {
	m := integers[0]
	for i, integer := range integers {
		if i == 0 {
			continue
		}
		if m > integer {
			m = integer
		}
	}
	return m
}

// max returns the max integer among input set.
// This function needs at least 1 argument (no argument causes panic).
func max(integers ...int) int {
	m := integers[0]
	for i, integer := range integers {
		if i == 0 {
			continue
		}
		if m < integer {
			m = integer
		}
	}
	return m
}

// chmin accepts a pointer of integer and a target value.
// If target value is SMALLER than the first argument,
//	then the first argument will be updated by the second argument.
func chmin(updatedValue *int, target int) bool {
	if *updatedValue > target {
		*updatedValue = target
		return true
	}
	return false
}

// chmax accepts a pointer of integer and a target value.
// If target value is LARGER than the first argument,
//	then the first argument will be updated by the second argument.
func chmax(updatedValue *int, target int) bool {
	if *updatedValue < target {
		*updatedValue = target
		return true
	}
	return false
}

// sum returns multiple integers sum.
func sum(integers ...int) int {
	var s int
	s = 0

	for _, i := range integers {
		s += i
	}

	return s
}

// abs is integer version of math.Abs
func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

// pow is integer version of math.Pow
// pow calculate a power by Binary Power (二分累乗法(O(log e))).
func pow(a, e int) int {
	if a < 0 || e < 0 {
		panic(errors.New("[argument error]: PowInt does not accept negative integers"))
	}

	if e == 0 {
		return 1
	}

	if e%2 == 0 {
		halfE := e / 2
		half := pow(a, halfE)
		return half * half
	}

	return a * pow(a, e-1)
}
