// !FIXME: UpdateRange 有问题，不要使用

package main

import (
	"fmt"
	"math/rand"
)

func main() {
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
		m.layerShifts = append(m.layerShifts, m.log*(preShift-1))
		preLevel = curLevel
		preShift++
		if curSize == 1 {
			break
		}
	}
}

func (m *RadixTreeLazy[E, Id]) UpdateRange(l, r int, value Id) {
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
				m.propagate(k, blockIndex, value)
				i += blockSize
				found = true
				break
			}
		}
		if !found {
			for k := len(m.layerShifts) - 1; k >= 1; k-- {
				blockIndex := i >> m.layerShifts[k]
				m.pushDown(k, blockIndex)
			}
			m.propagate(0, i, value)
			i++
		}
	}
}

func (m *RadixTreeLazy[E, Id]) QueryAll() E {
	if len(m.layers) == 0 {
		return m.e()
	}

	return m.layers[len(m.layers)-1][0]
}

func (m *RadixTreeLazy[E, Id]) pushDown(k, blockIndex int) {
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

func (m *RadixTreeLazy[E, Id]) propagate(k, blockIndex int, value Id) {
	m.lazys[k][blockIndex] = m.composition(value, m.lazys[k][blockIndex])
	blockSize := 1 << m.layerShifts[k]
	m.layers[k][blockIndex] = m.mapping(value, m.layers[k][blockIndex], blockSize)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// cross checking
type naive[E any, Id comparable] struct {
	e           func() E
	id          func() Id
	op          func(a, b E) E
	mapping     func(f E, x E, segLen int) E
	composition func(f, g E) E

	log       int
	n         int
	data      []E
	lazy      []Id
	blockSize int
}

func newNaive[E any, Id comparable](
	e func() E,
	id func() Id,
	op func(a, b E) E,
	mapping func(f E, x E, segLen int) E,
	composition func(f, g E) E,
	log int,
) *naive[E, Id] {
	if log < 1 {
		log = 6
	}
	return &naive[E, Id]{e: e, id: id, op: op, mapping: mapping, composition: composition, log: log}
}

func (m *naive[E, Id]) Build(n int, f func(i int) E) {
	m.n = n
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
	m.blockSize = 1 << m.log
}

func (m *naive[E, Id]) QueryAll() E {
	result := m.e()
	for i := 0; i < m.n; i++ {
		result = m.op(result, m.data[i])
	}
	return result
}

// 对区间 [l, r) 内的每个位置更新：用 mapping 处理 lazy 值 v
func (m *naive[E, Id]) UpdateRange(l, r int, v E) {
	for i := l; i < r; i++ {
		m.data[i] = m.mapping(v, m.data[i], 1)
	}
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
	rt2 := newNaive(e, id, op, mapping, composition, -1)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(2)
		switch op {
		case 0:
			// UpdateRange
			l, r := rand.Intn(N), rand.Intn(N)
			if l > r {
				l, r = r, l
			}
			v := rand.Intn(100)
			rt1.UpdateRange(l, r, v)
			rt2.UpdateRange(l, r, v)
		case 1:
			// QueryAll
			if res1, res2 := rt1.QueryAll(), rt2.QueryAll(); res1 != res2 {
				fmt.Println("QueryAll failed")
				fmt.Println(res1, res2)
				panic("")
			}
		}
	}
}
