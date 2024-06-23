package main

import (
	"fmt"
	"strings"
)

func main() {
	q := NewMonoQueue[int](func(a, b int) bool { return a < b })
	q.Append(5).Append(4).Append(3).Append(2).Append(1)
	fmt.Println(q.Min())
	q.Popleft()
	fmt.Println(q.Min())
}

// MonoQueueGeneric
// 单调队列维护滑动窗口最小值.
// 单调队列队头元素为当前窗口最小值，队尾元素为当前窗口最大值.
type MonoQueueGeneric[V comparable] struct {
	MinQueue       []V
	_minQueueCount []int32
	_less          func(a, b V) bool
	_len           int32
}

func NewMonoQueue[V comparable](less func(a, b V) bool) *MonoQueueGeneric[V] {
	return &MonoQueueGeneric[V]{_less: less}
}

func (q *MonoQueueGeneric[V]) Append(value V) *MonoQueueGeneric[V] {
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

func (q *MonoQueueGeneric[V]) Popleft() {
	q._minQueueCount[0]--
	if q._minQueueCount[0] == 0 {
		q.MinQueue = q.MinQueue[1:]
		q._minQueueCount = q._minQueueCount[1:]
	}
	q._len--
}

func (q *MonoQueueGeneric[V]) Head() V {
	return q.MinQueue[0]
}

func (q *MonoQueueGeneric[V]) Min() V {
	return q.MinQueue[0]
}

func (q *MonoQueueGeneric[V]) Len() int32 {
	return q._len
}

func (q *MonoQueueGeneric[V]) String() string {
	sb := []string{}
	for i := 0; i < len(q.MinQueue); i++ {
		sb = append(sb, fmt.Sprintf("%v", pair[V]{q.MinQueue[i], q._minQueueCount[i]}))
	}
	return fmt.Sprintf("MonoQueue{%v}", strings.Join(sb, ", "))
}

type pair[V comparable] struct {
	value V
	count int32
}

func (p pair[V]) String() string {
	return fmt.Sprintf("(value: %v, count: %v)", p.value, p.count)
}
