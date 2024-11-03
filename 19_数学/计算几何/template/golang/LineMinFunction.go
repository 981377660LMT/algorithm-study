// 一次函数最小值/一次函数最大值范围.
// from: https://maspypy.github.io/library/convex/line_min_function.hpp

package main

import (
	"fmt"
	"sort"
)

func main() {
	lines := [][2]float64{{1, 2}, {3, 1}, {4, 0}}
	res1 := LineMinFunction(lines)
	fmt.Println(res1) // [[-4e+18 0.6666666666666666 4 0] [0.6666666666666666 4e+18 1 2]]
	res2 := LineMaxFunction(lines)
	fmt.Println(res2) // [[-4e+18 0.5 -1 -2] [0.5 1 -3 -1] [1 4e+18 -4 -0]]
}

// 给定一组一次函数(斜率和截距)，求其最小值区间(起点、终点、斜率、截距).
func LineMinFunction(lines [][2]float64) (res [][4]float64) {
	points := make([][2]float64, len(lines))
	for i, line := range lines {
		points[i] = [2]float64{line[0], line[1]}
	}
	order := ConvexHull(points, Lower, false)
	points = reArrage(points, order)
	n := len(points)
	if n >= 2 && points[n-1][0] == points[n-2][0] {
		points = points[:n-1]
		n--
	}
	reverse(points)
	l := -INF
	for i := 0; i < n; i++ {
		r := INF
		a, b := points[i][0], points[i][1]
		if i+1 < n {
			c, d := points[i+1][0], points[i+1][1]
			if a == c {
				continue
			}
			r = (d - b) / (a - c)
			r = max64(r, l)
			r = min64(r, INF)
		}
		if l < r {
			res = append(res, [4]float64{l, r, a, b})
			l = r
		}
	}
	return
}

// 给定一组一次函数(斜率和截距)，求其最大值区间(起点、终点、斜率、截距).
func LineMaxFunction(lines [][2]float64) [][4]float64 {
	newLines := make([][2]float64, len(lines))
	for i, line := range lines {
		newLines[i] = [2]float64{-line[0], -line[1]}
	}
	res := LineMinFunction(newLines)
	for i := range res {
		res[i][2] = -res[i][2]
		res[i][3] = -res[i][3]
	}
	return res
}

type Mode uint8

const (
	Full Mode = iota
	Lower
	Upper
)

const INF float64 = 4e18

// (凸包/上凸包/下凸包).
// inclusive: 是否包含共线的点.
func ConvexHull(points [][2]float64, mode Mode, inclusive bool) []int32 {
	n := len(points)
	if n == 1 {
		return []int32{0}
	}

	compare := func(i, j int32) int8 {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[j][0], points[j][1]
		if x1 == x2 && y1 == y2 {
			if i < j {
				return -1
			}
			if i > j {
				return 1
			}
			return 0
		}
		if x1 < x2 || (x1 == x2 && y1 < y2) {
			return -1
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
		if inclusive {
			return []int32{0, 1}
		}
		return []int32{0}
	}

	order := make([]int32, n)
	for i := int32(0); i < int32(n); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return compare(order[i], order[j]) == -1 })

	check := func(i, j, k int32) bool {
		x1, y1 := points[i][0], points[i][1]
		x2, y2 := points[j][0], points[j][1]
		x3, y3 := points[k][0], points[k][1]
		dx1, dy1 := x2-x1, y2-y1
		dx2, dy2 := x3-x2, y3-y2
		det := dx1*dy2 - dx2*dy1
		if inclusive {
			return det >= 0
		}
		return det > 0
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

func argsortPoints(points [][2]float64) []int {
	n := len(points)
	order := make([]int, n)
	for i := 0; i < n; i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		p1, p2 := points[order[i]], points[order[j]]
		if p1 == p2 {
			return order[i] < order[j]
		}
		return p1[0] < p2[0] || (p1[0] == p2[0] && p1[1] < p2[1])
	})
	return order
}

func reArrage[T any](nums []T, order []int32) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
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

func min64(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max64(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}
