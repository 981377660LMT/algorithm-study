// Usage:
// !Attention: nums[i] >= 0; log 一般要开大一点.

// Count(start, end, a, b, xor) - 区间 [start, end) 中值在 [a, b) 之间的数的个数
// CountPrefix(start, end, x, xor) - 区间 [start, end) 中值在 [0, x) 之间的数的个数

// Kth(start, end, k, xor) - 区间 [start, end) 中第 k 小的数(0-indexed)

// Floor(start, end, x, xor) - 区间 [start, end) 中值小于等于 x 的最大值
// Ceil(start, end, x, xor) - 区间 [start, end) 中值大于等于 x 的最小值

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	M := NewWaveletMatrixXor(nums, 32)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var left, right, x int
		fmt.Fscan(in, &left, &right, &x)
		left--
		res := INF
		lower := M.Floor(left, right, x, 0) // 小于等于x的最大值
		if lower != -INF {
			res = min(res, abs(lower-x))
		}
		higher := M.Ceil(left, right, x, 0) // 大于等于x的最小值
		if higher != INF {
			res = min(res, abs(higher-x))
		}

		fmt.Fprintln(out, res)
	}

}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

const INF int = 1e18

type WaveletMatrixXor struct {
	n, log int
	mid    []int
	bv     []*BitVector
}

// log:所有数不超过 2^log-1 , 1e9时一般设为32
//  如果要支持异或,则需要按照异或的值来决定值域
func NewWaveletMatrixXor(data []int, log int) *WaveletMatrixXor {
	data = append(data[:0:0], data...)
	res := &WaveletMatrixXor{}
	n := len(data)
	mid := make([]int, log)
	bv := make([]*BitVector, log)
	for i := 0; i < log; i++ {
		bv[i] = NewBitVector(n)
	}
	a0, a1 := make([]int, n), make([]int, n)
	for d := log - 1; d >= 0; d-- {
		p0, p1 := 0, 0
		for i := 0; i < n; i++ {
			f := (data[i] >> d) & 1
			if f == 0 {
				a0[p0] = data[i]
				p0++
			} else {
				bv[d].Set(i)
				a1[p1] = data[i]
				p1++
			}
		}
		mid[d] = p0
		bv[d].Build()
		data, a0 = a0, data
		for i := 0; i < p1; i++ {
			data[p0+i] = a1[i]
		}
	}

	res.n = n
	res.log = log
	res.mid = mid
	res.bv = bv
	return res
}

// [left, right) 中位于 `[a,b)` 的数的個数
func (wm *WaveletMatrixXor) Count(left, right, a, b, xor int) int {
	return wm.CountPrefix(left, right, b, xor) - wm.CountPrefix(left, right, a, xor)
}

// [left, right) 中位于 `[0,x)` 的数的個数
func (wm *WaveletMatrixXor) CountPrefix(left, right, x, xor int) int {
	if x >= 1<<wm.log {
		return right - left
	}
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		add := x >> d & 1
		f := (x ^ xor) >> d & 1
		if add != 0 {
			res += wm.bv[d].Rank(right, f^1) - wm.bv[d].Rank(left, f^1)
		}
		left = wm.bv[d].Rank(left, f) + (f * wm.mid[d])
		right = wm.bv[d].Rank(right, f) + (f * wm.mid[d])
	}
	return res
}

// [left, right) 中的第 k 小的數(k>=0)
//  如果不存在則返回-1
func (wm *WaveletMatrixXor) Kth(left, right, k, xor int) int {
	if k < 0 || k >= right-left {
		return -1
	}
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0 := wm.bv[d].Rank(left, 0)
		r0 := wm.bv[d].Rank(right, 0)
		var kf int
		if f == 0 {
			kf = r0 - l0
		} else {
			kf = right - left - (r0 - l0)
		}
		if k < kf {
			if f == 0 {
				left = l0
				right = r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else {
			k -= kf
			res |= 1 << d
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left = l0
				right = r0
			}
		}
	}
	return res
}

// [left, right) 中小于等于 x 的数中最大的数
//  如果不存在则返回-INF
func (w *WaveletMatrixXor) Floor(start, end, value, xor int) int {
	less := w.CountPrefix(start, end, value, xor)
	if less == 0 {
		return -INF
	}
	return w.Kth(start, end, less-1, xor)
}

// [left, right) 中大于等于 x 的数中最小的数
//  如果不存在则返回INF
func (w *WaveletMatrixXor) Ceil(start, end, value, xor int) int {
	less := w.CountPrefix(start, end, value, xor)
	if less == end-start {
		return INF
	}
	return w.Kth(start, end, less, xor)
}

type BitVector struct {
	data [][2]int
}

func NewBitVector(n int) *BitVector {
	return &BitVector{data: make([][2]int, (n+63)>>5)}
}

func (bv *BitVector) Set(i int) {
	bv.data[i>>5][0] |= 1 << (i & 31)
}

func (bv *BitVector) Build() {
	for i := 0; i < len(bv.data)-1; i++ {
		bv.data[i+1][1] = bv.data[i][1] + bits.OnesCount(uint(bv.data[i][0]))
	}
}

// [0, k) 内1的個数
func (bv *BitVector) Rank(k int, f int) int {
	a, b := bv.data[k>>5][0], bv.data[k>>5][1]
	ret := b + bits.OnesCount(uint(a&((1<<(k&31))-1)))
	if f == 1 {
		return ret
	}
	return k - ret
}
