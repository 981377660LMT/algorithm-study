// DSU On Tree
// 这里是利用欧拉序处理子树

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	// https://judge.yosupo.jp/problem/vertex_add_subtree_sum
	// 0 u x : add x to u
	// 1 u : print sum of subtree of u
	const INF int = int(1e18)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	tree := make([][]Edge, n)
	for i := 1; i < n; i++ {
		var parent int
		fmt.Fscan(in, &parent)
		tree[parent] = append(tree[parent], Edge{parent, i, 1})
	}

	updates := make([][][2]int, n)
	queries := make([][]int, n)
	for i := 0; i < n; i++ {
		updates[i] = append(updates[i], [2]int{0, nums[i]})
	}
	for i := 0; i < q; i++ {
		var op, root int
		fmt.Fscan(in, &op, &root)
		if op == 0 {
			var add int
			fmt.Fscan(in, &add)
			updates[root] = append(updates[root], [2]int{i, add})
		} else {
			queries[root] = append(queries[root], i)
		}
	}

	bit := NewBitArray(q)
	res := make([]int, q)
	for i := 0; i < q; i++ {
		res[i] = INF
	}

	update := func(root int) {
		for _, v := range updates[root] {
			bit.Apply(v[0], v[1])
		}
	}
	query := func(root int) {
		for _, v := range queries[root] {
			res[v] = bit.Prod(v + 1)
		}
	}
	clear := func(root int) {
		for _, v := range updates[root] {
			bit.Apply(v[0], -v[1])
		}
	}
	reset := func() {}

	dsu := NewDSUonTree(tree, 0)
	dsu.Run(update, query, clear, reset)
	for _, v := range res {
		if v == INF {
			continue
		}
		fmt.Fprintln(out, v)
	}
}

type DSUonTree struct {
	g                        [][]Edge
	n                        int
	subSize, euler, down, up []int
	idx                      int
	root                     int
}

type Edge struct{ from, to, weight int }

func NewDSUonTree(tree [][]Edge, root int) *DSUonTree {
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

// update:添加root处的贡献
// query:查询root的子树的贡献并更新答案
// clear:退出轻儿子时清除root处的贡献
// reset:退出轻儿子时重置所有值(如果需要的话)
func (d *DSUonTree) Run(
	update func(root int),
	query func(root int),
	clear func(root int),
	reset func(),
) {
	var dsu func(cur, par int, keep bool)
	dsu = func(cur, par int, keep bool) {
		for i := 1; i < len(d.g[cur]); i++ {
			if to := d.g[cur][i].to; to != par {
				dsu(to, cur, false)
			}
		}

		if d.subSize[cur] != 1 {
			dsu(d.g[cur][0].to, cur, true)
		}

		if d.subSize[cur] != 1 {
			for i := d.up[d.g[cur][0].to]; i < d.up[cur]; i++ {
				update(d.euler[i])
			}
		}

		update(cur)
		query(cur)
		if !keep {
			for i := d.down[cur]; i < d.up[cur]; i++ {
				clear(d.euler[i])
			}
			reset()
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
	if len(d.g[cur]) >= 2 && d.g[cur][0].to == par {
		d.g[cur][0], d.g[cur][1] = d.g[cur][1], d.g[cur][0]
	}
	for i := range d.g[cur] {
		next := d.g[cur][i].to
		if next == par {
			continue
		}
		d.subSize[cur] += d.dfs1(next, cur)
		if d.subSize[next] > d.subSize[d.g[cur][0].to] {
			d.g[cur][0], d.g[cur][i] = d.g[cur][i], d.g[cur][0]
		}
	}
	return d.subSize[cur]
}

func (d *DSUonTree) dfs2(cur, par int) {
	d.euler[d.idx] = cur
	d.down[cur] = d.idx
	d.idx++
	for i := range d.g[cur] {
		next := d.g[cur][i].to
		if next == par {
			continue
		}
		d.dfs2(next, cur)
	}
	d.up[cur] = d.idx
}

type BitArray struct {
	n    int
	log  int
	data []int
}

// 長さ n の 0で初期化された配列で構築する.
func NewBitArray(n int) *BitArray {
	return &BitArray{n: n, log: bits.Len(uint(n)), data: make([]int, n+1)}
}

// 配列で構築する.
func NewBitArrayFrom(arr []int) *BitArray {
	res := NewBitArray(len(arr))
	res.Build(arr)
	return res
}

func (b *BitArray) Build(arr []int) {
	if b.n != len(arr) {
		panic("len of arr is not equal to n")
	}
	for i := 1; i <= b.n; i++ {
		b.data[i] = arr[i-1]
	}
	for i := 1; i <= b.n; i++ {
		j := i + (i & -i)
		if j <= b.n {
			b.data[j] += b.data[i]
		}
	}
}

// 要素 i に値 v を加える.
func (b *BitArray) Apply(i int, v int) {
	for i++; i <= b.n; i += i & -i {
		b.data[i] += v
	}
}

// [0, r) の要素の総和を求める.
func (b *BitArray) Prod(r int) int {
	res := int(0)
	for ; r > 0; r -= r & -r {
		res += b.data[r]
	}
	return res
}

// [l, r) の要素の総和を求める.
func (b *BitArray) ProdRange(l, r int) int {
	return b.Prod(r) - b.Prod(l)
}

// 区間[0,k]の総和がx以上となる最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) LowerBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] < x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

// 区間[0,k]の総和がxを上回る最小のkを求める.数列が単調増加であることを要求する.
func (b *BitArray) UpperBound(x int) int {
	i := 0
	for k := 1 << b.log; k > 0; k >>= 1 {
		if i+k <= b.n && b.data[i+k] <= x {
			x -= b.data[i+k]
			i += k
		}
	}
	return i
}

func (b *BitArray) String() string {
	sb := []string{}
	for i := 0; i < b.n; i++ {
		sb = append(sb, fmt.Sprintf("%d", b.ProdRange(i, i+1)))
	}
	return fmt.Sprintf("BitArray: [%v]", strings.Join(sb, ", "))
}
