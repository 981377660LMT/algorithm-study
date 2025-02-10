package main

type RadixTreeLazy[E any, Id comparable] struct {
	e           func() E
	id          func() Id
	op          func(a, b E) E    // 聚合操作
	mapping     func(f Id, x E) E // 应用标记到元素
	composition func(f, g Id) Id  // 合并标记
	log         int
	blockSize   int

	n           int
	data        []E
	levels      [][]E
	levelShifts []int
	lazy        [][]Id // 每层的懒标记
}

func NewRadixTreeLazy[E any, Id comparable](
	e func() E, id func() Id,
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
	shift := m.levelShifts[level]
	blockSize := 1 << shift
	firstBlock := l >> shift
	lastBlock := (r - 1) >> shift

	if level == -1 {
		for i := l; i < r; i++ {
			m.data[i] = m.mapping(f, m.data[i])
		}
		return
	}

	// Push down existing lazy
	m.pushDown(level, firstBlock)
	if firstBlock != lastBlock {
		m.pushDown(level, lastBlock)
	}

	if firstBlock == lastBlock {
		// Single block case
		blockStart := firstBlock << shift
		blockEnd := blockStart + blockSize
		innerL := max(l, blockStart)
		innerR := min(r, blockEnd)

		m.updateRange(k, level-1, innerL, innerR, f)
		m.pushUp(level, firstBlock)
	} else {
		// Update first partial block
		firstEnd := (firstBlock + 1) << shift
		if firstEnd > l {
			m.updateRange(k, level-1, l, firstEnd, f)
			m.pushUp(level, firstBlock)
		}

		// Update middle full blocks
		for bid := firstBlock + 1; bid < lastBlock; bid++ {
			m.levels[level][bid] = m.mapping(f, m.levels[level][bid])
			m.lazy[level][bid] = m.composition(f, m.lazy[level][bid])
		}

		// Update last partial block
		lastStart := lastBlock << shift
		if lastStart < r {
			m.updateRange(k, level-1, lastStart, r, f)
			m.pushUp(level, lastBlock)
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
	return m.queryRange(len(m.levels)-1, l, r)
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
