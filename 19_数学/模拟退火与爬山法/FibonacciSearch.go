// # FibonacciSearch 斐波那契搜索
package main

import "fmt"

func main() {
	fmt.Println(FibonacciSearch(func(x int) int { return x * x }, 0, 100, true))
	fmt.Println(FibonacciSearch(func(x int) int { return x * x }, 0, 100, false))
}

const INF int = 1e18

// 寻找[start,end)中的一个极值点,不要求单峰性质.
//  返回值: (极值点,极值)
func FibonacciSearch(f func(x int) int, start, end int, minimize bool) (int, int) {
	end--
	a, b, c, d := start, start+1, start+2, start+3
	n := 0
	for d < end {
		b = c
		c = d
		d = b + c - a
		n++
	}

	get := func(i int) int {
		if end < i {
			return INF
		}
		if minimize {
			return f(i)
		}
		return -f(i)
	}

	ya, yb, yc, yd := get(a), get(b), get(c), get(d)
	for i := 0; i < n; i++ {
		if yb < yc {
			d = c
			c = b
			b = a + d - c
			yd = yc
			yc = yb
			yb = get(b)
		} else {
			a = b
			b = c
			c = a + d - b
			ya = yb
			yb = yc
			yc = get(c)
		}
	}

	x := a
	y := ya
	if yb < y {
		x = b
		y = yb
	}
	if yc < y {
		x = c
		y = yc
	}
	if yd < y {
		x = d
		y = yd
	}

	if minimize {
		return x, y
	}
	return x, -y

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
