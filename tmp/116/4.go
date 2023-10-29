package main

import (
	"math/bits"
)

// 前置问题： [2262. 字符串的总引力
// ](https://leetcode.cn/problems/total-appeal-of-a-string/solutions/1461618/by-endlesscheng-g405/)

// ---
// 1. 从左往右遍历数组，考虑新加入一个数 $nums[i]$ 会对左侧数组产生什么影响。类似 2262，用一个哈希表 $last$ 记录上次 $nums[i]$ 出现的的位置 $j$(没有出现就是$-1$)，那么左端点在区间 $[j+1,i+1)$ 内的子数组不同元素个数都会加一。
// 2. 用线段树维护区间平方和即可。基于 atcoder library 的 [lazy segtree](https://github.com/atcoder/ac-library/blob/master/atcoder/lazysegtree.hpp) , 幺半群实现如下：
// ![lazysegtree.png](https://pic.leetcode.cn/1698550695-vTBhCS-QQ%E5%9B%BE%E7%89%8720231029113758.png)

// ---
// 线段树只是工具，难基本都是难在写幺半群，除了这道题之外比较难的：
// - [2213. 由单个字符重复的最长子字符串](https://leetcode.cn/problems/longest-substring-of-one-repeating-character/)
// - [2286. 以组为单位订音乐会的门票
// ](https://leetcode.cn/problems/booking-concert-tickets-in-groups/)

func sumCounts(nums []int) int {
	n := len(nums)
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		leaves[i] = FromElement(0)
	}
	seg := NewSegmentTreeRangeAddRangeSquareSum(leaves)
	last := make(map[int]int)
	res := 0
	for i, num := range nums {
		pre := -1
		if v, ok := last[num]; ok {
			pre = v
		}
		seg.Update(pre+1, i+1, 1)
		last[num] = i
		res = (seg.Query(0, i+1).sum2 + res) % MOD
	}
	return res
}

const INF = 1e18
const MOD int = 1e9 + 7

// SegmentTreeRangeAddRangeSquareSum-区间加区间平方和

type E = struct{ sum0, sum1, sum2 int } // !0次和(size),1次和(sum),2次和(square sum)
type Id = int

func FromElement(v int) E { return E{1, v, v * v} }

func (*SegTree) e() E   { return E{} }
func (*SegTree) id() Id { return 0 }
func (*SegTree) op(left, right E) E {
	return E{left.sum0 + right.sum0, left.sum1 + right.sum1, left.sum2 + right.sum2}
}

func (*SegTree) mapping(f Id, g E) E {
	return E{g.sum0, g.sum1 + f*g.sum0, g.sum2 + 2*g.sum1*f + g.sum0*f*f}
}

func (*SegTree) composition(f, g Id) Id {
	return f + g
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// !template
type SegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewSegmentTreeRangeAddRangeSquareSum(leaves []E) *SegTree {
	tree := &SegTree{}
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
func (tree *SegTree) Query(left, right int) E {
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
func (tree *SegTree) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegTree) Update(left, right int, f Id) {
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
func (tree *SegTree) MinLeft(right int, predicate func(data E) bool) int {
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
func (tree *SegTree) MaxRight(left int, predicate func(data E) bool) int {
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
func (tree *SegTree) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegTree) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *SegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *SegTree) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}
