// RangeAddRangeSquareSum

// 区间加区间平方和

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	nums := make([]int, 6)
	leaves := make([]E, 6)
	for i := range nums {
		leaves[i] = FromElement(nums[i])
	}
	tree := NewSegmentTreeRangeAddRangeSquareSum(leaves)
	fmt.Println(tree.Query(2, 4))
	tree.Set(1, FromElement(-4))
	fmt.Println(tree.Query(3, 5))
	tree.Update(2, 6, -2)
	tree.Set(4, FromElement(4))
	fmt.Println(tree.Query(0, 6))
}

func test() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		leaves[i] = E{1, nums[i], nums[i] * nums[i]}
	}
	tree := NewSegmentTreeRangeAddRangeSquareSum(leaves)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r, v int
			fmt.Fscan(in, &l, &r, &v)
			l--
			tree.Update(l, r, v)
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			res := tree.Query(l, r)
			fmt.Fprintln(out, res.sum2)
		}
	}
}

// 1. 从左往右遍历数组，考虑新加入一个数 $nums[i]$ 会对左侧数组产生什么影响。
// 类似 2262，用一个哈希表 $last$ 记录上次 $nums[i]$ 出现的的位置 $j$(没有出现就是$-1$)，
// 那么左端点在区间 $[j+1,i+1)$ 内的子数组不同元素个数都会加一。
// 2. 用线段树维护区间平方和即可。
// 2916. 子数组不同元素数目的平方和 II
// https://leetcode.cn/problems/subarrays-distinct-element-sum-of-squares-ii/
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

const MOD int = 1e9 + 7
const INF = 1e18

// SegmentTreeRangeAddRangeSquareSum-区间加区间平方和

type E = struct{ sum0, sum1, sum2 int } // !0次和(size),1次和(sum),2次和(square sum)
type Id = int

func FromElement(v int) E { return E{1, v, v * v} }

func (*SegmentTreeRangeAddRangeSquareSum) e() E   { return E{} }
func (*SegmentTreeRangeAddRangeSquareSum) id() Id { return 0 }
func (*SegmentTreeRangeAddRangeSquareSum) op(left, right E) E {
	return E{left.sum0 + right.sum0, left.sum1 + right.sum1, left.sum2 + right.sum2}
}

func (*SegmentTreeRangeAddRangeSquareSum) mapping(f Id, g E) E {
	return E{g.sum0, g.sum1 + f*g.sum0, g.sum2 + 2*g.sum1*f + g.sum0*f*f}
}

func (*SegmentTreeRangeAddRangeSquareSum) composition(f, g Id) Id {
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
type SegmentTreeRangeAddRangeSquareSum struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewSegmentTreeRangeAddRangeSquareSum(leaves []E) *SegmentTreeRangeAddRangeSquareSum {
	tree := &SegmentTreeRangeAddRangeSquareSum{}
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
func (tree *SegmentTreeRangeAddRangeSquareSum) Query(left, right int) E {
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
func (tree *SegmentTreeRangeAddRangeSquareSum) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAddRangeSquareSum) Update(left, right int, f Id) {
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
func (tree *SegmentTreeRangeAddRangeSquareSum) MinLeft(right int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAddRangeSquareSum) MaxRight(left int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAddRangeSquareSum) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegmentTreeRangeAddRangeSquareSum) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeRangeAddRangeSquareSum) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *SegmentTreeRangeAddRangeSquareSum) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *SegmentTreeRangeAddRangeSquareSum) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *SegmentTreeRangeAddRangeSquareSum) String() string {
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
