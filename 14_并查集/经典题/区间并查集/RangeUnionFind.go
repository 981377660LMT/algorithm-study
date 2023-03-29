// 区间并查集 RangeUnionFind/UnionFindRange

package main

// https://leetcode.cn/problems/amount-of-new-area-painted-each-day/
func amountPainted(paint [][]int) []int {
	uf := NewUnionFindRange(1e5 + 10)
	res := make([]int, len(paint))
	for i, v := range paint {
		left, right := v[0], v[1]
		res[i] = uf.UnionRange(left, right, func(big, small int) {})
	}
	return res
}

type UnionFindRange struct {
	Part   int
	n      int
	parent []int
	rank   []int
}

func NewUnionFindRange(n int) *UnionFindRange {
	uf := &UnionFindRange{
		Part:   n,
		n:      n,
		parent: make([]int, n),
		rank:   make([]int, n),
	}
	for i := 0; i < n; i++ {
		uf.parent[i] = i
		uf.rank[i] = 1
	}
	return uf
}

func (uf *UnionFindRange) Find(x int) int {
	for x != uf.parent[x] {
		uf.parent[x] = uf.parent[uf.parent[x]]
		x = uf.parent[x]
	}
	return x
}

// Union 后, 大的编号的组会指向小的编号的组.
func (uf *UnionFindRange) Union(x, y int, f func(big, small int)) bool {
	if x < y {
		x, y = y, x
	}
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		return false
	}
	uf.parent[rootX] = rootY
	uf.rank[rootY] += uf.rank[rootX]
	uf.Part--
	f(rootY, rootX)
	return true
}

// UnionRange 合并区间 [left, right] 的所有元素, 返回合并次数.
func (uf *UnionFindRange) UnionRange(left, right int, f func(big, small int)) int {
	if left > right {
		left, right = right, left
	}
	leftRoot := uf.Find(left)
	rightRoot := uf.Find(right)
	unionCount := 0
	for rightRoot != leftRoot {
		unionCount++
		uf.Union(rightRoot, rightRoot-1, f)
		rightRoot = uf.Find(rightRoot - 1)
	}
	return unionCount
}

func (uf *UnionFindRange) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindRange) GetSize(x int) int {
	return uf.rank[uf.Find(x)]
}

func (uf *UnionFindRange) GetGroups() map[int][]int {
	group := make(map[int][]int)
	for i := 0; i < uf.n; i++ {
		group[uf.Find(i)] = append(group[uf.Find(i)], i)
	}
	return group
}
