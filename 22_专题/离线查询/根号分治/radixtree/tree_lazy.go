// Lazy 线段树的空间优化版本，支持区间修改、区间查询.

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

// RadixTreeLazy 实现基于多级分块结构的 lazy 区间更新与聚合查询。
// 用户需提供 e、op、mapping、composition、id 这五个函数。
type RadixTreeLazy[E any, Id comparable] struct {
	n         int // 叶层长度
	log       int // 每块对数大小（即每块大小 = 1<<log）
	blockSize int // = 1<<log

	e           func() E
	op          func(a, b E) E
	mapping     func(f Id, x E, segLen int) E
	composition func(f, g Id) Id
	id          func() Id

	// 多级分块，layers[0] 为叶层，其它层每个元素存储下层一整块的聚合值
	layers [][]E
	// 每层（除叶层）存储 lazy 更新信息，对应层上每个块的待下传更新
	lazyLayers [][]Id
	// layerShifts[k] 表示层 k 中每块覆盖叶层的长度为 1 << layerShifts[k]（层0为0，层1为log，层2为2*log，…）
	layerShifts []int
}

// NewRadixTreeLazy 构造一个 RadixTreeLazy 实例，若 log < 1 则取默认值 3。
func NewRadixTreeLazy[E any, Id comparable](
	e func() E,
	op func(a, b E) E,
	mapping func(f Id, x E, segLen int) E,
	composition func(f, g Id) Id,
	id func() Id,
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

// Build 根据 n 及 f(i) 构造叶层，并自底向上构造各层聚合块，同时各层 lazy 值全部初始化为 id().
func (rt *RadixTreeLazy[E, Id]) Build(n int, f func(i int) E) {
	rt.n = n
	// 构造叶层
	level0 := make([]E, n)
	for i := 0; i < n; i++ {
		level0[i] = f(i)
	}
	rt.layers = [][]E{level0}
	rt.layerShifts = []int{0}
	// 叶层不维护 lazy 数组（lazy 只对内部节点有效），故设置 nil
	rt.lazyLayers = [][]Id{nil}

	// 自底向上构造内部层，每层每个块聚合下层 blockSize 个元素
	pre := level0
	shift := rt.log
	for len(pre) > 1 {
		size := (len(pre) + rt.blockSize - 1) >> rt.log
		cur := make([]E, size)
		lazyCur := make([]Id, size)
		for i := 0; i < size; i++ {
			start := i << rt.log
			end := start + rt.blockSize
			if end > len(pre) {
				end = len(pre)
			}
			val := pre[start]
			for j := start + 1; j < end; j++ {
				val = rt.op(val, pre[j])
			}
			cur[i] = val
			lazyCur[i] = rt.id()
		}
		rt.layers = append(rt.layers, cur)
		rt.lazyLayers = append(rt.lazyLayers, lazyCur)
		rt.layerShifts = append(rt.layerShifts, shift)
		pre = cur
		shift += rt.log
	}
}

// propagate 将层 k 中下标 b 处的 lazy 更新下传到下一层 (k-1) 中对应的块
func (rt *RadixTreeLazy[E, Id]) propagate(k int, b int) {
	if rt.lazyLayers[k][b] != rt.id() {
		// 层 k 中第 b 个块覆盖层 (k-1) 中下标 [start, end)
		start := b << rt.log
		end := start + rt.blockSize
		if end > len(rt.layers[k-1]) {
			end = len(rt.layers[k-1])
		}
		for i := start; i < end; i++ {
			// 计算 i 对应块在层 (k-1) 覆盖叶层的长度
			segLen := 1 << rt.layerShifts[k-1]
			// 对于最后一个块可能不足 segLen
			if (i+1)<<rt.layerShifts[k-1] > rt.n {
				segLen = rt.n - (i << rt.layerShifts[k-1])
			}
			rt.layers[k-1][i] = rt.mapping(rt.lazyLayers[k][b], rt.layers[k-1][i], segLen)
			// 若下层也维护 lazy，则更新之
			if k-1 > 0 {
				rt.lazyLayers[k-1][i] = rt.composition(rt.lazyLayers[k][b], rt.lazyLayers[k-1][i])
			}
		}
		rt.lazyLayers[k][b] = rt.id()
	}
}

// Get 返回下标 i 处的值（叶层），查询前沿着路径自顶向下传播 lazy 更新
func (rt *RadixTreeLazy[E, Id]) Get(i int) E {
	if i < 0 || i >= rt.n {
		return rt.e()
	}
	for k := len(rt.layerShifts) - 1; k >= 1; k-- {
		b := i >> rt.layerShifts[k]
		rt.propagate(k, b)
	}
	return rt.layers[0][i]
}

// Set 将下标 i 处的值设为 v（覆盖之前所有更新），先传播 lazy 更新，再更新叶层及沿途父节点
func (rt *RadixTreeLazy[E, Id]) Set(i int, v E) {
	if i < 0 || i >= rt.n {
		return
	}
	for k := len(rt.layerShifts) - 1; k >= 1; k-- {
		b := i >> rt.layerShifts[k]
		rt.propagate(k, b)
	}
	rt.layers[0][i] = v
	pre := rt.layers[0]
	for k := 1; k < len(rt.layers); k++ {
		b := i >> rt.layerShifts[k]
		start := b << rt.log
		end := start + rt.blockSize
		if end > len(pre) {
			end = len(pre)
		}
		val := pre[start]
		for j := start + 1; j < end; j++ {
			val = rt.op(val, pre[j])
		}
		rt.layers[k][b] = val
		pre = rt.layers[k]
	}
}

// Query 返回区间 [l, r) 内的聚合值。
// 查询时采用类似 RadixTree 的跳跃方式，对于整块直接取各层聚合值，对于边界部分则逐点查询 Get。
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
	// 如果查询整个区间，可直接返回最高层唯一值
	if l == 0 && r == rt.n {
		return rt.QueryAll()
	}
	res := rt.e()
	i := l
	for i < r {
		jumped := false
		// 尝试利用最高层整块跳跃
		for k := len(rt.layerShifts) - 1; k >= 0; k-- {
			blockSize := 1 << rt.layerShifts[k]
			if i&(blockSize-1) == 0 && i+blockSize <= r {
				b := i >> rt.layerShifts[k]
				if k > 0 {
					rt.propagate(k, b)
				}
				res = rt.op(res, rt.layers[k][b])
				i += blockSize
				jumped = true
				break
			}
		}
		if !jumped {
			res = rt.op(res, rt.Get(i))
			i++
		}
	}
	return res
}

// Update 对区间 [l, r) 施加 lazy 更新 f，更新时尽可能利用整块更新。
// 更新操作为：对于每个受影响区间，其新值 = mapping(f, 原值, 区间长度)；同时 lazy 值合并为 composition(f, 原 lazy)。
func (rt *RadixTreeLazy[E, Id]) Update(l, r int, f Id) {
	if l < 0 {
		l = 0
	}
	if r > rt.n {
		r = rt.n
	}
	if l >= r {
		return
	}
	i := l
	for i < r {
		updated := false
		for k := len(rt.layerShifts) - 1; k >= 1; k-- {
			blockSize := 1 << rt.layerShifts[k]
			// 若 i 对齐且整个块在 [l, r) 内，则在该层直接 lazy 更新
			if i&(blockSize-1) == 0 && i+blockSize <= r {
				b := i >> rt.layerShifts[k]
				rt.lazyLayers[k][b] = rt.composition(f, rt.lazyLayers[k][b])
				segLen := blockSize
				if i+blockSize > rt.n {
					segLen = rt.n - i
				}
				rt.layers[k][b] = rt.mapping(f, rt.layers[k][b], segLen)
				i += blockSize
				updated = true
				break
			}
		}
		if !updated {
			// 更新叶层单个元素
			rt.layers[0][i] = rt.mapping(f, rt.layers[0][i], 1)
			i++
		}
	}
	// 对于受部分更新的块，重新聚合下层数据（这里简单起见，对所有内部层进行一次重算）
	for k := 1; k < len(rt.layers); k++ {
		for b := 0; b < len(rt.layers[k]); b++ {
			// 若该块存在 pending lazy，则不更新（其值已包含 lazy 效果）
			if rt.lazyLayers[k][b] != rt.id() {
				continue
			}
			start := b << rt.log
			end := start + rt.blockSize
			if end > len(rt.layers[k-1]) {
				end = len(rt.layers[k-1])
			}
			val := rt.layers[k-1][start]
			for j := start + 1; j < end; j++ {
				val = rt.op(val, rt.layers[k-1][j])
			}
			rt.layers[k][b] = val
		}
	}
}

// QueryAll 返回整个区间 [0, n) 的聚合值，即最高层唯一值
func (rt *RadixTreeLazy[E, Id]) QueryAll() E {
	if len(rt.layers) == 0 {
		return rt.e()
	}
	return rt.layers[len(rt.layers)-1][0]
}

// GetAll 返回叶层所有值，更新前先将各层 lazy 全部下传
func (rt *RadixTreeLazy[E, Id]) GetAll() []E {
	for k := len(rt.layerShifts) - 1; k >= 1; k-- {
		for b := 0; b < len(rt.lazyLayers[k]); b++ {
			rt.propagate(k, b)
		}
	}
	return rt.layers[0]
}

func (rt *RadixTreeLazy[E, Id]) String() string {
	return fmt.Sprint(rt.GetAll())
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

	rt1 := NewRadixTreeLazy(e, op, mapping, composition, id, -1)
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
			// i := rand.Intn(N)
			// v := rand.Intn(100)
			// rt1.Update(i, i+1, v)
			// rt2.Update(i, v)
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
	op := func(a, b int) int { return a + b }

	N := int(2e5)
	randNums := make([]int, N)
	for i := 0; i < N; i++ {
		randNums[i] = rand.Intn(1000)
	}

	time1 := time.Now()
	rt1 := NewRadixTreeLazy(e, op, func(f int, x int, segLen int) int { return f }, func(f, g int) int { return g }, func() int { return 0 }, -1)
	rt1.Build(N, func(i int) int { return randNums[i] })

	for i := 0; i < N; i++ {

		rt1.Update(i, N, i)
		rt1.Get(i)
		rt1.Set(i, i)

	}

	time2 := time.Now()
	fmt.Println("RadixTree:", time2.Sub(time1))
}
