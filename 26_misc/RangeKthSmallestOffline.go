// 离线求区间第k小
// RangeKthSmallestOffline
// !速度不如 waveletMatrix

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
	yosupo()
}

func yosupo() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}

	rks := NewRangeKthSmallestOffline(nums)
	for i := 0; i < q; i++ {
		var l, r, k int
		fmt.Fscan(in, &l, &r, &k)
		rks.AddQuery(l, r, k)
	}

	indexes := rks.Run()
	for _, idx := range indexes {
		fmt.Fprintln(out, nums[idx])
	}
}

// Offline solver to find k-th smallest elements in continuous subsequences
// - Problem statement: <https://judge.yosupo.jp/problem/range_kth_smallest>
// - Complexity: O((N + Q) lg(N + Q))
type RangeKthSmallestOffline struct {
	n          int32
	order      []int32
	ql, qr, qk []int32
}

func NewRangeKthSmallestOffline(nums []int) *RangeKthSmallestOffline {
	res := &RangeKthSmallestOffline{}
	res.n = int32(len(nums))
	res.order = make([]int32, res.n)
	for i := int32(0); i < res.n; i++ {
		res.order[i] = i
	}
	sort.Slice(res.order, func(i, j int) bool {
		return nums[res.order[i]] < nums[res.order[j]]
	})
	return res
}

// Add a query to find (k + 1)-th smallest value in [l, r)
func (r *RangeKthSmallestOffline) AddQuery(start, end, k int) {
	if !(start >= 0 && start <= end && end <= int(r.n)) {
		panic("invalid query")
	}
	if !(end-start > k) {
		panic("invalid query")
	}
	r.ql = append(r.ql, int32(start))
	r.qr = append(r.qr, int32(end))
	r.qk = append(r.qk, int32(k))
}

// Solve all queries: ret[q] = (position of the K[q]-th smallest element in [L[q], R[q]))
func (r *RangeKthSmallestOffline) Run() (indexes []int) {
	q := int32(len(r.ql))
	lo := make([]int32, q)
	hi := make([]int32, q)
	for i := int32(0); i < q; i++ {
		hi[i] = r.n
	}
	vs := make([][]int32, r.n)

	for {
		stop := true
		for i := int32(0); i < q; i++ {
			if lo[i]+1 < hi[i] {
				stop = false
				vs[(lo[i]+hi[i])/2] = append(vs[(lo[i]+hi[i])/2], i)
			}
		}
		if stop {
			break
		}
		bit := NewBITArray0132(r.n)
		for t := int32(0); t < r.n; t++ {
			for _, i := range vs[t] {
				if bit.QueryRange(r.ql[i], r.qr[i]) <= r.qk[i] {
					lo[i] = t
				} else {
					hi[i] = t
				}
			}
			bit.Add(r.order[t])
			vs[t] = vs[t][:0]
		}
	}

	indexes = make([]int, q)
	for i := int32(0); i < q; i++ {
		indexes[i] = int(r.order[lo[i]])
	}
	return
}

// 01树状数组.
type BITArray0132 struct {
	n    int32
	size int32 // data、bit的长度
	data []uint64
	bit  *BITArray32
}

func NewBITArray0132(n int32) *BITArray0132 {
	size := int32(n>>6 + 1)
	data := make([]uint64, size)
	bit := NewBitArray(size)
	return &BITArray0132{n: n, size: size, data: data, bit: bit}
}

func (bit01 *BITArray0132) QueryPrefix(end int32) int32 {
	i, j := end>>6, end&63
	res := bit01.bit.QueryPrefix(i)
	res += int32(bits.OnesCount64(bit01.data[i] & ((1 << j) - 1)))
	return res
}

func (bit01 *BITArray0132) QueryRange(start, end int32) int32 {
	if start >= end {
		return 0
	}
	if start == 0 {
		return bit01.QueryPrefix(end)
	}
	res := int32(0)
	res -= int32(bits.OnesCount64(bit01.data[start>>6] & ((1 << (start & 63)) - 1)))
	res += int32(bits.OnesCount64(bit01.data[end>>6] & ((1 << (end & 63)) - 1)))
	res += bit01.bit.QueryRange(start>>6, end>>6)
	return res
}

func (bit01 *BITArray0132) Add(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 1 {
		return false
	}
	bit01.data[i] |= 1 << j
	bit01.bit.Add(i, 1)
	return true
}

func (bit01 *BITArray0132) Remove(index int32) bool {
	i, j := index>>6, index&63
	if (bit01.data[i]>>j)&1 == 0 {
		return false
	}
	bit01.data[i] ^= 1 << j
	bit01.bit.Add(i, -1)
	return true
}

func (bit01 *BITArray0132) Has(index int32) bool {
	i, j := index>>6, index&63
	return (bit01.data[i]>>j)&1 == 1
}

// !Point Add Range Sum, 0-based.
type BITArray32 struct {
	n    int32
	data []int32
}

func NewBitArray(n int32) *BITArray32 {
	res := &BITArray32{n: n, data: make([]int32, n)}
	return res
}

func (b *BITArray32) Add(index int32, v int32) {
	for index++; index <= b.n; index += index & -index {
		b.data[index-1] += v
	}
}

// [0, end).
func (b *BITArray32) QueryPrefix(end int32) int32 {
	res := int32(0)
	for ; end > 0; end -= end & -end {
		res += b.data[end-1]
	}
	return res
}

// [start, end).
func (b *BITArray32) QueryRange(start, end int32) int32 {
	if start >= end {
		return 0
	}
	if start == 0 {
		return b.QueryPrefix(end)
	}
	pos, neg := int32(0), int32(0)
	for end > start {
		pos += b.data[end-1]
		end &= end - 1
	}
	for start > end {
		neg += b.data[start-1]
		start &= start - 1
	}
	return pos - neg
}
