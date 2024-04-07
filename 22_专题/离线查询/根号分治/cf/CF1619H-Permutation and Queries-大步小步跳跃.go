// 根号倍增/分块倍增/SqrtLift：牺牲单次查询复杂度以减少单次修改复杂度.
// Permutation and Queries (CF1619H-Permutation and Queries-大步小步跳跃)
// https://www.luogu.com.cn/problem/CF1619H
// 给定一个长为n的排列P,要求支持q次操作:
// 1 x y: 交换P[x]和P[y].
// 2 i k: 给出i的初始值，执行k次操作i = P[i].

// 将原问题转化成有向图，发现是多个置换环
// 维护每个点的前驱后继、跳跃sqrt次后的位置.
// 1, 修改前驱后继时,jump也要修改,改变的只有x和y的前驱sqrt个.
// 2, 查询时,先大步用jump跳跃,再小步用前驱后继跳跃.
// O(nsqrt(n))预处理, O(qsqrt(n))查询.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	perm := make([]int, n)
	for i := range perm {
		fmt.Fscan(in, &perm[i])
		perm[i]--
	}

	// 预处理
	sqrt := int(math.Sqrt(float64(n)))
	prev, next := make([]int, n), make([]int, n)
	jumpSqrt := make([]int, n) // !跳跃sqrt次后的位置
	move := func(pos, k int) int {
		for i := 0; i < k; i++ {
			pos = next[pos]
		}
		return pos
	}
	for i := 0; i < n; i++ {
		next[i] = perm[i]
		prev[next[i]] = i
	}
	for i := 0; i < n; i++ {
		jumpSqrt[i] = move(i, sqrt)
	}

	adjustJumpSqrt := func(pos int) {
		posNext := next[pos]
		for i := 0; i < sqrt-1; i++ {
			pos = prev[pos]
		}
		// 修改前驱的jumpSqrt
		for i := 0; i < sqrt; i++ {
			jumpSqrt[pos] = posNext
			pos = next[pos]
			posNext = next[posNext]
		}
	}

	// 交换两个元素.
	swap := func(x, y int) {
		prev[next[x]], prev[next[y]] = prev[next[y]], prev[next[x]]
		next[x], next[y] = next[y], next[x]
		adjustJumpSqrt(x)
		adjustJumpSqrt(y)
	}

	// 从pos开始跳跃k次.
	query := func(pos, k int) int {
		for k >= sqrt {
			pos = jumpSqrt[pos]
			k -= sqrt
		}
		return move(pos, k)
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x, y int
			fmt.Fscan(in, &x, &y)
			x--
			y--
			swap(x, y)
		} else {
			var pos, k int
			fmt.Fscan(in, &pos, &k)
			pos--
			fmt.Fprintln(out, query(pos, k)+1)
		}
	}
}
