// [ARC125E] Snack
// https://atcoder.jp/contests/arc125/tasks/arc125_e
//
// 有 n 种零食 m 个小孩。
// 把 n 种零食依次喂给小孩，一共只有 a[i] 种零食。
// 每个小孩`一种`最多只能吃 b[i] 个，一共只能吃 c[i] 个。
// 问所有小孩吃到零食数量的最大值。
// n,m<=2e5
//
// !我们从源点向每个零食连流量为 a[i] 的边，每个零食向每个小孩连对应的流量为 b[i] 的边，所有小孩向汇点连流量为 c[i] 的边。
// !边权的特点在于`从所有(左侧)零食到(右侧)小孩i的容量都是 b[i]`(有很多冗余边).
// 最大流转最小割，
// 枚举源点到每个零食的边是否断掉，记有n-x个零食与S相连(割x条零食到S的边)。
// 那么每个人都有两种决策：把到n-x个零食的边全部断掉，或是断掉到汇点的边，代价为 `min((n-x)*b[i], c[i])`。
// !最小割答案为：选取x个A[i]割掉，加上 sum(min((n-x)*b[i], c[i]))。
// 解法：
// 零食从小到大排序，小孩按照B[i]/C[i]排序.
// 枚举割掉的零食数量x，三分小孩的决策点，O(nlogn)。

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)

	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &B[i])
	}
	C := make([]int, m)
	for i := 0; i < m; i++ {
		fmt.Fscan(in, &C[i])
	}

	child := make([][2]int, m)
	for i := 0; i < m; i++ {
		child[i] = [2]int{B[i], C[i]}
	}

	sort.Ints(A)
	sort.Slice(child, func(i, j int) bool {
		return child[i][0]*child[j][1] < child[j][0]*child[i][1]
	})
	preSumA := make([]int, n+1)
	for i := 1; i <= n; i++ {
		preSumA[i] = preSumA[i-1] + A[i-1]
	}
	preSumB := make([]int, m+1)
	for i := 1; i <= m; i++ {
		preSumB[i] = preSumB[i-1] + child[i-1][0]
	}
	preSumC := make([]int, m+1)
	for i := 1; i <= m; i++ {
		preSumC[i] = preSumC[i-1] + child[i-1][1]
	}

	res := INF
	for x := 0; x <= n; x++ {
		cur := preSumA[x]

		// 前i个小孩选择(n-x)*B[i]，后面的选择C[i]
		f := func(i int) int {
			return preSumB[i]*(n-x) + preSumC[m] - preSumC[i]
		}

		_, y := _FibonacciSearch(f, 0, m+1, true)
		cur += y
		res = min(res, cur)
	}

	fmt.Fprintln(out, res)
}

// 寻找[start,end)中的一个极值点,不要求单峰性质.
// GoldenSectionSearch, 黄金比搜索.
//
//	返回值: (极值点,极值)
func _FibonacciSearch(f func(x int) int, start, end int, minimize bool) (int, int) {
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
