// Dual 线段树的空间优化版本，支持区间修改、单点查询.
//
// - NewRadixTreeDual(e, op, log) 返回一个新的线段树对象，其中 e 是一个函数，返回一个单位元素，op 是一个函数，返回两个元素的组合结果，log 是线段树的块大小，如果 log < 1，则默认为 3.
// - (seg *RadixTreeDual) Build(n, f) 构建线段树，n 是线段树的长度，f 是一个函数，返回下标 i 处的元素.
// - (seg *RadixTreeDual) Get(index) 获取下标 index 处的元素.
// - (seg *RadixTreeDual) Set(index, value) 将下标 index 处的元素设置为 value.
// - (seg *RadixTreeDual) Update(start, end, value) 将区间 [start, end) 内的元素更新为 value.
// - (seg *RadixTreeDual) GetAll() 获取所有元素.

package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"os"
	"time"
)

func main() {
	// abc388_d()
	abc389_f()
}

func demo() {
	testTime()
	for i := 0; i < 1000; i++ {
		test()
	}
	fmt.Println("pass")
}

func abc388_d() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	seg := NewRadixTreeDual(func() int { return 0 }, func(a, b int) int { return a + b }, -1)
	seg.Build(n, func(i int) int { return 0 })
	for i := 0; i < n; i++ {
		nums[i] += seg.Get(i)
		cur := nums[i]
		k := min(cur, n-1-i)
		seg.Update(i+1, i+k+1, 1)
		nums[i] -= k
	}

	for i := 0; i < n; i++ {
		fmt.Fprint(out, nums[i], " ")
	}
}

func abc389_f() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const maxX int = 1 << 19

	var n int
	fmt.Fscan(in, &n)

	seg := NewRadixTreeDual(func() int { return 0 }, func(a, b int) int { return a + b }, -1)
	seg.Build(maxX, func(i int) int { return i })
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)

		a := MinLeft(maxX, func(x int) bool { return seg.Get(x) >= l }, 0)
		b := MinLeft(maxX, func(x int) bool { return seg.Get(x) > r }, a)
		seg.Update(a, b, 1)
	}

	res := seg.GetAll()

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var x int32
		fmt.Fscan(in, &x)
		fmt.Println(res[x])
	}
}

type RadixTreeDual[Id comparable] struct {
	id          func() Id
	composition func(a, b Id) Id
	log         int
	blockSize   int

	n           int
	layers      [][]Id // layers[0] 为叶层，存储最终值；layers[k] (k>=1) 存储懒更新（未下传）的值
	layerShifts []int
}

func NewRadixTreeDual[Id comparable](id func() Id, composition func(a, b Id) Id, log int) *RadixTreeDual[Id] {
	if log < 1 {
		log = 3
	}
	return &RadixTreeDual[Id]{
		id:          id,
		composition: composition,
		log:         log,
		blockSize:   1 << log,
	}
}

func (seg *RadixTreeDual[Id]) Build(n int, f func(i int) Id) {
	seg.n = n
	level0 := make([]Id, n)
	for i := 0; i < n; i++ {
		level0[i] = f(i)
	}
	seg.layers = [][]Id{level0}
	seg.layerShifts = []int{0}

	preLevel := level0
	shift := seg.log
	for len(preLevel) > 1 {
		sz := (len(preLevel) + seg.blockSize - 1) >> seg.log
		curLevel := make([]Id, sz)
		for i := 0; i < sz; i++ {
			curLevel[i] = seg.id()
		}
		seg.layers = append(seg.layers, curLevel)
		seg.layerShifts = append(seg.layerShifts, shift)
		preLevel = curLevel
		shift += seg.log
	}
}

func (seg *RadixTreeDual[Id]) propagate(k int, blockIndex int) {
	if seg.layers[k][blockIndex] != seg.id() {
		start := blockIndex << seg.log
		end := min(start+seg.blockSize, len(seg.layers[k-1]))
		for i := start; i < end; i++ {
			seg.layers[k-1][i] = seg.composition(seg.layers[k][blockIndex], seg.layers[k-1][i])
		}
		seg.layers[k][blockIndex] = seg.id()
	}
}

func (seg *RadixTreeDual[Id]) Get(index int) Id {
	if index < 0 || index >= seg.n {
		return seg.id()
	}
	for k := len(seg.layerShifts) - 1; k >= 1; k-- {
		blockIndex := index >> seg.layerShifts[k]
		seg.propagate(k, blockIndex)
	}
	return seg.layers[0][index]
}

func (seg *RadixTreeDual[Id]) Set(index int, value Id) {
	if index < 0 || index >= seg.n {
		return
	}
	for k := len(seg.layerShifts) - 1; k >= 1; k-- {
		blockIndex := index >> seg.layerShifts[k]
		seg.propagate(k, blockIndex)
	}
	seg.layers[0][index] = value
}

func (seg *RadixTreeDual[Id]) Update(start, end int, value Id) {
	if start < 0 {
		start = 0
	}
	if end > seg.n {
		end = seg.n
	}
	if start >= end {
		return
	}
	i := start
	for i < end {
		updated := false
		for k := len(seg.layerShifts) - 1; k >= 1; k-- {
			blockSize := 1 << seg.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= end {
				blockIndex := i >> seg.layerShifts[k]
				seg.layers[k][blockIndex] = seg.composition(value, seg.layers[k][blockIndex])
				i += blockSize
				updated = true
				break
			}
		}
		if !updated {
			seg.layers[0][i] = seg.composition(value, seg.layers[0][i])
			i++
		}
	}
}

func (seg *RadixTreeDual[Id]) GetAll() []Id {
	for k := len(seg.layerShifts) - 1; k >= 1; k-- {
		for i := 0; i < len(seg.layers[k]); i++ {
			seg.propagate(k, i)
		}
	}
	res := make([]Id, seg.n)
	copy(res, seg.layers[0])
	return res
}

func (seg *RadixTreeDual[Id]) String() string {
	return fmt.Sprintf("%v", seg.GetAll())
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

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
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
	e := func() int { return 0 }
	op := func(a, b int) int { return max(a, b) }
	N := rand.Intn(1000) + 1
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	rt1 := NewRadixTreeDual(e, op, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })
	rt2 := newNaive(e, op, -1)
	rt2.Build(N, func(i int) int { return randNums[i] })

	Q := int(1e4)
	for i := 0; i < Q; i++ {
		op := rand.Intn(5)
		switch op {
		case 0:
			l, r := rand.Intn(N), rand.Intn(N)
			if l > r {
				l, r = r, l
			}
			v := rand.Intn(100)
			rt1.Update(l, r, v)
			rt2.UpdateRange(l, r, v)
		case 1:
			i := rand.Intn(N)
			v := rand.Intn(100)
			rt1.Update(i, i+1, v)
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
			if len(nums1) != len(nums2) {
				panic("err GetAll")
			}
			for i := 0; i < len(nums1); i++ {
				if nums1[i] != nums2[i] {
					panic("err GetAll")
				}
			}

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
	op := func(a, b int) int { return a + b }

	N := int(2e5)
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	time1 := time.Now()
	rt1 := NewRadixTreeDual(e, op, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })

	for i := 0; i < N; i++ {

		rt1.Update(i, N, i)
		rt1.Get(i)
		rt1.Set(i, i)

	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
