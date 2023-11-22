package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

// No.1170 Never Want to Walk
// https://yukicoder.me/problems/no/1170
// 数轴上有n个车站,第i个位置在xi
// 如果两个车站之间的距离di与dj满足 A<=|di-dj|<=B,则这两个车站可以相互到达,否则不能相互到达
// 对每个车站,求出从该车站出发,可以到达的车站的数量
// 1<=n<=2e5 0<=A<=B<=1e9 0<=x1<=x2<...<=xn<=1e9
//
// !将序列搬到线段树上加速区间操作.
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, A, B int
	fmt.Fscan(in, &n, &A, &B)
	positions := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &positions[i])
	}

	D := NewDivideInterval(n)
	uf := NewUnionFindArray(D.Size())
	weights := make([]int, D.Size())
	for i := 0; i < n; i++ {
		weights[D.Id(i)] = 1
	}

	f := func(big, small int) {
		weights[big] += weights[small]
	}
	var dfs func(int) // 线段树上dfs
	dfs = func(cur int) {
		if D.IsLeaf(cur) {
			return
		}
		for k := 0; k < 2; k++ {
			child := cur<<1 | k
			if !uf.IsConnected(cur, child) {
				uf.UnionWithCallback(cur, child, f)
				dfs(child)
			}
		}
	}

	for i := 0; i < n; i++ {
		start := sort.SearchInts(positions, positions[i]+A)
		end := sort.SearchInts(positions, positions[i]+B+1)
		D.EnumerateSegment(start, end, func(segmentId int) {
			uf.UnionWithCallback(D.Id(i), segmentId, f)
			dfs(segmentId)
		})
	}

	for i := 0; i < n; i++ {
		fmt.Fprintln(out, weights[uf.Find(D.Id(i))])
	}
}

type DivideInterval struct {
	Offset int // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	n      int
}

// 线段树分割区间.
// 将长度为n的序列搬到长度为2*offset的线段树上, 以实现快速的区间操作.
func NewDivideInterval(n int) *DivideInterval {
	offset := 1
	for offset < n {
		offset <<= 1
	}
	return &DivideInterval{Offset: offset, n: n}
}

// 获取原下标为i的元素在树中的(叶子)编号.
func (d *DivideInterval) Id(rawIndex int) int {
	return rawIndex + d.Offset
}

// O(logn) 顺序遍历`[start,end)`区间对应的线段树节点.
func (d *DivideInterval) EnumerateSegment(start, end int, f func(segmentId int)) {
	if !(0 <= start && start <= end && end <= d.n) {
		panic("invalid range")
	}
	for _, i := range d.getSegmentIds(start, end) {
		f(i)
	}
}

// O(n) 从根向叶子方向push.
func (d *DivideInterval) PushDown(f func(parent, child int)) {
	for p := 1; p < d.Offset; p++ {
		f(p, p<<1)
		f(p, p<<1|1)
	}
}

// O(n) 从叶子向根方向update.
func (d *DivideInterval) PushUp(f func(parent, child1, child2 int)) {
	for p := d.Offset - 1; p > 0; p-- {
		f(p, p<<1, p<<1|1)
	}
}

// 线段树的节点个数.
func (d *DivideInterval) Size() int {
	return d.Offset + d.n
}

func (d *DivideInterval) IsLeaf(segmentId int) bool {
	return segmentId >= d.Offset
}

func (d *DivideInterval) Depth(u int) int {
	if u == 0 {
		return 0
	}
	return bits.Len(uint(u)) - 1
}

// 线段树(完全二叉树)中两个节点的最近公共祖先(两个二进制数字的最长公共前缀).
func (d *DivideInterval) Lca(u, v int) int {
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
	len := bits.Len(uint(diff))
	return u >> len
}

func (d *DivideInterval) getSegmentIds(start, end int) []int {
	if !(0 <= start && start <= end && end <= d.n) {
		panic("invalid range")
	}
	leftRes, rightRes := []int{}, []int{}
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
func NewUnionFindArray(n int) *_UnionFindArray {
	parent, rank := make([]int, n), make([]int, n)
	for i := 0; i < n; i++ {
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
	Part int

	rank   []int
	n      int
	parent []int
}

func (ufa *_UnionFindArray) Union(key1, key2 int) bool {
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

func (ufa *_UnionFindArray) UnionWithCallback(key1, key2 int, cb func(big, small int)) bool {
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
	cb(root2, root1)
	return true
}

func (ufa *_UnionFindArray) Find(key int) int {
	for ufa.parent[key] != key {
		ufa.parent[key] = ufa.parent[ufa.parent[key]]
		key = ufa.parent[key]
	}
	return key
}

func (ufa *_UnionFindArray) IsConnected(key1, key2 int) bool {
	return ufa.Find(key1) == ufa.Find(key2)
}

func (ufa *_UnionFindArray) GetGroups() map[int][]int {
	groups := make(map[int][]int)
	for i := 0; i < ufa.n; i++ {
		root := ufa.Find(i)
		groups[root] = append(groups[root], i)
	}
	return groups
}

func (ufa *_UnionFindArray) Size(key int) int {
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
