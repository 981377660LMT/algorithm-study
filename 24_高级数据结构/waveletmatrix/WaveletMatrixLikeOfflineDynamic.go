package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
)

func main() {
	// test
	n := rand.Intn(5000) + 1
	q := 5000
	allNums := make([]int, n)
	for i := range allNums {
		allNums[i] = rand.Intn(1e9) - 8e5
	}
	blockSize := int(math.Sqrt(float64(n))+1) * 4
	w := NewWaveletMatrixLikeOfflineDynamic(allNums, allNums, blockSize)
	m := NewMocker(allNums)
	for i := 0; i < q; i++ {
		start := rand.Intn(n)
		end := start + rand.Intn(n-start) + 1
		k := rand.Intn(end - start)
		target := allNums[rand.Intn(n)]
		lower := allNums[rand.Intn(n)]
		upper := allNums[rand.Intn(n)]

		resGet := w.Get(start)
		resGetMocker := m.Get(start)
		if resGet != resGetMocker {
			panic("Get")
		}

		resOk, _ := w.Kth(start, end, k)
		res, _ := m.Kth(start, end, k)
		if resOk != res {
			panic("Kth")
		}

		resLower, _ := w.Lower(start, end, target)
		resLowerMocker, _ := m.Lower(start, end, target)
		if resLower != resLowerMocker {
			panic("Lower")
		}

		resHigher, _ := w.Higher(start, end, target)
		resHigherMocker, _ := m.Higher(start, end, target)
		if resHigher != resHigherMocker {
			panic("Higher")
		}

		resFloor, _ := w.Floor(start, end, target)
		resFloorMocker, _ := m.Floor(start, end, target)
		if resFloor != resFloorMocker {
			panic("Floor")
		}

		resCeiling, _ := w.Ceiling(start, end, target)
		resCeilingMocker, _ := m.Ceiling(start, end, target)
		if resCeiling != resCeilingMocker {
			panic("Ceiling")
		}

		resCount := w.Count(start, end, target)
		resCountMocker := m.Count(start, end, target)
		if resCount != resCountMocker {
			panic("Count")
		}

		resCountRange := w.CountRange(start, end, lower, upper)
		resCountRangeMocker := m.CountRange(start, end, lower, upper)
		if resCountRange != resCountRangeMocker {
			panic("CountRange")
		}

		resCountFloor := w.CountFloor(start, end, target)
		resCountFloorMocker := m.CountFloor(start, end, target)
		if resCountFloor != resCountFloorMocker {
			panic("CountFloor")
		}

		resCountCeiling := w.CountCeiling(start, end, target)
		resCountCeilingMocker := m.CountCeiling(start, end, target)
		if resCountCeiling != resCountCeilingMocker {
			panic("CountCeiling")
		}

		resCountLower := w.CountLower(start, end, target)
		resCountLowerMocker := m.CountLower(start, end, target)
		if resCountLower != resCountLowerMocker {
			panic("CountLower")
		}

		resCountHigher := w.CountHigher(start, end, target)
		resCountHigherMocker := m.CountHigher(start, end, target)
		if resCountHigher != resCountHigherMocker {
			panic("CountHigher")
		}

		w.Set(start, target)
		m.Set(start, target)
	}

	fmt.Println("pass")
}

type WaveletMatrixLikeOfflineDynamic struct {
	nums      []int32
	valueToId map[int]int32
	idToValue []int

	belong1     []int32
	blockStart1 []int32
	blockEnd1   []int32
	blockCount1 int32
	belong2     []int32
	blockStart2 []int32
	blockEnd2   []int32
	blockCount2 int32

	preSum1          [][]int32 // preSum1[i][j] 表示前i个块中有多少个数属于第j个值域块.
	preSum2          [][]int32 // preSum2[i][j] 表示前i个块中有多少个数等于j.
	fragmentCounter1 []int32   // fragmentCounter1[j] 临时保存散块内出现在第j个值域块中的数的个数.
	fragmentCounter2 []int32   // fragmentCounter2[j] 临时保存散块内等于j的数的个数.
}

// 序列分块+值域分块模拟`wavelet matrix`.支持单点修改.
// 单次操作时间复杂度`O(sqrt(n))`.
//
//	initNums: 初始化的序列.
//	allowedNums: 允许出现的数, 包含修改和查询的数.
//	blockSize: 序列分块的大小.为了减少空间复杂度`O(n*块数)`，可以把块的大小增大一些以减少内存占用.
//	blockSize := int(math.Sqrt(float64(len(initNums)))+1) * 4
func NewWaveletMatrixLikeOfflineDynamic(initNums []int, allowedNums []int, blockSize int) *WaveletMatrixLikeOfflineDynamic {
	set := make(map[int]struct{})
	for _, v := range allowedNums {
		set[v] = struct{}{}
	}
	sorted := make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	valueToId := make(map[int]int32, len(sorted))
	idToValue := make([]int, len(sorted))
	for i, v := range sorted {
		valueToId[v] = int32(i)
		idToValue[i] = v
	}

	newNums := make([]int32, len(initNums))
	for i, v := range initNums {
		newNums[i] = valueToId[v]
	}

	block1 := createBlockInt32(int32(len(newNums)), int32(blockSize))
	belong1, blockStart1, blockEnd1, blockCount1 := block1.belong, block1.blockStart, block1.blockEnd, block1.blockCount
	block2 := createBlockInt32(int32(len(sorted)), int32(math.Sqrt(float64(len(sorted)))+1))
	belong2, blockStart2, blockEnd2, blockCount2 := block2.belong, block2.blockStart, block2.blockEnd, block2.blockCount
	preSum1 := make([][]int32, blockCount1+1)
	for i := range preSum1 {
		preSum1[i] = make([]int32, blockCount2)
	}
	preSum2 := make([][]int32, blockCount1+1)
	for i := range preSum2 {
		preSum2[i] = make([]int32, len(sorted))
	}
	fragmentCounter1 := make([]int32, blockCount2)
	fragmentCounter2 := make([]int32, len(sorted))

	for bid := int32(0); bid < blockCount1; bid++ {
		copy(preSum1[bid+1], preSum1[bid])
		copy(preSum2[bid+1], preSum2[bid])
		s1, s2 := preSum1[bid+1], preSum2[bid+1]
		for i := blockStart1[bid]; i < blockEnd1[bid]; i++ {
			num := newNums[i]
			vid := belong2[num]
			s1[vid]++
			s2[num]++
		}
	}

	return &WaveletMatrixLikeOfflineDynamic{
		nums:             newNums,
		valueToId:        valueToId,
		idToValue:        idToValue,
		belong1:          belong1,
		blockStart1:      blockStart1,
		blockEnd1:        blockEnd1,
		blockCount1:      blockCount1,
		belong2:          belong2,
		blockStart2:      blockStart2,
		blockEnd2:        blockEnd2,
		blockCount2:      blockCount2,
		preSum1:          preSum1,
		preSum2:          preSum2,
		fragmentCounter1: fragmentCounter1,
		fragmentCounter2: fragmentCounter2,
	}
}

func (w *WaveletMatrixLikeOfflineDynamic) Get(index int) int {
	if index < 0 || index >= len(w.nums) {
		panic(fmt.Sprintf("index %d is out of range [0, %d)", index, len(w.nums)))
	}
	return w.idToValue[w.nums[index]]
}

func (w *WaveletMatrixLikeOfflineDynamic) Set(index, value int) {
	if index < 0 || index >= len(w.nums) {
		panic(fmt.Sprintf("index %d is out of range [0, %d)", index, len(w.nums)))
	}
	target := w._getId(value)
	if w.nums[index] == target {
		return
	}
	preValue, curValue := w.nums[index], target
	preVid, curVid := w.belong2[preValue], w.belong2[curValue]
	for bid := w.belong1[index]; bid < w.blockCount1; bid++ {
		s1, s2 := w.preSum1[bid+1], w.preSum2[bid+1]
		s1[preVid]--
		s1[curVid]++
		s2[preValue]--
		s2[curValue]++
	}
	w.nums[index] = target
}

// 查询区间`[start, end)`中第`k(k>=0)`小的数.
func (w *WaveletMatrixLikeOfflineDynamic) Kth(start, end, k int) (res int, ok bool) {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end || k >= end-start {
		return
	}

	startInt32, endInt32, kInt32 := int32(start), int32(end), int32(k)
	bid1, bid2 := w.belong1[start], w.belong1[end-1]
	if bid1 == bid2 {
		w._addFragment(startInt32, endInt32)
		todo := kInt32 + 1
		for vid := int32(0); vid < w.blockCount2; vid++ {
			if tmp := todo - w.fragmentCounter1[vid]; tmp > 0 {
				todo = tmp
				continue
			}
			// 答案在这个值域中.
			for j := w.blockStart2[vid]; j < w.blockEnd2[vid]; j++ {
				todo -= w.fragmentCounter2[j]
				if todo <= 0 {
					w._removeFragment(startInt32, endInt32)
					return w.idToValue[j], true
				}
			}
		}
	} else {
		w._addFragment(startInt32, w.blockEnd1[bid1])
		w._addFragment(w.blockStart1[bid2], endInt32)
		todo := kInt32 + 1
		for vid := int32(0); vid < w.blockCount2; vid++ {
			curCount := w.fragmentCounter1[vid] + w.preSum1[bid2][vid] - w.preSum1[bid1+1][vid]
			if tmp := todo - curCount; tmp > 0 {
				todo = tmp
				continue
			}
			for j := w.blockStart2[vid]; j < w.blockEnd2[vid]; j++ {
				curCount := w.fragmentCounter2[j] + w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
				todo -= curCount
				if todo <= 0 {
					w._removeFragment(startInt32, w.blockEnd1[bid1])
					w._removeFragment(w.blockStart1[bid2], endInt32)
					return w.idToValue[j], true
				}
			}
		}
	}

	panic("unreachable")
}

func (w *WaveletMatrixLikeOfflineDynamic) Floor(start, end, value int) (res int, ok bool) {
	return w._prev(start, end, value, true)
}

func (w *WaveletMatrixLikeOfflineDynamic) Ceiling(start, end, value int) (res int, ok bool) {
	return w._next(start, end, value, true)
}

func (w *WaveletMatrixLikeOfflineDynamic) Lower(start, end, value int) (res int, ok bool) {
	return w._prev(start, end, value, false)
}

func (w *WaveletMatrixLikeOfflineDynamic) Higher(start, end, value int) (res int, ok bool) {
	return w._next(start, end, value, false)
}

func (w *WaveletMatrixLikeOfflineDynamic) Count(start, end, value int) int {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return 0
	}
	startInt32, endInt32 := int32(start), int32(end)
	target := w._getId(value)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	if bid1 == bid2 {
		res := 0
		for i := startInt32; i < endInt32; i++ {
			if w.nums[i] == target {
				res++
			}
		}
		return res
	}
	res := int32(0)
	for i := startInt32; i < w.blockEnd1[bid1]; i++ {
		if w.nums[i] == target {
			res++
		}
	}
	for i := w.blockStart1[bid2]; i < endInt32; i++ {
		if w.nums[i] == target {
			res++
		}
	}
	res += w.preSum2[bid2][target] - w.preSum2[bid1+1][target]
	return int(res)
}

func (w *WaveletMatrixLikeOfflineDynamic) CountRange(start, end int, lower, upper int) int {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return 0
	}
	if lower >= upper {
		return 0
	}
	startInt32, endInt32 := int32(start), int32(end)
	lowerInt32, upperInt32 := w._getId(lower), w._getId(upper)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	if bid1 == bid2 {
		res := 0
		for i := startInt32; i < endInt32; i++ {
			if w.nums[i] >= lowerInt32 && w.nums[i] < upperInt32 {
				res++
			}
		}
		return res
	}
	res := int32(0)
	for i := startInt32; i < w.blockEnd1[bid1]; i++ {
		if w.nums[i] >= lowerInt32 && w.nums[i] < upperInt32 {
			res++
		}
	}
	for i := w.blockStart1[bid2]; i < endInt32; i++ {
		if w.nums[i] >= lowerInt32 && w.nums[i] < upperInt32 {
			res++
		}
	}
	vid1, vid2 := w.belong2[lowerInt32], w.belong2[upperInt32]
	if vid1 == vid2 {
		for j := lowerInt32; j < upperInt32; j++ {
			res += w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
		}
	} else {
		for j := lowerInt32; j < w.blockEnd2[vid1]; j++ {
			res += w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
		}
		for j := vid1 + 1; j < vid2; j++ {
			res += w.preSum1[bid2][j] - w.preSum1[bid1+1][j]
		}
		for j := w.blockStart2[vid2]; j < upperInt32; j++ {
			res += w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
		}
	}
	return int(res)
}

func (w *WaveletMatrixLikeOfflineDynamic) CountFloor(start, end, value int) int {
	return w._countPrev(start, end, value, true)
}

func (w *WaveletMatrixLikeOfflineDynamic) CountCeiling(start, end, value int) int {
	return w._countNext(start, end, value, true)
}

func (w *WaveletMatrixLikeOfflineDynamic) CountLower(start, end, value int) int {
	return w._countPrev(start, end, value, false)
}

func (w *WaveletMatrixLikeOfflineDynamic) CountHigher(start, end, value int) int {
	return w._countNext(start, end, value, false)
}

func (w *WaveletMatrixLikeOfflineDynamic) _addFragment(start, end int32) {
	for i := start; i < end; i++ {
		num := w.nums[i]
		vid := w.belong2[num]
		w.fragmentCounter1[vid] += 1
		w.fragmentCounter2[num] += 1
	}
}

func (w *WaveletMatrixLikeOfflineDynamic) _removeFragment(start, end int32) {
	for i := start; i < end; i++ {
		num := w.nums[i]
		vid := w.belong2[num]
		w.fragmentCounter1[vid] -= 1
		w.fragmentCounter2[num] -= 1
	}
}

func (w *WaveletMatrixLikeOfflineDynamic) _getId(rawValue int) int32 {
	if id, ok := w.valueToId[rawValue]; ok {
		return id
	}
	panic(fmt.Sprintf("value %d is not allowed", rawValue))
}

func (w *WaveletMatrixLikeOfflineDynamic) _prev(start, end int, target int, inclusive bool) (res int, ok bool) {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return
	}
	targetInt32 := w._getId(target)
	if !inclusive && targetInt32 == 0 {
		return
	}
	startInt32, endInt32 := int32(start), int32(end)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	vid := w.belong2[targetInt32]
	if bid1 == bid2 {
		w._addFragment(startInt32, endInt32)
		var cand int32
		if inclusive {
			cand = targetInt32
		} else {
			cand = targetInt32 - 1
		}
		vStart := w.blockStart2[vid]
		for ; cand >= vStart && w.fragmentCounter2[cand] == 0; cand-- {
		}
		if cand >= vStart {
			w._removeFragment(startInt32, endInt32)
			return w.idToValue[cand], true
		}
		candVid := vid - 1
		for ; candVid >= 0 && w.fragmentCounter1[candVid] == 0; candVid-- {
		}
		if candVid == -1 {
			w._removeFragment(startInt32, endInt32)
			return
		}
		for j := w.blockEnd2[candVid] - 1; j >= w.blockStart2[candVid]; j-- {
			if w.fragmentCounter2[j] > 0 {
				w._removeFragment(startInt32, endInt32)
				return w.idToValue[j], true
			}
		}
	} else {
		w._addFragment(startInt32, w.blockEnd1[bid1])
		w._addFragment(w.blockStart1[bid2], endInt32)
		var cand int32
		if inclusive {
			cand = targetInt32
		} else {
			cand = targetInt32 - 1
		}
		vStart := w.blockStart2[vid]
		for ; cand >= vStart && w.fragmentCounter2[cand]+w.preSum2[bid2][cand]-w.preSum2[bid1+1][cand] == 0; cand-- {
		}
		if cand >= vStart {
			w._removeFragment(startInt32, w.blockEnd1[bid1])
			w._removeFragment(w.blockStart1[bid2], endInt32)
			return w.idToValue[cand], true
		}
		candVid := vid - 1
		for ; candVid >= 0 && w.fragmentCounter1[candVid]+w.preSum1[bid2][candVid]-w.preSum1[bid1+1][candVid] == 0; candVid-- {
		}
		if candVid == -1 {
			w._removeFragment(startInt32, w.blockEnd1[bid1])
			w._removeFragment(w.blockStart1[bid2], endInt32)
			return
		}
		for j := w.blockEnd2[candVid] - 1; j >= w.blockStart2[candVid]; j-- {
			if w.fragmentCounter2[j]+w.preSum2[bid2][j]-w.preSum2[bid1+1][j] > 0 {
				w._removeFragment(startInt32, w.blockEnd1[bid1])
				w._removeFragment(w.blockStart1[bid2], endInt32)
				return w.idToValue[j], true
			}
		}
	}

	panic("unreachable")
}

func (w *WaveletMatrixLikeOfflineDynamic) _next(start, end int, target int, inclusive bool) (res int, ok bool) {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return
	}
	targetInt32 := w._getId(target)
	if !inclusive && targetInt32 == int32(len(w.idToValue))-1 {
		return
	}
	startInt32, endInt32 := int32(start), int32(end)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	vid := w.belong2[targetInt32]
	if bid1 == bid2 {
		w._addFragment(startInt32, endInt32)
		var cand int32
		if inclusive {
			cand = targetInt32
		} else {
			cand = targetInt32 + 1
		}
		vEnd := w.blockEnd2[vid]
		for ; cand < vEnd && w.fragmentCounter2[cand] == 0; cand++ {
		}
		if cand < vEnd {
			w._removeFragment(startInt32, endInt32)
			return w.idToValue[cand], true
		}
		candVid := vid + 1
		for ; candVid < w.blockCount2 && w.fragmentCounter1[candVid] == 0; candVid++ {
		}
		if candVid == w.blockCount2 {
			w._removeFragment(startInt32, endInt32)
			return
		}
		for j := w.blockStart2[candVid]; j < w.blockEnd2[candVid]; j++ {
			if w.fragmentCounter2[j] > 0 {
				w._removeFragment(startInt32, endInt32)
				return w.idToValue[j], true
			}
		}
	} else {
		w._addFragment(startInt32, w.blockEnd1[bid1])
		w._addFragment(w.blockStart1[bid2], endInt32)
		var cand int32
		if inclusive {
			cand = targetInt32
		} else {
			cand = targetInt32 + 1
		}
		vEnd := w.blockEnd2[vid]
		for ; cand < vEnd && w.fragmentCounter2[cand]+w.preSum2[bid2][cand]-w.preSum2[bid1+1][cand] == 0; cand++ {
		}
		if cand < vEnd {
			w._removeFragment(startInt32, w.blockEnd1[bid1])
			w._removeFragment(w.blockStart1[bid2], endInt32)
			return w.idToValue[cand], true
		}
		candVid := vid + 1
		for ; candVid < w.blockCount2 && w.fragmentCounter1[candVid]+w.preSum1[bid2][candVid]-w.preSum1[bid1+1][candVid] == 0; candVid++ {
		}
		if candVid == w.blockCount2 {
			w._removeFragment(startInt32, w.blockEnd1[bid1])
			w._removeFragment(w.blockStart1[bid2], endInt32)
			return
		}
		for j := w.blockStart2[candVid]; j < w.blockEnd2[candVid]; j++ {
			if w.fragmentCounter2[j]+w.preSum2[bid2][j]-w.preSum2[bid1+1][j] > 0 {
				w._removeFragment(startInt32, w.blockEnd1[bid1])
				w._removeFragment(w.blockStart1[bid2], endInt32)
				return w.idToValue[j], true
			}
		}
	}

	panic("unreachable")
}

func (w *WaveletMatrixLikeOfflineDynamic) _countPrev(start, end int, target int, inclusive bool) int {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return 0
	}
	targetInt32 := w._getId(target)
	if !inclusive {
		targetInt32--
	}
	if targetInt32 < 0 {
		return 0
	}
	startInt32, endInt32 := int32(start), int32(end)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	if bid1 == bid2 {
		res := 0
		for i := startInt32; i < endInt32; i++ {
			if w.nums[i] <= targetInt32 {
				res++
			}
		}
		return res
	}
	res := int32(0)
	for i := startInt32; i < w.blockEnd1[bid1]; i++ {
		if w.nums[i] <= targetInt32 {
			res++
		}
	}
	for i := w.blockStart1[bid2]; i < endInt32; i++ {
		if w.nums[i] <= targetInt32 {
			res++
		}
	}
	vid := w.belong2[targetInt32]
	for j := int32(0); j < vid; j++ {
		res += w.preSum1[bid2][j] - w.preSum1[bid1+1][j]
	}
	for j := w.blockStart2[vid]; j <= targetInt32; j++ {
		res += w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
	}
	return int(res)
}

func (w *WaveletMatrixLikeOfflineDynamic) _countNext(start, end int, target int, inclusive bool) int {
	if start < 0 {
		start = 0
	}
	if end > len(w.nums) {
		end = len(w.nums)
	}
	if start >= end {
		return 0
	}
	targetInt32 := w._getId(target)
	if !inclusive {
		targetInt32++
	}
	if targetInt32 >= int32(len(w.idToValue)) {
		return 0
	}
	startInt32, endInt32 := int32(start), int32(end)
	bid1, bid2 := w.belong1[startInt32], w.belong1[endInt32-1]
	if bid1 == bid2 {
		res := 0
		for i := startInt32; i < endInt32; i++ {
			if w.nums[i] >= targetInt32 {
				res++
			}
		}
		return res
	} else {
		res := int32(0)
		for i := startInt32; i < w.blockEnd1[bid1]; i++ {
			if w.nums[i] >= targetInt32 {
				res++
			}
		}
		for i := w.blockStart1[bid2]; i < endInt32; i++ {
			if w.nums[i] >= targetInt32 {
				res++
			}
		}
		vid := w.belong2[targetInt32]
		for j := vid + 1; j < w.blockCount2; j++ {
			res += w.preSum1[bid2][j] - w.preSum1[bid1+1][j]
		}
		for j := targetInt32; j < w.blockEnd2[vid]; j++ {
			res += w.preSum2[bid2][j] - w.preSum2[bid1+1][j]
		}
		return int(res)
	}
}

func createBlockInt32(n int32, blockSize int32) struct {
	belong     []int32 // 下标所属的块.
	blockStart []int32 // 每个块的起始下标(包含).
	blockEnd   []int32 // 每个块的结束下标(不包含).
	blockCount int32   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int32, blockCount)
	blockEnd := make([]int32, blockCount)
	belong := make([]int32, n)
	for i := int32(0); i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := int32(0); i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int32
		blockStart []int32
		blockEnd   []int32
		blockCount int32
	}{belong, blockStart, blockEnd, blockCount}
}

// test
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//
//

type Mocker struct {
	nums []int
}

func NewMocker(nums []int) *Mocker {
	return &Mocker{nums: append([]int(nil), nums...)}
}

func (m *Mocker) Get(index int) int {
	return m.nums[index]
}

func (m *Mocker) Set(index, value int) {
	m.nums[index] = value
}

func (m *Mocker) Kth(start, end, k int) (res int, ok bool) {
	tmp := append([]int(nil), m.nums[start:end]...)
	sort.Ints(tmp)
	return tmp[k], true
}

func (m *Mocker) Floor(start, end, value int) (res int, ok bool) {
	tmp := append([]int(nil), m.nums[start:end]...)
	sort.Ints(tmp)
	for i := len(tmp) - 1; i >= 0; i-- {
		if tmp[i] <= value {
			return tmp[i], true
		}
	}
	return
}

func (m *Mocker) Ceiling(start, end, value int) (res int, ok bool) {
	tmp := append([]int(nil), m.nums[start:end]...)
	sort.Ints(tmp)
	for _, v := range tmp {
		if v >= value {
			return v, true
		}
	}
	return
}

func (m *Mocker) Lower(start, end, value int) (res int, ok bool) {
	tmp := append([]int(nil), m.nums[start:end]...)
	sort.Ints(tmp)
	for i := len(tmp) - 1; i >= 0; i-- {
		if tmp[i] < value {
			return tmp[i], true
		}
	}
	return
}

func (m *Mocker) Higher(start, end, value int) (res int, ok bool) {
	tmp := append([]int(nil), m.nums[start:end]...)
	sort.Ints(tmp)
	for _, v := range tmp {
		if v > value {
			return v, true
		}
	}
	return
}

func (m *Mocker) Count(start, end, value int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] == value {
			res++
		}
	}
	return res
}

func (m *Mocker) CountRange(start, end int, lower, upper int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] >= lower && m.nums[i] < upper {
			res++
		}
	}
	return res
}

func (m *Mocker) CountFloor(start, end, value int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] <= value {
			res++
		}
	}
	return res
}

func (m *Mocker) CountCeiling(start, end, value int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] >= value {
			res++
		}
	}
	return res
}

func (m *Mocker) CountLower(start, end, value int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] < value {
			res++
		}
	}
	return res
}

func (m *Mocker) CountHigher(start, end, value int) int {
	res := 0
	for i := start; i < end; i++ {
		if m.nums[i] > value {
			res++
		}
	}
	return res
}
