// https://qoj.ac/problem/382
// 二进制数加法/减法, 二进制数第k位, 二进制数的符号.

package main

import (
	"fmt"
	"math/bits"
	"strconv"
	"strings"
)

func main() {
	b := NewBinaryNumber(100)
	b.Add(0, 100)
	fmt.Println(b.String())
	b.Add(0, 2333)
	b.Add(0, -233)
	fmt.Println(b.Kth(5))
	fmt.Println(b.Kth(7))
	fmt.Println(b.Kth(15))
	b.Add(15, 5)
	fmt.Println(b.Kth(15))
	b.Add(12, -1)
	fmt.Println(b.Kth(15))
}

type BinaryNumber struct {
	bit  int32
	data []int32
	set  *fastSet32
}

func NewBinaryNumber(bit int32) *BinaryNumber {
	return &BinaryNumber{bit: bit, data: make([]int32, bit), set: newFastSet32(bit)}
}

func (bn *BinaryNumber) Sgn() int32 {
	k := bn.set.Prev(bn.bit - 1)
	if k == -1 {
		return 0
	}
	return bn.data[k]
}

// 二进制第k位.
func (bn *BinaryNumber) Kth(k int32) int32 {
	j := bn.set.Prev(k - 1)
	x := bn.data[k]
	y := int32(0)
	if j != -1 {
		y = bn.data[j]
	}
	if x == 0 {
		if y >= 0 {
			return 0
		}
		return 1
	}
	if y >= 0 {
		return 1
	}
	return 0
}

// 加上 2^k * x
func (bn *BinaryNumber) Add(k int32, x int) {
	for x != 0 {
		x += int(bn.data[k])
		bn.data[k] = int32(x & 1)
		if bn.data[k] == 0 {
			bn.set.Erase(k)
		} else {
			bn.set.Insert(k)
		}
		k++
		x >>= 1
	}
}

func (bn *BinaryNumber) String() string {
	res := make([]string, bn.bit)
	for i, x := range bn.data {
		if x == 0 {
			res[i] = "0"
		} else if x == 1 {
			res[i] = "+"
		} else {
			res[i] = "-"
		}
	}
	return strings.Join(res, "")
}

type fastSet32 struct {
	n, lg int32
	seg   [][]uint64
	size  int32
}

func newFastSet32(n int32) *fastSet32 {
	res := &fastSet32{n: n}
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

func newFastSet32From(n int32, f func(i int32) bool) *fastSet32 {
	res := newFastSet32(n)
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

func (fs *fastSet32) Has(i int32) bool {
	return (fs.seg[0][i>>6]>>(i&63))&1 != 0
}

func (fs *fastSet32) Insert(i int32) bool {
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

func (fs *fastSet32) Erase(i int32) bool {
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
func (fs *fastSet32) Next(i int32) int32 {
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
func (fs *fastSet32) Prev(i int32) int32 {
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
func (fs *fastSet32) Enumerate(start, end int32, f func(i int32)) {
	for x := fs.Next(start); x < end; x = fs.Next(x + 1) {
		f(x)
	}
}

func (fs *fastSet32) String() string {
	res := []string{}
	for i := int32(0); i < fs.n; i++ {
		if fs.Has(i) {
			res = append(res, strconv.Itoa(int(i)))
		}
	}
	return fmt.Sprintf("FastSet{%v}", strings.Join(res, ", "))
}

func (fs *fastSet32) Size() int32 {
	return fs.size
}

func (*fastSet32) bsr(x uint64) int32 {
	return 63 - int32(bits.LeadingZeros64(x))
}

func (*fastSet32) bsf(x uint64) int32 {
	return int32(bits.TrailingZeros64(x))
}
