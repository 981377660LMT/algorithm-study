// 维护区间贡献的 Wavelet Matrix
// !注意查询区间贡献时, 异或无效

// CountRange(start, end, a, b, xor) - 区间 [start, end) 中值在 [a, b) 之间的数的个数和这些数的和.
// CountPrefix(start, end, x, xor) - 区间 [start, end) 中值在 [0, x) 之间的数的个数和这些数的和.
// SumRange(start, end, k1, k2, xor) - 区间 [start, end) 中第 k1 到第k2小的数的和(0-based，左闭右开).
// SumPrefix(start, end, k, xor) - 区间 [start, end) 中前k小的数的和.

// Kth(start, end, k, xor) - 区间 [start, end) 中第 k(k>=0) 小的数.
// KthValueAndSum(start, end, k, xor) - 区间 [start, end) 中第 k 小的数(0-indexed) 和前 k 小的数的和(不包括这个数).
// Median(start, end, upper, xor) - 区间 [start, end) 中的中位数.
// SumAll(start, end) - 区间 [start, end) 中所有数的和.

// Floor(start, end, x, xor) - 区间 [start, end) 中值小于等于 x 的最大值
// Ceil(start, end, x, xor) - 区间 [start, end) 中值大于等于 x 的最小值

// MaxRight(start, end, xor, check) - 返回使得 check(count,prefixSum) 为 true 的最大(count,prefixSum), 其中prefixSum为[0,val)内的数的和.

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	// demo()
	// 区间最短距离和()
	// abc281_e()

	yuki919()
	// yuki1332()
	// yuki2065()
}

func demo() {
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	wm := NewWaveletMatrixWithSum(nums, nums, -1, false)
	fmt.Println(wm.CountRange(0, 3, 3, 7, 0))
	fmt.Println(wm.KthValueAndSum(1, 5, 0, 0))
	nums = []int{1, 2, 3, 4, 5}
	wm = NewWaveletMatrixWithSum(nums, nums, -1, false)
	fmt.Println(wm.KthValueAndSum(0, 5, 2, 0))
	fmt.Println(wm.KthValueAndSum(0, 6, 5, 0))
	fmt.Println(wm.SumAll(0, 4))

	fmt.Println(wm.Floor(0, 5, 4, 0))
	fmt.Println(wm.Kth(0, 5, 1, 0))
}

func abc281_e() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	var n, m, k int32
	fmt.Fscan(in, &n, &m, &k)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}
	wm := NewWaveletMatrixWithSum(nums, nums, -1, false)
	for i := int32(0); i < n-m+1; i++ {
		_, res := wm.KthValueAndSum(i, i+m, k, 0)
		fmt.Fprintln(out, res)
	}
}

// 子数组所有数到中位数的距离和(区间中位数距离和)
// https://yukicoder.me/problems/no/924
// n,q<=2e5
// -1e9 <= nums[i] <= 1e9
// 给定n个查询[l,r]
// !求区间[l,r]中位数到区间[l,r]中每个数的距离之和
// !也就求函数 f(x)= ∑|nums[i]-x| (l<=i<=right) 的最小值
// !区间中位数
func 区间最短距离和() {
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
	for i := 0; i < n; i++ {
		preSum[i+1] = preSum[i] + nums[i]
	}

	wm := NewWaveletMatrixWithSum(nums, nums, -1, false)
	for i := 0; i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--

		n := end - start
		lowerCount := n / 2
		ceilCount := n - lowerCount
		mid, lowerSum := wm.KthValueAndSum(start, end, lowerCount, 0)
		allSum := preSum[end] - preSum[start]
		ceilSum := allSum - lowerSum

		res := 0
		res += mid*int(lowerCount) - lowerSum
		res += ceilSum - mid*int(ceilCount)
		fmt.Fprintln(out, res)
	}
}

// No.919 You Are A Project Manager (区间中位数+前后缀分解)
// https://yukicoder.me/problems/no/919
// 给定一个数组nums，需要分为若干组.
// 初始时选择分组大小k，每次分组需要从最左边或最右边取k个数分为一组，
// 每个组的的得分为k*组内中位数(这里的中位数向下).
// 求最大化总得分.
//
// n<=1e4,-1e9<=nums[i]<1e9
func yuki919() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	wm := NewWaveletMatrixWithSum(nums, nil, -1, true)
	median := func(start, end int32) int {
		return wm.Median(start, end, false, 0)
	}

	res := -INF
	for size := int32(1); size <= n; size++ {
		count := n / size
		leftScore, rightScore := make([]int, count), make([]int, count)
		for i := int32(0); i < count; i++ {
			leftScore[i] = median(i*size, i*size+size)
			rightScore[i] = median(n-i*size-size, n-i*size)
		}
		leftPresum, rightPresum := make([]int, count+1), make([]int, count+1)
		for i := int32(0); i < count; i++ {
			leftPresum[i+1] = leftPresum[i] + leftScore[i]
			rightPresum[i+1] = rightPresum[i] + rightScore[i]
		}
		for i := int32(0); i < count; i++ {
			leftPresum[i+1] = max(leftPresum[i+1], leftPresum[i])
			rightPresum[i+1] = max(rightPresum[i+1], rightPresum[i])
		}
		for i := int32(0); i <= count; i++ {
			res = max(res, int(size)*(leftPresum[i]+rightPresum[count-i]))
		}
	}

	fmt.Fprintln(out, res)

}

// No.1332 Range Nearest Query
// https://yukicoder.me/problems/no/1332
// 区间最近值(无序数组区间前驱后继)查询
// !对每个查询[left,right,x],输出左闭右开区间内的数到给定数x的最小距离
// n<=3e5 q<=1e5 nums[i]<=1e9
func yuki1332() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	M := NewWaveletMatrixWithSum(nums, nil, -1, true)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int32
		var x int
		fmt.Fscan(in, &start, &end, &x)
		start--
		res := INF
		floor := M.Floor(start, end, x, 0) // 小于等于x的最大值
		if floor != -INF {
			res = min(res, abs(floor-x))
		}
		ceiling := M.Ceil(start, end, x, 0) // 大于等于x的最小值
		if ceiling != INF {
			res = min(res, abs(ceiling-x))
		}
		fmt.Fprintln(out, res)
	}
}

// No.2065 Sum of Min (区间chmin之和)
// https://yukicoder.me/problems/no/2065
// 小于x数的和+(end-start-小于x数的个数)*x
func yuki2065() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	M := NewWaveletMatrixWithSum(nums, nums, -1, true)
	for i := int32(0); i < q; i++ {
		var start, end int32
		var x int
		fmt.Fscan(in, &start, &end, &x)
		start--

		lessCount := M.CountRange(start, end, 0, x, 0)     // 小于x的数的个数
		lessSum := M.SumRange(start, end, 0, lessCount, 0) // 小于x的数的和
		res := lessSum + int((end-start-lessCount))*x
		fmt.Fprintln(out, res)
	}
}

const INF WmValue = 1e18

type WmValue = int
type WmSum = int

func (*WaveletMatrixWithSum) e() WmSum            { return 0 }
func (*WaveletMatrixWithSum) op(a, b WmSum) WmSum { return a + b }
func (*WaveletMatrixWithSum) inv(a WmSum) WmSum   { return -a }

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
	if k1 < 0 {
		k1 = 0
	}
	if k2 > end-start {
		k2 = end - start
	}
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

// 返回区间 [start, end) 中 [0, k) 的和.
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
