package segmenttree

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

// https://atcoder.jp/contests/practice2/tasks/practice2_l
func demo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]S, n)
	for i := range nums {
		var num int
		fmt.Fscan(in, &num)
		nums[i] = S{Zero: 1 - num, One: num}
	}

	tree := NewLazySegTree(nums)
	for i := 0; i < q; i++ {
		var op uint8
		var left, right int
		fmt.Fscan(in, &op, &left, &right)
		left--
		if op == 1 {
			tree.Update(left, right, true)
		} else {
			fmt.Fprintln(out, tree.Query(left, right).Inversion)
		}
	}
}

// !线段树维护的值的类型
type S struct {
	Zero, One, Inversion int
}

// !更新操作的值的类型/懒标记的值的类型
type F bool

// !线段树维护的值的幺元.
//  alias: e
func (tree *LazySegTree) dataUnit() S { return S{} }

// !更新操作/懒标记的幺元
//  alias: id
func (tree *LazySegTree) lazyUnit() F { return false }

// !合并左右区间的值
//  alias: op
func (tree *LazySegTree) mergeChildren(left, right S) S {
	return S{
		Zero:      left.Zero + right.Zero,
		One:       left.One + right.One,
		Inversion: left.Inversion + right.Inversion + left.One*right.Zero,
	}
}

// !父结点的懒标记更新子结点的值
//  alias: mapping
func (tree *LazySegTree) updateData(lazy F, data S) S {
	if !lazy {
		return data
	}
	return S{
		Zero:      data.One,
		One:       data.Zero,
		Inversion: data.One*data.Zero - data.Inversion,
	}
}

// !合并父结点的懒标记和子结点的懒标记
//  alias: composition
func (tree *LazySegTree) updateLazy(parentLazy, childLazy F) F {
	return (parentLazy && !childLazy) || (!parentLazy && childLazy)
}

func NewLazySegTree(
	v []S,
) *LazySegTree {
	tree := &LazySegTree{}

	n := int(len(v))
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]S, 2*tree.size)
	tree.lazy = make([]F, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.dataUnit()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.lazyUnit()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = v[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// !template
type LazySegTree struct {
	n    int
	log  int
	size int
	data []S
	lazy []F
}

// 查询切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) S {
	if left == right {
		return tree.dataUnit()
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
	sml, smr := tree.dataUnit(), tree.dataUnit()
	for left < right {
		if left&1 != 0 {
			sml = tree.mergeChildren(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.mergeChildren(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.mergeChildren(sml, smr)
}

func (tree *LazySegTree) QueryAll() S {
	return tree.data[1]
}

// 更新切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f F) {
	if left == right {
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
			tree.allApply(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.allApply(right, f)
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
func (tree *LazySegTree) MinLeft(right int, predicate func(data S) bool) int {
	if right == 0 {
		return 0
	}

	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}

	res := tree.dataUnit()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}

		if !predicate(tree.mergeChildren(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right*2 + 1
				if predicate(tree.mergeChildren(tree.data[right], res)) {
					res = tree.mergeChildren(tree.data[right], res)
					right--
				}
			}

			return right + 1 - tree.size
		}

		res = tree.mergeChildren(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}

	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MaxRight(left int, predicate func(data S) bool) int {
	if left == tree.n {
		return tree.n
	}

	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}

	res := tree.dataUnit()
	for {
		for left%2 == 0 {
			left >>= 1
		}
		if !predicate(tree.mergeChildren(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left *= 2
				if predicate(tree.mergeChildren(res, tree.data[left])) {
					res = tree.mergeChildren(res, tree.data[left])
					left--
				}
			}

			return left - tree.size
		}

		res = tree.mergeChildren(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}

	return tree.n
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.mergeChildren(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	tree.allApply(2*root, tree.lazy[root])
	tree.allApply(2*root+1, tree.lazy[root])
	tree.lazy[root] = tree.lazyUnit()
}

func (tree *LazySegTree) allApply(root int, f F) {
	tree.data[root] = tree.updateData(f, tree.data[root])
	if root < tree.size {
		tree.lazy[root] = tree.updateLazy(f, tree.lazy[root])
	}
}
