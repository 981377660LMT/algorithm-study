package main

import (
	"fmt"
	"sort"
)

func main() {
	points := [][2]int{{1, 2}, {3, 4}, {1, 2}, {3, 4}, {500000, 500000}, {-500000, -500000}}
	newXs, newYs, originX, originY, _, _ :=
		Disceretize2D(int32(len(points)), func(i int32) (int, int) { return points[i][0], points[i][1] })
	newPoints := make([][2]int, len(points))
	for i := range points {
		newPoints[i] = [2]int{int(newXs[i]), int(newYs[i])}
	}

	compress := func(x, y int) (int32, int32) {
		return int32(sort.SearchInts(originX, x)), int32(sort.SearchInts(originY, y))
	}

	fmt.Println(newPoints)
	fmt.Println(compress(1, 2))
	fmt.Println(compress(-30987654, 4))
}

// 3027. 人员站位的方案数 II
// https://leetcode.cn/problems/find-the-number-of-ways-to-place-people-ii/description/
func numberOfPairs(points [][]int) int {
	_, _, originX, originY, fx, fy :=
		Disceretize2D(int32(len(points)), func(i int32) (int, int) { return points[i][0], points[i][1] })

	row, col := int32(len(originX)), int32(len(originY))
	grid := make([][]int, row)
	for i := range grid {
		grid[i] = make([]int, col)
	}
	for i := range points {
		x, y := fx(points[i][0]), fy(points[i][1])
		points[i][0], points[i][1] = int(x), int(y)
		grid[x][y]++
	}

	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	S := NewPresumDense2D(e, op, inv)
	S.Build(row, col, func(r, c int32) int { return grid[r][c] })

	res := 0
	for _, p1 := range points {
		for _, p2 := range points {
			if S.Query(int32(p1[0]), int32(p2[0])+1, int32(p2[1]), int32(p1[1])+1) == 2 {
				res++
			}
		}
	}
	return res
}

// 二维离散化.
func Disceretize2D(n int32, f func(i int32) (x, y int)) (
	newXs, newYs []int32,
	originX, originY []int,
	fx, fy func(v int) int32,
) {
	xs, ys := make([]int, n), make([]int, n)
	for i := int32(0); i < n; i++ {
		xs[i], ys[i] = f(i)
	}
	newXs, originX = discretize1D(xs)
	newYs, originY = discretize1D(ys)
	fx = func(x int) int32 { return int32(sort.SearchInts(originX, x)) }
	fy = func(y int) int32 { return int32(sort.SearchInts(originY, y)) }
	return
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

type PresumDense2D[E any] struct {
	n, m   int32
	presum [][]E
	e      func() E
	op     func(a, b E) E
	inv    func(a E) E
}

func NewPresumDense2D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *PresumDense2D[E] {
	return &PresumDense2D[E]{e: e, op: op, inv: inv}
}

func (p *PresumDense2D[E]) Build(n int32, m int32, f func(r, c int32) E) {
	e, op := p.e, p.op
	presum := make([][]E, n+1)
	for i := int32(0); i < n+1; i++ {
		presum[i] = make([]E, m+1)
		for j := int32(0); j < m+1; j++ {
			presum[i][j] = e()
		}
	}
	for i := int32(1); i < n+1; i++ {
		for j := int32(1); j < m+1; j++ {
			presum[i][j] = op(f(i-1, j-1), presum[i][j-1])
		}
		for j := int32(1); j < m+1; j++ {
			presum[i][j] = op(presum[i][j], presum[i-1][j])
		}
	}
	p.n, p.m = n, m
	p.presum = presum
}

// [x1, x2) x [y1, y2)
func (p *PresumDense2D[E]) Query(x1, x2 int32, y1, y2 int32) E {
	x1, y1 = max32(x1, 0), max32(y1, 0)
	x2, y2 = min32(x2, p.n), min32(y2, p.m)
	if x1 >= x2 || y1 >= y2 {
		return p.e()
	}
	res := p.presum[x2][y2]
	res = p.op(res, p.inv(p.presum[x1][y2]))
	res = p.op(res, p.inv(p.presum[x2][y1]))
	res = p.op(res, p.presum[x1][y1])
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
