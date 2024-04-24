// 区间最近值(无序数组区间前驱后继)查询
// !对每个查询[left,right,x],输出左闭右开区间内的数到给定数x的最小距离
// n<=3e5 q<=1e5 nums[i]<=1e9

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

	M := NewWaveletMatrix(nums)
	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var left, right, x int
		fmt.Fscan(in, &left, &right, &x)
		left--
		res := INF
		lower := M.Floor(left, right, x) // 小于等于x的最大值
		if lower != -INF {
			res = min(res, abs(lower-x))
		}
		higher := M.Ceil(left, right, x) // 大于等于x的最小值
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

// 指定された配列から WaveletMatrix を構築する.
//
//	data:変換する配列(data[i]は0以上)
func NewWaveletMatrix(data []int) *WaveletMatrix {
	dataCopy := make([]int, len(data))
	max_ := 0
	for i, v := range data {
		if v > max_ {
			max_ = v
		}
		dataCopy[i] = v
	}
	maxLog := bits.Len(uint(max_)) + 1
	n := len(dataCopy)
	mat := make([]*BitVector, maxLog)
	zs := make([]int, maxLog)
	buff1 := make([]int, maxLog)
	buff2 := make([]int, maxLog)

	ls, rs := make([]int, n), make([]int, n)
	for dep := 0; dep < maxLog; dep++ {
		mat[dep] = NewBitVector(n + 1)
		p, q := 0, 0
		for i := 0; i < n; i++ {
			k := (dataCopy[i] >> uint(maxLog-dep-1)) & 1
			if k == 1 {
				rs[q] = dataCopy[i]
				mat[dep].Set(i)
				q++
			} else {
				ls[p] = dataCopy[i]
				p++
			}
		}

		zs[dep] = p
		mat[dep].Build()
		ls = dataCopy
		for i := 0; i < q; i++ {
			dataCopy[p+i] = rs[i]
		}
	}

	return &WaveletMatrix{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		buff1:  buff1,
		buff2:  buff2,
	}
}

type WaveletMatrix struct {
	n            int
	maxLog       int
	mat          []*BitVector
	zs           []int
	buff1, buff2 []int
}

// [start, end) に含まれる value の個数を求める.
//
//	alias: Rank
func (w *WaveletMatrix) Count(start, end, value int) int {
	return w.count(value, end) - w.count(value, start)
}

// [start, end) に含まれる [value, upper) の個数を求める.
// alias: RangeFreq
func (w *WaveletMatrix) CountRange(start, end, lower, upper int) int {
	return w.freqDfs(0, start, end, 0, lower, upper)
}

// k(0-indexed) 番目の value の位置を求める.
//
//	alias: Select
func (w *WaveletMatrix) Index(value, k int) int {
	w.count(value, w.n)
	for dep := w.maxLog - 1; dep >= 0; dep-- {
		bit := (value >> uint(w.maxLog-dep-1)) & 1
		k = w.mat[dep].IndexWithStart(bit, k, w.buff1[dep])
		if k < 0 || k >= w.buff2[dep] {
			return -1
		}
		k -= w.buff1[dep]
	}
	return k
}

func (w *WaveletMatrix) IndexWithStart(value, k, start int) int {
	return w.Index(value, k+w.count(value, start))
}

// [start, end) に含まれる要素の中で k(0-indexed) 番目に大きいものを求める.
//
//	alias: Quantile
func (w *WaveletMatrix) KthMax(start, end, k int) int {
	if k < 0 || k >= end-start {
		return -1
	}
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		p, q := w.mat[dep].Count(1, start), w.mat[dep].Count(1, end)
		if k < q-p {
			start = w.zs[dep] + p
			end = w.zs[dep] + q
			res |= 1 << uint(w.maxLog-dep-1)
		} else {
			k -= q - p
			start -= p
			end -= q
		}
	}
	return res
}

// [start, end) に含まれる要素の中で k(0-indexed) 番目に小さいものを求める.
//
//	alias: Rquantile
func (w *WaveletMatrix) KthMin(start, end, k int) int {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) に含まれる要素の中で value の次に小さいものを求める.存在しない場合は -INF を返す.
//
//	value >= 0
//	alias: Pred
func (w *WaveletMatrix) Lower(start, end, value int) int {
	k := w.lt(start, end, value)
	if k != 0 {
		return w.KthMin(start, end, k-1)
	}
	return -INF
}

// [start, end) に含まれる要素の中で value より大きいものを求める.存在しない場合は INF を返す.
//
//	value >= 0
//	alias: Succ
func (w *WaveletMatrix) Higher(start, end, value int) int {
	k := w.le(start, end, value)
	if k == end-start {
		return INF
	}
	return w.KthMin(start, end, k)
}

// [start, end) に含まれる要素の中で value 以下のものを求める.存在しない場合は -INF を返す.
func (w *WaveletMatrix) Floor(start, end, value int) int {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Lower(start, end, value)
}

// [start, end) に含まれる要素の中で value 以上のものを求める.存在しない場合は INF を返す.
func (w *WaveletMatrix) Ceil(start, end, value int) int {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Higher(start, end, value)
}

func (w *WaveletMatrix) access(k int) int {
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		bit := w.mat[dep].Get(k)
		res = (res << 1) | bit
		k = w.mat[dep].Count(bit, k) + w.zs[dep]*dep
	}
	return res
}

func (w *WaveletMatrix) count(value, end int) int {
	left, right := 0, end
	for dep := 0; dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := (value >> uint(w.maxLog-dep-1)) & 1
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return right - left
}

func (w *WaveletMatrix) freqDfs(d, left, right, val, a, b int) int {
	if left == right {
		return 0
	}
	if d == w.maxLog {
		if a <= val && val < b {
			return right - left
		}
		return 0
	}

	nv := (1 << uint(w.maxLog-d-1)) | val
	nnv := ((1 << uint(w.maxLog-d-1)) - 1) | nv
	if nnv < a || b <= val {
		return 0
	}
	if a <= val && nnv < b {
		return right - left
	}
	lc, rc := w.mat[d].Count(1, left), w.mat[d].Count(1, right)
	return w.freqDfs(d+1, left-lc, right-rc, val, a, b) + w.freqDfs(d+1, lc+w.zs[d], rc+w.zs[d], nv, a, b)
}

func (w *WaveletMatrix) ll(left, right, v int) (int, int) {
	res := 0
	for dep := 0; dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := v >> uint(w.maxLog-dep-1) & 1
		if bit == 1 {
			res += right - left + w.mat[dep].Count(1, left) - w.mat[dep].Count(1, right)
		}
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return res, right - left
}

func (w *WaveletMatrix) lt(left, right, v int) int {
	a, _ := w.ll(left, right, v)
	return a
}

func (w *WaveletMatrix) le(left, right, v int) int {
	a, b := w.ll(left, right, v)
	return a + b
}

type BitVector struct {
	n     int
	block []int
	sum   []int
}

func NewBitVector(n int) *BitVector {
	blockCount := (n + 63) >> 6
	return &BitVector{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
	}
}

func (f *BitVector) Set(i int) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *BitVector) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

func (f *BitVector) Get(i int) int {
	return (f.block[i>>6] >> uint(i&63)) & 1
}

// [0,end) に含まれる 1 の個数.
func (f *BitVector) Count(value, end int) int {
	if value == 1 {
		return f.count1(end)
	}
	return end - f.count1(end)
}

// [0,end) に k(0-indexed) 番目の value の位置を求める.
// 存在しない場合は -1 を返す.
func (f *BitVector) Index(value, k int) int {
	if k < 0 || f.Count(value, f.n) <= k {
		return -1
	}

	left, right := 0, f.n
	for right-left > 1 {
		mid := (left + right) >> 1
		if f.Count(value, mid) >= k+1 {
			right = mid
		} else {
			left = mid
		}
	}
	return right - 1
}

func (f *BitVector) IndexWithStart(value, k, start int) int {
	return f.Index(value, k+f.Count(value, start))
}

func (f *BitVector) count1(k int) int {
	mask := (1 << uint(k&63)) - 1
	return f.sum[k>>6] + bits.OnesCount(uint(f.block[k>>6]&mask))
}
