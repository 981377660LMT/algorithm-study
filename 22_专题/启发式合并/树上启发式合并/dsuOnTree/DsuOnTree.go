package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func U41492() {
	// https://www.luogu.com.cn/problem/U41492
	// 给一棵根为0的树，每次询问子树颜色种类数
	const INF int = int(1e18)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
	}

	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	var q int
	fmt.Fscan(in, &q)
	queries := make([][]int, n)
	for i := 0; i < q; i++ {
		var root int
		fmt.Fscan(in, &root)
		root--
		queries[root] = append(queries[root], i)
	}

	id := make(map[int]int, n)
	newValues := make([]int, n)
	for i, v := range values {
		if _, ok := id[v]; !ok {
			id[v] = len(id)
		}
		newValues[i] = id[v]
	}

	res := make([]int, q)
	counter := make([]int, len(id))
	count := 0
	add := func(root int) {
		counter[newValues[root]]++
		if counter[newValues[root]] == 1 {
			count++
		}
	}
	query := func(root int) {
		for _, qi := range queries[root] {
			res[qi] = count
		}
	}
	remove := func(root int) {
		counter[newValues[root]]--
		if counter[newValues[root]] == 0 {
			count--
		}
	}
	reset := func() {}

	dsu := NewDSUonTree(tree, 0)
	dsu.Run(add, remove, query, reset)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type DSUonTree struct {
	g                        [][]int
	n                        int
	subSize, euler, down, up []int
	idx                      int
	root                     int
}

func NewDSUonTree(tree [][]int, root int) *DSUonTree {
	res := &DSUonTree{
		g:       tree,
		n:       len(tree),
		subSize: make([]int, len(tree)),
		euler:   make([]int, len(tree)),
		down:    make([]int, len(tree)),
		up:      make([]int, len(tree)),
		root:    root,
	}

	res.dfs1(root, -1)
	res.dfs2(root, -1)
	return res
}

// add:添加root处的贡献
// remove:退出轻儿子时清除root处的贡献
// query:查询root的子树的贡献并更新答案
// reset:退出轻儿子时重置所有值(如果需要的话)
func (d *DSUonTree) Run(
	add func(root int),
	remove func(root int),
	query func(root int),
	reset func(),
) {
	var dsu func(cur, par int, keep bool)
	dsu = func(cur, par int, keep bool) {
		nexts := d.g[cur]
		for i := 1; i < len(nexts); i++ {
			if to := nexts[i]; to != par {
				dsu(to, cur, false)
			}
		}

		if d.subSize[cur] != 1 {
			dsu(nexts[0], cur, true)
		}

		if d.subSize[cur] != 1 {
			for i := d.up[nexts[0]]; i < d.up[cur]; i++ {
				add(d.euler[i])
			}
		}

		add(cur)
		query(cur)
		if !keep {
			for i := d.down[cur]; i < d.up[cur]; i++ {
				remove(d.euler[i])
			}
			if reset != nil {
				reset()
			}
		}
	}

	dsu(d.root, -1, false)
}

// 每个结点的欧拉序起点.
func (d *DSUonTree) Id(root int) int {
	return d.down[root]
}

func (d *DSUonTree) dfs1(cur, par int) int {
	d.subSize[cur] = 1
	nexts := d.g[cur]
	if len(nexts) >= 2 && nexts[0] == par {
		nexts[0], nexts[1] = nexts[1], nexts[0]
	}
	for i, next := range nexts {
		if next == par {
			continue
		}
		d.subSize[cur] += d.dfs1(next, cur)
		if d.subSize[next] > d.subSize[nexts[0]] {
			nexts[0], nexts[i] = nexts[i], nexts[0]
		}
	}
	return d.subSize[cur]
}

func (d *DSUonTree) dfs2(cur, par int) {
	d.euler[d.idx] = cur
	d.down[cur] = d.idx
	d.idx++
	for _, next := range d.g[cur] {
		if next == par {
			continue
		}
		d.dfs2(next, cur)
	}
	d.up[cur] = d.idx
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

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step`.
func (d *Doubling) MaxStep(from int, check func(e E) bool) int {
	if !d.isPrepared {
		panic("Doubling is not prepared")
	}

	res := d.e()
	step := 0

	for k := d.log - 1; k >= 0; k-- {
		pos := k*d.n + from
		to := d.to[pos]
		next := d.op(res, d.dp[pos])
		if check(next) {
			step |= 1 << k
			from = to
			res = next
		}
	}

	return step
}
