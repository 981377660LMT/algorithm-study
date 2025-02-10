package main

import (
	"fmt"
)

// 多级分块结构的隐式树，用于查询区间聚合值.
// 可以传入log来控制每个块的大小，以平衡时间与空间复杂度.
// log=1 => 线段树，log=10 => 朴素分块.
// golang 用于内存管理的 pageAlloc 基数树中，log=3.
type RadixTree[E any] struct {
	// props
	e         func() E
	op        func(a, b E) E
	log       int
	blockSize int

	// data
	n           int
	data        []E
	levels      [][]E
	levelShifts []int
}

// log: 每个块的大小B=1<<log.
// e: 幺元.
// op: 结合律.
func NewRadixTree[E any](e func() E, op func(a, b E) E, log int) *RadixTree[E] {
	if log < 1 {
		log = 1
	}
	return &RadixTree[E]{
		e:         e,
		op:        op,
		log:       log,
		blockSize: 1 << log,
	}
}

func (m *RadixTree[E]) Build(n int, f func(i int) E) {
	m.n = n
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
	m.levels = [][]E{}
	m.levelShifts = []int{}

	buildFrom := func(pre []E) []E {
		numPre := len(pre)
		cur := make([]E, (numPre+m.blockSize-1)>>m.log)
		for i := range cur {
			start := i << m.log
			end := min(start+m.blockSize, len(pre))
			v := m.e()
			for j := start; j < end; j++ {
				v = m.op(v, pre[j])
			}
			cur[i] = v
		}
		return cur
	}

	preLevel := m.data
	preShift := 0
	for len(preLevel) > 1 {
		curLevel := buildFrom(preLevel)
		m.levels = append(m.levels, curLevel)
		m.levelShifts = append(m.levelShifts, m.log*(preShift+1))
		preLevel = curLevel
		preShift++
	}
}

func (m *RadixTree[E]) QueryRange(l, r int) E {
	if l < 0 {
		l = 0
	}
	if r > m.n {
		r = m.n
	}
	if l >= r {
		return m.e()
	}

	res := m.e()
	type item struct{ l, r, k int }
	stack := []item{{l, r, len(m.levels) - 1}}

	for len(stack) > 0 {
		top := stack[len(stack)-1]
		stack = stack[:len(stack)-1]
		l, r, k := top.l, top.r, top.k
		if l >= r {
			continue
		}

		if k < 0 {
			for i := l; i < r; i++ {
				res = m.op(res, m.data[i])
			}
			continue
		}

		shift := m.levelShifts[k]
		startBlock := l >> shift
		endBlock := (r - 1) >> shift

		if startBlock >= endBlock {
			stack = append(stack, item{l, r, k - 1})
			continue
		}

		midStart, midEnd := startBlock+1, endBlock-1
		if midStart <= midEnd {
			if midStart < len(m.levels[k]) {
				midEnd = min(midEnd, len(m.levels[k])-1)
				for i := midStart; i <= midEnd; i++ {
					res = m.op(res, m.levels[k][i])
				}
			}
		}

		leftEnd := (startBlock + 1) << shift
		if leftEnd > r {
			leftEnd = r
		}
		if leftEnd > l {
			stack = append(stack, item{l, leftEnd, k - 1})
		}

		rightStart := endBlock << shift
		if rightStart < l {
			rightStart = l
		}
		if rightStart < r {
			stack = append(stack, item{rightStart, r, k - 1})
		}
	}
	return res
}

func (m *RadixTree[E]) QueryAll() E {
	if len(m.levels) == 0 || len(m.levels[len(m.levels)-1]) == 0 {
		return m.e()
	}
	return m.levels[len(m.levels)-1][0]
}

func (m *RadixTree[E]) Get(i int) E {
	if i < 0 || i >= m.n {
		return m.e()
	}
	return m.data[i]
}

// O(1).
func (m *RadixTree[E]) GetAll() []E {
	return m.data
}

// A[i] = op(A[i], v).
func (m *RadixTree[E]) Update(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.data[i] = m.op(m.data[i], v)
	for k := 0; k < len(m.levels); k++ {
		shift := m.levelShifts[k]
		bid := i >> shift
		start := bid << shift
		end := min(start+1<<shift, m.n)

		var v E
		if k == 0 {
			v = m.e()
			for j := start; j < end; j++ {
				v = m.op(v, m.data[j])
			}
		} else {
			preLevel := m.levels[k-1]
			preStart := bid * m.blockSize
			preEnd := min(preStart+m.blockSize, len(preLevel))
			v = m.e()
			for j := preStart; j < preEnd; j++ {
				v = m.op(v, preLevel[j])
			}
		}
		if bid < len(m.levels[k]) {
			m.levels[k][bid] = v
		}
	}
}

// A[i] = v.
func (m *RadixTree[E]) Set(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.data[i] = v
	for k := 0; k < len(m.levels); k++ {
		shift := m.levelShifts[k]
		bid := i >> shift
		start := bid << shift
		end := min(start+1<<shift, m.n)

		var v E
		if k == 0 {
			v = m.e()
			for j := start; j < end; j++ {
				v = m.op(v, m.data[j])
			}
		} else {
			preLevel := m.levels[k-1]
			preStart := bid * m.blockSize
			preEnd := min(preStart+m.blockSize, len(preLevel))
			v = m.e()
			for j := preStart; j < preEnd; j++ {
				v = m.op(v, preLevel[j])
			}
		}
		if bid < len(m.levels[k]) {
			m.levels[k][bid] = v
		}
	}
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (m *RadixTree[E]) MaxRight(left int, predicate func(E) bool) int {
	if left == m.n {
		return m.n
	}
	cur := m.e()
	if !predicate(cur) {
		return left
	}

	pos := left
	for k := 0; k < len(m.levels); k++ {
		shift := m.levelShifts[k]
		if pos&(1<<shift-1) == 0 && pos+1<<shift <= m.n {
			bid := pos >> shift
			if bid >= len(m.levels[k]) {
				continue
			}
			nv := m.op(cur, m.levels[k][bid])
			if predicate(nv) {
				cur = nv
				pos += 1 << shift
				k = -1 // Restart from level 0
			}
		}
	}

	for pos < m.n {
		nv := m.op(cur, m.Get(pos))
		if !predicate(nv) {
			break
		}
		cur = nv
		pos++
	}
	return pos
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (m *RadixTree[E]) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	cur := m.e()
	if !predicate(cur) {
		return right
	}

	pos := right
	for k := 0; k < len(m.levels); k++ {
		shift := m.levelShifts[k]
		if pos&(1<<shift-1) == 0 && pos >= 1<<shift {
			bid := (pos - 1) >> shift
			if bid >= len(m.levels[k]) {
				continue
			}
			nv := m.op(m.levels[k][bid], cur)
			if predicate(nv) {
				cur = nv
				pos -= 1 << shift
				k = -1 // Restart from level 0
			}
		}
	}

	for pos > 0 {
		nv := m.op(m.Get(pos-1), cur)
		if !predicate(nv) {
			break
		}
		cur = nv
		pos--
	}
	return pos
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
