package main

import (
	"fmt"
	"time"
)

func main() {
	time1 := time.Now()
	count := 0
	EnumeratePytagoreanTriple(1e8, func(a, b, c int) { count++ }, false)
	fmt.Println(count)
	fmt.Println(time.Since(time1))
}

// EnumeratePytagoreanTriple 遍历勾股数(勾股定理).
//
//	cLimit 限制勾股数的最大值.
//	f 回调函数.
//	coprimeOnly 是否只枚举两两互质的勾股数对.
//	cLimit = 1e8：互素对有 1.59*1e7 个, 0.2s;
//	cLimit = 1e8：全部有 2.71*1e8 个, 0.8s.
func EnumeratePytagoreanTriple(cLimit int, f func(a, b, c int), coprimeOnly bool) {
	stack := [][3]int{}
	add := func(a, b, c int) {
		if c <= cLimit {
			stack = append(stack, [3]int{a, b, c})
		}
	}
	add(3, 4, 5)
	for len(stack) > 0 {
		cur := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		a, b, c := cur[0], cur[1], cur[2]
		add(a-2*b+2*c, 2*a-b+2*c, 2*a-2*b+3*c)
		add(a+2*b+2*c, 2*a+b+2*c, 2*a+2*b+3*c)
		add(-a+2*b+2*c, -2*a+b+2*c, -2*a+2*b+3*c)
		if coprimeOnly {
			f(a, b, c)
		} else {
			x, y, z := a, b, c
			for z <= cLimit {
				f(x, y, z)
				x += a
				y += b
				z += c
			}
		}
	}
}
