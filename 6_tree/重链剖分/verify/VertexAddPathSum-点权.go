// VertexAddPathSum

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {

	// https://judge.yosupo.jp/problem/vertex_add_path_sum
	// 单点加/路径和查询
	// 0 vertex add => 顶点加
	// 1 root1 root2 => 路径和查询
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	hld := NewHeavyLightDecomposition(n)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	for i := 0; i < n-1; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		hld.AddEdge(u, v)
	}
	hld.Build(0)

	bit := NewBITArray(n)
	for i := 0; i < n; i++ {
		start, _ := hld.Id(i)
		bit.Add(start, start+1, values[i])
	}

	for i := 0; i < q; i++ {
		var op, v, add, root1, root2 int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &v, &add)
			start, _ := hld.Id(v)
			bit.Add(start, start+1, add)
			// hld.QueryPath(v, v, true, func(start, end int) {
			// 	bit.Add(start, end, add)
			// })
		} else {
			fmt.Fscan(in, &root1, &root2)
			res := 0
			hld.QueryPath(root1, root2, true, func(start, end int) {
				res += bit.Query(start, end)
			})
			fmt.Fprintln(out, res)
		}
	}
}

type HeavyLightDecomposition struct {
	Parent   []int
	Depth    []int
	Size     []int
	g        [][]int
	id       int
	down, up []int
	nxt      []int // heavy pathの先頭
}

func NewHeavyLightDecomposition(n int) *HeavyLightDecomposition {
	return &HeavyLightDecomposition{g: make([][]int, n)}
}

// 無向辺 u <-> v を追加する.
func (hld *HeavyLightDecomposition) AddEdge(u, v int) {
	hld.g[u] = append(hld.g[u], v)
	hld.g[v] = append(hld.g[v], u)
}

// 有向辺 u -> v を追加する.
func (hld *HeavyLightDecomposition) AddDirectedEdge(u, v int) {
	hld.g[u] = append(hld.g[u], v)
}

// rootを根とした重軽分解を構築する.
func (hld *HeavyLightDecomposition) Build(root int) {
	n := len(hld.g)
	hld.Size = make([]int, n)
	hld.Depth = make([]int, n)
	hld.down = make([]int, n)
	hld.up = make([]int, n)
	hld.nxt = make([]int, n)
	hld.Parent = make([]int, n)
	for i := 0; i < n; i++ {
		hld.down[i] = -1
		hld.up[i] = -1
		hld.nxt[i] = root
		hld.Parent[i] = root
	}

	hld.dfsSize(root, -1)
	hld.dfsHld(root, -1)
}

// 頂点 i のオイラーツアー順を [down,up) の形で返す.
//  0 <= down < up <= n.
func (hld *HeavyLightDecomposition) Id(u int) (down, up int) {
	down, up = hld.down[u], hld.up[u]
	return
}

// 可換なパスクエリを処理する.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QueryPath(u, v int, vertex bool, f func(start, end int)) {
	lca_ := hld.LCA(u, v)
	for _, p := range hld.ascend(u, lca_) {
		s, t := p[0]+1, p[1]
		if s > t {
			f(t, s)
		} else {
			f(s, t)
		}
	}
	if vertex {
		f(hld.down[lca_], hld.down[lca_]+1)
	}
	for _, p := range hld.descend(lca_, v) {
		s, t := p[0], p[1]+1
		if s > t {
			f(t, s)
		} else {
			f(s, t)
		}
	}
}

// 非可換なパスクエリを処理する.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QueryNonCommutativePath(u, v int, vertex bool, f func(start, end int)) {
	lca_ := hld.LCA(u, v)
	for _, p := range hld.ascend(u, lca_) {
		f(p[0]+1, p[1])
	}
	if vertex {
		f(hld.down[lca_], hld.down[lca_]+1)
	}
	for _, p := range hld.descend(lca_, v) {
		f(p[0], p[1]+1)
	}
}

// 部分木クエリを処理する.
//   0 <= start <= end <= n, [start,end).
func (hld *HeavyLightDecomposition) QuerySubTree(u int, vertex bool, f func(start, end int)) {
	if vertex {
		f(hld.down[u], hld.up[u])
	} else {
		f(hld.down[u]+1, hld.up[u])
	}
}

func (hld *HeavyLightDecomposition) LCA(u, v int) int {
	for hld.nxt[u] != hld.nxt[v] {
		if hld.down[u] < hld.down[v] {
			u, v = v, u
		}
		u = hld.Parent[hld.nxt[u]]
	}
	if hld.Depth[u] < hld.Depth[v] {
		return u
	}
	return v
}

func (hld *HeavyLightDecomposition) Dist(u, v int) int {
	return hld.Depth[u] + hld.Depth[v] - hld.Depth[hld.LCA(u, v)]*2
}

func (hld *HeavyLightDecomposition) dfsSize(cur, pre int) {
	hld.Size[cur] = 1
	for i, to := range hld.g[cur] {
		if to == pre {
			continue
		}
		// if to == hld.Parent[cur] {
		// 	if len(hld.g[cur]) >= 2 && hld.g[cur][0] == to {
		// 		hld.g[cur][0], hld.g[cur][1] = hld.g[cur][1], hld.g[cur][0]
		// 	} else {
		// 		continue
		// 	}
		// }

		hld.Depth[to] = hld.Depth[cur] + 1
		hld.Parent[to] = cur
		hld.dfsSize(to, cur)
		hld.Size[cur] += hld.Size[to]
		if hld.Size[to] > hld.Size[hld.g[cur][0]] {
			hld.g[cur][0], hld.g[cur][i] = hld.g[cur][i], hld.g[cur][0]
		}
	}

}

func (hld *HeavyLightDecomposition) dfsHld(cur, pre int) {
	hld.down[cur] = hld.id
	hld.id++
	for _, to := range hld.g[cur] {
		if to == pre {
			continue
		}
		if to == hld.g[cur][0] {
			hld.nxt[to] = hld.nxt[cur]
		} else {
			hld.nxt[to] = to
		}
		hld.dfsHld(to, cur)
	}
	hld.up[cur] = hld.id
}

// [u, v)
func (hld *HeavyLightDecomposition) ascend(u, v int) [][2]int {
	var res [][2]int
	for hld.nxt[u] != hld.nxt[v] {
		res = append(res, [2]int{hld.down[u], hld.down[hld.nxt[u]]})
		u = hld.Parent[hld.nxt[u]]
	}
	if u != v {
		res = append(res, [2]int{hld.down[u], hld.down[v] + 1})
	}
	return res
}

// (u, v]
func (hld *HeavyLightDecomposition) descend(u, v int) [][2]int {
	if u == v {
		return nil
	}
	if hld.nxt[u] == hld.nxt[v] {
		return [][2]int{{hld.down[u] + 1, hld.down[v]}}
	}
	res := hld.descend(u, hld.Parent[hld.nxt[v]])
	res = append(res, [2]int{hld.down[hld.nxt[v]], hld.down[v]})
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// !Range Add Range Sum, 0-based.
type BITArray struct {
	n     int
	tree1 []int
	tree2 []int
}

func NewBITArray(n int) *BITArray {
	return &BITArray{
		n:     n,
		tree1: make([]int, n+1),
		tree2: make([]int, n+1),
	}
}

// 切片内[start, end)的每个元素加上delta.
//  0<=start<=end<=n
func (b *BITArray) Add(start, end, delta int) {
	end--
	b.add(start, delta)
	b.add(end+1, -delta)
}

// 求切片内[start, end)的和.
//  0<=start<=end<=n
func (b *BITArray) Query(start, end int) int {
	end--
	return b.query(end) - b.query(start-1)
}

func (b *BITArray) add(index, delta int) {
	index++
	rawIndex := index
	for index <= b.n {
		b.tree1[index] += delta
		b.tree2[index] += (rawIndex - 1) * delta
		index += index & -index
	}
}

func (b *BITArray) query(index int) (res int) {
	index++
	if index > b.n {
		index = b.n
	}
	rawIndex := index
	for index > 0 {
		res += rawIndex*b.tree1[index] - b.tree2[index]
		index -= index & -index
	}
	return
}
