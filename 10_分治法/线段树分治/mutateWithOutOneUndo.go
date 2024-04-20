// https://leetcode.cn/contest/tianchi2022/problems/tRZfIV/

package main

import (
	"fmt"
	"sort"
	"strings"
)

func main() {
	// [[1,0,0,0,0,0,0,0,0],[0,1,0,0,0,0,0,0,0],[0,0,1,0,1,0,1,0,0],[0,0,0,1,0,0,0,0,0],[0,0,1,0,1,0,0,0,0],[0,0,0,0,0,1,0,0,0],[0,0,1,0,0,0,1,0,0],[0,0,0,0,0,0,0,1,0],[0,0,0,0,0,0,0,0,1]]

	fmt.Println(minMalwareSpread([][]int{
		{1, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 1, 0, 0},
		{0, 0, 0, 1, 0, 0, 0, 0, 0},
		{0, 0, 1, 0, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 1, 0, 0, 0},
		{0, 0, 1, 0, 0, 0, 1, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1},
	}, []int{6, 0, 4}))
}

// 238. 除自身以外数组的乘积
// https://leetcode.cn/problems/product-of-array-except-self/
func productExceptSelf(nums []int) []int {
	n := int32(len(nums))
	res := make([]int32, n)
	for i := int32(0); i < n; i++ {
		res[i] = 1
	}

	cur := int32(1)
	history := make([]int32, 0, n)
	MutateWithoutOneUndo(
		0, n,
		func(index int32) {
			history = append(history, cur)
			cur *= int32(nums[index])
		},
		func() {
			cur = history[len(history)-1]
			history = history[:len(history)-1]
		},
		func(index int32) {
			res[index] = cur
		},
	)

	res32 := make([]int, n)
	for i := int32(0); i < n; i++ {
		res32[i] = int(res[i])
	}
	return res32
}

// 928. 尽量减少恶意软件的传播 II
// https://leetcode.cn/problems/minimize-malware-spread-ii/description/
// 分治加点(分治删点).
//
// !1.初始时，将每个非病毒节点与其相邻的非病毒节点合并.
// !2.然后，对每个病毒节点，将其与相邻的非病毒块合并，记录新感染的节点数.
func minMalwareSpread(adjMatrix [][]int, virus []int) int {
	if len(virus) == 0 {
		return 0
	}
	if len(virus) == 1 {
		return virus[0]
	}

	n := int32(len(adjMatrix))
	bad := make([]bool, n)
	for _, v := range virus {
		bad[v] = true
	}

	adjList := make([][]int32, n)
	uf := NewUnionFindArrayWithUndo(n)
	for i := int32(0); i < n; i++ {
		for j := i + 1; j < n; j++ {
			if adjMatrix[i][j] == 1 {
				adjList[i] = append(adjList[i], j)
				adjList[j] = append(adjList[j], i)
				if !bad[i] && !bad[j] {
					uf.Union(i, j)
				}
			}
		}
	}

	minVirusCount, resId := n, int32(-1)
	newVirusCount := int32(0)
	historyStack := [][2]int32{} // 历史栈 (newVirusCount, newBadRoot)
	states := [][2]int32{}       // 数据结构状态 (ufState, historyStackState)
	MutateWithoutOneUndo(
		0, int32(len(virus)),
		func(index int32) {
			cur := int32(virus[index])
			states = append(states, [2]int32{uf.GetState(), int32(len(historyStack))})

			for _, next := range adjList[cur] {
				if bad[next] {
					continue
				}
				nextRoot := uf.Find(next)
				if bad[nextRoot] {
					continue
				}
				nextSize := uf.GetSizeByRoot(nextRoot)
				uf.Union(cur, nextRoot)
				newVirusCount += nextSize
				bad[nextRoot] = true
				historyStack = append(historyStack, [2]int32{newVirusCount, nextRoot})
			}
		},
		func() {
			preState := states[len(states)-1]
			ufState, historyState := preState[0], preState[1]
			states = states[:len(states)-1]
			uf.Rollback(ufState)
			for int32(len(historyStack)) > historyState {
				cur := historyStack[len(historyStack)-1]
				historyStack = historyStack[:len(historyStack)-1]
				newVirusCount = cur[0]
				bad[cur[1]] = false
			}
		},
		func(index int32) {
			curId := int32(virus[index])
			if newVirusCount < minVirusCount || (newVirusCount == minVirusCount && curId < resId) {
				minVirusCount = newVirusCount
				resId = curId
			}
		},
	)

	return int(resId)
}

// 线段树分治的特殊情形.
// 调用 `query` 时，`state` 为对除了 `index` 以外所有点均调用过了 `mutate` 的状态。但不保证调用 `mutate` 的顺序。
// 总计会调用 $O(NlgN)$ 次的 `mutate` 和 `undo`, 以及 $O(N)$ 次的 `query`.
func MutateWithoutOneUndo(
	start, end int32,
	/** 这里的 index 也就是 time. */
	mutate func(index int32),
	undo func(),
	query func(index int32),
) {
	var dfs func(curStart, curEnd int32)
	dfs = func(curStart, curEnd int32) {
		if curEnd == curStart+1 {
			query(curStart)
			return
		}

		mid := (curStart + curEnd) >> 1
		for i := curStart; i < mid; i++ {
			mutate(i)
		}
		dfs(mid, curEnd)
		for i := curStart; i < mid; i++ {
			undo()
		}

		for i := mid; i < curEnd; i++ {
			mutate(i)
		}
		dfs(curStart, mid)
		for i := mid; i < curEnd; i++ {
			undo()
		}
	}

	dfs(start, end)
}

type historyItem struct{ a, b int32 }

type UnionFindArrayWithUndo struct {
	Part    int32
	n       int32
	data    []int32
	history []historyItem // (root,data)
}

func NewUnionFindArrayWithUndo(n int32) *UnionFindArrayWithUndo {
	data := make([]int32, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: n, data: data}
}

// !撤销上一次合并操作，没合并成功也要撤销.
func (uf *UnionFindArrayWithUndo) Undo() bool {
	if len(uf.history) == 0 {
		return false
	}
	small, smallData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	big, bigData := uf.history[len(uf.history)-1].a, uf.history[len(uf.history)-1].b
	uf.history = uf.history[:len(uf.history)-1]
	uf.data[small] = smallData
	uf.data[big] = bigData
	if big != small {
		uf.Part++
	}
	return true
}

// 回滚到指定的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int32) bool {
	state <<= 1
	if state < 0 || state > int32(len(uf.history)) {
		return false
	}
	for state < int32(len(uf.history)) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//
//	也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int32 {
	return int32(len(uf.history) >> 1)
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int32) bool {
	x, y = uf.Find(x), uf.Find(y)
	uf.history = append(uf.history, historyItem{x, uf.data[x]})
	uf.history = append(uf.history, historyItem{y, uf.data[y]})
	if x == y {
		return false
	}
	if uf.data[x] > uf.data[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.data[x] += uf.data[y]
	uf.data[y] = x
	uf.Part--
	return true
}

func (uf *UnionFindArrayWithUndo) Find(x int32) int32 {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}
func (ufa *UnionFindArrayWithUndo) SetPart(part int32) { ufa.Part = part }

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int32) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int32) int32 { return -uf.data[uf.Find(x)] }

func (uf *UnionFindArrayWithUndo) GetSizeByRoot(root int32) int32 { return -uf.data[root] }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int32][]int32 {
	groups := make(map[int32][]int32)
	for i := int32(0); i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) String() string {
	sb := []string{"UnionFindArray:"}
	groups := ufa.GetGroups()
	keys := make([]int32, 0, len(groups))
	for k := range groups {
		keys = append(keys, k)
	}
	sort.Slice(keys, func(i, j int) bool { return keys[i] < keys[j] })
	for _, root := range keys {
		member := groups[root]
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

func mins(nums []int) int {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}
