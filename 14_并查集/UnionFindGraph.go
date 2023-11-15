// 无向图中：
// 联通分量数(part) = 树的个数(treeCount) + 环的个数
// 边的个数: 总点数 - 树的个数
// 树的性质: 联通分量中的点数 = 边数 + 1
// 环的性质: 联通分量中的点数 = 边数

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {

}

// 可撤销并查集，维护连通分量为树的联通分量个数.
type UnionFindGraphArray struct {
	Part      int // 联通分量数
	TreeCount int // 联通分量为树的联通分量个数(孤立点也算树)
	n         int
	data      []int
	edge      []int
	history   [][3]int // (root,data,edge)
}

func NewUnionFindGraphArray(n int) *UnionFindGraphArray {
	edge := make([]int, n)
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindGraphArray{
		Part:      n,
		TreeCount: n,
		n:         n,
		data:      data,
		edge:      edge,
	}
}

// 添加边对(u,v).
func (uf *UnionFindGraphArray) Union(u, v int) bool {
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
func (uf *UnionFindGraphArray) Find(u int) int {
	cur := u
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindGraphArray) Undo() bool {
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

// !从每条边中恰好选一个点, 最多能选出多少个不同的点.
//
//	对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
//	因此答案为`总点数-树的个数`.
func (uf *UnionFindGraphArray) Solve() int {
	return uf.n - uf.TreeCount
}

// u所在的联通分量是否为树.
func (uf *UnionFindGraphArray) IsTree(u int) bool {
	root := uf.Find(u)
	vertex := uf.GetSize(root)
	return vertex == uf.edge[root]+1
}

func (uf *UnionFindGraphArray) IsCycle(u int) bool {
	root := uf.Find(u)
	vertex := uf.GetSize(root)
	return vertex == uf.edge[root]
}

// 联通分量为树的联通分量个数(孤立点也算树).
func (uf *UnionFindGraphArray) CountTree() int {
	return uf.TreeCount
}

// 联通分量为环的联通分量个数(孤立点不算环).
func (uf *UnionFindGraphArray) CountCycle() int {
	return uf.Part - uf.TreeCount
}

func (uf *UnionFindGraphArray) CountEdge() int {
	return uf.n - uf.TreeCount
}

// x所在的连通分量的点数.
func (uf *UnionFindGraphArray) GetSize(x int) int { return -uf.data[uf.Find(x)] }

// x所在的连通分量的边数.
func (uf *UnionFindGraphArray) GetEdge(x int) int { return uf.edge[uf.Find(x)] }

func (ufa *UnionFindGraphArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindGraphArray) String() string {
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

// 可撤销并查集，维护连通分量为树的联通分量个数.
type UnionFindGraphMap struct {
	Part      int // 联通分量数
	TreeCount int // 联通分量为树的联通分量个数(孤立点也算树)
	data      map[int]int
	edge      map[int]int
	history   [][3]int // (root,data,edge)
}

func NewUnionFindGraphMap() *UnionFindGraphMap {
	return &UnionFindGraphMap{
		data: make(map[int]int),
		edge: make(map[int]int),
	}
}

// 添加边对(u,v).
func (uf *UnionFindGraphMap) Union(u, v int) bool {
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
func (uf *UnionFindGraphMap) Find(u int) int {
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

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindGraphMap) Undo() bool {
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

// !从每条边中恰好选一个点, 最多能选出多少个不同的点.
//
//	对每个大小为m的连通块,树的贡献为m-1,环的贡献为m.
//	因此答案为`总点数-树的个数`.
func (s *UnionFindGraphMap) Solve() int {
	return len(s.data) - s.TreeCount
}

// u所在的联通分量是否为树.
func (s *UnionFindGraphMap) IsTree(u int) bool {
	root := s.Find(u)
	vertex := s.GetSize(root)
	return vertex == s.edge[root]+1
}

func (uf *UnionFindGraphMap) IsCycle(u int) bool {
	root := uf.Find(u)
	vertex := uf.GetSize(root)
	return vertex == uf.edge[root]
}

func (s *UnionFindGraphMap) CountTree() int {
	return s.TreeCount
}

// 联通分量为环的联通分量个数(孤立点不算环).
func (uf *UnionFindGraphMap) CountCycle() int {
	return uf.Part - uf.TreeCount
}

func (uf *UnionFindGraphMap) CountEdge() int {
	return len(uf.data) - uf.TreeCount
}

func (uf *UnionFindGraphMap) GetSize(x int) int { return -uf.data[uf.Find(x)] }
func (uf *UnionFindGraphMap) GetEdge(x int) int { return uf.edge[uf.Find(x)] }

func (ufa *UnionFindGraphMap) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for k := range ufa.data {
		root := ufa.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}
func (uf *UnionFindGraphMap) Add(key int) bool {
	if _, ok := uf.data[key]; ok {
		return false
	}
	uf.data[key] = -1
	uf.edge[key] = 0
	uf.Part++
	uf.TreeCount++
	return true
}
func (ufa *UnionFindGraphMap) String() string {
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
