package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math/bits"
	"os"
	"sort"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

// # 長さ
// # N の正整数列
// # H=(H
// # 1
// # ​
// #  ,H
// # 2
// # ​
// #  ,…,H
// # N
// # ​
// #  ) が与えられます。

// # 長さ
// # N+1 の非負整数列
// # A=(A
// # 0
// # ​
// #  ,A
// # 1
// # ​
// #  ,…,A
// # N
// # ​
// #  ) があります。 はじめ、
// # A
// # 0
// # ​
// #  =A
// # 1
// # ​
// #  =⋯=A
// # N
// # ​
// #  =0 です。

// # A に対して、次の操作を繰り返します。

// # A
// # 0
// # ​
// #   の値を
// # 1 増やす。
// # i=1,2,…,N に対して、この順に次の操作を行う。
// # A
// # i−1
// # ​
// #  >A
// # i
// # ​
// #   かつ
// # A
// # i−1
// # ​
// #  >H
// # i
// # ​
// #   のとき、
// # A
// # i−1
// # ​
// #   の値を
// # 1 減らし、
// # A
// # i
// # ​
// #   の値を
// # 1 増やす。
// # i=1,2,…,N のそれぞれに対して、初めて
// # A
// # i
// # ​
// #  >0 が成り立つのは何回目の操作の後か求めてください。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	n := io.NextInt()
	heights := make([]int, n)
	for i := 0; i < n; i++ {
		heights[i] = io.NextInt()
	}

	wm := NewWaveletMatrixWithSum(heights, heights, -1, true)

}

const INF WmValue = 1e18

type WmValue = int
type WmSum = int

func (*WaveletMatrixWithSum) e() WmSum            { return 0 }
func (*WaveletMatrixWithSum) op(a, b WmSum) WmSum { return a + b }
func (*WaveletMatrixWithSum) inv(a WmSum) WmSum   { return -a }

type WaveletMatrixWithSum struct {
	n, log   int32
	setLog   bool
	compress bool
	useSum   bool
	mid      []int32
	bv       []*BitVector
	key      []WmValue
	presum   [][]WmSum
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
	wm.setLog = log != -1
	wm.compress = compress
	wm.useSum = len(sumData) > 0
	if wm.n == 0 {
		wm.log = 0
		wm.presum = [][]WmSum{{wm.e()}}
		return
	}

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
		wm.key = wm.key[:len(wm.key):len(wm.key)]
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
	if wm.useSum {
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
		if wm.useSum {
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
				if wm.useSum {
					S0[p0] = S[i]
				}
				A0[p0] = A[i]
				p0++
			} else {
				if wm.useSum {
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

// 返回区间 [start, end) 中 值在 [a, b) 中的元素个数以及这些元素的和.
func (wm *WaveletMatrixWithSum) RangeCountAndSum(start, end int32, a, b WmValue, xorValue WmValue) (int32, WmSum) {
	if xorValue != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || a >= b {
		return 0, wm.e()
	}
	if wm.compress {
		a = wm._lowerBound(wm.key, a)
		b = wm._lowerBound(wm.key, b)
	}
	count, sum := int32(0), wm.e()
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if rx <= a || b <= lx {
			return
		}
		if a <= lx && rx <= b {
			count += r - l
			if wm.useSum {
				sum = wm.op(sum, wm._get(d, l, r))
			}
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xorValue>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l0, r0, lx, mx)
		dfs(d, l1, r1, mx, rx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	return count, sum
}

// 返回区间 [start, end) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果).
// 如果k >= end-start, 返回 (-1, 区间 op 的结果).
func (wm *WaveletMatrixWithSum) KthValueAndSum(start, end, k int32, xorVal WmValue) (WmValue, WmSum) {
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
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

	sum, val := wm.e(), WmValue(0)
	for d := wm.log - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		l1, r1 := start+wm.mid[d]-l0, end+wm.mid[d]-r0
		if (xorVal>>d)&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		if k < r0-l0 {
			start, end = l0, r0
		} else {
			k -= r0 - l0
			val |= 1 << d
			start, end = l1, r1
			if wm.useSum {
				sum = wm.op(sum, wm._get(d, l0, r0))
			}
		}
	}
	if wm.useSum {
		sum = wm.op(sum, wm._get(0, start, start+k))
	}
	if wm.compress {
		val = wm.key[val]
	}
	return val, sum
}

// [start, end)区间内第k(k>=0)小的元素.
func (wm *WaveletMatrixWithSum) Kth(start, end, k int32, xorVal WmValue) WmValue {
	if k < 0 {
		k = 0
	}
	if n := end - start - 1; k > n {
		k = n
	}
	v, _ := wm.KthValueAndSum(start, end, k, xorVal)
	return v
}

// upper: 向上取中位数还是向下取中位数.
func (wm *WaveletMatrixWithSum) Median(start, end int32, upper bool, xorVal WmValue) WmValue {
	n := end - start
	var k int32
	if upper {
		k = n >> 1
	} else {
		k = (n - 1) >> 1
	}
	return wm.Kth(start, end, k, xorVal)
}

// [start, end) 中小于等于 x 的数中最大的数.
//
//	如果不存在则返回-INF.
func (wm *WaveletMatrixWithSum) Floor(start, end int32, x WmValue, xor WmValue) WmValue {
	if xor != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return -INF
	}
	res := -INF
	x++
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if rx-1 <= res || l == r || x <= lx {
			return
		}
		if d == 0 {
			res = max(res, lx)
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xor>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l1, r1, mx, rx)
		dfs(d, l0, r0, lx, mx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	if wm.compress && res != -INF {
		res = wm.key[res]
	}
	return res
}

// [start, end) 中大于等于 x 的数中最小的数
//
//	如果不存在则返回INF
func (wm *WaveletMatrixWithSum) Ceil(start, end int32, x WmValue, xor WmValue) int {
	if xor != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return INF
	}
	if wm.compress {
		x = wm._lowerBound(wm.key, x)
	}
	res := INF
	var dfs func(d, l, r int32, lx, rx WmValue)
	dfs = func(d, l, r int32, lx, rx WmValue) {
		if res <= lx || l == r || rx <= x {
			return
		}
		if d == 0 {
			res = min(res, lx)
			return
		}
		d--
		mx := (lx + rx) >> 1
		l0, r0 := wm.bv[d].Rank(l, false), wm.bv[d].Rank(r, false)
		l1, r1 := l+wm.mid[d]-l0, r+wm.mid[d]-r0
		if xor>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		dfs(d, l0, r0, lx, mx)
		dfs(d, l1, r1, mx, rx)
	}
	dfs(wm.log, start, end, 0, 1<<wm.log)
	if wm.compress && res < INF {
		res = wm.key[res]
	}
	return res
}

// 返回区间 [start, end) 中 范围在 [a, b) 中的元素的和.
func (wm *WaveletMatrixWithSum) SumRange(start, end int32, a, b WmValue, xorVal WmValue) WmSum {
	if !wm.useSum {
		panic("sum data must be provided")
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || a >= b {
		return wm.e()
	}
	_, sum := wm.RangeCountAndSum(start, end, a, b, xorVal)
	return sum
}

// 返回区间 [start, end) 中 排名在 [k1, k2) 中的元素的和.
func (wm *WaveletMatrixWithSum) SumSlice(start, end, k1, k2 int32, xorVal WmValue) WmSum {
	if !wm.useSum {
		panic("sum data must be provided")
	}
	if k1 < 0 {
		k1 = 0
	}
	if k2 > end-start {
		k2 = end - start
	}
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end || k1 >= k2 {
		return wm.e()
	}
	_, sum1 := wm.KthValueAndSum(start, end, k1, xorVal)
	_, sum2 := wm.KthValueAndSum(start, end, k2, xorVal)
	return wm.op(sum2, wm.inv(sum1))
}

func (wm *WaveletMatrixWithSum) SumAll(start, end int32) WmSum {
	if start < 0 {
		start = 0
	}
	if end > wm.n {
		end = wm.n
	}
	if start >= end {
		return wm.e()
	}
	return wm._get(wm.log, start, end)
}

// 使得predicate(count, sum)为true的最大的(count, sum).
func (wm *WaveletMatrixWithSum) MaxRight(predicate func(int32, WmSum) bool, start, end int32, xorVal WmValue) (int32, WmSum) {
	if xorVal != 0 {
		if !wm.setLog {
			panic("log should be set when xor is used")
		}
	}
	if start >= end {
		return 0, wm.e()
	}
	if s := wm._get(wm.log, start, end); predicate(end-start, s) {
		return end - start, s
	}
	count, sum := int32(0), wm.e()
	for d := wm.log - 1; d >= 0; d-- {
		l0, r0 := wm.bv[d].Rank(start, false), wm.bv[d].Rank(end, false)
		l1, r1 := start+wm.mid[d]-l0, end+wm.mid[d]-r0
		if xorVal>>d&1 == 1 {
			l0, l1 = l1, l0
			r0, r1 = r1, r0
		}
		if s := wm.op(sum, wm._get(d, l0, r0)); predicate(count+r0-l0, s) {
			count += r0 - l0
			sum = s
			start, end = l1, r1
		} else {
			start, end = l0, r0
		}
	}
	k := wm._binarySearch(func(k int32) bool {
		return predicate(count+k, wm.op(sum, wm._get(0, start, start+k)))
	}, 0, end-start)
	count += k
	sum = wm.op(sum, wm._get(0, start, start+k))
	return count, sum
}

func (wm *WaveletMatrixWithSum) _get(d, l, r int32) WmSum {
	if wm.useSum {
		return wm.op(wm.presum[d][r], wm.inv(wm.presum[d][l]))
	}
	return wm.e()
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
