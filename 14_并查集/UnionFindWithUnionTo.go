package main

import "fmt"

func main() {
	uf := NewUnionFindArrayWithUnionTo(10)
	uf.UnionTo(1, 2)
	fmt.Println(uf.GetGroups())
}

type UnionFindArrayWithUnionTo struct {
	Part   int
	n      int
	parent []int32
	rank   []int32
}

func NewUnionFindArrayWithUnionTo(n int) *UnionFindArrayWithUnionTo {
	parent := make([]int32, n)
	rank := make([]int32, n)
	for i := 0; i < n; i++ {
		parent[i] = int32(i)
		rank[i] = 1
	}
	return &UnionFindArrayWithUnionTo{Part: n, n: n, parent: parent, rank: rank}
}

func (u *UnionFindArrayWithUnionTo) Find(x int) int {
	x32 := int32(x)
	for u.parent[x32] != x32 {
		u.parent[x32] = u.parent[u.parent[x32]]
		x32 = u.parent[x32]
	}
	return int(x32)
}

// 按秩合并.
func (u *UnionFindArrayWithUnionTo) Union(x, y int, f func(big, small int)) bool {
	rootX, rootY := u.Find(x), u.Find(y)
	if rootX == rootY {
		return false
	}
	if u.rank[rootX] > u.rank[rootY] {
		rootX, rootY = rootY, rootX
	}
	u.parent[rootX] = int32(rootY)
	u.rank[rootY] += u.rank[rootX]
	u.Part--
	if f != nil {
		f(rootY, rootX)
	}
	return true
}

// 定向合并.
func (u *UnionFindArrayWithUnionTo) UnionTo(child, parent int) bool {
	rootX, rootY := u.Find(child), u.Find(parent)
	if rootX == rootY {
		return false
	}
	u.parent[rootX] = int32(rootY)
	u.rank[rootY] += u.rank[rootX]
	u.Part--
	return true
}

func (u *UnionFindArrayWithUnionTo) GetSize(x int) int {
	return int(u.rank[u.Find(x)])
}

func (u *UnionFindArrayWithUnionTo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for key := 0; key < u.n; key++ {
		root := u.Find(key)
		groups[root] = append(groups[root], key)
	}
	return groups
}

type UnionFindMapWithUnionTo struct {
	Part   int
	parent map[int32]int32
	rank   map[int32]int32
}

func NewUnionFindMapWithUnionTo() *UnionFindMapWithUnionTo {
	return &UnionFindMapWithUnionTo{Part: 0, parent: make(map[int32]int32), rank: make(map[int32]int32)}
}

func (u *UnionFindMapWithUnionTo) Find(key int) int {
	key32 := int32(key)
	if _, ok := u.parent[key32]; !ok {
		u.add(key32)
		return key
	}
	for u.parent[key32] != key32 {
		u.parent[key32] = u.parent[u.parent[key32]]
		key32 = u.parent[key32]
	}
	return int(key32)
}

// 按秩合并.
func (u *UnionFindMapWithUnionTo) Union(key1, key2 int, f func(big, small int)) bool {
	root1, root2 := int32(u.Find(key1)), int32(u.Find(key2))
	if root1 == root2 {
		return false
	}
	if u.rank[root1] > u.rank[root2] {
		root1, root2 = root2, root1
	}
	u.parent[root1] = root2
	u.rank[root2] += u.rank[root1]
	u.Part--
	if f != nil {
		f(int(root2), int(root1))
	}
	return true
}

// 定向合并.
func (u *UnionFindMapWithUnionTo) UnionTo(child, parent int) bool {
	root1, root2 := int32(u.Find(child)), int32(u.Find(parent))
	if root1 == root2 {
		return false
	}
	u.parent[root1] = root2
	u.rank[root2] += u.rank[root1]
	u.Part--
	return true
}

func (u *UnionFindMapWithUnionTo) GetSize(key int) int {
	return int(u.rank[int32(u.Find(key))])
}

func (u *UnionFindMapWithUnionTo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for key := range u.parent {
		root := u.Find(int(key))
		groups[root] = append(groups[root], int(key))
	}
	return groups
}

func (u *UnionFindMapWithUnionTo) add(key32 int32) {
	u.parent[key32] = key32
	u.rank[key32] = 1
	u.Part++
}
