// TODO 有问题
package main

import (
	"fmt"
	"math/rand"
	"sort"
)

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

func main() {
	rand.Seed(0)
	gen := func(L, N, R int) []int {
		A := make([]int, N)
		for i := 0; i < N; i++ {
			A[i] = rand.Intn(200) - 100
		}
		sort.Sort(sort.Reverse(sort.IntSlice(A)))
		preSum := make([]int, N+1)
		for i := 0; i < N; i++ {
			preSum[i+1] = preSum[i] + A[i]
		}
		for i := 0; i < L; i++ {
			A = append([]int{INF}, A...)
		}
		for i := 0; i < R; i++ {
			A = append(A, INF)
		}
		return A
	}

	equal := func(A, B []int) bool {
		if len(A) != len(B) {
			return false
		}
		for i := range A {
			if A[i] != B[i] {
				return false
			}
		}
		return true
	}

	naive := func(A, B []int) []int {
		N := len(A)
		M := len(B)
		if N == 0 || M == 0 {
			return []int{}
		}
		C := make([]int, N+M-1)
		for i := range C {
			C[i] = INF
		}
		for i := 0; i < N; i++ {
			for j := 0; j < M; j++ {
				if A[i] == INF || B[j] == INF {
					continue
				}
				C[i+j] = min(C[i+j], A[i]+B[j])
			}
		}
		return C
	}
	// [93 52] [69 -38 -74 -1000000000000000000 -1000000000000000000 -1000000000000000000 -1000000000000000000]
	fmt.Println(MinPlusConvolutionConvex([]int{93, 52}, []int{69, -38, -74, -INF, -INF, -INF, -INF}, true, true))
	for a1 := 0; a1 < 5; a1++ {
		for b1 := 0; b1 < 10; b1++ {
			for c1 := 0; c1 < 5; c1++ {
				A := gen(a1, b1, c1)
				for a2 := 0; a2 < 5; a2++ {
					for b2 := 0; b2 < 10; b2++ {
						for c2 := 0; c2 < 5; c2++ {
							B := gen(a2, b2, c2)
							C := MinPlusConvolutionConvex(A, B, true, true)
							if !equal(naive(A, B), C) {
								fmt.Println(A, B, C, naive(A, B))
								panic("error")
							}
						}
					}
				}
			}
		}
	}
}

const INF int = 1e18

func MinPlusConvolutionConvex(A, B []int, convexA, convexB bool) (C []int) {
	if len(A) == 0 || len(B) == 0 {
		return
	}
	if !convexA && !convexB {
		panic("at least one of A and B must be convex")
	}
	if !convexB {
		A, B = B, A
	}
	NA := len(A)
	NB := len(B)
	N := NA + NB - 1
	L := 0
	R := NB
	for L < R && B[L] == INF {
		L++
	}
	if L == R {
		C = make([]int, N)
		for i := range C {
			C[i] = INF
		}
		return
	}
	for B[R-1] == INF {
		R--
	}
	B = B[L:R]
	nB := R - L
	n := NA + nB - 1

	choose := func(i, j, k int) int {
		if i < k {
			return j
		}
		if i-j >= nB {
			return k
		}
		if A[j]+B[i-j] > A[k]+B[i-k] {
			return k
		}
		return j
	}

	J := _SMAWK(n, NA, choose)
	C = make([]int, N)
	for i := range C {
		C[i] = INF
	}
	for i := 0; i < n; i++ {
		if A[J[i]] == INF {
			C[L+i] = INF
		} else {
			C[L+i] = A[J[i]] + B[i-J[i]]
		}
	}
	return
}

// choose: func(i, j, k int) int 选择(i,j)和(i,k)中的哪一个(j or k)
//  返回值: argMin[i] 表示第i行的最小值的列号
func _SMAWK(H, W int, choose func(i, j, k int) int) (argMin []int) {
	var dfs func(X, Y []int) []int
	dfs = func(X, Y []int) []int {
		N := len(X)
		if N == 0 {
			return nil
		}
		YY := []int{}
		for _, y := range Y {
			for len(YY) > 0 {
				py := YY[len(YY)-1]
				x := X[len(YY)-1]
				if choose(x, py, y) == py {
					break
				}
				YY = YY[:len(YY)-1]
			}
			if len(YY) < len(X) {
				YY = append(YY, y)
			}
		}

		XX := []int{}
		for i := 1; i < len(X); i += 2 {
			XX = append(XX, X[i])
		}
		II := dfs(XX, YY)
		I := make([]int, N)
		for i := range II {
			I[i+i+1] = II[i]
		}
		p := 0

		for i := 0; i < N; i += 2 {
			var LIM int
			if i+1 == N {
				LIM = Y[len(Y)-1]
			} else {
				LIM = I[i+1]
			}
			best := Y[p]
			for Y[p] < LIM {
				p++
				best = choose(X[i], best, Y[p])
			}
			I[i] = best
		}
		return I
	}

	X := make([]int, H)
	Y := make([]int, W)
	for i := range X {
		X[i] = i
	}
	for i := range Y {
		Y[i] = i
	}
	return dfs(X, Y)
}
