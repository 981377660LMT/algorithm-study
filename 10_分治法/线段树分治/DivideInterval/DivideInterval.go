package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

func main() {
	// abc339g()
	// abc342g()
	SP11470()
	// yuki1170()
}

type InnerTreeAbc339g struct {
	data   []int
	preSum []int
}

func NewInnerTreeAbc339g() *InnerTreeAbc339g {
	return &InnerTreeAbc339g{}
}

func (stl *InnerTreeAbc339g) Add(x int) {
	stl.data = append(stl.data, x)
}

func (stl *InnerTreeAbc339g) Build() {
	sort.Ints(stl.data)
	stl.preSum = make([]int, len(stl.data)+1)
	for i, x := range stl.data {
		stl.preSum[i+1] = stl.preSum[i] + x
	}
}

// 小于等于upper的元素之和.
func (stl *InnerTreeAbc339g) Query(upper int) int {
	pos := sort.SearchInts(stl.data, upper+1)
	return stl.preSum[pos]
}

// G - Smaller Sum
// https://atcoder.jp/contests/abc339/tasks/abc339_g
// 二维数点问题.
// 给定长为n的数组nums和q次查询,每次查询给定区间[l,r]和x,求区间[l,r]中小于等于x的元素之和.
// 本题强制在线。后一次询问给出的 l,r,x 需异或上前一次询问的答案。
//
// !二维区间查询不下推区间，而是采用标记永久化，减少空间使用.
func abc339g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	D := NewDivideInterval(int32(n))
	innerTree := make([]*InnerTreeAbc339g, D.Size())
	for i := range innerTree {
		innerTree[i] = NewInnerTreeAbc339g()
	}

	add := func(index, delta int) {
		D.EnumeratePoint(int32(index), func(segmentId int32) {
			innerTree[segmentId].Add(delta)
		})
	}

	query := func(start, end int, x int) int {
		res := 0
		D.EnumerateSegment(
			int32(start), int32(end),
			func(segmentId int32) {
				res += innerTree[segmentId].Query(x)
			},
			false,
		)
		return res
	}

	for i, num := range nums {
		add(i, num)
	}
	for i := range innerTree {
		innerTree[i].Build()
	}

	var q int
	fmt.Fscan(in, &q)
	preRes := 0
	for i := 0; i < q; i++ {
		var a, b, c int
		fmt.Fscan(in, &a, &b, &c)
		a ^= preRes
		b ^= preRes
		c ^= preRes
		a--
		res := query(a, b, c)
		preRes = res
		fmt.Fprintln(out, res)
	}
}

// G - Retroactive Range Chmax (可追溯区间最大值修改)
// https://atcoder.jp/contests/abc342/tasks/abc342_g
// 维护一个数列，有以下三个操作：
// 1 l r x: 将区间[l,r]中的所有元素与x取最大值.
// 2 i: 将第i次操作删除，保证第i次操作是操作1.
// 3 i: 查询当前数列中第i个元素的值.
//
// 树套树，内层树维护添加过的操作.
// !更新/删除操作相当于更新/删除内层树中存储的信息.
// !查询操作相当于找到所有更新操作取最大值.
func abc342g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	var q int
	fmt.Fscan(in, &q)
	queries := make([][4]int, q)
	for i := range queries {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 1 {
			var start, end, x int
			fmt.Fscan(in, &start, &end, &x)
			start--
			queries[i] = [4]int{kind, start, end, x}
		} else if kind == 2 {
			var time int
			fmt.Fscan(in, &time)
			time--
			queries[i] = [4]int{kind, time, 0, 0}
		} else {
			var index int
			fmt.Fscan(in, &index)
			index--
			queries[i] = [4]int{kind, index, 0, 0}
		}
	}

	D := NewDivideInterval(int32(n))
	innerTree := make([]*ErasableHeapGeneric[[2]int], D.Size()) // (value, time)
	for i := range innerTree {
		innerTree[i] = NewErasableHeapGeneric[[2]int](func(a, b [2]int) bool { return a[0] > b[0] })
	}

	apply := func(start, end int, x int, time int) {
		D.EnumerateSegment(
			int32(start), int32(end),
			func(segmentId int32) {
				innerTree[segmentId].Push([2]int{x, time})
			},
			false,
		)
	}

	revoke := func(time int) {
		start, end, x := queries[time][1], queries[time][2], queries[time][3]
		D.EnumerateSegment(
			int32(start), int32(end),
			func(segmentId int32) {
				innerTree[segmentId].Erase([2]int{x, time})
			},
			false,
		)
	}

	query := func(index int) int {
		res := nums[index]
		D.EnumeratePoint(
			int32(index),
			func(segmentId int32) {
				if tree := innerTree[segmentId]; tree.Len() > 0 {
					res = max(res, tree.Peek()[0])
				}
			},
		)
		return res
	}

	for i, item := range queries {
		kind := item[0]
		if kind == 1 {
			start, end, x := item[1], item[2], item[3]
			apply(start, end, x, i)
		} else if kind == 2 {
			time := item[1]
			revoke(time)
		} else {
			index := item[1]
			fmt.Fprintln(out, query(index))
		}
	}
}

// TTM - To the moon
// https://www.luogu.com.cn/problem/SP11470
// 给定一个数组nums和q次操作.操作有四种，初始时时间为0.
// C l r d: 将区间[l,r]中的所有元素加上d，同时当前的时间戳加 1。
// Q l r: 查询此时区间[l,r]中的所有元素之和。
// H l r t: 查询时间戳为t时区间[l,r]中的所有元素之和。
// B t : 将当前时间戳置为t.
//
// 翻译：
// - 区间加法并创建新版本.
// - 查询某个版本的区间和.
// - 版本回溯.
//
// !可持久化线段树的单点修改 => 一般操作即可.
// !可持久化线段树的区间修改 =>
// 不能下传懒标记，共用子结点直接把标记传递下去会影响答案的正确性，需要标记永久化.
// TODO
func SP11470() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

}

// No.1170 Never Want to Walk
// https://yukicoder.me/problems/no/1170
// 数轴上有n个车站,第i个位置在xi
// 如果两个车站之间的距离di与dj满足 A<=|di-dj|<=B,则这两个车站可以相互到达,否则不能相互到达
// 对每个车站,求出从该车站出发,可以到达的车站的数量
// 1<=n<=2e5 0<=A<=B<=1e9 0<=x1<=x2<...<=xn<=1e9
//
// !将序列搬到线段树上加速区间操作(线段树优化建图).
func yuki1170() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	var A, B int
	fmt.Fscan(in, &n, &A, &B)
	positions := make([]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &positions[i])
	}

	D := NewDivideInterval(n)
	uf := NewUnionFindArray(D.Size())
	weights := make([]int, D.Size())
	for i := int32(0); i < n; i++ {
		weights[D.Id(i)] = 1
	}

	f := func(big, small int32) {
		weights[big] += weights[small]
	}
	var dfs func(int32) // 线段树上dfs
	dfs = func(cur int32) {
		if D.IsLeaf(cur) {
			return
		}
		for k := int32(0); k < 2; k++ {
			child := cur<<1 | k
			if !uf.IsConnected(cur, child) {
				uf.UnionWithCallback(cur, child, f)
				dfs(child)
			}
		}
	}

	for i := int32(0); i < n; i++ {
		start := int32(sort.SearchInts(positions, positions[i]+A))
		end := int32(sort.SearchInts(positions, positions[i]+B+1))
		D.EnumerateSegment(start, end, func(segmentId int32) {
			uf.UnionWithCallback(D.Id(i), segmentId, f)
			dfs(segmentId)
		}, false)
	}

	for i := int32(0); i < n; i++ {
		fmt.Fprintln(out, weights[uf.Find(D.Id(i))])
	}
}

type DivideInterval struct {
	Offset int32 // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	n      int32
}

// 线段树分割区间.
// 将长度为n的序列搬到长度为offset+n的线段树上, 以实现快速的区间操作.
func NewDivideInterval(n int32) *DivideInterval {
	offset := int32(1)
	for offset < n {
		offset <<= 1
	}
	return &DivideInterval{Offset: offset, n: n}
}

// 获取原下标为i的元素在树中的(叶子)编号.
func (d *DivideInterval) Id(rawIndex int32) int32 {
	return rawIndex + d.Offset
}

// O(logn) 顺序遍历`[start,end)`区间对应的线段树节点.
// sorted表示是否按照节点编号从小到大的顺序遍历.
func (d *DivideInterval) EnumerateSegment(start, end int32, f func(segmentId int32), sorted bool) {
	if start < 0 {
		start = 0
	}
	if end > d.n {
		end = d.n
	}
	if start >= end {
		return
	}

	if sorted {
		for _, i := range d.getSegmentIds(start, end) {
			f(i)
		}
	} else {
		for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
			if start&1 == 1 {
				f(start)
				start++
			}
			if end&1 == 1 {
				end--
				f(end)
			}
		}
	}
}

func (d *DivideInterval) EnumeratePoint(index int32, f func(segmentId int32)) {
	if index < 0 || index >= d.n {
		return
	}
	index += d.Offset
	for index > 0 {
		f(index)
		index >>= 1
	}
}

// O(n) 从根向叶子方向push.
func (d *DivideInterval) PushDown(f func(parent, child int32)) {
	for p := int32(1); p < d.Offset; p++ {
		f(p, p<<1)
		f(p, p<<1|1)
	}
}

// O(n) 从叶子向根方向update.
func (d *DivideInterval) PushUp(f func(parent, child1, child2 int32)) {
	for p := d.Offset - 1; p > 0; p-- {
		f(p, p<<1, p<<1|1)
	}
}

// 线段树的节点个数.
func (d *DivideInterval) Size() int32 {
	return d.Offset + d.n
}

func (d *DivideInterval) IsLeaf(segmentId int32) bool {
	return segmentId >= d.Offset
}

func (d *DivideInterval) Depth(u int32) int32 {
	if u == 0 {
		return 0
	}
	return int32(bits.LeadingZeros32(uint32(u))) - 1
}

// 线段树(完全二叉树)中两个节点的最近公共祖先(两个二进制数字的最长公共前缀).
func (d *DivideInterval) Lca(u, v int32) int32 {
	if u == v {
		return u
	}
	if u > v {
		u, v = v, u
	}
	depth1 := d.Depth(u)
	depth2 := d.Depth(v)
	diff := u ^ (v >> (depth2 - depth1))
	if diff == 0 {
		return u
	}
	len := bits.Len32(uint32(diff))
	return u >> len
}

func (d *DivideInterval) getSegmentIds(start, end int32) []int32 {
	if !(0 <= start && start <= end && end <= d.n) {
		return nil
	}
	var leftRes, rightRes []int32
	for start, end = start+d.Offset, end+d.Offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftRes = append(leftRes, start)
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = append(rightRes, end)
		}
	}
	for i := len(rightRes) - 1; i >= 0; i-- {
		leftRes = append(leftRes, rightRes[i])
	}
	return leftRes
}

//
//

// NewUnionFindWithCallback ...
func NewUnionFindArray(n int32) *_UnionFindArray {
	parent, rank := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		parent[i] = i
		rank[i] = 1
	}

	return &_UnionFindArray{
		Part:   n,
		rank:   rank,
		n:      n,
		parent: parent,
	}
}

type _UnionFindArray struct {
	// 连通分量的个数
	Part int32

	rank   []int32
	n      int32
	parent []int32
}

func (ufa *_UnionFindArray) Union(key1, key2 int32) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}

	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) UnionWithCallback(key1, key2 int32, preMerge func(big, small int32)) bool {
	root1, root2 := ufa.Find(key1), ufa.Find(key2)
	if root1 == root2 {
		return false
	}
	if ufa.rank[root1] > ufa.rank[root2] {
		root1, root2 = root2, root1
	}
	preMerge(root2, root1)
	ufa.parent[root1] = root2
	ufa.rank[root2] += ufa.rank[root1]
	ufa.Part--
	return true
}

func (ufa *_UnionFindArray) Find(key int32) int32 {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int32) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int32][]int32 {
	groups := make(map[int32][]int32)
	for i := int32(0); i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int32) int32 {
	return ufa.rank[ufa.Find(key)]
}

func (ufa *_UnionFindArray) String() string {
	sb := []string{"UnionFindArray:"}
	for root, member := range ufa.GetGroups() {
		cur := fmt.Sprintf("%d: %v", root, member)
		sb = append(sb, cur)
	}
	sb = append(sb, fmt.Sprintf("Part: %d", ufa.Part))
	return strings.Join(sb, "\n")
}

type ErasableHeapGeneric[H comparable] struct {
	data   *HeapGeneric[H]
	erased *HeapGeneric[H]
	size   int
}

func NewErasableHeapGeneric[H comparable](less func(a, b H) bool, nums ...H) *ErasableHeapGeneric[H] {
	return &ErasableHeapGeneric[H]{NewHeap(less, nums...), NewHeap(less), len(nums)}
}

// 从堆中删除一个元素,要保证堆中存在该元素.
func (h *ErasableHeapGeneric[H]) Erase(value H) {
	h.erased.Push(value)
	h.normalize()
	h.size--
}

func (h *ErasableHeapGeneric[H]) Push(value H) {
	h.data.Push(value)
	h.normalize()
	h.size++
}

func (h *ErasableHeapGeneric[H]) Pop() (value H) {
	value = h.data.Pop()
	h.normalize()
	h.size--
	return
}

func (h *ErasableHeapGeneric[H]) Peek() (value H) {
	value = h.data.Top()
	return
}

func (h *ErasableHeapGeneric[H]) Len() int {
	return h.size
}

func (h *ErasableHeapGeneric[H]) Clear() {
	h.data.Clear()
	h.erased.Clear()
	h.size = 0
}

func (h *ErasableHeapGeneric[H]) normalize() {
	for h.data.Len() > 0 && h.erased.Len() > 0 && h.data.Top() == h.erased.Top() {
		h.data.Pop()
		h.erased.Pop()
	}
}

type HeapGeneric[H comparable] struct {
	data []H
	less func(a, b H) bool
}

func NewHeap[H comparable](less func(a, b H) bool, nums ...H) *HeapGeneric[H] {
	nums = append(nums[:0:0], nums...)
	heap := &HeapGeneric[H]{less: less, data: nums}
	if len(nums) > 1 {
		heap.heapify()
	}
	return heap
}

func (h *HeapGeneric[H]) Push(value H) {
	h.data = append(h.data, value)
	h.pushUp(h.Len() - 1)
}

func (h *HeapGeneric[H]) Pop() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}

	value = h.data[0]
	h.data[0] = h.data[h.Len()-1]
	h.data = h.data[:h.Len()-1]
	h.pushDown(0)
	return
}

func (h *HeapGeneric[H]) Top() (value H) {
	if h.Len() == 0 {
		panic("heap is empty")
	}
	value = h.data[0]
	return
}

func (h *HeapGeneric[H]) Len() int { return len(h.data) }

func (h *HeapGeneric[H]) Clear() {
	h.data = h.data[:0]
}

func (h *HeapGeneric[H]) heapify() {
	n := h.Len()
	for i := (n >> 1) - 1; i > -1; i-- {
		h.pushDown(i)
	}
}

func (h *HeapGeneric[H]) pushUp(root int) {
	for parent := (root - 1) >> 1; parent >= 0 && h.less(h.data[root], h.data[parent]); parent = (root - 1) >> 1 {
		h.data[root], h.data[parent] = h.data[parent], h.data[root]
		root = parent
	}
}

func (h *HeapGeneric[H]) pushDown(root int) {
	n := h.Len()
	for left := (root<<1 + 1); left < n; left = (root<<1 + 1) {
		right := left + 1
		minIndex := root

		if h.less(h.data[left], h.data[minIndex]) {
			minIndex = left
		}

		if right < n && h.less(h.data[right], h.data[minIndex]) {
			minIndex = right
		}

		if minIndex == root {
			return
		}

		h.data[root], h.data[minIndex] = h.data[minIndex], h.data[root]
		root = minIndex
	}
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
