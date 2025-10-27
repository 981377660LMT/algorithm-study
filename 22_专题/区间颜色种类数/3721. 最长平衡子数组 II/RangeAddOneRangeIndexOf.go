// RangeAddOneRangeIndexOf

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

// 3721. 最长平衡子数组 II
// https://leetcode.cn/problems/longest-balanced-subarray-ii/
// 给你一个整数数组 nums。
// 如果子数组中 不同偶数 的数量等于 不同奇数 的数量，则称该 子数组 是 平衡的 。
// 返回 最长 平衡子数组的长度。
// 子数组 是数组中连续且 非空 的一段元素序列。
func longestBalanced2(nums []int) int {
	n := len(nums)
	res := 0
	last := make(map[int]int)
	presum := NewRangeAddOneRangeIndexOf(n, func(i int) int { return 0 })
	cursum := 0
	for i := 1; i <= n; i++ {
		v := nums[i-1]
		b := v%2*2 - 1
		if j := last[v]; j == 0 {
			cursum += b
			presum.Update(i, n+1, b)
		} else {
			presum.Update(j, i, -b)
		}
		last[v] = i

		if tmp := presum.IndexOf(cursum, 0, i-res); tmp != -1 {
			if i-tmp > res {
				res = i - tmp
			}
		}
	}
	return res
}

type RangeAddOneRangeIndexOf struct {
	n   int
	seg *LazySegTree
}

func NewRangeAddOneRangeIndexOf(
	n int, f func(i int) int,
) *RangeAddOneRangeIndexOf {
	seg := NewLazySegTree(n, func(i int) E {
		v := f(i)
		return E{min: v, max: v}
	})
	res := &RangeAddOneRangeIndexOf{
		n:   n,
		seg: seg,
	}
	return res
}

// v=1/-1
func (riar *RangeAddOneRangeIndexOf) Update(start, end int, v int) {
	if start < 0 {
		start = 0
	}
	if end > riar.n {
		end = riar.n
	}
	if start >= end {
		return
	}
	riar.seg.Update(start, end, v)
}

// 返回区间[start,end)中第一个等于target的下标, 若不存在则返回-1.
func (riar *RangeAddOneRangeIndexOf) IndexOf(target int, start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > riar.n {
		end = riar.n
	}
	if start >= end {
		return -1
	}

	seg := riar.seg
	var dfs func(root, l, r int) int
	dfs = func(root, l, r int) int {
		if r <= start || end <= l {
			return -1
		}
		node := seg.data[root]
		if target < node.min || target > node.max {
			return -1
		}
		if r-l == 1 {
			if node.min == target {
				return l
			}
			return -1
		}
		seg.pushDown(root)
		mid := (l + r) >> 1
		if res := dfs(root<<1, l, mid); res != -1 {
			return res
		}
		return dfs(root<<1|1, mid, r)
	}
	return dfs(1, 0, seg.size)
}

const INF = 1e18

// RangeAddRangeMinMax
type E = struct{ min, max int }
type Id = int

func (*LazySegTree) e() E   { return E{min: INF, max: -INF} }
func (*LazySegTree) id() Id { return 0 }
func (*LazySegTree) op(left, right E) E {
	left.min = min(left.min, right.min)
	left.max = max(left.max, right.max)
	return left
}
func (*LazySegTree) mapping(f Id, g E, size int) E {
	if f == 0 {
		return g
	}
	g.min += f
	g.max += f
	return g
}
func (*LazySegTree) composition(f, g Id) Id {
	return f + g
}

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
