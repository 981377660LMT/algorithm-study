package main

func maximumSegmentSum(nums []int, removeQueries []int) []int64 {
	n := len(nums)
	curMax := 0
	res := make([]int64, n)
	uf := NewUndoDSU(n)
	visited := make([]bool, n)
	for i := n - 1; i > 0; i-- {
		addIndex := removeQueries[i]
		visited[addIndex] = true
		uf.SetWeight(addIndex, nums[addIndex])
		if addIndex > 0 && visited[addIndex-1] {
			uf.Union(addIndex, addIndex-1)
		}
		if addIndex < n-1 && visited[addIndex+1] {
			uf.Union(addIndex, addIndex+1)
		}
		curMax = max(curMax, uf.GetWeight(uf.Find(addIndex)))
		res[i-1] = int64(curMax)
	}

	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

type S = int

func (*UndoDSU) op(s1, s2 S) S { return s1 + s2 }

func NewUndoDSU(n int) *UndoDSU { return NewUndoDSUWithWeights(make([]S, n)) }
func NewUndoDSUWithWeights(weights []S) *UndoDSU {
	n := len(weights)
	parent, rank, ws := make([]int, n), make([]int, n), make([]S, n)
	for i := 0; i < n; i++ {
		parent[i], rank[i], ws[i] = i, 1, weights[i]
	}
	history := []historyItem{}
	return &UndoDSU{rank: rank, parents: parent, weights: ws, history: history}
}

type historyItem struct {
	a, b int
	c    S
}

type UndoDSU struct {
	rank    []int
	parents []int
	weights []S
	history []historyItem
}

// 将下标为index元素`所在集合`的权值置为value.
func (uf *UndoDSU) SetWeight(index int, value S) {
	index = uf.Find(index)
	uf.history = append(uf.history, historyItem{index, uf.rank[index], uf.weights[index]})
	uf.weights[index] = value
}

// 获取下标为index元素`所在集合`的权值.
func (uf *UndoDSU) GetWeight(index int) S { return uf.weights[uf.Find(index)] }

// 撤销上一次合并(Union)或者修改权值(SetWeight)操作
func (uf *UndoDSU) Undo() {
	uf.weights[uf.parents[uf.history[len(uf.history)-1].a]] = uf.history[len(uf.history)-1].c
	uf.rank[uf.parents[uf.history[len(uf.history)-1].a]] = uf.history[len(uf.history)-1].b
	uf.parents[uf.history[len(uf.history)-1].a] = uf.history[len(uf.history)-1].a
	uf.history = uf.history[:len(uf.history)-1]
}

// 撤销所有操作
func (uf *UndoDSU) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UndoDSU) Find(x int) int {
	if uf.parents[x] == x {
		return x
	}
	return uf.Find(uf.parents[x])
}

func (uf *UndoDSU) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	if uf.rank[x] < uf.rank[y] {
		x, y = y, x
	}
	uf.history = append(uf.history, historyItem{y, uf.rank[x], uf.weights[x]})
	if x != y {
		uf.parents[y] = x
		uf.rank[x] += uf.rank[y]
		uf.weights[x] = uf.op(uf.weights[x], uf.weights[y])
		return true
	}
	return false
}

func (uf *UndoDSU) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }
