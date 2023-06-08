package main

import (
	"fmt"
	"sort"
	"strings"
)

func countComponents(n int, edges [][]int) int {
	uf := NewUnionFindArray(n)
	for _, edge := range edges {
		uf.Union(edge[0], edge[1])
	}
	return uf.Part
}

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

// 定向合并.
func (ufa *UnionFindArray) UnionTo(child, parent int) bool {
	root1, root2 := ufa.Find(child), ufa.Find(parent)
	if root1 == root2 {
		return false
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

//
//
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

// 定向合并.
func (ufa *UnionFindMap) UnionTo(child, parent int) bool {
	root1, root2 := ufa.Find(child), ufa.Find(parent)
	if root1 == root2 {
		return false
	}
	ufa.data[root1] += ufa.data[root2]
	ufa.data[root2] = root1
	ufa.Part--
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
