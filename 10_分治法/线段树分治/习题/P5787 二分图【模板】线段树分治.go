// https://www.luogu.com.cn/problem/P5787
// 给定n个节点的图，在k个时间内有m条边会出现后消失。问第i时刻是否是二分图。
// 二分图判定使用可撤销的扩展域并查集维护。

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)

	mutations := make([][4]int, m) // (u,v,start,end)
	for i := 0; i < m; i++ {
		var u, v, start, end int
		fmt.Fscan(in, &u, &v, &start, &end)
		u--
		v--
		mutations[i] = [4]int{u, v, start, end}
	}

	queries := make([]int, k)
	for i := range queries {
		queries[i] = i
	}

	S := NewSegmentTreeDivideAndConquerUndo()
	for id := range mutations {
		start, end := mutations[id][2], mutations[id][3]
		S.AddMutation(start, end, id)
	}
	for id, time := range queries {
		S.AddQuery(time, id)
	}

	res := make([]bool, k)
	uf := NewUnionFindWithDistAndUndo(n) // 带权并查集
	history := []bool{true}
	S.Run(
		func(mutationId int) {
			u, v := mutations[mutationId][0], mutations[mutationId][1]
			uf.SnapShot()
			if !history[len(history)-1] {
				history = append(history, false)
				return
			}
			root1, dist1 := uf.Find(u)
			root2, dist2 := uf.Find(v)
			if root1 == root2 {
				history = append(history, (dist1-dist2+1)&1 != 1)
			} else {
				uf.Union(u, v, 1)
				history = append(history, true)
			}
		},
		func() {
			uf.Rollback(-1)
			history = history[:len(history)-1]
		},
		func(queryId int) {
			res[queryId] = history[len(history)-1]
		},
	)

	for _, v := range res {
		if v {
			fmt.Fprintln(out, "Yes")
		} else {
			fmt.Fprintln(out, "No")
		}
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

func NewSegmentTreeDivideAndConquerUndo() *SegmentTreeDivideAndConquerUndo {
	return &SegmentTreeDivideAndConquerUndo{}
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

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	mutate: 添加编号为`mutationId`的变更后产生的副作用.
//	undo: 撤销一次`mutate`操作.
//	query: 响应编号为`queryId`的查询.
//	一共调用**O(nlogn)**次`mutate`和`undo`，**O(q)**次`query`.
func (o *SegmentTreeDivideAndConquerUndo) Run(mutate func(mutationId int), undo func(), query func(queryId int)) {
	if len(o.queries) == 0 {
		return
	}
	o.mutate, o.undo, o.query = mutate, undo, query
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

type T = int

func e() T        { return 0 }
func op(x, y T) T { return x + y }
func inv(x T) T   { return -x }

// 维护到每个组根节点距离的可撤销并查集.
// 用于维护环的权值，树上的距离等.
type UnionFindWithDistAndUndo struct {
	data      *RollbackArray
	snapShots []int
}

func NewUnionFindWithDistAndUndo(n int) *UnionFindWithDistAndUndo {
	return &UnionFindWithDistAndUndo{
		data: NewRollbackArray(n, func(index int) arrayItem { return arrayItem{parent: -1, dist: e()} }),
	}
}

// 将当前快照加入栈顶.
func (uf *UnionFindWithDistAndUndo) SnapShot() int {
	res := uf.data.GetTime()
	uf.snapShots = append(uf.snapShots, res)
	return res
}

func (uf *UnionFindWithDistAndUndo) GetTime() int {
	return uf.data.GetTime()
}

// time=-1表示回滚到栈顶(上一次)快照的时间，并删除该快照.
func (uf *UnionFindWithDistAndUndo) Rollback(time int) {
	if time != -1 {
		uf.data.Rollback(time)
	} else {
		if len(uf.snapShots) == 0 {
			return
		}
		time = uf.snapShots[len(uf.snapShots)-1]
		uf.snapShots = uf.snapShots[:len(uf.snapShots)-1]
		uf.data.Rollback(time)
	}
}

// distToRoot(parent) + dist = distToRoot(child).
func (uf *UnionFindWithDistAndUndo) Union(parent int, child int, dist T) bool {
	v1, x1 := uf.Find(parent)
	v2, x2 := uf.Find(child)
	if v1 == v2 {
		return dist == op(x2, inv(x1))
	}
	s1, s2 := -uf.data.Get(v1).parent, -uf.data.Get(v2).parent
	if s1 < s2 {
		s1, s2 = s2, s1
		v1, v2 = v2, v1
		x1, x2 = x2, x1
		dist = inv(dist)
	}
	// v1 <- v2
	dist = op(x1, dist)
	dist = op(dist, inv(x2))
	uf.data.Set(v2, arrayItem{parent: v1, dist: dist})
	uf.data.Set(v1, arrayItem{parent: -(s1 + s2), dist: e()})
	return true
}

// 返回v所在组的根节点和到v到根节点的距离.
func (uf *UnionFindWithDistAndUndo) Find(v int) (root int, distToRoot T) {
	root, distToRoot = v, e()
	for {
		item := uf.data.Get(root)
		if item.parent < 0 {
			break
		}
		distToRoot = op(distToRoot, item.dist)
		root = item.parent
	}
	return
}

// Dist(x, y) = DistToRoot(x) - DistToRoot(y).
// 如果x和y不在同一个集合,抛出错误.
func (uf *UnionFindWithDistAndUndo) Dist(x int, y int) T {
	vx, dx := uf.Find(x)
	vy, dy := uf.Find(y)
	if vx != vy {
		panic("x and y are not in the same set")
	}
	return op(dx, inv(dy))
}

func (uf *UnionFindWithDistAndUndo) DistToRoot(x int) T {
	_, dx := uf.Find(x)
	return dx
}

func (uf *UnionFindWithDistAndUndo) GetSize(x int) int {
	root, _ := uf.Find(x)
	return -uf.data.Get(root).parent
}

func (uf *UnionFindWithDistAndUndo) GetGroups() map[int][]int {
	res := make(map[int][]int)
	for i := 0; i < uf.data.Len(); i++ {
		root, _ := uf.Find(i)
		res[root] = append(res[root], i)
	}
	return res
}

type arrayItem struct {
	parent int
	dist   T
}

type historyItem struct {
	index int
	value arrayItem
}

type RollbackArray struct {
	n       int
	data    []arrayItem
	history []historyItem
}

func NewRollbackArray(n int, f func(index int) arrayItem) *RollbackArray {
	data := make([]arrayItem, n)
	for i := 0; i < n; i++ {
		data[i] = f(i)
	}
	return &RollbackArray{
		n:    n,
		data: data,
	}
}

func (r *RollbackArray) GetTime() int {
	return len(r.history)
}

func (r *RollbackArray) Rollback(time int) {
	for len(r.history) > time {
		pair := r.history[len(r.history)-1]
		r.history = r.history[:len(r.history)-1]
		r.data[pair.index] = pair.value
	}
}

func (r *RollbackArray) Undo() bool {
	if len(r.history) == 0 {
		return false
	}
	pair := r.history[len(r.history)-1]
	r.history = r.history[:len(r.history)-1]
	r.data[pair.index] = pair.value
	return true
}

func (r *RollbackArray) Get(index int) arrayItem {
	return r.data[index]
}

func (r *RollbackArray) Set(index int, value arrayItem) {
	r.history = append(r.history, historyItem{index: index, value: r.data[index]})
	r.data[index] = value
}

func (r *RollbackArray) GetAll() []arrayItem {
	return append(r.data[:0:0], r.data...)
}

func (r *RollbackArray) Len() int {
	return r.n
}
