// https://oi-wiki.org/graph/block-forest/
// 例题:
// 有一张 n 个点 m 条边的无向连通图，还有 q 个点对，
// 你需要输出每个点是多少给定点对的必经点（即如果点对为 (u,v)，
// 那么如果 u 到 v 无论如何都要经过 x ，那么 x 是该点对的必经点）
// !直接建出圆方树，发现(u,v) 在圆方树路径上的圆点都是必经点，lca 树上差分一下就可以了。
// !如果原图中某个连通分量只有一个点，则需要具体情况具体分析，我们在后续讨论中不考虑孤立点。

// https://ei1333.hateblo.jp/entry/2020/03/25/010057
// !实装与oiwiki有些不同，这里是把每个点双缩成了一个点
// Block-cut tree上の頂点番号の割り当て方なんですが,
// !二重頂点連結成分ならその連結成分の番号,
// !関節点なら (二重頂点連結成分の個数)+何番目の関節点かをしています.
//  !即编号大于等于二重连通分量数的点都是割点 (即 bct.Rev[i] >= len(bct.BCC))
// Rev: 原图的每个顶点对应的圆方树的顶点编号
// Group[i]: 圆方树的每个顶点i对应的原图的顶点编号们.(如果是割点就对应一个点,如果是双连通分量就对应多个点)
// Tree: 圆方树

package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// https://judge.u-aizu.ac.jp/onlinejudge/description.jsp?id=3022
	// 网络连接,对一个网络集群维修,每天维修一台机器,需要断开它和其他机器的连接
	// !对每个顶点,移除他连接的所有边后,求剩下的连通分量的权值和的最大值(総合性能値)
	// 圆方树上dp
	// 如果i不是割点的话,就是原图中的所有点的权值和减去这个点的权值
	// 如果i是割点的话,就是max(子树和,allSum-子树和-这个点的权值)

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	values := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &values[i])
	}
	graph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u--
		v--
		graph[u] = append(graph[u], Edge{u, v, 1, i})
		graph[v] = append(graph[v], Edge{v, u, 1, i})
	}

	bct := NewBlockCutTree(graph)
	bct.Build()
	weight := make([]int, len(bct.Tree)) // 每个圆方树的顶点的权值和
	all := 0
	for i := 0; i < n; i++ {
		weight[bct.Get(i)] += values[i]
		all += values[i]
	}

	res := make([]int, n)
	for i := 0; i < n; i++ {
		res[i] = all - values[i] // 如果i不是割点的话,就是原图中的所有点的权值和减去这个点的权值
	}

	var dfs func(cur, parent int) int
	dfs = func(cur, parent int) int {
		max_, sum := 0, 0
		for _, e := range bct.Tree[cur] {
			if e.to == parent {
				continue
			}
			nextRes := dfs(e.to, cur)
			sum += nextRes
			max_ = max(max_, nextRes)
		}
		if bct.IsCut(cur) {
			res[bct.Group[cur][0]] = max(max_, all-sum-weight[cur])
		}
		return sum + weight[cur]
	}
	dfs(0, -1)

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, res[i])
	}
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type Edge = struct{ from, to, cost, index int }
type Graph = [][]Edge
type BlockCutTree struct {
	Tree  Graph   // 圆方树
	Rev   []int   // 原图的每个顶点对应的圆方树的顶点编号
	Group [][]int // 圆方树的每个顶点i对应的原图的顶点编号们.(如果是割点就对应一个点,如果是双连通分量就对应多个点)
	g     [][]Edge
	bcc   *BiConnectedComponents
}

func NewBlockCutTree(g Graph) *BlockCutTree {
	return &BlockCutTree{g: g, bcc: NewBiConnectedComponents(g)}
}

func (bct *BlockCutTree) Build() {
	bct.bcc.Build()
	bct.Rev = make([]int, len(bct.g))
	for i := range bct.Rev {
		bct.Rev[i] = -1
	}
	ptr := len(bct.bcc.BCC)
	for _, idx := range bct.bcc.lowLink.Articulation {
		bct.Rev[idx] = ptr
		ptr++
	}
	last := make([]int, ptr)
	for i := range last {
		last[i] = -1
	}
	bct.Tree = make(Graph, ptr)
	for i := range bct.bcc.BCC {
		for _, e := range bct.bcc.BCC[i] {
			for _, ver := range []int{e.from, e.to} {
				if bct.Rev[ver] >= len(bct.bcc.BCC) {
					tmp := last[bct.Rev[ver]]
					last[bct.Rev[ver]] = i
					if tmp != i {
						bct.Tree[bct.Rev[ver]] = append(bct.Tree[bct.Rev[ver]], Edge{bct.Rev[ver], i, e.cost, e.index})
						bct.Tree[i] = append(bct.Tree[i], Edge{i, bct.Rev[ver], e.cost, e.index})
					}
				} else {
					bct.Rev[ver] = i
				}
			}
		}
	}
	bct.Group = make([][]int, ptr)
	for i := range bct.g {
		bct.Group[bct.Rev[i]] = append(bct.Group[bct.Rev[i]], i)
	}
}

// 原图的顶点k对应的圆方树的顶点编号.
func (bct *BlockCutTree) Get(k int) int { return bct.Rev[k] }

// 圆方树的顶点k是否是原图的割点.
func (bct *BlockCutTree) IsCut(k int) bool { return k >= len(bct.bcc.BCC) }

type BiConnectedComponents struct {
	BCC     [][]Edge // 每个边双连通分量中的边
	g       [][]Edge
	lowLink *LowLink
	used    []bool
	tmp     []Edge
}

func NewBiConnectedComponents(g [][]Edge) *BiConnectedComponents {
	return &BiConnectedComponents{
		g:       g,
		lowLink: NewLowLink(g),
	}
}

func (bcc *BiConnectedComponents) Build() {
	bcc.lowLink.Build()
	bcc.used = make([]bool, len(bcc.g))
	for i := 0; i < len(bcc.used); i++ {
		if !bcc.used[i] {
			bcc.dfs(i, -1)
		}
	}
}

func (bcc *BiConnectedComponents) dfs(idx, par int) {
	bcc.used[idx] = true
	beet := false
	for _, next := range bcc.g[idx] {
		if next.to == par {
			b := beet
			beet = true
			if !b {
				continue
			}
		}

		if !bcc.used[next.to] || bcc.lowLink.ord[next.to] < bcc.lowLink.ord[idx] {
			bcc.tmp = append(bcc.tmp, next)
		}

		if !bcc.used[next.to] {
			bcc.dfs(next.to, idx)
			if bcc.lowLink.low[next.to] >= bcc.lowLink.ord[idx] {
				bcc.BCC = append(bcc.BCC, []Edge{})
				for {
					e := bcc.tmp[len(bcc.tmp)-1]
					bcc.BCC[len(bcc.BCC)-1] = append(bcc.BCC[len(bcc.BCC)-1], e)
					bcc.tmp = bcc.tmp[:len(bcc.tmp)-1]
					if e.index == next.index {
						break
					}
				}
			}
		}
	}
}

type LowLink struct {
	Articulation []int  // 関節点
	Bridge       []Edge // 橋
	g            [][]Edge
	ord, low     []int
	used         []bool
}

func NewLowLink(g [][]Edge) *LowLink {
	return &LowLink{g: g}
}

func (ll *LowLink) Build() {
	ll.used = make([]bool, len(ll.g))
	ll.ord = make([]int, len(ll.g))
	ll.low = make([]int, len(ll.g))
	k := 0
	for i := 0; i < len(ll.g); i++ {
		if !ll.used[i] {
			k = ll.dfs(i, k, -1)
		}
	}
}

func (ll *LowLink) dfs(idx, k, par int) int {
	ll.used[idx] = true
	ll.ord[idx] = k
	k++
	ll.low[idx] = ll.ord[idx]
	isArticulation := false
	beet := false
	cnt := 0
	for _, e := range ll.g[idx] {
		if e.to == par {
			tmp := beet
			beet = true
			if !tmp {
				continue
			}
		}
		if !ll.used[e.to] {
			cnt++
			k = ll.dfs(e.to, k, idx)
			ll.low[idx] = min(ll.low[idx], ll.low[e.to])
			if par >= 0 && ll.low[e.to] >= ll.ord[idx] {
				isArticulation = true
			}
			if ll.ord[idx] < ll.low[e.to] {
				ll.Bridge = append(ll.Bridge, e)
			}
		} else {
			ll.low[idx] = min(ll.low[idx], ll.ord[e.to])
		}
	}

	if par == -1 && cnt > 1 {
		isArticulation = true
	}
	if isArticulation {
		ll.Articulation = append(ll.Articulation, idx)
	}
	return k
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
