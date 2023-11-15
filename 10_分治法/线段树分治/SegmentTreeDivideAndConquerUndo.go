package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	// demo()
	DynamicGraphVertexAddComponentSum()

}

func demo() {
	dc := NewSegmentTreeDivideAndConquerUndo(
		func(edgeId int) {
			fmt.Println(fmt.Sprintf("add %d", edgeId))
		},
		func() {
			fmt.Println("undo")
		},
		func(queryId int) {
			fmt.Println(fmt.Sprintf("query %d", queryId))
		},
	)

	adds := [][2]int{{0, 1}, {1, 2}, {2, 3}}
	queries := []int{0, 1, 2}
	for i, e := range adds {
		start, end := e[0], e[1]
		dc.AddMutation(start, end, i)
	}
	for i, pos := range queries {
		dc.AddQuery(pos, i)
	}

	dc.Run()
}

const INF int = 1e18

// https://judge.yosupo.jp/problem/dynamic_graph_vertex_add_component_sum
// 0 u v 连接u v (保证u v不连接)
// 1 u v 断开u v  (保证u v连接)
// 2 u x 将u的值加上x
// 3 u 输出u所在连通块的值
func DynamicGraphVertexAddComponentSum() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	sums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &sums[i])
	}
	operations := make([][3]int, 0, q)
	for i := 0; i < q; i++ {
		var op, u, v, add int
		fmt.Fscan(in, &op)
		if op == 0 || op == 1 {
			fmt.Fscan(in, &u, &v)
			operations = append(operations, [3]int{op, u, v})
		} else if op == 2 {
			fmt.Fscan(in, &u, &add)
			operations = append(operations, [3]int{op, u, add})
		} else {
			fmt.Fscan(in, &u)
			operations = append(operations, [3]int{op, u})
		}
	}

	// !solution
	edges := [][2]int{}               // (u, v), u > v
	existEdge := make(map[int][2]int) // (u, v) -> (id, startTime)
	adds := [][2]int{}                // (pos, add)
	queries := []int{}                // pos
	res := []int{}
	uf := NewUnionFindArrayWithUndoAndWeight(sums)
	dc := NewSegmentTreeDivideAndConquerUndo(
		func(mutationId int) {
			if mutationId >= 0 {
				e := edges[mutationId]
				u, v := e[0], e[1]
				uf.Union(u, v)
			} else {
				mutationId = ^mutationId
				a := adds[mutationId]
				pos, add := a[0], a[1]
				uf.SetGroupWeight(pos, uf.GetGroupWeight(pos)+add)
			}
		},
		func() {
			uf.Undo()
		},
		func(queryId int) {
			pos := queries[queryId]
			res[queryId] = uf.GetGroupWeight(pos)
		},
	)

	for time, operation := range operations {
		op := operation[0]
		if op == 0 {
			u, v := operation[1], operation[2]
			if u < v {
				u, v = v, u
			}
			hash := u*n + v
			existEdge[hash] = [2]int{len(edges), time}
			edges = append(edges, [2]int{u, v})
		} else if op == 1 {
			u, v := operation[1], operation[2]
			if u < v {
				u, v = v, u
			}
			hash := u*n + v
			item := existEdge[hash]
			id, startTime := item[0], item[1]
			dc.AddMutation(startTime, time, id)
			delete(existEdge, hash)
		} else if op == 2 {
			pos, add := operation[1], operation[2]
			id := ^len(adds) // 加权值操作的mutationId取反
			dc.AddMutation(time, INF, id)
			adds = append(adds, [2]int{pos, add})
		} else {
			pos := operation[1]
			dc.AddQuery(time, len(queries))
			queries = append(queries, pos)
			res = append(res, 0)
		}
	}

	for _, item := range existEdge {
		id, startTime := item[0], item[1]
		dc.AddMutation(startTime, INF, id)
	}

	dc.Run()

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治undo版.
// 线段树分治是一种处理动态修改和询问的离线算法.
// 通过将某一元素的出现时间段在线段树上保存到`log(n)`个结点中,
// 我们可以 dfs 遍历整棵线段树，运用可撤销数据结构维护来得到每个时间点的答案.
type SegmentTreeDivideAndConquerUndo struct {
	mutate    func(mutationId int)
	undo      func()
	query     func(queryId int)
	mutations []segMutation
	queries   []segQuery
	nodes     [][]int
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	mutate: 添加编号为`mutationId`的变更后产生的副作用.
//	undo: 撤销一次`mutate`操作.
//	query: 响应编号为`queryId`的查询.
//	一共调用**O(nlogn)**次`mutate`和`undo`，**O(q)**次`query`.
func NewSegmentTreeDivideAndConquerUndo(mutate func(mutationId int), undo func(), query func(queryId int)) *SegmentTreeDivideAndConquerUndo {
	return &SegmentTreeDivideAndConquerUndo{mutate: mutate, undo: undo, query: query}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerUndo) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, segMutation{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerUndo) AddQuery(time int, id int) {
	o.queries = append(o.queries, segQuery{time, id})
}

func (o *SegmentTreeDivideAndConquerUndo) Run() {
	if len(o.queries) == 0 {
		return
	}
	times := make([]int, len(o.queries))
	for i := range o.queries {
		times[i] = o.queries[i].time
	}
	sort.Ints(times)
	uniqueInplace(&times)
	usedTimes := make([]bool, len(times)+1)
	usedTimes[0] = true
	for _, e := range o.mutations {
		usedTimes[lowerBound(times, e.start)] = true
		usedTimes[lowerBound(times, e.end)] = true
	}
	for i := 1; i < len(times); i++ {
		if !usedTimes[i] {
			times[i] = times[i-1]
		}
	}
	uniqueInplace(&times)

	n := len(times)
	offset := 1
	for offset < n {
		offset <<= 1
	}
	o.nodes = make([][]int, offset+n)
	for _, e := range o.mutations {
		left := offset + lowerBound(times, e.start)
		right := offset + lowerBound(times, e.end)
		eid := e.id << 1
		for left < right {
			if left&1 == 1 {
				o.nodes[left] = append(o.nodes[left], eid)
				left++
			}
			if right&1 == 1 {
				right--
				o.nodes[right] = append(o.nodes[right], eid)
			}
			left >>= 1
			right >>= 1
		}
	}

	for _, q := range o.queries {
		pos := offset + upperBound(times, q.time) - 1
		o.nodes[pos] = append(o.nodes[pos], (q.id<<1)|1)
	}

	o.dfs(1)
}

func (o *SegmentTreeDivideAndConquerUndo) dfs(now int) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(id >> 1)
		} else {
			o.mutate(id >> 1)
		}
	}
	if now<<1 < len(o.nodes) {
		o.dfs(now << 1)
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now << 1) | 1)
	}
	for i := len(curNodes) - 1; i >= 0; i-- {
		if curNodes[i]&1 == 0 {
			o.undo()
		}
	}
}

func uniqueInplace(sorted *[]int) {
	tmp := *sorted
	slow := 0
	for fast := 0; fast < len(tmp); fast++ {
		if tmp[fast] != tmp[slow] {
			slow++
			tmp[slow] = tmp[fast]
		}
	}
	*sorted = tmp[:slow+1]
}

func lowerBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func upperBound(arr []int, target int) int {
	left, right := 0, len(arr)-1
	for left <= right {
		mid := (left + right) >> 1
		if arr[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

//
//

type S = int

func op(s1, s2 S) S { return s1 + s2 }

func NewUnionFindArrayWithUndoAndWeight(weight []S) *UnionFindArrayWithUndoAndWeight {
	n := len(weight)
	parent, rank, ws := make([]int, n), make([]int, n), make([]S, n)
	for i := 0; i < n; i++ {
		parent[i], rank[i], ws[i] = i, 1, weight[i]
	}
	history := []historyItem{}
	return &UnionFindArrayWithUndoAndWeight{Part: n, rank: rank, parent: parent, weight: ws, history: history}
}

type historyItem struct {
	root, rank int
	weight     S
}

type UnionFindArrayWithUndoAndWeight struct {
	Part    int
	rank    []int
	parent  []int
	weight  []S
	history []historyItem
}

// 将下标为index元素`所在集合`的权值置为value.
func (uf *UnionFindArrayWithUndoAndWeight) SetGroupWeight(index int, value S) {
	index = uf.Find(index)
	uf.history = append(uf.history, historyItem{index, uf.rank[index], uf.weight[index]})
	uf.weight[index] = value
}

// 获取下标为index元素`所在集合`的权值.
func (uf *UnionFindArrayWithUndoAndWeight) GetGroupWeight(index int) S {
	return uf.weight[uf.Find(index)]
}

// 撤销上一次合并(Union)或者修改权值(Set)操作
func (uf *UnionFindArrayWithUndoAndWeight) Undo() {
	if len(uf.history) == 0 {
		return
	}
	last := len(uf.history) - 1
	small := uf.history[last].root
	ps := uf.parent[small]
	uf.weight[ps] = uf.history[last].weight
	uf.rank[ps] = uf.history[last].rank
	if ps != small {
		uf.parent[small] = small
		uf.Part++
	}
	uf.history = uf.history[:last]
}

// 撤销所有操作
func (uf *UnionFindArrayWithUndoAndWeight) Reset() {
	for len(uf.history) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndoAndWeight) Find(x int) int {
	if uf.parent[x] == x {
		return x
	}
	return uf.Find(uf.parent[x])
}

func (uf *UnionFindArrayWithUndoAndWeight) Union(x, y int) bool {
	x, y = uf.Find(x), uf.Find(y)
	if uf.rank[x] < uf.rank[y] {
		x ^= y
		y ^= x
		x ^= y
	}
	uf.history = append(uf.history, historyItem{y, uf.rank[x], uf.weight[x]})
	if x != y {
		uf.parent[y] = x
		uf.rank[x] += uf.rank[y]
		uf.weight[x] = op(uf.weight[x], uf.weight[y])
		uf.Part--
		return true
	}
	return false
}

func (uf *UnionFindArrayWithUndoAndWeight) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindArrayWithUndoAndWeight) Size(x int) int { return uf.rank[uf.Find(x)] }
