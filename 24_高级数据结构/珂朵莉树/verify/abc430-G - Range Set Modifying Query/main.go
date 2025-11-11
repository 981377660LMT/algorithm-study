// G - Range Set Modifying Query
// https://atcoder.jp/contests/abc430/tasks/abc430_g
//
// 有 N 个集合 S1..SN，初始都为空。
// 共 Q 个操作：
// 类型1：对区间 [L,R] 的每个集合 Si 加入元素 x。
// 类型2：对区间 [L,R] 的每个集合 Si 删除元素 x。
// 类型3：对区间 [L,R] 的集合，求其中元素个数的最大值，以及达到该最大值的集合个数。
// 限制：x 在 1..60。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	data := make([]*ODT, 61)
	for i := range data {
		data[i] = NewODT(n, false)
	}
	seg := NewLazySegTree(n, func(i int) E { return E{min: 0, minCount: 1} })

	add := func(l, r, x int) {
		data[x].EnumerateRange(
			l, r,
			func(start, end int, v Value) {
				if !v {
					seg.Update(start, end, -1)
				}
			},
			true,
		)
		data[x].Set(l, r, true)
	}

	remove := func(l, r, x int) {
		data[x].EnumerateRange(
			l, r,
			func(start, end int, v Value) {
				if v {
					seg.Update(start, end, 1)
				}
			},
			true,
		)
		data[x].Set(l, r, false)
	}

	query := func(l, r int) (int, int) {
		res := seg.Query(l, r)
		return -res.min, res.minCount
	}

	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var l, r, x int
			fmt.Fscan(in, &l, &r, &x)
			l--
			add(l, r, x)
		} else if op == 2 {
			var l, r, x int
			fmt.Fscan(in, &l, &r, &x)
			l--
			remove(l, r, x)
		} else {
			var l, r int
			fmt.Fscan(in, &l, &r)
			l--
			max, maxCount := query(l, r)
			fmt.Fprintln(out, max, maxCount)
		}
	}
}

type Value = bool

type ODT struct {
	Len        int // 区间数
	Count      int // 区间元素个数之和
	llim, rlim int
	noneValue  Value
	data       []Value
	ss         *_fastSet
}

// 指定区间长度 n 和哨兵 noneValue 建立一个 ODT.
//
//	区间为[0,n).
func NewODT(n int, noneValue Value) *ODT {
	res := &ODT{}
	dat := make([]Value, n)
	for i := 0; i < n; i++ {
		dat[i] = noneValue
	}
	ss := _newFastSet(n)
	ss.Insert(0)

	res.rlim = n
	res.noneValue = noneValue
	res.data = dat
	res.ss = ss
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODT) Get(x int, erase bool) (start, end int, value Value) {
	start, end = odt.ss.Prev(x), odt.ss.Next(x+1)
	value = odt.data[start]
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.data[start] = odt.noneValue
		// 尝试合并相邻的区间.start] = odt.noneValue
		odt.Merge(start)
		// 尝试合并相邻的区间.(start)
		odt.Merge(end)
	}
	return
}

func (odt *ODT) Set(start, end int, value Value) {
	odt.EnumerateRange(start, end, func(l, r int, x Value) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	// 尝试合并相邻的区间.
	odt.Merge(start)
	// 尝试合并相邻的区间.(start)
	odt.Merge(end)
}

func (odt *ODT) EnumerateAll(f func(start, end int, value Value)) {
	odt.EnumerateRange(0, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODT) EnumerateRange(start, end int, f func(start, end int, value Value), erase bool) {
	if !(odt.llim <= start && start <= end && end <= odt.rlim) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue
	if !erase {
		l := odt.ss.Prev(start)
		for l < end {
			r := odt.ss.Next(l + 1)
			f(max(l, start), min(r, end), odt.data[l])
			l = r
		}
		return
	}

	// 分割区间
	p := odt.ss.Prev(start)
	if p < start {
		odt.ss.Insert(start)
		v := odt.data[p]
		odt.data[start] = v
		if v != NONE {
			odt.Len++
		}
	}
	p = odt.ss.Next(end)
	if end < p {
		v := odt.data[odt.ss.Prev(end)]
		odt.data[end] = v
		odt.ss.Insert(end)
		if v != NONE {
			odt.Len++
		}
	}
	p = start
	for p < end {
		q := odt.ss.Next(p + 1)
		x := odt.data[p]
		f(p, q, x)
		if x != NONE {
			odt.Len--
			odt.Count -= q - p
		}
		odt.ss.Erase(p)
		p = q
	}
	odt.ss.Insert(start)
	odt.data[start] = NONE
}

// 在位置 pos 处分割区间，返回包含 pos 的区间的起始位置.
// 如果 pos 已经是区间的起始位置，则不进行分割.
func (odt *ODT) Split(pos int) int {
	if pos >= odt.rlim {
		return odt.rlim
	}
	if pos <= odt.llim {
		return odt.llim
	}
	if odt.ss.Has(pos) {
		return pos
	}
	start := odt.ss.Prev(pos)
	odt.ss.Insert(pos)
	odt.data[pos] = odt.data[start]
	if odt.data[pos] != odt.noneValue {
		odt.Len++
	}
	return pos
}

// 尝试合并相邻的区间.
func (odt *ODT) Merge(p int) {
	if p <= 0 || odt.rlim <= p {
		return
	}
	q := odt.ss.Prev(p - 1)
	if dataP, dataQ := odt.data[p], odt.data[q]; dataP == dataQ {
		if dataP != odt.noneValue {
			odt.Len--
		}
		odt.ss.Erase(p)
	}
}

func (odt *ODT) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int, value Value) {
		var v interface{} = value
		if value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
}

type _fastSet struct {
	n, lg int
	seg   [][]int
}

func _newFastSet(n int) *_fastSet {
	res := &_fastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *_fastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *_fastSet) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
}

func (fs *_fastSet) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] &= ^(1 << (i & 63))
		if fs.seg[h][i>>6] != 0 {
			break
		}
		i >>= 6
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *_fastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i>>6 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i>>6] >> (i & 63)
		if d == 0 {
			i = i>>6 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsf(fs.seg[g][i>>6])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *_fastSet) Prev(i int) int {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := 0; h < fs.lg; h++ {
		if i == -1 {
			break
		}
		d := fs.seg[h][i>>6] << (63 - i&63)
		if d == 0 {
			i = i>>6 - 1
			continue
		}
		// find
		i += fs.bsr(d) - 63
		for g := h - 1; g >= 0; g-- {
			i <<= 6
			i += fs.bsr(fs.seg[g][i>>6])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *_fastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *_fastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("_fastSet{%v}", strings.Join(res, ", "))
}

func (*_fastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*_fastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}

const INF = 1e18

// RangeAddRangeMinCount

type E = struct{ min, minCount int }
type Id = int

func (*LazySegTree) e() E   { return E{min: INF} }
func (*LazySegTree) id() Id { return 0 }
func (*LazySegTree) op(left, right E) E {
	if left.min < right.min {
		return left
	}
	if left.min > right.min {
		return right
	}
	left.minCount += right.minCount
	return left
}
func (*LazySegTree) mapping(f Id, g E, size int) E {
	if g.min == INF {
		return g
	}
	return E{min: g.min + f, minCount: g.minCount}
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

func NewLazySegTree(n int, f func(int) E) *LazySegTree {
	tree := &LazySegTree{}
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

func NewLazySegTreeFrom(leaves []E) *LazySegTree {
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
//
//	0<=left<=right<=len(tree.data)
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
func (tree *LazySegTree) GetAll() []E {
	for i := 1; i < tree.size; i++ {
		tree.pushDown(i)
	}
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
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
	size := 1 << (tree.log - (bits.Len32(uint32(root)) - 1) /**topbit**/)
	tree.data[root] = tree.mapping(f, tree.data[root], size)
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
