// floorRange/enumerateFloor

package main

import (
	"fmt"
)

func main() {
	EnumerateFloor(20, func(left, right int, div int) {
		fmt.Println(left, right, div)
	})

	EnumerateFloorInterval(2, 20, func(left, right, div int) {
		fmt.Println(left, right, div)
	})

	EnumerateFloor2D(10, 10, func(x1, x2, y1, y2 int, div1, div2 int) {
		fmt.Println(x1, x2, y1, y2, div1, div2)
	})
}

// 将 [1,n] 内的数分成O(2*sqrt(n))段, 每段内的 n//i 相同
func EnumerateFloor(n int, f func(left, right int, div int)) {
	for l, r := 1, 0; l <= n; l = r + 1 {
		h := n / l
		r = n / h
		f(l, r, h)
	}
}

// 将 [lower,upper] 内的数分成O(2*sqrt(upper))段, 每段内的 upper//i 相同
func EnumerateFloorInterval(lower, upper int, f func(left, right, div int)) {
	for l, r := lower, 0; l <= upper; l = r + 1 {
		h := upper / l
		if h == 0 {
			break
		}
		r = min(upper/h, upper)
		f(l, r, h)
	}
}

// 将 [1,n] x [1,m] 内的数分成O(2*sqrt(n)*2*sqrt(m))段, 每段内的 (n//i, m//i) 相同
func EnumerateFloor2D(n, m int, f func(x1, x2, y1, y2 int, div1, div2 int)) {
	for x1, x2 := 1, 0; x1 <= n; x1 = x2 + 1 {
		hn := n / x1
		x2 = n / hn
		for y1, y2 := 1, 0; y1 <= m; y1 = y2 + 1 {
			hm := m / y1
			y2 = m / hm
			f(x1, x2, y1, y2, hn, hm)
		}
	}
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
