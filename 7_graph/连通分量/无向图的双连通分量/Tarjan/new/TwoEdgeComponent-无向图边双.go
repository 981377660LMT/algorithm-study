package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	Yuki1983()
	// yosupo()
}

func Yuki1983() {
	// https://yukicoder.me/problems/no/1983
	// 给定一个无向图
	// 对每个查询a,b，问从a到b的路径是否存在且唯一
	// !用并查集把所有的桥的端点合并
	// !如果起点终点在一个连通分量内，那么答案就是Yes(只能通过桥往来)
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, q int
	fmt.Fscan(in, &n, &m, &q)
	graph := make([][]Neighbor, n)
	edges := make([][2]int, 0, m)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		u, v = u-1, v-1
		edges = append(edges, [2]int{u, v})
		graph[u] = append(graph[u], Neighbor{to: v, id: i})
		graph[v] = append(graph[v], Neighbor{to: u, id: i})
	}

	_, belong := TwoEdgeComponent(n, m, graph)

	uf := NewUnionFindArray(n)
	for _, e := range edges {
		u, v := e[0], e[1]
		if belong[u] != belong[v] { // (u,v)是桥
			uf.Union(u, v)
		}
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

// https://judge.yosupo.jp/problem/two_edge_connected_components
func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	g := make([][]Neighbor, n)
	for i := 0; i < m; i++ {
		var u, v int
		fmt.Fscan(in, &u, &v)
		g[u] = append(g[u], Neighbor{v, i})
		g[v] = append(g[v], Neighbor{u, i})
	}

	count, belong := TwoEdgeComponent(n, m, g)
	group := make([][]int, count)
	for i := 0; i < n; i++ {
		group[belong[i]] = append(group[belong[i]], i)
	}
	fmt.Fprintln(out, count)
	for _, p := range group {
		fmt.Fprint(out, len(p))
		for _, v := range p {
			fmt.Fprint(out, " ", v)
		}
		fmt.Fprintln(out)
	}
}

type Neighbor = struct{ to, id int }

// 求无向图的边双连通分量.
// 返回值为 (边双连通分量数, 每个点所属的边双连通分量编号).
// !如果某条边(u,v)满足 `belong[u]!=belong[v]`，则称该边为桥.
func TwoEdgeComponent(n, m int, graph [][]Neighbor) (count int, belong []int) {
	path := make([]int, 0, n)
	parent := make([]int, n)
	for i := range parent {
		parent[i] = -2
	}
	dp := make([]int, n)
	belong = make([]int, n)
	used := make([]bool, m)

	var dfs func(int)
	dfs = func(v int) {
		path = append(path, v)
		for _, e := range graph[v] {
			if used[e.id] {
				continue
			}

			used[e.id] = true
			if parent[e.to] == -2 {
				parent[e.to] = v
				dfs(e.to)
			} else {
				dp[v]++
				dp[e.to]--
			}
		}
	}

	for v := 0; v < n; v++ {
		if parent[v] == -2 {
			parent[v] = -1
			dfs(v)
		}
	}
	for i := n - 1; i >= 0; i-- {
		v := path[i]
		if parent[v] != -1 {
			dp[parent[v]] += dp[v]
		}
	}
	for _, v := range path {
		if dp[v] == 0 {
			belong[v] = count
			count++
		} else {
			belong[v] = belong[parent[v]]
		}
	}
	return
}

// 边双缩点成树.
// !各个 e-BCC 的连接成一棵树，每个 e-BCC 为树的一个节点
// !ebcc1 - ebcc2 - ebcc3 - ...
func ToTree(edges [][2]int, count int, belong []int) (tree [][]int) {
	tree = make([][]int, count)
	for _, e := range edges {
		u, v := belong[e[0]], belong[e[1]]
		if u != v { // (u,v)是桥
			tree[u] = append(tree[u], v)
			tree[v] = append(tree[v], u)
		}
	}
	return
}

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
