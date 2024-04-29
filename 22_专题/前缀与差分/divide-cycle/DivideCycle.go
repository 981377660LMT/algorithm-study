// 环区间分解(环分解，环遍历)
// DivideCycle/CycleDivide

package main

import "fmt"

func main() {
	DivideCycle(10, 8, 90, func(start, end int, times int) {
		fmt.Println(start, end, times)
	})
}

// 将环上的区间分解为形如`[start, end)`的区间，每个区间遍历`times`次.
func DivideCycle(n int, start int, end int, f func(start, end int, times int)) {
	if start >= end || n <= 0 {
		return
	}
	loop := (end - start) / n
	if loop > 0 {
		f(0, n, loop)
	}
	if (end-start)%n == 0 {
		return
	}
	start %= n
	end %= n
	if start < end {
		f(start, end, 1)
	} else {
		f(start, n, 1)
		if end > 0 {
			f(0, end, 1)
		}
	}
}
