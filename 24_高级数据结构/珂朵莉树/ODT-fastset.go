// 珂朵莉树(ODT)/Intervals

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"strconv"
	"strings"
)

const INF int = 1e18

func main() {
	// UnionOfInterval()
	CF1638E()
}

func demo() {
	odt := NewODT(10, -INF)
	odt.Set(0, 3, 1)
	odt.Set(3, 5, 2)
	fmt.Println(odt.Len, odt.Count, odt)
}

// https://www.luogu.com.cn/problem/CF1638E
// 给定一个长为n的数组.初始时每个元素的值为0，颜色为1.
// 进行q次操作:
// Color l r c : 将区间[l,r]内的元素的颜色变为c.
// Add c x : 将颜色为c的元素的值加上x.
// Query i : 查询第i个元素的值.
// 其中颜色为1-N.
//
// 操作2：只能维护懒标记.
// 操作1: 当颜色修改的时候，要把当前颜色之前没有加上的 tag 加上，需要一个数据结构来高效维护区间加，单点查。
// 操作3: 单点查。
func CF1638E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	lazyAdd := make([]int, n+1)
	bit := NewSimpleBitRangeAddPointGetArray(n + 1)
	odt := NewODT(n, -1)
	odt.Set(0, n, 1)

	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "Color" {
			var start, end, newColor int
			fmt.Fscan(in, &start, &end, &newColor)
			start--
			odt.EnumerateRange(
				start, end,
				func(l, r int, x int) { bit.AddRange(l, r, lazyAdd[x]) }, // 加回来
				true,
			)
			bit.AddRange(start, end, -lazyAdd[newColor]) // 覆盖完之后，当前的颜色之前加的数他并不需要加
			odt.Set(start, end, newColor)
		} else if op == "Add" {
			var color, x int
			fmt.Fscan(in, &color, &x)
			lazyAdd[color] += x
		} else if op == "Query" {
			var index int
			fmt.Fscan(in, &index)
			index--
			_, _, color := odt.Get(index, false)
			res := bit.Get(index) + lazyAdd[color]
			fmt.Fprintln(out, res)
		}
	}

}

func UnionOfInterval() {
	// https://atcoder.jp/contests/abc256/tasks/abc256_d
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	odt := NewODT(2e5+10, -INF)
	for i := 0; i < n; i++ {
		var l, r int
		fmt.Fscan(in, &l, &r)
		odt.Set(l, r, 1)
	}
	odt.EnumerateAll(func(l, r int, x int) {
		if x == 1 {
			fmt.Fprintln(out, l, r)
		}
	})
}

type Value = int

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
		odt.mergeAt(start)
		odt.mergeAt(end)
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
	odt.mergeAt(start)
	odt.mergeAt(end)
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

func (odt *ODT) mergeAt(p int) {
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

type SimpleBitRangeAddPointGetArray struct {
	bit *SimpleBit
}

func NewSimpleBitRangeAddPointGetArray(n int) *SimpleBitRangeAddPointGetArray {
	return &SimpleBitRangeAddPointGetArray{bit: NewSimpleBit(n)}
}

func (b *SimpleBitRangeAddPointGetArray) AddRange(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > b.bit.n {
		end = b.bit.n
	}
	if start >= end {
		return
	}
	b.bit.Add(start, delta)
	b.bit.Add(end, -delta)
}

func (b *SimpleBitRangeAddPointGetArray) Get(index int) int {
	return b.bit.QueryPrefix(index + 1)
}

// !Point Add Range Sum, 0-based.
type SimpleBit struct {
	n    int
	data []int
}

func NewSimpleBit(n int) *SimpleBit {
	res := &SimpleBit{n: n, data: make([]int, n)}
	return res
}

func (b *SimpleBit) Add(index int, v int) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *SimpleBit) QueryPrefix(end int) int {
	if end > b.n {
		end = b.n
	}
	res := 0
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *SimpleBit) QueryRange(start, end int) int {
	if start < 0 {
		start = 0
	}
	if end > b.n {
		end = b.n
	}
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := 0, 0
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
