package main

import (
	"fmt"
)

type RadixTree[E any] struct {
	e          func() E
	op         func(a, b E) E
	log        int
	blockSize  int
	data       []E
	levels     [][]E
	levelSizes []int
	n          int
}

func NewRadixTree[E any](e func() E, op func(a, b E) E, log int) *RadixTree[E] {
	if log < 1 {
		log = 1
	} else if log > 16 {
		log = 16
	}
	return &RadixTree[E]{
		e:   e,
		op:  op,
		log: log,
	}
}

func (m *RadixTree[E]) Build(n int, f func(i int) E) {
	m.n = n
	m.data = make([]E, n)
	for i := 0; i < n; i++ {
		m.data[i] = f(i)
	}
	m.blockSize = 1 << m.log
	m.levels = [][]E{}
	m.levelSizes = []int{}

	// Build level 0
	level0Blocks := (n + m.blockSize - 1) >> m.log
	level0 := make([]E, level0Blocks)
	for i := range level0 {
		start := i * m.blockSize
		end := min(start+m.blockSize, n)
		val := m.e()
		for j := start; j < end; j++ {
			val = m.op(val, m.data[j])
		}
		level0[i] = val
	}
	m.levels = append(m.levels, level0)
	m.levelSizes = append(m.levelSizes, m.blockSize)

	// Build higher levels
	for len(m.levels[len(m.levels)-1]) > 1 {
		prevLevel := m.levels[len(m.levels)-1]
		prevLevelBlocks := len(prevLevel)
		currentLevelBlocks := (prevLevelBlocks + m.blockSize - 1) >> m.log
		currentLevel := make([]E, currentLevelBlocks)
		for i := range currentLevel {
			start := i * m.blockSize
			end := min(start+m.blockSize, prevLevelBlocks)
			val := m.e()
			for j := start; j < end; j++ {
				val = m.op(val, prevLevel[j])
			}
			currentLevel[i] = val
		}
		m.levels = append(m.levels, currentLevel)
		m.levelSizes = append(m.levelSizes, m.levelSizes[len(m.levelSizes)-1]<<m.log)
	}
}

func (m *RadixTree[E]) QueryRange(l, r int) E {
	if l >= r || l < 0 || r > m.n {
		return m.e()
	}
	res := m.e()
	type task struct{ l, r, k int }
	tasks := []task{{l, r, len(m.levels) - 1}}

	for len(tasks) > 0 {
		t := tasks[len(tasks)-1]
		tasks = tasks[:len(tasks)-1]
		l, r, k := t.l, t.r, t.k

		if l >= r {
			continue
		}

		if k < 0 {
			for i := l; i < r; i++ {
				res = m.op(res, m.data[i])
			}
			continue
		}

		s := m.levelSizes[k]
		startBlock, endBlock := l/s, (r-1)/s

		if startBlock > endBlock {
			tasks = append(tasks, task{l, r, k - 1})
			continue
		}

		if startBlock == endBlock {
			tasks = append(tasks, task{l, r, k - 1})
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

		leftEnd := (startBlock + 1) * s
		if leftEnd > r {
			leftEnd = r
		}
		if leftEnd > l {
			tasks = append(tasks, task{l, leftEnd, k - 1})
		}

		rightStart := endBlock * s
		if rightStart < l {
			rightStart = l
		}
		if rightStart < r {
			tasks = append(tasks, task{rightStart, r, k - 1})
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

func (m *RadixTree[E]) GetAll() []E {
	res := make([]E, len(m.data))
	copy(res, m.data)
	return res
}

func (m *RadixTree[E]) Update(i int, v E) {
	newVal := m.op(m.Get(i), v)
	m.Set(i, newVal)
}

func (m *RadixTree[E]) Set(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.data[i] = v
	for k := 0; k < len(m.levels); k++ {
		s := m.levelSizes[k]
		blockIdx := i / s
		start := blockIdx * s
		end := min(start+s, m.n)

		var val E
		if k == 0 {
			val = m.e()
			for j := start; j < end; j++ {
				val = m.op(val, m.data[j])
			}
		} else {
			prevLevel := m.levels[k-1]
			prevStart := blockIdx * m.blockSize
			prevEnd := min(prevStart+m.blockSize, len(prevLevel))
			val = m.e()
			for j := prevStart; j < prevEnd; j++ {
				val = m.op(val, prevLevel[j])
			}
		}
		if blockIdx < len(m.levels[k]) {
			m.levels[k][blockIdx] = val
		}
	}
}

func (m *RadixTree[E]) MaxRight(left int, predicate func(E) bool) int {
	if left == m.n {
		return m.n
	}
	current := m.e()
	if !predicate(current) {
		return left
	}

	pos := left
	for k := 0; k < len(m.levels); k++ {
		s := m.levelSizes[k]
		if pos%s == 0 && pos+s <= m.n {
			blockIdx := pos / s
			if blockIdx >= len(m.levels[k]) {
				continue
			}
			nextVal := m.op(current, m.levels[k][blockIdx])
			if predicate(nextVal) {
				current = nextVal
				pos += s
				k = -1 // Restart from level 0
			}
		}
	}

	for pos < m.n {
		nextVal := m.op(current, m.Get(pos))
		if !predicate(nextVal) {
			break
		}
		current = nextVal
		pos++
	}
	return pos
}

func (m *RadixTree[E]) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	current := m.e()
	if !predicate(current) {
		return right
	}

	pos := right
	for k := 0; k < len(m.levels); k++ {
		s := m.levelSizes[k]
		if pos%s == 0 && pos-s >= 0 {
			blockIdx := (pos - 1) / s
			if blockIdx >= len(m.levels[k]) {
				continue
			}
			nextVal := m.op(m.levels[k][blockIdx], current)
			if predicate(nextVal) {
				current = nextVal
				pos -= s
				k = -1 // Restart from level 0
			}
		}
	}

	for pos > 0 {
		nextVal := m.op(m.Get(pos-1), current)
		if !predicate(nextVal) {
			break
		}
		current = nextVal
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

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }

	// Test RadixTree
	rt := NewRadixTree(e, op, 3)
	rt.Build(10, func(i int) int { return i + 1 })
	fmt.Println("RadixTree QueryRange [2,5):", rt.QueryRange(2, 5)) // 3+4+5=12
	fmt.Println("RadixTree GetAll:", rt.GetAll())                   // [1 2 3 4 5 6 7 8 9 10]

	rt.Update(3, 5)
	fmt.Println("After update index 3:", rt.GetAll()) // 3 becomes 8 (3+5)

	rt.Set(6, 10)
	fmt.Println("After set index 6 to 10:", rt.GetAll())

	fmt.Println("Query range [2,5):", rt.QueryRange(2, 5)) // 3+8+5=16

	rt.Build(10, func(i int) int { return i + 1 })
	fmt.Println("MinLeft:", rt.MinLeft(10, func(x int) bool { return x > 5 }))  // Expected 5
	fmt.Println("MaxRight:", rt.MaxRight(0, func(x int) bool { return x < 5 })) // Expected 4
}
