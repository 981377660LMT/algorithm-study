package main

import (
	"bufio"
	"fmt"
	"os"
)

// E - Sensor Optimization Dilemma 2
// https://atcoder.jp/contests/abc374/tasks/abc374_e
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, X int
	fmt.Fscan(in, &N, &X)

	A, P, B, Q := make([]int, N), make([]int, N), make([]int, N), make([]int, N)
	for i := 0; i < N; i++ {
		fmt.Fscan(in, &A[i], &P[i], &B[i], &Q[i])
	}

	check := func(mid int) bool {
		res := 0
		for i := 0; i < N; i++ {
			_, _, c := OptimizeTwoSelection(A[i], P[i], B[i], Q[i], mid)
			res += c
		}
		return res <= X
	}

	left, right := 1, int(2e9)
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}

	fmt.Fprintln(out, right)
}

const INF int = 1e18

// 给定两种物品,每种物品分数为vi,价格为ci.
// 选择任意数量的物品,使得总分数>=w,且总价格最小.
// 返回最小价格 cost, 以及选择的物品数量 x1, x2.
//
// !要么第一种物品选择个数<=v2, 要么第二种物品选择个数<=v1.
// 时间复杂度 O(v1+v2).
func OptimizeTwoSelection(v1 int, c1 int, v2 int, c2 int, w int) (x1 int, x2 int, cost int) {
	x1, x2, cost = 0, 0, INF

	for i1 := 0; i1 <= v2; i1++ {
		curCost := i1 * c1
		remain := max(0, w-i1*v1)
		i2 := (remain + v2 - 1) / v2
		curCost += i2 * c2
		if curCost < cost {
			x1, x2, cost = i1, i2, curCost
		}
	}

	for i2 := 0; i2 <= v1; i2++ {
		curCost := i2 * c2
		remain := max(0, w-i2*v2)
		i1 := (remain + v1 - 1) / v1
		curCost += i1 * c1
		if curCost < cost {
			x1, x2, cost = i1, i2, curCost
		}
	}

	return
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
