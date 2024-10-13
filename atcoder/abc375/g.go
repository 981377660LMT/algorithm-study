package main

import (
	"bufio"
	"bytes"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
	"strconv"
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
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, M := io.NextInt(), io.NextInt()
	edges := make([][3]int, M)
	for i := 0; i < M; i++ {
		u, v, w := io.NextInt()-1, io.NextInt()-1, io.NextInt()
		edges[i] = [3]int{u, v, w}
	}

	sp := NewUndirectedModifiableShortestPath(int32(N), int32(M), 0, int32(N)-1)
	for i := 0; i < M; i++ {
		u, v, w := edges[i][0], edges[i][1], edges[i][2]
		sp.AddEdge(int32(u), int32(v), int(w))
	}
	d := sp.MinDist()

	for i := 0; i < M; i++ {
		cur := sp.QueryOnDeleteEdge(int32(i))
		if cur != d {
			io.Println("Yes")
		} else {
			io.Println("No")
		}
	}

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

func (sp *UndirectedModifiableShortestPath) QueryOnAddEdge(a, b int32, dist int) int {
	sp.prepare()
	res := sp.dst.distToSrc
	res = min(res, sp.nodes[a].distToSrc+dist+sp.nodes[b].distToDst)
	res = min(res, sp.nodes[b].distToSrc+dist+sp.nodes[a].distToDst)
	return res
}

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

const _LOAD int32 = 200

type S = *node

var EMPTY S

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
		blocks = append(blocks, elements[start:end:end])
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
	sl.blocks[pos] = append(sl.blocks[pos][:index], append([]S{value}, sl.blocks[pos][index:]...)...)
	sl.mins[pos] = sl.blocks[pos][0]

	if n := int32(len(sl.blocks[pos])); _LOAD+_LOAD < n {
		sl.blocks = append(sl.blocks, nil)
		copy(sl.blocks[pos+2:], sl.blocks[pos+1:])
		sl.blocks[pos+1] = sl.blocks[pos][_LOAD:]
		sl.blocks[pos] = sl.blocks[pos][:_LOAD:_LOAD]
		sl.mins = append(sl.mins, EMPTY)
		copy(sl.mins[pos+2:], sl.mins[pos+1:])
		sl.mins[pos+1] = sl.blocks[pos+1][0]
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

func (sl *sortedList32) Range(min, max S) []S {
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

func (sl *sortedList32) Len() int32 {
	return sl.size
}

func (sl *sortedList32) _delete(pos, index int32) {
	sl.size--
	sl._updateTree(pos, -1)
	copy(sl.blocks[pos][index:], sl.blocks[pos][index+1:])
	sl.blocks[pos] = sl.blocks[pos][:len(sl.blocks[pos])-1]
	if len(sl.blocks[pos]) > 0 {
		sl.mins[pos] = sl.blocks[pos][0]
		return
	}

	copy(sl.blocks[pos:], sl.blocks[pos+1:])
	sl.blocks = sl.blocks[:len(sl.blocks)-1]
	copy(sl.mins[pos:], sl.mins[pos+1:])
	sl.mins = sl.mins[:len(sl.mins)-1]
	sl.shouldRebuildTree = true
}

func (sl *sortedList32) _locLeft(value S) (pos, index int32) {
	if sl.size == 0 {
		return
	}

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
