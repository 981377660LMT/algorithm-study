// TODO 树上差分
// https://blog.csdn.net/justidle/article/details/104508212
// 边差分、点差分

// 树上差分，就是利用差分的性质，对路径上的重要节点进行修改（而不是暴力全改），
// 作为其差分数组的值，最后在求值时，利用 dfs 遍历求出差分数组的前缀和，
// 就可以达到降低复杂度的目的。树上差分时需要求 LCA

// update
// build
// get

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {

}

// 每个点的子树中与它距离小于等于k的点有多少个

// https://www.luogu.com.cn/problem/P3066
// 给出以0号点为根的一棵有根树,问每个点的子树中与它距离小于等于k的点有多少个
// 倍增+差分区间加
// TODO
func P3066() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)

	tree := make([][]int, n)
	db := NewDoubling(n, n)
	for i := 1; i < n; i++ {
		var parent, w int
		fmt.Fscan(in, &parent, &w)
		parent--
		tree[parent] = append(tree[parent], i)
		db.Add(i, parent, w)
	}
	db.Build()

	diff := make([]int, n)
	var dfs func(cur, pre int, dist int)
	dfs = func(cur, pre int, dist int) {
		_, maxUp, _ := db.MaxStep(cur, func(e int) bool { return e <= k })
	}

}

type DiffOnTree struct {
	tree [][]int
}

func NewDiffOnTree(n int, tree [][]int) *DiffOnTree {
	d := &DiffOnTree{make([][]int, n)}
	return d
}

type E = int

// monoidAdd
func (*Doubling) e() E          { return 0 }
func (*Doubling) op(e1, e2 E) E { return e1 + e2 }

type Doubling struct {
	n          int
	log        int
	isPrepared bool
	to         []int
	dp         []E
}

func NewDoubling(n, maxStep int) *Doubling {
	res := &Doubling{}
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
func (d *Doubling) MaxStep(from int, check func(e E) bool) (step int, to int, res E) {
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
