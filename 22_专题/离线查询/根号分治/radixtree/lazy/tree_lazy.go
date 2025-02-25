package main

import (
	"fmt"
	"math/rand"
	"slices"
	"strings"
	"time"
)

func main() {
	testTime()
	for i := 0; i < 1000; i++ {
		test()
	}
	fmt.Println("pass")
}

// RadixTreeLazy 是一个支持区间修改和区间查询的基数树实现。
// 结合了RadixTree的查询能力和RadixTreeDual的修改能力。
type RadixTreeLazy[E any, Id comparable] struct {
	e           func() E                    // 聚合操作的幺元
	id          func() Id                   // 懒标记的幺元（恒等变换）
	op          func(a, b E) E              // 两个元素的聚合操作
	mapping     func(f Id, x E, size int) E // 将懒标记应用到节点值
	composition func(a, b Id) Id            // 两个懒标记的组合
	log         int                         // 控制每个块的大小 B=1<<log
	blockSize   int                         // 1 << log

	n           int    // 数组长度
	data        [][]E  // 数据层，data[0]是叶节点层
	lazy        [][]Id // 懒标记层
	layerShifts []int  // 每层块大小的移位量
}

// NewRadixTreeLazy 创建一个新的RadixTreeLazy实例
func NewRadixTreeLazy[E any, Id comparable](
	e func() E,
	id func() Id,
	op func(a, b E) E,
	mapping func(Id, E, int) E,
	composition func(a, b Id) Id,
	log int,
) *RadixTreeLazy[E, Id] {
	if log < 1 {
		log = 3
	}
	return &RadixTreeLazy[E, Id]{
		e:           e,
		op:          op,
		id:          id,
		composition: composition,
		mapping:     mapping,
		log:         log,
		blockSize:   1 << log,
	}
}

// Build 根据初始值构建基数树
func (rt *RadixTreeLazy[E, Id]) Build(n int, f func(i int) E) {
	rt.n = n

	// 初始化叶节点层
	level0 := make([]E, n)
	for i := 0; i < n; i++ {
		level0[i] = f(i)
	}
	rt.data = [][]E{level0}
	rt.layerShifts = []int{0}

	// 构建数据层
	build := func(pre []E) []E {
		cur := make([]E, (len(pre)+rt.blockSize-1)>>rt.log)
		for i := range cur {
			start := i << rt.log
			end := min(start+rt.blockSize, len(pre))
			v := pre[start]
			for j := start + 1; j < end; j++ {
				v = rt.op(v, pre[j])
			}
			cur[i] = v
		}
		return cur
	}

	preLevel := level0
	preShift := 1
	for len(preLevel) > 1 {
		curLevel := build(preLevel)
		rt.data = append(rt.data, curLevel)
		rt.layerShifts = append(rt.layerShifts, rt.log*preShift)
		preLevel = curLevel
		preShift++
	}

	// 初始化懒标记层
	rt.lazy = make([][]Id, len(rt.data))
	for k := 0; k < len(rt.data); k++ {
		rt.lazy[k] = make([]Id, len(rt.data[k]))
		for i := range rt.lazy[k] {
			rt.lazy[k][i] = rt.id()
		}
	}
}

// pushDown 将懒标记从上层下推到下层
func (rt *RadixTreeLazy[E, Id]) pushDown(k, blockIndex int) {
	if rt.lazy[k][blockIndex] == rt.id() {
		return
	}

	if k > 0 { // 非叶子层
		start := blockIndex << rt.log
		end := min(start+rt.blockSize, len(rt.data[k-1]))
		blockSize := 1 << rt.layerShifts[k-1]

		// 更新子块的值和懒标记
		for i := start; i < end; i++ {
			rt.data[k-1][i] = rt.mapping(rt.lazy[k][blockIndex], rt.data[k-1][i], blockSize)
			rt.lazy[k-1][i] = rt.composition(rt.lazy[k][blockIndex], rt.lazy[k-1][i])
		}
	} else { // 叶子层
		// 叶子节点直接应用懒标记，无需下推
		rt.data[0][blockIndex] = rt.mapping(rt.lazy[0][blockIndex], rt.data[0][blockIndex], 1)
	}

	// 重置当前懒标记
	rt.lazy[k][blockIndex] = rt.id()
}

// pushUp 更新上层节点值（从下往上）
func (rt *RadixTreeLazy[E, Id]) pushUp(k, blockIndex int) {
	if k <= 0 {
		return
	}

	start := blockIndex << rt.log
	end := min(start+rt.blockSize, len(rt.data[k-1]))

	// 确保下层懒标记已经下推
	for i := start; i < end; i++ {
		rt.pushDown(k-1, i)
	}

	// 计算新的聚合值
	v := rt.data[k-1][start]
	for i := start + 1; i < end; i++ {
		v = rt.op(v, rt.data[k-1][i])
	}
	rt.data[k][blockIndex] = v
}

// Get 获取单点值
func (rt *RadixTreeLazy[E, Id]) Get(index int) E {
	if index < 0 || index >= rt.n {
		return rt.e()
	}

	// 从上到下推懒标记
	for k := len(rt.layerShifts) - 1; k >= 0; k-- {
		blockIndex := index >> rt.layerShifts[k]
		rt.pushDown(k, blockIndex)
	}

	return rt.data[0][index]
}

// Set 设置单点值
func (rt *RadixTreeLazy[E, Id]) Set(index int, value E) {
	if index < 0 || index >= rt.n {
		return
	}

	// 从上到下推懒标记
	for k := len(rt.layerShifts) - 1; k >= 0; k-- {
		blockIndex := index >> rt.layerShifts[k]
		rt.pushDown(k, blockIndex)
	}

	// 设置值
	rt.data[0][index] = value

	// 从下到上更新聚合值
	for k := 1; k < len(rt.data); k++ {
		blockIndex := index >> rt.layerShifts[k]
		rt.pushUp(k, blockIndex)
	}
}

// Query 查询区间[l, r)的聚合值
func (rt *RadixTreeLazy[E, Id]) Query(l, r int) E {
	if l < 0 {
		l = 0
	}
	if r > rt.n {
		r = rt.n
	}
	if l >= r {
		return rt.e()
	}
	if l == 0 && r == rt.n {
		return rt.QueryAll()
	}

	return rt.query(l, r, len(rt.layerShifts)-1, 0, rt.n)
}

// query 内部查询实现
func (rt *RadixTreeLazy[E, Id]) query(l, r, k, start, end int) E {
	// 完全不相交
	if end <= l || r <= start {
		return rt.e()
	}

	blockIndex := start >> rt.layerShifts[k]

	// 完全包含
	if l <= start && end <= r {
		rt.pushDown(k, blockIndex)
		if k == 0 {
			return rt.data[0][start]
		}
		return rt.data[k][blockIndex]
	}

	// 需要下推后进一步查询
	rt.pushDown(k, blockIndex)

	if k == 0 {
		// 已经到叶节点层
		return rt.data[0][start]
	}

	// 分块查询
	mid := min(start+(1<<rt.layerShifts[k])>>1, end)
	leftVal := rt.query(l, r, k-1, start, mid)
	rightVal := rt.query(l, r, k-1, mid, end)
	return rt.op(leftVal, rightVal)
}

// Update 更新区间[l, r)的值
func (rt *RadixTreeLazy[E, Id]) Update(l, r int, value Id) {
	if l < 0 {
		l = 0
	}
	if r > rt.n {
		r = rt.n
	}
	if l >= r {
		return
	}

	rt.update(l, r, value, len(rt.layerShifts)-1, 0, rt.n)
}

// update 内部更新实现
func (rt *RadixTreeLazy[E, Id]) update(l, r int, value Id, k, start, end int) {
	// 完全不相交
	if end <= l || r <= start {
		return
	}

	blockIndex := start >> rt.layerShifts[k]

	// 完全包含
	if l <= start && end <= r {
		// 更新当前块的值和懒标记
		blockSize := 1 << rt.layerShifts[k]
		rt.data[k][blockIndex] = rt.mapping(value, rt.data[k][blockIndex], blockSize)
		rt.lazy[k][blockIndex] = rt.composition(value, rt.lazy[k][blockIndex])
		return
	}

	// 部分包含，需要先下推再更新子块
	rt.pushDown(k, blockIndex)

	if k == 0 {
		// 已经到叶节点层
		if l <= start && start < r {
			rt.data[0][start] = rt.mapping(value, rt.data[0][start], 1)
		}
		return
	}

	mid := min(start+(1<<rt.layerShifts[k])>>1, end)
	rt.update(l, r, value, k-1, start, mid)
	rt.update(l, r, value, k-1, mid, end)

	// 更新当前块的聚合值
	rt.pushUp(k, blockIndex)
}

// QueryAll 查询整个数组的聚合值
func (rt *RadixTreeLazy[E, Id]) QueryAll() E {
	if len(rt.data) == 0 {
		return rt.e()
	}
	// 根节点可能有懒标记
	topLevel := len(rt.data) - 1
	rt.pushDown(topLevel, 0)
	return rt.data[topLevel][0]
}

// GetAll 获取所有叶节点值
func (rt *RadixTreeLazy[E, Id]) GetAll() []E {
	// 把所有懒标记下推
	for k := len(rt.layerShifts) - 1; k >= 1; k-- {
		for i := 0; i < len(rt.data[k]); i++ {
			rt.pushDown(k, i)
		}
	}
	result := make([]E, rt.n)
	copy(result, rt.data[0])
	return result
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
				cand := rt.op(res, rt.data[k][blockIndex])
				if predicate(cand) {
					res = cand
					i += blockSize
					jumped = true
					break
				}
			}
		}
		if !jumped {
			res = rt.op(res, rt.data[0][i])
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
				cand := rt.op(rt.data[k][blockIndex], res)
				if predicate(cand) {
					res = cand
					i -= blockSize
					jumped = true
					break
				}
			}
		}
		if !jumped {
			res = rt.op(rt.data[0][i], res)
			if !predicate(res) {
				return i + 1
			}
			i--
		}
	}
	return 0
}

// String 返回树的字符串表示
func (rt *RadixTreeLazy[E, Id]) String() string {
	values := rt.GetAll()
	var sb strings.Builder
	sb.WriteString("[")
	for i, v := range values {
		if i > 0 {
			sb.WriteString(", ")
		}
		sb.WriteString(fmt.Sprintf("%v", v))
	}
	sb.WriteString("]")
	return sb.String()
}

// 工具函数
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

		rt1.Query(i, N)
		rt1.QueryAll()
		rt1.Get(i)
		rt1.Set(i, i)
		rt1.MaxRight(i, func(x int) bool { return x < int(1e18) })
		rt1.MinLeft(i, func(x int) bool { return x < int(1e18) })

	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
