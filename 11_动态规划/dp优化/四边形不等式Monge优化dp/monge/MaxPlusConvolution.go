package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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

	minusA := make([]int, n)
	for i := 0; i < n; i++ {
		minusA[i] = -numsA[i]
	}
	minusB := make([]int, m)
	for i := 0; i < m; i++ {
		minusB[i] = -numsB[i]
	}

	res := MaxPlusConvolution(minusA, minusB, true, true)
	for _, v := range res {
		fmt.Fprint(out, -v, " ")
	}
}

const INF int = 1e18

// 两个数组`和的卷积`最大值, 即 `C[k] = max{A[i]+B[j]} (i+j==k, 0<=k<=n-1+m-1)`
// 至少一个数组是凹函数.
// 这里凹函数的定义是: f(i+2)+f(i) <= f(i+1)+f(i+1) (0<=i<=n-2)
//
//	concaveA: numsA是否是凹函数
//	concaveB: numsB是否是凹函数
func MaxPlusConvolution(numsA, numsB []int, concaveA, concaveB bool) []int {
	if !concaveA && !concaveB {
		panic("at least one of A and B must be concave")
	}
	if concaveA && concaveB {
		return _maxPlusConvolutionConcaveConcave(numsA, numsB)
	}
	if concaveA && !concaveB {
		return _maxPlusConvolutionArbitraryConcave(numsB, numsA)
	}
	if concaveB && !concaveA {
		return _maxPlusConvolutionArbitraryConcave(numsA, numsB)
	}
	panic("unreachable")
}

// n,m => 5e5 , 570ms
func _maxPlusConvolutionConcaveConcave(concaveA, concaveB []int) []int {
	n, m := len(concaveA), len(concaveB)
	if n == 0 && m == 0 {
		return nil
	}

	res := make([]int, n+m-1)
	for i := range res {
		res[i] = -INF
	}
	for n > 0 && concaveA[n-1] == -INF {
		n--
	}
	for m > 0 && concaveB[m-1] == -INF {
		m--
	}
	if n == 0 && m == 0 {
		return res
	}

	a, b := 0, 0
	for a < n && concaveA[a] == -INF {
		a++
	}
	for b < m && concaveB[b] == -INF {
		b++
	}
	res[a+b] = concaveA[a] + concaveB[b]

	for i := a + b + 1; i < n+m-1; i++ {
		if b == m-1 || (a != n-1 && concaveA[a+1]+concaveB[b] > concaveA[a]+concaveB[b+1]) {
			a++
			_chmax(&res[i], concaveA[a]+concaveB[b])
		} else {
			b++
			_chmax(&res[i], concaveA[a]+concaveB[b])
		}
	}
	return res
}

// n,m => 5e5 , 630ms
func _maxPlusConvolutionArbitraryConcave(arbitrary, concave []int) []int {
	n, m := len(arbitrary), len(concave)
	if n == 0 && m == 0 {
		return nil
	}
	res := make([]int, n+m-1)
	for i := range res {
		res[i] = -INF
	}
	for m > 0 && concave[m-1] == -INF {
		m--
	}
	if m == 0 {
		return res
	}
	b := 0
	for b < m && concave[b] == -INF {
		b++
	}

	choose := func(i, j, k int) bool {
		if i < k {
			return false
		}
		if i-j >= m-b {
			return true
		}
		return arbitrary[j]+concave[b+i-j] <= arbitrary[k]+concave[b+i-k]
	}

	maxArg := _monotoneMinima(n+m-b-1, n, choose)
	for i := 0; i < n+m-b-1; i++ {
		x, y := arbitrary[maxArg[i]], concave[b+i-maxArg[i]]
		if x > -INF && y > -INF {
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

func _chmax(a *int, b int) bool {
	if *a < b {
		*a = b
		return true
	}
	return false
}
