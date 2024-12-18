package main

import (
	"fmt"
)

func main() {
	a := "ABC"
	b := "ABD"
	d := MyersDiff(a, b)
	fmt.Printf("The minimum edit distance between \"%s\" and \"%s\" is %d\n", a, b, d)
}

// 允许插入、删除操作.
// 时间复杂度O((N+M)D), 空间复杂度O(N+M)
func MyersDiff(a, b string) int32 {
	aLen := int32(len(a))
	bLen := int32(len(b))
	max := aLen + bLen

	type point struct {
		x int32
		y int32
	}
	v := make([]point, 2*max+1)

	// Initialize
	v[max+1].x = 0
	v[max+1].y = 0

	for d := int32(0); d <= max; d++ {
		for k := -d; k <= d; k += 2 {
			idx := max + k
			if k == -d || (k != d && v[idx-1].x < v[idx+1].x) {
				v[idx].x = v[idx+1].x
			} else {
				v[idx].x = v[idx-1].x + 1
			}
			v[idx].y = v[idx].x - k

			// Extend the snake
			for v[idx].x < aLen && v[idx].y < bLen && a[v[idx].x] == b[v[idx].y] {
				v[idx].x++
				v[idx].y++
			}

			if v[idx].x >= aLen && v[idx].y >= bLen {
				// Found the shortest path
				return d
			}
		}
	}

	return -1
}
