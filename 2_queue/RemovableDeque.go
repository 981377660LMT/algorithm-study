// 可根据值删除元素的双端队列.(一次删除会删除deque中所有值为value的元素.)

package main

import (
	"fmt"
	"time"
)

func main() {
	rq := NewRemovableDeque(10)
	rq.Append(1)
	rq.Append(2)
	rq.Remove(2)
	fmt.Println(rq.Pop())
	fmt.Println(rq.Len())
	time1 := time.Now()
	Q := NewRemovableDeque(2e5)

	for i := 0; i < 2e5; i++ {
		Q.Append(i)
		Q.Append(i)
		Q.PopLeft()
		Q.Remove(i)
	}
	time2 := time.Now()
	fmt.Println(time2.Sub(time1))
}

type Value = int

type Pair = struct {
	value     Value
	addedTime int
}

type RemovableDeque struct {
	queue       *dq
	counter     map[Value]int
	removedTime map[Value]int
	length      int
	time        int
}

func NewRemovableDeque(cap int) *RemovableDeque {
	return &RemovableDeque{
		queue:       newDq(cap),
		counter:     make(map[Value]int),
		removedTime: make(map[Value]int),
		length:      0,
		time:        0,
	}
}

func (rq *RemovableDeque) Append(value Value) {
	rq.length++
	rq.queue.Append(Pair{value, rq.time})
	rq.counter[value]++
}

func (rq *RemovableDeque) AppendLeft(value Value) {
	rq.length++
	rq.queue.AppendLeft(Pair{value, rq.time})
	rq.counter[value]++
}

func (rq *RemovableDeque) Pop() Value {
	rq.length--
	rq._normalizeTail()
	res := rq.queue.Pop().value
	if _, ok := rq.counter[res]; ok {
		rq.counter[res]--
		if rq.counter[res] == 0 {
			delete(rq.counter, res)
		}
	}
	return res
}

func (rq *RemovableDeque) PopLeft() Value {
	rq.length--
	rq._normalizeHead()
	res := rq.queue.PopLeft().value
	if _, ok := rq.counter[res]; ok {
		rq.counter[res]--
		if rq.counter[res] == 0 {
			delete(rq.counter, res)
		}
	}
	return res
}

func (rq *RemovableDeque) Head() Value {
	rq._normalizeHead()
	return rq.queue.Head().value
}

func (rq *RemovableDeque) Tail() Value {
	rq._normalizeTail()
	return rq.queue.Tail().value
}

// 删除deque中所有值为value的元素.
func (rq *RemovableDeque) Remove(value Value) {
	if _, ok := rq.counter[value]; ok {
		rq.length -= rq.counter[value]
		delete(rq.counter, value)
		rq.removedTime[value] = rq.time
		rq.time++
	}
}

func (rq *RemovableDeque) Count(value Value) int {
	return rq.counter[value]
}

func (rq *RemovableDeque) Empty() bool {
	return rq.length == 0
}

func (rq *RemovableDeque) Len() int {
	return rq.length
}

func (rq *RemovableDeque) String() string {
	res := make([]Value, 0, rq.length)
	for i := 0; i < rq.length; i++ {
		p := rq.queue.At(i)
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			continue
		}
		res = append(res, v)
	}
	return fmt.Sprint(res)
}

func (rq *RemovableDeque) _normalizeHead() {
	for rq.queue.Size() > 0 {
		p := rq.queue.Head()
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			rq.queue.PopLeft()
		} else {
			break
		}
	}
}

func (rq *RemovableDeque) _normalizeTail() {
	for rq.queue.Size() > 0 {
		p := rq.queue.Tail()
		v, t := p.value, p.addedTime
		if _, ok := rq.removedTime[v]; ok && t <= rq.removedTime[v] {
			rq.queue.Pop()
		} else {
			break
		}
	}
}

type dq struct{ l, r []Pair }

func newDq(cap int) *dq { return &dq{make([]Pair, 0, 1+cap/2), make([]Pair, 0, 1+cap/2)} }

func (q *dq) Empty() bool {
	return len(q.l) == 0 && len(q.r) == 0
}

func (q *dq) Size() int {
	return len(q.l) + len(q.r)
}

func (q *dq) AppendLeft(v Pair) {
	q.l = append(q.l, v)
}

func (q *dq) Append(v Pair) {
	q.r = append(q.r, v)
}

func (q *dq) PopLeft() Pair {
	var v Pair
	if len(q.l) > 0 {
		q.l, v = q.l[:len(q.l)-1], q.l[len(q.l)-1]
	} else {
		v, q.r = q.r[0], q.r[1:]
	}
	return v
}

func (q *dq) Pop() Pair {
	var v Pair
	if len(q.r) > 0 {
		q.r, v = q.r[:len(q.r)-1], q.r[len(q.r)-1]
	} else {
		v, q.l = q.l[0], q.l[1:]
	}
	return v
}

func (q *dq) Head() Pair {
	if len(q.l) > 0 {
		return q.l[len(q.l)-1]
	}
	return q.r[0]
}

func (q *dq) Tail() Pair {
	if len(q.r) > 0 {
		return q.r[len(q.r)-1]
	}
	return q.l[0]
}

// 0 <= i < q.Size()
func (q *dq) At(i int) Pair {
	if i < len(q.l) {
		return q.l[len(q.l)-1-i]
	}
	return q.r[i-len(q.l)]
}
