package main

import (
	"fmt"
	"sort"
)

func main() {
	demo()
}

func demo() {
	dc := NewOffLineDynamicConnectivity(
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
		dc.AddEdge(start, end, i)
	}
	for i, pos := range queries {
		dc.AddQuery(pos, i)
	}

	dc.Run()
}

// 线段树分治.
// 线段树分治是一种处理动态修改和询问的离线算法.
// 通过将某一元素的出现时间段在线段树上保存，我们可以 dfs 遍历整棵线段树，
// 运用可撤销数据结构维护来得到每个时间点的答案.
type OffLineDynamicConnectivity struct {
	add     func(edgeId int)
	undo    func()
	query   func(queryId int)
	edges   []struct{ start, end, id int }
	queries []struct{ time, id int }
	nodes   [][]int
	edgeId  int
	queryId int
}

// dfs 遍历整棵线段树来得到每个时间点的答案.
//
//	add: 添加编号为`edgeId`的边后产生的副作用.
//	undo: 撤销一次`add`操作.
//	query: 响应编号为`queryId`的查询.
func NewOffLineDynamicConnectivity(add func(edgeId int), undo func(), query func(queryId int)) *OffLineDynamicConnectivity {
	return &OffLineDynamicConnectivity{add: add, undo: undo, query: query}
}

// 在时间范围`[startTime, endTime)`内添加一条编号为`edgeId`的边.
func (o *OffLineDynamicConnectivity) AddEdge(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	o.edges = append(o.edges, struct{ start, end, id int }{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`queryId`的查询.
func (o *OffLineDynamicConnectivity) AddQuery(time int, id int) {
	o.queries = append(o.queries, struct{ time, id int }{time, id})
}

func (o *OffLineDynamicConnectivity) Run() {
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
	for _, e := range o.edges {
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
	for _, e := range o.edges {
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

func (o *OffLineDynamicConnectivity) dfs(now int) {
	curNodes := o.nodes[now]
	for _, id := range curNodes {
		if id&1 == 1 {
			o.query(id >> 1)
		} else {
			o.add(id >> 1)
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
