// https://37zigen.com/bankers-queue/

package main

type PersistentStack struct {
	next *PersistentStack
	val  int
}

func NewPersistentStack(val int, next *PersistentStack) *PersistentStack {
	return &PersistentStack{val: val, next: next}
}

func (s *PersistentStack) Top() int {
	return s.val
}

func (s *PersistentStack) Pop() *PersistentStack {
	if s.next != nil {
		return s.next
	}
	return nil
}

func Push(x *PersistentStack, v int) *PersistentStack {
	return NewPersistentStack(v, x)
}

func Concat(x, y *PersistentStack) *PersistentStack {
	if x == nil {
		return y
	} else {
		return NewPersistentStack(x.val, Concat(x.Pop(), y))
	}
}

func Reverse(head *PersistentStack) *PersistentStack {
	var ret *PersistentStack
	for x := head; x != nil; x = x.Pop() {
		ret = Push(ret, x.Top())
	}
	return ret
}

type PersistentQueue struct {
	fsize, rsize int
	f, r         *PersistentStack
}

func NewPersistentQueue(f *PersistentStack, fsize int, r *PersistentStack, rsize int) *PersistentQueue {
	return &PersistentQueue{fsize: fsize, rsize: rsize, f: f, r: r}
}

func (q *PersistentQueue) IsEmpty() bool {
	return q.fsize == 0
}

func (q *PersistentQueue) Top() int {
	return q.f.Top()
}

func (q *PersistentQueue) Push(x int) *PersistentQueue {
	return NewPersistentQueue(q.f, q.fsize, Push(q.r, x), q.rsize+1).Normalize()
}

func (q *PersistentQueue) Pop() *PersistentQueue {
	return NewPersistentQueue(q.f.Pop(), q.fsize-1, q.r, q.rsize).Normalize()
}

func (q *PersistentQueue) Normalize() *PersistentQueue {
	if q.fsize >= q.rsize {
		return q
	} else {
		return NewPersistentQueue(Concat(q.f, Reverse(q.r)), q.fsize+q.rsize, nil, 0)
	}
}
