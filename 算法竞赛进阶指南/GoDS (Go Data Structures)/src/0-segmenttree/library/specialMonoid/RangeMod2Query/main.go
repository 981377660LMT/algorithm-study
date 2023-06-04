// RangeMod2Query(RangeMod2RangeAddRangeSum)
// https://yukicoder.me/problems/no/879
// 三种操作：
// 1 l r :如果ai是奇数,变为1,否则变为0.(Mod2)
// 2 l r x:区间[l,r]中的每个数都加上x.
// 3 l r:查询区间[l,r]的和.
//
// 线段树维护(奇数个数,偶数个数,和,大小)

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func RangeMod2Query(nums []int, queries [][]int) []int {
	leaves := make([]E, len(nums))
	for i, v := range nums {
		if nums[i]&1 == 1 {
			leaves[i] = E{sum: v, size: 1, odd: 1}
		} else {
			leaves[i] = E{sum: v, size: 1, even: 1}
		}
	}

	seg := NewLazySegTree(leaves)
	res := []int{}
	for _, q := range queries {
		op := q[0]
		if op == 1 {
			l, r := q[1], q[2]
			seg.Update(l, r, Id{kind: 2})
		} else if op == 2 {
			l, r, x := q[1], q[2], q[3]
			seg.Update(l, r, Id{add: x})
		} else {
			l, r := q[1], q[2]
			res = append(res, seg.Query(l, r).sum)
		}
	}

	return res
}

const INF = 1e18

// RangeMod2RangeAddRangeSum
type E = struct{ sum, size, odd, even int }
type Id = struct{ kind, add int }

func (*LazySegTree) e() E   { return E{} }
func (*LazySegTree) id() Id { return Id{} }
func (*LazySegTree) op(left, right E) E {
	return E{sum: left.sum + right.sum, size: left.size + right.size, odd: left.odd + right.odd, even: left.even + right.even}
}
func (*LazySegTree) mapping(f Id, g E) E {
	res := g
	add, mod := f.add, f.kind
	if mod == 1 {
		res.odd, res.even = res.even, res.odd
	}
	if mod != 0 {
		res.sum = res.odd
	}
	res.sum += add * res.size
	if add&1 == 1 {
		res.odd, res.even = res.even, res.odd
	}
	return res
}
func (*LazySegTree) composition(f, g Id) Id {
	res := g
	if f.kind != 0 {
		if res.kind == 0 {
			res.kind = 2
		}
		if f.kind == 1 {
			res.kind = 3 - res.kind
		}
		if res.add&1 == 1 {
			res.kind = 3 - res.kind
		}
		res.add = 0
	}
	res.add += f.add
	return res
}

// !template
type LazySegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewLazySegTree(leaves []E) *LazySegTree {
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
//   0<=left<=right<=len(tree.data)
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

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
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

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MinLeft(right int, predicate func(data E) bool) int {
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
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
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
func (tree *LazySegTree) MaxRight(left int, predicate func(data E) bool) int {
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
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
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
	tree.data[root] = tree.mapping(f, tree.data[root])
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
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	queries := make([][]int, q)
	for i := range queries {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			queries[i] = []int{op, l, r}
		} else if op == 2 {
			var l, r, x int
			fmt.Fscan(in, &l, &r, &x)
			l--
			queries[i] = []int{op, l, r, x}
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			queries[i] = []int{op, l, r}
		}
	}

	res := RangeMod2Query(nums, queries)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}

}
