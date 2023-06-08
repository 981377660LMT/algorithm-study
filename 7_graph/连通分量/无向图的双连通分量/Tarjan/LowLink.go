// 橋や関節点などを効率的に求める際に有効なアルゴリズム.
// グラフをDFSして各頂点 idx について, ord[idx] := DFS で頂点に訪れた順番,
// low[idx] := 頂点 idxからDFS木の葉方向の辺を 0回以上,
// 後退辺を 1回以下通って到達可能な頂点の ord の最小値 を求める.

// build(): LowLink を構築する.
// !構築後, Articulation には関節点, Bridge には橋が格納される.
// 非連結でも多重辺を含んでいてもOK.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func criticalConnections(n int, connections [][]int) [][]int {
	g := make([][]Edge, n)
	for i := 0; i < len(connections); i++ {
		u, v := connections[i][0], connections[i][1]
		g[u] = append(g[u], Edge{from: u, to: v})
		g[v] = append(g[v], Edge{from: v, to: u})
	}
	lowLink := NewLowLink(g)
	lowLink.Build()
	res := [][]int{}
	for i := 0; i < len(lowLink.Bridge); i++ {
		res = append(res, []int{lowLink.Bridge[i].from, lowLink.Bridge[i].to})
	}
	return res
}

func main() {
	yuki1983()
}

func yuki1983() {
	// https://yukicoder.me/problems/no/1983
	// 给定一个无向图
	// 对每个查询a,b，问从a到b的路径是否存在且唯一
	// !用并查集把所有的桥的端点合并
	// !如果起点终点在一个连通分量内，那么答案就是Yes(只能走桥)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	graph := make([][]Edge, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		graph[u] = append(graph[u], Edge{u, v})
		graph[v] = append(graph[v], Edge{v, u})
	}

	lowlink := NewLowLink(graph)
	lowlink.Build()

	uf := NewUnionFindArray(n)
	for _, e := range lowlink.Bridge {
		uf.Union(e.from, e.to)
	}

	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		if uf.IsConnected(a, b) {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
	}
}

type Edge = struct{ from, to int }
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

//
//
func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
