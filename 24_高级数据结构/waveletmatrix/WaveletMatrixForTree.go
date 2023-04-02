// CountRange: 返回区间 [left, right) 中 范围在 [a, b) 中的 元素的个数.
// Kth: 返回区间 [left, right) 中的 第k小的元素
// KthValueAndSum: 返回区间 [left, right) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
// Sum: 返回区间 [left, right) 中第[k1, k2)个元素的 op 的结果
// SumAll: 返回区间 [left, right) 的 op 的结果
// Median: 返回区间 [left, right) 的中位数

// CountRangeSegments: 返回所有区间中 范围在 [a, b) 中的 元素的个数.
// KthSegments: 返回所有区间中的 第k小的元素
// KthValueAndSumSegments: 返回所有区间中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
// SumSegments: 返回所有区间中所有数的 op 的结果
// SumAllSegments: 返回所有区间的 op 的结果
// MedianSegments: 返回所有区间的中位数

// MaxRight: 返回使得 check(count,prefixSum) 为 true 的最大 (count, prefixSum) 对.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func demo() {
	// https://yukicoder.me/problems/no/1332
	// 区间前驱后继
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	wm := NewWaveletMatrixSegments(nums, -1, nil)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var l, r, x int
		fmt.Fscan(in, &l, &r, &x)
		l--
		res := INF
		n := wm.CountRange(l, r, 0, x, 0)
		if n > 0 {
			res = min(res, abs(x-wm.Kth(l, r, n-1, 0)))
		}
		if n < r-l {
			res = min(res, abs(x-wm.Kth(l, r, n, 0)))
		}
		fmt.Fprintln(out, res)
	}

}

func main() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wm := NewWaveletMatrixSegments(nums, -1, nums)

	fmt.Println(wm.CountRange(0, 10, 3, 7, 0))
	fmt.Println(wm.Kth(0, 10, 9, 0))
	fmt.Println(wm.KthValueAndSum(0, 10, 10, 0))
	fmt.Println(wm.Sum(0, 10, 1, 3, 0))
	fmt.Println(wm.SumAll(0, 10))
	fmt.Println(wm.Median(false, 0, 5, 0))

	fmt.Println(wm.CountRangeSegments([][2]int{{0, 1}, {5, 10}}, 3, 90, 0))
	fmt.Println(wm.KthSegments([][2]int{{0, 1}, {5, 10}}, 3, 0))
	fmt.Println(wm.KthValueAndSumSegments([][2]int{{0, 1}, {5, 10}}, 30, 0))
	fmt.Println(wm.SumSegments([][2]int{{0, 1}, {5, 10}}, 1, 3, 0))
	fmt.Println(wm.SumAllSegments([][2]int{{0, 1}, {5, 10}}))
	fmt.Println(wm.MedianSegments(true, [][2]int{{0, 5}}, 0))

	count, sum := wm.MaxRight(0, 10, 0, func(count int, sum int) bool {
		return count <= 3
	})
	fmt.Println(count, sum)
}

const INF int = 1e18

type E = int

func (*WaveletMatrixSegments) e() E        { return 0 }
func (*WaveletMatrixSegments) op(a, b E) E { return a + b }
func (*WaveletMatrixSegments) inv(a E) E   { return -a }

type WaveletMatrixSegments struct {
	n, log int
	mid    []int
	bv     []*BitVector
	preSum [][]int
	unit   E
}

//
// log:如果要支持异或,则需要按照异或的值来决定值域
//  设为-1时表示不使用异或
// sumData:如果要支持区间和,则需要传入数组
//  设为nil时表示不使用区间和
func NewWaveletMatrixSegments(nums []E, log int, sumData []E) *WaveletMatrixSegments {
	res := &WaveletMatrixSegments{}
	res.build(nums, log, sumData)
	return res
}

// 返回区间 [left, right) 中 范围在 [a, b) 中的 元素的个数.
func (wm *WaveletMatrixSegments) CountRange(left, right, a, b, xor int) int {
	return wm.prefixCount(left, right, b, xor) - wm.prefixCount(left, right, a, xor)
}

func (wm *WaveletMatrixSegments) CountRangeSegments(segments [][2]int, a, b, xor int) int {
	res := 0
	for _, seg := range segments {
		res += wm.CountRange(seg[0], seg[1], a, b, xor)
	}
	return res
}

// 返回区间 [left, right) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
//  如果k < 0, 返回 (-1, 0); 如果k >= right-left, 返回 (-1, 区间 op 的结果)
func (wm *WaveletMatrixSegments) KthValueAndSum(left, right, k, xor int) (int, E) {
	if k < 0 {
		return -1, 0
	}
	if right-left <= k {
		return -1, wm.get(wm.log, left, right)
	}
	res, sum := 0, wm.unit
	count := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if count+kf > k {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else {
			var s E
			if f == 0 {
				s = wm.get(d, l0, r0)
			} else {
				s = wm.get(d, wm.mid[d]-l0+left, wm.mid[d]-r0+right)
			}
			count += kf
			res |= 1 << d
			sum = wm.op(sum, s)
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		}
	}
	sum = wm.op(sum, wm.get(0, left, left+k-count))
	return res, sum
}

//  如果k < 0, 返回 (-1, 0); 如果k >= segments总长, 返回 (-1, 区间 op 的结果)
func (wm *WaveletMatrixSegments) KthValueAndSumSegments(segments [][2]int, k, xor int) (int, E) {
	if k < 0 {
		return -1, 0
	}
	totalLen := 0
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1, wm.SumAllSegments(segments)
	}
	count := 0
	sum := wm.unit
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		c := 0
		for _, seg := range segments {
			L, R := seg[0], seg[1]
			l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
			c += f*(R-L-r0+l0) + (f^1)*(r0-l0)
		}
		if count+c > k {
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				var s E
				if f == 0 {
					s = wm.get(d, l0, r0)
				} else {
					s = wm.get(d, wm.mid[d]-l0+L, wm.mid[d]-r0+R)
				}
				sum = wm.op(sum, s)
				if f == 0 {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}

	for _, seg := range segments {
		L, R := seg[0], seg[1]
		t := min(R-L, k-count)
		sum = wm.op(sum, wm.get(0, L, L+t))
		count += t
	}
	return res, sum
}

// 如果不存在,返回-1.
func (wm *WaveletMatrixSegments) Kth(left, right, k, xor int) E {
	if k < 0 || k >= right-left {
		return -1
	}
	count := 0
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		c := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if count+c > k {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else {
			count += c
			res |= 1 << d
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		}
	}
	return res
}

// 如果不存在,返回-1.
func (wm *WaveletMatrixSegments) KthSegments(segments [][2]int, k, xor int) E {
	if k < 0 {
		return -1
	}
	totalLen := 0
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1
	}
	count := 0
	res := 0
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		c := 0
		for i := range segments {
			seg := &segments[i]
			L, R := seg[0], seg[1]
			l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
			c += f*(R-L-r0+l0) + (f^1)*(r0-l0)
		}
		if count+c > k {
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for i := range segments {
				seg := &segments[i]
				L, R := seg[0], seg[1]
				l0, r0 := wm.bv[d].Rank(L, 0), wm.bv[d].Rank(R, 0)
				if f == 0 {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}

	return res
}

// 区间中位数.
//  upper: true表示上中位数, false表示下中位数.
func (wm *WaveletMatrixSegments) Median(upper bool, left, right, xor int) E {
	n := right - left
	var k int
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.Kth(left, right, k, xor)
}

func (wm *WaveletMatrixSegments) MedianSegments(upper bool, segments [][2]int, xor int) E {
	n := 0
	for _, seg := range segments {
		n += seg[1] - seg[0]
	}
	var k int
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.KthSegments(segments, k, xor)
}

func (wm *WaveletMatrixSegments) Sum(left, right, k1, k2, xor int) E {
	return wm.prefixSum(left, right, k2, xor) - wm.prefixSum(left, right, k1, xor)
}

func (wm *WaveletMatrixSegments) SumSegments(segments [][2]int, k1, k2, xor int) E {
	return wm.prefixSumSegments(segments, k2, xor) - wm.prefixSumSegments(segments, k1, xor)
}

func (wm *WaveletMatrixSegments) SumAll(left, right int) E {
	return wm.get(wm.log, left, right)
}

func (wm *WaveletMatrixSegments) SumAllSegments(segments [][2]int) E {
	res := wm.unit
	for _, seg := range segments {
		res = wm.op(res, wm.get(wm.log, seg[0], seg[1]))
	}
	return res
}

// 返回使得 check(count,prefixSum) 为 true 的最大 (count, prefixSum) 对.
//  !(即区间内小于 val 的数的个数count和 和 prefixSum 满足 check函数, 找到这样的最大的 (count, prefixSum).
//  eg: val = 5 => 即区间内值域在 [0,5) 中的数的和满足 check 函数.
func (wm *WaveletMatrixSegments) MaxRight(left, right, xor int, check func(count int, preSum E) bool) (int, E) {
	if tmp := wm.get(wm.log, left, right); check(right-left, tmp) {
		return right - left, tmp
	}
	count := 0
	res := wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		c := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		var s E
		if f == 0 {
			s = wm.get(d, l0, r0)
		} else {
			s = wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0)
		}
		if tmp := wm.op(res, s); check(count+c, tmp) {
			count += c
			res = tmp
			if f == 0 {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			} else {
				left, right = l0, r0
			}
		} else {
			if f == 0 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		}
	}
	k := wm.binarySearch(func(k int) bool {
		return check(count+k, wm.op(res, wm.get(0, left, left+k)))
	}, 0, right-left)
	count += k
	res = wm.op(res, wm.get(0, left, left+k))
	return count, res
}

func (w *WaveletMatrixSegments) build(nums []E, log int, sumData []E) {
	numsCopy := make([]E, len(nums))
	max_ := 1
	for i, v := range nums {
		numsCopy[i] = v
		if v > max_ {
			max_ = v
		}
	}
	if log == -1 {
		log = bits.Len(uint(max_))
	}

	makeSum := sumData != nil
	sumData = append(sumData[:0:0], sumData...)

	w.unit = w.e()
	n := len(numsCopy)
	mid := make([]int, log)
	bv := make([]*BitVector, log)
	for i := 0; i < log; i++ {
		bv[i] = NewBitVector(n)
	}

	var preSum [][]E
	if makeSum {
		preSum = make([][]E, log+1)
		for i := range preSum {
			preSum[i] = make([]E, n+1)
			for j := range preSum[i] {
				preSum[i][j] = w.unit
			}
		}
	}

	a0, a1 := make([]E, n), make([]E, n)
	s0, s1 := make([]E, n), make([]E, n)
	for d := log - 1; d >= -1; d-- {
		p0, p1 := 0, 0
		if makeSum {
			for i := 0; i < n; i++ {
				preSum[d+1][i+1] = w.op(preSum[d+1][i], sumData[i])
			}
		}
		if d == -1 {
			break
		}
		for i := 0; i < n; i++ {
			f := (numsCopy[i] >> d) & 1
			if f == 0 {
				if makeSum {
					s0[p0] = sumData[i]
				}
				a0[p0] = numsCopy[i]
				p0++
			} else {
				if makeSum {
					s1[p1] = sumData[i]
				}
				bv[d].Set(i)
				a1[p1] = numsCopy[i]
				p1++
			}
		}
		mid[d] = p0
		bv[d].Build()
		numsCopy, a0 = a0, numsCopy
		sumData, s0 = s0, sumData
		for i := 0; i < p1; i++ {
			numsCopy[p0+i] = a1[i]
			if makeSum {
				sumData[p0+i] = s1[i]
			}
		}
	}

	w.n, w.log = n, log
	w.mid, w.bv, w.preSum = mid, bv, preSum
}

func (wm *WaveletMatrixSegments) binarySearch(f func(E) bool, ok, ng int) int {
	for abs(ok-ng) > 1 {
		x := (ok + ng) >> 1
		if f(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

// 返回区间 [left, right) 中 范围在 [0, x) 中的 元素的个数.
func (wm *WaveletMatrixSegments) prefixCount(left, right, x, xor int) int {
	if x == 0 {
		return 0
	}
	if x >= 1<<wm.log {
		return right - left
	}
	count := 0
	for d := wm.log - 1; d >= 0; d-- {
		add := (x >> d) & 1
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if add == 1 {
			count += kf
			if f == 1 {
				left, right = l0, r0
			} else {
				left += wm.mid[d] - l0
				right += wm.mid[d] - r0
			}
		} else if f == 0 {
			left, right = l0, r0
		} else {
			left += wm.mid[d] - l0
			right += wm.mid[d] - r0
		}
	}
	return count
}

func (wm *WaveletMatrixSegments) prefixSum(left, right, k, xor int) E {
	_, res := wm.KthValueAndSum(left, right, k, xor)
	return res
}

func (wm *WaveletMatrixSegments) prefixSumSegments(segments [][2]int, k, xor int) E {
	_, res := wm.KthValueAndSumSegments(segments, k, xor)
	return res
}

func (wm *WaveletMatrixSegments) get(d, l, r int) E {
	return wm.op(wm.inv(wm.preSum[d][l]), wm.preSum[d][r])
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}

func min(a, b int) int {
	if a > b {
		return b
	}
	return a
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

func (bv *BitVector) Rank(k int, f int) int {
	a, b := bv.data[k>>5][0], bv.data[k>>5][1]
	ret := b + bits.OnesCount(uint(a&((1<<(k&31))-1)))
	if f == 1 {
		return ret
	}
	return k - ret
}
