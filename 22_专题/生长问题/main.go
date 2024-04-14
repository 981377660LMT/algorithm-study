// 生长问题

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	P5579()
	// CF453E()
}

// P5579 [PA2015] Siano
// https://www.luogu.com.cn/problem/P5579
// 每一亩土地上都种植了一种独一无二的草，其中，第 i 亩土地的草每天会长高grows[i]厘米。
// Byteasar 一共会进行 m 次收割，其中第i 次收割在第day[i]天，并把所有高度大于等于cut[i]的部分全部割去。
// 每次收割得到的草的高度总和是多少。
//
// https://blog.csdn.net/Maxwei_wzj/article/details/82889218
// !1.将草按照grow[i]排序，每次被割的草一定是一段后缀，可以二分定位.
// !2.第day天草i的高度为：
//
//		lastHeight[i]+(day-lastDay[i])*grows[i]
//
//	 其中lastHeight[i]表示上次被割的高度，lastDay[i]表示上次被割的时间，变形为
//
//		day*grows[i] + (lastHeight[i]-lastDay[i]*grows[i])，
//
// !lastHeight[i]和lastDay[i]可以在线段树上维护.
// !3.查询时，答案为 区间草的高度之和 - 区间长度*cut[i]
func P5579() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	grows := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &grows[i])
	}
	sort.Ints(grows)

	query := func(day int, cut int) int {
		return 0
	}

	// harvest query
	for i := 0; i < m; i++ {
	}
}

// Little Pony and Lord Tirek (珂朵莉树分割区间 + waveletMatrix 求和)
// https://www.luogu.com.cn/problem/CF453E
//
// !给定一个数组，每秒nums[i]会变为min(nums[i]+grows[i],maxHeights[i])
// 每次操作询问区间和，并将区间清零.
//
// 思路：
// 1.珂朵莉树维护收割时间，将区间分割为多段"上一次收割时间相同"的区间.
// 2.另一个数据结构维护区间信息.
// 3.收割时遍历珂朵莉树每一段区间，根据生长时间、区间和信息计算收割答案.
//
//	这里将植物分为两类：已经长到最大高度的植物和未长到最大高度的植物.
//	已经长到最大高度的植物的和为maxHeights[i]之和，未长到最大高度的植物的和为grows[i]*天数之和.
func CF453E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)

	nums := make([]int, n)
	maxHeights := make([]int, n)
	grows := make([]int, n)
	for i := 0; i < n; i++ {
		var v, h, g int
		fmt.Fscan(in, &v, &h, &g)
		nums[i] = v
		maxHeights[i] = h
		grows[i] = g
	}

	time := make([]int, n) // 生长到最大高度的天数
	for i := 0; i < n; i++ {
		if grows[i] == 0 {
			time[i] = 1 << 30
		} else {
			time[i] = (maxHeights[i] + grows[i] - 1) / grows[i] // ceil(maxHeights[i], grows[i])
		}
	}

	last := NewODT(n, -1) // 上一次被收割的时间，-1表示从未被收割

	sumData := make([][2]int, n) // grow[i]之和、maxHeights[i]之和
	for i := 0; i < n; i++ {
		sumData[i] = [2]int{grows[i], maxHeights[i]}
	}
	wm := NewWaveletMatrixWithSum(time, sumData, -1, true)

	query := func(day int, start, end int) int {
		res := 0
		f := func(l, r int, preDay int) {
			if preDay == -1 {
				for i := l; i < r; i++ {
					now := min(nums[i]+grows[i]*day, maxHeights[i])
					res += now
				}
				return
			}

			elapse := day - preDay
			count1 := wm.CountPrefix(int32(l), int32(r), elapse+1, 0) // 多少个植物已经长到最大高度
			sum1 := wm.SumPrefix(int32(l), int32(r), count1, 0)
			res += sum1[1] // 已经长到最大高度的植物的和(maxHeights[i]之和)
			sum2 := wm.SumRange(int32(l), int32(r), count1, int32(r-l), 0)
			res += sum2[0] * elapse // 未长到最大高度的植物的和(grows[i]*天数之和)
		}

		last.EnumerateRange(start, end, f, true)
		last.Set(start, end, day)
		return res
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var day, start, end int
		fmt.Fscan(in, &day, &start, &end)
		start--
		fmt.Println(query(day, start, end))
	}
}

const INF WmValue = 1e9 + 10

type WmValue = int
type WmSum = [2]int

func (*WaveletMatrixWithSum) e() WmSum            { return WmSum{0, 0} }
func (*WaveletMatrixWithSum) op(a, b WmSum) WmSum { return WmSum{a[0] + b[0], a[1] + b[1]} }
func (*WaveletMatrixWithSum) inv(a WmSum) WmSum   { return WmSum{-a[0], -a[1]} }

type WaveletMatrixWithSum struct {
	n, log   int32
	mid      []int32
	bv       []*BitVector
	key      []WmValue
	setLog   bool
	presum   [][]WmSum
	compress bool
}

// nums: 数组元素.
// sumData: 和数据，nil表示不需要和数据.
// log: 如果需要支持异或查询则需要传入log，-1表示默认.
// compress: 是否对nums进行离散化(值域较大(1e9)时可以离散化加速).
func NewWaveletMatrixWithSum(nums []WmValue, sumData []WmSum, log int32, compress bool) *WaveletMatrixWithSum {
	wm := &WaveletMatrixWithSum{}
	wm.build(nums, sumData, log, compress)
	return wm
}

func (wm *WaveletMatrixWithSum) build(nums []WmValue, sumData []WmSum, log int32, compress bool) {
	numsCopy := append(nums[:0:0], nums...)
	sumDataCopy := append(sumData[:0:0], sumData...)

	wm.n = int32(len(numsCopy))
	wm.log = log
	wm.compress = compress
	wm.setLog = log != -1
	if wm.n == 0 {
		wm.log = 0
		return
	}
	makeSum := len(sumData) > 0
	if compress {
		if wm.setLog {
			panic("compress and log should not be set at the same time")
		}
		wm.key = make([]WmValue, 0, wm.n)
		order := wm._argSort(numsCopy)
		for _, i := range order {
			if len(wm.key) == 0 || wm.key[len(wm.key)-1] != numsCopy[i] {
				wm.key = append(wm.key, numsCopy[i])
			}
			numsCopy[i] = WmValue(len(wm.key) - 1)
		}
	}
	if wm.log == -1 {
		tmp := wm._maxs(numsCopy)
		if tmp < 1 {
			tmp = 1
		}
		wm.log = int32(bits.Len(uint(tmp)))
	}
	wm.mid = make([]int32, wm.log)
	wm.bv = make([]*BitVector, wm.log)
	for i := range wm.bv {
		wm.bv[i] = NewBitVector(wm.n)
	}
	if makeSum {
		wm.presum = make([][]WmSum, 1+wm.log)
		for i := range wm.presum {
			sums := make([]WmSum, wm.n+1)
			for j := range sums {
				sums[j] = wm.e()
			}
			wm.presum[i] = sums
		}
	}
	if len(sumDataCopy) == 0 {
		sumDataCopy = make([]WmSum, len(numsCopy))
	}

	A, S := numsCopy, sumDataCopy
	A0, A1 := make([]WmValue, wm.n), make([]WmValue, wm.n)
	S0, S1 := make([]WmSum, wm.n), make([]WmSum, wm.n)
	for d := wm.log - 1; d >= -1; d-- {
		p0, p1 := int32(0), int32(0)
		if makeSum {
			tmp := wm.presum[d+1]
			for i := int32(0); i < wm.n; i++ {
				tmp[i+1] = wm.op(tmp[i], S[i])
			}
		}
		if d == -1 {
			break
		}
		for i := int32(0); i < wm.n; i++ {
			f := (A[i] >> d & 1) == 1
			if !f {
				if makeSum {
					S0[p0] = S[i]
				}
				A0[p0] = A[i]
				p0++
			} else {
				if makeSum {
					S1[p1] = S[i]
				}
				wm.bv[d].Set(i)
				A1[p1] = A[i]
				p1++
			}
		}
		wm.mid[d] = p0
		wm.bv[d].Build()
		A, A0 = A0, A
		S, S0 = S0, S
		for i := int32(0); i < p1; i++ {
			A[p0+i] = A1[i]
			S[p0+i] = S1[i]
		}
	}
}

// 区间 [start, end) 中范围在 [a, b) 中的元素的个数.
func (wm *WaveletMatrixWithSum) CountRange(start, end int32, a, b WmValue, xorVal WmValue) int32 {
	return wm.CountPrefix(start, end, b, xorVal) - wm.CountPrefix(start, end, a, xorVal)
}

func (wm *WaveletMatrixWithSum) SumRange(start, end, k1, k2 int32, xorVal WmValue) WmSum {
	if k1 >= k2 {
		return wm.e()
	}
	add := wm.SumPrefix(start, end, k2, xorVal)
	sub := wm.SumPrefix(start, end, k1, xorVal)
	return wm.op(add, wm.inv(sub))
}

// 返回区间 [start, end) 中 范围在 [0, x) 中的元素的个数.
func (wm *WaveletMatrixWithSum) CountPrefix(start, end int32, x WmValue, xor WmValue) int32 {
	if xor != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	if x == 0 {
		return 0
	}
	if x >= 1<<wm.log {
		return end - start
	}
	count := int32(0)
	for d := wm.log - 1; d >= 0; d-- {
		add := (x>>d)&1 == 1
		f := (xor>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		kf := int32(0)
		if f {
			kf = (end - start - r0 + l0)
		} else {
			kf = (r0 - l0)
		}
		if add {
			count += kf
			if f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		}
	}
	return count
}

// 返回区间 [start, end) 中前k个(k>=0)元素的和.
func (wm *WaveletMatrixWithSum) SumPrefix(start, end, k int32, xor WmValue) WmSum {
	_, sum := wm.KthValueAndSum(start, end, k, xor)
	return sum
}

// [start, end)区间内第k(k>=0)小的元素.
func (wm *WaveletMatrixWithSum) Kth(start, end, k int32, xorVal WmValue) WmValue {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set")
		}
	}
	count := int32(0)
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		var c int32
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = r0 - l0
		}
		if count+c > k {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			count += c
			res |= 1 << d
			if !f {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			} else {
				start, end = l0, r0
			}
		}
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res
}

// 返回区间 [start, end) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果).
// 如果k >= end-start, 返回 (-1, 区间 op 的结果).
func (wm *WaveletMatrixWithSum) KthValueAndSum(start, end, k int32, xorVal WmValue) (WmValue, WmSum) {
	if start >= end {
		return -1, wm.e()
	}
	if k >= end-start {
		return -1, wm.SumAll(start, end)
	}
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if len(wm.presum) == 0 {
		panic("sumData is not provided")
	}
	count := int32(0)
	sum := wm.e()
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		c := int32(0)
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = r0 - l0
		}
		if count+c > k {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			var s WmSum
			if f {
				s = wm._get(d, start+wm.mid[d]-l0, end+wm.mid[d]-r0)
			} else {
				s = wm._get(d, l0, r0)
			}
			count += c
			sum = wm.op(sum, s)
			res |= 1 << d
			if !f {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			} else {
				start, end = l0, r0
			}
		}
	}
	sum = wm.op(sum, wm._get(0, start, start+k-count))
	if wm.compress {
		res = wm.key[res]
	}
	return res, sum
}

// upper: 向上取中位数还是向下取中位数.
func (wm *WaveletMatrixWithSum) Median(start, end int32, upper bool, xorVal WmValue) WmValue {
	n := end - start
	var k int32
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.Kth(start, end, k, xorVal)
}

func (wm *WaveletMatrixWithSum) SumAll(start, end int32) WmSum {
	return wm._get(wm.log, start, end)
}

// 使得predicate(count, sum)为true的最大的(count, sum).
func (wm *WaveletMatrixWithSum) MaxRight(predicate func(int32, WmSum) bool, start, end int32, xorVal WmValue) (int32, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start == end {
		return end - start, wm.e()
	}
	if s := wm._get(wm.log, start, end); predicate(end-start, s) {
		return end - start, s
	}
	count := int32(0)
	sum := wm.e()
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		c := int32(0)
		if f {
			c = (end - start) - (r0 - l0)
		} else {
			c = (r0 - l0)
		}
		var s WmSum
		if f {
			s = wm._get(d, start+wm.mid[d]-l0, end+wm.mid[d]-r0)
		} else {
			s = wm._get(d, l0, r0)
		}
		if tmp := wm.op(sum, s); predicate(count+c, tmp) {
			count += c
			sum = tmp
			if f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		} else {
			if !f {
				start, end = l0, r0
			} else {
				start += wm.mid[d] - l0
				end += wm.mid[d] - r0
			}
		}
	}
	k := wm._binarySearch(func(k int32) bool {
		return predicate(count+k, wm.op(sum, wm._get(0, start, start+k)))
	}, 0, end-start)
	count += k
	sum = wm.op(sum, wm._get(0, start, start+k))
	return count, sum
}

// [start, end) 中小于等于 x 的数中最大的数
//
//	如果不存在则返回-INF
func (wm *WaveletMatrixWithSum) Floor(start, end int32, x WmValue, xor WmValue) WmValue {
	less := wm.CountPrefix(start, end, x, xor)
	if less == 0 {
		return -INF
	}
	res := wm.Kth(start, end, less-1, xor)
	return res
}

// [start, end) 中大于等于 x 的数中最小的数
//
//	如果不存在则返回INF
func (wm *WaveletMatrixWithSum) Ceil(start, end int32, x WmValue, xor WmValue) int {
	less := wm.CountPrefix(start, end, x, xor)
	if less == end-start {
		return INF
	}
	res := wm.Kth(start, end, less, xor)
	return res
}

func (wm *WaveletMatrixWithSum) CountSegments(segments [][2]int32, a, b WmValue, xorVal WmValue) int32 {
	res := int32(0)
	for _, seg := range segments {
		res += wm.CountRange(seg[0], seg[1], a, b, xorVal)
	}
	return res
}

func (wm *WaveletMatrixWithSum) KthSegments(segments [][2]int32, k int32, xorVal WmValue) WmValue {
	totalLen := int32(0)
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	count := int32(0)
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		c := int32(0)
		for _, seg := range segments {
			l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
			if f {
				c += (seg[1] - seg[0]) - (r0 - l0)
			} else {
				c += r0 - l0
			}
		}
		if count+c > k {
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if !f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		}
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res
}

func (wm *WaveletMatrixWithSum) KthValueAndSumSegments(segments [][2]int32, k int32, xorVal WmValue) (WmValue, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	totalLen := int32(0)
	for _, seg := range segments {
		totalLen += seg[1] - seg[0]
	}
	if k >= totalLen {
		return -1, wm.SumAllSegments(segments)
	}
	count := int32(0)
	sum := wm.e()
	res := WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		f := (xorVal>>d)&1 == 1
		c := int32(0)
		for _, seg := range segments {
			l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
			if f {
				c += (seg[1] - seg[0]) - (r0 - l0)
			} else {
				c += r0 - l0
			}
		}
		if count+c > k {
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				if !f {
					seg[0], seg[1] = l0, r0
				} else {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				}
			}
		} else {
			count += c
			res |= 1 << d
			for _, seg := range segments {
				l0, r0 := wm.bv[d].Rank(seg[0], false), wm.bv[d].Rank(seg[1], false)
				var s WmSum
				if f {
					s = wm._get(d, seg[0]+wm.mid[d]-l0, seg[1]+wm.mid[d]-r0)
				} else {
					s = wm._get(d, l0, r0)
				}
				sum = wm.op(sum, s)
				if !f {
					seg[0] += wm.mid[d] - l0
					seg[1] += wm.mid[d] - r0
				} else {
					seg[0], seg[1] = l0, r0
				}
			}
		}
	}
	for _, seg := range segments {
		t := min32(seg[1]-seg[0], k-count)
		sum = wm.op(sum, wm._get(0, seg[0], seg[0]+t))
		count += t
	}
	if wm.compress {
		res = wm.key[res]
	}
	return res, sum
}

func (wm *WaveletMatrixWithSum) MedianSegments(segments [][2]int32, upper bool, xorVal WmValue) WmValue {
	n := int32(0)
	for _, seg := range segments {
		n += seg[1] - seg[0]
	}
	var k int32
	if upper {
		k = n / 2
	} else {
		k = (n - 1) / 2
	}
	return wm.KthSegments(segments, k, xorVal)
}

func (wm *WaveletMatrixWithSum) SumAllSegments(segments [][2]int32) WmSum {
	sum := wm.e()
	for _, seg := range segments {
		sum = wm.op(sum, wm._get(wm.log, seg[0], seg[1]))
	}
	return sum
}

func (wm *WaveletMatrixWithSum) _prefixSumSegments(segments [][2]int32, k int32, xor WmValue) WmSum {
	_, sum := wm.KthValueAndSumSegments(segments, k, xor)
	return sum
}

func (wm *WaveletMatrixWithSum) _get(d, l, r int32) WmSum {
	return wm.op(wm.inv(wm.presum[d][l]), wm.presum[d][r])
}

func (wm *WaveletMatrixWithSum) _argSort(nums []WmValue) []int32 {
	order := make([]int32, len(nums))
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}

func (wm *WaveletMatrixWithSum) _maxs(nums []WmValue) WmValue {
	res := nums[0]
	for _, v := range nums {
		if v > res {
			res = v
		}
	}
	return res
}

func (wm *WaveletMatrixWithSum) _lowerBound(nums []WmValue, target WmValue) WmValue {
	left, right := int32(0), int32(len(nums)-1)
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] < target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return WmValue(left)
}

func (wm *WaveletMatrixWithSum) _binarySearch(f func(int32) bool, ok, ng int32) int32 {
	for abs32(ok-ng) > 1 {
		x := (ok + ng) >> 1
		if f(x) {
			ok = x
		} else {
			ng = x
		}
	}
	return ok
}

type BitVector struct {
	bits   []uint64
	preSum []int32
}

func NewBitVector(n int32) *BitVector {
	return &BitVector{bits: make([]uint64, n>>6+1), preSum: make([]int32, n>>6+1)}
}

func (bv *BitVector) Set(i int32) {
	bv.bits[i>>6] |= 1 << (i & 63)
}

func (bv *BitVector) Build() {
	for i := 0; i < len(bv.bits)-1; i++ {
		bv.preSum[i+1] = bv.preSum[i] + int32(bits.OnesCount64(bv.bits[i]))
	}
}

func (bv *BitVector) Rank(k int32, f bool) int32 {
	m, s := bv.bits[k>>6], bv.preSum[k>>6]
	res := s + int32(bits.OnesCount64(m&((1<<(k&63))-1)))
	if f {
		return res
	}
	return k - res
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
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

// RangeAssignRangeSumMin

type E = struct{ sum, size, min int }
type Id = int

func (*LazySegTree) e() E   { return E{min: INF} }
func (*LazySegTree) id() Id { return INF }
func (*LazySegTree) op(left, right E) E {
	return E{left.sum + right.sum, left.size + right.size, min(left.min, right.min)}
}
func (*LazySegTree) mapping(f Id, g E) E {
	if f == INF {
		return g
	}
	return E{f * g.size, g.size, f}
}
func (*LazySegTree) composition(f, g Id) Id {
	if f == INF {
		return g
	}
	return f
}

// !template
type LazySegTree struct {
	n    int
	size int
	log  int
	data []E
	lazy []Id
}

func NewLazySegTree(n int, f func(int) E) *LazySegTree {
	tree := &LazySegTree{}
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = f(i)
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

func NewLazySegTreeFrom(leaves []E) *LazySegTree {
	tree := &LazySegTree{}
	n := len(leaves)
	tree.n = n
	tree.log = int(bits.Len(uint(n - 1)))
	tree.size = 1 << tree.log
	tree.data = make([]E, tree.size<<1)
	tree.lazy = make([]Id, tree.size)
	for i := range tree.data {
		tree.data[i] = tree.e()
	}
	for i := range tree.lazy {
		tree.lazy[i] = tree.id()
	}
	for i := 0; i < n; i++ {
		tree.data[tree.size+i] = leaves[i]
	}
	for i := tree.size - 1; i >= 1; i-- {
		tree.pushUp(i)
	}
	return tree
}

// 查询切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Query(left, right int) E {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return tree.e()
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	sml, smr := tree.e(), tree.e()
	for left < right {
		if left&1 != 0 {
			sml = tree.op(sml, tree.data[left])
			left++
		}
		if right&1 != 0 {
			right--
			smr = tree.op(tree.data[right], smr)
		}
		left >>= 1
		right >>= 1
	}
	return tree.op(sml, smr)
}
func (tree *LazySegTree) QueryAll() E {
	return tree.data[1]
}
func (tree *LazySegTree) GetAll() []E {
	res := make([]E, tree.n)
	copy(res, tree.data[tree.size:tree.size+tree.n])
	return res
}

// 更新切片[left:right]的值
//
//	0<=left<=right<=len(tree.data)
func (tree *LazySegTree) Update(left, right int, f Id) {
	if left < 0 {
		left = 0
	}
	if right > tree.n {
		right = tree.n
	}
	if left >= right {
		return
	}
	left += tree.size
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		if ((left >> i) << i) != left {
			tree.pushDown(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushDown((right - 1) >> i)
		}
	}
	l2, r2 := left, right
	for left < right {
		if left&1 != 0 {
			tree.propagate(left, f)
			left++
		}
		if right&1 != 0 {
			right--
			tree.propagate(right, f)
		}
		left >>= 1
		right >>= 1
	}
	left = l2
	right = r2
	for i := 1; i <= tree.log; i++ {
		if ((left >> i) << i) != left {
			tree.pushUp(left >> i)
		}
		if ((right >> i) << i) != right {
			tree.pushUp((right - 1) >> i)
		}
	}
}

// 二分查询最小的 left 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MinLeft(right int, predicate func(data E) bool) int {
	if right == 0 {
		return 0
	}
	right += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown((right - 1) >> i)
	}
	res := tree.e()
	for {
		right--
		for right > 1 && right&1 != 0 {
			right >>= 1
		}
		if !predicate(tree.op(tree.data[right], res)) {
			for right < tree.size {
				tree.pushDown(right)
				right = right<<1 | 1
				if tmp := tree.op(tree.data[right], res); predicate(tmp) {
					res = tmp
					right--
				}
			}
			return right + 1 - tree.size
		}
		res = tree.op(tree.data[right], res)
		if (right & -right) == right {
			break
		}
	}
	return 0
}

// 二分查询最大的 right 使得切片 [left:right] 内的值满足 predicate
func (tree *LazySegTree) MaxRight(left int, predicate func(data E) bool) int {
	if left == tree.n {
		return tree.n
	}
	left += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(left >> i)
	}
	res := tree.e()
	for {
		for left&1 == 0 {
			left >>= 1
		}
		if !predicate(tree.op(res, tree.data[left])) {
			for left < tree.size {
				tree.pushDown(left)
				left <<= 1
				if tmp := tree.op(res, tree.data[left]); predicate(tmp) {
					res = tmp
					left++
				}
			}
			return left - tree.size
		}
		res = tree.op(res, tree.data[left])
		left++
		if (left & -left) == left {
			break
		}
	}
	return tree.n
}

// 单点查询(不需要 pushUp/op 操作时使用)
func (tree *LazySegTree) Get(index int) E {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	return tree.data[index]
}

// 单点赋值
func (tree *LazySegTree) Set(index int, e E) {
	index += tree.size
	for i := tree.log; i >= 1; i-- {
		tree.pushDown(index >> i)
	}
	tree.data[index] = e
	for i := 1; i <= tree.log; i++ {
		tree.pushUp(index >> i)
	}
}

func (tree *LazySegTree) pushUp(root int) {
	tree.data[root] = tree.op(tree.data[root<<1], tree.data[root<<1|1])
}
func (tree *LazySegTree) pushDown(root int) {
	if tree.lazy[root] != tree.id() {
		tree.propagate(root<<1, tree.lazy[root])
		tree.propagate(root<<1|1, tree.lazy[root])
		tree.lazy[root] = tree.id()
	}
}
func (tree *LazySegTree) propagate(root int, f Id) {
	tree.data[root] = tree.mapping(f, tree.data[root])
	// !叶子结点不需要更新lazy
	if root < tree.size {
		tree.lazy[root] = tree.composition(f, tree.lazy[root])
	}
}

func (tree *LazySegTree) String() string {
	var sb []string
	sb = append(sb, "[")
	for i := 0; i < tree.n; i++ {
		if i != 0 {
			sb = append(sb, ", ")
		}
		sb = append(sb, fmt.Sprintf("%v", tree.Get(i)))
	}
	sb = append(sb, "]")
	return strings.Join(sb, "")
}
