// 更新:区间染色(update)
// 查询:区间和、区间最小值
// !思路:在lazy里用一个变量false表示是否更新过

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://www.luogu.com.cn/problem/CF1439C 饥饿的男人
	// !给定一个大小为 n 的`非升序`列 a 。现在有两类操作：
	// 1 x y 对 1<=i<=x 的元素取 max(ai,y) <=> 前缀`Chmax`
	// 2 x y 从下标x开始,从左到右遍历,如果ai<=y,则res+=1,y-=ai,最后输出res
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		leaves[i] = E{nums[i], 1, nums[i]}
	}
	tree := NewLazySegTree(leaves)

	for i := 0; i < q; i++ {
		var op, x, y int
		fmt.Fscan(in, &op, &x, &y)
		if op == 1 {
			// !前缀chMax的操作可以看成找到第一个小于等于y的位置,然后`区间更新`[first,x]的值为y
			right := tree.MaxRight(0, func(e E) bool { return e.minVal > y })
			tree.Update(right, x, Id{true, y})
		} else {
			x--
			res := 0
			for x < n {
				right := tree.MaxRight(x, func(e E) bool { return e.sum <= y })
				y -= tree.Query(x, right).sum
				res += right - x
				x = tree.MaxRight(right, func(e E) bool { return e.minVal > y })
			}
			fmt.Fprintln(out, res)
		}
	}
}

const INF int = 1e18

type E = struct{ sum, size, minVal int }
type Id = struct {
	updated bool
	target  int
}

func (*LazySegTree) e() E   { return E{minVal: INF} }
func (*LazySegTree) id() Id { return Id{target: INF} }
func (*LazySegTree) op(left, right E) E {
	return E{left.sum + right.sum, left.size + right.size, min(left.minVal, right.minVal)}
}
func (*LazySegTree) mapping(lazy Id, data E) E {
	if !lazy.updated {
		return data
	}
	return E{lazy.target * data.size, data.size, lazy.target}
}
func (*LazySegTree) composition(parentLazy, childLazy Id) Id {
	if !parentLazy.updated {
		return childLazy
	}
	return parentLazy
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

//
//
//
//
// !template
type LazySegTree struct {
	n    int
	log  int
	size int
	data []E
	lazy []Id
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
	tree.data[root] = tree.op(tree.data[2*root], tree.data[2*root+1])
}

func (tree *LazySegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(2*root, tree.lazy[root])
		tree.propagate(2*root+1, tree.lazy[root])
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
