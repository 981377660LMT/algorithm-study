// https://ei1333.github.io/library/graph/mst/boruvka.hpp
// Boruvka(最小全域木) 在线最小生成树
// 不预先给出图，
// 而是给定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
	"unsafe"
)

const INF int = 2e18

func main() {
	SpeedRunMSTEasy()
	// SpeedRunMSTHard()
}

// P - MST (Easy)
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_g
// 给定数组A，长度为N。
// 给定一张n个顶点的完全图，边(u,v)的权值为A[u]*A[v]。
// 求这张图的最小生成树的权值。
// n<=2e5.
func SpeedRunMSTEasy() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	sort.Ints(A) // !从小到大排序

	set := NewFastSetFrom(n, func(i int) bool { return true })
	setUsed := func(u int) { set.Erase(u) }
	setUnused := func(u int) { set.Insert(u) }
	findUnused := func(u int) (v int, cost int) {
		min_ := set.Next(0)
		if min_ == n {
			return -1, -1
		}
		max_ := set.Prev(n)
		best := min_
		if A[max_]*A[u] < A[best]*A[u] {
			best = max_
		}
		return best, A[best] * A[u]
	}

	edges := OnlineMST(n, setUsed, setUnused, findUnused)
	res := 0
	for _, e := range edges {
		res += e[2]
	}
	fmt.Fprintln(out, res)
}

// P - MST (Hard)
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
// 给定两个数组A和B，长度为N。
// 给定一张n个顶点的完全图，边(u,v)的权值为A[u]*A[v]+B[u]*B[v]。
// 求这张图的最小生成树的权值。
// n<=5e4.
//
// cht维护二维点集，每次求出给定(x,y)条件下使得ax+by最小的点(a,b).
// 总时间复杂度O(n(logn)^2)
func SpeedRunMSTHard() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	A := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &A[i])
	}
	B := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &B[i])
	}

	var cht *LineContainer2DWithId
	init := func() {
		if cht == nil {
			cht = NewLineContainer2DWithId(n)
		} else {
			cht.Clear()
		}
	}
	add := func(u int) { cht.Add(A[u], B[u], int32(u)) }
	find := func(u int) (v int, cost int) {
		f, id := cht.QueryMin(A[u], B[u])
		return int(id), f
	}
	edges := OnlineMSTIncremental(n, init, add, find)
	res := 0
	for _, e := range edges {
		res += e[2]
	}
	fmt.Fprintln(out, res)
}

// Brouvka
//
//	不预先给出图，而是指定一个函数 findUnused 来找到未使用过的点中与u权值最小的点。
//	findUnused(u)：返回 unused 中与 u 权值最小的点 v 和边权 cost
//	              如果不存在，返回 (-1,*)
func OnlineMST(
	n int,
	setUsed func(u int), setUnused func(u int), findUnused func(u int) (v int, cost int),
) (res [][3]int) {
	uf := newUnionFindArraySimple(n)
	res = make([][3]int, 0, n-1)
	for {
		updated := false
		groups := make([][]int, n)
		cand := make([][3]int, n) // [u, v, cost]
		for v := 0; v < n; v++ {
			cand[v] = [3]int{-1, -1, -1}
		}

		for v := 0; v < n; v++ {
			leader := uf.Find(v)
			groups[leader] = append(groups[leader], v)
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			for _, x := range groups[v] {
				setUsed(x)
			}
			for _, x := range groups[v] {
				y, cost := findUnused(x)
				if y == -1 {
					continue
				}
				a, c := cand[v][0], cand[v][2]
				if a == -1 || cost < c {
					cand[v] = [3]int{x, y, cost}
				}
			}
			for _, x := range groups[v] {
				setUnused(x)
			}
		}

		for v := 0; v < n; v++ {
			if uf.Find(v) != v {
				continue
			}
			a, b, c := cand[v][0], cand[v][1], cand[v][2]
			if a == -1 {
				continue
			}
			updated = true
			if uf.Union(a, b) {
				res = append(res, [3]int{a, b, c})
			}
		}

		if !updated {
			break
		}
	}

	return res
}

// Brouvka 在线最小生成树，配合incremental(可增量计算的)数据结构使用.
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
func OnlineMSTIncremental(
	n int,
	init func(), // 初始化数据结构，大约调用 2logN 次.
	add func(u int), // 添加搜索点
	find func(u int) (v int, cost int), // (v,cost) or (-1,*)
) (res [][3]int) {
	uf := newUnionFindArraySimple(n)
	res = make([][3]int, 0, n-1)

	for uf.Part > 1 {
		updated := false
		groups := make([][]int, n)
		for v := 0; v < n; v++ {
			leader := uf.Find(v)
			groups[leader] = append(groups[leader], v)
		}
		weight := make([]int, n)
		for i := 0; i < n; i++ {
			weight[i] = INF
		}
		who := make([][2]int, n)
		for i := 0; i < n; i++ {
			who[i] = [2]int{-1, -1}
		}

		init()
		for i := 0; i < n; i++ {
			for _, v := range groups[i] {
				w, x := find(v)
				if w != -1 && x < weight[i] {
					weight[i] = x
					who[i] = [2]int{v, w}
				}
			}
			for _, v := range groups[i] {
				add(v)
			}
		}

		init()
		for i := n - 1; i >= 0; i-- {
			for _, v := range groups[i] {
				w, x := find(v)
				if w != -1 && x < weight[i] {
					weight[i] = x
					who[i] = [2]int{v, w}
				}
			}
			for _, v := range groups[i] {
				add(v)
			}
		}

		for i := 0; i < n; i++ {
			a, b := who[i][0], who[i][1]
			if a == -1 {
				continue
			}
			if uf.Union(a, b) {
				updated = true
				res = append(res, [3]int{a, b, weight[i]})
			}
		}

		if !updated {
			break
		}
	}

	return res
}

// O(n^2) Prim求完全图最小生成树.
// https://atcoder.jp/contests/pakencamp-2023-day1/tasks/pakencamp_2023_day1_p
func Prim(n int, cost func(u, v int) int) [][3]int {
	res := make([][3]int, 0, n-1)
	weight := make([]int, n)
	for i := 0; i < n; i++ {
		weight[i] = INF
	}
	to := make([]int, n)
	add := func(v int) {
		for w := 0; w < n; w++ {
			if to[w] != -1 && chmin(&weight[w], cost(v, w)) {
				to[w] = v
			}
		}
		weight[v] = INF
		to[v] = -1
	}
	add(0)
	for i := 0; i < n-1; i++ {
		argMin := 0
		for j := 1; j < n; j++ {
			if weight[j] < weight[argMin] {
				argMin = j
			}
		}
		res = append(res, [3]int{to[argMin], argMin, weight[argMin]})
		add(argMin)
	}
	return res
}

type unionFindArraySimple struct {
	Part int
	n    int
	data []int32
}

func newUnionFindArraySimple(n int) *unionFindArraySimple {
	data := make([]int32, n)
	for i := 0; i < n; i++ {
		data[i] = -1
	}
	return &unionFindArraySimple{Part: n, n: n, data: data}
}

func (u *unionFindArraySimple) Union(key1 int, key2 int) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = int32(root1)
	u.Part--
	return true
}

func (u *unionFindArraySimple) Find(key int) int {
	if u.data[key] < 0 {
		return key
	}
	u.data[key] = int32(u.Find(int(u.data[key])))
	return int(u.data[key])
}

func (u *unionFindArraySimple) GetSize(key int) int {
	return int(-u.data[u.Find(key)])
}

type FastSet struct {
	n, lg int
	seg   [][]int
	size  int
}

func NewFastSet(n int) *FastSet {
	res := &FastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func NewFastSetFrom(n int, f func(i int) bool) *FastSet {
	res := NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *FastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *FastSet) Enumerate(start, end int, f func(i int)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet) Size() int {
	return fs.size
}

func (*FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

func chmin(a *int, b int) bool {
	if *a > b {
		*a = b
		return true
	}
	return false
}

type Line struct {
	k, b   int
	p1, p2 int // p=p1/p2
}

type LineContainer2DWithId struct {
	minCHT, maxCHT       *_LineContainer
	kMax, kMin           int
	bMax, bMin           int
	kMaxIndex, kMinIndex int32
	bMaxIndex, bMinIndex int32
	mp                   map[[2]int]int32
}

func NewLineContainer2DWithId(capacity int) *LineContainer2DWithId {
	return &LineContainer2DWithId{
		minCHT: _NewLineContainer(true),
		maxCHT: _NewLineContainer(false),
		kMax:   -INF, kMin: INF, bMax: -INF, bMin: INF,
		kMaxIndex: -1, kMinIndex: -1, bMaxIndex: -1, bMinIndex: -1,
		mp: make(map[[2]int]int32, capacity),
	}
}

// 追加 a*x + b*y.
func (lc *LineContainer2DWithId) Add(a, b int, id int32) {
	lc.minCHT.Add(b, a)
	lc.maxCHT.Add(b, a)
	pair := [2]int{a, b}
	lc.mp[pair] = id

	if a > lc.kMax {
		lc.kMax = a
		lc.kMaxIndex = id
	}
	if a < lc.kMin {
		lc.kMin = a
		lc.kMinIndex = id
	}
	if b > lc.bMax {
		lc.bMax = b
		lc.bMaxIndex = id
	}
	if b < lc.bMin {
		lc.bMin = b
		lc.bMinIndex = id
	}
}

// 查询 x=xi,y=yi 时的最大值 max_{a,b} (ax + by)和对应的点id.
func (lc *LineContainer2DWithId) QueryMax(x, y int) (int, int32) {
	if lc.minCHT.Size() == 0 {
		return -INF, -1
	}

	if x == 0 {
		if y > 0 {
			return lc.bMax * y, lc.bMaxIndex
		}
		return lc.bMin * y, lc.bMinIndex
	}
	if y == 0 {
		if x > 0 {
			return lc.kMax * x, lc.kMaxIndex
		}
		return lc.kMin * x, lc.kMinIndex
	}

	// y/x
	if x > 0 {
		line := lc.maxCHT.sl.BisectLeftByPairForValue(y, x)
		a := line.b
		b := line.k
		return a*x + b*y, lc.mp[[2]int{a, b}]
	}
	line := lc.minCHT.sl.BisectLeftByPairForValue(y, x)
	a := -line.b
	b := -line.k
	return a*x + b*y, lc.mp[[2]int{a, b}]
}

// 查询 x=xi,y=yi 时的最小值 min_{a,b} (ax + by).
func (lc *LineContainer2DWithId) QueryMin(x, y int) (int, int32) {
	v, i := lc.QueryMax(-x, -y)
	return -v, i
}

func (lc *LineContainer2DWithId) Clear() {
	lc.minCHT.Clear()
	lc.maxCHT.Clear()
	lc.kMax, lc.kMin = -INF, INF
	lc.bMax, lc.bMin = -INF, INF
	lc.kMaxIndex, lc.kMinIndex = -1, -1
	lc.bMaxIndex, lc.bMinIndex = -1, -1
	lc.mp = make(map[[2]int]int32)
}

type _LineContainer struct {
	minimize bool
	sl       *SpecializedSortedList
}

func _NewLineContainer(minimize bool) *_LineContainer {
	return &_LineContainer{
		minimize: minimize,
		sl:       NewSpecializedSortedList(func(a, b S) bool { return a.k < b.k }),
	}
}

func (lc *_LineContainer) Add(k, m int) {
	if lc.minimize {
		k, m = -k, -m
	}

	newLine := &Line{k: k, b: m}
	lc.sl.Add(newLine)

	iter := lc.sl.BisectRightByKForIterator(k)
	iter.Prev()

	{
		probe := iter.Copy()
		probe.Next()
		start := probe.ToIndex()
		removeCount := int32(0)
		for lc.insect(iter, probe) {
			probe.Next()
			removeCount++
		}
		lc.sl.Erase(start, start+removeCount)
	}

	{
		probe := iter.Copy()
		if !iter.IsBegin() {
			iter.Prev()
			if lc.insect(iter, probe) {
				probIndex := probe.ToIndex()
				probe.Next()
				lc.insect(iter, probe)
				lc.sl.Pop(probIndex)
			}
		}
	}

	if iter.IsBegin() {
		return
	}

	{
		var pivot *Line
		if iter.HasNext() {
			pivot = iter.NextValue()
		}
		end := iter.ToIndex() + 1
		removeCount := int32(0)
		for !iter.IsBegin() {
			iter.Prev()
			if lessLine(iter.Value(), iter.NextValue()) {
				break
			}
			lc.insectLine(iter.Value(), pivot)
			removeCount++
		}
		lc.sl.Erase(end-removeCount, end)
	}
}

// 查询 kx + m 的最小值（或最大值).
func (lc *_LineContainer) Query(x int) int {
	if lc.sl.Len() == 0 {
		panic("empty container")
	}
	line := lc.sl.BisectLeftByPairForValue(x, 1)
	v := line.k*x + line.b
	if lc.minimize {
		return -v
	}
	return v
}

func (lc *_LineContainer) Size() int32 { return lc.sl.Len() }

func (lc *_LineContainer) Clear() { lc.sl.Clear() }

// 这个函数在向集合添加新线或删除旧线时用于计算交点。
// 计算线性函数x和y的交点，并将结果存储在x->p中。
func (lc *_LineContainer) insect(iterX, iterY *Iterator) bool {
	if iterY.IsEnd() {
		line1 := iterX.Value()
		line1.p1 = INF
		line1.p2 = 1
		return false
	}
	line1, line2 := iterX.Value(), iterY.Value()
	if line1.k == line2.k {
		if line1.b > line2.b {
			line1.p1 = INF
			line1.p2 = 1
		} else {
			line1.p1 = INF
			line1.p2 = -1
		}
	} else {
		// lc_div
		line1.p1 = line2.b - line1.b
		line1.p2 = line1.k - line2.k
	}
	return !lessPair(line1.p1, line1.p2, line2.p1, line2.p2)
}

func (lc *_LineContainer) insectLine(line1, line2 *Line) bool {
	if line2 == nil {
		line1.p1 = INF
		line1.p2 = 1
		return false
	}
	if line1.k == line2.k {
		if line1.b > line2.b {
			line1.p1 = INF
			line1.p2 = 1
		} else {
			line1.p1 = INF
			line1.p2 = -1
		}
	} else {
		// lc_div
		line1.p1 = line2.b - line1.b
		line1.p2 = line1.k - line2.k
	}
	return !lessPair(line1.p1, line1.p2, line2.p1, line2.p2)
}
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// 分母不为0的分数比较大小
//
//	a1/b1 < a2/b2
func lessPair(a1, b1, a2, b2 int) bool {
	if a1 == INF || a2 == INF { // 有一个是+-INF
		return a1/b1 < a2/b2
	}
	diff := a1*b2 - a2*b1
	mul := b1 * b2
	return diff^mul < 0
}

func lessLine(a, b *Line) bool {
	return lessPair(a.p1, a.p2, b.p1, b.p2)
}

const _LOAD int32 = 20 // 75/100/150/200

type S = *Line

type SpecializedSortedList struct {
	less              func(a, b S) bool
	size              int32
	blocks            [][]S
	mins              []S
	tree              []int32
	shouldRebuildTree bool
}

func NewSpecializedSortedList(less func(a, b S) bool, elements ...S) *SpecializedSortedList {
	elements = append(elements[:0:0], elements...)
	res := &SpecializedSortedList{less: less}
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

func (sl *SpecializedSortedList) Erase(start, end int32) {
	sl.Enumerate(start, end, nil, true)
}

func (sl *SpecializedSortedList) Enumerate(start, end int32, f func(value S), erase bool) {
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
				sl.blocks = Replace(sl.blocks, int(pos), int(pos+1))
				sl.mins = Replace(sl.mins, int(pos), int(pos+1))
				sl.shouldRebuildTree = true
				pos--
			} else {
				// !delete [index, end)
				sl._updateTree(pos, -deleted)
				sl.blocks[pos] = Replace(sl.blocks[pos], int(startIndex), int(endIndex))
				sl.mins[pos] = sl.blocks[pos][0]
			}
			sl.size -= deleted
		}

		count -= deleted
		startIndex = 0
	}
}

func (sl *SpecializedSortedList) Add(value S) *SpecializedSortedList {
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

func (sl *SpecializedSortedList) Pop(index int32) {
	pos, startIndex := sl._findKth(index)
	sl._delete(pos, startIndex)
}

func (sl *SpecializedSortedList) At(index int32) S {
	if index < 0 || index >= sl.size {
		return nil
	}
	pos, startIndex := sl._findKth(index)
	return sl.blocks[pos][startIndex]
}

func (sl *SpecializedSortedList) BisectRightByK(k int) int32 {
	pos, index := sl._locRightByK(k)
	return sl._queryTree(pos) + index
}

// 返回一个迭代器，指向键值> key的第一个元素.
// UpperBoundByK.
func (sl *SpecializedSortedList) BisectRightByKForIterator(k int) *Iterator {
	pos, index := sl._locRightByK(k)
	return &Iterator{sl: sl, pos: pos, index: index}
}

func (sl *SpecializedSortedList) BisectLeftByPair(a, b int) int32 {
	pos, index := sl._locLeftByPair(a, b)
	return sl._queryTree(pos) + index
}

func (sl *SpecializedSortedList) BisectLeftByPairForValue(a, b int) S {
	pos, index := sl._locLeftByPair(a, b)
	return sl.blocks[pos][index]
}

func (sl *SpecializedSortedList) Clear() {
	sl.size = 0
	sl.blocks = sl.blocks[:0]
	sl.mins = sl.mins[:0]
	sl.tree = sl.tree[:0]
	sl.shouldRebuildTree = true
}

func (sl *SpecializedSortedList) Len() int32 {
	return sl.size
}

func (sl *SpecializedSortedList) _delete(pos, index int32) {
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

func (sl *SpecializedSortedList) _locLeftByPair(a, b int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(-1)
	right := int32(len(sl.blocks) - 1)
	for left+1 < right {
		mid := (left + right) >> 1
		if !lessPair(sl.mins[mid].p1, sl.mins[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}
	if right > 0 {
		block := sl.blocks[right-1]
		last := block[len(block)-1]
		if !lessPair(last.p1, last.p2, a, b) {
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
		if !lessPair(cur[mid].p1, cur[mid].p2, a, b) {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _locRight(value S) (pos, index int32) {
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

func (sl *SpecializedSortedList) _locRightByK(k int) (pos, index int32) {
	if sl.size == 0 {
		return
	}

	// find pos
	left := int32(0)
	right := int32(len(sl.blocks))
	for left+1 < right {
		mid := (left + right) >> 1
		if k < sl.mins[mid].k {
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
		if k < cur[mid].k {
			right = mid
		} else {
			left = mid
		}
	}

	index = right
	return
}

func (sl *SpecializedSortedList) _buildTree() {
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

func (sl *SpecializedSortedList) _updateTree(index, delta int32) {
	if sl.shouldRebuildTree {
		return
	}
	tree := sl.tree
	for i := index; i < int32(len(tree)); i |= i + 1 {
		tree[i] += delta
	}
}

func (sl *SpecializedSortedList) _queryTree(end int32) int32 {
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

func (sl *SpecializedSortedList) _findKth(k int32) (pos, index int32) {
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
	m := int32(len(tree))
	bitLength := bits.Len32(uint32(m))
	for d := bitLength - 1; d >= 0; d-- {
		next := pos + (1 << d)
		if next < m && k >= tree[next] {
			pos = next
			k -= tree[pos]
		}
	}
	return pos + 1, k
}

type Iterator struct {
	sl         *SpecializedSortedList
	pos, index int32
}

func (it *Iterator) HasNext() bool {
	b := it.sl.blocks
	m := int32(len(b))
	if it.pos < m-1 {
		return true
	}
	return it.pos == m-1 && it.index < int32(len(b[it.pos]))-1
}

func (it *Iterator) Next() {
	it.index++
	if it.index == int32(len(it.sl.blocks[it.pos])) {
		it.pos++
		it.index = 0
	}
}

func (it *Iterator) HasPrev() bool {
	if it.pos > 0 {
		return true
	}
	return it.pos == 0 && it.index > 0
}

func (it *Iterator) Prev() {
	it.index--
	if it.index == -1 {
		it.pos--
		it.index = int32(len(it.sl.blocks[it.pos]) - 1)
	}
}

// GetMut
func (it *Iterator) Value() S {
	return it.sl.blocks[it.pos][it.index]
}

func (it *Iterator) NextValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex++
	if newIndex == int32(len(it.sl.blocks[it.pos])) {
		newPos++
		newIndex = 0
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) PrevValue() S {
	newPos, newIndex := it.pos, it.index
	newIndex--
	if newIndex == -1 {
		newPos--
		newIndex = int32(len(it.sl.blocks[newPos]) - 1)
	}
	return it.sl.blocks[newPos][newIndex]
}

func (it *Iterator) ToIndex() int32 {
	res := it.sl._queryTree(it.pos)
	return res + it.index
}

func (it *Iterator) Copy() *Iterator {
	return &Iterator{sl: it.sl, pos: it.pos, index: it.index}
}

func (it *Iterator) IsBegin() bool {
	return it.pos == 0 && it.index == 0
}

func (it *Iterator) IsEnd() bool {
	m := int32(len(it.sl.blocks))
	return it.pos == m && it.index == 0
}
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

// Replace replaces the elements s[i:j] by the given v, and returns the modified slice.
// !Like JavaScirpt's Array.prototype.splice.
func Replace[S ~[]E, E any](s S, i, j int, v ...E) S {
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

func Insert[S ~[]E, E any](s S, i int, v ...E) S {
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
