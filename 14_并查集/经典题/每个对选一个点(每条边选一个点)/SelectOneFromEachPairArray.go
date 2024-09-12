package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

func main() {
	BallCollector()

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
		pairs[i] = [2]int{id(a), id(b)}
	}

	tree := make([][]int, n)
	for i := 0; i < n-1; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		tree[a] = append(tree[a], b)
		tree[b] = append(tree[b], a)
	}

	uf := NewSelectOneFromEachPairArray(len(_pool))
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
		pairs[i] = [2]int{id(a), id(b)}
	}

	uf := NewSelectOneFromEachPairArray(len(_pool))
	for _, p := range pairs {
		uf.Union(p[0], p[1])
	}
	fmt.Fprintln(out, uf.Solve())
}

var _pool = make(map[interface{}]int)

func id(o interface{}) int {
	if v, ok := _pool[o]; ok {
		return v
	}
	v := len(_pool)
	_pool[o] = v
	return v
}

// 可撤销并查集，维护连通分量为树的联通分量个数.
type SelectOneFromEachPairArray struct {
	Part      int // 联通分量数
	TreeCount int // 联通分量为树的联通分量个数(孤立点也算树)
	n         int
	data      []int
	edge      []int
	history   [][3]int // (root,data,edge)
}

// 指定`点权最大值n`建立并查集.
func NewSelectOneFromEachPairArray(n int) *SelectOneFromEachPairArray {
	edge := make([]int, n)
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &SelectOneFromEachPairArray{
		Part:      n,
		TreeCount: n,
		n:         n,
		data:      data,
		edge:      edge,
	}
}

// !从每条边中恰好选一个点, 最多能选出多少个不同的点.
//
//	对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
//	因此答案为`总点数-树的个数`.
func (s *SelectOneFromEachPairArray) Solve() int {
	return s.n - s.TreeCount
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *SelectOneFromEachPairArray) Undo() bool {
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
func (s *SelectOneFromEachPairArray) IsTree(u int) bool {
	root := s.Find(u)
	vertex := s.GetSize(root)
	return vertex == s.edge[root]+1
}

func (s *SelectOneFromEachPairArray) CountTree() int {
	return s.TreeCount
}

// 添加边对(u,v).
func (uf *SelectOneFromEachPairArray) Union(u, v int) bool {
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
func (uf *SelectOneFromEachPairArray) Find(u int) int {
	cur := u
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

func (uf *SelectOneFromEachPairArray) GetSize(x int) int { return -uf.data[uf.Find(x)] }
func (uf *SelectOneFromEachPairArray) GetEdge(x int) int { return uf.edge[uf.Find(x)] }

func (ufa *SelectOneFromEachPairArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *SelectOneFromEachPairArray) String() string {
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
