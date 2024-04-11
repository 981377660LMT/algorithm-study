package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	P1972()
}

// P1972 [SDOI2009] HH的项链
// https://www.luogu.com.cn/problem/P1972
// 区间颜色种类数
func P1972() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	D := NewDictionary()
	for i, v := range nums {
		nums[i] = D.Id(v)
	}

	// last[i] 表示第i个数之前与nums[i]相同的数中最大的下标
	last := make([]int32, n)
	lastIndex := make([]int32, D.Size())
	for i := range lastIndex {
		lastIndex[i] = -1
	}
	for i, v := range nums {
		last[i] = lastIndex[v]
		lastIndex[v] = int32(i)
	}

	for i := 0; i < len(last); i++ {
		last[i]++
	}

	// 那么现在查询询问的实际上是区间[start,end)中有多少个数满足pre[v]<start.
	wm := NewWaveletMatrix32(int32(n), func(i int32) int32 { return last[i] })
	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var start, end int32
		fmt.Fscan(in, &start, &end)
		start--
		fmt.Fprintln(out, start, end, 0, start)
		fmt.Fprintln(out, wm.CountRange(start, end, 0, start+1))
	}
}

type V = int32
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int32
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int32{},
	}
}
func (d *Dictionary) Id(value V) int32 {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := int32(len(d._idToValue))
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int32 {
	return int32(len(d._idToValue))
}

const INF int32 = 1e9 + 10

// 注意f(i)>=0.
func NewWaveletMatrix32(n int32, f func(i int32) int32) *WaveletMatrix32 {
	dataCopy := make([]int32, n)
	max_ := int32(0)
	for i := int32(0); i < n; i++ {
		v := f(i)
		if v > max_ {
			max_ = v
		}
		dataCopy[i] = v
	}
	maxLog := int32(bits.Len32(uint32(max_)) + 1)
	mat := make([]*BitVector32, maxLog)
	zs := make([]int32, maxLog)
	buff1 := make([]int32, maxLog)
	buff2 := make([]int32, maxLog)

	ls, rs := make([]int32, n), make([]int32, n)
	for dep := int32(0); dep < maxLog; dep++ {
		mat[dep] = NewBitVector32(n + 1)
		p, q := int32(0), int32(0)
		for i := int32(0); i < n; i++ {
			k := (dataCopy[i] >> (maxLog - dep - 1)) & 1
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
		for i := int32(0); i < q; i++ {
			dataCopy[p+i] = rs[i]
		}
	}

	return &WaveletMatrix32{
		n:      n,
		maxLog: maxLog,
		mat:    mat,
		zs:     zs,
		buff1:  buff1,
		buff2:  buff2,
	}
}

type WaveletMatrix32 struct {
	n            int32
	maxLog       int32
	mat          []*BitVector32
	zs           []int32
	buff1, buff2 []int32
}

// [start, end) 内的 value 的個数.
func (w *WaveletMatrix32) Count(start, end, value int32) int32 {
	return w.count(value, end) - w.count(value, start)
}

// [start, end) 内 [lower, upper) 的个数.
func (w *WaveletMatrix32) CountRange(start, end, lower, upper int32) int32 {
	return w.freqDfs(0, start, end, 0, lower, upper)
}

// 第k(0-indexed)个value的位置.
func (w *WaveletMatrix32) Index(value, k int32) int32 {
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

func (w *WaveletMatrix32) IndexWithStart(value, k, start int32) int32 {
	return w.Index(value, k+w.count(value, start))
}

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix32) Kth(start, end, k int32) int32 {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 内第k(0-indexed)大的数.
func (w *WaveletMatrix32) KthMax(start, end, k int32) int32 {
	if k < 0 || k >= end-start {
		return -1
	}
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
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

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix32) KthMin(start, end, k int32) int32 {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 中比 value 严格小的数, 不存在返回 -INF.
func (w *WaveletMatrix32) Lower(start, end, value int32) int32 {
	k := w.lt(start, end, value)
	if k != 0 {
		return w.KthMin(start, end, k-1)
	}
	return -INF
}

// [start, end) 中比 value 严格大的数, 不存在返回 INF.
func (w *WaveletMatrix32) Higher(start, end, value int32) int32 {
	k := w.le(start, end, value)
	if k == end-start {
		return INF
	}
	return w.KthMin(start, end, k)
}

// [start, end) 中不超过 value 的最大值, 不存在返回 -INF.
func (w *WaveletMatrix32) Floor(start, end, value int32) int32 {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Lower(start, end, value)
}

// [start, end) 中不小于 value 的最小值, 不存在返回 INF.
func (w *WaveletMatrix32) Ceiling(start, end, value int32) int32 {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Higher(start, end, value)
}

func (w *WaveletMatrix32) access(k int32) int32 {
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
		bit := w.mat[dep].Get(k)
		res = (res << 1) | bit
		k = w.mat[dep].Count(bit, k) + w.zs[dep]*dep
	}
	return res
}

func (w *WaveletMatrix32) count(value, end int32) int32 {
	left, right := int32(0), end
	for dep := int32(0); dep < w.maxLog; dep++ {
		w.buff1[dep] = left
		w.buff2[dep] = right
		bit := (value >> (w.maxLog - dep - 1)) & 1
		left = w.mat[dep].Count(bit, left) + w.zs[dep]*bit
		right = w.mat[dep].Count(bit, right) + w.zs[dep]*bit
	}
	return right - left
}

func (w *WaveletMatrix32) freqDfs(d, left, right, val, a, b int32) int32 {
	if left == right {
		return 0
	}
	if d == w.maxLog {
		if a <= val && val < b {
			return right - left
		}
		return 0
	}

	nv := (1 << (w.maxLog - d - 1)) | val
	nnv := ((1 << (w.maxLog - d - 1)) - 1) | nv
	if nnv < a || b <= val {
		return 0
	}
	if a <= val && nnv < b {
		return right - left
	}
	lc, rc := w.mat[d].Count(1, left), w.mat[d].Count(1, right)
	return w.freqDfs(d+1, left-lc, right-rc, val, a, b) + w.freqDfs(d+1, lc+w.zs[d], rc+w.zs[d], nv, a, b)
}

func (w *WaveletMatrix32) ll(left, right, v int32) (int32, int32) {
	res := int32(0)
	for dep := int32(0); dep < w.maxLog; dep++ {
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

func (w *WaveletMatrix32) lt(left, right, v int32) int32 {
	a, _ := w.ll(left, right, v)
	return a
}

func (w *WaveletMatrix32) le(left, right, v int32) int32 {
	a, b := w.ll(left, right, v)
	return a + b
}

type BitVector32 struct {
	n     int32
	block []int
	sum   []int
}

func NewBitVector32(n int32) *BitVector32 {
	blockCount := (n + 63) >> 6
	return &BitVector32{
		n:     n,
		block: make([]int, blockCount+1),
		sum:   make([]int, blockCount+1),
	}
}

func (f *BitVector32) Set(i int32) {
	f.block[i>>6] |= 1 << uint(i&63)
}

func (f *BitVector32) Build() {
	for i := 0; i < len(f.block)-1; i++ {
		f.sum[i+1] = f.sum[i] + bits.OnesCount(uint(f.block[i]))
	}
}

func (f *BitVector32) Get(i int32) int32 {
	return int32((f.block[i>>6] >> (i & 63))) & 1
}

func (f *BitVector32) Count(value, end int32) int32 {
	mask := (1 << uint(end&63)) - 1
	res := int32(f.sum[end>>6] + bits.OnesCount(uint(f.block[end>>6]&mask)))
	if value == 1 {
		return res
	}
	return end - res
}

func (f *BitVector32) Index(value, k int32) int32 {
	if k < 0 || f.Count(value, f.n) <= k {
		return -1
	}

	left, right := int32(0), f.n
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

func (f *BitVector32) IndexWithStart(value, k, start int32) int32 {
	return f.Index(value, k+f.Count(value, start))
}
