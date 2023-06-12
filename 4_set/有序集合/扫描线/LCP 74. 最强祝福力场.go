// https://leetcode.cn/problems/xepqZ5/
// 经过不断的勘测记录，小扣将所有力场的分布都记录了下来。
// !forceField[i] = [x,y,side] 表示第 i 片力场将覆盖以坐标 (x,y) 为中心，边长为 side 的正方形区域。
// 若任意一点的 力场强度 等于覆盖该点的力场数量，请求出在这片地带中 力场强度 最强处的 力场强度。
// 注意：
// 力场范围的边缘同样被力场覆盖。

// 按照x坐标排序，然后扫描线，线段树维护y坐标的最大值(RangeAddRangeMax)

package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
)

func main() {
	forceField := [][]int{{4, 4, 6}, {7, 5, 3}, {1, 6, 2}, {5, 6, 3}}
	fmt.Println(fieldOfGreatestBlessing(forceField))
}

func fieldOfGreatestBlessing(forceField [][]int) int {
	for i := range forceField {
		forceField[i][0] *= 2
		forceField[i][1] *= 2
	}

	ys := make([]int, 0, len(forceField)*2) // 离散化
	for _, v := range forceField {
		y, s := v[1], v[2]
		ys = append(ys, y-s, y+s)
	}
	rank, _ := sortedSet(ys)

	type event struct{ x, op, y1, y2 int } // 先进，op=-1，后进，op=1
	events := make([]event, 0, len(forceField)*2)
	for _, v := range forceField {
		x, y, s := v[0], v[1], v[2]
		y1, y2 := rank[y-s], rank[y+s]
		events = append(events, event{x - s, -1, y1, y2}, event{x + s, 1, y1, y2})
	}
	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		return events[i].op < events[j].op
	})

	res := 0
	leaves := make([]E, len(ys))
	tree := NewLazySegTree(leaves)
	for i := range events {
		op, y1, y2 := events[i].op, events[i].y1, events[i].y2
		if op == -1 {
			tree.Update(y1, y2+1, 1)
			res = max(res, tree.QueryAll())
		} else {
			tree.Update(y1, y2+1, -1)
		}
	}

	return res
}

// (紧)离散化.
//  rank: 给定一个在 nums 中的值,返回它的排名(0~len(rank)-1).
//  sorted: 离散化后的数组.
func sortedSet(nums []int) (rank map[int]int, sorted []int) {
	set := make(map[int]struct{})
	for _, v := range nums {
		set[v] = struct{}{}
	}
	allNums := make([]int, 0, len(set))
	for k := range set {
		allNums = append(allNums, k)
	}
	sort.Ints(allNums)
	rank = make(map[int]int, len(allNums))
	for i, v := range allNums {
		rank[v] = i
	}
	sorted = make([]int, len(nums))
	for i, v := range nums {
		sorted[i] = rank[v]
	}
	return
}

const INF = 1e18

// RangeAddRangeMax
type E = int
type Id = int

func (*LazySegTree) e() E   { return 0 }
func (*LazySegTree) id() Id { return 0 }
func (*LazySegTree) op(left, right E) E {
	return max(left, right)
}
func (*LazySegTree) mapping(f Id, g E) E {
	return f + g
}
func (*LazySegTree) composition(f, g Id) Id {
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
