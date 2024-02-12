// https://maspypy.github.io/library/ds/doubling.hpp
// 倍增doubling

// !状態 a から 1 回操作すると、状態 b に遷移し、モノイドの元 x を加える。
//  行き先がない場合：-1 （add 不要）

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	CF1175E()
}

// Minimal Segment Cover (线段包含/线段覆盖)
// https://www.luogu.com.cn/problem/CF1175E
// 给定 n 个线段，每个线段形如 [l,r]，
// 有 q 次查询，每次查询形如 [L,R]，求出所有线段的并集中包含 [L,R] 的最小线段数。
// n,m<=2e5,0<=l<=r<=5e5
//
// !预处理出从一个左端点用一条线段最远能覆盖到哪个右端点,注意线段起点可能在这个左端点左边
// 然后就可以倍增预处理出每个左端点用x条线段最远能覆盖到哪个右端点
func CF1175E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	segments := make([][2]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &segments[i][0], &segments[i][1])
	}
	queries := make([][2]int, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &queries[i][0], &queries[i][1])
	}

	MAX := int(5e5 + 10)
	maxRight := make([]int, MAX) // 每个左端点用一条线段最远能覆盖到哪个右端点
	for i := 0; i < n; i++ {
		maxRight[i] = i
	}
	for i := 0; i < n; i++ {
		a, b := segments[i][0], segments[i][1]
		maxRight[a] = max(maxRight[a], b)
	}
	for i := 1; i < MAX; i++ {
		maxRight[i] = max(maxRight[i], maxRight[i-1]) // 注意线段起点可能在这个左端点左边
	}

	D := NewDoubling(MAX, MAX, func() E { return 0 }, func(e1, e2 E) E { return e1 + e2 })
	for i := 0; i < MAX; i++ {
		D.Add(i, maxRight[i], maxRight[i]-i)
	}
	D.Build()

	for i := 0; i < q; i++ {
		a, b := queries[i][0], queries[i][1]
		k, _, _ := D.MaxStep(a, func(e E) bool { return a+e < b })
		k++
		if k > n {
			k = -1
		}
		fmt.Fprintln(out, k)
	}
}

func yuki1097() {
	// https://yukicoder.me/problems/no/1097
	// 给定一个数组和q次查询
	// 初始时res为0，每次查询会执行k次操作:
	//  将res加上nums[res%n]的值
	// 求每次查询后res的值
	// (n,q<=2e5,k<=1e12)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	db := NewDoubling(
		n, 1e12+10,
		func() E { return 0 },
		func(e1, e2 E) E { return e1 + e2 },
	)
	for i := 0; i < n; i++ {
		db.Add(i, (i+nums[i])%n, nums[i]) // res的模从i变为(i+nums[i])%n，res加上nums[i]
	}
	db.Build()

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var k int
		fmt.Fscan(in, &k)
		_, res := db.Jump(0, k)
		fmt.Fprintln(out, res)
	}
}

// 8027. 在传球游戏中最大化函数值
func getMaxFunctionValue(receiver []int, k int64) int64 {
	n := len(receiver)
	intK := int(k)
	db := NewDoubling(
		n, intK+1,
		func() E { return 0 },
		func(e1, e2 E) E { return e1 + e2 },
	)
	for i := 0; i < n; i++ {
		db.Add(i, receiver[i], i)
	}
	db.Build()

	res := 0
	for i := 0; i < n; i++ {
		_, v := db.Jump(i, intK+1)
		if v > res {
			res = v
		}
	}

	return int64(res)
}

type E = int

type Doubling struct {
	n          int
	log        int
	isPrepared bool
	to         []int
	dp         []E
	e          func() E
	op         func(e1, e2 E) E
}

func NewDoubling(n, maxStep int, e func() E, op func(e1, e2 E) E) *Doubling {
	res := &Doubling{e: e, op: op}
	res.n = n
	res.log = bits.Len(uint(maxStep))
	size := n * res.log
	res.to = make([]int, size)
	res.dp = make([]E, size)
	for i := 0; i < size; i++ {
		res.to[i] = -1
		res.dp[i] = res.e()
	}
	return res
}

// 初始状态(leaves):从 `from` 状态到 `to` 状态，边权为 `weight`.
//
//	0 <= from, to < n
func (d *Doubling) Add(from, to int, weight E) {
	if d.isPrepared {
		panic("Doubling is prepared")
	}
	if to < -1 || to >= d.n {
		panic("to is out of range")
	}

	d.to[from] = to
	d.dp[from] = weight
}

func (d *Doubling) Build() {
	if d.isPrepared {
		panic("Doubling is prepared")
	}

	d.isPrepared = true
	n := d.n
	for k := 0; k < d.log-1; k++ {
		for v := 0; v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				d.dp[next] = d.dp[k*n+v]
				continue
			}
			d.to[next] = d.to[k*n+w]
			d.dp[next] = d.op(d.dp[k*n+v], d.dp[k*n+w])
		}
	}
}

// 从 `from` 状态开始，执行 `step` 次操作，返回最终状态的编号和操作的结果。
//
//	0 <= from < n
//	如果最终状态不存在，返回 -1, e()
func (d *Doubling) Jump(from, step int) (to int, res E) {
	if !d.isPrepared {
		panic("Doubling is not prepared")
	}
	if step >= 1<<d.log {
		panic("step is over max step")
	}

	res = d.e()
	to = from
	for k := 0; k < d.log; k++ {
		if to == -1 {
			break
		}

		if step&(1<<k) != 0 {
			pos := k*d.n + to
			res = d.op(res, d.dp[pos])
			to = d.to[pos]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号和操作的结果。
func (d *Doubling) MaxStep(from int, check func(weight E) bool) (step int, to int, res E) {
	if !d.isPrepared {
		panic("Doubling is not prepared")
	}

	res = d.e()
	for k := d.log - 1; k >= 0; k-- {
		pos := k*d.n + from
		tmp := d.to[pos]
		next := d.op(res, d.dp[pos])
		if check(next) {
			step |= 1 << k
			from = tmp
			res = next
		}
	}

	to = from
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
