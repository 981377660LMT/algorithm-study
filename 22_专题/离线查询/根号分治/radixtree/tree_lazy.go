package main

import "fmt"

func main() {
	// rangeAddrangeMax
	e := func() int { return 0 }
	id := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	mapping := func(f, x int) int { return f + x }
	composition := func(f, g int) int { return f + g }

	seg := NewRadixTreeLazy(e, id, op, mapping, composition, -1)
	seg.Build(10, func(i int) int { return 0 })
	fmt.Println(seg.QueryRange(0, 10)) // 0
	seg.UpdateRange(0, 10, 1)
	fmt.Println(seg.QueryRange(0, 10)) // 10

}

type RadixTreeLazy[E any, Id comparable] struct {
	e           func() E
	id          func() Id
	op          func(a, b E) E
	mapping     func(f Id, x E) E
	composition func(f, g Id) Id
	log         int
	blockSize   int

	n           int
	data        []E
	levels      [][]E
	levelShifts []int
	lazy        [][]Id
}

func NewRadixTreeLazy[E any, Id comparable](
	e func() E,
	id func() Id,
	op func(a, b E) E,
	mapping func(f Id, x E) E,
	composition func(f, g Id) Id,
	log int,
) *RadixTreeLazy[E, Id] {
	if log < 1 {
		log = 6
	}
	return &RadixTreeLazy[E, Id]{
		e:           e,
		id:          id,
		op:          op,
		mapping:     mapping,
		composition: composition,
		log:         log,
		blockSize:   1 << log,
	}
}

func (m *RadixTreeLazy[E, Id]) Build(n int, f func(i int) E) {
	m.n = n
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
	m.levels = [][]E{}
	m.levelShifts = []int{}
	m.lazy = [][]Id{}

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
		curLazy := make([]Id, len(curLevel))
		for i := range curLazy {
			curLazy[i] = m.id()
		}
		m.lazy = append(m.lazy, curLazy)
		preLevel = curLevel
		preShift++
	}
}

func (m *RadixTreeLazy[E, Id]) UpdateRange(l, r int, f Id) {
	if l < 0 {
		l = 0
	}
	if r > m.n {
		r = m.n
	}
	if l >= r {
		return
	}
	m.updateRange(0, len(m.levels)-1, l, r, f)
}

func (m *RadixTreeLazy[E, Id]) updateRange(k, level, l, r int, f Id) {

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
	return m.queryRange(len(m.levels)-1, l, r)
}
func (m *RadixTreeLazy[E, Id]) queryRange(level, l, r int) E {
	if level < 0 {
		res := m.e()
		for i := l; i < r; i++ {
			res = m.op(res, m.data[i])
		}
		return res
	}

	shift := m.levelShifts[level]
	blockSize := 1 << shift
	firstBlock := l >> shift
	lastBlock := (r - 1) >> shift

	// Push down lazy
	m.pushDown(level, firstBlock)
	if firstBlock != lastBlock {
		m.pushDown(level, lastBlock)
	}

	if firstBlock == lastBlock {
		blockStart := firstBlock << shift
		return m.queryRange(level-1, max(l, blockStart), min(r, blockStart+blockSize))
	}

	res := m.e()

	// Left partial
	leftEnd := (firstBlock + 1) << shift
	if leftEnd > l {
		res = m.op(res, m.queryRange(level-1, l, leftEnd))
	}

	// Middle full
	for bid := firstBlock + 1; bid < lastBlock; bid++ {
		res = m.op(res, m.levels[level][bid])
	}

	// Right partial
	rightStart := lastBlock << shift
	if rightStart < r {
		res = m.op(res, m.queryRange(level-1, rightStart, r))
	}

	return res
}

func (m *RadixTreeLazy[E, Id]) Get(i int) E {
	if i < 0 || i >= m.n {
		return m.e()
	}
	for level := len(m.levels) - 1; level >= 0; level-- {
		bid := i >> m.levelShifts[level]
		m.pushDown(level, bid)
	}
	return m.data[i]
}

func (m *RadixTreeLazy[E, Id]) GetAll() []E {
	for level := len(m.levels) - 1; level >= 0; level-- {
		for i := range m.levels[level] {
			m.pushDown(level, i)
		}
	}
	return m.data
}

func (m *RadixTreeLazy[E, Id]) Set(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	for level := len(m.levels) - 1; level >= 0; level-- {
		bid := i >> m.levelShifts[level]
		m.pushDown(level, bid)
	}
	m.data[i] = v
	for level := 0; level < len(m.levels); level++ {
		bid := i >> m.levelShifts[level]
		m.pushUp(level, bid)
	}
}

func (m *RadixTreeLazy[E, Id]) Update(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	for level := len(m.levels) - 1; level >= 0; level-- {
		bid := i >> m.levelShifts[level]
		m.pushDown(level, bid)
	}
	m.data[i] = m.op(m.data[i], v)
	for level := 0; level < len(m.levels); level++ {
		bid := i >> m.levelShifts[level]
		m.pushUp(level, bid)
	}
}

func (m *RadixTreeLazy[E, Id]) pushUp(level, bid int) {
	if level == 0 {
		return
	}

	shift := m.levelShifts[level] - m.levelShifts[level-1]
	start := bid << shift
	end := min(start+(1<<shift), len(m.levels[level-1]))

	v := m.e()
	for i := start; i < end; i++ {
		v = m.op(v, m.levels[level-1][i])
	}
	m.levels[level][bid] = v
}

func (m *RadixTreeLazy[E, Id]) pushDown(k, bid int) {
	if k >= len(m.levels) || m.lazy[k][bid] == m.id() {
		return
	}

	currentLazy := m.lazy[k][bid]
	m.lazy[k][bid] = m.id()

	if k == 0 {
		// Leaf level: apply to data
		start := bid << m.log
		end := min(start+m.blockSize, m.n)
		for i := start; i < end; i++ {
			m.data[i] = m.mapping(currentLazy, m.data[i])
		}
	} else {
		// Internal level: propagate to children
		childShift := m.levelShifts[k-1]
		childrenPerBlock := 1 << (m.levelShifts[k] - childShift)
		startChild := bid * childrenPerBlock
		endChild := min(startChild+childrenPerBlock, len(m.levels[k-1]))

		// Update children's values and lazy
		for i := startChild; i < endChild; i++ {
			m.levels[k-1][i] = m.mapping(currentLazy, m.levels[k-1][i])
			m.lazy[k-1][i] = m.composition(currentLazy, m.lazy[k-1][i])
		}
	}
}

func (m *RadixTreeLazy[E, Id]) propagate(k, bid int) {}

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
