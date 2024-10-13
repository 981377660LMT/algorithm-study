// Dynamic Shortest Path
// 指定起点终点的动态最短路，带修改的无向图最短路/带修改的最短路问题

package main

import (
	"bufio"
	"bytes"
	"fmt"
	stdio "io"
	"math/bits"
	"math/rand"
	"os"
	"sort"
	"strconv"
	"time"
	"unsafe"
)

var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

func main() {
	// CF843D()
	CF1163F()
	// demo()
}

// DynamicShortestPath
// https://www.luogu.com.cn/problem/CF843D
// 1 1 v: 查询从1到v的最短路
// 2 m e1...em: e1到em边的权值+1
// n,m<=1e5,q<=1e3
func CF843D() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()
}

// Indecisive Taxi Fee
// https://www.luogu.com.cn/problem/CF1163F
// 给定n个顶点m条无向带权边构成的连通图，指定起点和终点。
// 之后q个请求，第i请求询问，假如修改编号为ei的边的权重为wi，
// 要求回答起点到终点的最短距离（注意每个请求不会对之后的请求产生影响）。
// 其中1≤n,m,q≤106，且每条边的权重为1到1e9之间
func CF1163F() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n, m, q := int32(io.NextInt()), int32(io.NextInt()), int32(io.NextInt())
	sp := NewUndirectedModifiableShortestPath(n, m, 0, n-1)
	for i := int32(0); i < m; i++ {
		u, v, w := int32(io.NextInt()), int32(io.NextInt()), int32(io.NextInt())
		u, v = u-1, v-1
		sp.AddEdge(u, v, int(w))
	}
	for i := int32(0); i < q; i++ {
		ei, w := int32(io.NextInt()), int32(io.NextInt())
		ei--
		io.Println(sp.QueryOnModifyEdge(ei, int(w)))
	}
}

func demo() {
	n := int32(2e5)
	m := int32(2e5)
	time1 := time.Now()
	S := NewUndirectedModifiableShortestPath(n, m, 0, n-1)
	for i := int32(0); i < m; i++ {
		S.AddEdge(int32(rand.Intn(int(n))), int32(rand.Intn(int(n))), rand.Intn(1000000000)+1)
	}
	for i := int32(0); i < m; i++ {
		S.QueryOnModifyEdge(m-i-1, rand.Intn(1000000000)+1)
	}
	fmt.Println(time.Since(time1))
}

const INF int = 2e18
const INF32 int32 = 1e9 + 10

type UndirectedModifiableShortestPath struct {
	curTag   int32
	seg      *segmentTreeDual32
	nodes    []*node
	edges    []*edge
	src      *node
	dst      *node
	prepared bool
}

func NewUndirectedModifiableShortestPath(n, m, s, t int32) *UndirectedModifiableShortestPath {
	nodes := make([]*node, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = newNode()
		nodes[i].id = i
		nodes[i].distToSrc = INF
		nodes[i].distToDst = INF
	}
	edges := make([]*edge, 0, m)
	src := nodes[s]
	dst := nodes[t]
	return &UndirectedModifiableShortestPath{nodes: nodes, edges: edges, src: src, dst: dst}
}

func (sp *UndirectedModifiableShortestPath) AddEdge(u, v int32, cost int) {
	if sp.prepared {
		panic("can't add edge after prepare")
	}
	e := newEdge()
	e.a = sp.nodes[u]
	e.b = sp.nodes[v]
	e.w = cost
	e.a.adj = append(e.a.adj, e)
	e.b.adj = append(e.b.adj, e)
	sp.edges = append(sp.edges, e)
}

func (sp *UndirectedModifiableShortestPath) MinDist() int {
	sp.prepare()
	return sp.dst.distToSrc
}

// O(1)
func (sp *UndirectedModifiableShortestPath) QueryOnAddEdge(a, b int32, dist int) int {
	sp.prepare()
	res := sp.dst.distToSrc
	res = min(res, sp.nodes[a].distToSrc+dist+sp.nodes[b].distToDst)
	res = min(res, sp.nodes[b].distToSrc+dist+sp.nodes[a].distToDst)
	return res
}

// O(log n)
func (sp *UndirectedModifiableShortestPath) QueryOnModifyEdge(i int32, cost int) int {
	sp.prepare()
	e := sp.edges[i]
	w := cost
	res := min(e.a.distToSrc+e.b.distToDst+w, e.a.distToDst+e.b.distToSrc+w)
	if e.tag == -1 {
		res = min(res, sp.dst.distToSrc)
	} else {
		val := sp.seg.Get(e.tag)
		res = min(res, val)
	}
	return res
}

// O(log n)
func (sp *UndirectedModifiableShortestPath) QueryOnDeleteEdge(i int32) int {
	return sp.QueryOnModifyEdge(i, INF)
}

func (sp *UndirectedModifiableShortestPath) update(a, b *node, w int) {
	dist := a.distToSrc + b.distToDst + w
	l := sp.prev(a)
	r := sp.post(b)
	sp.seg.Update(l+1, r, dist)
}

func (sp *UndirectedModifiableShortestPath) post(root *node) int32 {
	if root.r == -INF32 {
		root.r = sp.curTag + 1
		for _, e := range root.adj {
			node := e.Other(root)
			if node.distToDst+e.w == root.distToDst {
				if e.tag != -1 {
					root.r = e.tag
				} else {
					root.r = sp.post(node)
				}
				break
			}
		}
	}
	return root.r
}

func (sp *UndirectedModifiableShortestPath) prev(root *node) int32 {
	if root.l == -INF32 {
		root.l = -1
		for _, e := range root.adj {
			node := e.Other(root)
			if node.distToSrc+e.w == root.distToSrc {
				if e.tag != -1 {
					root.l = e.tag
				} else {
					root.l = sp.prev(node)
				}
				break
			}
		}
	}
	return root.l
}

func (sp *UndirectedModifiableShortestPath) prepare() {
	if sp.prepared {
		return
	}
	sp.prepared = true

	pq := newSortedList32(func(a, b S) bool {
		if a.distToSrc == b.distToSrc {
			return a.id < b.id
		}
		return a.distToSrc < b.distToSrc
	})

	sp.src.distToSrc = 0
	pq.Add(sp.src)
	for pq.Len() > 0 {
		head := pq._popFirst()
		for _, e := range head.adj {
			node := e.Other(head)
			if tmp := head.distToSrc + e.w; node.distToSrc > tmp {
				pq.Discard(node)
				node.distToSrc = tmp
				pq.Add(node)
			}
		}
	}

	sp.dst.distToDst = 0
	pq = newSortedList32(func(a, b S) bool {
		if a.distToDst == b.distToDst {
			return a.id < b.id
		}
		return a.distToDst < b.distToDst
	})
	pq.Add(sp.dst)
	for pq.Len() > 0 {
		head := pq._popFirst()
		for _, e := range head.adj {
			node := e.Other(head)
			if tmp := head.distToDst + e.w; node.distToDst > tmp {
				pq.Discard(node)
				node.distToDst = tmp
				pq.Add(node)
			}
		}
	}

	for trace := sp.src; trace != sp.dst; {
		var next *edge
		for _, e := range trace.adj {
			node := e.Other(trace)
			if node.distToDst+e.w == trace.distToDst {
				next = e
				break
			}
		}
		if next == nil {
			return
		}
		sp.curTag++
		next.tag = sp.curTag
		trace = next.Other(trace)
	}

	sp.seg = newSegmentTreeDual32(sp.curTag + 1)
	for _, e := range sp.edges {
		if e.tag != -1 {
			continue
		}
		sp.update(e.a, e.b, e.w)
		sp.update(e.b, e.a, e.w)
	}
}

type edge struct {
	a, b *node
	w    int
	tag  int32
}

func newEdge() *edge {
	return &edge{tag: -1}
}

func (e *edge) Other(x *node) *node {
	if x == e.a {
		return e.b
	}
	return e.a
}

type node struct {
	adj                  []*edge
	distToSrc, distToDst int
	l, r                 int32
	id                   int32
}

func newNode() *node {
	return &node{l: -INF32, r: -INF32}
}

// RangeChminPointGet

type Id = int

const COMMUTATIVE = true

func (*segmentTreeDual32) id() Id                 { return INF }
func (*segmentTreeDual32) composition(f, g Id) Id { return min(f, g) }

type segmentTreeDual32 struct {
	n            int32
	size, height int32
	lazy         []Id
	unit         Id
}

func newSegmentTreeDual32(n int32) *segmentTreeDual32 {
	res := &segmentTreeDual32{}
	size := int32(1)
	height := int32(0)
	for size < n {
		size <<= 1
		height++
	}
	lazy := make([]Id, 2*size)
	unit := res.id()
	for i := int32(0); i < 2*size; i++ {
		lazy[i] = unit
	}
	res.n = n
	res.size = size
	res.height = height
	res.lazy = lazy
	res.unit = unit
	return res
}
func (seg *segmentTreeDual32) Get(index int32) Id {
	index += seg.size
	for i := seg.height; i > 0; i-- {
		seg.propagate(index >> i)
	}
	return seg.lazy[index]
}
func (seg *segmentTreeDual32) GetAll() []Id {
	for i := int32(0); i < seg.size; i++ {
		seg.propagate(i)
	}
	res := make([]Id, seg.n)
	copy(res, seg.lazy[seg.size:seg.size+seg.n])
	return res
}
func (seg *segmentTreeDual32) Update(left, right int32, value Id) {
	if left < 0 {
		left = 0
	}
	if right > seg.n {
		right = seg.n
	}
	if left >= right {
		return
	}
	left += seg.size
	right += seg.size
	if !COMMUTATIVE {
		for i := seg.height; i > 0; i-- {
			if (left>>i)<<i != left {
				seg.propagate(left >> i)
			}
			if (right>>i)<<i != right {
				seg.propagate((right - 1) >> i)
			}
		}
	}
	for left < right {
		if left&1 > 0 {
			seg.lazy[left] = seg.composition(value, seg.lazy[left])
			left++
		}
		if right&1 > 0 {
			right--
			seg.lazy[right] = seg.composition(value, seg.lazy[right])
		}
		left >>= 1
		right >>= 1
	}
}
func (seg *segmentTreeDual32) propagate(k int32) {
	if seg.lazy[k] != seg.unit {
		seg.lazy[k<<1] = seg.composition(seg.lazy[k], seg.lazy[k<<1])
		seg.lazy[k<<1|1] = seg.composition(seg.lazy[k], seg.lazy[k<<1|1])
		seg.lazy[k] = seg.unit
	}
}
func (st *segmentTreeDual32) String() string {
	var buf bytes.Buffer
	buf.WriteByte('[')
	for i := int32(0); i < st.n; i++ {
		if i > 0 {
			buf.WriteByte(' ')
		}
		buf.WriteString(fmt.Sprint(st.Get(i)))
	}
	buf.WriteByte(']')
	return buf.String()
}

// 1e5 -> 200, 2e5 -> 400
const _LOAD int32 = 200

type S = *node

var EMPTY S

// 使用分块+树状数组维护的有序序列.
type sortedList32 struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func newSortedList32(less func(a, b S) bool, elements ...S) *sortedList32 {
	elements = append(elements[:0:0], elements...)
	res := &sortedList32{less: less}
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

func (sl *sortedList32) Add(value S) *sortedList32 {
	sl.size++
	if len(sl.blocks) == 0 {
		sl.blocks = append(sl.blocks, []S{value})
		sl.mins = append(sl.mins, value)
		sl.shouldRebuildTree = true
		return sl
	}

	pos, index := sl._locRight(value)

	sl._updateTree(pos, 1)
	sl.blocks[pos] = Insert(sl.blocks[pos], int(index), value)
	sl.mins[pos] = sl.blocks[pos][0]

	// n -> load + (n - load)
	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		left := append([]S(nil), sl.blocks[pos][:_LOAD]...)
		right := append([]S(nil), sl.blocks[pos][_LOAD:]...)
		sl.blocks = Replace(sl.blocks, int(pos), int(pos)+1, left, right)
		sl.mins = Insert(sl.mins, int(pos)+1, right[0])
		sl.shouldRebuildTree = true
	}

	return sl
}

func (sl *sortedList32) Discard(value S) bool {
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

func (sl *sortedList32) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *sortedList32) Len() int32 {
	return sl.size
}

func (sl *sortedList32) _delete(pos, index int32) {
	// !delete element
	sl.size--
	sl._updateTree(pos, -1)
	sl.blocks[pos] = Replace(sl.blocks[pos], int(index), int(index+1))
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	// !delete block
	sl.blocks = Replace(sl.blocks, int(pos), int(pos)+1)
	sl.mins = Replace(sl.mins, int(pos), int(pos)+1)
	sl.shouldRebuildTree = true
}

func (sl *sortedList32) _locLeft(value S) (pos, index int32) {
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

func (sl *sortedList32) _locRight(value S) (pos, index int32) {
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

func (sl *sortedList32) _locBlock(value S) int32 {
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

func (sl *sortedList32) _buildTree() {
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

func (sl *sortedList32) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	m := int32(len(tree))
	for i := index; i < m; i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *sortedList32) _queryTree(end int32) int32 {
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

func (sl *sortedList32) _findKth(k int32) (pos, index int32) {
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

func (sl *sortedList32) _popFirst() S {
	pos, startIndex := int32(0), int32(0)
	value := sl.blocks[pos][startIndex]
	sl._delete(pos, startIndex)
	return value
}

// Replace replaces the elements s[i:j] by the given v, and returns the modified slice.
// !Like JavaScirpt's Array.prototype.splice.
func Replace[S []E, E any](s S, i, j int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if j > len(s) {
		j = len(s)
	}
	if i == j {
		return Insert(s, i, v...)
	}
	if j == len(s) {
		return append(s[:i], v...)
	}
	tot := len(s[:i]) + len(v) + len(s[j:])
	if tot > cap(s) {
		s2 := append(s[:i], make(S, tot-i)...)
		copy(s2[i:], v)
		copy(s2[i+len(v):], s[j:])
		return s2
	}
	r := s[:tot]
	if i+len(v) <= j {
		copy(r[i:], v)
		copy(r[i+len(v):], s[j:])
		// clear(s[tot:])
		return r
	}
	if !overlaps(r[i+len(v):], v) {
		copy(r[i+len(v):], s[j:])
		copy(r[i:], v)
		return r
	}
	y := len(v) - (j - i)
	if !overlaps(r[i:j], v) {
		copy(r[i:j], v[y:])
		copy(r[len(s):], v[:y])
		rotateRight(r[i:], y)
		return r
	}
	if !overlaps(r[len(s):], v) {
		copy(r[len(s):], v[:y])
		copy(r[i:j], v[y:])
		rotateRight(r[i:], y)
		return r
	}
	k := startIdx(v, s[j:])
	copy(r[i:], v)
	copy(r[i+len(v):], r[i+k:])
	return r
}

func rotateLeft[E any](s []E, r int) {
	for r != 0 && r != len(s) {
		if r*2 <= len(s) {
			swap(s[:r], s[len(s)-r:])
			s = s[:len(s)-r]
		} else {
			swap(s[:len(s)-r], s[r:])
			s, r = s[len(s)-r:], r*2-len(s)
		}
	}
}

func rotateRight[E any](s []E, r int) {
	rotateLeft(s, len(s)-r)
}

func swap[E any](x, y []E) {
	for i := 0; i < len(x); i++ {
		x[i], y[i] = y[i], x[i]
	}
}

func overlaps[E any](a, b []E) bool {
	if len(a) == 0 || len(b) == 0 {
		return false
	}
	elemSize := unsafe.Sizeof(a[0])
	if elemSize == 0 {
		return false
	}
	return uintptr(unsafe.Pointer(&a[0])) <= uintptr(unsafe.Pointer(&b[len(b)-1]))+(elemSize-1) &&
		uintptr(unsafe.Pointer(&b[0])) <= uintptr(unsafe.Pointer(&a[len(a)-1]))+(elemSize-1)
}

func startIdx[E any](haystack, needle []E) int {
	p := &needle[0]
	for i := range haystack {
		if p == &haystack[i] {
			return i
		}
	}
	panic("needle not found")
}

func Insert[S []E, E any](s S, i int, v ...E) S {
	if i < 0 {
		i = 0
	}
	if i > len(s) {
		i = len(s)
	}

	m := len(v)
	if m == 0 {
		return s
	}
	n := len(s)
	if i == n {
		return append(s, v...)
	}
	if n+m > cap(s) {
		s2 := append(s[:i], make(S, n+m-i)...)
		copy(s2[i:], v)
		copy(s2[i+m:], s[i:])
		return s2
	}
	s = s[:n+m]
	if !overlaps(v, s[i+m:]) {
		copy(s[i+m:], s[i:])
		copy(s[i:], v)
		return s
	}
	copy(s[n:], v)
	rotateRight(s[i:], m)
	return s
}
