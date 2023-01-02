// https://judge.yosupo.jp/problem/point_set_range_composite
// 单点更新,区间聚合
// 0 index a b  => 点index变为 a*x + b
// 1 left right x  => 求从左到右 f(f(f(x))) mod 998244353的值

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const MOD int = 998244353

func main() {

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	initNums := make([]S, n)
	for i := range initNums {
		var a, b int
		fmt.Fscan(in, &a, &b)
		initNums[i] = S{mul: a, add: b}
	}
	tree := NewLazySegTree(initNums)

	for i := 0; i < q; i++ {
		var op, index, mul, add, left, right, x int
		fmt.Fscan(in, &op)
		if op == 0 {
			fmt.Fscan(in, &index, &mul, &add)
			tree.Update(index, index+1, F{mul: mul, add: add})
		} else {
			fmt.Fscan(in, &left, &right, &x)
			res := tree.Query(left, right)
			fmt.Fprintln(out, (res.mul*x+res.add)%MOD)
		}
	}

}

// !线段树维护的值的类型
type S = struct{ mul, add int }

// !更新操作的值的类型/懒标记的值的类型
type F = struct{ mul, add int }

// !线段树维护的值的幺元
//  alias: e
func (tree *LazySegTree) dataUnit() S { return S{mul: 1} }

// !更新操作/懒标记的幺元
//  alias: id
func (tree *LazySegTree) lazyUnit() F { return F{mul: 1} }

// !合并左右区间的值
//  alias: op
func (tree *LazySegTree) mergeChildren(left, right S) S {
	return S{
		mul: (left.mul * right.mul) % MOD,
		add: (left.add*right.mul + right.add) % MOD,
	}
}

// !父结点的懒标记更新子结点的值
//  alias: mapping
func (tree *LazySegTree) updateData(lazy F, data S) S {
	if lazy == tree.lazyUnit() {
		return data
	}
	return lazy
}

// !合并父结点的懒标记和子结点的懒标记
//  alias: composition
func (tree *LazySegTree) updateLazy(parentLazy, childLazy F) F {
	if parentLazy == tree.lazyUnit() {
		return childLazy
	}
	return parentLazy
}

func NewLazySegTree(
	leaves []S,
) *LazySegTree {
	tree := &LazySegTree{}

	n := int(len(leaves))
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
		tree.data[tree.size+i] = leaves[i]
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

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.mergeChildren(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	tree.propagate(2*root, tree.lazy[root])
	tree.propagate(2*root+1, tree.lazy[root])
	tree.lazy[root] = tree.lazyUnit()
}

func (tree *LazySegTree) propagate(root int, f F) {
	tree.data[root] = tree.updateData(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.updateLazy(f, tree.lazy[root])
	}
}
