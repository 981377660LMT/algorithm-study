// FibonacciSearch 斐波那契搜索

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yuki2276()
}

func demo() {
	fmt.Println(FibonacciSearch(func(x int) int { return x * x }, 0, 100, true))
	fmt.Println(FibonacciSearch(func(x int) int { return x * x }, 0, 100, false))
}

// 2448. 使数组相等的最小开销
// https://leetcode.cn/problems/minimum-cost-to-make-array-equal/
func minCost(nums []int, cost []int) int64 {
	abs := func(x int) int {
		if x < 0 {
			return -x
		}
		return x
	}

	_, y := (FibonacciSearch(func(pos int) int {
		sum := 0
		for i, num := range nums {
			sum += abs(num-pos) * cost[i]
		}
		return sum
	}, 0, 1e6, true))

	return int64(y)
}

func yuki2276() {
	// 替换问号序列后得到AC子序列的最大方案数
	// https://yukicoder.me/problems/no/2276
	//
	// A, C, ? からなる長さ N の文字列 S が与えられます。
	// あなたは、S に対して次の操作を行います。
	// !S に含まれる ? を、それぞれ A または C のいずれかで書き換える。
	// !操作後、S から（連続するとは限らない）部分列として AC を取り出す方法の数が得点となります。
	// ここで、部分列を取り出す方法が異なるとは、取り出す位置が異なることをいいます。
	// 得点の最大値を求めてください。

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var s string
	fmt.Fscan(in, &s)

	sb := []byte(s)

	replacePos := make([]int, 0)
	for i := 0; i < n; i++ {
		if s[i] == '?' {
			replacePos = append(replacePos, i)
		}
	}

	// 将前n个?替换为A,后n个?替换为C,求AC的个数
	fun := func(n int) int {
		if !(0 <= n && n <= len(replacePos)) {
			panic("n out of range")
		}

		for i := 0; i < n; i++ {
			sb[replacePos[i]] = 'A'
		}
		for i := n; i < len(replacePos); i++ {
			sb[replacePos[i]] = 'C'
		}

		countA := 0
		countAC := 0
		for _, v := range sb {
			if v == 'A' {
				countA++
			}
			if v == 'C' {
				countAC += countA
			}
		}

		return countAC
	}

	_, y := FibonacciSearch(fun, 0, len(replacePos), false)
	fmt.Fprintln(out, y)
}

const INF int = 1e18

// 寻找[left,right]中的一个极值点,不要求单峰性质.
//
//	返回值: (极值点,极值)
func FibonacciSearch(f func(x int) int, left, right int, minimize bool) (int, int) {
	a, b, c, d := left, left+1, left+2, left+3
	step := 0
	for d < right {
		b = c
		c = d
		d = b + c - a
		step++
	}

	get := func(i int) int {
		if right < i {
			return INF
		}
		if minimize {
			return f(i)
		}
		return -f(i)
	}

	ya, yb, yc, yd := get(a), get(b), get(c), get(d)
	for i := 0; i < step; i++ {
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
