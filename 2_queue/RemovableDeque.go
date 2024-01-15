package main

import "fmt"

func main() {
	rq := NewRemovableQueue(1, 2, 3)
	rq.Remove(2)
	fmt.Println(rq.PopLeft())
	fmt.Println(rq.Len())
	fmt.Println(rq.PopLeft())
}

type Value = int

type RemovableQueue struct {
	queue, removedQueue []Value
}

func NewRemovableQueue(values ...Value) *RemovableQueue {
	return &RemovableQueue{
		queue: append(values[:0:0], values...),
	}
}

func (rq *RemovableQueue) Append(value Value) {
	rq.queue = append(rq.queue, value)
}

func (rq *RemovableQueue) PopLeft() Value {
	rq.refresh()
	res := rq.queue[0]
	rq.queue = rq.queue[1:]
	return res
}

func (rq *RemovableQueue) Top() Value {
	rq.refresh()
	return rq.queue[0]
}

// 删除前必须保证value存在于队列.
func (rq *RemovableQueue) Remove(value Value) {
	rq.removedQueue = append(rq.removedQueue, value)
}

func (rq *RemovableQueue) Empty() bool {
	return rq.Len() == 0
}

func (rq *RemovableQueue) Len() int {
	return len(rq.queue) - len(rq.removedQueue)
}

func (rq *RemovableQueue) refresh() {
	for len(rq.removedQueue) > 0 && rq.removedQueue[0] == rq.queue[0] {
		rq.removedQueue = rq.removedQueue[1:]
		rq.queue = rq.queue[1:]
	}
}
