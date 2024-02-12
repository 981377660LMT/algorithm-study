// Update
// QueryRange
// QueryAll
// QueryPrefix
// MaxRight
// MaxRightWithIndex
// MinLeft
// Kth
// 树状数组树上二分

package main

import (
	"fmt"
	"math/bits"
)

func main() {
	bit := NewBITGroupFrom(5, func(index int) E { return E(index) })
	fmt.Println(bit)
	fmt.Println(bit.MaxRight(0, func(s E) bool { return s <= 100 })) // 5
	fmt.Println(bit.MaxRight(1, func(s E) bool { return s <= 3 }))   // 3
	fmt.Println(bit.MinLeft(5, func(s E) bool { return s <= 7 }))    // 3
}

// https://leetcode.cn/problems/longest-uploaded-prefix/description/
type LUPrefix struct {
	bit *BITGroup
}

func Constructor(n int) LUPrefix {
	return LUPrefix{bit: NewBITGroup(n + 1)}
}

func (this *LUPrefix) Upload(video int) {
	this.bit.Update(video-1, 1)
}

func (this *LUPrefix) Longest() int {
	return this.bit.MaxRightWithIndex(0, func(index int, sum E) bool { return sum == index })
}

/**
 * Your LUPrefix object will be instantiated and called as such:
 * obj := Constructor(n);
 * obj.Upload(video);
 * param_2 := obj.Longest();
 */
type E = int

func e() E           { return 0 }
func op(e1, e2 E) E  { return e1 + e2 }
func inv(e E) E      { return -e }    // 如果只查询前缀, 可以不需要是群
func mul(e E, k E) E { return e * k } // 如果不需要区间修改, 可以不需要是环

type BITGroup struct {
	n     int
	data  []E
	total E
}

func NewBITGroup(n int) *BITGroup {
	data := make([]E, n)
	for i := range data {
		data[i] = e()
	}
	return &BITGroup{n: n, data: data, total: e()}
}

func NewBITGroupFrom(n int, f func(index int) E) *BITGroup {
	total := e()
	data := make([]E, n)
	for i := range data {
		data[i] = f(i)
		total = op(total, data[i])
	}
	for i := 1; i <= n; i++ {
		j := i + (i & -i)
		if j <= n {
			data[j-1] = op(data[j-1], data[i-1])
		}
	}
	return &BITGroup{n: n, data: data, total: total}
}

func (fw *BITGroup) Update(i int, x E) {
	fw.total = op(fw.total, x)
	for i++; i <= fw.n; i += i & -i {
		fw.data[i-1] = op(fw.data[i-1], x)
	}
}

func (fw *BITGroup) QueryAll() E { return fw.total }

// [0, end)
func (fw *BITGroup) QueryPrefix(end int) E {
	if end > fw.n {
		end = fw.n
	}
	res := e()
	for end > 0 {
		res = op(res, fw.data[end-1])
		end &= end - 1
	}
	return res
}

// [start, end)
func (fw *BITGroup) QueryRange(start, end int) E {
	if start < 0 {
		start = 0
	}
	if end > fw.n {
		end = fw.n
	}
	if start == 0 {
		return fw.QueryPrefix(end)
	}
	if start > end {
		return e()
	}
	pos, neg := e(), e()
	for end > start {
		pos = op(pos, fw.data[end-1])
		end &= end - 1
	}
	for start > end {
		neg = op(neg, fw.data[start-1])
		start &= start - 1
	}
	return op(pos, inv(neg))
}

func (fw *BITGroup) MaxRight(start int, predicate func(sum E) bool) int {
	s := e()
	i := start
	getK := func() int {
		for {
			if i&1 == 1 {
				s = op(s, inv(fw.data[i-1]))
				i--
			}
			if i == 0 {
				return bits.Len32(uint32(fw.n))
			}
			k := bits.TrailingZeros32(uint32(i)) - 1
			if i+(1<<k) > fw.n {
				return k
			}
			t := op(s, fw.data[i+(1<<k)-1])
			if !predicate(t) {
				return k
			}
			s = op(s, inv(fw.data[i-1]))
			i -= i & -i
		}
	}
	k := getK()
	for k > 0 {
		k--
		if i+(1<<k)-1 < fw.n {
			t := op(s, fw.data[i+(1<<k)-1])
			if predicate(t) {
				i += 1 << k
				s = t
			}
		}
	}
	return i
}

// MaxRightWithIndex
func (fw *BITGroup) MaxRightWithIndex(start int, predicate func(index int, sum E) bool) int {
	s := e()
	i := start
	getK := func() int {
		for {
			if i&1 == 1 {
				s = op(s, inv(fw.data[i-1]))
				i--
			}
			if i == 0 {
				return bits.Len32(uint32(fw.n))
			}
			k := bits.TrailingZeros32(uint32(i)) - 1
			if i+(1<<k) > fw.n {
				return k
			}
			t := op(s, fw.data[i+(1<<k)-1])
			if !predicate(i+(1<<k), t) {
				return k
			}
			s = op(s, inv(fw.data[i-1]))
			i -= i & -i
		}
	}
	k := getK()
	for k > 0 {
		k--
		if i+(1<<k)-1 < fw.n {
			t := op(s, fw.data[i+(1<<k)-1])
			if predicate(i+(1<<k), t) {
				i += 1 << k
				s = t
			}
		}
	}
	return i
}

func (fw *BITGroup) MinLeft(end int, predicate func(sum E) bool) int {
	s := e()
	i := end
	k := 0
	for i > 0 && predicate(s) {
		s = op(s, fw.data[i-1])
		k = bits.TrailingZeros32(uint32(i))
		i -= i & -i
	}
	if predicate(s) {
		return 0
	}
	for k > 0 {
		k--
		t := op(s, inv(fw.data[i+(1<<k)-1]))
		if !predicate(t) {
			i += 1 << k
			s = t
		}
	}
	return i + 1
}

func (fw *BITGroup) Kth(k int, start int) int {
	return fw.MaxRight(start, func(x E) bool { return x <= k })
}

func (fw *BITGroup) String() string {
	res := []string{}
	for i := 0; i < fw.n; i++ {
		res = append(res, fmt.Sprintf("%d", fw.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITGroup%v", res)
}

type BITGroupRangeAdd struct {
	n    int
	bit0 *BITGroup
	bit1 *BITGroup
}

func NewBITGroupRangeAdd(n int) *BITGroupRangeAdd {
	return &BITGroupRangeAdd{n: n, bit0: NewBITGroup(n), bit1: NewBITGroup(n)}
}

func NewBITGroupRangeAddFrom(n int, f func(index int) E) *BITGroupRangeAdd {
	return &BITGroupRangeAdd{n: n, bit0: NewBITGroupFrom(n, f), bit1: NewBITGroup(n)}
}

func (fw *BITGroupRangeAdd) Update(index int, x E) {
	fw.bit0.Update(index, x)
}

func (fw *BITGroupRangeAdd) UpdateRange(start, end int, x E) {
	fw.bit0.Update(start, mul(x, -start))
	fw.bit0.Update(end, mul(x, end))
	fw.bit1.Update(start, x)
	fw.bit1.Update(end, inv(x))
}

func (fw *BITGroupRangeAdd) QueryRange(start, end int) E {
	rightRes := op(mul(fw.bit1.QueryPrefix(end), end), fw.bit0.QueryPrefix(end))
	leftRes := op(mul(fw.bit1.QueryPrefix(start), start), fw.bit0.QueryPrefix(start))
	return op(inv(leftRes), rightRes)
}

func (fw *BITGroupRangeAdd) String() string {
	res := []string{}
	for i := 0; i < fw.n; i++ {
		res = append(res, fmt.Sprintf("%d", fw.QueryRange(i, i+1)))
	}
	return fmt.Sprintf("BITGroupRangeAdd%v", res)
}
