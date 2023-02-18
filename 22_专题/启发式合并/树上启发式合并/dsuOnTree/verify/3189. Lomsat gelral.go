package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/CF600E#submit
	// 求子树内`出现颜色次数最多`的颜色的和
	const INF int = int(1e18)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}

	tree := make([][]Edge, n)
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		tree[u] = append(tree[u], Edge{u, v, 1})
		tree[v] = append(tree[v], Edge{v, u, 1})
	}

	res := make([]int, n)
	majorSum := NewMajorSum()
	update := func(root int) {
		c := values[root]
		majorSum.Add(c)
	}
	query := func(root int) {
		res[root] = majorSum.Query()
	}
	clear := func(root int) {
		c := values[root]
		majorSum.Discard(c)
	}
	reset := func() {}

	dsu := NewDSUonTree(tree, 0)
	dsu.Run(update, query, clear, reset)
	for _, v := range res {
		fmt.Fprint(out, v, " ")
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

type MajorSum struct {
	maxFreq   int
	res       int
	counter   map[int]int
	freqSum   map[int]int
	freqTypes map[int]int
}

func NewMajorSum() *MajorSum {
	return &MajorSum{
		counter:   make(map[int]int),
		freqSum:   make(map[int]int),
		freqTypes: make(map[int]int),
	}
}

func (ms *MajorSum) Add(x int) {
	ms.counter[x]++
	xFreq := ms.counter[x]
	ms.freqSum[xFreq] += x
	ms.freqSum[xFreq-1] -= x
	ms.freqTypes[xFreq]++
	ms.freqTypes[xFreq-1]--
	if xFreq > ms.maxFreq {
		ms.maxFreq = xFreq
		ms.res = x
	} else if xFreq == ms.maxFreq {
		ms.res += x
	}
}

func (ms *MajorSum) Discard(x int) {
	if ms.counter[x] == 0 {
		return
	}
	ms.counter[x]--
	xFreq := ms.counter[x]
	ms.freqSum[xFreq] += x
	ms.freqSum[xFreq+1] -= x
	ms.freqTypes[xFreq]++
	ms.freqTypes[xFreq+1]--
	if xFreq+1 == ms.maxFreq {
		ms.res -= x
		if ms.freqTypes[ms.maxFreq] == 0 {
			ms.maxFreq--
			ms.res = ms.freqSum[ms.maxFreq]
		}
	}
	if ms.counter[x] == 0 {
		delete(ms.counter, x)
	}
}

func (ms *MajorSum) Query() int {
	return ms.res
}
