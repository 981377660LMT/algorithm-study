package main

import (
	"fmt"
	"sort"
	"time"
)

type Event struct{ pos, kind, count, wage int }
type Pair struct{ value, count int }

const INF int = 1e18

func minimumTotalDistance(robot []int, factory [][]int) int64 {
	n, m := len(robot), len(factory)
	events := make([]Event, 0, n+m)
	factoryCount := 0
	for i := 0; i < n; i++ {
		pos := robot[i]
		events = append(events, Event{pos, 0, 1, 0})
	}
	for i := 0; i < m; i++ {
		pos, count := factory[i][0], factory[i][1]
		events = append(events, Event{pos, 1, count, 0})
		factoryCount += count
	}
	if factoryCount < n {
		return -1
	}
	sort.Slice(events, func(i, j int) bool { return events[i].pos < events[j].pos })

	mousePq := NewBinaryHeap(func(a, b interface{}) int { return a.(Pair).value - b.(Pair).value }, nil)
	holePq := NewBinaryHeap(func(a, b interface{}) int { return a.(Pair).value - b.(Pair).value }, nil)
	holePq.Push(Pair{INF, INF}) // !一开始让所有老鼠都和一个距离它 INF 的洞匹配
	res := 0
	for i := 0; i < len(events); i++ {
		event := events[i]
		if event.kind == 0 { // !老鼠
			top := holePq.Pop().(Pair)
			res += top.value + event.pos
			mousePq.Push(Pair{-event.pos*2 - top.value, 1})
			top.count--
			if top.count > 0 {
				holePq.Push(top)
			}
		} else { // !洞
			rollback := 0
			for mousePq.Len() > 0 && event.count > 0 && mousePq.Peek().(Pair).value+event.pos+event.wage < 0 {
				top := mousePq.Pop().(Pair)
				now := min(top.count, event.count)
				res += now*top.value + event.pos + event.wage
				holePq.Push(Pair{-event.pos*2 - top.value, now})
				rollback += now
				event.count -= now
				top.count -= now
				if top.count > 0 {
					mousePq.Push(top)
				}
			}
			if rollback > 0 {
				mousePq.Push(Pair{-event.pos - event.wage, rollback})
			}
			if event.count > 0 {
				holePq.Push(Pair{-event.pos + event.wage, event.count})
			}
		}
	}

	return int64(res)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func NewBinaryHeap(comparator Comparator, nums []interface{}) *BinaryHeap {
	if nums == nil {
		nums = []interface{}{}
	}
	numsCopy := append([]interface{}{}, nums...)
	heap := &BinaryHeap{comparator: comparator, data: numsCopy}
	heap.heapify()
	return heap
}

// Should return a number:
//    negative , if a < b
//    zero     , if a == b
//    positive , if a > b
type Comparator func(a, b interface{}) int

type BinaryHeap struct {
	data       []interface{}
	comparator Comparator
}

func (h *BinaryHeap) Push(value interface{}) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *BinaryHeap) Pop() (value interface{}) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *BinaryHeap) Peek() (value interface{}) {
	if h.Len() == 0 {
		return
	}

	value = h.data[0]
	return
}

func (h *BinaryHeap) Len() int { return len(h.data) }

func (h *BinaryHeap) heapify() {
	for i := (h.Len() >> 1) - 1; i >= 0; i-- {
		h.pushDown(i)
	}
}

func (h *BinaryHeap) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.comparator(h.data[root], h.data[parent]) < 0; parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *BinaryHeap) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.comparator(h.data[left], h.data[minIndex]) < 0 {
			minIndex = left
		}

		if right < n && h.comparator(h.data[right], h.data[minIndex]) < 0 {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
}

// test
func main() {
	// 1e5
	robots := make([]int, 1e5)
	for i := 0; i < 1e5; i++ {
		robots[i] = i
	}
	// 1e5
	factory := make([][]int, 1e5)
	for i := 0; i < 1e5; i++ {
		factory[i] = []int{i, i}
	}

	time1 := time.Now()
	fmt.Println(minimumTotalDistance(robots, factory))
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))
}
