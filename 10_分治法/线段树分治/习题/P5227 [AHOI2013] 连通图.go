package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strings"
)

const INF int = 1e9

// P5227 [AHOI2013] 连通图
// https://www.luogu.com.cn/problem/P5227
// 给定一个无向连通图和若干个小集合，每个小集合包含一些边(不超过4条)，
// !对于每个集合，你需要确定将集合中的边删掉后改图是否保持联通。
// !注意不等价于"不考虑这个集合".
// 集合间的询问相互独立
//
// !每条边会在若干个时间区间内出现。预处理每条边删除的时间.
func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	edges := make([][2]int, m)
	for i := range edges {
		fmt.Fscan(in, &edges[i][0], &edges[i][1])
		edges[i][0]--
		edges[i][1]--
	}

	var k int
	fmt.Fscan(in, &k)
	edgeGroups := make([][]int, k)
	for i := range edgeGroups {
		var g int
		fmt.Fscan(in, &g)
		edgeGroups[i] = make([]int, g)
		for j := range edgeGroups[i] {
			fmt.Fscan(in, &edgeGroups[i][j])
			edgeGroups[i][j]--
		}
	}

	removedTimes := make([][]int, m) // 每条边在哪些时间消失
	for i, g := range edgeGroups {
		for _, e := range g {
			removedTimes[e] = append(removedTimes[e], i)
		}
	}

	mutations := [][3]int{} // (start,end,edgeId)
	for eid, curTimes := range removedTimes {
		if len(curTimes) == 0 {
			mutations = append(mutations, [3]int{-INF, INF, eid})
			continue
		}
		startTime := -INF
		for _, endTime := range curTimes {
			mutations = append(mutations, [3]int{startTime, endTime, eid})
			startTime = endTime + 1
		}
		mutations = append(mutations, [3]int{curTimes[len(curTimes)-1] + 1, INF, eid})
	}

	queries := make([]int, k)
	for i := range queries {
		queries[i] = i
	}

	seg := NewSegmentTreeDivideAndConquerUndo()
	for id, item := range mutations {
		seg.AddMutation(item[0], item[1], id)
	}
	for id, time := range queries {
		seg.AddQuery(time, id)
	}

	res := make([]bool, k)
	uf := NewUnionFindArrayWithUndo(n)
	seg.Run(
		func(index int) {
			item := mutations[index]
			edge := item[2]
			uf.Union(edges[edge][0], edges[edge][1])
		},
		func() {
			uf.Undo()
		},
		func(index int) {
			res[index] = uf.Part == 1
		},
	)

	for _, b := range res {
		if b {
			fmt.Fprintln(out, "Connected")
		} else {
			fmt.Fprintln(out, "Disconnected")
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

type historyItem struct{ small, big, smallRank int }

type UnionFindArrayWithUndo2 struct {
	Part      int
	_n        int
	_parent   []int
	_rank     []int
	_optStack []historyItem
}

func NewUnionFindArrayWithUndo(n int) *UnionFindArrayWithUndo2 {
	parent := make([]int, n)
	rank := make([]int, n)
	for i := 0; i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}
	return &UnionFindArrayWithUndo2{
		_n:      n,
		_parent: parent,
		_rank:   rank,
		Part:    n,
	}
}

func (uf *UnionFindArrayWithUndo2) Find(x int) int {
	for uf._parent[x] != x {
		x = uf._parent[x]
	}
	return x
}

func (uf *UnionFindArrayWithUndo2) Union(x, y int) bool {
	rootX := uf.Find(x)
	rootY := uf.Find(y)
	if rootX == rootY {
		uf._optStack = append(uf._optStack, historyItem{-1, -1, -1})
		return false
	}
	if uf._rank[rootX] > uf._rank[rootY] {
		rootX, rootY = rootY, rootX
	}
	uf._parent[rootX] = rootY
	uf._rank[rootY] += uf._rank[rootX]
	uf.Part--
	uf._optStack = append(uf._optStack, historyItem{rootX, rootY, uf._rank[rootX]})
	return true
}

func (uf *UnionFindArrayWithUndo2) Undo() {
	if len(uf._optStack) == 0 {
		return
	}
	opt := uf._optStack[len(uf._optStack)-1]
	uf._optStack = uf._optStack[:len(uf._optStack)-1]
	if opt.small == -1 {
		return
	}
	uf._parent[opt.small] = opt.small
	uf._rank[opt.big] -= opt.smallRank
	uf.Part++
}

func (uf *UnionFindArrayWithUndo2) Reset() {
	for len(uf._optStack) > 0 {
		uf.Undo()
	}
}

func (uf *UnionFindArrayWithUndo2) IsConnected(x, y int) bool {
	return uf.Find(x) == uf.Find(y)
}

func (uf *UnionFindArrayWithUndo2) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < uf._n; i++ {
		root := uf.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (uf *UnionFindArrayWithUndo2) GetPart() int {
	return uf.Part
}

func (ufa *UnionFindArrayWithUndo2) String() string {
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
func (adq *AddRemoveQuery) Work(lastTime int) []Event {
	if adq.monotone {
		return adq.workMonotone(lastTime)
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

func (adq *AddRemoveQuery) workMonotone(lastTime int) []Event {
	for value, startTime := range adq.mp {
		if startTime == lastTime {
			continue
		}
		adq.events = append(adq.events, Event{startTime, lastTime, value})
	}
	return adq.events
}
