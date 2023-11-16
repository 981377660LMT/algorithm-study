// P4219 [BJOI2014] 大融合
// https://www.luogu.com.cn/problem/P4219

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const INF int = 1e9

// 0 x y: 在结点 x 和 y 之间连接，保证之前没有边 x-y.
// 1 x y: 给定一条已经存在的边，求有多少条简单路径包含边 x-y.
// 保证在任意时刻，图的形态都是一棵森林。
//
// !负载操作就是删掉后这条边两边连通块个数乘起来，那么直接可撤销并查集维护一下即可
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	operations := make([][3]int, q)
	for i := range operations {
		var kind string
		var x, y int
		fmt.Fscan(in, &kind, &x, &y)
		x--
		y--
		if kind == "A" {
			operations[i] = [3]int{0, x, y}
		} else {
			operations[i] = [3]int{1, x, y}
		}
	}

	queries := []int{} // edgeHash
	Q := NewAddRemoveQuery(true)
	seg := NewSegmentTreeDivideAndConquerUndo()
	for i, op := range operations {
		kind, x, y := op[0], op[1], op[2]
		if x < y {
			x, y = y, x
		}
		hash := x*n + y
		if kind == 0 {
			Q.Add(i, hash)
		} else {
			Q.Remove(i, hash)
			Q.Add(i+1, hash) // 下个时刻加回来
			queries = append(queries, hash)
			seg.AddQuery(i, len(queries)-1)
		}
	}
	mutations := Q.GetEvents(INF) // start,end,edgeHash
	for _, m := range mutations {
		seg.AddMutation(m.start, m.end, m.value)
	}

	res := make([]int, len(queries))
	uf := NewUnionFindArrayWithUndo(n)
	seg.Run(
		func(edgeHash int) {
			u, v := edgeHash/n, edgeHash%n
			uf.Union(u, v)
		},
		func() {
			uf.Undo()
		},
		func(queryId int) {
			hash := queries[queryId]
			u, v := hash/n, hash%n
			res[queryId] = uf.GetSize(u) * uf.GetSize(v)
		},
	)

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

type historyItem struct{ a, b int }
type UnionFindArrayWithUndo struct {
	Part      int
	n         int
	innerSnap int
	data      []int
	history   []historyItem // (root,data)
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo {
	data := make([]int, n)
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

type V = int
type Event = struct {
	start int
	end   int
	value V
}

// 将时间轴上单点的 add 和 remove 转化为区间上的 add.
// !不能加入相同的元素，删除时元素必须要存在。
// 如果 add 和 remove 按照时间顺序严格单增，那么可以使用 monotone = true 来加速。
type AddRemoveQuery struct {
	mp       map[V]int
	events   []Event
	adds     map[V][]int
	removes  map[V][]int
	monotone bool
}

func NewAddRemoveQuery(monotone bool) *AddRemoveQuery {
	return &AddRemoveQuery{
		mp:       map[V]int{},
		events:   []Event{},
		adds:     map[V][]int{},
		removes:  map[V][]int{},
		monotone: monotone,
	}
}

func (adq *AddRemoveQuery) Add(time int, value V) {
	if adq.monotone {
		adq.addMonotone(time, value)
	} else {
		adq.adds[value] = append(adq.adds[value], time)
	}
}

func (adq *AddRemoveQuery) Remove(time int, value V) {
	if adq.monotone {
		adq.removeMonotone(time, value)
	} else {
		adq.removes[value] = append(adq.removes[value], time)
	}
}

// lastTime: 所有变更都结束的时间.例如INF.
func (adq *AddRemoveQuery) GetEvents(lastTime int) []Event {
	if adq.monotone {
		return adq.getMonotone(lastTime)
	}
	res := []Event{}
	for value, addTimes := range adq.adds {
		removeTimes := []int{}
		if cand, ok := adq.removes[value]; ok {
			removeTimes = cand
			delete(adq.removes, value)
		}
		if len(removeTimes) < len(addTimes) {
			removeTimes = append(removeTimes, lastTime)
		}
		sort.Ints(addTimes)
		sort.Ints(removeTimes)
		for i, t := range addTimes {
			if t < removeTimes[i] {
				res = append(res, Event{t, removeTimes[i], value})
			}
		}
	}
	return res
}

func (adq *AddRemoveQuery) addMonotone(time int, value V) {
	if _, ok := adq.mp[value]; ok {
		panic("can't add a value that already exists")
	}
	adq.mp[value] = time
}

func (adq *AddRemoveQuery) removeMonotone(time int, value V) {
	if startTime, ok := adq.mp[value]; !ok {
		panic("can't remove a value that doesn't exist")
	} else {
		delete(adq.mp, value)
		if startTime != time {
			adq.events = append(adq.events, Event{startTime, time, value})
		}
	}
}

func (adq *AddRemoveQuery) getMonotone(lastTime int) []Event {
	for value, startTime := range adq.mp {
		if startTime == lastTime {
			continue
		}
		adq.events = append(adq.events, Event{startTime, lastTime, value})
	}
	return adq.events
}
