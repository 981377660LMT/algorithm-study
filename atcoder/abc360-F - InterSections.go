// abc360-F - InterSections (最大相交区间数)
// 给定n个区间，形如[L[i], R[i]]，
// 区间i和区间j`相交`的定义为：L[i]<L[j]<R[i]<R[j]或L[j]<L[i]<R[j]<R[i]。
// 记f(l,r)为区间i和[l,r]相交的区间数.
// !对0<=l<r<=1e9，求出一个(l,r)满足f(l,r)最大.存在多个解时，输出(l,r)最小的解。
//
// !区间问题，考虑将区间放到二维平面上，变为点。
// 将区间放到二维平面上，变为点(L[i],R[i])。
// 区间[l,r]和第i个区间相交 => 两个不重合的矩形区域。
// ![0, L[i]) x [L[i]+1, R[i]) 和 [L[i]+1, R[i]) x (R[i]+1, +∞)
// 然后用扫描线+线段树，看值最大的区间。
//
// 注意[0,1e9]的区间，此时答案为[0,1].
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strings"
)

const INF32 int32 = 1e9

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	L, R := make([]int32, n), make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &L[i], &R[i])
	}

	// !离散化可能用到的L和R
	allX := func() []int32 {
		var res []int32
		for i := int32(0); i < n; i++ {
			res = append(res, L[i]-1, L[i], L[i]+1)
			res = append(res, R[i]-1, R[i], R[i]+1)
		}
		res = append(res, -1, 0, INF32, INF32+1)
		return res
	}()
	Unique2(&allX)

	// !扫描线，遍历x轴上的点，查询对应的y使得f(x,y)最大
	size := int32(len(allX))
	add, sub := make([][][2]int32, size), make([][][2]int32, size)
	for i := int32(0); i < n; i++ {
		a, b := L[i], R[i]
		l1 := BisectLeft(allX, 0)
		add[l1] = append(add[l1], [2]int32{a + 1, b})
		l2 := BisectLeft(allX, a)
		sub[l2] = append(sub[l2], [2]int32{a + 1, b})
		l3 := BisectLeft(allX, a+1)
		add[l3] = append(add[l3], [2]int32{b + 1, INF32 + 1})
		l4 := BisectLeft(allX, b)
		sub[l4] = append(sub[l4], [2]int32{b + 1, INF32 + 1})
	}

	res, resL, resR := int32(-1), int32(-1), int32(-1)
	seg := NewLazySegTree32(size, func(i int32) E { return E{index: i} })
	for i := int32(0); i < size; i++ {
		for _, lr := range add[i] {
			l, r := BisectLeft(allX, lr[0]), BisectLeft(allX, lr[1])
			seg.Update(l, r, 1)
		}
		for _, lr := range sub[i] {
			l, r := BisectLeft(allX, lr[0]), BisectLeft(allX, lr[1])
			seg.Update(l, r, -1)
		}
		// !题目要求0<=l<r<=1e9
		if !(0 <= allX[i] && allX[i] <= INF32) {
			continue
		}
		e := seg.Query(i+1, size)
		if e.index == -1 {
			continue
		}
		if e.value > res {
			res = e.value
			resL, resR = allX[i], allX[e.index]
		}
	}

	fmt.Fprintln(out, resL, resR)
}

// RangeAddRangeMaxIndex

type E = struct {
	index int32
	value int32
}
type Id = int32

func (*LazySegTree32) e() E   { return E{index: -1, value: -INF32} }
func (*LazySegTree32) id() Id { return 0 }
func (*LazySegTree32) op(left, right E) E {
	if left.value > right.value {
		return left
	}
	if left.value < right.value {
		return right
	}
	if left.index < right.index {
		return left
	}
	return right
}
func (*LazySegTree32) mapping(f Id, g E, size int) E {
	g.value += f
	return g
}
func (*LazySegTree32) composition(f, g Id) Id {
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
type LazySegTree32 struct {
	n    int32
	size int32
	log  int32
	data []E
	lazy []Id
}

func NewLazySegTree32(n int32, f func(int32) E) *LazySegTree32 {
	tree := &LazySegTree32{}
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewLazySegTree32From(leaves []E) *LazySegTree32 {
	tree := &LazySegTree32{}
	n := int32(len(leaves))
	tree.n = n
	tree.log = int32(bits.Len32(uint32(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := int32(0); i < n; i++ {
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
func (tree *LazySegTree32) Query(left, right int32) E {
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
func (tree *LazySegTree32) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree32) GetAll() []E {
	for i := int32(1); i < tree.size; i++ {
		tree.pushDown(i)
	}
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree32) Update(left, right int32, f Id) {
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
	for i := int32(1); i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree32) MinLeft(right int32, predicate func(data E) bool) int32 {
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
func (tree *LazySegTree32) MaxRight(left int32, predicate func(data E) bool) int32 {
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
func (tree *LazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree32) Set(index int32, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := int32(1); i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree32) pushUp(root int32) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree32) pushDown(root int32) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree32) propagate(root int32, f Id) {
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)) /**topbit**/)
	tree.data[root] = tree.mapping(f, tree.data[root], size)
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *LazySegTree32) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := int32(0); i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}

// Sorts the slice, removes duplicates, and clips the slice.
//
// T: cmp.Ordered
//
//	allX := []int{3, 1, 4, 1, 5, 9, 2, 6, 5, 3, 5}
//	Unique(&allX)
//	size := len(allX) // 7
//	newX := sort.SearchInts(allX, 5) // 4
func Unique2[T int | int32 | uint | uint32](nums *[]T) {
	value := *nums
	sort.Slice(value, func(i, j int) bool { return value[i] < value[j] })
	slow := 0
	for fast := 1; fast < len(value); fast++ {
		if value[fast] != value[slow] {
			slow++
			value[slow] = value[fast]
		}
	}
	value = value[: slow+1 : slow+1]
	*nums = value
}

// Find the index of the first element that is not less than x.
//
//	sort.SearchInts/LowerBound.
func BisectLeft[T int | int32](nums []T, x T) int32 {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < x {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}
