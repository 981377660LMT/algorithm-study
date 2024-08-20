package main

import (
	"fmt"
	"sort"
)

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumSparse2D(e, op, inv)
	points := [][2]int{{1, 2}, {3, 4}, {1, 2}, {3, 4}, {500000, 500000}, {-500000, -500000}}
	S.Build(int32(len(points)), func(i int32) (int, int, int) { return points[i][0], points[i][1], 1 })
	fmt.Println(S.Query(0, 3, 0, 3))
	fmt.Println(S.Query(1, 2, 1, 2))
}

// 3027. 人员站位的方案数 II
// https://leetcode.cn/problems/find-the-number-of-ways-to-place-people-ii/description/
func numberOfPairs(points [][]int) int {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumSparse2D(e, op, inv)
	S.Build(int32(len(points)), func(i int32) (int, int, int) { return points[i][0], points[i][1], 1 })
	res := 0
	for _, p1 := range points {
		for _, p2 := range points {
			if S.Query(p1[0], p2[0]+1, p2[1], p1[1]+1) == 2 {
				res++
			}
		}
	}
	return res
}

type PresumSparse2D[E any] struct {
	n, m    int32
	presum  [][]E
	originX []int
	originY []int
	e       func() E
	op      func(a, b E) E
	inv     func(a E) E
}

func NewPresumSparse2D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumSparse2D[E] {
	return &PresumSparse2D[E]{e: e, op: op, inv: inv}
}

// 二维平面上的前缀和.
//
//	f: 返回第i个元素的横纵坐标值和权值.
func (p *PresumSparse2D[E]) Build(n int32, f func(i int32) (x, y int, e E)) {
	xs, ys, es := make([]int, n), make([]int, n), make([]E, n)
	for i := int32(0); i < n; i++ {
		xs[i], ys[i], es[i] = f(i)
	}
	newXs, originX := discretize1D(xs)
	newYs, originY := discretize1D(ys)
	e, op := p.e, p.op
	row, col := int32(len(originX)), int32(len(originY))
	presum := make([][]E, row+1)
	for i := int32(0); i < row+1; i++ {
		presum[i] = make([]E, col+1)
		for j := int32(0); j < col+1; j++ {
			presum[i][j] = e()
		}
	}
	for i := int32(0); i < n; i++ {
		x, y, e := newXs[i], newYs[i], es[i]
		presum[x+1][y+1] = op(presum[x+1][y+1], e)
	}
	for i := int32(1); i < row+1; i++ {
		for j := int32(1); j < col+1; j++ {
			presum[i][j] = op(presum[i][j], presum[i][j-1])
		}
		for j := int32(1); j < col+1; j++ {
			presum[i][j] = op(presum[i][j], presum[i-1][j])
		}
	}

	p.n, p.m = row, col
	p.presum = presum
	p.originX, p.originY = originX, originY
}

// [x1, x2) x [y1, y2)
func (p *PresumSparse2D[E]) Query(x1, x2 int, y1, y2 int) E {
	if x1 >= x2 || y1 >= y2 {
		return p.e()
	}
	newX1, newX2 := p.compressX(x1), p.compressX(x2)
	newY1, newY2 := p.compressY(y1), p.compressY(y2)
	res := p.presum[newX2][newY2]
	res = p.op(res, p.inv(p.presum[newX1][newY2]))
	res = p.op(res, p.inv(p.presum[newX2][newY1]))
	res = p.op(res, p.presum[newX1][newY1])
	return res
}

func (p *PresumSparse2D[E]) compressX(x int) int32 {
	return int32(sort.SearchInts(p.originX, x))
}

func (p *PresumSparse2D[E]) compressY(y int) int32 {
	return int32(sort.SearchInts(p.originY, y))
}

func discretize1D(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
