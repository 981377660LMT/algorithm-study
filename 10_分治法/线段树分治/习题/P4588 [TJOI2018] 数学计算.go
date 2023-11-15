package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e9

// P4588 [TJOI2018] 数学计算
// https://www.luogu.com.cn/problem/P4588
// 初识时x=1，执行q次操作：
// 1 mul: 将x乘以mul，并输出 x 模 MOD;
// 2 pos: 将x除以第pos次操作所乘的数，并输出 x 模 MOD;
//
// !难以删除 => logn代价将删除变为撤销。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	solve := func() {
		var q, MOD int
		fmt.Fscan(in, &q, &MOD)

		operations := make([][2]int, q)
		for i := range operations {
			var op int
			fmt.Fscan(in, &op)
			if op == 1 {
				var mul int
				fmt.Fscan(in, &mul)
				operations[i] = [2]int{1, mul}
			} else {
				var pos int
				fmt.Fscan(in, &pos)
				pos--
				operations[i] = [2]int{2, pos}
			}
		}

		Q := NewAddRemoveQuery(true)
		for i, item := range operations {
			op := item[0]
			if op == 1 {
				Q.Add(i, i)
			} else {
				pos := item[1]
				Q.Remove(i, pos)
			}
		}

		mutations := Q.Work(q) // start, end, pos
		queries := make([]int, q)
		for i := range queries {
			queries[i] = i
		}

		seg := NewSegmentTreeDivideAndConquerCopy()
		for id, item := range mutations {
			seg.AddMutation(item.start, item.end, id)
		}
		for id, time := range queries {
			seg.AddQuery(time, id)
		}

		res := make([]int, q)
		initState := &State{value: 1}
		seg.Run(
			initState,
			func(state *State, mutationId int) {
				pos := mutations[mutationId].value
				mul := operations[pos][1]
				state.value = (state.value * mul) % MOD
			},
			func(state *State) *State {
				return &State{state.value}
			},
			func(state *State, queryId int) {
				res[queryId] = state.value % MOD
			},
		)

		for _, v := range res {
			fmt.Fprintln(out, v)
		}
	}

	var T int
	fmt.Fscan(in, &T)
	for i := 0; i < T; i++ {
		solve()
	}
}

type State = struct{ value int }

type segMutation struct{ start, end, id int }
type segQuery struct{ time, id int }

// 线段树分治copy版.
// 如果修改操作难以撤销，可以在每个节点处保存一份副本.
// !调用O(n)次拷贝注意不要超出内存.
type SegmentTreeDivideAndConquerCopy struct {
	initState *State
	mutate    func(state *State, mutationId int)
	copy      func(state *State) *State
	query     func(state *State, queryId int)
	mutations []segMutation
	queries   []segQuery
	nodes     [][]int
}

func NewSegmentTreeDivideAndConquerCopy() *SegmentTreeDivideAndConquerCopy {
	return &SegmentTreeDivideAndConquerCopy{}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (o *SegmentTreeDivideAndConquerCopy) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.mutations = append(o.mutations, segMutation{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (o *SegmentTreeDivideAndConquerCopy) AddQuery(time int, id int) {
	o.queries = append(o.queries, segQuery{time, id})
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	 initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func (o *SegmentTreeDivideAndConquerCopy) Run(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) {
	if len(o.queries) == 0 {
		return
	}
	o.initState = initState
	o.mutate, o.copy, o.query = mutate, copy, query
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

	o.dfs(1, o.initState)
}

func (o *SegmentTreeDivideAndConquerCopy) dfs(now int, state *State) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(state, id>>1)
		} else {
			o.mutate(state, id>>1)
		}
	}
	if now<<1 < len(o.nodes) {
		o.dfs(now<<1, o.copy(state))
	}
	if (now<<1)|1 < len(o.nodes) {
		o.dfs((now<<1)|1, o.copy(state))
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
