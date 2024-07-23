// https://www.luogu.com.cn/problem/CF580E
// Kefa and Watch
// 1 l r c: 区间赋值, 将区间[l,r]的值全部赋值为c
// 2 l r d: 区间查询, 判断区间[l,r]是否以d为循环节
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
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m, k int32
	fmt.Fscan(in, &n, &m, &k)
	var s string
	fmt.Fscan(in, &s)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		nums[i] = int32(s[i] - '0')
	}

	seg := NewLazySegTree32(n, func(i int32) E { return FromElement(uint(nums[i])) })
	for i := int32(0); i < m+k; i++ {
		var op, l, r, v int32
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
const N int = 1e5 + 10 // TODO

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
	len          int32
	hash1, hash2 uint
}

func FromElement(h uint) E {
	return E{len: 1, hash1: h, hash2: h}
}

type Id = uint

func (*LazySegTree32) e() E   { return E{} }
func (*LazySegTree32) id() Id { return INF }
func (*LazySegTree32) op(a, b E) E {
	return E{
		len:   a.len + b.len,
		hash1: a.hash1*BASEPOW0[b.len] + b.hash1,
		hash2: a.hash2*BASEPOW1[b.len] + b.hash2,
	}
}

func (*LazySegTree32) mapping(f Id, g E) E {
	if f == INF {
		return g
	}
	return E{len: g.len, hash1: f * BASEPOW2[g.len-1], hash2: f * BASEPOW3[g.len-1]}
}

func (*LazySegTree32) composition(f, g Id) Id {
	if f == INF {
		return g
	}
	return f
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
	tree.data[root] = tree.mapping(f, tree.data[root])
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
