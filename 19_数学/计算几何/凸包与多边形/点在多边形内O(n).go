package main

import "math"

func main() {

}

const EPS = 1e-8

type State int

const (
	OUT State = 0
	ON  State = 1
	IN  State = 2
)

// 判断点 p 是否在多边形 ps 内部（不保证凸性）
// https://oi-wiki.org/geometry/2d/#ray-casting-algorithm
// https://leetcode.cn/contest/sf-tech/problems/uWWzsv/
// http://acm.hdu.edu.cn/showproblem.php?pid=1756
// 法一：射线法（光线投射算法 Ray casting algorithm）  奇内偶外
// 由于转角法更方便，这里省略射线法的代码
// 法二：转角法（绕数法、回转数法）
// 【配图】https://blog.csdn.net/Form_/article/details/77855163
// 代码参考《训练指南》
// 这里统计绕数 Winding Number
// 从 p 出发向右作射线，统计多边形穿过这条射线正反多少次
// 【输入 ps 不要求是逆时针还是顺时针】
func inAnyPolygon(ps []point, p point) State {
	sign := func(x float64) int {
		if x < -EPS {
			return -1
		}
		if x < EPS {
			return 0
		}
		return 1
	}
	ps = append(ps, ps[0]) // 额外补一个点，方便枚举所有边
	wn := 0
	for i := 1; i < len(ps); i++ {
		p1, p2 := ps[i-1], ps[i]
		if p.onSeg(line{p1, p2}) {
			return ON // 在边界上
		}
		// det: 正左负右
		k := sign(float64(p2.sub(p1).det(p.sub(p1)))) // 适配 int 和 float64
		d1 := sign(float64(p1.y - p.y))
		d2 := sign(float64(p2.y - p.y))
		if k > 0 && d1 <= 0 && d2 > 0 { // 逆时针穿过射线（p 需要在 p1-p2 左侧）
			wn++
		} else if k < 0 && d2 <= 0 && d1 > 0 { // 顺时针穿过射线（p 需要在 p1-p2 右侧）
			wn--
		}
	}
	if wn != 0 {
		return IN // 在内部
	}
	return OUT // 在外部
}

/* 二维向量（点）*/
type point struct{ x, y int }

func (a point) add(b point) point   { return point{a.x + b.x, a.y + b.y} }
func (a point) sub(b point) point   { return point{a.x - b.x, a.y - b.y} }
func (a point) dot(b point) int     { return a.x*b.x + a.y*b.y }
func (a point) det(b point) int     { return a.x*b.y - a.y*b.x }
func (a point) len2() int           { return a.x*a.x + a.y*a.y }
func (a point) dis2(b point) int    { return a.sub(b).len2() }
func (a point) len() float64        { return math.Sqrt(float64(a.x*a.x + a.y*a.y)) }
func (a point) dis(b point) float64 { return a.sub(b).len() }

/* 二维直线（线段）*/
type line struct{ p1, p2 point }

// 方向向量 directional vector
func (a line) vec() point { return a.p2.sub(a.p1) }

// 点 a 是否在直线 l 上
// 判断方法：a-p1 与 a-p2 共线
func (a point) onLine(l line) bool {
	p1, p2 := l.p1.sub(a), l.p2.sub(a)
	return p1.det(p2) == 0
}

// 点 a 是否在线段 l 上
// 判断方法：a-p1 与 a-p2 共线且方向相反
func (a point) onSeg(l line) bool {
	p1, p2 := l.p1.sub(a), l.p2.sub(a)
	return p1.det(p2) == 0 && p1.dot(p2) <= 0 // 含端点
	//return math.Abs(p1.det(p2)) < eps && p1.dot(p2) < eps
}
