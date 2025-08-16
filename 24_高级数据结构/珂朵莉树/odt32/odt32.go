package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

func main() {
	odt := NewODT32(10, -1)
	odt.Set(0, 3, 1)
	odt.Set(3, 5, 2)
	fmt.Println(odt.Len, odt.Count, odt)
	odt.EnumerateAll(func(start, end int32, value *int) {
		fmt.Printf("[%d,%d):%d ", start, end, *value)
	})
	fmt.Println()
}

type ODT32[V comparable] struct {
	Len        int32 // 区间数
	Count      int32 // 区间元素个数之和
	llim, rlim int32
	noneValue  V
	data       []V
	ss         *FastSet32
}

// 指定区间长度 n 和哨兵 noneValue 建立一个 ODT.
//
//	区间为[0,n).
func NewODT32[V comparable](n int32, noneValue V) *ODT32[V] {
	res := &ODT32[V]{}
	dat := make([]V, n)
	for i := int32(0); i < n; i++ {
		dat[i] = noneValue
	}
	ss := NewFastSet32(n)
	ss.Insert(0)

	res.rlim = n
	res.noneValue = noneValue
	res.data = dat
	res.ss = ss
	return res
}

// 返回包含 x 的区间的信息.
func (odt *ODT32[V]) Get(x int32, erase bool) (start, end int32, value V) {
	start, end = odt.ss.Prev(x), odt.ss.Next(x+1)
	value = odt.data[start]
	if erase && value != odt.noneValue {
		odt.Len--
		odt.Count -= end - start
		odt.data[start] = odt.noneValue
		odt.Merge(start)
		odt.Merge(end)
	}
	return
}

func (odt *ODT32[V]) Set(start, end int32, value V) {
	odt.EnumerateRange(start, end, func(l, r int32, x *V) {}, true)
	odt.ss.Insert(start)
	odt.data[start] = value
	if value != odt.noneValue {
		odt.Len++
		odt.Count += end - start
	}
	odt.Merge(start)
	odt.Merge(end)
}

func (odt *ODT32[V]) EnumerateAll(f func(start, end int32, value *V)) {
	odt.EnumerateRange(0, odt.rlim, f, false)
}

// 遍历范围 [L, R) 内的所有数据.
func (odt *ODT32[V]) EnumerateRange(start, end int32, f func(start, end int32, value *V), erase bool) {
	if !(odt.llim <= start && start <= end && end <= odt.rlim) {
		panic(fmt.Sprintf("invalid range [%d, %d)", start, end))
	}

	NONE := odt.noneValue
	if !erase {
		l := odt.ss.Prev(start)
		for l < end {
			r := odt.ss.Next(l + 1)
			f(max32(l, start), min32(r, end), &odt.data[l])
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
		f(p, q, &x)
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

// 在位置 pos 处分割区间.
// 如果 pos 已经是区间的起始位置，则不进行分割.
func (odt *ODT32[V]) Split(pos int32) {
	if pos >= odt.rlim || pos <= odt.llim || odt.ss.Has(pos) {
		return
	}
	start := odt.ss.Prev(pos)
	odt.ss.Insert(pos)
	odt.data[pos] = odt.data[start]
	if odt.data[pos] != odt.noneValue {
		odt.Len++
	}
}

func (odt *ODT32[V]) Merge(p int32) {
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

func (odt *ODT32[V]) String() string {
	sb := []string{}
	odt.EnumerateAll(func(start, end int32, value *V) {
		var v interface{} = value
		if *value == odt.noneValue {
			v = "nil"
		}
		sb = append(sb, fmt.Sprintf("[%d,%d):%v", start, end, v))
	})
	return fmt.Sprintf("ODT{%v}", strings.Join(sb, ", "))
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

type FastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func NewFastSet32(n int32) *FastSet32 {
	res := &FastSet32{n: n}
	seg := [][]uint64{}
	n_ := n
	for {
		seg = append(seg, make([]uint64, (n_+63)>>6))
		n_ = (n_ + 63) >> 6
		if n_ <= 1 {
			break
		}
	}
	res.seg = seg
	res.lg = int32(len(seg))
	return res
}

func NewFastSet32From(n int32, f func(i int32) bool) *FastSet32 {
	res := NewFastSet32(n)
	for i := int32(0); i < n; i++ {
		if f(i) {
			res.seg[0][i>>6] |= 1 << (i & 63)
			res.size++
		}
	}
	for h := int32(0); h < res.lg-1; h++ {
		for i := 0; i < len(res.seg[h]); i++ {
			if res.seg[h][i] != 0 {
				res.seg[h+1][i>>6] |= 1 << (i & 63)
			}
		}
	}
	return res
}

func (fs *FastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *FastSet32) Insert(i int32) bool {
	if fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		fs.seg[h][i>>6] |= 1 << (i & 63)
		i >>= 6
	}
	fs.size++
	return true
}

func (fs *FastSet32) Erase(i int32) bool {
	if !fs.Has(i) {
		return false
	}
	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		cache[i>>6] &= ^(1 << (i & 63))
		if cache[i>>6] != 0 {
			break
		}
		i >>= 6
	}
	fs.size--
	return true
}

// 返回大于等于i的最小元素.如果不存在,返回n.
func (fs *FastSet32) Next(i int32) int32 {
	if i < 0 {
		i = 0
	}
	if i >= fs.n {
		return fs.n
	}

	for h := int32(0); h < fs.lg; h++ {
		cache := fs.seg[h]
		if i>>6 == int32(len(cache)) {
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
func (fs *FastSet32) Prev(i int32) int32 {
	if i < 0 {
		return -1
	}
	if i >= fs.n {
		i = fs.n - 1
	}

	for h := int32(0); h < fs.lg; h++ {
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
func (fs *FastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *FastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *FastSet32) Size() int32 {
	return fs.size
}

func (*FastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*FastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}
