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

		removedTimes := make([][]int, q) // 1操作的mul在哪些时间消失
		for i, operation := range operations {
			op := operation[0]
			if op == 2 {
				pos := operation[1]
				removedTimes[pos] = append(removedTimes[pos], i)
			}
		}

		mutations := [][3]int{} // (start,end,posId)
		for eid, operation := range operations {
			op := operation[0]
			if op == 2 {
				continue
			}
			curTimes := removedTimes[eid]
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

		queries := make([]int, q)
		for i := range queries {
			queries[i] = i
		}

		res := make([]int, q)
		initState := &State{value: 1}
		seg := NewSegmentTreeDivideAndConquerCopy(
			initState,
			func(state *State, mutationId int) {
				pos := mutations[mutationId][2]
				mul := operations[pos][1]
				state.value = (state.value * mul) % MOD
				fmt.Println("mutate", state.value, pos, mul)
			},
			func(state *State) *State {
				return &State{state.value}
			},
			func(state *State, queryId int) {
				fmt.Println("query,", state.value, queryId)
				res[queryId] = state.value % MOD
			},
		)
		for id, item := range mutations {
			seg.AddMutation(item[0], item[1], id)
		}
		for id, time := range queries {
			seg.AddQuery(time, id)
		}
		seg.Run()

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

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	 initState: 数据结构的初始状态.
//		mutate: 添加编号为`mutationId`的变更后产生的副作用.
//		copy: 拷贝一份数据结构的副本.
//		query: 响应编号为`queryId`的查询.
//		一共调用 **O(nlogn)** 次`mutate`，**O(n)** 次`copy` 和 **O(q)** 次`query`.
func NewSegmentTreeDivideAndConquerCopy(
	initState *State,
	mutate func(state *State, mutationId int),
	copy func(state *State) *State,
	query func(state *State, queryId int),
) *SegmentTreeDivideAndConquerCopy {
	return &SegmentTreeDivideAndConquerCopy{
		initState: initState,
		mutate:    mutate, copy: copy, query: query,
	}
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

func (o *SegmentTreeDivideAndConquerCopy) Run() {
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
