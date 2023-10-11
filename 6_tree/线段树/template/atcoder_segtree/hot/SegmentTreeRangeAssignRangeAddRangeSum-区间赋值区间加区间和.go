// SegmentTreeRangeAssignRangeAddRangeSum-区间赋值区间加区间和
// 相当于仿射变换的特殊情况
// 1. 区间赋值 -> {mul: 0, add: x}
// 2. 区间加 -> {mul: 1, add: x}

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	leaves := make([]E, len(nums))
	for i := 0; i < len(nums); i++ {
		leaves[i] = FromElement(nums[i])
	}
	tree := NewSegmentTreeRangeAffineRangeSum(leaves)
	fmt.Println(tree.Query(0, 10))
	tree.Update(0, 10, RangeAssignId(1))
	fmt.Println(tree.Query(0, 10))
	tree.Update(0, 10, RangeAddId(1))
	fmt.Println(tree.Query(0, 10))
}

const INF = 1e18

// RangeAssignRangeAddRangeSum

type E = struct{ sum, size int }
type Id = struct{ mul, add int }

func FromElement(v int) E    { return E{sum: v, size: 1} }
func RangeAssignId(v int) Id { return Id{mul: 0, add: v} }
func RangeAddId(v int) Id    { return Id{mul: 1, add: v} }

func (*SegmentTreeRangeAffineRangeSum) e() E   { return E{} }
func (*SegmentTreeRangeAffineRangeSum) id() Id { return Id{mul: 1} }
func (*SegmentTreeRangeAffineRangeSum) op(left, right E) E {
	return E{sum: left.sum + right.sum, size: left.size + right.size}
}
func (*SegmentTreeRangeAffineRangeSum) mapping(f Id, g E) E {
	return E{sum: f.mul*g.sum + f.add*g.size, size: g.size}
}
func (*SegmentTreeRangeAffineRangeSum) composition(f, g Id) Id {
	return Id{mul: f.mul * g.mul, add: (f.mul*g.add + f.add)}
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
type SegmentTreeRangeAffineRangeSum struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewSegmentTreeRangeAffineRangeSum(leaves []E) *SegmentTreeRangeAffineRangeSum {
	tree := &SegmentTreeRangeAffineRangeSum{}
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
func (tree *SegmentTreeRangeAffineRangeSum) Query(left, right int) E {
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
func (tree *SegmentTreeRangeAffineRangeSum) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAffineRangeSum) Update(left, right int, f Id) {
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
func (tree *SegmentTreeRangeAffineRangeSum) MinLeft(right int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAffineRangeSum) MaxRight(left int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAffineRangeSum) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegmentTreeRangeAffineRangeSum) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeRangeAffineRangeSum) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *SegmentTreeRangeAffineRangeSum) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *SegmentTreeRangeAffineRangeSum) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *SegmentTreeRangeAffineRangeSum) String() string {
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
