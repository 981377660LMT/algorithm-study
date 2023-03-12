// https://leetcode-cn.com/problems/t3fKg1/solution/10xing-jie-jue-zhan-dou-by-foxtail-ke2e/
// LCP 32. 批量处理任务
// 2 <= tasks.length <= 10^5
// 0 <= tasks[i][0] <= tasks[i][1] <= 10^9
// 某实验室计算机待处理任务以 [start,end,period] 格式记于二维数组 tasks，
// 表示完成该任务的时间范围为起始时间 start 至结束时间 end 之间，需要计算机投入 period 的时长
// !处于开机状态的计算机可同时处理任意多个任务，请返回电脑最少开机多久，可处理完所有任务。

package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func processTasks(tasks [][]int) int {
	allNums := make(map[int]struct{})
	for _, task := range tasks {
		allNums[task[0]-1] = struct{}{}
		allNums[task[1]] = struct{}{}
	}
	nums := make([]int, 0, len(allNums))
	for num := range allNums {
		nums = append(nums, num)
	}
	sort.Ints(nums)

	n := len(nums)
	mp := make(map[int]int, n)
	for i, num := range nums {
		mp[num] = i
	}
	D := NewDualShortestPath(n+10, true)
	for _, task := range tasks {
		u, v := mp[task[0]-1], mp[task[1]] // v - u >= period
		D.AddEdge(u, v, -task[2])
	}
	for i := 1; i < n; i++ {
		D.AddEdge(i-1, i, 0)                 // Si>=Si-1
		D.AddEdge(i, i-1, nums[i]-nums[i-1]) // Si-Si-1<=allNums[i]-allNums[i-1]
	}
	res, ok := D.Run()
	if !ok {
		return -1
	}
	return res[n-1]
}

const INF int = 1e18

type DualShortestPath struct {
	n   int
	g   [][][2]int
	min bool
}

func NewDualShortestPath(n int, min bool) *DualShortestPath {
	return &DualShortestPath{
		n:   n,
		g:   make([][][2]int, n),
		min: min,
	}
}

// f(j) <= f(i) + w
func (d *DualShortestPath) AddEdge(i, j, w int) {
	if d.min {
		d.g[i] = append(d.g[i], [2]int{j, w})
	} else {
		d.g[j] = append(d.g[j], [2]int{i, w})
	}
}

// 求 `f(i) - f(0)` 的最小值/最大值, 并检测是否有负环/正环
func (d *DualShortestPath) Run() (dist []int, ok bool) {
	if d.min {
		return d.spfaMin()
	}
	return d.spfaMax()
}

func (d *DualShortestPath) spfaMin() (dist []int, ok bool) {
	dist = make([]int, d.n)
	queue := NewDeque2(d.n)
	count := make([]int, d.n)
	inQueue := make([]bool, d.n)
	for i := 0; i < d.n; i++ {
		queue.Append(i)
		inQueue[i] = true
		count[i] = 1
	}
	for queue.Size() > 0 {
		cur := queue.PopLeft()
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue.AppendLeft(next)
				}
			}
		}
	}

	for i := 0; i < d.n; i++ {
		dist[i] = -dist[i]
	}
	ok = true
	return
}

func (d *DualShortestPath) spfaMax() (dist []int, ok bool) {
	dist = make([]int, d.n)
	inQueue := make([]bool, d.n)
	count := make([]int, d.n)
	for i := 0; i < d.n; i++ {
		dist[i] = INF
	}

	queue := []int{0}
	dist[0] = 0
	inQueue[0] = true
	count[0] = 1
	for len(queue) > 0 {
		cur := queue[0]
		queue = queue[1:]
		inQueue[cur] = false
		for _, e := range d.g[cur] {
			next, weight := e[0], e[1]
			cand := dist[cur] + weight
			if cand < dist[next] {
				dist[next] = cand
				if !inQueue[next] {
					count[next]++
					if count[next] >= d.n+1 {
						return nil, false
					}
					inQueue[next] = true
					queue = append(queue, next)
				}
			}
		}
	}

	ok = true
	return
}

type D = int
type Deque struct{ l, r []D }

func NewDeque2(cap int) *Deque { return &Deque{make([]D, 0, 1+cap/2), make([]D, 0, 1+cap/2)} }

func (q Deque) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q Deque) Size() int {
	return len(q.l) + len(q.r)
}

func (q *Deque) AppendLeft(v D) {
	q.l = append(q.l, v)
}

func (q *Deque) Append(v D) {
	q.r = append(q.r, v)
}

func (q *Deque) PopLeft() (v D) {
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return
}

func (q *Deque) Pop() (v D) {
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return
}

func (q Deque) Front() D {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q Deque) Back() D {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q Deque) At(i int) D {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	var m int
	fmt.Fscan(in, &m)
	tasks := make([][]int, m)
	for i := 0; i < m; i++ {
		tasks[i] = make([]int, 3)
		fmt.Fscan(in, &tasks[i][0], &tasks[i][1], &tasks[i][2])
	}

	res := processTasks(tasks)
	fmt.Fprintln(out, res)
}
