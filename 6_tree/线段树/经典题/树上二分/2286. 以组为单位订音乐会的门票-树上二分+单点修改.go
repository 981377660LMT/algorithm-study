package main

import (
	"math/bits"
)

type BookMyShow struct {
	row, col int
	tree     *LazySegTree
}

func Constructor(n int, m int) BookMyShow {
	nums := make([]E, n)
	for i := range nums {
		nums[i] = E{sum: m, max: m}
	}
	tree := NewLazySegTree(nums)
	return BookMyShow{row: n, col: m, tree: tree}
}

// !同一组的 k 位观众坐在 同一排座位，且座位连续 。
//  返回最小可能的 r 和 c 满足第 r 排中 [c, c + k - 1] 的座位都是空的，且 r <= maxRow.
//  如果 无法 安排座位，返回 [].
func (this *BookMyShow) Gather(k int, maxRow int) []int {
	first := this.tree.MaxRight(0, func(e E) bool { return e.max < k }) // !找到第一个空座位>=k的行
	if first > maxRow {
		return nil
	}
	used := this.col - this.tree.Query(first, first+1).sum
	this.tree.Update(first, first+1, -k)
	return []int{first, used}
}

// !k 位观众中 每一位 都有座位坐，但他们 不一定 坐在一起。
//  如果组里所有 k 个成员 不一定 要坐在一起的前提下，都能在第 0 排到第 maxRow 排之间找到座位，那么请返回 true.
//  这种情况下，每个成员都优先找排数 最小 ，然后是座位编号最小的座位。如果不能安排所有 k 个成员的座位，请返回 false 。
func (this *BookMyShow) Scatter(k int, maxRow int) bool {
	remain := this.tree.Query(0, maxRow+1).sum
	if remain < k {
		return false
	}

	first := this.tree.MaxRight(0, func(e E) bool { return e.sum == 0 }) // !找到第一个未坐满的行
	for k > 0 {
		remain := this.tree.Query(first, first+1).sum
		min_ := min(k, remain)
		this.tree.Update(first, first+1, -min_)
		k -= min_
		first++
	}

	return true
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

/**
 * Your BookMyShow object will be instantiated and called as such:
 * obj := Constructor(n, m);
 * param_1 := obj.Gather(k,maxRow);
 * param_2 := obj.Scatter(k,maxRow);
 */

type E = struct{ max, sum int } // 维护每行的剩余座位数和最大剩余座位数
type Id = int                   // delta

func (tree *LazySegTree) e() E   { return E{} }
func (tree *LazySegTree) id() Id { return 0 }
func (tree *LazySegTree) op(left, right E) E {
	return E{
		max: max(left.max, right.max),
		sum: left.sum + right.sum,
	}
}
func (tree *LazySegTree) mapping(lazy Id, data E) E { // !单点修改
	return E{
		max: data.max + lazy,
		sum: data.sum + lazy,
	}
}
func (tree *LazySegTree) composition(parentLazy, childLazy Id) Id {
	return parentLazy + childLazy
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
	v []E,
) *LazySegTree {
	tree := &LazySegTree{}

	n := int(len(v))
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = int(1) << tree.log
	tree.data = make([]E, 2*tree.size)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = v[i]
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
	for i := int(1); i <= tree.log; i++ {
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

// func main() {
// 	// 	["BookMyShow","gather","gather","scatter","scatter"]
// 	// [[2,5],[4,0],[2,0],[5,1],[5,1]]
// 	bookShow := Constructor(2, 5)
// 	fmt.Println(bookShow.Gather(4, 0))
// 	fmt.Println(bookShow.Gather(2, 0))
// 	fmt.Println(bookShow.tree.Query(0, 1), bookShow.tree.Query(1, 2))
// 	fmt.Println(bookShow.Scatter(5, 1))
// 	fmt.Println(bookShow.Scatter(5, 1))
// }
