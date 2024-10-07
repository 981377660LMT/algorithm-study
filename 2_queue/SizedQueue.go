package main

import "fmt"

func main() {
	Q := NewSizedQueue(3)
	Q.Append(1)
	Q.Append(2)
	Q.Append(3)
	Q.Append(4)
	fmt.Println(Q.Sum()) // 9
	fmt.Println(Q)
}

type SizedQueue struct {
	sum   int
	size  int
	queue []int
}

func NewSizedQueue(size int) *SizedQueue {
	return &SizedQueue{size: size}
}

func (q *SizedQueue) Append(v int) {
	q.queue = append(q.queue, v)
	q.sum += v
	if len(q.queue) > q.size {
		q.sum -= q.queue[0]
		q.queue = q.queue[1:]
	}
}

func (q *SizedQueue) PopLeft() int {
	res := q.queue[0]
	q.queue = q.queue[1:]
	q.sum -= res
	return res
}

func (q *SizedQueue) Sum() int {
	return q.sum
}

func (q *SizedQueue) At(i int) int {
	return q.queue[i]
}

func (q *SizedQueue) Len() int {
	return len(q.queue)
}

func (q *SizedQueue) String() string {
	res := "SizedQueue{"
	for i := 0; i < len(q.queue); i++ {
		res += fmt.Sprint(q.queue[i])
		if i != len(q.queue)-1 {
			res += ","
		}
	}
	res += "}"
	return res
}
