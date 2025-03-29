// F - Variety Split Hard
// 将数组分割成三个非空连续子数组，使得三个子数组中不同元素的数量之和最大。
// https://atcoder.jp/contests/abc397/tasks/abc397_f
// 给定一个长度为 N 的整数序列
//   A = (A₁, A₂, …, Aₙ)。
//
// 将 A 在两个位置切分成三个连续且非空的子数组。
// 对于满足 1 ≤ i < j ≤ N−1 的整数对 (i, j)，
// 将数组划分为三个部分：
//   (A₁, A₂, …, Aᵢ),
//   (Aᵢ₊₁, Aᵢ₊₂, …, Aⱼ),
//   (Aⱼ₊₁, Aⱼ₊₂, …, Aₙ)。
//
// 每个子数组中可能包含多个整数，定义“种类数”为子数组中不同整数的个数。
// 请问在所有可能的切分方案中，三个子数组的种类数之和能够取得的最大值是多少？
//
// 对于每个分割点 j（第二个分割点的位置），
// !我们需要找到最优的第一个分割点 i，使得 pre[i]+mid(i+1,j)+suf[j+1]的值最大。
package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	pool := make(map[int32]int32)
	id := func(o int32) int32 {
		if v, ok := pool[o]; ok {
			return v
		}
		v := int32(len(pool))
		pool[o] = v
		return v
	}

	var N int32
	fmt.Fscan(in, &N)
	A := make([]int32, N)
	for i := int32(0); i < N; i++ {
		fmt.Fscan(in, &A[i])
	}
	for i := int32(0); i < N; i++ {
		A[i] = id(A[i])
	}

	left := make([]int32, N+1)
	{
		counter := make([]int32, len(pool))
		k := int32(0)
		for i := int32(0); i < N; i++ {
			if counter[A[i]] == 0 {
				k++
			}
			counter[A[i]]++
			left[i+1] = k
		}
	}

	right := make([]int32, N+1)
	{
		counter := make([]int32, len(pool))
		k := int32(0)
		for i := N - 1; i >= 0; i-- {
			if counter[A[i]] == 0 {
				k++
			}
			counter[A[i]]++
			right[i] = k
		}
	}

	res := int32(0)
	seg := NewLazySegTree32(N, func(i int32) int32 { return 0 }) // 维护分割成两个子数组的最大种类数
	last := make([]int32, len(pool))
	for i := range last {
		last[i] = -1
	}
	for j := int32(0); j < N; j++ {
		cur := A[j]
		seg.Update(last[cur]+1, j, 1)
		last[cur] = j
		seg.Set(j, left[j]+1)
		tmp := seg.QueryAll() + right[j+1]
		res = max32(res, tmp)
	}

	fmt.Fprintln(out, res)
}

const INF32 int32 = 1 << 30

// RangeAddRangeMax

type E = int32

type Id = int32

func (*LazySegTree32) e() E   { return 0 }
func (*LazySegTree32) id() Id { return 0 }
func (*LazySegTree32) op(left, right E) E {
	return max32(left, right)
}
func (*LazySegTree32) mapping(f Id, g E, size int) E {
	return f + g
}
func (*LazySegTree32) composition(f, g Id) Id {
	return f + g
}
func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}
func max32(a, b int32) int32 {
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
