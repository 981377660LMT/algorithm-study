package main

import (
	"fmt"
	"math/rand"
	"slices"
	"time"
)

func main() {
	testTime()
	for i := 0; i < 1000; i++ {
		test()
	}
	fmt.Println("pass")
}

// 多级分块结构的隐式树，用于查询区间聚合值.
// 可以传入log来控制每个块的大小，以平衡时间与空间复杂度.
// log=1 => 线段树，log=10 => 朴素分块.
// golang 用于内存管理的 pageAlloc 基数树中，log=3.
type RadixTree[E any] struct {
	e         func() E
	op        func(a, b E) E
	log       int
	blockSize int

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
		log = 6
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

	build := func(pre []E) []E {
		cur := make([]E, (len(pre)+m.blockSize-1)>>m.log)
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
	preShift := 1
	for len(preLevel) > 1 {
		curLevel := build(preLevel)
		m.levels = append(m.levels, curLevel)
		m.levelShifts = append(m.levelShifts, m.log*preShift)
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
	return m.queryRangeRecursive(l, r, len(m.levels)-1)
}

func (m *RadixTree[E]) queryRangeRecursive(l, r, k int) E {
	if l >= r {
		return m.e()
	}

	if k < 0 {
		res := m.e()
		for i := l; i < r; i++ {
			res = m.op(res, m.data[i])
		}
		return res
	}

	shift := m.levelShifts[k]
	startBlock := l >> shift
	endBlock := (r - 1) >> shift
	if startBlock == endBlock {
		return m.queryRangeRecursive(l, r, k-1)
	}

	res := m.e()

	for i := startBlock + 1; i < endBlock; i++ {
		res = m.op(res, m.levels[k][i])
	}

	leftEnd := (startBlock + 1) << shift
	if leftEnd > l {
		res = m.op(m.queryRangeRecursive(l, leftEnd, k-1), res)
	}

	rightStart := endBlock << shift
	if rightStart < r {
		res = m.op(res, m.queryRangeRecursive(rightStart, r, k-1))
	}

	return res
}

func (m *RadixTree[E]) QueryAll() E {
	if m.n == 0 {
		return m.e()
	}
	if m.n == 1 {
		return m.data[0]
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
	pre := m.data
	for k := 0; k < len(m.levels); k++ {
		bid := i >> m.levelShifts[k]
		start := bid << m.log
		end := min(start+m.blockSize, len(pre))
		v = m.e()
		for j := start; j < end; j++ {
			v = m.op(v, pre[j])
		}
		m.levels[k][bid] = v
		pre = m.levels[k]
	}
}

// A[i] = v.
func (m *RadixTree[E]) Set(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.data[i] = v
	pre := m.data
	for k := 0; k < len(m.levels); k++ {
		bid := i >> m.levelShifts[k]
		start := bid << m.log
		end := min(start+m.blockSize, len(pre))
		v = m.e()
		for j := start; j < end; j++ {
			v = m.op(v, pre[j])
		}
		m.levels[k][bid] = v
		pre = m.levels[k]
	}
}

func (m *RadixTree[E]) MaxRight(l int, f func(E) bool) int {
	if l < 0 {
		l = 0
	}
	if l >= m.n {
		return m.n
	}
	return m.maxRightRecursive(l, f, len(m.levels)-1)
}

func (m *RadixTree[E]) maxRightRecursive(l int, f func(E) bool, k int) int {
	panic("TODO")
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

func (m *naive[E]) MaxRight(l int, f func(E) bool) int {
	sum := m.e()
	for i := l; i < m.n; i++ {
		sum = m.op(sum, m.data[i])
		if !f(sum) {
			return i
		}
	}
	return m.n
}

func (m *naive[E]) MinLeft(r int, f func(E) bool) int {
	sum := m.e()
	for i := r - 1; i >= 0; i-- {
		sum = m.op(m.data[i], sum)
		if !f(sum) {
			return i + 1
		}
	}
	return 0
}

func test() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }

	N := rand.Intn(1000) + 1
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	rt1 := NewRadixTree(e, op, 2)
	rt1.Build(N, func(i int) int { return randNums[i] })
	rt2 := newNaive(e, op, 2)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(8)
		switch op {
		case 0:
			l, r := rand.Intn(N), rand.Intn(N)

			if rt1.QueryRange(l, r) != rt2.QueryRange(l, r) {
				panic("err QueryRange")
			}
		case 1:
			if r1, r2 := rt1.QueryAll(), rt2.QueryAll(); r1 != r2 {
				fmt.Println(rt1.GetAll(), rt2.GetAll())
				panic(fmt.Sprintf("err QueryAll: %v %v", r1, r2))
			}
		case 2:
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Update(i, v)
			rt2.Update(i, v)
		case 3:
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Set(i, v)
			rt2.Set(i, v)

		case 4:
			// Get
			i := rand.Intn(N)
			if rt1.Get(i) != rt2.Get(i) {
				panic("err Get")
			}
		case 5:
			// GetAll
			nums1, nums2 := rt1.GetAll(), rt2.GetAll()
			if slices.Compare(nums1, nums2) != 0 {
				panic("err GetAll")
			}
		case 6:
			// QueryAll
			if rt1.QueryAll() != rt2.QueryAll() {
				panic("err QueryAll")
			}
		case 7:
			// MaxRight
			l := rand.Intn(N)
			f := func(v int) bool { return v < 100 }
			if rt1.MaxRight(l, f) != rt2.MaxRight(l, f) {
				panic("err MaxRight")
			}

		}
	}
}

func testTime() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }

	N := int(2e5)
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	time1 := time.Now()
	rt1 := NewRadixTree(e, op, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })

	for i := 0; i < N; i++ {
		rt1.QueryRange(i, N)
		rt1.QueryAll()
		rt1.Get(i)
		rt1.Set(i, i)
	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
