package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc348g()
	// yosupo()
}

// G - Max(Sum - Max)
// https://atcoder.jp/contests/abc348/tasks/abc348_g
// 给定两个长为n的数组A和B.
// 对k=1,2,...,n，选择k个下标，使得∑A[i] - max{B[j]}最大化。
// n<=2e5
func abc348g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	A, B := make([]int, n), make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &A[i])
		fmt.Fscan(in, &B[i])
	}

	order := argSort(B)
	A = reArrage(A, order)
	B = reArrage(B, order)

	var dfs func(int32, int32) ([]int, []int)
	dfs = func(left, right int32) ([]int, []int) {
		if left+1 == right {
			X := []int{0, A[left]}
			Y := []int{-INF, A[left] - B[left]}
			return X, Y
		}
		mid := (left + right) >> 1
		X1, Y1 := dfs(left, mid)
		X2, Y2 := dfs(mid, right)
		n := right - left
		X := make([]int, n+1)
		Y := make([]int, n+1)
		for i := range X {
			X[i] = -INF
			Y[i] = -INF
		}
		for i := range X1 {
			X[i] = X1[i]
		}
		for i := range Y1 {
			Y[i] = Y1[i]
		}

		P := MaxPlusConvolution(X1, X2, true, true)
		Q := MaxPlusConvolution(X1, Y2, true, false)
		for i := range P {
			X[i] = max(X[i], P[i])
		}
		for i := range Q {
			Y[i] = max(Y[i], Q[i])
		}
		return X, Y
	}

	_, Y := dfs(0, n)
	for i := int32(1); i <= n; i++ {
		fmt.Fprintln(out, Y[i])
	}
}

func yosupo() {
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

const INF int = 4e18

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
			res[i] = max(res[i], concaveA[a]+concaveB[b])
		} else {
			b++
			res[i] = max(res[i], concaveA[a]+concaveB[b])
		}
	}
	return res
}

// n,m => 5e5 , 630ms
func _maxPlusConvolutionArbitraryConcave(arbitrary, concave []int) []int {
	n, m := int32(len(arbitrary)), int32(len(concave))
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
	b := int32(0)
	for b < m && concave[b] == -INF {
		b++
	}

	choose := func(i, j, k int32) bool {
		if i < k {
			return false
		}
		if i-j >= m-b {
			return true
		}
		return arbitrary[j]+concave[b+i-j] <= arbitrary[k]+concave[b+i-k]
	}

	maxArg := _monotoneMinima(n+m-b-1, n, choose)
	for i := int32(0); i < n+m-b-1; i++ {
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
func _monotoneMinima(H, W int32, choose func(i, j, k int32) bool) []int32 {
	minCol := make([]int32, H)
	var dfs func(x1, x2, y1, y2 int32)
	dfs = func(x1, x2, y1, y2 int32) {
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func argSort(nums []int) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func reArrage(nums []int, order []int32) []int {
	res := make([]int, len(order))
	for i := range order {
		res[i] = nums[order[i]]
	}
	return res
}
