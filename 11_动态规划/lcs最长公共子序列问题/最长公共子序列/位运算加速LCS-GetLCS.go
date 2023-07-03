// lcsbit/bitlcs
// https://atcoder.jp/contests/dp/submissions/34604402
// https://loj.ac/s/1633431

package main

import (
	"fmt"
	"strings"
)

func main() {
	bs := NewBS(100)
	bs.Set(1).Set(2).Set(3)
	fmt.Println(bs.Get(1), bs.Get(2), bs.Get(3))
	fmt.Println(bs)
	bs.Shift()
	fmt.Println(bs)

	// test minus
	bs1 := NewBS(100)
	bs1.Set(1).Set(2).Set(3)
	bs2 := NewBS(100)
	bs2.Set(1).Set(2).Set(3).Set(4)
	fmt.Println(Minus(bs2, bs1))
}

func GetLCS(nums1, nums2 []int) []int {
	n, m := len(nums1), len(nums2)
}

type BS struct {
	data []uint64
}

func NewBS(n int) *BS {
	return &BS{data: make([]uint64, 1+n>>6)}
}

func (bs *BS) Set(i int) *BS {
	bs.data[i>>6] |= 1 << (i & 63)
	return bs
}

func (bs *BS) Get(i int) bool {
	return bs.data[i>>6]&(1<<(i&63)) != 0
}

// 所有位向高位移动一位.
//
//	[1,2,3] -> [2,3,4]
func (bs *BS) Shift() *BS {
	last := uint64(0)
	for i := range bs.data {
		cur := (bs.data[i] >> 63) // 保留最高位
		bs.data[i] <<= 1          // 所有位向高位移动一位
		bs.data[i] |= last        // 将上一位的最高位放到当前位的最低位
		last = cur
	}
	return bs
}

func (bs *BS) SetAll(rhs *BS) *BS {
	copy(bs.data, rhs.data)
	return bs
}

func (bs *BS) AndAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] &= rhs.data[i]
	}
	return bs
}

func (bs *BS) OrAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] |= rhs.data[i]
	}
	return bs
}

func (bs *BS) XorAll(rhs *BS) *BS {
	for i := range bs.data {
		bs.data[i] ^= rhs.data[i]
	}
	return bs
}

func (bs *BS) String() string {
	sb := strings.Builder{}
	sb.Write([]byte{'['})
	nums := strings.Builder{}
	for i := range bs.data {
		start, end := i<<6, (i+1)<<6
		for j := start; j < end; j++ {
			if bs.Get(j) {
				nums.WriteString(fmt.Sprintf("%d,", j))
			}
		}
	}
	sb.WriteString(nums.String()[:nums.Len()-1])
	sb.Write([]byte{']'})
	return sb.String()
}

// Minus 计算两个集合的差集.
//
//	注意要用大的(lhs)减小的(rhs).
func Minus(lhs, rhs *BS) *BS {
	last := uint64(0)
	res := &BS{data: make([]uint64, len(lhs.data))}
	for i := range lhs.data {
		cur := (lhs.data[i] < (rhs.data[i] + last))
		res.data[i] = lhs.data[i] - rhs.data[i] - last
		if cur {
			last = 1
		} else {
			last = 0
		}
	}

	return res
}
