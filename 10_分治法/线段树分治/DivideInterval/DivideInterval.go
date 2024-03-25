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
	abc342g()
	// SP11470()
	// yuki1170()
}

type InnerTreeRangeSum struct {
	data   []int
	preSum []int
}

func NewInnerTreeRangeSum() *InnerTreeRangeSum {
	return &InnerTreeRangeSum{}
}

func (stl *InnerTreeRangeSum) Add(x int) {
	stl.data = append(stl.data, x)
}

func (stl *InnerTreeRangeSum) Build() {
	sort.Ints(stl.data)
	stl.preSum = make([]int, len(stl.data)+1)
	for i, x := range stl.data {
		stl.preSum[i+1] = stl.preSum[i] + x
	}
}

// 小于等于upper的元素之和.
func (stl *InnerTreeRangeSum) Query(upper int) int {
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
	innerTree := make([]*InnerTreeRangeSum, D.Size())
	for i := range innerTree {
		innerTree[i] = NewInnerTreeRangeSum()
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

type InnerTreeAbc342g struct {
}

// G - Retroactive Range Chmax (可追溯区间最大值修改)
// https://atcoder.jp/contests/abc342/tasks/abc342_g
// 维护一个数列，有以下三个操作：
// 1 l r x: 将区间[l,r]中的所有元素与x取最大值.
// 2 i: 将第i次操作删除，保证第i次操作是操作1.
// 3 i: 查询当前数列中第i个元素的值.
func abc342g() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
}

// https://www.luogu.com.cn/problem/SP11470
func SP11470() {

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
