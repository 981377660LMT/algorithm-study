// 固定半径的圆最多能覆盖多少个点（圆边上也算覆盖） len(ps)>0 && r>0
package main

import (
	"math"
	"sort"
)

func main() {

}

const EPS = 1e-8

type point struct{ x, y int }

// 求一固定半径的圆最多能覆盖多少个点（圆边上也算覆盖） len(ps)>0 && r>0
// Angular Sweep 算法 O(n^2logn)
// https://www.geeksforgeeks.org/angular-sweep-maximum-points-can-enclosed-circle-given-radius/
// LC1453 https://leetcode-cn.com/problems/maximum-number-of-darts-inside-of-a-circular-dartboard/solution/python3-angular-sweepsuan-fa-by-lih/
func maxCoveredPoints(points []point, r int) int {
	type event struct {
		angle float64
		delta int
	}

	n := len(points)
	res := 1
	for i, a := range points {
		events := make([]event, 0, 2*n-2)
		for j, b := range points {
			if j == i {
				continue
			}
			ab := b.sub(a)
			if ab.len2() > 4*r*r {
				continue
			}
			at := math.Atan2(float64(ab.y), float64(ab.x))
			ac := math.Acos(ab.len() / float64(2*r))
			events = append(events, event{at - ac, 1}, event{at + ac, -1})
		}

		sort.Slice(events, func(i, j int) bool {
			a, b := events[i], events[j]
			return a.angle+EPS < b.angle || a.angle < b.angle+EPS && a.delta > b.delta
		})

		max_, count := 0, 1 // 1 指当前固定的点 a
		for _, e := range events {
			count += e.delta
			max_ = max(max_, count)
		}
		res = max(res, max_)
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func (a point) sub(b point) point { return point{a.x - b.x, a.y - b.y} }
func (a point) len2() int         { return a.x*a.x + a.y*a.y }
func (a point) len() float64      { return math.Sqrt(float64(a.x*a.x + a.y*a.y)) }
