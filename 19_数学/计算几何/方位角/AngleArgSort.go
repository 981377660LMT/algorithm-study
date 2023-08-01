package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	demo()
}

func demo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	points := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &points[i][0], &points[i][1])
	}
	res := AngleSort(points)
	for _, p := range res {
		fmt.Fprintln(out, p[0], p[1])
	}
}

type V = int

// 极角排序，返回值为点的下标
func AngleArgSort(points [][2]V) []int {
	lower, origin, upper := []int{}, []int{}, []int{}
	O := [2]V{0, 0}
	for i, p := range points {
		if p == O {
			origin = append(origin, i)
		} else if p[1] < 0 || (p[1] == 0 && p[0] > 0) {
			lower = append(lower, i)
		} else {
			upper = append(upper, i)
		}
	}

	sort.Slice(lower, func(i, j int) bool {
		oi, oj := lower[i], lower[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})
	sort.Slice(upper, func(i, j int) bool {
		oi, oj := upper[i], upper[j]
		pi, pj := points[oi], points[oj]
		return pi[0]*pj[1]-pi[1]*pj[0] > 0
	})

	res := lower
	res = append(res, origin...)
	res = append(res, upper...)
	return res
}

// 极角排序，返回值为排序后的点
func AngleSort(points [][2]V) [][2]V {
	order := AngleArgSort(points)
	res := make([][2]V, len(points))
	for i, o := range order {
		res[i] = points[o]
	}
	return res
}
