// https://www.luogu.com.cn/problem/P5443

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
	"strings"
)

// 给定一张 n个点，m条边的无向带权图。
// 每次询问给定一个二元组 (x,y)，从 x 号节点开始出发，只允许通过边权 ≥y 的边。
// 问能够到达的联通块最大的大小。
// 操作1:
// 1 ei wi 将第 ei 条边的边权修改为 wi。
// 2 si li 询问从 si 号节点出发，只允许通过边权 ≥li 的边，能够到达的联通块的最大大小。
// !要求动态修改边权
//
// - 如果问题是静态的，可以使用kruskal重构树在线维护;
// - 如果问题是动态的，发现更新成本大，考虑操作序列分块.
type edge struct{ u, v, w, id int }
type mutation = struct{ eid, weight, mid int }
type query = struct{ start, lower, qid int }

func P5443桥梁(n int, edges []edge, operations [][3]int) (res []int) {
	m, q := len(edges), len(operations)
	res = make([]int, q)
	for i := range res {
		res[i] = -1
	}

	// 操作序列分块.
	block := UseBlock(q, 4*int(math.Sqrt(float64(q))+1)) // 减少分块个数
	blockStart, blockEnd, blockCount := block.blockStart, block.blockEnd, block.blockCount
	for bid := 0; bid < blockCount; bid++ {
		var curMutation []mutation
		var curQuery []query
		isEidMutated := make([]bool, m) // 记录边(id)是否被修改过
		for i := blockStart[bid]; i < blockEnd[bid]; i++ {
			kind := operations[i][0]
			if kind == 1 {
				ei, weight := operations[i][1], operations[i][2]
				curMutation = append(curMutation, mutation{eid: ei, weight: weight, mid: i})
				isEidMutated[ei] = true
			} else {
				start, heavy := operations[i][1], operations[i][2]
				curQuery = append(curQuery, query{start: start, lower: heavy, qid: i})
			}
		}

		// 边按照边权从大到小排序，查询按照汽车重量从大到小排序.
		sort.Slice(edges, func(i, j int) bool { return edges[i].w > edges[j].w })
		sort.Slice(curQuery, func(i, j int) bool { return curQuery[i].lower > curQuery[j].lower })
		eidToEdgeIndex := make([]int, m)
		eidToWeight := make([]int, m)
		for i := range edges {
			eidToEdgeIndex[edges[i].id] = i
			eidToWeight[edges[i].id] = edges[i].w
		}

		uf := NewUnionFindArrayWithUndo(n)
		edgePtr := 0
		for _, query := range curQuery { // 处理块内查询
			for edgePtr < m && edges[edgePtr].w >= query.lower {
				if !isEidMutated[edges[edgePtr].id] {
					uf.Union(edges[edgePtr].u, edges[edgePtr].v) // !合并未被修改的已有边
				}
				edgePtr++
			}

			state := uf.GetState()
			// 处理被修改的边
			for _, mutation := range curMutation {
				eid := mutation.eid
				eidToWeight[eid] = edges[eidToEdgeIndex[eid]].w // !reset
			}
			for _, mutation := range curMutation {
				eid, weight, mid := mutation.eid, mutation.weight, mutation.mid
				if mid < query.qid {
					eidToWeight[eid] = weight // !update
				}
			}
			for _, mutation := range curMutation {
				eid := mutation.eid
				if eidToWeight[eid] >= query.lower {
					uf.Union(edges[eidToEdgeIndex[eid]].u, edges[eidToEdgeIndex[eid]].v) // !合并被修改的边
				}
			}
			res[query.qid] = uf.GetSize(query.start)
			uf.Rollback(state)
		}

		// 更新边权
		for _, mutation := range curMutation {
			eid, weight := mutation.eid, mutation.weight
			edges[eidToEdgeIndex[eid]].w = weight
		}
	}

	removeMinusOne := func(nums []int) []int {
		ptr := 0
		for _, v := range nums {
			if v != -1 {
				nums[ptr] = v
				ptr++
			}
		}
		return nums[:ptr]
	}
	return removeMinusOne(res)
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([]edge, m)
	for i := 0; i < m; i++ {
		var u, v, w int
		fmt.Fscan(in, &u, &v, &w)
		u--
		v--
		edges[i] = edge{u: u, v: v, w: w, id: i}
	}

	var q int
	fmt.Fscan(in, &q)
	ops := make([][3]int, q)
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var ei, weight int
			fmt.Fscan(in, &ei, &weight)
			ei--
			ops[i] = [3]int{kind, ei, weight}
		} else {
			var start, heavy int
			fmt.Fscan(in, &start, &heavy)
			start--
			ops[i] = [3]int{kind, start, heavy}
		}
	}

	res := P5443桥梁(n, edges, ops)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo {
	data := make([]int, n)
	for i := range data {
		data[i] = -1
	}
	return &UnionFindArrayWithUndo{Part: n, n: n, data: data}
}

type historyItem struct{ a, b int }

type UnionFindArrayWithUndo struct {
	Part      int
	n         int
	innerSnap int
	data      []int
	history   []historyItem // (root,data)
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

// 保存并查集当前的状态.
//
//	!Snapshot() 之后可以调用 Rollback(-1) 回滚到这个状态.
func (uf *UnionFindArrayWithUndo) Snapshot() {
	uf.innerSnap = len(uf.history) >> 1
}

// 回滚到指定的状态.
//
//	state 为 -1 表示回滚到上一次 `SnapShot` 时保存的状态.
//	其他值表示回滚到状态id为state时的状态.
func (uf *UnionFindArrayWithUndo) Rollback(state int) bool {
	if state == -1 {
		state = uf.innerSnap
	}
	state <<= 1
	if state < 0 || state > len(uf.history) {
		return false
	}
	for state < len(uf.history) {
		uf.Undo()
	}
	return true
}

// 获取当前并查集的状态id.
//
//	也就是当前合并(Union)被调用的次数.
func (uf *UnionFindArrayWithUndo) GetState() int {
	return len(uf.history) >> 1
}

func (uf *UnionFindArrayWithUndo) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo) Union(x, y int) bool {
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

func (uf *UnionFindArrayWithUndo) Find(x int) int {
	cur := x
	for uf.data[cur] >= 0 {
		cur = uf.data[cur]
	}
	return cur
}
func (ufa *UnionFindArrayWithUndo) SetPart(part int) { ufa.Part = part }

func (uf *UnionFindArrayWithUndo) IsConnected(x, y int) bool { return uf.Find(x) == uf.Find(y) }

func (uf *UnionFindArrayWithUndo) GetSize(x int) int { return -uf.data[uf.Find(x)] }

func (ufa *UnionFindArrayWithUndo) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *UnionFindArrayWithUndo) String() string {
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
