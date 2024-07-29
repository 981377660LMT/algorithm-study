package main

import "fmt"

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }
	bit := NewFenwickTree(e, op, inv)
	bit.Build(10, func(i int32) int { return int(i) })
	fmt.Println(bit.GetAll())
	bit.Update(1, 8)
	fmt.Println(bit.GetAll())
}

type FenwickTree[E any] struct {
	n     int32
	total E
	data  []E
	e     func() E
	op    func(e1, e2 E) E
	inv   func(e E) E
}

func NewFenwickTree[E any](e func() E, op func(e1, e2 E) E, inv func(e E) E) *FenwickTree[E] {
	return &FenwickTree[E]{e: e, op: op, inv: inv}
}

func (fw *FenwickTree[E]) Build(n int32, f func(i int32) E) {
	data := make([]E, n)
	for i := int32(0); i < n; i++ {
		data[i] = f(i)
	}
	for i := int32(1); i <= n; i++ {
		if j := i + (i & -i); j <= n {
			data[j-1] = fw.op(data[i-1], data[j-1])
		}
	}
	fw.n = n
	fw.data = data
	fw.total = fw.QueryPrefix(n)
}

func (fw *FenwickTree[E]) QueryAll() E { return fw.total }

// [0, end)
func (fw *FenwickTree[E]) QueryPrefix(end int32) E {
	if end > fw.n {
		end = fw.n
	}
	res := fw.e()
	for ; end > 0; end &= end - 1 {
		res = fw.op(res, fw.data[end-1])
	}
	return res
}

// [start, end)
func (fw *FenwickTree[E]) QueryRange(start, end int32) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start > end {
		return fw.e()
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	pos, neg := fw.e(), fw.e()
	for end > start {
		pos = fw.op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = fw.op(neg, fw.data[start-1])
		start &= start - 1
	}
	return fw.op(pos, fw.inv(neg))
}

// 要求op满足交换律(commute).
func (fw *FenwickTree[E]) Update(i int32, x E) {
	fw.total = fw.op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = fw.op(fw.data[i-1], x)
	}
}

func (fw *FenwickTree[E]) GetAll() []E {
	res := make([]E, fw.n)
	for i := int32(0); i < fw.n; i++ {
		res[i] = fw.QueryRange(i, i+1)
	}
	return res
}
