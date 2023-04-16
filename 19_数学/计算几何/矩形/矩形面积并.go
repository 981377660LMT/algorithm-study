// 求矩形的面积并
//  rectangle[i] = [x1, y1, x2, y2]
//  0<=x1<x2<=10^9
//  0<=y1<y2<=10^9
//  1<=rectangle.length<=1e5
// https://leetcode.cn/problems/rectangle-area-ii/solution/ju-xing-mian-ji-ii-by-leetcode-solution-ulqz/

//
// Area of Union of Rectangles (Bentley)
//
// Description:
//   For a given set of rectangles, it gives the area of the union.
//   This problem is sometines called the Klee's measure problem [Klee'77].
//
// Algorithm:
//   Bentley's plane-sweep algorithm [Bentley'77].
//   We first apply the coordinate compression technique.
//   Then the y-structure, which is called measure tree, is simply implemented
//   by using segment tree data structure.
//
// Complexity:
//   O(n log n) time and O(n) space.
//
// Verify:
//   LightOJ 1120: Rectangle Union
//
// References:
//
//   V. Klee (1977):
//   Can the measure of \cup[a_i, b_i] be computed in less than O(n \log n) steps?
//   American Mathematical Monthly, vol.84, pp. 284--285.
//
//   J. L. Bentley (1977):
//   Algorithms for Klee's rectangle problems.
//   Unpublished notes, Computer Science Department, Carnegie Mellon University.
//

package main

import "sort"

func rectangleArea(rectangles [][]int) int {
	rs := make([]Rectangle, 0, len(rectangles))
	for _, r := range rectangles {
		rs = append(rs, Rectangle{r[0], r[1], r[2], r[3]})
	}
	return RectangleUnion(rs)
}

const MOD int = 1e9 + 7

type Rectangle struct {
	x1, y1, x2, y2 int // (x1,y1): 左下角坐标, (x2,y2): 右上角坐标.
}

// 矩形面积并.
func RectangleUnion(rectangles []Rectangle) int {
	ys := make([]int, 0, len(rectangles)*2)
	for _, r := range rectangles {
		ys = append(ys, r.y1, r.y2)
	}
	ys = sortedSet(ys)
	n := len(ys)
	C, A := make([]int, 8*n), make([]int, 8*n)
	var aux func(a, b, c, l, r, k int)
	aux = func(a, b, c, l, r, k int) {
		a, b = max(a, l), min(b, r)
		if a >= b {
			return
		}
		if a == l && b == r {
			C[k] += c
		} else {
			aux(a, b, c, l, (l+r)/2, 2*k+1)
			aux(a, b, c, (l+r)/2, r, 2*k+2)
		}
		if C[k] != 0 {
			A[k] = ys[r] - ys[l]
		} else {
			A[k] = A[2*k+1] + A[2*k+2]
		}
	}

	type event struct{ x, l, h, c int }
	es := make([]event, 0, len(rectangles)*2)
	for _, r := range rectangles {
		l, h := sort.SearchInts(ys, r.y1), sort.SearchInts(ys, r.y2)
		es = append(es, event{r.x1, l, h, +1}, event{r.x2, l, h, -1})
	}
	sort.Slice(es, func(i, j int) bool {
		if es[i].x != es[j].x {
			return es[i].x < es[j].x
		}
		return es[i].c > es[j].c
	})
	area, prev := 0, 0
	for _, e := range es {
		area += (e.x - prev) * A[0]
		area %= MOD
		prev = e.x
		aux(e.l, e.h, e.c, 0, n, 0)
	}
	return area % MOD
}

func sortedSet(xs []int) (sorted []int) {
	set := make(map[int]struct{}, len(xs))
	for _, v := range xs {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	return
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
