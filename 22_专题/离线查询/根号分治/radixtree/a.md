参考 RadixTree、RadixTreeDual、SegmentTreeLazy 代码，实现一个
区间修改、区间查询的 RadixTreeLazy 版本.

```go
// RadixTree.go

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
```

```go
// RadixTreeDual.go

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
```

```go
// SegmentTreeLazy

// !template
type LazySegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewLazySegTree(n int, f func(int) E) *LazySegTree {
	tree := &LazySegTree{}
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewLazySegTreeFrom(leaves []E) *LazySegTree {
	tree := &LazySegTree{}
	n := len(leaves)
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree) GetAll() []E {
	for i := 1; i < tree.size; i++ {
		tree.pushDown(i)
	}
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if tmp := tree.op(tree.data[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if tmp := tree.op(res, tree.data[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree) propagate(root int, f Id) {
	size := 1 << (tree.log - (bits.Len32(uint32(root)) - 1) /**topbit**/)
	tree.data[root] = tree.mapping(f, tree.data[root], size)
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *LazySegTree) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}
```
