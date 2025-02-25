// !FIXME: UpdateRange 有问题，不要使用

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

// 2213. 由单个字符重复的最长子字符串
// https://leetcode.cn/problems/longest-substring-of-one-repeating-character/
// ![l,r]区间的最大连续长度就是
// !左区间的最大连续长度,右区间最大连续长度,以及左右两区间结合在一起中间的最大连续长度.
func longestRepeating(s string, queryCharacters string, queryIndices []int) []int {
	type E = struct {
		size                int
		preMax, sufMax, max int  // 前缀最大值,后缀最大值,区间最大值
		lc, rc              byte // 区间左端点字符,右端点字符
	}

	const INF int = 1e18
	e := func() E {
		return E{}
	}
	op := func(a, b E) E {
		res := E{lc: a.lc, rc: b.rc, size: a.size + b.size}
		if a.rc == b.lc {
			res.preMax = a.preMax
			if a.preMax == a.size {
				res.preMax += b.preMax
			}
			res.sufMax = b.sufMax
			if b.sufMax == b.size {
				res.sufMax += a.sufMax
			}
			res.max = max(max(a.max, b.max), a.sufMax+b.preMax)
		} else {
			res.preMax = a.preMax
			res.sufMax = b.sufMax
			res.max = max(a.max, b.max)
		}
		return res
	}

	n := len(s)
	seg := NewRadixTreeLazy(e, e, op, func(f E, x E, segSize int) E {
		return E{}
	}, func(f E, g E) E {
		return E{}
	}, -1)
	seg.Build(n, func(i int) E { return E{1, 1, 1, 1, s[i], s[i]} })
	res := make([]int, len(queryIndices))
	for i := 0; i < len(queryIndices); i++ {
		pos := queryIndices[i]
		char := queryCharacters[i]
		seg.Set(pos, E{1, 1, 1, 1, char, char})
		res[i] = seg.QueryRange(0, n).max
	}
	return res
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
				m.pushDown(k, blockIndex)
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
		m.pushDown(k, blockIndex)
	}
	return m.layers[0][i]
}

func (m *RadixTreeLazy[E, Id]) Set(i int, value E) {
	if i < 0 || i >= m.n {
		return
	}
	for k := len(m.layerShifts) - 1; k >= 1; k-- {
		blockIndex := i >> m.layerShifts[k]
		m.pushDown(k, blockIndex)
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

func (m *RadixTreeLazy[E, Id]) Update(i int, value E) {
	if i < 0 || i >= m.n {
		return
	}
	for k := len(m.layerShifts) - 1; k >= 1; k-- {
		blockIndex := i >> m.layerShifts[k]
		m.pushDown(k, blockIndex)
	}
	m.layers[0][i] = m.op(m.layers[0][i], value)
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

func (m *RadixTreeLazy[E, Id]) GetAll() []E {
	for k := len(m.layers) - 1; k >= 1; k-- {
		for i := 0; i < len(m.lazys[k]); i++ {
			m.pushDown(k, i)
		}
	}
	return m.layers[0]
}

// MaxRight 二分查询满足条件的最右位置
func (rt *RadixTreeLazy[E, Id]) MaxRight(left int, predicate func(E) bool) int {
	if left == rt.n {
		return rt.n
	}
	if left < 0 {
		left = 0
	}

	// 确保路径上的懒标记都已下推
	for k := len(rt.layerShifts) - 1; k >= 0; k-- {
		blockIndex := left >> rt.layerShifts[k]
		rt.pushDown(k, blockIndex)
	}

	res := rt.e()
	i := left
	for i < rt.n {
		jumped := false
		for k := len(rt.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << rt.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= rt.n {
				// 检查是否可以整块跳过
				blockIndex := i >> rt.layerShifts[k]
				rt.pushDown(k, blockIndex)
				cand := rt.op(res, rt.layers[k][blockIndex])
				if predicate(cand) {
					res = cand
					i += blockSize
					jumped = true
					break
				}
			}
		}
		if !jumped {
			res = rt.op(res, rt.layers[0][i])
			if !predicate(res) {
				return i
			}
			i++
		}
	}
	return rt.n
}

// MinLeft 二分查询满足条件的最左位置
func (rt *RadixTreeLazy[E, Id]) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	if right > rt.n {
		right = rt.n
	}

	// 确保路径上的懒标记都已下推
	for k := len(rt.layerShifts) - 1; k >= 0; k-- {
		blockIndex := (right - 1) >> rt.layerShifts[k]
		rt.pushDown(k, blockIndex)
	}

	res := rt.e()
	i := right - 1
	for i >= 0 {
		jumped := false
		for k := len(rt.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << rt.layerShifts[k]
			if (i+1)&(blockSize-1) == 0 && i+1-blockSize >= 0 {
				// 检查是否可以整块跳过
				blockIndex := (i + 1 - blockSize) >> rt.layerShifts[k]
				rt.pushDown(k, blockIndex)
				cand := rt.op(rt.layers[k][blockIndex], res)
				if predicate(cand) {
					res = cand
					i -= blockSize
					jumped = true
					break
				}
			}
		}
		if !jumped {
			res = rt.op(rt.layers[0][i], res)
			if !predicate(res) {
				return i + 1
			}
			i--
		}
	}
	return 0
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

func (m *naive[E, Id]) QueryRange(l, r int) E {
	result := m.e()
	for i := l; i < r; i++ {
		result = m.op(result, m.data[i])
	}
	return result
}

func (m *naive[E, Id]) QueryAll() E {
	result := m.e()
	for i := 0; i < m.n; i++ {
		result = m.op(result, m.data[i])
	}
	return result
}

func (m *naive[E, Id]) Get(i int) E {
	return m.data[i]
}

func (m *naive[E, Id]) GetAll() []E {
	return m.data
}

// 对单个位置 i 进行更新：用 mapping 处理 lazy 值 v
func (m *naive[E, Id]) Update(i int, v E) {
	m.data[i] = m.op(m.data[i], v)
}

// 对区间 [l, r) 内的每个位置更新：用 mapping 处理 lazy 值 v
func (m *naive[E, Id]) UpdateRange(l, r int, v E) {
	for i := l; i < r; i++ {
		m.data[i] = m.mapping(v, m.data[i], 1)
	}
}

func (m *naive[E, Id]) Set(i int, v E) {
	m.data[i] = v
}

func (m *naive[E, Id]) MaxRight(l int, f func(E) bool) int {
	sum := m.e()
	for i := l; i < m.n; i++ {
		sum = m.op(sum, m.data[i])
		if !f(sum) {
			return i
		}
	}
	return m.n
}

func (m *naive[E, Id]) MinLeft(r int, f func(E) bool) int {
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
	rt2 := newNaive(e, id, op, mapping, composition, -1)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(10)
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
			// Update
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Update(i, v)
			rt2.Update(i, v)
		case 2:
			// Get
			i := rand.Intn(N)
			if rt1.Get(i) != rt2.Get(i) {
				panic("err Get")
			}
		case 3:
			// GetAll
			nums1, nums2 := rt1.GetAll(), rt2.GetAll()
			if slices.Compare(nums1, nums2) != 0 {
				panic("err GetAll")
			}
		case 4:
			// Set
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Set(i, v)
			rt2.Set(i, v)
		case 5:
			// QueryRange
			l, r := rand.Intn(N), rand.Intn(N)
			if l > r {
				l, r = r, l
			}
			if rt1.QueryRange(l, r) != rt2.QueryRange(l, r) {
				panic("err QueryRange")
			}
		case 6:
			// QueryAll
			// if rt1.QueryAll() != rt2.QueryAll() {
			// 	panic("err QueryAll")
			// }
		case 7:
			// MaxRight
			i := rand.Intn(N)
			v := rand.Intn(100)
			if rt1.MaxRight(i, func(x int) bool { return x < v }) != rt2.MaxRight(i, func(x int) bool { return x < v }) {
				panic("err MaxRight")
			}
		case 8:
			// MinLeft
			i := rand.Intn(N)
			v := rand.Intn(100)
			if rt1.MinLeft(i, func(x int) bool { return x < v }) != rt2.MinLeft(i, func(x int) bool { return x < v }) {
				panic("err MinLeft")
			}
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
