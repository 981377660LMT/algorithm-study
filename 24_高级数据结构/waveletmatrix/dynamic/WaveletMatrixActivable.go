// // https://github.com/programming-team-code/programming_team_code/blob/0c4cacf62f1d1e8e727610f4d92a8238e95743e0/data_structures/wavelet_merge/wavelet_tree_updates.hpp#L1
// // TODO 有问题
// // WaveletMatrixActivable.go
// // Api:
// //
// //	NewWaveletMatrixActivable(nums []int32, minV, maxV int32, active func(i int32) bool) *WaveletMatrixActivable
// //	SetActive(i int32, isActive bool)
// //	RectCount(le, ri, x, y int32) int32
// //	KthSmallest(le, ri, k int32) int32

// package main

// import (
// 	"bufio"
// 	"fmt"
// 	"math/bits"
// 	"os"
// 	"sort"
// )

// func main() {
// 	yosupo()
// 	// demo()
// }

// func demo() {
// 	nums := []int{1, 4, 0, 1, 3}
// 	newNums, origin := Discretize(nums)
// 	W := NewWaveletMatrixActivable(newNums, 0, int32(len(origin)), func(i int32) bool { return true })
// 	fmt.Println(newNums, origin)
// 	// fmt.Println(W.RectCount(0, 5, 1, 3))
// 	// fmt.Println(W.KthSmallest(0, 5, 2))
// 	fmt.Println(origin[W.KthSmallest(1, 3, 2)])
// 	// fmt.Println(W.KthSmallest(3, 4, 1))
// }

// // https://judge.yosupo.jp/problem/range_kth_smallest
// func yosupo() {
// 	in := bufio.NewReader(os.Stdin)
// 	out := bufio.NewWriter(os.Stdout)
// 	defer out.Flush()

// 	var n, q int
// 	fmt.Fscan(in, &n, &q)
// 	nums := make([]int, n)
// 	for i := 0; i < n; i++ {
// 		fmt.Fscan(in, &nums[i])
// 	}

// 	newNums, origin := Discretize(nums)
// 	wm := NewWaveletMatrixActivable(
// 		newNums, 0, int32(len(origin)),
// 		func(i int32) bool { return true },
// 	)
// 	for i := 0; i < q; i++ {
// 		var start, end, x int32
// 		fmt.Fscan(in, &start, &end, &x)
// 		res := wm.KthSmallest(start, end, x+1)
// 		fmt.Fprintln(out, origin[res])
// 	}
// }

// type WaveletMatrixActivable struct {
// 	n, minV, maxV int32
// 	boolPresums   []*boolPresum
// 	boolBits      []*boolBit
// }

// type pair = struct {
// 	first  int32
// 	second bool
// }

// // 传入前需要离散化.
// // minV <= nums[i] < maxV
// func NewWaveletMatrixActivable(
// 	nums []int32,
// 	minV, maxV int32,
// 	active func(i int32) bool,
// ) *WaveletMatrixActivable {
// 	n := int32(len(nums))
// 	res := &WaveletMatrixActivable{
// 		n: n, minV: minV, maxV: maxV,
// 		boolPresums: make([]*boolPresum, maxV-minV),
// 		boolBits:    make([]*boolBit, 2*(maxV-minV)),
// 	}
// 	cpy := make([]pair, n)
// 	for i := int32(0); i < n; i++ {
// 		cpy[i] = pair{first: nums[i], second: active(i)}
// 	}
// 	res.build(cpy, 0, n, minV, maxV, 1)
// 	return res
// }

// // O(log(maxV-minV) * log(n/64)
// func (w *WaveletMatrixActivable) SetActive(i int32, isActive bool) {
// 	if w.boolBits[1].On(i) == isActive {
// 		return
// 	}
// 	w.setActiveImpl(i, isActive, w.minV, w.maxV, 1)
// }

// // O(log(maxV-minV) * log(n/64)
// func (w *WaveletMatrixActivable) RectCount(le, ri, x, y int32) int32 {
// 	return w.rectCountImpl(le, ri, x, y, w.minV, w.maxV, 1)
// }

// // O(log(maxV-minV) * log(n/64)
// // !1<=k<=ri-le
// func (w *WaveletMatrixActivable) KthSmallest(le, ri, k int32) int32 {
// 	return w.kthSmallestImpl(le, ri, k, w.minV, w.maxV, 1)
// }

// func (w *WaveletMatrixActivable) build(cpy []pair, le, ri, tl, tr, v int32) {
// 	bools := make([]bool, ri-le)
// 	for i := le; i < ri; i++ {
// 		bools[i-le] = cpy[i].second
// 	}
// 	w.boolBits[v] = newBoolBit(ri-le, func(i int32) bool { return bools[i] })
// 	if tr-tl <= 1 {
// 		return
// 	}
// 	tm := split(tl, tr)
// 	for i := le; i < ri; i++ {
// 		bools[i-le] = cpy[i].first < tm
// 	}
// 	w.boolPresums[v] = newBoolPresum(ri-le, func(i int32) bool { return bools[i] })
// 	mi := stablePartition(cpy, le, ri, func(p pair) bool { return p.first < tm })
// 	w.build(cpy, le, mi, tl, tm, 2*v)
// 	w.build(cpy, mi, ri, tm, tr, 2*v+1)
// }

// func (w *WaveletMatrixActivable) setActiveImpl(i int32, isActive bool, tl, tr, v int32) {
// 	w.boolBits[v].Set(i, isActive)
// 	if tr-tl == 1 {
// 		return
// 	}
// 	tm := split(tl, tr)
// 	pi := w.boolPresums[v].Popcount(i)
// 	if w.boolPresums[v].On(i) {
// 		w.setActiveImpl(pi, isActive, tl, tm, 2*v)
// 	} else {
// 		w.setActiveImpl(i-pi, isActive, tm, tr, 2*v+1)
// 	}
// }

// func (w *WaveletMatrixActivable) rectCountImpl(le, ri, x, y, tl, tr, v int32) int32 {
// 	if y <= tl || tr <= x {
// 		return 0
// 	}
// 	if x <= tl && tr <= y {
// 		return w.boolBits[v].PopcountRange(le, ri)
// 	}
// 	tm := split(tl, tr)
// 	pl := w.boolPresums[v].Popcount(le)
// 	pr := w.boolPresums[v].Popcount(ri)
// 	return w.rectCountImpl(pl, pr, x, y, tl, tm, 2*v) +
// 		w.rectCountImpl(le-pl, ri-pr, x, y, tm, tr, 2*v+1)
// }

// func (w *WaveletMatrixActivable) kthSmallestImpl(le, ri, k, tl, tr, v int32) int32 {
// 	if tr-tl == 1 {
// 		return tl
// 	}
// 	tm := split(tl, tr)
// 	pl := w.boolPresums[v].Popcount(le)
// 	pr := w.boolPresums[v].Popcount(ri)
// 	cntLeft := w.boolBits[2*v].PopcountRange(pl, pr)
// 	if k <= cntLeft {
// 		return w.kthSmallestImpl(pl, pr, k, tl, tm, 2*v)
// 	}
// 	return w.kthSmallestImpl(le-pl, ri-pr, k-cntLeft, tm, tr, 2*v+1)
// }

// type boolPresum struct {
// 	n      int32
// 	presum []int32
// 	mask   []uint64
// }

// func newBoolPresum(n int32, f func(int32) bool) *boolPresum {
// 	mask := make([]uint64, n>>6+1)
// 	for i := int32(0); i < n; i++ {
// 		if f(i) {
// 			mask[i>>6] |= 1 << (i & 63)
// 		}
// 	}
// 	presum := make([]int32, n>>6+1)
// 	for i := 0; i < len(mask)-1; i++ {
// 		presum[i+1] = int32(bits.OnesCount64(mask[i])) + presum[i]
// 	}
// 	return &boolPresum{n: n, presum: presum, mask: mask}
// }

// func (b *boolPresum) Popcount(i int32) int32 {
// 	return b.presum[i>>6] + int32(bits.OnesCount64(b.mask[i>>6]&((1<<(i&63))-1)))
// }

// func (b *boolPresum) On(i int32) bool {
// 	return (b.mask[i>>6]>>(i&63))&1 == 1
// }

// type boolBit struct {
// 	n      int32
// 	mask   []uint64
// 	presum *bITArray32
// }

// func newBoolBit(n int32, f func(int32) bool) *boolBit {
// 	mask := make([]uint64, n>>6+1)
// 	for i := int32(0); i < n; i++ {
// 		if f(i) {
// 			mask[i>>6] |= 1 << (i & 63)
// 		}
// 	}
// 	presum := newBitArrayFrom(int32(len(mask)), func(i int32) int32 {
// 		return int32(bits.OnesCount64(mask[i]))
// 	})
// 	return &boolBit{n: n, mask: mask, presum: presum}
// }

// func (b *boolBit) Popcount(i int32) int32 {
// 	return b.presum.QueryPrefix(i>>6) + int32(bits.OnesCount64(b.mask[i>>6]&((1<<(i&63))-1)))
// }

// func (b *boolBit) PopcountRange(le, ri int32) int32 {
// 	return b.Popcount(ri) - b.Popcount(le)
// }

// func (b *boolBit) On(i int32) bool {
// 	return (b.mask[i>>6]>>(i&63))&1 == 1
// }

// func (b *boolBit) Set(i int32, newNum bool) {
// 	if b.On(i) == newNum {
// 		return
// 	}
// 	b.mask[i>>6] ^= 1 << (i & 63)
// 	if newNum {
// 		b.presum.Add(i>>6, 1)
// 	} else {
// 		b.presum.Add(i>>6, -1)
// 	}
// }

// func split(tl, tr int32) int32 {
// 	pw2 := int32(1 << (bits.Len32(uint32(tr-tl)) - 1))
// 	return min32(tl+pw2, tr-pw2/2)
// }

// // !Point Add Range Sum, 0-based.
// type bITArray32 struct {
// 	n    int32
// 	data []int32
// }

// func newBitArray(n int32) *bITArray32 {
// 	res := &bITArray32{n: n, data: make([]int32, n)}
// 	return res
// }

// func newBitArrayFrom(n int32, f func(i int32) int32) *bITArray32 {
// 	data := make([]int32, n)
// 	for i := int32(0); i < n; i++ {
// 		data[i] = f(i)
// 	}
// 	for i := int32(1); i <= n; i++ {
// 		j := i + (i & -i)
// 		if j <= n {
// 			data[j-1] += data[i-1]
// 		}
// 	}
// 	return &bITArray32{n: n, data: data}
// }

// func (b *bITArray32) Add(index int32, v int32) {
// 	for index++; index <= b.n; index += index & -index {
// 		b.data[index-1] += v
// 	}
// }

// // [0, end).
// func (b *bITArray32) QueryPrefix(end int32) int32 {
// 	if end > b.n {
// 		end = b.n
// 	}
// 	res := int32(0)
// 	for ; end > 0; end &= end - 1 {
// 		res += b.data[end-1]
// 	}
// 	return res
// }

// // [start, end).
// func (b *bITArray32) QueryRange(start, end int32) int32 {
// 	if start < 0 {
// 		start = 0
// 	}
// 	if end > b.n {
// 		end = b.n
// 	}
// 	if start >= end {
// 		return 0
// 	}
// 	if start == 0 {
// 		return b.QueryPrefix(end)
// 	}
// 	pos, neg := int32(0), int32(0)
// 	for end > start {
// 		pos += b.data[end-1]
// 		end &= end - 1
// 	}
// 	for start > end {
// 		neg += b.data[start-1]
// 		start &= start - 1
// 	}
// 	return pos - neg
// }

// func min(a, b int) int {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func max(a, b int) int {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// func min32(a, b int32) int32 {
// 	if a < b {
// 		return a
// 	}
// 	return b
// }

// func max32(a, b int32) int32 {
// 	if a > b {
// 		return a
// 	}
// 	return b
// }

// // 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// // origin[newNums[i]] == nums[i]
// func Discretize(nums []int) (newNums []int32, origin []int) {
// 	newNums = make([]int32, len(nums))
// 	origin = make([]int, 0, len(newNums))
// 	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
// 	for _, i := range order {
// 		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
// 			origin = append(origin, nums[i])
// 		}
// 		newNums[i] = int32(len(origin) - 1)
// 	}
// 	origin = origin[:len(origin):len(origin)]
// 	return
// }

// func argSort(n int32, less func(i, j int32) bool) []int32 {
// 	order := make([]int32, n)
// 	for i := range order {
// 		order[i] = int32(i)
// 	}
// 	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
// 	return order
// }
// func stablePartition[T any](arr []T, start, end int32, predicate func(T) bool) int32 {
// 	ptr := start
// 	buffer := make([]T, end-start)
// 	bufferPtr := int32(0)
// 	for i := start; i < end; i++ {
// 		if predicate(arr[i]) {
// 			arr[ptr] = arr[i]
// 			ptr++
// 		} else {
// 			buffer[bufferPtr] = arr[i]
// 			bufferPtr++
// 		}
// 	}
// 	for i := int32(0); i < bufferPtr; i++ {
// 		arr[ptr] = buffer[i]
// 	}
// 	return ptr
// }
