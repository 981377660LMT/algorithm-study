package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	yosupoAngleSort()
}

// https://judge.yosupo.jp/problem/sort_points_by_argument
func yosupoAngleSort() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	points := make([][2]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	order := AngleArgSort2(points)
	newPoints := ReArrange(func(i int32) [2]int { return points[i] }, order)
	for _, p := range order {
		fmt.Fprintln(out, newPoints[p][0], newPoints[p][1])
	}
}

type Num interface {
	int | uint | int8 | uint8 | int16 | uint16 | int32 | uint32 | int64 | uint64 | float32 | float64
}

type Point[T Num] struct {
	x, y T
}

// 叉积.
func (p Point[T]) Det(q Point[T]) T {
	return p.x*q.y - p.y*q.x
}

// 点积.
func (p Point[T]) Dot(q Point[T]) T {
	return p.x*q.x + p.y*q.y
}

// 极角排序，返回值为点的下标.
func AngleArgSort[T Num](points []Point[T]) []int32 {
	var origin, lower, upper []int32
	O := Point[T]{0, 0}
	for i := int32(0); i < int32(len(points)); i++ {
		p := points[i]
		if p == O {
			origin = append(origin, i)
		} else if p.y < 0 || (p.y == 0 && p.x > 0) {
			lower = append(lower, i)
		} else {
			upper = append(upper, i)
		}
	}

	sort.Slice(lower, func(i, j int) bool {
		oi, oj := lower[i], lower[j]
		return points[oi].Det(points[oj]) > 0
	})
	sort.Slice(upper, func(i, j int) bool {
		oi, oj := upper[i], upper[j]
		return points[oi].Det(points[oj]) > 0
	})

	res := lower
	res = append(res, origin...)
	res = append(res, upper...)
	return res
}

func AngleArgSort2[T Num](points [][2]T) []int32 {
	order := make([]Point[T], len(points))
	for i := 0; i < len(points); i++ {
		order[i] = Point[T]{points[i][0], points[i][1]}
	}
	return AngleArgSort(order)
}

func ReArrange[T any](f func(i int32) T, order []int32) []T {
	res := make([]T, len(order))
	for _, v := range order {
		res[v] = f(v)
	}
	return res
}
