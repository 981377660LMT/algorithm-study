// 维护区间贡献的 Wavelet Matrix
// !注意查询区间贡献时, 异或无效
// Count/CountPrefix/Kth/MaxRightCount

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	// https://yukicoder.me/problems/no/924
	// n,q<=2e5
	// -1e9 <= nums[i] <= 1e9
	// 给定n个查询[l,r]
	// !求区间[l,r]中位数到区间[l,r]中每个数的距离之和
	// !也就求函数 f(x)= ∑|nums[i]-x| (l<=i<=right) 的最小值
	// !区间中位数/区间最短距离和/区间曼哈顿距离
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	OFFSET := int(1e9 + 10)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
		nums[i] += OFFSET
	}
	preSum := make([]int, n+1)
	for i := range nums {
		preSum[i+1] = preSum[i] + nums[i]
	}

	wm := NewWaveletMatrixSum(nums, 32+2)
	for i := 0; i < q; i++ {
		var left, right int // 1<=l<=r<=n
		fmt.Fscan(in, &left, &right)
		left--

		n := right - left
		lowerCount := n / 2
		ceilCount := n - lowerCount
		mid, lowerSum := wm.Kth(left, right, lowerCount, 0)
		_, allSum := wm.Kth(left, right, n, 0)
		ceilSum := allSum - lowerSum

		res := 0
		res += mid*lowerCount - lowerSum
		res += ceilSum - mid*ceilCount
		fmt.Fprintln(out, res)
	}
}

const INF int = 1e18

type E = int

func (*WaveletMatrixSum) e() E        { return 0 }
func (*WaveletMatrixSum) op(a, b E) E { return a + b }
func (*WaveletMatrixSum) inv(a E) E   { return -a }

type WaveletMatrixSum struct {
	n, log int
	mid    []int
	bv     []*BitVector
	preSum [][]int
	unit   E
}

// log:如果要支持异或,则需要按照异或的值来决定值域
//  设为-1时表示不使用异或
func NewWaveletMatrixSum(nums []int, log int) *WaveletMatrixSum {
	numsCopy := make([]int, len(nums))
	max_ := 0
	for i, v := range nums {
		numsCopy[i] = v
		if v > max_ {
			max_ = v
		}
	}
	if log == -1 {
		log = bits.Len(uint(max_))
	}
	res := &WaveletMatrixSum{}
	res.unit = res.e()
	n := len(numsCopy)
	mid := make([]int, log)
	bv := make([]*BitVector, log)
	for i := 0; i < log; i++ {
		bv[i] = NewBitVector(n)
	}
	preSum := make([][]int, log+1)
	for i := range preSum {
		preSum[i] = make([]int, n+1)
		for j := range preSum[i] {
			preSum[i][j] = res.unit
		}
	}

	a0, a1 := make([]int, n), make([]int, n)
	for d := log - 1; d >= -1; d-- {
		p0, p1 := 0, 0
		for i := 0; i < n; i++ {
			preSum[d+1][i+1] = res.op(preSum[d+1][i], numsCopy[i])
		}
		if d == -1 {
			break
		}
		for i := 0; i < n; i++ {
			f := (numsCopy[i] >> d) & 1
			if f == 0 {
				a0[p0] = numsCopy[i]
				p0++
			} else {
				bv[d].Set(i)
				a1[p1] = numsCopy[i]
				p1++
			}
		}
		mid[d] = p0
		bv[d].Build()
		numsCopy, a0 = a0, numsCopy
		for i := 0; i < p1; i++ {
			numsCopy[p0+i] = a1[i]
		}
	}

	res.n, res.log = n, log
	res.mid, res.bv, res.preSum = mid, bv, preSum
	return res
}

// 返回区间 [left, right) 中 范围在 [a, b) 中的 (元素的个数, op 的结果)
func (wm *WaveletMatrixSum) CountRange(left, right, a, b, xor int) (int, E) {
	c1, s1 := wm.CountPrefix(left, right, a, xor)
	c2, s2 := wm.CountPrefix(left, right, b, xor)
	return c2 - c1, wm.op(wm.inv(s1), s2)
}

// 返回区间 [left, right) 中 范围在 [0, x) 中的 (元素的个数, op 的结果)
func (wm *WaveletMatrixSum) CountPrefix(left, right, x, xor int) (int, E) {
	if x >= 1<<wm.log {
		return right - left, wm.get(wm.log, left, right)
	}
	count := 0
	sum := wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		add := (x >> d) & 1
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if add == 1 {
			count += kf
			if f == 1 {
				sum = wm.op(sum, wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0))
				left, right = l0, r0
			} else {
				sum = wm.op(sum, wm.get(d, l0, r0))
				left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
			}
		} else if f == 0 {
			left, right = l0, r0
		} else {
			left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
		}
	}

	return count, sum
}

// 返回区间 [left, right) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
//  如果k < 0, 返回 (-1, 0); 如果k >= right-left, 返回 (-1, 区间 op 的结果)
func (wm *WaveletMatrixSum) Kth(left, right, k, xor int) (int, E) {
	if k < 0 {
		return -1, 0
	}
	if right-left <= k {
		return -1, wm.get(wm.log, left, right)
	}
	res, sum := 0, wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		kf := f*(right-left-r0+l0) + (f^1)*(r0-l0)
		if k < kf {
			if f == 0 {
				left, right = l0, r0
			} else {
				left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
			}
		} else {
			k -= kf
			res |= 1 << d
			if f == 1 {
				sum = wm.op(sum, wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0))
				left, right = l0, r0
			} else {
				sum = wm.op(sum, wm.get(d, l0, r0))
				left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
			}
		}
	}

	if k != 0 {
		sum = wm.op(sum, wm.get(0, left, left+k))
	}
	return res, sum
}

// 返回使得 check(prefixSum) 为 true 的最大值 val.
//  !(即区间内小于 val 的数的和 prefixSum 满足 check函数, 找到这样的最大的 val)
//  如果整个区间都满足, 返回 INF.
//  eg: val = 5 => 即区间内值域在 [0,5) 中的数的和满足 check 函数.
func (wm *WaveletMatrixSum) MaxRightValue(left, right, xor int, check func(preSum E) bool) E {
	if check(wm.get(wm.log, left, right)) {
		return INF
	}
	res := 0
	sum := wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		var loSum E
		if f == 0 {
			loSum = wm.get(d, l0, r0)
		} else {
			loSum = wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0)
		}
		if check(wm.op(sum, loSum)) {
			sum = wm.op(sum, loSum)
			res |= 1 << d
			if f == 1 {
				left, right = l0, r0
			} else {
				left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
			}
		} else if f == 0 {
			left, right = l0, r0
		} else {
			left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
		}
	}

	return res
}

// 返回使得 check(prefixSum) 为 true 的区间前缀个数的最大值.
//  eg: count = 4 => 即区间内的数排序后, 前4个数的和满足 check 函数.
func (wm *WaveletMatrixSum) MaxRightCount(left, right, xor int, check func(preSum E) bool) int {
	if check(wm.get(wm.log, left, right)) {
		return right - left
	}

	res := 0
	sum := wm.unit
	for d := wm.log - 1; d >= 0; d-- {
		f := (xor >> d) & 1
		l0, r0 := wm.bv[d].Rank(left, 0), wm.bv[d].Rank(right, 0)
		var kf int
		var loSum E
		if f == 0 {
			kf = r0 - l0
			loSum = wm.get(d, l0, r0)
		} else {
			kf = (right - left) - (r0 - l0)
			loSum = wm.get(d, left+wm.mid[d]-l0, right+wm.mid[d]-r0)
		}

		if check(wm.op(sum, loSum)) {
			sum = wm.op(sum, loSum)
			res += kf
			if f == 1 {
				left, right = l0, r0
			} else {
				left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
			}
		} else if f == 0 {
			left, right = l0, r0
		} else {
			left, right = left+wm.mid[d]-l0, right+wm.mid[d]-r0
		}
	}

	res += wm.binarySearch(func(k int) bool {
		return check(wm.op(sum, wm.get(0, left, left+k)))
	}, 0, right-left)

	return res
}

// [left, right) 中小于等于 x 的数中最大的数
//  如果不存在则返回-INF
func (w *WaveletMatrixSum) Floor(start, end, value, xor int) int {
	less, _ := w.CountPrefix(start, end, value, xor)
	if less == 0 {
		return -INF
	}
	res, _ := w.Kth(start, end, less-1, xor)
	return res
}

// [left, right) 中大于等于 x 的数中最小的数
//  如果不存在则返回INF
func (w *WaveletMatrixSum) Ceiling(start, end, value, xor int) int {
	less, _ := w.CountPrefix(start, end, value, xor)
	if less == end-start {
		return INF
	}
	res, _ := w.Kth(start, end, less, xor)
	return res
}

func (wm *WaveletMatrixSum) binarySearch(f func(E) bool, ok, ng int) int {
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

func (wm *WaveletMatrixSum) get(d, l, r int) E {
	return wm.op(wm.inv(wm.preSum[d][l]), wm.preSum[d][r])
}

func abs(a int) int {
	if a < 0 {
		return -a
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

// [0, k) 内の 1 の個数
func (bv *BitVector) Rank(k int, f int) int {
	a, b := bv.data[k>>5][0], bv.data[k>>5][1]
	ret := b + bits.OnesCount(uint(a&((1<<(k&31))-1)))
	if f == 1 {
		return ret
	}
	return k - ret
}
