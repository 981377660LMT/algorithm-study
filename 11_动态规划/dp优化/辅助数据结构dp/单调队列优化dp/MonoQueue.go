package main

import (
	"fmt"
	"strings"
)

// https://leetcode.cn/problems/minimum-number-of-coins-for-fruits/description/
func main() {
	queue := NewMonoQueue(func(a, b V) bool { return a.v < b.v })
	queue.Append(V{1, 1})
	fmt.Println(queue)
}

type V = struct{ v, i int }

// 单调队列维护滑动窗口最小值.
// 单调队列队头元素为当前窗口最小值，队尾元素为当前窗口最大值.
type MonoQueue struct {
	MinQueue       []V
	_minQueueCount []int32
	_less          func(a, b V) bool
	_len           int
}

func NewMonoQueue(less func(a, b V) bool) *MonoQueue {
	return &MonoQueue{
		_less: less,
	}
}

func (q *MonoQueue) Append(value V) *MonoQueue {
	count := int32(1)
	for len(q.MinQueue) > 0 && q._less(value, q.MinQueue[len(q.MinQueue)-1]) {
		q.MinQueue = q.MinQueue[:len(q.MinQueue)-1]
		count += q._minQueueCount[len(q._minQueueCount)-1]
		q._minQueueCount = q._minQueueCount[:len(q._minQueueCount)-1]
	}
	q.MinQueue = append(q.MinQueue, value)
	q._minQueueCount = append(q._minQueueCount, count)
	q._len++
	return q
}

func (q *MonoQueue) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueue) Head() V {
	return q.MinQueue[0]
}

func (q *MonoQueue) Min() V {
	return q.MinQueue[0]
}

func (q *MonoQueue) Len() int {
	return q._len
}

func (q *MonoQueue) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair struct {
	value V
	count int32
}

func (p pair) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}
