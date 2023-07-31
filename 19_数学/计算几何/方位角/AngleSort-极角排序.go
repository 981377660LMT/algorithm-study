package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	fmt.Println(AngleSort([][2]int{{0, 0}, {1, 0}, {0, 1}, {-1, 0}, {0, -1}}))
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
	order := AngleSort(points)
	for _, i := range order {
		fmt.Fprintln(out, points[i][0], points[i][1])
	}
}

// 极角排序，返回每个点在极角排序后的位置(order)
func AngleSort(points [][2]int) (order []int) {
	lower, origin, upper := []int{}, []int{}, []int{}
	O := [2]int{0, 0}
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
		a, b := points[lower[i]], points[lower[j]]
		return a[0]*b[1] > a[1]*b[0]
	})
	sort.Slice(upper, func(i, j int) bool {
		a, b := points[upper[i]], points[upper[j]]
		return a[0]*b[1] > a[1]*b[0]
	})
	order = lower
	order = append(order, origin...)
	order = append(order, upper...)
	return
}
