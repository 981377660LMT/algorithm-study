// 更新:区间染色/区间gcd(ai,x)
// 查询:区间最大值、区间和
// RangeGcdRangeMaxSum
// https://hitonanode.github.io/cplib-cpp/segmenttree/trees/acl_range-update-gcd-range-max-sum.hpp
// https://rsm9.hatenablog.com/entry/2021/02/01/220408

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/880
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]E, n)
	for i := 0; i < n; i++ {
		var tmp int
		fmt.Fscan(in, &tmp)
		nums[i] = NewE(tmp)
	}

	seg := NewSegmentTreeBeats(nums)
	for i := 0; i < q; i++ {
		var t, l, r, x int
		fmt.Fscan(in, &t, &l, &r)
		l--
		if t <= 2 {
			fmt.Fscan(in, &x)
			if t == 1 {
				seg.Update(l, r, Id{updVal: x}) // RangeAssign
			} else {
				seg.Update(l, r, Id{gcdVal: x}) // RangeUpdateGcd
			}
		} else {
			res := seg.Query(l, r)
			if t == 3 {
				fmt.Fprintln(out, res.max) // RangeMax
			} else {
				fmt.Fprintln(out, res.sum) // RangeSum
			}
		}
	}

}

const INF int = 2e9

type E struct {
	sum, max, lcm, size int
	fail                bool
}

func NewE(v int) E {
	return E{
		sum:  v,
		max:  v,
		lcm:  v,
		size: 1,
	}
}

type Id struct{ updVal, gcdVal int }

func (*SegmentTreeBeats) e() E   { return E{lcm: 1} }
func (*SegmentTreeBeats) id() Id { return Id{} }
func (*SegmentTreeBeats) op(x, y E) E {
	return E{
		sum:  x.sum + y.sum,
		max:  max(x.max, y.max),
		lcm:  min(lcm(x.lcm, y.lcm), INF),
		size: x.size + y.size,
	}
}

func (*SegmentTreeBeats) mapping(f Id, x E) E {
	if f.updVal != 0 {
		return E{
			sum:  f.updVal * x.size,
			max:  f.updVal,
			lcm:  f.updVal,
			size: x.size,
		}
	}
	if f.gcdVal != 0 {
		if x.size == 1 {
			v := gcd(x.max, f.gcdVal)
			return E{
				sum:  v,
				max:  v,
				lcm:  v,
				size: 1,
			}
		} else if f.gcdVal%x.lcm != 0 {
			x.fail = true // !Special
		}
	}
	return x
}

func (*SegmentTreeBeats) composition(f, g Id) Id {
	if f.updVal != 0 {
		return f
	}
	if g.updVal != 0 {
		return Id{updVal: gcd(g.updVal, f.gcdVal)}
	}
	return Id{gcdVal: gcd(f.gcdVal, g.gcdVal)}
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

func gcd(a, b int) int {
	for b != 0 {
		a, b = b, a%b
	}
	return a
}

func lcm(a, b int) int {
	return (a / gcd(a, b)) * b
}

//
//
//
//
// !template
type SegmentTreeBeats struct {
	n    int
	log  int
	size int
	data []E
	lazy []Id
}

func NewSegmentTreeBeats(
	leaves []E,
) *SegmentTreeBeats {
	tree := &SegmentTreeBeats{}
	n := int(len(leaves))
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, 2*tree.size)
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
func (tree *SegmentTreeBeats) Query(left, right int) E {
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

func (tree *SegmentTreeBeats) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *SegmentTreeBeats) Update(left, right int, f Id) {
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
func (tree *SegmentTreeBeats) MinLeft(right int, predicate func(data E) bool) int {
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
				right = right*2 + 1
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
func (tree *SegmentTreeBeats) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}

	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}

	res := tree.e()
	for {
		for left%2 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left *= 2
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
func (tree *SegmentTreeBeats) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegmentTreeBeats) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeBeats) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[2*root], tree.data[2*root+1])
}

func (tree *SegmentTreeBeats) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(2*root, tree.lazy[root])
		tree.propagate(2*root+1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}

func (tree *SegmentTreeBeats) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
		// !Special
		if tree.data[root].fail {
			tree.pushDown(root)
			tree.pushUp(root)
		}
	}
}
