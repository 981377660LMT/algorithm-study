package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
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
	update := func(root int) {
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
	clear := func(root int) {
		counter[newValues[root]]--
		if counter[newValues[root]] == 0 {
			count--
		}
	}
	reset := func() {}

	dsu := NewDSUonTree(tree, 0)
	dsu.Run(update, query, clear, reset)
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
				update(d.euler[i])
			}
		}

		update(cur)
		query(cur)
		if !keep {
			for i := d.down[cur]; i < d.up[cur]; i++ {
				clear(d.euler[i])
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
