package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

// 逆时针给出 n 个凸多边形的顶点坐标，求它们交的面积
func main() {
	// https://www.luogu.com.cn/problem/P4196
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	polys := make([][]point, n)
	for i := range polys {
		var m int
		fmt.Fscan(in, &m)
		polys[i] = make([]point, m)
		for j := range polys[i] {
			fmt.Fscan(in, &polys[i][j].x, &polys[i][j].y)
		}
	}

	res := intersectedArea(polys)
	fmt.Fprintf(out, "%.3f", res) // 保留三位小数
}

// 求凸多边形的相交面积(凸多边形交)
func intersectedArea(polys [][]point) float64 {
	ls := []line{}
	for i := 0; i < len(polys); i++ {
		for j := 1; j < len(polys[i]); j++ {
			ls = append(ls, line{polys[i][j-1], polys[i][j]})
		}
		ls = append(ls, line{polys[i][len(polys[i])-1], polys[i][0]})
	}

	lps := HalfPlanesIntersection(ls)
	area := 0.0
	p0 := lps[0].p
	for i := 2; i < len(lps); i++ {
		area += lps[i-1].p.sub(p0).det(lps[i].p.sub(p0)) // 三角形面积
	}
	return area / 2
}

const EPS = 1e-8

type point struct{ x, y float64 }
type line struct{ p1, p2 point }
type lp struct {
	l line
	p point // l 与下一条直线的交点
}

// 半平面交
// 直线的一侧的平面就叫做半平面.
// 左半平面:y>=kx+b,右半平面:y<kx+b
// !应用:求凸多边形的面积交,求线性规划的可行域
// O(nlogn)，时间开销主要在排序上
// 大致思路：首先极角排序，然后用一个队列维护半平面交的顶点，每添加一个半平面，就不断检查队首队尾是否在半平面外，是就剔除
// 注意要先剔除队尾再剔除队首
// 注：凸包的对偶问题很接近半平面交，所以二者算法很接近
// https://oi-wiki.org/geometry/half-plane/
// https://www.luogu.com.cn/blog/105254/dui-ban-ping-mian-jiao-suan-fa-zheng-que-xing-xie-shi-di-tan-suo
// 模板题 https://www.luogu.com.cn/problem/P4196 https://www.acwing.com/problem/content/2805/
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
func HalfPlanesIntersection(ls []line) []lp { // !规定左侧为半平面
	sort.Slice(ls, func(i, j int) bool { return ls[i].vec().polarAngle() < ls[j].vec().polarAngle() })
	q := []lp{{l: ls[0]}}
	for i := 1; i < len(ls); i++ {
		l := ls[i]
		for len(q) > 1 && !q[len(q)-2].p.onLeft(l) {
			q = q[:len(q)-1]
		}
		for len(q) > 1 && !q[0].p.onLeft(l) {
			q = q[1:]
		}
		if math.Abs(l.vec().det(q[len(q)-1].l.vec())) < EPS {
			// 由于极角排序过，此时两有向直线平行且同向，取更靠内侧的直线
			if l.p1.onLeft(q[len(q)-1].l) {
				q[len(q)-1].l = l
			}
		} else {
			q = append(q, lp{l: l})
		}
		if len(q) > 1 {
			q[len(q)-2].p = q[len(q)-2].l.intersection(q[len(q)-1].l)
		}
	}
	// 最后用队首检查下队尾，删除无用半平面
	for len(q) > 1 && !q[len(q)-2].p.onLeft(q[0].l) {
		q = q[:len(q)-1]
	}

	// if len(q) < 3 {
	// 	!半平面交不足三个点的特殊情况，根据题意来返回
	// 	如果需要避免这种情况，可以先加入一个无穷大矩形对应的四个半平面，再求半平面交
	// 	return nil
	// }

	// 补上首尾半平面的交点
	q[len(q)-1].p = q[len(q)-1].l.intersection(q[0].l)
	return q
}

func (p point) add(b point) point    { return point{p.x + b.x, p.y + b.y} }
func (p point) sub(b point) point    { return point{p.x - b.x, p.y - b.y} }
func (p point) det(b point) float64  { return p.x*b.y - p.y*b.x }
func (p point) mul(k float64) point  { return point{p.x * k, p.y * k} }
func (p point) polarAngle() float64  { return math.Atan2(p.y, p.x) }
func (p point) onLeft(l line) bool   { return l.vec().det(p.sub(l.p1)) > EPS }
func (p line) vec() point            { return p.p2.sub(p.p1) }
func (p line) point(t float64) point { return p.p1.add(p.vec().mul(t)) }
func (p line) intersection(b line) point {
	va, vb, u := p.vec(), b.vec(), p.p1.sub(b.p1)
	t := vb.det(u) / va.det(vb)
	return p.point(t)
}
