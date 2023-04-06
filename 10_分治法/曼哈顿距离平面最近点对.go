package main

import (
	"sort"
)

// !曼哈顿距离平面最近点对(如果存在距离相等的,返回字典序最小的点对)
// n>=2
func beautifulPair(xs []int, ys []int) []int {
	points := make([][2]int, len(xs))
	mp := make(map[[2]int][]int)
	for i := range points {
		points[i] = [2]int{xs[i], ys[i]}
		mp[points[i]] = append(mp[points[i]], i)
	}

	samePair := [][]int{}
	for _, ps := range mp {
		if len(ps) > 1 {
			samePair = append(samePair, []int{ps[0], ps[1]})
		}
	}

	if len(samePair) > 0 {
		// 字典序最小的点对
		sort.Slice(samePair, func(i, j int) bool {
			if samePair[i][0] != samePair[j][0] {
				return samePair[i][0] < samePair[j][0]
			}
			return samePair[i][1] < samePair[j][1]
		})
		return samePair[0]
	}

	_, id1, id2 := ClosestPair(points)
	return []int{id1, id2}
}

const INF int = 1e18

// 计算距离的平方的公式.
func calDist2(x1, y1, x2, y2 int) int {
	res := (abs(x1-x2) + abs(y1-y2))
	return res * res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
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
