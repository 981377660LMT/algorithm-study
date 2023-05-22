// 维护图中桥的个数，支持增加边的操作
// https://ei1333.github.io/library/graph/connected-components/incremental-bridge-connectivity.hpp
// 概要
// 辺の追加クエリのみ存在するとき, 二重辺連結成分を効率的に管理するデータ構造.

// 使い方
// IncrementalBridgeConnectivity(sz): sz 頂点で初期化する.
// Find(k): 頂点 k が属する二重辺連結成分(の代表元)を求める.
// GetBridgeSize(): 現在の橋の個数を返す.
// AddEdge(x, y): 頂点 x と y との間に無向辺を追加する.

package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	// https://judge.yosupo.jp/submission/125538
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	ibc := NewIncrementalBridgeConnectivity(n)
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		ibc.AddEdge(a, b)
	}
	id := make([]int, n)
	for i := 0; i < n; i++ {
		id[i] = -1
	}
	k := 0
	res := make([][]int, n)
	for i := 0; i < n; i++ {
		r := ibc.Find(i)
		if id[r] == -1 {
			id[r] = k
			k++
		}
		res[id[r]] = append(res[id[r]], i)
	}
	fmt.Fprintln(out, k)
	for i := 0; i < k; i++ {
		fmt.Fprint(out, len(res[i]))
		for _, v := range res[i] {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

type IncrementalBridgeConnectivity struct {
	cc, bcc *_UnionFindArray
	bbf     []int
	bridge  int
}

func NewIncrementalBridgeConnectivity(sz int) *IncrementalBridgeConnectivity {
	bbf := make([]int, sz)
	for i := 0; i < sz; i++ {
		bbf[i] = sz
	}
	return &IncrementalBridgeConnectivity{
		cc:  NewUnionFindArray(sz),
		bcc: NewUnionFindArray(sz),
		bbf: bbf,
	}
}

func (ibc *IncrementalBridgeConnectivity) Find(k int) int {
	return ibc.bcc.Find(k)
}

func (ibc *IncrementalBridgeConnectivity) GetBridgeSize() int {
	return ibc.bridge
}

func (ibc *IncrementalBridgeConnectivity) AddEdge(x, y int) {
	x, y = ibc.bcc.Find(x), ibc.bcc.Find(y)
	if ibc.cc.Find(x) == ibc.cc.Find(y) {
		w := ibc.lca(x, y)
		ibc.compress(x, w)
		ibc.compress(y, w)
	} else {
		if ibc.cc.Size(x) > ibc.cc.Size(y) {
			x, y = y, x
		}
		ibc.link(x, y)
		ibc.cc.Union(x, y)
		ibc.bridge++
	}
}

func (ibc *IncrementalBridgeConnectivity) size() int {
	return len(ibc.bbf)
}

func (ibc *IncrementalBridgeConnectivity) par(x int) int {
	if ibc.bbf[x] == ibc.size() {
		return ibc.size()
	}
	return ibc.bcc.Find(ibc.bbf[x])
}

func (ibc *IncrementalBridgeConnectivity) lca(x, y int) int {
	used := make(map[int]struct{})
	for {
		if x != ibc.size() {
			if _, ok := used[x]; ok {
				return x
			}
			used[x] = struct{}{}
			x = ibc.par(x)
		}
		x, y = y, x
	}
}

func (ibc *IncrementalBridgeConnectivity) compress(x, y int) {
	for ibc.bcc.Find(x) != ibc.bcc.Find(y) {
		nxt := ibc.par(x)
		ibc.bbf[x] = ibc.bbf[y]
		ibc.bcc.Union(x, y)
		x = nxt
		ibc.bridge--
	}
}

func (ibc *IncrementalBridgeConnectivity) link(x, y int) {
	v, pre := x, y
	for v != ibc.size() {
		nxt := ibc.par(v)
		ibc.bbf[v] = pre
		pre = v
		v = nxt
	}
}

func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		Rank:   rank,
		size:   n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int
	// 每个连通分量的大小
	Rank []int

	size   int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.Rank[root1] > ufa.Rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.Rank[root2] += ufa.Rank[root1]
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
	for i := 0; i < ufa.size; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
	return ufa.Rank[ufa.Find(key)]
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
