package main

import "fmt"

// 多级分块结构的隐式树，用于查询区间聚合值.
// 可以传入log来控制每个块的大小，以平衡时间与空间复杂度.
// log=1 => 线段树，log=10 => 朴素分块.
// golang 用于内存管理的 pageAlloc 基数树中，log=3.
type RadixTree[E any] struct {
	e    func() E
	op   func(a, b E) E
	log  int
	mask int // (1<<log)-1
}

// log: 每个块的大小B=1<<log.
// e: 幺元.
// op: 结合律.
func NewRadixTree[E any](e func() E, op func(a, b E) E, log int) *RadixTree[E] {
	if log < 1 {
		log = 1
	}
	return &RadixTree[E]{e: e, op: op, log: log}
}

func (m *RadixTree[E]) Build(n int, f func(i int) E) {
	panic("todo")
}

func (m *RadixTree[E]) QueryRange(l, r int) E {
	panic("todo")
}

func (m *RadixTree[E]) QueryAll() E {
	panic("todo")

}

func (m *RadixTree[E]) Get(i int) E {
	panic("todo")
}

func (m *RadixTree[E]) GetAll() []E {
	panic("todo")
}

// A[i] = op(A[i], v).
func (m *RadixTree[E]) Update(i int, v E) {
	panic("todo")
}

// A[i] = v.
func (m *RadixTree[E]) Set(i int, v E) {
	panic("todo")
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (m *RadixTree[E]) MaxRight(left int, predicate func(E) bool) int {
	panic("todo")
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (m *RadixTree[E]) MinLeft(right int, predicate func(E) bool) int {
	panic("todo")
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

//
// cross checking
type naive[E any] struct {
	e         func() E
	op        func(a, b E) E
	log       int
	n         int
	data      []E
	blockSize int
}

func newNaive[E any](e func() E, op func(a, b E) E, log int) *naive[E] {
	if log < 1 {
		log = 1
	}
	return &naive[E]{e: e, op: op, log: log}
}

func (m *naive[E]) Build(n int, f func(i int) E) {
	m.n = n
	m.blockSize = 1 << m.log
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
}

func (m *naive[E]) QueryRange(l, r int) E {
	result := m.e() // start with the identity element
	for i := l; i < r; i++ {
		result = m.op(result, m.data[i])
	}
	return result
}

func (m *naive[E]) QueryAll() E {
	result := m.e()
	for i := 0; i < m.n; i++ {
		result = m.op(result, m.data[i])
	}
	return result
}

func (m *naive[E]) Get(i int) E {
	return m.data[i]
}

func (m *naive[E]) GetAll() []E {
	return m.data
}

func (m *naive[E]) Update(i int, v E) {
	m.data[i] = m.op(m.data[i], v)
}

func (m *naive[E]) Set(i int, v E) {
	m.data[i] = v
}

func (m *naive[E]) MaxRight(left int, predicate func(E) bool) int {
	sum := m.e()
	for right := left; right < m.n; right++ {
		sum = m.op(sum, m.data[right])
		if !predicate(sum) {
			return right
		}
	}
	return m.n
}

func (m *naive[E]) MinLeft(right int, predicate func(E) bool) int {
	sum := m.e()
	for left := right - 1; left >= 0; left-- {
		sum = m.op(m.data[left], sum)
		if !predicate(sum) {
			return left + 1
		}
	}
	return 0
}

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }

	tree := newNaive(e, op, 3)
	tree.Build(10, func(i int) int { return i + 1 })

	// Querying the range [2, 5)
	fmt.Println("Query range [2, 5):", tree.QueryRange(2, 5))

	// Get all values
	fmt.Println("Get All:", tree.GetAll())

	// Update value at index 3
	tree.Update(3, 5)
	fmt.Println("After update at index 3:", tree.GetAll())

	// Set value at index 6
	tree.Set(6, 10)
	fmt.Println("After set at index 6:", tree.GetAll())

	// Querying the range [2, 5)
	fmt.Println("Query range [2, 5):", tree.QueryRange(2, 5))

	tree.Build(10, func(i int) int { return i + 1 })
	fmt.Println("Query all:", tree.GetAll())
	// min left
	fmt.Println("Min left:", tree.MinLeft(10, func(x int) bool { return x < 27 }))
	// max right
	fmt.Println("Max right:", tree.MaxRight(0, func(x int) bool { return x <= 6 }))
}
