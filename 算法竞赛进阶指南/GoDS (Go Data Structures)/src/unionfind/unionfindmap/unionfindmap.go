package unionfindmap

import (
	"fmt"
	"strings"
)

func demo() {
	uf := NewUnionFindMap(true)
	uf.Union(0, 1)
	uf.Union(1, 2)
	uf.Union(1, 999)
	fmt.Println(uf)
}

type U = int
type UnionFindMap struct {
	// 连通分量的个数
	Part int

	rank map[U]int
	// 当key不存在于并查集中时，是否自动添加
	autoAdd bool
	parent  map[U]U
}

func NewUnionFindMap(autoAdd bool) *UnionFindMap {
	return &UnionFindMap{
		Part:    0,
		rank:    make(map[U]int),
		autoAdd: autoAdd,
		parent:  make(map[U]U),
	}
}

func (ufm *UnionFindMap) Union(key1, key2 U) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}

	absent := !ufm.contains(key1) || !ufm.contains(key2)
	if absent {
		return false
	}

	if ufm.rank[root1] > ufm.rank[root2] {
		root1, root2 = root2, root1
	}

	ufm.parent[root1] = root2
	ufm.rank[root2] += ufm.rank[root1]
	ufm.Part--
	return true
}

func (ufm *UnionFindMap) Find(key U) U {
	if !ufm.contains(key) {
		if ufm.autoAdd {
			ufm.Add(key)
		}
		return key
	}

	for ufm.parent[key] != key {
		ufm.parent[key] = ufm.parent[ufm.parent[key]]
		key = ufm.parent[key]
	}
	return key
}

func (ufm *UnionFindMap) Add(key U) bool {
	if ufm.contains(key) {
		return false
	}

	ufm.parent[key] = key
	ufm.rank[key] = 1
	ufm.Part++
	return true
}

func (ufm *UnionFindMap) IsConnected(key1, key2 U) bool {
	absent := !ufm.contains(key1) || !ufm.contains(key2)
	if absent && !ufm.autoAdd {
		return false
	}

	return ufm.Find(key1) == ufm.Find(key2)
}

func (ufm *UnionFindMap) GetGroups() map[U][]U {
	groups := make(map[U][]U)
	for k := range ufm.parent {
		root := ufm.Find(k)
		groups[root] = append(groups[root], k)
	}
	return groups
}

func (ufm *UnionFindMap) String() string {
	sb := []string{"UnionFindMap:"}
	for root, member := range ufm.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufm.Part))
	return strings.Join(sb, "\n")
}

func (ufm *UnionFindMap) contains(key U) bool {
	_, ok := ufm.parent[key]
	return ok
}
