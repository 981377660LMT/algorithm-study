// G - No Cross Matching
// https://atcoder.jp/contests/abc373/tasks/abc373_g
// 二维平面，给定两组点p,q各n个，打乱q的顺序，使得n个线段pi→qi互不相交。给定一种q的顺序或告知不可能。
// nlognlogn算法.
// 距离之和最小的完美匹配.
// https://codeforces.com/contest/958/problem/E3
// https://codeforces.com/blog/entry/43463

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	CF958E3()
}

// CF958E3
// Guard Duty (hard)
// https://www.luogu.com.cn/problem/CF958E3
// 给定平面上n个红点和n个黑点，保证互不相同，无三点共线，
// 求一个红点到黑点的完美匹配，使得匹配之间的连线互不交叉。数据保证有解。
func CF958E3() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N int
	fmt.Fscan(in, &N)
	points := make([][2]int, 2*N)
	for i := 0; i < 2*N; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	idx := make([][2]int, 2*N)
	for i := 0; i < N; i++ {
		idx[i] = [2]int{0, i}
		idx[N+i] = [2]int{1, i}
	}
	order := argsortPoints(points)
	points = reArrage(points, order)
	idx = reArrage(idx, order)

	res := make([]int, N)

	match := func(x, y [2]int) {
		a, i := x[0], x[1]
		b, j := y[0], y[1]
		if a != 0 {
			a, b = b, a
			i, j = j, i
		}
		res[i] = j
	}

	type pair struct {
		points [][2]int
		idx    [][2]int
	}
	queue := []pair{{points, idx}}
	for len(queue) > 0 {
		points := queue[len(queue)-1].points
		idx := queue[len(queue)-1].idx
		queue = queue[:len(queue)-1]
		if len(idx) == 0 {
			continue
		}
		if len(idx) == 2 {
			match(idx[0], idx[1])
			continue
		}
		n := len(idx)
		done := make([]bool, n)
		ord := ConvexHull(points, Full, true)
		m := len(ord)
		upd := false
		for k := 0; k < m; k++ {
			i, j := ord[k], ord[(k+1)%m]
			if done[i] || done[j] {
				continue
			}
			if idx[i][0] != idx[j][0] {
				match(idx[i], idx[j])
				upd = true
				done[i], done[j] = true, true
			}
		}
		if upd {
			var I []int
			for i := 0; i < n; i++ {
				if !done[i] {
					I = append(I, i)
				}
			}
			points = reArrage(points, I)
			idx = reArrage(idx, I)
			queue = append(queue, pair{points, idx})
			continue
		}
		m = func() int {
			sm := 0
			for i := 0; i < n; i++ {
				sm += 2*idx[i][0] - 1
				if sm == 0 {
					return 1 + i
				}
			}
			panic("unreachable")
		}()

		queue = append(queue, pair{points[:m], idx[:m]})
		queue = append(queue, pair{points[m:n], idx[m:n]})
	}

	for _, x := range res {
		fmt.Fprintln(out, 1+x)
	}
}

type Mode uint8

const (
	Full Mode = iota
	Lower
	Upper
)

const INF int = 4e18

// (凸包/上凸包/下凸包).
// inclusive: 是否包含共线的点.
func ConvexHull(points [][2]int, mode Mode, inclusive bool) []int32 {
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

func argsortPoints(points [][2]int) []int {
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

func reArrage[T any](nums []T, order []int) []T {
	res := make([]T, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
