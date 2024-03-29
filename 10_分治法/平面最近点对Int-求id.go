package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

func main() {
	// 平面最近点对
	// 返回最近点对距离的平方
	// https://algs4.cs.princeton.edu/code/edu/princeton/cs/algs4/ClosestPair.java.html
	// 模板题 https://www.luogu.com.cn/problem/P1429 https://codeforces.com/problemset/problem/429/D
	// bichromatic closest pair 有两种类型的点，只需要额外判断类型是否不同即可 https://www.acwing.com/problem/content/121/ http://poj.org/problem?id=3714
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		var x, y int
		fmt.Fscan(in, &x, &y)
		points[i] = [2]int{x, y}
	}
	_, id1, id2 := ClosestPair(points)
	fmt.Fprintf(out, "%.4f", math.Sqrt(float64(calDist2(points[id1][0], points[id1][1], points[id2][0], points[id2][1]))))
}

const INF int = 1e18

// 距离的平方的公式.
func calDist2(x1, y1, x2, y2 int) int {
	return (x1-x2)*(x1-x2) + (y1-y2)*(y1-y2)
}

// 平面最近点对.返回(最近点对距离的平方, 点1的id, 点2的id).
//  如果存在多个相等的点对，返回字典序最小的对.
//  !调用 closestPair 前需保证没有重复的点.
//  https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
func ClosestPair(points [][2]int) (dist2, pid1, pid2 int) {
	if len(points) <= 1 {
		return INF, 0, 0
	}
	pWithId := make([][3]int, len(points))
	for i, p := range points {
		pWithId[i] = [3]int{p[0], p[1], i}
	}
	sort.Slice(pWithId, func(i, j int) bool { return pWithId[i][0] < pWithId[j][0] })
	return _closestPair(pWithId)
}

func _closestPair(ps [][3]int) (dist2, pid1, pid2 int) {
	n := len(ps)
	if n <= 1 {
		return INF, -1, -1
	}
	m := n >> 1
	x := ps[m][0]
	d1, pid1, pid2 := _closestPair(ps[:m])
	d2, pid3, pid4 := _closestPair(ps[m:])
	if d1 > d2 {
		d1, pid1, pid2 = d2, pid3, pid4
	} else if d1 == d2 && (pid1 > pid3 || (pid1 == pid3 && pid2 > pid4)) { // 字典序最小的点对
		pid1, pid2 = pid3, pid4
	}
	copy(ps, merge(ps[:m], ps[m:]))
	checkPs := [][3]int{}
	for _, pi := range ps {
		if (pi[0]-x)*(pi[0]-x) > d1 {
			continue
		}
		for j := len(checkPs) - 1; j >= 0; j-- {
			pj := checkPs[j]
			dy := pi[1] - pj[1]
			if dy*dy >= d1 { // > ?
				break
			}
			dx := pi[0] - pj[0]
			cand := calDist2(0, 0, dx, dy)
			if cand < d1 {
				d1, pid1, pid2 = cand, pi[2], pj[2]
			} else if cand == d1 && (pid1 > pi[2] || (pid1 == pi[2] && pid2 > pj[2])) { // 字典序最小的点对
				pid1, pid2 = pi[2], pj[2]
			}
		}
		checkPs = append(checkPs, pi)
	}
	if pid1 > pid2 {
		pid1, pid2 = pid2, pid1
	}
	return d1, pid1, pid2
}

func merge(a, b [][3]int) [][3]int {
	i, n := 0, len(a)
	j, m := 0, len(b)
	res := make([][3]int, 0, n+m)
	for {
		if i == n {
			return append(res, b[j:]...)
		}
		if j == m {
			return append(res, a[i:]...)
		}
		if a[i][1] < b[j][1] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
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
