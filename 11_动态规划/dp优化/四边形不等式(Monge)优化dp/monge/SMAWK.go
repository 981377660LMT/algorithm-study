// https://maspypy.github.io/library/convex/smawk.hpp
// https://noshi91.github.io/Library/algorithm/smawk.cpp
// SMAWK算法是一种用于在隐式定义的完全单调矩阵(totally monge)的每一行中查找最小值的算法
// monotone: 每一行取得最值的列号是单调的
// totally monotone: 对任意的2x2子矩阵, 取得最值的列号是单调的

package main

import "fmt"

func main() {
	M := [][]int{
		{0, 1, 3, 2, 4},
		{0, 2, 4, 3, 1},
		{1, 3, 4, 2, 0},
		{4, 2, 3, 1, 0},
	}
	H, W := len(M), len(M[0])
	choose := func(i, j, k int) int {
		if M[i][j] > M[i][k] {
			return k
		}
		return j
	}
	I := SMAWK(H, W, choose)
	if fmt.Sprint(I) != "[0 0 4 4]" {
		panic("I != [0 0 4 4]")
	}

}

// choose: func(i, j, k int) int 选择(i,j)和(i,k)中的哪一个(j or k)
//  返回值: argMin[i] 表示第i行的最小值的列号
func SMAWK(H, W int, choose func(i, j, k int) int) (argMin []int) {

	var dfs func(X, Y []int) []int
	dfs = func(X, Y []int) (I []int) {
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
		I = make([]int, N)
		for i := range II {
			I[i+i+1] = II[i]
		}
		p := 0
		for i := 0; i < N; i += 2 {
			LIM := Y[len(Y)-1]
			if i+1 < N {
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
