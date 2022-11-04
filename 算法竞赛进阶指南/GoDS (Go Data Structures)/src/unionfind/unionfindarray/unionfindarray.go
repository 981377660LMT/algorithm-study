package unionfindarray

import (
	"fmt"
	"strings"
)

func demo() {
	uf := NewUnionFindArray(10)
	uf.Union(0, 1)
	uf.Union(1, 2)
	fmt.Println(uf)
}

func NewUnionFindArray(n int) *UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &UnionFindArray{
		Part:   n,
		Rank:   rank,
		size:   n,
		parent: parent,
	}
}

type UnionFindArray struct {
	// 连通分量的个数
	Part int
	// 每个连通分量的大小
	Rank []int

	size   int
	parent []int
}

func (ufa *UnionFindArray) Union(key1, key2 int) bool {
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

func (ufa *UnionFindArray) Find(key int) int {
	for p := ufa.parent[key]; key != p; {
		ufa.parent[key] = ufa.parent[p]
		key = p
	}
	return key
}

func (ufa *UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.size; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}
