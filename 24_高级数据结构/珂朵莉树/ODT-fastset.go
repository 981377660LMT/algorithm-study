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

func demo() {
	odt := NewODT(10, -INF)
	odt.Set(0, 3, 1)
	odt.Set(3, 5, 2)
	fmt.Println(odt.Len, odt.Count, odt)
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

func main() {
	UnionOfInterval()
}

type Item = int

type ODT struct {
	Len        int // 区间数
	Count      int // 区间元素个数之和
	llim, rlim int
	noneValue  Item
	data       []Item
	ss         *fastSet
}

// 指定区间长度 n 和哨兵 noneValue 建立一个 ODT.
//  区间为[0,n).
func NewODT(n int, noneValue Item) *ODT {
	res := &ODT{}
	dat := make([]Item, n)
	for i := 0; i < n; i++ {
		dat[i] = noneValue
	}
	ss := newFastSet(n)
	ss.Insert(0)

	res.rlim = n
	res.noneValue = noneValue
	res.data = dat
	res.ss = ss
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODT) Get(x int, erase bool) (start, end int, value Item) {
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

func (odt *ODT) Set(start, end int, value Item) {
	odt.EnumerateRange(start, end, func(l, r int, x Item) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	odt.mergeAt(start)
	odt.mergeAt(end)
}

func (odt *ODT) EnumerateAll(f func(start, end int, value Item)) {
	odt.EnumerateRange(0, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODT) EnumerateRange(start, end int, f func(start, end int, value Item), erase bool) {
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
		odt.data[start] = odt.data[p]
		if odt.data[start] != NONE {
			odt.Len++
		}
	}
	p = odt.ss.Next(end)
	if end < p {
		odt.data[end] = odt.data[odt.ss.Prev(end)]
		odt.ss.Insert(end)
		if odt.data[end] != NONE {
			odt.Len++
		}
	}
	p = start
	for p < end {
		q := odt.ss.Next(p + 1)
		x := odt.data[p]
		f(p, q, x)
		if odt.data[p] != NONE {
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
	odt.EnumerateAll(func(start, end int, value Item) {
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
	if odt.data[p] == odt.data[q] {
		if odt.data[p] != odt.noneValue {
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

type fastSet struct {
	n, lg int
	seg   [][]int
}

func newFastSet(n int) *fastSet {
	res := &fastSet{n: n}
	seg := [][]int{}
	n_ := n
	for {
		seg = append(seg, make([]int, (n_+63)/64))
		n_ = (n_ + 63) / 64
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = len(seg)
	return res
}

func (fs *fastSet) Has(i int) bool {
	return (fs.seg[0][i/64]>>(i%64))&1 != 0
}

func (fs *fastSet) Insert(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] |= 1 << (i % 64)
		i /= 64
	}
}

func (fs *fastSet) Erase(i int) {
	for h := 0; h < fs.lg; h++ {
		fs.seg[h][i/64] &= ^(1 << (i % 64))
		if fs.seg[h][i/64] != 0 {
			break
		}
		i /= 64
	}
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *fastSet) Next(i int) int {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := 0; h < fs.lg; h++ {
		if i/64 == len(fs.seg[h]) {
			break
		}
		d := fs.seg[h][i/64] >> (i % 64)
		if d == 0 {
			i = i/64 + 1
			continue
		}
		// find
		i += fs.bsf(d)
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsf(fs.seg[g][i/64])
		}

		return i
	}

	return fs.n
}

// 返回小于等于i的最大元素.如果不存在,返回-1.
func (fs *fastSet) Prev(i int) int {
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
		d := fs.seg[h][i/64] << (63 - i%64)
		if d == 0 {
			i = i/64 - 1
			continue
		}
		// find
		i += fs.bsr(d) - (64 - 1)
		for g := h - 1; g >= 0; g-- {
			i *= 64
			i += fs.bsr(fs.seg[g][i/64])
		}

		return i
	}

	return -1
}

// 遍历[start,end)区间内的元素.
func (fs *fastSet) Enumerate(start, end int, f func(i int)) {
	x := start - 1
	for {
		x = fs.Next(x + 1)
		if x >= end {
			break
		}
		f(x)
	}
}

func (fs *fastSet) String() string {
	res := []string{}
	for i := 0; i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(i))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (*fastSet) bsr(x int) int {
	return 63 - bits.LeadingZeros(uint(x))
}

func (*fastSet) bsf(x int) int {
	return bits.TrailingZeros(uint(x))
}
