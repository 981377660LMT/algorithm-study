package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// P5656()
	// TheFootballSeason()
}

// P5656 【模板】二元一次不定方程 (exgcd)
// https://www.luogu.com.cn/problem/P5656
func P5656() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		n, x1, y1, x2, y2 := SolveLinearEquation(a, b, c, false)
		if n == -1 {
			// 无整数解,输出 -1。
			fmt.Fprintln(out, -1)
		} else if n == 0 {
			// 有整数解但无正整数解,输出所有整数解中 x 的最小正整数值， y 的最小正整数值。
			fmt.Fprintln(out, x1, y2)
		} else {
			// 有正整数解,输出所有正整数解的个数，以及 x 的最小正整数值， y 的最小正整数值， x 的最大正整数值， y 的最大正整数值。
			fmt.Fprintln(out, n, x1, y2, x2, y1)
		}
	}
}

// https://www.luogu.com.cn/problem/CF1244C
// (x,y,z)表示一个为非负整数三元组，满足
// x+y+z=w，a*x+b*y=c
// 无解输出 -1，否则输出任意一组正整数解.
// !让非负解 x+y 尽量小,最简单的做法就是 min(x1+y1, x2+y2).
// !注意会爆longlong,用python最好
func TheFootballSeason() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var w, c, a, b int
	fmt.Fscan(in, &w, &c, &a, &b)
	n, x1, y1, x2, y2 := SolveLinearEquation(a, b, c, true)
	if n == -1 {
		fmt.Fprintln(out, -1)
	} else {
		res1, res2 := x1, y1
		if x1+y1 > x2+y2 || res1 < 0 || res2 < 0 {
			res1, res2 = x2, y2
		}
		if res1+res2 > w || res1 < 0 || res2 < 0 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, res1, res2, w-res1-res2)
		}
	}
}

// a*x + b*y = c 的通解为
// x = (c/g)*x0 + (b/g)*k
// y = (c/g)*y0 - (a/g)*k
// 其中 g = gcd(a,b) 且需要满足 g|c
// x0 和 y0 是 ax+by=g 的一组特解（即 exgcd(a,b) 的返回值）
//
// 为方便讨论，这里要求输入的 a b c 必须为正整数
// 注意若输入超过 1e9 可能要用高精
// 返回：正整数解的个数（无解时为 -1，无正整数解时为 0）
//
//	x 取最小正整数时的解 x1 y1，此时 y1 是最大正整数解
//	y 取最小正整数时的解 x2 y2，此时 x2 是最大正整数解
func SolveLinearEquation(a, b, c int, allowZero bool) (n, x1, y1, x2, y2 int) {
	g, x0, y0 := exgcd(a, b)

	// 无解
	if c%g != 0 {
		n = -1
		return
	}

	a /= g
	b /= g
	c /= g
	x0 *= c
	y0 *= c

	x1 = x0 % b
	if allowZero {
		if x1 < 0 {
			x1 += b
		}
	} else {
		if x1 <= 0 {
			x1 += b
		}
	}

	k1 := (x1 - x0) / b
	y1 = y0 - k1*a

	y2 = y0 % a
	if allowZero {
		if y2 < 0 {
			y2 += a
		}
	} else {
		if y2 <= 0 {
			y2 += a
		}
	}

	k2 := (y0 - y2) / a
	x2 = x0 + k2*b

	// 无正整数解
	if y1 <= 0 {
		return
	}

	// k 越大 x 越大
	n = k2 - k1 + 1
	return
}

// 二元一次不定方程（线性丢番图方程中的一种）https://en.wikipedia.org/wiki/Diophantine_equation
// exgcd solve equation ax+by=gcd(a,b), gcd 可能为负数.
// 特解满足 |x|<=|b|, |y|<=|a|
// https://cp-algorithms.com/algebra/extended-euclid-algorithm.html
func exgcd(a, b int) (gcd, x, y int) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, y, x = exgcd(b, a%b)
	y -= a / b * x
	return
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
