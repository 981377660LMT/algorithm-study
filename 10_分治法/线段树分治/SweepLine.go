package main

import "sort"

// 2251. 花期内花的数目
// https://leetcode.cn/problems/number-of-flowers-in-full-bloom/submissions/
func fullBloomFlowers(flowers [][]int, people []int) []int {
	res := make([]int, len(people))
	count := 0
	sweepLine := NewSweepLine()

	for i, e := range flowers {
		sweepLine.AddMutation(e[0], e[1]+1, i)
	}
	for i, e := range people {
		sweepLine.AddQuery(e, i)
	}

	sweepLine.Run(func(mutationId int) {
		count++
	},
		func(mutationId int) {
			count--
		},
		func(queryId int) {
			res[queryId] = count
		})
	return res
}

// 2747. 统计没有收到请求的服务器数目
// https://leetcode.cn/problems/count-zero-request-servers/
func countServers(n int, logs [][]int, x int, queries []int) []int {
	res := make([]int, len(queries))
	count := n
	serverCounter := make([]int, n+1)
	sweepLine := NewSweepLine()

	for i, e := range logs {
		start := e[1]
		sweepLine.AddMutation(start, start+x+1, i)
	}
	for i, q := range queries {
		sweepLine.AddQuery(q, i)
	}

	sweepLine.Run(func(mutationId int) {
		serverId := logs[mutationId][0]
		serverCounter[serverId]++
		if serverCounter[serverId] == 1 {
			count--
		}
	},
		func(mutationId int) {
			serverId := logs[mutationId][0]
			serverCounter[serverId]--
			if serverCounter[serverId] == 0 {
				count++
			}
		},
		func(queryId int) {
			res[queryId] = count
		})
	return res
}

// 给定一个时间轴（或者设想一个），
// 有若干个操作(可交换,commutative)在时间 [start,end) 中起作用。
// 询问某一个时间某个值是什么.
// 如果修改操作可删除，那么可以使用'扫描线'来解决.
type SweepLine struct {
	mutate    func(mutationId int)
	remove    func(mutationId int)
	query     func(queryId int)
	mutations []struct{ start, end, id int }
	queries   []struct{ time, id int }
	nodes     [][]int
}

func NewSweepLine() *SweepLine {
	return &SweepLine{}
}

// 在时间范围`[startTime, endTime)`内添加一个编号为`id`的变更.
func (s *SweepLine) AddMutation(startTime, endTime int, id int) {
	if startTime >= endTime {
		return
	}
	s.mutations = append(s.mutations, struct{ start, end, id int }{startTime, endTime, id})
}

// 在时间`time`时添加一个编号为`id`的查询.
func (s *SweepLine) AddQuery(time int, id int) {
	s.queries = append(s.queries, struct{ time, id int }{time, id})
}

// 使用扫描线得到每个时间点的答案.
//
//	mutate: 添加编号为`mutationId`的变更.
//	remove: 删除编号为`mutationId`的变更.
//	query: 响应编号为`queryId`的查询.
//	一共调用 **O(n)** 次`mutate`、`remove` 和 **O(q)** 次`query`.
func (s *SweepLine) Run(
	mutate func(mutationId int),
	remove func(mutationId int),
	query func(queryId int),
) {
	if len(s.queries) == 0 {
		return
	}
	s.mutate, s.remove, s.query = mutate, remove, query
	times := make([]int, len(s.queries))
	for i := range s.queries {
		times[i] = s.queries[i].time
	}
	sort.Ints(times)
	dedup(&times)
	usedTimes := make([]bool, len(times)+1)
	usedTimes[0] = true
	for _, e := range s.mutations {
		usedTimes[lowerBound(times, e.start)] = true
		usedTimes[lowerBound(times, e.end)] = true
	}
	for i := 1; i < len(times); i++ {
		if !usedTimes[i] {
			times[i] = times[i-1]
		}
	}
	dedup(&times)

	s.nodes = make([][]int, len(times)+1)
	for _, e := range s.mutations {
		left := lowerBound(times, e.start)
		right := lowerBound(times, e.end)
		s.nodes[left] = append(s.nodes[left], e.id*2+1)
		s.nodes[right] = append(s.nodes[right], -e.id*2-1)
	}

	for _, q := range s.queries {
		pos := upperBound(times, q.time) - 1
		s.nodes[pos] = append(s.nodes[pos], q.id*2)
	}

	s.doSweep()
}

func (s *SweepLine) doSweep() {
	for _, events := range s.nodes {
		for _, id := range events {
			if id >= 0 {
				if id&1 == 1 {
					s.mutate(id >> 1)
				} else {
					s.query(id >> 1)
				}
			} else {
				s.remove(-id >> 1)
			}
		}
	}
}

func dedup(sorted *[]int) {
	if len(*sorted) == 0 {
		return
	}
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
