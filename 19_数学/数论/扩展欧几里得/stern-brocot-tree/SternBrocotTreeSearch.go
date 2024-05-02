package main

import (
	"bufio"
	"fmt"
	"os"
)

// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=1208
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	for {
		var a, b int
		fmt.Fscan(in, &a, &b)
		if a == 0 && b == 0 {
			break
		}
		num1, deno1, num2, deno2 := SternBrocotTreeSearch(b, func(x, y int) bool { return x*x < y*y*a })
		fmt.Fprintf(out, "%d/%d %d/%d\n", num2, deno2, num1, deno1)
	}
}

// SternBrocotTreeSearch 搜索.
// 返回分子和分母都不超过n的最简分数num/deno中,
// 满足条件predicate(num/deno)的最大值num1/deno1,以及不满足条件predicate(num/deno)的最小值num2/deno2。
// predicate(num/deno)是单调的。
// 时间复杂度为O(f(n)logn),其中f(n)是计算predicate(num/deno)的时间复杂度。
// !可以用在01分数规划问题，避免浮点数二分精度
func SternBrocotTreeSearch(n int, predicate func(num, deno int) bool) (num1, deno1, num2, deno2 int) {
	a, b, c, d := 0, 1, 1, 0
	for {
		num := a + c
		den := b + d
		if num > n || den > n {
			break
		}
		if predicate(num, den) {
			k := 2
			for {
				num = a + k*c
				if num > n {
					break
				}
				den = b + k*d
				if den > n {
					break
				}
				if !predicate(num, den) {
					break
				}
				k *= 2
			}
			k /= 2
			a += c * k
			b += d * k
		} else {
			k := 2
			for {
				num = a*k + c
				if num > n {
					break
				}
				den = b*k + d
				if den > n {
					break
				}
				if predicate(num, den) {
					break
				}
				k *= 2
			}
			k /= 2
			c += a * k
			d += b * k
		}
	}
	return a, b, c, d
}
