package main

import (
	"fmt"
	"sort"
	"strings"
)

// 给你一个 n 个节点的带权无向图，节点编号为 0 到 n - 1 。

// 给你一个整数 n 和一个数组 edges ，其中 edges[i] = [ui, vi, wi] 表示节点 ui 和 vi 之间有一条权值为 wi 的无向边。

// 在图中，一趟旅途包含一系列节点和边。旅途开始和结束点都是图中的节点，且图中存在连接旅途中相邻节点的边。注意，一趟旅途可能访问同一条边或者同一个节点多次。

// 如果旅途开始于节点 u ，结束于节点 v ，我们定义这一趟旅途的 代价 是经过的边权按位与 AND 的结果。换句话说，如果经过的边对应的边权为 w0, w1, w2, ..., wk ，那么代价为w0 & w1 & w2 & ... & wk ，其中 & 表示按位与 AND 操作。

// 给你一个二维数组 query ，其中 query[i] = [si, ti] 。对于每一个查询，你需要找出从节点开始 si ，在节点 ti 处结束的旅途的最小代价。如果不存在这样的旅途，答案为 -1 。

// 返回数组 answer ，其中 answer[i] 表示对于查询 i 的 最小 旅途代价。

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	n    int
	data []int
}

func NewUnionFindArray(n int) *UnionFindArray {
	data := make([]int, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArray{
		Part: n,
		n:    n,
		data: data,
	}
}

// 按秩合并.
func (ufa *UnionFindArray) Union(key1, key2 int) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	return true
}
func (ufa *UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.data[root1] > ufa.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufa *UnionFindArray) Find(key int) int {
	if ufa.data[key] < 0 {
		return key
	}
	ufa.data[key] = ufa.Find(ufa.data[key])
	return ufa.data[key]
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetSize(key int) int {
	return -ufa.data[ufa.Find(key)]
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
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

var _pool = make(map[interface{}]int)

func id(o interface{}) int {
	if v, ok := _pool[o]; ok {
		return v
	}
	v := len(_pool)
	_pool[o] = v
	return v
}

type UnionFindMap struct {
	Part int
	data map[int]int
}

func NewUnionFindMap() *UnionFindMap {
	return &UnionFindMap{
		data: make(map[int]int),
	}
}

func (ufm *UnionFindMap) Union(key1, key2 int) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufm.data[root1] > ufm.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufm.data[root1] += ufm.data[root2]
	ufm.data[root2] = root1
	ufm.Part--
	return true
}
func (ufm *UnionFindMap) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufm.data[root1] > ufm.data[root2] {
		root1 ^= root2
		root2 ^= root1
		root1 ^= root2
	}
	ufm.data[root1] += ufm.data[root2]
	ufm.data[root2] = root1
	ufm.Part--
	if cb != nil {
		cb(root1, root2)
	}
	return true
}

func (ufm *UnionFindMap) Find(key int) int {
	if _, ok := ufm.data[key]; !ok {
		ufm.Add(key)
		return key
	}
	if ufm.data[key] < 0 {
		return key
	}
	ufm.data[key] = ufm.Find(ufm.data[key])
	return ufm.data[key]
}

func (ufm *UnionFindMap) IsConnected(key1, key2 int) bool {
	return ufm.Find(key1) == ufm.Find(key2)
}

func (ufm *UnionFindMap) GetSize(key int) int {
	return -ufm.data[ufm.Find(key)]
}

func (ufm *UnionFindMap) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for k := range ufm.data {
		root := ufm.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}

func (ufm *UnionFindMap) Has(key int) bool {
	_, ok := ufm.data[key]
	return ok
}

func (ufm *UnionFindMap) Add(key int) bool {
	if _, ok := ufm.data[key]; ok {
		return false
	}
	ufm.data[key] = -1
	ufm.Part++
	return true
}

func (ufm *UnionFindMap) String() string {
	sb := []string{"UnionFindMap:"}
	groups := ufm.GetGroups()
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
	sb = append(sb, fmt.Sprintf("Part: %d", ufm.Part))
	return strings.Join(sb, "\n")
}

func minimumCost(n int, edges [][]int, query [][]int) []int {
	q := len(query)
	res := make([]int, q)
	uf := NewUnionFindArray(n)
	groupAnd := make([]int, n)
	for i := 0; i < n; i++ {
		groupAnd[i] = -1
	}
	edgeWeights := make(map[[2]int][]int)
	for _, edge := range edges {
		u, v, w := edge[0], edge[1], edge[2]
		if u > v {
			u, v = v, u
		}
		pair := [2]int{u, v}
		edgeWeights[pair] = append(edgeWeights[pair], w)
	}
	for _, edge := range edges {
		u, v := edge[0], edge[1]
		if u > v {
			u, v = v, u
		}
		root1, root2 := uf.Find(u), uf.Find(v)
		if root1 == root2 {
			weights := edgeWeights[[2]int{u, v}]
			for _, w := range weights {
				groupAnd[root1] &= w
			}
		} else {
			uf.UnionWithCallback(u, v, func(big, small int) {
				groupAnd[big] &= groupAnd[small]
				weights := edgeWeights[[2]int{u, v}]
				for _, w := range weights {
					groupAnd[big] &= w
				}
			})
		}

	}
	for i, query := range query {
		si, ti := query[0], query[1]
		if si == ti {
			res[i] = 0
			continue
		}
		root1, root2 := uf.Find(si), uf.Find(ti)
		if root1 != root2 {
			res[i] = -1
			continue
		}
		res[i] = groupAnd[root1]
	}
	return res
}

func main() {
	// n = 3, edges = [[0,2,7],[0,1,15],[1,2,6],[1,2,1]], query = [[1,2]]
	fmt.Println(minimumCost(3, [][]int{{0, 2, 7}, {0, 1, 15}, {1, 2, 6}, {1, 2, 1}}, [][]int{{1, 2}})) // [6]

}
