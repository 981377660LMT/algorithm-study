package main

import (
	"sort"
)

// https://leetcode.cn/problems/minimum-time-to-visit-all-houses/
func minTotalTime(forward []int, backward []int, queries []int) int64 {
	n := len(forward)
	C, _ := NewWeightedCycle(
		n,
		func(i int) int {
			return forward[i]
		},
		func(i int) int {
			if i == n-1 {
				return backward[0]
			}
			return backward[i+1]
		},
	)
	res := 0
	pre := 0
	for _, cur := range queries {
		res += C.Dist(pre, cur)
		pre = cur
	}
	return int64(res)
}

type WeightedCycle struct {
	n         int
	weightCw  []int
	weightCcw []int
	prefixCw  []int
	prefixCcw []int
}

// NewWeightedCycle 创建一个新的带权重的环
//
//	n: 环上的节点数量
//	weightCw: 顺时针方向边的权重， weightCw(i)表示从i顺时针到(i+1)%n的边权重.
//	weightCcw: 逆时针方向边的权重，weightCcw(i)表示从(i+1)%n逆时针到i的边权重.
func NewWeightedCycle(n int, weightCw, weightCcw func(i int) int) (*WeightedCycle, error) {
	wc := &WeightedCycle{
		n:         n,
		weightCw:  make([]int, n),
		weightCcw: make([]int, n),
	}
	for i := 0; i < n; i++ {
		wc.weightCw[i] = weightCw(i)
		wc.weightCcw[i] = weightCcw(i)
	}

	wc.buildPrefixSums()

	return wc, nil
}

func (wc *WeightedCycle) buildPrefixSums() {
	wc.prefixCw = make([]int, wc.n+1)
	wc.prefixCcw = make([]int, wc.n+1)
	for i := range wc.n {
		wc.prefixCw[i+1] = wc.prefixCw[i] + wc.weightCw[i]
		wc.prefixCcw[i+1] = wc.prefixCcw[i] + wc.weightCcw[i]
	}
}

func (wc *WeightedCycle) Dist(u, v int) int {
	return min(wc.DistCcw(u, v), wc.DistCw(u, v))
}

// 返回逆时针从from到to的带权距离.
func (wc *WeightedCycle) DistCcw(from, to int) int {
	if from >= to {
		return wc.prefixCcw[from] - wc.prefixCcw[to]
	}
	return wc.prefixCcw[from] + (wc.prefixCcw[wc.n] - wc.prefixCcw[to])
}

// 返回顺时针从from到to的带权距离.
func (wc *WeightedCycle) DistCw(from, to int) int {
	if to >= from {
		return wc.prefixCw[to] - wc.prefixCw[from]
	}
	return (wc.prefixCw[wc.n] - wc.prefixCw[from]) + wc.prefixCw[to]
}

// 返回环上两点间最短路径的线段表示.
func (wc *WeightedCycle) Segment(u, v int) [][2]int {
	if wc.DistCcw(u, v) <= wc.DistCw(u, v) {
		return wc.SegmentCcw(u, v)
	}
	return wc.SegmentCw(u, v)
}

// 返回逆时针从from到to的路径段.
func (wc *WeightedCycle) SegmentCcw(from, to int) [][2]int {
	if from >= to {
		return [][2]int{{from, to}}
	}
	return [][2]int{{from, 0}, {wc.n - 1, to}}
}

// 返回顺时针从from到to的路径段.
func (wc *WeightedCycle) SegmentCw(from, to int) [][2]int {
	if to >= from {
		return [][2]int{{from, to}}
	}
	return [][2]int{{from, wc.n - 1}, {0, to}}
}

// 返回环上两点间的最短路径.
func (wc *WeightedCycle) Path(u, v int) []int {
	if wc.DistCcw(u, v) <= wc.DistCw(u, v) {
		return wc.PathCcw(u, v)
	}
	return wc.PathCw(u, v)
}

// 返回逆时针从from到to的路径经过的点.
func (wc *WeightedCycle) PathCcw(from, to int) []int {
	var path []int
	if from >= to {
		for i := from; i >= to; i-- {
			path = append(path, i)
		}
	} else {
		for i := from; i >= 0; i-- {
			path = append(path, i)
		}
		for i := wc.n - 1; i >= to; i-- {
			path = append(path, i)
		}
	}
	return path
}

// 返回顺时针从from到to的路径经过的点.
func (wc *WeightedCycle) PathCw(from, to int) []int {
	var path []int
	if to >= from {
		for i := from; i <= to; i++ {
			path = append(path, i)
		}
	} else {
		for i := from; i < wc.n; i++ {
			path = append(path, i)
		}
		for i := 0; i <= to; i++ {
			path = append(path, i)
		}
	}
	return path
}

// 判断x是否在from到to的逆时针路径上.
func (wc *WeightedCycle) OnPathCcw(from, to, x int) bool {
	if x < to {
		x += wc.n
	}
	if from < to {
		from += wc.n
	}
	return to <= x && x <= from
}

// 判断x是否在from到to的顺时针路径上.
func (wc *WeightedCycle) OnPathCw(from, to, x int) bool {
	if from > to {
		to += wc.n
	}
	if from > x {
		x += wc.n
	}
	return from <= x && x <= to
}

// 逆时针从from出发走特定距离到达的位置.
func (wc *WeightedCycle) JumpCcw(from, distance int) int {
	if distance == 0 {
		return from
	}

	totalWeight := wc.prefixCcw[wc.n]
	distance = distance % totalWeight
	if distance == 0 {
		return from
	}

	target := wc.prefixCcw[from] - distance
	if target < 0 {
		target += totalWeight
	}

	pos := bisectLeft(wc.prefixCcw, target)
	if pos < wc.n {
		return pos
	}
	return 0
}

// 顺时针从from出发走特定距离到达的位置.
func (wc *WeightedCycle) JumpCw(from, distance int) int {
	if distance == 0 {
		return from
	}

	totalWeight := wc.prefixCw[wc.n]
	distance = distance % totalWeight
	if distance == 0 {
		return from
	}

	target := wc.prefixCw[from] + distance
	if target >= totalWeight {
		target -= totalWeight
	}

	pos := bisectRight(wc.prefixCw, target)
	return pos - 1
}

func bisectRight(a []int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] > x })
}

func bisectLeft(a []int, x int) int {
	return sort.Search(len(a), func(i int) bool { return a[i] >= x })
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
