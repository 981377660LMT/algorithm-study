// Usage:
// !Attention: nums[i] >= 0

// Count(start, end, value) – [start, end) 中值为 value 的数的个数.
// CountRange(start, end, lower, upper) – [start, end) 中值在 [lower, upper) 内的数的个数.

// Kth(start, end, k) – [start, end) 中第 k(0-indexed) 小的数.
// KthMax(start, end, k) – [start, end) 中第 k(0-indexed) 大的数.

// Lower(start, end, value) – [start, end) 中值小于 value 的最大值.不存在的话返回 -INF.
// Higher(start, end, value) – [start, end) 中值大于 value 的最小值.不存在的话返回 INF.
// Floor(start, end, value) – [start, end) 中值不超过 value 的最大值.不存在的话返回 -INF.
// Ceiling(start, end, value) – [start, end) 中值不小于 value 的最小值.不存在的话返回 INF.

// Referece:
// https://beet-aizu.github.io/library/datastructure/waveletmatrix.cpp
// https://blog.hamayanhamayan.com/entry/2017/06/13/103352

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

func main() {
	CF813E()
}

func demo() {
	nums := []int{1, 2, 3, 1, 5, 6, 7, 8, 9, 10}
	M := NewWaveletMatrix(nums)

	fmt.Println(M.Count(0, 1, 1))
	fmt.Println(M.CountRange(0, 10, 1, 5))
	fmt.Println(M.Index(1, 1))
	fmt.Println(M.IndexWithStart(1, 0, 2))
	fmt.Println(M.Kth(0, 10, 2))
	fmt.Println(M.KthMax(0, 10, 2))
	fmt.Println(M.Lower(0, 3, 2))
	fmt.Println(M.Floor(0, 3, 2))
	fmt.Println(M.Higher(0, 10, 1))
	fmt.Println(M.Ceiling(0, 10, 1))
}

const INF int = 1e18

// Army Creation
// https://www.luogu.com.cn/problem/CF813E
// 有n个数，q次询问，每次询问[l,r]中最多可以选多少个数使得相同的数最多有k个，强制在线.
// !k不变，可以预处理出每个数前面的第k个相同的数在哪.
// !对于询问区间[l,r]，若前面的第k个数的位置小于l，则可以选择这个数
func CF813E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	preRes := 0
	parse := func(v int) int {
		return (v+preRes)%n + 1
	}

	mp := make(map[int][]int)
	for i, v := range nums {
		mp[v] = append(mp[v], i)
	}
	maxStart := make([]int, n) // 选择下标为i的数时，区间的左端点需要严格小于maxStart[i]
	for _, pos := range mp {
		for i := k; i < len(pos); i++ {
			maxStart[pos[i]] = pos[i-k] + 1
		}
	}
	wm := NewWaveletMatrix(maxStart)
	query := func(start, end int) int {
		return wm.CountRange(start, end, 0, start+1)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var start, end int
		fmt.Fscan(in, &start, &end)
		start = parse(start)
		end = parse(end)
		if start > end {
			start, end = end, start
		}
		start--
		preRes = query(start, end)
		fmt.Fprintln(out, preRes)
	}
}

// 给你一个长度为 n 下标从 0 开始的整数数组 nums ，它包含 1 到 n 的所有数字，请你返回上升四元组的数目。
// 如果一个四元组 (i, j, k, l) 满足以下条件，我们称它是上升的：
// 0 <= i < j < k < l < n 且
// nums[i] < nums[k] < nums[j] < nums[l] 。
// 4 <= nums.length <= 4000, nums 中所有数字 互不相同 ，nums 是一个排列。
// https://leetcode.cn/problems/count-increasing-quadruplets/
func countQuadruplets(nums []int) int64 {
	W := NewWaveletMatrix(nums)
	res := 0
	for j := 1; j < len(nums)-2; j++ {
		for k := j + 1; k < len(nums)-1; k++ {
			if !(nums[k] < nums[j]) {
				continue
			}
			left := W.CountRange(0, j, 0, nums[k])
			right := W.CountRange(k+1, len(nums), nums[j]+1, INF)
			res += (left * right)
		}
	}
	return int64(res)
}

// 给定非负整数数组 nums 构建一个 WaveletMatrix.
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

// [start, end) 内的 value 的個数.
func (w *WaveletMatrix) Count(start, end, value int) int {
	return w.count(value, end) - w.count(value, start)
}

// [start, end) 内 [lower, upper) 的个数.
func (w *WaveletMatrix) CountRange(start, end, lower, upper int) int {
	return w.freqDfs(0, start, end, 0, lower, upper)
}

// 第k(0-indexed)个value的位置.
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

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix) Kth(start, end, k int) int {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 内第k(0-indexed)大的数.
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

// [start, end) 内第k(0-indexed)小的数(kth).
func (w *WaveletMatrix) KthMin(start, end, k int) int {
	return w.KthMax(start, end, end-start-k-1)
}

// [start, end) 中比 value 严格小的数, 不存在返回 -INF.
func (w *WaveletMatrix) Lower(start, end, value int) int {
	k := w.lt(start, end, value)
	if k != 0 {
		return w.KthMin(start, end, k-1)
	}
	return -INF
}

// [start, end) 中比 value 严格大的数, 不存在返回 INF.
func (w *WaveletMatrix) Higher(start, end, value int) int {
	k := w.le(start, end, value)
	if k == end-start {
		return INF
	}
	return w.KthMin(start, end, k)
}

// [start, end) 中不超过 value 的最大值, 不存在返回 -INF.
func (w *WaveletMatrix) Floor(start, end, value int) int {
	count := w.Count(start, end, value)
	if count > 0 {
		return value
	}
	return w.Lower(start, end, value)
}

// [start, end) 中不小于 value 的最小值, 不存在返回 INF.
func (w *WaveletMatrix) Ceiling(start, end, value int) int {
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

func (f *BitVector) Count(value, end int) int {
	mask := (1 << uint(end&63)) - 1
	res := f.sum[end>>6] + bits.OnesCount(uint(f.block[end>>6]&mask))
	if value == 1 {
		return res
	}
	return end - res
}

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
