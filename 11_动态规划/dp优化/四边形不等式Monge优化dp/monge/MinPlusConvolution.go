package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	yosupo1()
}

func yosupo1() {
	// https://judge.yosupo.jp/problem/min_plus_convolution_convex_convex
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	numsA := make([]int, n)
	for i := range numsA {
		fmt.Fscan(in, &numsA[i])
	}
	numsB := make([]int, m)
	for i := range numsB {
		fmt.Fscan(in, &numsB[i])
	}

	res := MinPlusConvolution(numsA, numsB, true, false)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

const INF int = 1e18

// 两个数组`和的卷积`最小值, 即 `C[k] = min{A[i]+B[j]} (i+j==k, 0<=k<=n-1+m-1)`
// 至少一个数组是凸函数.
// 这里凸函数的定义是: f(i+2)+f(i) >= f(i+1)+f(i+1) (0<=i<=n-2)
//
//	convexA: numsA是否是凸函数
//	convexB: numsB是否是凸函数
func MinPlusConvolution(numsA, numsB []int, convexA, convexB bool) []int {
	if !convexA && !convexB {
		panic("at least one of A and B must be convex")
	}
	if convexA && convexB {
		return _minPlusConvolutionConvexConvex(numsA, numsB)
	}
	if convexA && !convexB {
		return _minPlusConvolutionArbitraryConvex(numsB, numsA)
	}
	if convexB && !convexA {
		return _minPlusConvolutionArbitraryConvex(numsA, numsB)
	}
	panic("unreachable")
}

// n,m => 5e5 , 570ms
func _minPlusConvolutionConvexConvex(convexA, convexB []int) []int {
	n, m := len(convexA), len(convexB)
	if n == 0 && m == 0 {
		return nil
	}

	res := make([]int, n+m-1)
	for i := range res {
		res[i] = INF
	}
	for n > 0 && convexA[n-1] == INF {
		n--
	}
	for m > 0 && convexB[m-1] == INF {
		m--
	}
	if n == 0 && m == 0 {
		return res
	}

	a, b := 0, 0
	for a < n && convexA[a] == INF {
		a++
	}
	for b < m && convexB[b] == INF {
		b++
	}
	res[a+b] = convexA[a] + convexB[b]

	for i := a + b + 1; i < n+m-1; i++ {
		if b == m-1 || (a != n-1 && convexA[a+1]+convexB[b] < convexA[a]+convexB[b+1]) {
			a++
			_chmin(&res[i], convexA[a]+convexB[b])
		} else {
			b++
			_chmin(&res[i], convexA[a]+convexB[b])
		}
	}
	return res
}

// n,m => 5e5 , 630ms
func _minPlusConvolutionArbitraryConvex(arbitrary, convex []int) []int {
	n, m := len(arbitrary), len(convex)
	if n == 0 && m == 0 {
		return nil
	}
	res := make([]int, n+m-1)
	for i := range res {
		res[i] = INF
	}
	for m > 0 && convex[m-1] == INF {
		m--
	}
	if m == 0 {
		return res
	}
	b := 0
	for b < m && convex[b] == INF {
		b++
	}

	choose := func(i, j, k int) bool {
		if i < k {
			return false
		}
		if i-j >= m-b {
			return true
		}
		return arbitrary[j]+convex[b+i-j] >= arbitrary[k]+convex[b+i-k]
	}

	minArg := _monotoneMinima(n+m-b-1, n, choose)
	for i := 0; i < n+m-b-1; i++ {
		x, y := arbitrary[minArg[i]], convex[b+i-minArg[i]]
		if x < INF && y < INF {
			res[b+i] = x + y
		}
	}
	return res
}

// 寻找二维矩阵中每一行的最小值.
//
//	choose(i,j,k) : (i,j) -> (i,k) 是否可以转移极小值.
func _monotoneMinima(H, W int, choose func(i, j, k int) bool) []int {
	minCol := make([]int, H)
	var dfs func(x1, x2, y1, y2 int)
	dfs = func(x1, x2, y1, y2 int) {
		if x1 == x2 {
			return
		}
		x := (x1 + x2) >> 1
		bestY := y1
		for y := y1 + 1; y < y2; y++ {
			if choose(x, bestY, y) {
				bestY = y
			}
		}
		minCol[x] = bestY
		dfs(x1, x, y1, bestY+1)
		dfs(x+1, x2, bestY, y2)
	}
	dfs(0, H, 0, W)
	return minCol
}

func _chmin(a *int, b int) bool {
	if *a > b {
		*a = b
		return true
	}
	return false
}
