// 动态最小生成树(边权值为树上两点距离)
// 动态加、删点，求连接点集的最小生成树的点的个数.
//
//    0
//   / \
//  1   2
//     / \
//    3   4
//       / \
//      5   6

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

func main() {
	edges := [][]int32{{0, 1}, {0, 2}, {2, 3}, {2, 4}, {4, 5}, {4, 6}}
	n := int32(len(edges) + 1)
	tree := make([][]int32, n)
	for _, e := range edges {
		tree[e[0]] = append(tree[e[0]], e[1])
		tree[e[1]] = append(tree[e[1]], e[0])
	}
	mst := NewMinSpanningGraphOnTree(n)
	mst.Init(tree, 0)
	mst.Add(0)
	mst.Add(3)
	mst.Add(5)
	fmt.Println(mst.SizeOfMinimumSpanningGraph()) // 5
	mst.Remove(3)
	fmt.Println(mst.SizeOfMinimumSpanningGraph()) // 4
}

type MinSpanningGraphOnTree struct {
	tree      [][]int32
	depth     []int32
	lca       *lcaOnTreeBySchieberVishkin
	index2dfn []int32
	dfn2index []int32
	order     int32
	ss        *SortedSet32
	size      int32
}

func NewMinSpanningGraphOnTree(n int32) *MinSpanningGraphOnTree {
	res := &MinSpanningGraphOnTree{
		depth:     make([]int32, n),
		index2dfn: make([]int32, n),
		dfn2index: make([]int32, n),
		lca:       newLcaOnTreeBySchieberVishkin(n),
	}
	ss := NewSortedSet32(func(a, b int32) bool { return res.index2dfn[a] < res.index2dfn[b] })
	res.ss = ss
	return res
}

func (m *MinSpanningGraphOnTree) Init(tree [][]int32, root int32) {
	m.tree = tree
	m.lca.Init(tree, root)
	m.order = 0
	m.dfs(root, -1)
	m.Clear()
}

func (m *MinSpanningGraphOnTree) SizeOfMinimumSpanningGraph() int32 {
	if m.Size() == 0 {
		return 0
	}
	res := m.size + m.Dist(m.ss.Min(), m.ss.Max())
	return res/2 + 1
}

// 加入前必须保证v不在集合中.
func (m *MinSpanningGraphOnTree) Add(v int32) {
	floor, ok1 := m.ss.Floor(v)
	if ok1 {
		m.size += m.Dist(floor, v)
	}
	ceil, ok2 := m.ss.Ceiling(v)
	if ok2 {
		m.size += m.Dist(ceil, v)
	}
	if ok1 && ok2 {
		m.size -= m.Dist(floor, ceil)
	}
	m.ss.Add(v)
}

// 移除前必须保证v在集合中.
func (m *MinSpanningGraphOnTree) Remove(v int32) {
	m.ss.Discard(v)
	floor, ok1 := m.ss.Floor(v)
	if ok1 {
		m.size -= m.Dist(floor, v)
	}
	ceil, ok2 := m.ss.Ceiling(v)
	if ok2 {
		m.size -= m.Dist(ceil, v)
	}
	if ok1 && ok2 {
		m.size += m.Dist(floor, ceil)
	}
}

func (m *MinSpanningGraphOnTree) Has(v int32) bool {
	return m.ss.Has(v)
}

func (m *MinSpanningGraphOnTree) Size() int32 {
	return m.ss.Len()
}

func (m *MinSpanningGraphOnTree) Dist(a, b int32) int32 {
	c := m.lca.Lca(a, b)
	return m.depth[a] + m.depth[b] - m.depth[c]*2
}

func (m *MinSpanningGraphOnTree) Clear() {
	m.ss.Clear()
	m.size = 0
}

func (m *MinSpanningGraphOnTree) dfs(cur, pre int32) {
	m.index2dfn[cur] = m.order
	m.dfn2index[m.order] = cur
	m.order++
	if pre == -1 {
		m.depth[cur] = 0
	} else {
		m.depth[cur] = m.depth[pre] + 1
	}
	for _, to := range m.tree[cur] {
		if to != pre {
			m.dfs(to, cur)
		}
	}
}

// O(n)时空间预处理，O(1)查询LCA。
type lcaOnTreeBySchieberVishkin struct {
	preOrder []int32
	i        []int32
	head     []int32
	a        []int32
	parent   []int32
	time     int32
}

func newLcaOnTreeBySchieberVishkin(n int32) *lcaOnTreeBySchieberVishkin {
	res := &lcaOnTreeBySchieberVishkin{
		preOrder: make([]int32, n),
		i:        make([]int32, n),
		head:     make([]int32, n),
		a:        make([]int32, n),
		parent:   make([]int32, n),
	}
	return res
}

func (l *lcaOnTreeBySchieberVishkin) Init(tree [][]int32, root int32) {
	l.time = 0
	l._dfs1(tree, root, -1)
	l._dfs2(tree, root, -1, 0)
}

func (l *lcaOnTreeBySchieberVishkin) InitWithIsRoot(tree [][]int32, isRoot func(i int32) bool) {
	l.time = 0
	for i := int32(0); i < int32(len(tree)); i++ {
		if isRoot(i) {
			l._dfs1(tree, i, -1)
			l._dfs2(tree, i, -1, 0)
		}
	}
}

// floorLog: bits.Len32(uint32(n)) - 1
func (l *lcaOnTreeBySchieberVishkin) Lca(x, y int32) int32 {
	var hb int32
	if a, b := l.i[x], l.i[y]; a == b {
		hb = a & -a
	} else {
		hb = 1 << (bits.Len32(uint32(a^b)) - 1)
	}
	tmp := l.a[x] & l.a[y] & -hb
	hz := tmp & -tmp
	ex := l._enterIntoStrip(x, hz)
	ey := l._enterIntoStrip(y, hz)
	if l.preOrder[ex] < l.preOrder[ey] {
		return ex
	} else {
		return ey
	}
}

func (l *lcaOnTreeBySchieberVishkin) _dfs1(tree [][]int32, u, p int32) {
	l.parent[u] = p
	l.i[u] = l.time
	l.preOrder[u] = l.time
	l.time++
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		l._dfs1(tree, v, u)
		if a, b := l.i[u], l.i[v]; a&-a < b&-b {
			l.i[u] = b
		}
	}
	l.head[l.i[u]] = u
}

func (l *lcaOnTreeBySchieberVishkin) _dfs2(tree [][]int32, u, p, up int32) {
	l.a[u] = up | l.i[u]&-l.i[u]
	for _, v := range tree[u] {
		if v == p {
			continue
		}
		l._dfs2(tree, v, u, l.a[u])
	}
}

func (l *lcaOnTreeBySchieberVishkin) _enterIntoStrip(x, hz int32) int32 {
	if a := l.i[x]; a&-a == hz {
		return x
	}
	tmp := l.a[x] & (hz - 1)
	hw := int32(1 << (bits.Len32(uint32(tmp)) - 1))
	return l.parent[l.head[l.i[x]&-hw|hw]]
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int32 = 200

type S = int32

// 使用分块+树状数组维护的有序序列.
type SortedSet32 struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSortedSet32(less func(a, b S) bool, elements ...S) *SortedSet32 {
	elements = append(elements[:0:0], elements...)
	res := &SortedSet32{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := int32(len(elements))
	blocks := [][]S{}
	for start := int32(0); start < n; start += _LOAD {
		end := min32(start+_LOAD, n)
		blocks = append(blocks, elements[start:end:end]) // !各个块互不影响, max参数也需要指定为end
	}
	mins := make([]S, len(blocks))
	for i, cur := range blocks {
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = blocks
	res.mins = mins
	res.shouldRebuildTree = true
	return res
}

func (sl *SortedSet32) Add(value S) bool {
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		sl.size++
		return true
	}

	pos, index := sl._locLeft(value)
	if index < int32(len(sl.blocks[pos])) && sl.blocks[pos][index] == value {
		return false
	}

	sl.size++
	sl._updateTree(pos, 1)
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks[:pos+1], append([][]S{sl.blocks[pos][_LOAD:]}, sl.blocks[pos+1:]...)...)
		sl.mins = append(sl.mins[:pos+1], append([]S{sl.blocks[pos][_LOAD]}, sl.mins[pos+1:]...)...)
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD] // !注意max的设置(为了让左右互不影响)
		sl.shouldRebuildTree = true
	}

	return true
}

func (sl *SortedSet32) Has(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locLeft(value)
	return index < int32(len(sl.blocks[pos])) && sl.blocks[pos][index] == value
}

func (sl *SortedSet32) Discard(value S) bool {
	if len(sl.blocks) == 0 {
		return false
	}
	pos, index := sl._locRight(value)
	if index > 0 && sl.blocks[pos][index-1] == value {
		sl._delete(pos, index-1)
		return true
	}
	return false
}

func (sl *SortedSet32) Pop(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

func (sl *SortedSet32) At(index int32) S {
	if index < 0 {
		index += sl.size
	}
	if index < 0 || index >= sl.size {
		panic("index out of range")
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SortedSet32) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SortedSet32) Lower(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedSet32) Higher(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

func (sl *SortedSet32) Floor(value S) (res S, ok bool) {
	pos := sl.BisectRight(value)
	if pos == 0 {
		return
	}
	return sl.At(pos - 1), true
}

func (sl *SortedSet32) Ceiling(value S) (res S, ok bool) {
	pos := sl.BisectLeft(value)
	if pos == sl.size {
		return
	}
	return sl.At(pos), true
}

// 返回第一个大于等于 `value` 的元素的索引/严格小于 `value` 的元素的个数.
func (sl *SortedSet32) BisectLeft(value S) int32 {
	pos, index := sl._locLeft(value)
	return sl._queryTree(pos) + index
}

// 返回第一个严格大于 `value` 的元素的索引/小于等于 `value` 的元素的个数.
func (sl *SortedSet32) BisectRight(value S) int32 {
	pos, index := sl._locRight(value)
	return sl._queryTree(pos) + index
}

func (sl *SortedSet32) Count(value S) int32 {
	return sl.BisectRight(value) - sl.BisectLeft(value)
}

func (sl *SortedSet32) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SortedSet32) ForEach(f func(value S, index int32) bool, reverse bool) {
	if !reverse {
		count := int32(0)
		for i := 0; i < len(sl.blocks); i++ {
			block := sl.blocks[i]
			for j := 0; j < len(block); j++ {
				if f(block[j], count) {
					return
				}
				count++
			}
		}
		return
	}
	count := int32(0)
	for i := len(sl.blocks) - 1; i >= 0; i-- {
		block := sl.blocks[i]
		for j := len(block) - 1; j >= 0; j-- {
			if f(block[j], count) {
				return
			}
			count++
		}
	}
}

func (sl *SortedSet32) Enumerate(start, end int32, f func(value S), erase bool) {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return
	}

	pos, startIndex := sl._findKth(start)
	count := end - start
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endIndex := min32(int32(len(block)), startIndex+count)
		if f != nil {
			for j := startIndex; j < endIndex; j++ {
				f(block[j])
			}
		}
		deleted := endIndex - startIndex

		if erase {
			if deleted == int32(len(block)) {
				// !delete block
				sl.blocks = append(sl.blocks[:pos], sl.blocks[pos+1:]...)
				sl.mins = append(sl.mins[:pos], sl.mins[pos+1:]...)
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = append(sl.blocks[pos][:startIndex], sl.blocks[pos][endIndex:]...)
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SortedSet32) Slice(start, end int32) []S {
	if start < 0 {
		start = 0
	}
	if end > sl.size {
		end = sl.size
	}
	if start >= end {
		return nil
	}
	count := end - start
	res := make([]S, 0, count)
	pos, index := sl._findKth(start)
	m := int32(len(sl.blocks))
	for ; count > 0 && pos < m; pos++ {
		block := sl.blocks[pos]
		endPos := min32(int32(len(block)), index+count)
		curCount := endPos - index
		res = append(res, block[index:endPos]...)
		count -= curCount
		index = 0
	}
	return res
}

func (sl *SortedSet32) Range(min, max S) []S {
	if sl.less(max, min) {
		return nil
	}
	res := []S{}
	pos := sl._locBlock(min)
	m := int32(len(sl.blocks))
	for i := pos; i < m; i++ {
		block := sl.blocks[i]
		for j := 0; j < len(block); j++ {
			x := block[j]
			if sl.less(max, x) {
				return res
			}
			if !sl.less(x, min) {
				res = append(res, x)
			}
		}
	}
	return res
}

func (sl *SortedSet32) Min() S {
	if sl.size == 0 {
		panic("Min() called on empty SortedList")
	}
	return sl.mins[0]
}

func (sl *SortedSet32) Max() S {
	if sl.size == 0 {
		panic("Max() called on empty SortedList")
	}
	lastBlock := sl.blocks[len(sl.blocks)-1]
	return lastBlock[len(lastBlock)-1]
}

func (sl *SortedSet32) String() string {
	sb := strings.Builder{}
	sb.WriteString("SortedList{")
	sl.ForEach(func(value S, index int32) bool {
		if index > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(fmt.Sprintf("%v", value))
		return false
	}, false)
	sb.WriteByte('}')
	return sb.String()
}

func (sl *SortedSet32) Len() int32 {
	return sl.size
}

func (sl *SortedSet32) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	copy(sl.blocks[pos][index:], sl.blocks[pos][index+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	copy(sl.blocks[pos:], sl.blocks[pos+1:])
	sl.blocks = sl.blocks[:len(sl.blocks)-1]
	copy(sl.mins[pos:], sl.mins[pos+1:])
	sl.mins = sl.mins[:len(sl.mins)-1]
	sl.shouldRebuildTree = true
}

func (sl *SortedSet32) _locLeft(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	pos = right

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(cur[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedSet32) _locRight(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, sl.mins[mid]) {
			right = mid
		} else {
			left = mid
		}
	}
	pos = left

	// find index
	cur := sl.blocks[pos]
	left = -1
	right = int32(len(cur))
	for left+1 < right {
		mid := (left + right) >> 1
		if sl.less(value, cur[mid]) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SortedSet32) _locBlock(value S) int32 {
	left, right := int32(-1), int32(len(sl.blocks)-1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !sl.less(sl.mins[mid], value) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		if !sl.less(block[len(block)-1], value) {
			right--
		}
	}
	return right
}

func (sl *SortedSet32) _buildTree() {
	sl.tree = make([]int32, len(sl.blocks))
	for i := 0; i < len(sl.blocks); i++ {
		sl.tree[i] = int32(len(sl.blocks[i]))
	}
	tree := sl.tree
	for i := 0; i < len(tree); i++ {
		j := i | (i + 1)
		if j < len(tree) {
			tree[j] += tree[i]
		}
	}
	sl.shouldRebuildTree = false
}

func (sl *SortedSet32) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := index; i < m; i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SortedSet32) _queryTree(end int32) int32 {
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	sum := int32(0)
	for end > 0 {
		sum += tree[end-1]
		end &= end - 1
	}
	return sum
}

func (sl *SortedSet32) _findKth(k int32) (pos, index int32) {
	if k < int32(len(sl.blocks[0])) {
		return 0, k
	}
	last := int32(len(sl.blocks) - 1)
	lastLen := int32(len(sl.blocks[last]))
	if k >= sl.size-lastLen {
		return last, k + lastLen - sl.size
	}
	if sl.shouldRebuildTree {
		sl._buildTree()
	}
	tree := sl.tree
	pos = -1
	bitLength := bits.Len32(uint32(len(tree)))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < int32(len(tree)) && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
