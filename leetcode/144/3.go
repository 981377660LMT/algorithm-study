package main

import (
	"fmt"
	"math/bits"
	"sort"
	"strings"
)






func maxRemoval(nums []int, queries [][]int) int {
	n := len(nums)
	m := len(queries)

	counts := make([]int, n)
	diff := make([]int, n+1)
	for _, q := range queries {
		l, r := q[0], q[1]
		diff[l] += 1
		if r+1 < n {
			diff[r+1] -= 1
		}
	}
	counts[0] = diff[0]
	for i := 1; i < n; i++ {
		counts[i] = counts[i-1] + diff[i]
	}

	leaves := make([]int, n)
	for i := 0; i < n; i++ {
		if counts[i] < nums[i] {
			return -1
		}
		leaves[i] = counts[i] - nums[i]
	}

	st := NewLazySegTree32From(leaves)

	queryIndices := make([]int, m)
	for i := 0; i < m; i++ {
		queryIndices[i] = i
	}
	sort.Slice(queryIndices, func(i, j int) bool {
		a := queryIndices[i]
		b := queryIndices[j]
		li, ri := queries[a][0], queries[a][1]
		lj, rj := queries[b][0], queries[b][1]
		leni := ri - li + 1
		lenj := rj - lj + 1
		return leni < lenj
	})

	res := 0
	for _, idx := range queryIndices {
		l, r := int32(queries[idx][0]), int32(queries[idx][1]+1)
		min_ := st.Query(l, r)
		if min_ >= 1 {
			res++
			st.Update(l, r, -1)
		}
	}

	return res
}

const INF int = 1e18
const INF32 int32 = 1 << 30


type E = int

type Id = int

func (*LazySegTree32) e() E   { return INF }
func (*LazySegTree32) id() Id { return 0 }
func (*LazySegTree32) op(left, right E) E {
	return min(left, right)
}
func (*LazySegTree32) mapping(f Id, g E, size int) E {
	return f + g
}
func (*LazySegTree32) composition(f, g Id) Id {
	return f + g
}

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

func (tree *LazySegTree32) Get(index int32) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

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
	size := 1 << (tree.log - int32((bits.Len32(uint32(root)) - 1)))
	tree.data[root] = tree.mapping(f, tree.data[root], size)
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func mins(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num < res {
			res = num
		}
	}
	return res
}

func maxs(nums ...int) int {
	res := nums[0]
	for _, num := range nums {
		if num > res {
			res = num
		}
	}
	return res
}
