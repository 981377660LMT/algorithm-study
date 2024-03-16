// 745. 前缀和后缀搜索
// https://leetcode.cn/problems/prefix-and-suffix-search/description/

package main

import "sort"

type WordFilter struct {
	prefixTrie, suffixTrie *Trie
	lid1, rid1             []int32
	lid2, rid2             []int32
	maxSeg                 *SegmentTree2DSparse32Fast
}

// 使用词典中的单词 words 初始化对象。
func Constructor(words []string) WordFilter {
	prefixTrie, suffixTrie := NewTrie(), NewTrie()
	points := make(map[[2]int32]int32)
	for i, word := range words {
		pos1 := prefixTrie.Insert(word)
		pos2 := suffixTrie.Insert(ReverseString(word))
		points[[2]int32{pos1, pos2}] = int32(i) // 最大索引
	}
	getLidRid := func(t *Trie) (lid, rid []int32) {
		lid, rid = make([]int32, t.Size()), make([]int32, t.Size())
		dfn := int32(0)
		var dfs func(int32)
		dfs = func(pos int32) {
			lid[pos] = dfn
			dfn++
			for _, next := range t.nodes[pos].children {
				dfs(next)
			}
			rid[pos] = dfn
		}
		dfs(0)
		return
	}
	lid1, rid1 := getLidRid(prefixTrie)
	lid2, rid2 := getLidRid(suffixTrie)
	xs, ys, ws := make([]int32, 0, len(points)), make([]int32, 0, len(points)), make([]int32, 0, len(points))
	for k, v := range points {
		xs = append(xs, lid1[k[0]])
		ys = append(ys, lid2[k[1]])
		ws = append(ws, v)
	}
	maxSeg := NewSegmentTree2DSparse32FastWithWeights(xs, ys, ws, false)
	return WordFilter{prefixTrie: prefixTrie, suffixTrie: suffixTrie, lid1: lid1, rid1: rid1, lid2: lid2, rid2: rid2, maxSeg: maxSeg}
}

// 返回词典中具有前缀 prefix 和后缀 suff 的单词的下标。
// 如果存在不止一个满足要求的下标，返回其中 最大的下标 。
// 如果不存在这样的单词，返回 -1 。
func (this *WordFilter) F(prefix string, suffix string) int {
	pos1, ok1 := this.prefixTrie.Search(prefix)
	if !ok1 {
		return -1
	}
	pos2, ok2 := this.suffixTrie.Search(ReverseString(suffix))
	if !ok2 {
		return -1
	}
	x1, x2 := this.lid1[pos1], this.rid1[pos1]
	y1, y2 := this.lid2[pos2], this.rid2[pos2]
	return int(this.maxSeg.Query(x1, x2, y1, y2))
}

type TrieNode struct {
	children map[byte]int32
	endCount int32
}

func NewTrieNode() *TrieNode {
	return &TrieNode{children: make(map[byte]int32)}
}

type Trie struct {
	nodes []*TrieNode
}

func NewTrie() *Trie {
	return &Trie{nodes: []*TrieNode{NewTrieNode()}}
}

func (t *Trie) Insert(s string) (pos int32) {
	for i := 0; i < len(s); i++ {
		char := s[i]
		if next, ok := t.nodes[pos].children[char]; ok {
			pos = next
		} else {
			next = int32(len(t.nodes))
			t.nodes[pos].children[char] = next
			t.nodes = append(t.nodes, NewTrieNode())
			pos = next
		}
	}
	t.nodes[pos].endCount++
	return pos
}

func (t *Trie) Search(s string) (pos int32, ok bool) {
	nodes := t.nodes
	for i := 0; i < len(s); i++ {
		char := s[i]
		if next, ok := nodes[pos].children[char]; ok {
			pos = next
		} else {
			return 0, false
		}
	}
	return pos, true
}

// Trie 中的节点数(包含根节点).
func (t *Trie) Size() int32 {
	return int32(len(t.nodes))
}

// 需要满足交换律.
type E = int32

func e() E        { return 0 }
func op(a, b E) E { return max32(a, b) } // TODO

type SegmentTree2DSparse32Fast struct {
	n          int32
	keyX       []int32
	keyY       []int32
	minX       int32
	allY       []int32
	pos        []int32
	indptr     []int32
	size       int32
	data       []E
	discretize bool
	unit       E
	toLeft     []int32
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(自动所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewSegmentTree2DSparse32Fast(xs, ys []int32, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	ws := make([]E, len(xs))
	for i := range ws {
		ws[i] = res.unit
	}
	res._build(xs, ys, ws)
	return res
}

// discretize:
//
//	为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
//	为 false 时不对x维度二分离散化,而是直接用x的值作为下标(自动所有x给一个偏移量minX),
//	x 维度数组长度为最大值减最小值.
func NewSegmentTree2DSparse32FastWithWeights(xs, ys []int32, ws []E, discretize bool) *SegmentTree2DSparse32Fast {
	res := &SegmentTree2DSparse32Fast{discretize: discretize, unit: e()}
	res._build(xs, ys, ws)
	return res
}

func (t *SegmentTree2DSparse32Fast) Update(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._update(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

func (t *SegmentTree2DSparse32Fast) Set(rawIndex int32, value E) {
	i := int32(1)
	p := t.pos[rawIndex]
	indPtr, toLeft := t.indptr, t.toLeft
	for {
		t._set(i, p-indPtr[i], value)
		if i >= t.size {
			break
		}
		lc := toLeft[p] - toLeft[indPtr[i]]
		rc := p - indPtr[i] - lc
		if toLeft[p+1] > toLeft[p] {
			p = indPtr[i<<1] + lc
			i <<= 1
		} else {
			p = indPtr[i<<1|1] + rc
			i = i<<1 | 1
		}
	}
}

// [lx,rx) * [ly,ry)
func (t *SegmentTree2DSparse32Fast) Query(lx, rx, ly, ry int32) E {
	L := t._xtoi(lx)
	R := t._xtoi(rx)
	res := t.unit
	indPtr, toLeft := t.indptr, t.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res = op(res, t._query(i, a, b))
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, t.size, bisectLeft(t.allY, ly, 0, int32(len(t.allY)-1)), bisectLeft(t.allY, ry, 0, int32(len(t.allY)-1)))
	return res
}

// nlogn.
func (seg *SegmentTree2DSparse32Fast) Count(lx, rx, ly, ry int32) int32 {
	L := seg._xtoi(lx)
	R := seg._xtoi(rx)
	res := int32(0)
	indPtr, toLeft := seg.indptr, seg.toLeft
	var dfs func(i, l, r, a, b int32)
	dfs = func(i, l, r, a, b int32) {
		if a == b || R <= l || r <= L {
			return
		}
		if L <= l && r <= R {
			res += b - a
			return
		}
		la := toLeft[indPtr[i]+a] - toLeft[indPtr[i]]
		ra := a - la
		lb := toLeft[indPtr[i]+b] - toLeft[indPtr[i]]
		rb := b - lb
		m := (l + r) >> 1
		dfs(i<<1, l, m, la, lb)
		dfs(i<<1|1, m, r, ra, rb)
	}
	dfs(1, 0, seg.size, bisectLeft(seg.allY, ly, 0, int32(len(seg.allY)-1)), bisectLeft(seg.allY, ry, 0, int32(len(seg.allY)-1)))
	return res
}

func (t *SegmentTree2DSparse32Fast) _update(i int32, y int32, val E) {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	y += n
	for y > 0 {
		t.data[offset+y] = op(t.data[offset+y], val)
		y >>= 1
	}
}

func (seg *SegmentTree2DSparse32Fast) _set(i, y int32, val E) {
	lid := seg.indptr[i]
	n := seg.indptr[i+1] - seg.indptr[i]
	off := lid << 1
	y += n
	seg.data[off+y] = val
	for y > 1 {
		y >>= 1
		seg.data[off+y] = op(seg.data[off+y<<1], seg.data[off+y<<1|1])
	}
}

func (t *SegmentTree2DSparse32Fast) _query(i int32, ly, ry int32) E {
	lid := t.indptr[i]
	n := t.indptr[i+1] - t.indptr[i]
	offset := lid << 1
	left, right := n+ly, n+ry
	val := t.unit
	for left < right {
		if left&1 == 1 {
			val = op(val, t.data[offset+left])
			left++
		}
		if right&1 == 1 {
			right--
			val = op(t.data[offset+right], val)
		}
		left >>= 1
		right >>= 1
	}
	return val
}

func (seg *SegmentTree2DSparse32Fast) _build(X, Y []int32, wt []E) {
	if len(X) != len(Y) || len(X) != len(wt) {
		panic("Lengths of X, Y, and wt must be equal.")
	}

	if seg.discretize {
		seg.keyX = unique(X)
		seg.n = int32(len(seg.keyX))
	} else {
		if len(X) > 0 {
			min_, max_ := int32(0), int32(0)
			for _, x := range X {
				if x < min_ {
					min_ = x
				}
				if x > max_ {
					max_ = x
				}
			}
			seg.minX = min_
			seg.n = max_ - min_ + 1
		}
	}

	log := int32(0)
	for 1<<log < seg.n {
		log++
	}
	size := int32(1 << log)
	seg.size = size

	orderX := make([]int32, len(X))
	for i := range orderX {
		orderX[i] = seg._xtoi(X[i])
	}
	seg.indptr = make([]int32, 2*size+1)
	for _, i := range orderX {
		i += size
		for i > 0 {
			seg.indptr[i+1]++
			i >>= 1
		}
	}
	for i := int32(1); i <= 2*size; i++ {
		seg.indptr[i] += seg.indptr[i-1]
	}
	seg.data = make([]E, 2*seg.indptr[2*size])
	for i := range seg.data {
		seg.data[i] = seg.unit
	}

	seg.toLeft = make([]int32, seg.indptr[size]+1)
	ptr := append([]int32(nil), seg.indptr...)
	order := argSort(Y)
	seg.pos = make([]int32, len(X))
	for i, v := range order {
		seg.pos[v] = int32(i)
	}
	for _, rawIdx := range order {
		i := orderX[rawIdx] + size
		j := int32(-1)
		for i > 0 {
			p := ptr[i]
			ptr[i]++
			seg.data[seg.indptr[i+1]+p] = wt[rawIdx]
			if j != -1 && j&1 == 0 {
				seg.toLeft[p+1] = 1
			}
			j = i
			i >>= 1
		}
	}
	for i := int32(1); i < int32(len(seg.toLeft)); i++ {
		seg.toLeft[i] += seg.toLeft[i-1]
	}

	for i := int32(0); i < 2*size; i++ {
		off := 2 * seg.indptr[i]
		n := seg.indptr[i+1] - seg.indptr[i]
		for j := n - 1; j >= 1; j-- {
			seg.data[off+j] = op(seg.data[off+j<<1], seg.data[off+j<<1|1])
		}
	}

	allY := append([]int32(nil), Y...)
	sort.Slice(allY, func(i, j int) bool { return allY[i] < allY[j] })
	seg.allY = allY
}

func (seg *SegmentTree2DSparse32Fast) _xtoi(x int32) int32 {
	if seg.discretize {
		return bisectLeft(seg.keyX, x, 0, int32(len(seg.keyX)-1))
	}
	tmp := x - seg.minX
	if tmp < 0 {
		tmp = 0
	} else if tmp > seg.n {
		tmp = seg.n
	}
	return tmp
}

func bisectLeft(nums []int32, x int32, left, right int32) int32 {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

func unique(nums []int32) (sorted []int32) {
	set := make(map[int32]struct{}, len(nums))
	for _, v := range nums {
		set[v] = struct{}{}
	}
	sorted = make([]int32, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Slice(sorted, func(i, j int) bool { return sorted[i] < sorted[j] })
	return
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

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func argSort(nums []int32) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func ReverseString(s string) string {
	n := len(s)
	runes := make([]rune, n)
	for _, r := range s {
		n--
		runes[n] = r
	}
	return string(runes)
}
