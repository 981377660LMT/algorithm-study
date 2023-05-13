// !当区间[l+d,r]的哈希值与[l,r-d]的哈希值相等时，那么该区间[l,r]是以 d 为循环节的**
// (前缀,后缀)
// 数学归纳法易证。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strings"
)

func main() {
	// 3 1 2
	// 112
	// 2 2 3 1
	// 1 1 3 8
	// 2 1 2 1
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int
	fmt.Fscan(in, &n, &m, &k)
	var s string
	fmt.Fscan(in, &s)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		nums[i] = int(s[i] - '0')
	}
	leaves := make([]E, n)
	for i := 0; i < n; i++ {
		leaves[i] = CreateE(uint(nums[i]))
	}

	seg := NewLazySegTree(leaves)
	for i := 0; i < m+k; i++ {
		var op, l, r, v int
		fmt.Fscan(in, &op, &l, &r, &v)
		l--
		if op == 1 { // RangeAssign
			seg.Update(l, r, uint(v))
		} else { // RangeQuery
			res1 := seg.Query(l+v, r)
			res2 := seg.Query(l, r-v)
			if res1.hash1 == res2.hash1 && res1.hash2 == res2.hash2 {
				fmt.Fprintln(out, "YES")
			} else {
				fmt.Fprintln(out, "NO")
			}
		}
	}
}

const INF uint = 1e18
const N int = 1e5 + 10

var BASEPOW0 [N]uint
var BASEPOW1 [N]uint
var BASEPOW2 [N]uint
var BASEPOW3 [N]uint

func init() {
	BASEPOW0[0] = 1
	BASEPOW1[0] = 1
	BASEPOW2[0] = 1
	BASEPOW3[0] = 1
	for i := 1; i < N; i++ {
		BASEPOW0[i] = BASEPOW0[i-1] * BASE0
		BASEPOW1[i] = BASEPOW1[i-1] * BASE1
		BASEPOW2[i] = BASEPOW2[i-1]*BASE0 + 1
		BASEPOW3[i] = BASEPOW3[i-1]*BASE1 + 1
	}
}

// !RangeAssignRangeHash

// 131/13331/1713302033171(回文素数)
const BASE0 = 131
const BASE1 = 13331

type E = struct {
	len          int
	hash1, hash2 uint
}

func CreateE(h uint) E {
	return E{len: 1, hash1: h, hash2: h}
}

type Id = uint

func (*LazySegTree) e() E   { return E{} }
func (*LazySegTree) id() Id { return INF }
func (*LazySegTree) op(a, b E) E {
	return E{
		len:   a.len + b.len,
		hash1: a.hash1*BASEPOW0[b.len] + b.hash1,
		hash2: a.hash2*BASEPOW1[b.len] + b.hash2,
	}
}

func (*LazySegTree) mapping(f Id, g E) E {
	if f == INF {
		return g
	}
	return E{len: g.len, hash1: f * BASEPOW2[g.len-1], hash2: f * BASEPOW3[g.len-1]}
}

func (*LazySegTree) composition(f, g Id) Id {
	if f == INF {
		return g
	}
	return f
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
