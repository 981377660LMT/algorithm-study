package main

import (
	"fmt"
	"math/rand"
	"time"
)

func main() {
	testTime()
	for i := 0; i < 1000; i++ {
		test()
	}
	fmt.Println("pass")
}

type RadixTreeLazy[E any, Id comparable] struct {
	e           func() E
	id          func() Id
	op          func(a, b E) E
	mapping     func(f Id, e E, size int) E
	composition func(f, g Id) Id
	log         int
	blockSize   int

	n           int
	layers      [][]E
	layerShifts []int
	lazys       [][]Id
}

func NewRadixTreeLazy[E any, Id comparable](
	e func() E,
	id func() Id,
	op func(a, b E) E,
	mapping func(f Id, e E, size int) E,
	composition func(f, g Id) Id,
	log int,
) *RadixTreeLazy[E, Id] {
	if log < 1 {
		log = 3
	}
	return &RadixTreeLazy[E, Id]{
		e:           e,
		op:          op,
		mapping:     mapping,
		composition: composition,
		id:          id,
		log:         log,
		blockSize:   1 << log,
	}
}

func (m *RadixTreeLazy[E, Id]) Build(n int, f func(int) E) {
	m.n = n
	level0 := make([]E, n)
	for i := 0; i < n; i++ {
		level0[i] = f(i)
	}
	m.layers = [][]E{level0}
	m.layerShifts = []int{0}

	m.lazys = [][]Id{make([]Id, n)}
	for i := range m.lazys[0] {
		m.lazys[0][i] = m.id()
	}

	preLevel := level0
	preShift := 1
	for {
		curSize := (len(preLevel) + m.blockSize - 1) / m.blockSize
		if curSize == 0 {
			break
		}
		curLevel := make([]E, curSize)
		curLazy := make([]Id, curSize)
		for i := range curLazy {
			curLazy[i] = m.id()
		}
		for i := 0; i < curSize; i++ {
			start := i * m.blockSize
			end := min(start+m.blockSize, len(preLevel))
			if start >= len(preLevel) {
				curLevel[i] = m.e()
				continue
			}
			val := preLevel[start]
			for j := start + 1; j < end; j++ {
				val = m.op(val, preLevel[j])
			}
			curLevel[i] = val
		}
		m.layers = append(m.layers, curLevel)
		m.lazys = append(m.lazys, curLazy)
		m.layerShifts = append(m.layerShifts, m.log*preShift)
		preLevel = curLevel
		preShift++
		if curSize == 1 {
			break
		}
	}
}

func (m *RadixTreeLazy[E, Id]) propagate(k, blockIndex int) {
	if k == 0 || m.lazys[k][blockIndex] == m.id() {
		return
	}

	blockSizePrev := 1 << m.layerShifts[k-1]
	startSubBlock := blockIndex * m.blockSize
	endSubBlock := startSubBlock + m.blockSize
	if endSubBlock > len(m.layers[k-1]) {
		endSubBlock = len(m.layers[k-1])
	}

	f := m.lazys[k][blockIndex]
	subLazys := m.lazys[k-1]
	subLayers := m.layers[k-1]

	for sb := startSubBlock; sb < endSubBlock; sb++ {
		subLazys[sb] = m.composition(f, subLazys[sb])
		subLayers[sb] = m.mapping(f, subLayers[sb], blockSizePrev)
	}

	m.lazys[k][blockIndex] = m.id()

	if startSubBlock < len(subLayers) {
		val := subLayers[startSubBlock]
		for sb := startSubBlock + 1; sb < endSubBlock; sb++ {
			val = m.op(val, subLayers[sb])
		}
		m.layers[k][blockIndex] = val
	} else {
		m.layers[k][blockIndex] = m.e()
	}
}

func (m *RadixTreeLazy[E, Id]) apply(k, blockIndex int, value Id) {
	m.lazys[k][blockIndex] = m.composition(value, m.lazys[k][blockIndex])
	blockSize := 1 << m.layerShifts[k]
	m.layers[k][blockIndex] = m.mapping(value, m.layers[k][blockIndex], blockSize)
}

func (m *RadixTreeLazy[E, Id]) Update(l, r int, value Id) {
	if l < 0 {
		l = 0
	}
	if r > m.n {
		r = m.n
	}
	if l >= r {
		return
	}

	i := l
	for i < r {
		found := false
		for k := len(m.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << m.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= r {
				blockIndex := i >> m.layerShifts[k]
				m.apply(k, blockIndex, value)
				i += blockSize
				found = true
				break
			}
		}
		if !found {
			for k := len(m.layerShifts) - 1; k >= 1; k-- {
				blockIndex := i >> m.layerShifts[k]
				m.propagate(k, blockIndex)
			}
			m.apply(0, i, value)
			i++
		}
	}
}

func (m *RadixTreeLazy[E, Id]) QueryRange(l, r int) E {
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
	i := l
	for i < r {
		jumped := false
		for k := len(m.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << m.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= r {
				blockIndex := i >> m.layerShifts[k]
				res = m.op(res, m.layers[k][blockIndex])
				i += blockSize
				jumped = true
				break
			}
		}
		if !jumped {
			for k := len(m.layerShifts) - 1; k >= 1; k-- {
				blockIndex := i >> m.layerShifts[k]
				m.propagate(k, blockIndex)
			}
			res = m.op(res, m.layers[0][i])
			i++
		}
	}
	return res
}

func (m *RadixTreeLazy[E, Id]) QueryAll() E {
	if len(m.layers) == 0 {
		return m.e()
	}
	return m.layers[len(m.layers)-1][0]
}

func (m *RadixTreeLazy[E, Id]) Get(i int) E {
	if i < 0 || i >= m.n {
		return m.e()
	}
	for k := len(m.layerShifts) - 1; k >= 1; k-- {
		blockIndex := i >> m.layerShifts[k]
		m.propagate(k, blockIndex)
	}
	return m.layers[0][i]
}

func (m *RadixTreeLazy[E, Id]) Set(i int, value E) {
	if i < 0 || i >= m.n {
		return
	}
	for k := len(m.layerShifts) - 1; k >= 1; k-- {
		blockIndex := i >> m.layerShifts[k]
		m.propagate(k, blockIndex)
	}
	m.layers[0][i] = value
	for k := 1; k < len(m.layers); k++ {
		blockIndex := i >> m.layerShifts[k]
		start := blockIndex * m.blockSize
		end := min(start+m.blockSize, len(m.layers[k-1]))
		val := m.layers[k-1][start]
		for j := start + 1; j < end; j++ {
			val = m.op(val, m.layers[k-1][j])
		}
		m.layers[k][blockIndex] = val
	}
}

func min(a, b int) int {
	if a < b {
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
		log = 6
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

func (m *naive[E]) UpdateRange(l, r int, v E) {
	for i := l; i < r; i++ {
		m.data[i] = m.op(m.data[i], v)
	}
}

func (m *naive[E]) Set(i int, v E) {
	m.data[i] = v
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
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

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
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
	const INF int = 1e18
	e := func() int { return 0 }
	id := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	mapping := func(f int, x int, segLen int) int {
		return f + x
	}
	composition := func(f, g int) int {
		return f + g
	}

	N := rand.Intn(1000) + 1
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	rt1 := NewRadixTreeLazy(e, id, op, mapping, composition, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })
	rt2 := newNaive(e, op, -1)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(5)
		switch op {
		case 0:
			// l, r := rand.Intn(N), rand.Intn(N)
			// if l > r {
			// 	l, r = r, l
			// }
			// v := rand.Intn(100)
			// rt1.Update(l, r, v)
			// rt2.UpdateRange(l, r, v)
		case 1:
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Update(i, i+1, v)
			rt2.UpdateRange(i, i+1, v)
		case 2:
			// Get
			i := rand.Intn(N)
			if rt1.Get(i) != rt2.Get(i) {
				panic("err Get")
			}
		case 3:
			// GetAll
			// nums1, nums2 := rt1.GetAll(), rt2.GetAll()
			// if slices.Compare(nums1, nums2) != 0 {
			// 	panic("err GetAll")
			// }

		case 4:
			// Set
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Set(i, v)
			rt2.Set(i, v)

		}
	}
}

func testTime() {
	e := func() int { return 0 }
	id := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	mapping := func(f int, x int, segLen int) int {
		return f + x
	}
	composition := func(f, g int) int {
		return f + g
	}

	N := int(2e5)
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	time1 := time.Now()
	rt1 := NewRadixTreeLazy(e, id, op, mapping, composition, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })

	for i := 0; i < N; i++ {

		rt1.QueryAll()
		rt1.Get(i)
		rt1.Set(i, i)

	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
