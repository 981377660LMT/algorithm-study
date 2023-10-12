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

// UnionRange 合并`闭区间`[left, right] 的所有元素, 返回合并次数.
func (uf *UnionFindRange) UnionRange(left, right int, f func(big, small int)) int {
	if left >= right {
		return 0
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

// 维护每个分组左右边界的区间并查集.
// 按秩合并.
type UnionFindRange2 struct {
	GroupStart []int
	GroupEnd   []int
	Part       int
	_data      []int
}

func NewUnionFindRange2(n int) *UnionFindRange2 {
	start, end, data := make([]int, n), make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
		start[i] = i
		end[i] = i + 1
		data[i] = -1
	}
	return &UnionFindRange2{GroupStart: start, GroupEnd: end, Part: n, _data: data}
}

func (uf *UnionFindRange2) Find(x int) int {
	if uf._data[x] < 0 {
		return x
	}
	uf._data[x] = uf.Find(uf._data[x])
	return uf._data[x]
}

func (uf *UnionFindRange2) Union(x, y int, f func(big, small int)) bool {
	x = uf.Find(x)
	y = uf.Find(y)
	if x == y {
		return false
	}
	if uf._data[x] > uf._data[y] {
		x, y = y, x
	}
	uf._data[x] += uf._data[y]
	uf._data[y] = x
	uf.Part--
	if f != nil {
		f(x, y)
	}
	uf.GroupStart[x] = min(uf.GroupStart[x], uf.GroupStart[y])
	uf.GroupEnd[x] = max(uf.GroupEnd[x], uf.GroupEnd[y])
	return true
}

// UnionRange 合并`左闭右开`区间 [start, end) 的所有元素, 返回合并次数.
func (uf *UnionFindRange2) UnionRange(start, end int, f func(big, small int)) int {
	start = max(start, 0)
	end = min(end, len(uf._data))
	if start >= end {
		return 0
	}
	res := 0
	cur := 0
	for cur = uf.GroupEnd[uf.Find(start)]; cur < end; cur = uf.GroupEnd[uf.Find(start)] {
		uf.Union(start, cur, f)
		res++
	}
	return res
}

// 每个点所在分组的左右边界[start,end).
func (uf *UnionFindRange2) GetRange(x int) (start, end int) {
	return uf.GroupStart[uf.Find(x)], uf.GroupEnd[uf.Find(x)]
}

func (uf *UnionFindRange2) GetSize(x int) int {
	return -uf._data[uf.Find(x)]
}

func (uf *UnionFindRange2) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := 0; i < len(uf._data); i++ {
		res[uf.Find(i)] = append(res[uf.Find(i)], i)
	}
	return res
}

func (uf *UnionFindRange2) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}
