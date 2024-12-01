/* eslint-disable no-inner-declarations */

// !推荐使用，更快.
// 对每个下标，查询 最右侧/最左侧/右侧第一个/左侧第一个 lower/floor/ceiling/higher 的元素.
// 动态单调栈(DynamicMonoStack、MonoStackDynamic).
// 线段树实现.
// !左侧右侧包含当前位置.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	abc382c()
}

// C - Kaiten Sushi
// https://atcoder.jp/contests/abc382/tasks/abc382_c
// 回转寿司
//
// 在某个回转寿司店，有 N 人从 1 到 N 的编号到访。人 i 的美食度是 A i ​ 。
// 传送带上流动着 M 个寿司。第 j 个流动的寿司的美味程度是 B j ​ 。
// 每个寿司依次流过每个人 1,2,…,N 的面前。
// 每个人在美味程度>=自己美食度的寿司流到自己面前时会取走并食用该寿司，其他情况下则不做任何事情。
// 请确定每个寿司由谁食用，或者是否没有人食用。
//
// 等价于：对于每个寿司，找到第一个美食度大于等于它的人。
func abc382c() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var N, M int32
	fmt.Fscan(in, &N, &M)
	A, B := make([]int, N), make([]int, M)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := int32(0); i < M; i++ {
		fmt.Fscan(in, &B[i])
	}

	res := make([]int32, M)
	Q := NewRightMostLeftMostQuery(A)
	for i := int32(0); i < M; i++ {
		res[i] = int32(Q.RightNearestFloor(0, B[i]))
	}
	for _, v := range res {
		if v == -1 {
			fmt.Fprintln(out, -1)
		} else {
			fmt.Fprintln(out, v+1)
		}
	}
}

const INF int = 1e18

type RightMostLeftMostQuery struct {
	_n    int
	_tree *SegmentTreeRangeAddRangeMinMax
}

func NewRightMostLeftMostQuery(arr []int) *RightMostLeftMostQuery {
	n := len(arr)
	tree := NewSegmentTreeRangeAddRangeFrom(n, func(i int) E { return E{min: arr[i], max: arr[i]} })
	return &RightMostLeftMostQuery{_n: n, _tree: tree}
}

func (rm *RightMostLeftMostQuery) Get(index int) int {
	if index < 0 || index >= rm._n {
		panic(fmt.Sprintf("index %v out of range [0, %v)", index, rm._n))
	}
	return rm._tree.Get(index).min
}

func (rm *RightMostLeftMostQuery) Set(index int, value int) {
	if 0 <= index && index < rm._n {
		rm._tree.Set(index, E{min: value, max: value})
	}
}

func (rm *RightMostLeftMostQuery) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > rm._n {
		end = rm._n
	}
	if start >= end {
		return
	}
	rm._tree.Update(start, end, delta)
}

func (rm *RightMostLeftMostQuery) RightMostLower(index int, target int) int {
	cand := rm._tree.MinLeft(rm._n, func(e E) bool { return e.min >= target }) - 1
	if cand >= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightMostFloor(index int, target int) int {
	cand := rm._tree.MinLeft(rm._n, func(e E) bool { return e.min > target }) - 1
	if cand >= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightMostCeiling(index int, target int) int {
	cand := rm._tree.MinLeft(rm._n, func(e E) bool { return e.max < target }) - 1
	if cand >= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightMostHigher(index int, target int) int {
	cand := rm._tree.MinLeft(rm._n, func(e E) bool { return e.max <= target }) - 1
	if cand >= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) LeftMostLower(index int, target int) int {
	cand := rm._tree.MaxRight(0, func(e E) bool { return e.min >= target })
	if cand <= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) LeftMostFloor(index int, target int) int {
	cand := rm._tree.MaxRight(0, func(e E) bool { return e.min > target })
	if cand <= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) LeftMostCeiling(index int, target int) int {
	cand := rm._tree.MaxRight(0, func(e E) bool { return e.max < target })
	if cand <= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) LeftMostHigher(index int, target int) int {
	cand := rm._tree.MaxRight(0, func(e E) bool { return e.max <= target })
	if cand <= index {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightNearestLower(index int, target int) int {
	cand := rm._tree.MaxRight(index, func(e E) bool { return e.min >= target })
	if cand < rm._n {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightNearestFloor(index int, target int) int {
	cand := rm._tree.MaxRight(index, func(e E) bool { return e.min > target })
	if cand < rm._n {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightNearestCeiling(index int, target int) int {
	cand := rm._tree.MaxRight(index, func(e E) bool { return e.max < target })
	if cand < rm._n {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) RightNearestHigher(index int, target int) int {
	cand := rm._tree.MaxRight(index, func(e E) bool { return e.max <= target })
	if cand < rm._n {
		return cand
	}
	return -1
}

func (rm *RightMostLeftMostQuery) LeftNearestLower(index int, target int) int {
	cand := rm._tree.MinLeft(index+1, func(e E) bool { return e.min >= target }) - 1
	return cand
}

func (rm *RightMostLeftMostQuery) LeftNearestFloor(index int, target int) int {
	cand := rm._tree.MinLeft(index+1, func(e E) bool { return e.min > target }) - 1
	return cand
}

func (rm *RightMostLeftMostQuery) LeftNearestCeiling(index int, target int) int {
	cand := rm._tree.MinLeft(index+1, func(e E) bool { return e.max < target }) - 1
	return cand
}

func (rm *RightMostLeftMostQuery) LeftNearestHigher(index int, target int) int {
	cand := rm._tree.MinLeft(index+1, func(e E) bool { return e.max <= target }) - 1
	return cand
}

// RangeAddRangeMinMax-区间加区间最大最小值
type E = struct{ min, max int }
type Id = int

func (*SegmentTreeRangeAddRangeMinMax) e() E   { return E{min: INF, max: -INF} }
func (*SegmentTreeRangeAddRangeMinMax) id() Id { return 0 }
func (*SegmentTreeRangeAddRangeMinMax) op(left, right E) E {
	left.min = min(left.min, right.min)
	left.max = max(left.max, right.max)
	return left
}
func (*SegmentTreeRangeAddRangeMinMax) mapping(f Id, g E) E {
	if f == 0 {
		return g
	}
	g.min += f
	g.max += f
	return g
}
func (*SegmentTreeRangeAddRangeMinMax) composition(f, g Id) Id {
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
type SegmentTreeRangeAddRangeMinMax struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewSegmentTreeRangeAddRangeFrom(n int, f func(i int) E) *SegmentTreeRangeAddRangeMinMax {
	tree := &SegmentTreeRangeAddRangeMinMax{}
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
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAddRangeMinMax) Query(left, right int) E {
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
func (tree *SegmentTreeRangeAddRangeMinMax) QueryAll() E {
	return tree.data[1]
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *SegmentTreeRangeAddRangeMinMax) Update(left, right int, f Id) {
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
func (tree *SegmentTreeRangeAddRangeMinMax) MinLeft(right int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAddRangeMinMax) MaxRight(left int, predicate func(data E) bool) int {
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
func (tree *SegmentTreeRangeAddRangeMinMax) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *SegmentTreeRangeAddRangeMinMax) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeRangeAddRangeMinMax) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *SegmentTreeRangeAddRangeMinMax) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *SegmentTreeRangeAddRangeMinMax) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *SegmentTreeRangeAddRangeMinMax) String() string {
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
