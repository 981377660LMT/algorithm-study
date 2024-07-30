package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

// RangeAssignRangeComposite
// https://judge.yosupo.jp/problem/range_set_range_composite
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	leaves := make([]E, n)
	for i := range leaves {
		var a, b int
		fmt.Fscan(in, &a, &b)
		leaves[i] = E{a, b}
	}

	seg := NewSegmentTreeRangeAssignFrom(n, func(i int) E { return leaves[i] })
	for i := 0; i < q; i++ {
		var kind int
		fmt.Fscan(in, &kind)
		if kind == 0 {
			var start, end, b, c int
			fmt.Fscan(in, &start, &end, &b, &c)
			seg.Assign(start, end, E{b, c})
		} else {
			var start, end, x int
			fmt.Fscan(in, &start, &end, &x)
			comp := seg.Query(start, end)
			fmt.Println(eval(comp, x))
		}
	}
}

// 区间赋值的线段树.
type SegmentTreeRangeAssign struct {
	n    int
	seg  *_SegmentTree
	cut  *_FastSet
	data []E
}

const INF int = 1e18

const MOD int = 998244353

type E = struct{ mul, add int }

const IS_COMMUTATIVE = false // 仿射变换群不满足交换律
func e() E {
	return E{1, 0}
}
func op(e1, e2 E) E {
	return E{e1.mul * e2.mul % MOD, (e1.add*e2.mul + e2.add) % MOD}
}
func inv(e E) E { // 仿射变换逆元
	mul, add := e.mul, e.add
	mul = modPow(mul, MOD-2, MOD) // modInv of mul
	return E{mul, mul * (MOD - add) % MOD}
}
func pow(e E, x int) E {
	// res := E{1, 0}
	resMul, resAdd := 1, 0
	eMul, eAdd := e.mul, e.add
	for x > 0 {
		if x&1 == 1 {
			// res = op(res, e)
			resMul, resAdd = resMul*eMul%MOD, (resAdd*eMul+eAdd)%MOD
		}
		// e = op(e, e)
		eMul, eAdd = eMul*eMul%MOD, (eAdd*eMul+eAdd)%MOD
		x >>= 1
	}
	return E{resMul, resAdd}
}
func eval(e E, x int) int {
	return (e.mul*x + e.add) % MOD
}

func modPow(x, n, mod int) int {
	res := 1
	for n > 0 {
		if n&1 == 1 {
			res = res * x % mod
		}
		x = x * x % mod
		n >>= 1
	}
	return res
}

// template
func NewSegmentTreeRangeAssign(n int) *SegmentTreeRangeAssign {
	return NewSegmentTreeRangeAssignFrom(n, func(i int) E { return e() })
}

func NewSegmentTreeRangeAssignFrom(n int, f func(i int) E) *SegmentTreeRangeAssign {
	res := &SegmentTreeRangeAssign{
		n:   n,
		seg: _NewSegmentTree(n, func(i int) E { return f(i) }),
		cut: _NewFastSetFrom(n, func(i int) bool { return true }),
	}
	res.data = res.seg.GetAll()
	return res
}

func (st *SegmentTreeRangeAssign) Query(start, end int) E {
	a, b, c := st.cut.Prev(start), st.cut.Next(start), st.cut.Prev(end)
	if a == c {
		return pow(st.data[a], end-start)
	}
	x := pow(st.data[a], b-start)
	y := st.seg.Query(b, c)
	z := pow(st.data[c], end-c)
	return op(op(x, y), z)
}

func (st *SegmentTreeRangeAssign) Assign(start, end int, value E) {
	a, b := st.cut.Prev(start), st.cut.Next(end)
	if a < start {
		st.seg.Set(a, pow(st.data[a], start-a))
	}
	if end < b {
		y := st.data[st.cut.Prev(end)]
		st.data[end] = y
		st.cut.Insert(end)
		st.seg.Set(end, pow(y, b-end))
	}
	st.cut.Enumerate(start+1, end, func(i int) { st.seg.Set(i, e()); st.cut.Erase(i) })
	st.data[start] = value
	st.cut.Insert(start)
	st.seg.Set(start, pow(value, end-start))
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

type _SegmentTree struct {
	n, size int
	seg     []E
}

func _NewSegmentTree(n int, f func(int) E) *_SegmentTree {
	res := &_SegmentTree{}
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = f(i)
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func _NewSegmentTreeFrom(leaves []E) *_SegmentTree {
	res := &_SegmentTree{}
	n := len(leaves)
	size := 1
	for size < n {
		size <<= 1
	}
	seg := make([]E, size<<1)
	for i := range seg {
		seg[i] = e()
	}
	for i := 0; i < n; i++ {
		seg[i+size] = leaves[i]
	}
	for i := size - 1; i > 0; i-- {
		seg[i] = op(seg[i<<1], seg[i<<1|1])
	}
	res.n = n
	res.size = size
	res.seg = seg
	return res
}
func (st *_SegmentTree) Get(index int) E {
	if index < 0 || index >= st.n {
		return e()
	}
	return st.seg[index+st.size]
}
func (st *_SegmentTree) Set(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = value
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}
func (st *_SegmentTree) Update(index int, value E) {
	if index < 0 || index >= st.n {
		return
	}
	index += st.size
	st.seg[index] = op(st.seg[index], value)
	for index >>= 1; index > 0; index >>= 1 {
		st.seg[index] = op(st.seg[index<<1], st.seg[index<<1|1])
	}
}

// [start, end)
func (st *_SegmentTree) Query(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > st.n {
		end = st.n
	}
	if start >= end {
		return e()
	}
	leftRes, rightRes := e(), e()
	start += st.size
	end += st.size
	for start < end {
		if start&1 == 1 {
			leftRes = op(leftRes, st.seg[start])
			start++
		}
		if end&1 == 1 {
			end--
			rightRes = op(st.seg[end], rightRes)
		}
		start >>= 1
		end >>= 1
	}
	return op(leftRes, rightRes)
}

func (st *_SegmentTree) QueryAll() E { return st.seg[1] }

func (st *_SegmentTree) GetAll() []E {
	res := make([]E, st.n)
	copy(res, st.seg[st.size:st.size+st.n])
	return res
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (st *_SegmentTree) MaxRight(left int, predicate func(E) bool) int {
	if left == st.n {
		return st.n
	}
	left += st.size
	res := e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(op(res, st.seg[left])) {
			for left < st.size {
				left <<= 1
				if tmp := op(res, st.seg[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - st.size
		}
		res = op(res, st.seg[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return st.n
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (st *_SegmentTree) MinLeft(right int, predicate func(E) bool) int {
	if right == 0 {
		return 0
	}
	right += st.size
	res := e()
	for {
		right--
		for right > 1 && right&1 == 1 {
			right >>= 1
		}
		if !predicate(op(st.seg[right], res)) {
			for right < st.size {
				right = right<<1 | 1
				if tmp := op(st.seg[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - st.size
		}
		res = op(st.seg[right], res)
		if right&-right == right {
			break
		}
	}
	return 0
}

type _FastSet struct {
	n, lg int
	seg   [][]int
}

func _NewFastSet(n int) *_FastSet {
	res := &_FastSet{n: n}
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

func _NewFastSetFrom(n int, f func(i int) bool) *_FastSet {
	res := _NewFastSet(n)
	for i := 0; i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
		}
	}
	for h := 0; h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *_FastSet) Has(i int) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *_FastSet) Insert(i int) bool {
	if fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	return true
}

func (fs *_FastSet) Erase(i int) bool {
	if !fs.Has(i) {
		return false
	}
	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *_FastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == len(cache) {
			break
		}
		d := cache[i>>6] >> (i & 63)
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
func (fs *_FastSet) Prev(i int) int {
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
func (fs *_FastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *_FastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *_FastSet) Size() int {
	return fs.n
}

func (*_FastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*_FastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
