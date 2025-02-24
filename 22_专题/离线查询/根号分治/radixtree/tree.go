// 线段树的空间优化版本，支持区间查询、单点修改:
//
// - NewRadixTree(e, op, log): 构造一个新的线段树，e 是幺元，op 是结合律，log 是每个块的大小。
// - Build(n, f): 构建长度为 n 的线段树，f 是初始化函数。
// - QueryRange(l, r): 查询 [l,r) 区间的聚合值。
// - QueryAll(): 查询所有元素的聚合值。
// - Get(i): 查询第 i 个元素的值。
// - GetAll(): 查询所有元素的值。
// - Update(i, v): 将第 i 个元素更新为 v。
// - Set(i, v): 将第 i 个元素设置为 v。

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
	seg := NewRadixTree(e, op, -1)
	seg.Build(n, func(i int) E { return E{1, 1, 1, 1, s[i], s[i]} })
	res := make([]int, len(queryIndices))
	for i := 0; i < len(queryIndices); i++ {
		pos := queryIndices[i]
		char := queryCharacters[i]
		seg.Set(pos, E{1, 1, 1, 1, char, char})
		res[i] = seg.QueryAll().max
	}
	return res
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
	layers      [][]E // layers[k][i] 表示第k层第i个块的聚合值.
	layerShifts []int // layers[k] 表示第k层的块大小为1<<layerShifts[k].
}

// log: 每个块的大小B=1<<log.默认log=3.
// e: 幺元.
// op: 结合律.
func NewRadixTree[E any](e func() E, op func(a, b E) E, log int) *RadixTree[E] {
	if log < 1 {
		log = 3
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
	level0 := make([]E, n)
	for i := 0; i < n; i++ {
		level0[i] = f(i)
	}
	m.layers = [][]E{level0}
	m.layerShifts = []int{0}

	build := func(pre []E) []E {
		cur := make([]E, (len(pre)+m.blockSize-1)>>m.log)
		for i := range cur {
			start := i << m.log
			end := min(start+m.blockSize, len(pre))
			v := pre[start]
			for j := start + 1; j < end; j++ {
				v = m.op(v, pre[j])
			}
			cur[i] = v
		}
		return cur
	}

	preLevel := level0
	preShift := 1
	for len(preLevel) > 1 {
		curLevel := build(preLevel)
		m.layers = append(m.layers, curLevel)
		m.layerShifts = append(m.layerShifts, m.log*preShift)
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
	if l == 0 && r == m.n {
		return m.QueryAll()
	}
	res := m.e()
	i := l
	for i < r {
		jumped := false
		// 从最高层开始尝试找到最大的对齐块
		for k := len(m.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << m.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= r {
				res = m.op(res, m.layers[k][i>>m.layerShifts[k]])
				i += blockSize
				jumped = true
				break
			}
		}
		if !jumped {
			res = m.op(res, m.layers[0][i])
			i++
		}
	}
	return res
}

func (m *RadixTree[E]) QueryAll() E {
	if len(m.layers) == 0 {
		return m.e()
	}
	return m.layers[len(m.layers)-1][0]
}

func (m *RadixTree[E]) Get(i int) E {
	if i < 0 || i >= m.n {
		return m.e()
	}
	return m.layers[0][i]
}

// O(1).
func (m *RadixTree[E]) GetAll() []E {
	return m.layers[0]
}

// A[i] = op(A[i], v).
func (m *RadixTree[E]) Update(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.layers[0][i] = m.op(m.layers[0][i], v)
	pre := m.layers[0]
	for k := 1; k < len(m.layers); k++ {
		bid := i >> m.layerShifts[k]
		start := bid << m.log
		end := min(start+m.blockSize, len(pre))
		cur := pre[start]
		for j := start + 1; j < end; j++ {
			cur = m.op(cur, pre[j])
		}
		m.layers[k][bid] = cur
		pre = m.layers[k]
	}
}

// A[i] = v.
func (m *RadixTree[E]) Set(i int, v E) {
	if i < 0 || i >= m.n {
		return
	}
	m.layers[0][i] = v
	pre := m.layers[0]
	for k := 1; k < len(m.layers); k++ {
		bid := i >> m.layerShifts[k]
		start := bid << m.log
		end := min(start+m.blockSize, len(pre))
		cur := pre[start]
		for j := start + 1; j < end; j++ {
			cur = m.op(cur, pre[j])
		}
		m.layers[k][bid] = cur
		pre = m.layers[k]
	}
}

// 返回最大的 r，使得区间 [l, r) 的聚合值满足 f.
// O(logN).
func (rt *RadixTree[E]) MaxRight(left int, predicate func(E) bool) int {
	if left == rt.n {
		return rt.n
	}
	res := rt.e()
	i := left
	for i < rt.n {
		jumped := false
		// 尝试利用整块跳跃
		for k := len(rt.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << rt.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= rt.n {
				cand := rt.op(res, rt.layers[k][i>>rt.layerShifts[k]])
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

// MinLeft 返回最小的 l，使得区间 [l, r) 的聚合值满足 f.
// O(logN).
func (rt *RadixTree[E]) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	res := rt.e()
	i := right - 1
	for i >= 0 {
		jumped := false
		for k := len(rt.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << rt.layerShifts[k]
			// 判断当前下标是否正好为某块的最右端
			if (i+1)&(blockSize-1) == 0 && i+1-blockSize >= 0 {
				cand := rt.op(rt.layers[k][(i+1-blockSize)>>rt.layerShifts[k]], res)
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
	e := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	N := rand.Intn(1000) + 1
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	rt1 := NewRadixTree(e, op, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })
	rt2 := newNaive(e, op, -1)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(9)
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
			f := func(x int) bool { return x < 500 }
			if rt1.MaxRight(l, f) != rt2.MaxRight(l, f) {
				panic("err MaxRight")
			}
		case 8:
			// MinLeft
			r := rand.Intn(N)
			f := func(x int) bool { return x < 500 }
			if rt1.MinLeft(r, f) != rt2.MinLeft(r, f) {
				panic("err MinLeft")
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
	rt1 := NewRadixTree(e, op, 3)
	rt1.Build(N, func(i int) int { return randNums[i] })

	for i := 0; i < N; i++ {
		rt1.QueryRange(i, N)
		rt1.QueryAll()
		rt1.Get(i)
		rt1.Set(i, i)
		rt1.MaxRight(i, func(x int) bool { return x < int(1e18) })
		rt1.MinLeft(i, func(x int) bool { return x < int(1e18) })
	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
