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

type UnionFindMap struct {
	// 连通分量的个数
	Part int
	// 每个连通分量的大小
	Rank map[interface{}]int

	// 当key不存在于并查集中时，是否自动添加
	autoAdd bool
	parent  map[interface{}]interface{}
}

func NewUnionFindMap(autoAdd bool) *UnionFindMap {
	return &UnionFindMap{
		Part:    0,
		Rank:    make(map[interface{}]int),
		autoAdd: autoAdd,
		parent:  make(map[interface{}]interface{}),
	}
}

func (ufm *UnionFindMap) Union(key1, key2 interface{}) bool {
	root1, root2 := ufm.Find(key1), ufm.Find(key2)
	if root1 == root2 {
		return false
	}

	absent := !ufm.contains(key1) || !ufm.contains(key2)
	if absent {
		return false
	}

	if ufm.Rank[root1] > ufm.Rank[root2] {
		root1, root2 = root2, root1
	}

	ufm.parent[root1] = root2
	ufm.Rank[root2] += ufm.Rank[root1]
	ufm.Part--
	return true
}

func (ufm *UnionFindMap) Find(key interface{}) interface{} {
	if !ufm.contains(key) {
		if ufm.autoAdd {
			ufm.Add(key)
		}
		return key
	}

	for p := ufm.parent[key]; key != p; {
		ufm.parent[key] = ufm.parent[p]
		key = p
	}
	return key
}

func (ufm *UnionFindMap) Add(key interface{}) bool {
	if ufm.contains(key) {
		return false
	}

	ufm.parent[key] = key
	ufm.Rank[key] = 1
	ufm.Part++
	return true
}

func (ufm *UnionFindMap) IsConnected(key1, key2 interface{}) bool {
	absent := !ufm.contains(key1) || !ufm.contains(key2)
	if absent && !ufm.autoAdd {
		return false
	}

	return ufm.Find(key1) == ufm.Find(key2)
}

func (ufm *UnionFindMap) GetGroups() map[interface{}][]interface{} {
	groups := make(map[interface{}][]interface{})
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

func (ufm *UnionFindMap) contains(key interface{}) bool {
	_, ok := ufm.parent[key]
	return ok
}
