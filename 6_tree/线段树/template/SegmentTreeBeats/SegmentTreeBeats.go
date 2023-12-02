// https://rsm9.hatenablog.com/entry/2021/02/01/220408
// 吉老师线段树.
//
// !区别于普通线段树的地方:
// mapping函数返回的第二个参数表示是否更新成功.如果信息不足更新失败, 那么就需要递继续递归pushDown和pushUp.
// 注意叶子结点不能失败.
//
// 例子:
// 1.RangeChminChmaxAdd RangeSum, ChminChmax会导致Sum的计算失败.
// 2.RangeDivideAssignRangeSum, 如果区间中元素种类不唯一，会导致divide查询失败.

package main

import (
	"fmt"
	"math/bits"
	"strings"
)

func main() {

}

const INF int = 1e18
const MOD int = 1e9 + 7

type E = struct{ size, sum int }
type Id = struct{ mul, add int }

func (tree *SegmentTreeBeats) e() E   { return E{size: 1} }
func (tree *SegmentTreeBeats) id() Id { return Id{mul: 1} }
func (tree *SegmentTreeBeats) op(left, right E) E {
	return E{
		size: left.size + right.size,
		sum:  (left.sum + right.sum) % MOD,
	}
}

func (tree *SegmentTreeBeats) mapping(lazy Id, data E, size int) (newData E, success bool) {
	return E{
		size: data.size * lazy.mul % MOD,
		sum:  data.sum * lazy.mul % MOD,
	}, true
}

func (tree *SegmentTreeBeats) composition(parentLazy, childLazy Id) Id {
	return Id{
		mul: (parentLazy.mul * childLazy.mul) % MOD,
		add: (parentLazy.mul*childLazy.add + parentLazy.add) % MOD,
	}
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

// atcoder::lazy_segtree に1行書き足すだけの抽象化 Segment Tree Beats
// https://rsm9.hatenablog.com/entry/2021/02/01/220408
type SegmentTreeBeats struct {
	n, log, offset int
	data           []E
	lazy           []Id
}

func NewSegmentTreeBeats(leaves []E) *SegmentTreeBeats {
	res := &SegmentTreeBeats{}
	n := len(leaves)
	res.n = n
	res.log = 1
	for 1<<res.log < n {
		res.log++
	}
	res.offset = 1 << res.log
	res.data = make([]E, res.offset<<1)
	res.lazy = make([]Id, res.offset)
	for i := range res.data {
		res.data[i] = res.e()
	}
	for i := range res.lazy {
		res.lazy[i] = res.id()
	}
	for i := 0; i < n; i++ {
		res.data[res.offset+i] = leaves[i]
	}
	for i := res.offset - 1; i >= 1; i-- {
		res.pushUp(i)
	}
	return res
}

func (tree *SegmentTreeBeats) Get(index int) E {
	index += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

func (tree *SegmentTreeBeats) Set(index int, e E) {
	index += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *SegmentTreeBeats) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return tree.e()
	}
	start += tree.offset
	end += tree.offset
	for i := tree.log; i >= 1; i-- {
		if ((start >> i) << i) != start {
			tree.pushDown(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushDown((end - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for start < end {
		if start&1 != 0 {
			sml = tree.op(sml, tree.data[start])
			start++
		}
		if end&1 != 0 {
			end--
			smr = tree.op(tree.data[end], smr)
		}
		start >>= 1
		end >>= 1
	}
	return tree.op(sml, smr)
}

func (tree *SegmentTreeBeats) QueryAll() E {
	return tree.data[1]
}

func (tree *SegmentTreeBeats) Update(start, end int, f Id) {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return
	}
	start += tree.offset
	end += tree.offset
	for i := tree.log; i >= 1; i-- {
		if ((start >> i) << i) != start {
			tree.pushDown(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushDown((end - 1) >> i)
		}
	}
	l2, r2 := start, end
	for start < end {
		if start&1 != 0 {
			tree.propagate(start, f)
			start++
		}
		if end&1 != 0 {
			end--
			tree.propagate(end, f)
		}
		start >>= 1
		end >>= 1
	}
	start = l2
	end = r2
	for i := 1; i <= tree.log; i++ {
		if ((start >> i) << i) != start {
			tree.pushUp(start >> i)
		}
		if ((end >> i) << i) != end {
			tree.pushUp((end - 1) >> i)
		}
	}
}

func (tree *SegmentTreeBeats) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}
	right += tree.offset
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
			for right < tree.offset {
				tree.pushDown(right)
				right = right<<1 | 1
				if predicate(tree.op(tree.data[right], res)) {
					res = tree.op(tree.data[right], res)
					right--
				}
			}
			return right + 1 - tree.offset
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

func (tree *SegmentTreeBeats) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}
	left += tree.offset
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.offset {
				tree.pushDown(left)
				left <<= 1
				if predicate(tree.op(res, tree.data[left])) {
					res = tree.op(res, tree.data[left])
					left++
				}
			}
			return left - tree.offset
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

func (tree *SegmentTreeBeats) GetAll() []E {
	for i := 1; i < tree.offset; i++ {
		tree.pushDown(i)
	}
	res := tree.data[tree.offset : tree.offset+tree.n]
	return append(res[:0:0], res...)
}

func (tree *SegmentTreeBeats) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}

func (tree *SegmentTreeBeats) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}

func (tree *SegmentTreeBeats) propagate(root int, lazy Id) {
	size := 1 << (tree.log - topbit(root))
	var success bool
	tree.data[root], success = tree.mapping(lazy, tree.data[root], size)
	if root < tree.offset {
		tree.lazy[root] = tree.composition(lazy, tree.lazy[root])
		if !success { // !区别于普通线段树的地方.
			tree.pushDown(root)
			tree.pushUp(root)
		}
	}
}

func (tree *SegmentTreeBeats) String() string {
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

func topbit(x int) int {
	if x == 0 {
		return -1
	}
	return 31 - bits.LeadingZeros32(uint32(x))
}
