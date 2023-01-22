// !仿射变换 Range Affine Range Query
// 区间赋值 区间加 区间和查询

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]E, n)
	for i := 0; i < n; i++ {
		var num int
		fmt.Fscan(in, &num)
		nums[i] = E{size: 1, sum: num}
	}

	tree := NewLazySegTree(nums)

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var x int
			fmt.Fscan(in, &x)
			tree.Update(0, n, Id{mul: 0, add: x}) // 区间赋值
		} else if op == 2 {
			var i, x int
			fmt.Fscan(in, &i, &x)
			i--
			tree.Update(i, i+1, Id{mul: 1, add: x}) // 单点加
		} else {
			var i int
			fmt.Fscan(in, &i)
			i--
			fmt.Fprintln(out, tree.Query(i, i+1).sum) // 单点查询
		}
	}
}

type E = struct{ size, sum int }
type Id = struct{ mul, add int }

func (tree *LazySegTree) e() E   { return E{size: 1} }
func (tree *LazySegTree) id() Id { return Id{mul: 1} }
func (tree *LazySegTree) op(left, right E) E {
	return E{
		size: left.size + right.size,
		sum:  left.sum + right.sum,
	}
}

func (tree *LazySegTree) mapping(lazy Id, data E) E {
	return E{
		size: data.size,
		sum:  data.sum*lazy.mul + data.size*lazy.add,
	}
}

func (tree *LazySegTree) composition(parentLazy, childLazy Id) Id {
	return Id{
		mul: parentLazy.mul * childLazy.mul,
		add: parentLazy.mul*childLazy.add + parentLazy.add,
	}
}

func NewLazySegTree(
	leaves []E,
) *LazySegTree {
	tree := &LazySegTree{}

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

// !template
type LazySegTree struct {
	n    int
	log  int
	size int
	data []E
	lazy []Id
}

// 查询切片[left:right]的值
//   0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) E {
	if left == right {
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
		for left%2 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left *= 2
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left--
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

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	tree.propagate(2*root, tree.lazy[root])
	tree.propagate(2*root+1, tree.lazy[root])
	tree.lazy[root] = tree.id()
}

func (tree *LazySegTree) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}
