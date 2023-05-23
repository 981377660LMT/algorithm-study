package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	// BallCollector()
	ReversibleCards()
}

// https://atcoder.jp/contests/abc302/tasks/abc302_h
// 给定一棵树，每个点有两个值。
// 对于v=2,3,...,n，问从点1到点 v的最短路径途径的每个点中，
// 各选一个数，其不同数的个数的最大值。
func BallCollector() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	pairs := make([][2]int, n)
	for i := range pairs {
		var a, b int
		fmt.Fscan(in, &a, &b)
		pairs[i] = [2]int{a, b}
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	uf := NewSelectOneFromEachPairMap()
	res := make([]int, n)

	var dfs func(int, int)
	dfs = func(cur, pre int) {
		a, b := pairs[cur][0], pairs[cur][1]
		uf.Union(a, b)
		res[cur] = uf.Solve()
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur)
			}
		}
		uf.Undo()
	}

	dfs(0, -1)
	for i := 1; i < n; i++ {
		fmt.Fprint(out, res[i], " ")
	}
}

// https://atcoder.jp/contests/arc111/tasks/arc111_b
func ReversibleCards() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	pairs := make([][2]int, n)
	for i := range pairs {
		var a, b int
		fmt.Fscan(in, &a, &b)
		pairs[i] = [2]int{a, b}
	}

	uf := NewSelectOneFromEachPairMap()
	for _, p := range pairs {
		uf.Union(p[0], p[1])
	}
	fmt.Fprintln(out, uf.Solve())
}

// 可撤销并查集，维护连通分量为树的联通分量个数.
type SelectOneFromEachPairMap struct {
	Part      int // 联通分量数
	TreeCount int // 联通分量为树的联通分量个数(孤立点也算树)
	data      map[int]int
	edge      map[int]int
	history   [][3]int // (root,data,edge)
}

func NewSelectOneFromEachPairMap() *SelectOneFromEachPairMap {
	return &SelectOneFromEachPairMap{
		data: make(map[int]int),
		edge: make(map[int]int),
	}
}

// !从每条边中恰好选一个点, 最多能选出多少个不同的点.
//  对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
//  因此答案为`总点数-树的个数`.
func (s *SelectOneFromEachPairMap) Solve() int {
	return len(s.data) - s.TreeCount
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *SelectOneFromEachPairMap) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}

	small, smallData, smallEdge := uf.history[len(uf.history)-1][0], uf.history[len(uf.history)-1][1], uf.history[len(uf.history)-1][2]
	uf.history = uf.history[:len(uf.history)-1]
	big, bigData, bigEdge := uf.history[len(uf.history)-1][0], uf.history[len(uf.history)-1][1], uf.history[len(uf.history)-1][2]
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[big], uf.data[small] = bigData, smallData
	uf.edge[big], uf.edge[small] = bigEdge, smallEdge

	// update part/treeCount
	if big == small {
		if uf.IsTree(big) {
			uf.TreeCount++
		}
	} else {
		uf.Part++
		if uf.IsTree(big) || uf.IsTree(small) {
			uf.TreeCount++
		}
	}

	return true
}

// u所在的联通分量是否为树.
func (s *SelectOneFromEachPairMap) IsTree(u int) bool {
	root := s.Find(u)
	vertex := s.GetSize(root)
	return vertex == s.edge[root]+1
}

func (s *SelectOneFromEachPairMap) CountTree() int {
	return s.TreeCount
}

// 添加边对(u,v).
func (uf *SelectOneFromEachPairMap) Union(u, v int) bool {
	u, v = uf.Find(u), uf.Find(v)
	uf.history = append(uf.history, [3]int{u, uf.data[u], uf.edge[u]}) // big
	uf.history = append(uf.history, [3]int{v, uf.data[v], uf.edge[v]}) // small
	if u == v {
		if uf.IsTree(u) {
			uf.TreeCount--
		}
		uf.edge[u]++
		return false
	}

	if uf.data[u] > uf.data[v] {
		u ^= v
		v ^= u
		u ^= v
	}

	// 四种情况的简写
	if uf.IsTree(u) || uf.IsTree(v) {
		uf.TreeCount--
	}

	uf.data[u] += uf.data[v]
	uf.data[v] = u
	uf.edge[u] += uf.edge[v] + 1
	uf.Part--
	return true
}

// 不能路径压缩.
func (uf *SelectOneFromEachPairMap) Find(u int) int {
	if _, ok := uf.data[u]; !ok {
		uf.Add(u)
		return u
	}
	cur := u
	for {
		v, ok := uf.data[cur]
		if !ok || v < 0 {
			return cur
		}
		cur = v
	}
}

func (uf *SelectOneFromEachPairMap) GetSize(x int) int { return -uf.data[uf.Find(x)] }
func (uf *SelectOneFromEachPairMap) GetEdge(x int) int { return uf.edge[uf.Find(x)] }

func (ufa *SelectOneFromEachPairMap) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for k := range ufa.data {
		root := ufa.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}

func (uf *SelectOneFromEachPairMap) Add(key int) bool {
	if _, ok := uf.data[key]; ok {
		return false
	}
	uf.data[key] = -1
	uf.edge[key] = 0
	uf.Part++
	uf.TreeCount++
	return true
}

func (ufa *SelectOneFromEachPairMap) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
