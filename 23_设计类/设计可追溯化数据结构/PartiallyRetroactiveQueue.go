// // https://noshi91.github.io/Library/data_structure/partially_retroactive_queue.cpp
// // PartiallyRetroactiveQueue
// // 部分可追溯队列
// // !NotVerified.

// package main

// import (
// 	"fmt"
// )

// func main() {
// 	queue := NewPartiallyRetroactiveQueue[int]()
// 	fmt.Println(queue.Empty()) // true
// 	fmt.Println(queue.GetAll())
// 	queue.InsertEnqueue(1, nil)
// 	fmt.Println(queue.GetAll())
// 	queue.InsertDequeue(nil)
// 	fmt.Println(queue.GetAll())
// 	time2 := queue.InsertEnqueue(2, nil)
// 	fmt.Println(queue.GetAll())
// 	queue.Delete(time2)
// 	// fmt.Println(queue.Front(), queue.Back())
// }

// // !NotVerified.
// type PartiallyRetroactiveQueue[T any] struct {
// 	inited bool
// 	front  *node[T]
// 	back   *node[T]
// }

// func NewPartiallyRetroactiveQueue[T any]() *PartiallyRetroactiveQueue[T] {
// 	return &PartiallyRetroactiveQueue[T]{}
// }

// func (q *PartiallyRetroactiveQueue[T]) Front() T {
// 	if q.front == nil {
// 		panic("empty queue")
// 	}
// 	return q.front.Value
// }

// func (q *PartiallyRetroactiveQueue[T]) Back() T {
// 	if q.back == nil {
// 		panic("empty queue")
// 	}
// 	return q.back.Value
// }

// func (q *PartiallyRetroactiveQueue[T]) Empty() bool {
// 	return !q.inited
// }

// // before 为nil时，表示在当前时间插入.
// func (q *PartiallyRetroactiveQueue[T]) InsertEnqueue(value T, beforeTime *timePointer[T]) *timePointer[T] {
// 	if !q.inited {
// 		node := &node[T]{Value: value}
// 		q.front = node
// 		q.back = node
// 		node.isBeforeF = true
// 		q.inited = true
// 		return &timePointer[T]{Pointer: node, isEnqueue: true}
// 	}
// 	if beforeTime == nil {
// 		node := &node[T]{Prev: q.back, Value: value}
// 		q.back.Next = node
// 		q.back = node
// 		return &timePointer[T]{Pointer: node, isEnqueue: true}
// 	} else {
// 		ptr := beforeTime.Pointer
// 		node := &node[T]{Prev: ptr.Prev, Next: ptr, Value: value}
// 		if ptr.Prev != nil {
// 			ptr.Prev.Next = node
// 		}
// 		ptr.Prev = node
// 		if ptr.isBeforeF {
// 			q.front.isBeforeF = false
// 			q.front = q.front.Prev
// 		}
// 		return &timePointer[T]{Pointer: node, isEnqueue: true}
// 	}
// }

// // before 为nil时，表示在当前时间插入.
// func (q *PartiallyRetroactiveQueue[T]) InsertDequeue(beforeTime *timePointer[T]) *timePointer[T] {
// 	if q.front == nil && q.back == nil {
// 		q.inited = false
// 		return &timePointer[T]{Pointer: q.front, isEnqueue: false}
// 	}
// 	q.front = q.front.Next
// 	if q.front != nil {
// 		q.front.isBeforeF = true
// 	} else {
// 		q.back = nil
// 		q.inited = false
// 	}
// 	return &timePointer[T]{Pointer: q.front, isEnqueue: false}
// }

// // 删除操作.
// func (q *PartiallyRetroactiveQueue[T]) Delete(time *timePointer[T]) {
// 	ptr, isEnq := time.Pointer, time.isEnqueue
// 	if isEnq {
// 		if ptr.Next != nil {
// 			ptr.Next.Prev = ptr.Prev
// 		} else {
// 			q.back = ptr.Prev
// 		}
// 		if ptr.Prev != nil {
// 			ptr.Prev.Next = ptr.Next
// 		}
// 		if ptr.isBeforeF {
// 			q.front = q.front.Next
// 			q.front.isBeforeF = true
// 		}
// 	} else {
// 		q.front.isBeforeF = false
// 		q.front = q.front.Prev
// 	}
// }

// func (q *PartiallyRetroactiveQueue[T]) GetAll() []T {
// 	res := []T{}
// 	for ptr := q.back; ptr != nil; ptr = ptr.Prev {
// 		res = append(res, ptr.Value)
// 		if ptr == q.front {
// 			break
// 		}
// 	}
// 	return res
// }

// type node[T any] struct {
// 	isBeforeF  bool
// 	Prev, Next *node[T]
// 	Value      T
// }

// func (n *node[T]) String() string {
// 	if n == nil {
// 		return "nil"
// 	}
// 	return fmt.Sprintf("%v", n.Value)
// }

// type timePointer[T any] struct {
// 	isEnqueue bool
// 	Pointer   *node[T]
// }
