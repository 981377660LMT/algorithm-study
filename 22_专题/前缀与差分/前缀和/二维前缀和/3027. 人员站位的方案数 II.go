// 3027. 人员站位的方案数 II (二维离散化+前缀和)
// https://leetcode.cn/problems/find-the-number-of-ways-to-place-people-ii/description/

package main

import "sort"

func numberOfPairs(points [][]int) int {
	xs, xtoi := SortedSet(int32(len(points)), func(i int32) int { return int(points[i][0]) })
	ys, ytoi := SortedSet(int32(len(points)), func(i int32) int { return int(points[i][1]) })

	m, n := int32(len(xs)), int32(len(ys))
	grid := make([][]int32, m)
	for i := range grid {
		grid[i] = make([]int32, n)
	}
	allPoints := make([][2]int32, 0, len(points))
	for _, p := range points {
		x, y := xtoi[p[0]], ytoi[p[1]]
		grid[x][y]++
		allPoints = append(allPoints, [2]int32{x, y})
	}

	preSum := NewPreSum2DDenseFrom(grid)
	res := int32(0)
	for _, p1 := range allPoints {
		for _, p2 := range allPoints {
			if p1[0] > p2[0] || p2[1] > p1[1] {
				continue
			}
			if preSum.SumRegion(p1[0], p2[0], p2[1], p1[1]) == 2 {
				res++
			}
		}
	}
	return int(res)

}

type Int interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64
}

func SortedSet[T Int](n int32, getX func(i int32) T) (newX []T, xtoi map[T]int32) {
	allNums := make(map[T]struct{})
	for i := int32(0); i < n; i++ {
		allNums[getX(i)] = struct{}{}
	}
	newX = make([]T, 0, len(allNums))
	for x := range allNums {
		newX = append(newX, x)
	}
	sort.Slice(newX, func(i, j int) bool { return newX[i] < newX[j] })
	xtoi = make(map[T]int32, len(newX))
	for i, x := range newX {
		xtoi[x] = int32(i)
	}
	return
}

type PreSum2DDense[T Int] struct {
	preSum [][]T
}

func NewPreSum2DDense[T Int](row, col int32, f func(int32, int32) T) *PreSum2DDense[T] {
	preSum := make([][]T, row+1)
	for i := range preSum {
		preSum[i] = make([]T, col+1)
	}
	for r := int32(0); r < row; r++ {
		for c := int32(0); c < col; c++ {
			preSum[r+1][c+1] = f(r, c) + preSum[r][c+1] + preSum[r+1][c] - preSum[r][c]
		}
	}
	return &PreSum2DDense[T]{preSum}
}

func NewPreSum2DDenseFrom[T Int](matrix [][]T) *PreSum2DDense[T] {
	return NewPreSum2DDense(int32(len(matrix)), int32(len(matrix[0])), func(r, c int32) T { return matrix[r][c] })
}

// 查询sum(A[x1:x2+1, y1:y2+1])的值(包含边界).
// 0 <= x1 <= x2 < row, 0 <= y1 <= y2 < col.
func (ps *PreSum2DDense[T]) SumRegion(x1, x2, y1, y2 int32) T {
	if x1 > x2 || y1 > y2 {
		return 0
	}
	return ps.preSum[x2+1][y2+1] - ps.preSum[x2+1][y1] -
		ps.preSum[x1][y2+1] + ps.preSum[x1][y1]
}
