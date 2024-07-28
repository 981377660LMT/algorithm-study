// 动态区间频率查询
//
// 1.区间加，查询区间某元素出现次数(RangeAddRangeFreq)
// 2.区间赋值，查询区间某元素出现次数(RangeAssignRangeFreq)
//
// !如果要支持insert/pop操作，使用"WaveletMatrixDynamic"
//
// TODO: Rust

package main

import (
	"fmt"
	"math"
	"math/rand"
	"sort"
	"strings"
	"time"
)

// 区间加，区间频率查询.
// 单次修改、查询时间复杂度`O(sqrt(n)logn)`，空间复杂度`O(n)`.
type RangeAddRangeFreq struct {
	nums        []int
	belong      []int
	blockStart  []int
	blockEnd    []int
	blockLazy   []int
	blockSorted [][]int
}

// ps := NewPointSetRangeFreq(arr, int(math.Sqrt(float64(len(arr)))+1))
func NewRangeAddRangeFreq(nums []int, blockSize int) *RangeAddRangeFreq {
	nums = append(nums[:0:0], nums...)
	res := &RangeAddRangeFreq{}
	block := UseBlock(len(nums), blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockLazy := make([]int, blockCount)
	blockSorted := make([][]int, blockCount)
	res.nums = nums
	res.belong = belong
	res.blockStart = blockStart
	res.blockEnd = blockEnd
	res.blockLazy = blockLazy
	res.blockSorted = blockSorted
	for bid := 0; bid < blockCount; bid++ {
		res.rebuild(bid)
	}
	return res
}

func (ra *RangeAddRangeFreq) Get(index int) int {
	return ra.nums[index] + ra.blockLazy[ra.belong[index]]
}

func (ra *RangeAddRangeFreq) Set(index, value int) {
	if index < 0 || index >= len(ra.nums) {
		return
	}
	bid := ra.belong[index]
	target := value - ra.blockLazy[bid]
	if target == ra.nums[index] {
		return
	}
	ra.nums[index] = target
	ra.rebuild(bid)
}

func (ra *RangeAddRangeFreq) Add(start, end int, delta int) {
	if start < 0 {
		start = 0
	}
	if end > len(ra.nums) {
		end = len(ra.nums)
	}
	if start >= end {
		return
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			ra.nums[i] += delta
		}
		ra.rebuild(bid1)
	} else {
		for i := start; i < ra.blockEnd[bid1]; i++ {
			ra.nums[i] += delta
		}
		ra.rebuild(bid1)
		for bid := bid1 + 1; bid < bid2; bid++ {
			ra.blockLazy[bid] += delta
		}
		for i := ra.blockStart[bid2]; i < end; i++ {
			ra.nums[i] += delta
		}
		ra.rebuild(bid2)
	}
}

// 统计 [start, end) 中等于 target 的元素个数.
func (ra *RangeAddRangeFreq) RangeFreq(start, end int, target int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.nums) {
		end = len(ra.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.nums[i] + ra.blockLazy[bid1]; cur == target {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid1]; cur == target {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += Count(ra.blockSorted[bid], target-ra.blockLazy[bid], 0, len(ra.blockSorted[bid])-1)
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid2]; cur == target {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中严格小于 higher 的元素个数.
func (ra *RangeAddRangeFreq) RangeFreqHigher(start, end int, higher int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.nums) {
		end = len(ra.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.nums[i] + ra.blockLazy[bid1]; cur < higher {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid1]; cur < higher {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += BisectLeft(ra.blockSorted[bid], higher-ra.blockLazy[bid], 0, len(ra.blockSorted[bid])-1)
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid2]; cur < higher {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中小于等于 ceiling 的元素个数.
func (ra *RangeAddRangeFreq) RangeFreqCeiling(start, end int, ceiling int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.nums) {
		end = len(ra.nums)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.nums[i] + ra.blockLazy[bid1]; cur <= ceiling {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid1]; cur <= ceiling {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		res += BisectRight(ra.blockSorted[bid], ceiling-ra.blockLazy[bid], 0, len(ra.blockSorted[bid])-1)
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.nums[i] + ra.blockLazy[bid2]; cur <= ceiling {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中大于等于 floor 的元素个数.
func (ra *RangeAddRangeFreq) RangeFreqFloor(start, end int, floor int) int {
	return (end - start) - ra.RangeFreqHigher(start, end, floor)
}

// 统计 [start, end) 中严格大于 lower 的元素个数.
func (ra *RangeAddRangeFreq) RangeFreqLower(start, end int, lower int) int {
	return (end - start) - ra.RangeFreqCeiling(start, end, lower)
}

func (ra *RangeAddRangeFreq) String() string {
	sb := []string{}
	sb = append(sb, "RangeAddRangeFreq{")
	for i := 0; i < len(ra.nums); i++ {
		sb = append(sb, fmt.Sprintf("%d", ra.Get(i)))
		if i != len(ra.nums)-1 {
			sb = append(sb, ",")
		}
	}
	sb = append(sb, "}")
	return strings.Join(sb, "")
}

func (ra *RangeAddRangeFreq) rebuild(bid int) {
	curSorted := make([]int, ra.blockEnd[bid]-ra.blockStart[bid])
	copy(curSorted, ra.nums[ra.blockStart[bid]:ra.blockEnd[bid]])
	sort.Ints(curSorted)
	ra.blockSorted[bid] = curSorted
}

const INF int = 1e18

// 区间赋值，区间频率查询.
// 单次修改、查询时间复杂度`O(sqrt(n)logn)`，空间复杂度`O(n)`.
type RangeAssignRangeFreq struct {
	color           []int
	belong          []int
	blockStart      []int
	blockEnd        []int
	blockSorted     [][]int
	blockColor      []int // 整块赋值后的颜色.INF表示未赋值过.
	updateTime      []int // 单点赋值时间戳.-1表示未赋值过.
	blockUpdateTime []int // 整块赋值时间戳.-1表示未赋值过.
	time            int   // 每次赋值操作时间戳+1
}

// ps := NewRangeAssignRangeFreq(arr, int(0.75*math.Sqrt(float64(len(arr)))+1))
func NewRangeAssignRangeFreq(nums []int, blockSize int) *RangeAssignRangeFreq {
	nums = append(nums[:0:0], nums...)
	res := &RangeAssignRangeFreq{}
	block := UseBlock(len(nums), blockSize)
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount
	blockSorted := make([][]int, blockCount)
	blockColor := make([]int, blockCount)
	updateTime := make([]int, len(nums))
	blockUpdateTime := make([]int, blockCount)
	for i := range blockColor {
		blockColor[i] = INF
		blockUpdateTime[i] = -1
	}
	for i := range updateTime {
		updateTime[i] = -1
	}
	res.color = nums
	res.belong = belong
	res.blockStart = blockStart
	res.blockEnd = blockEnd
	res.blockSorted = blockSorted
	res.blockColor = blockColor
	res.updateTime = updateTime
	res.blockUpdateTime = blockUpdateTime
	for bid := 0; bid < blockCount; bid++ {
		res.rebuild(bid)
	}
	return res
}

func (ra *RangeAssignRangeFreq) Get(index int) int {
	bid := ra.belong[index]
	if ra.blockUpdateTime[bid] > ra.updateTime[index] {
		return ra.blockColor[bid]
	}
	return ra.color[index]
}

func (ra *RangeAssignRangeFreq) Set(index, value int) {
	ra.Assign(index, index+1, value)
}

func (ra *RangeAssignRangeFreq) Assign(start, end int, value int) {
	if start < 0 {
		start = 0
	}
	if end > len(ra.color) {
		end = len(ra.color)
	}
	if start >= end {
		return
	}
	time := ra.time
	ra.time++
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			ra.color[i] = value
			ra.updateTime[i] = time
		}
		ra.rebuild(bid1)
	} else {
		for i := start; i < ra.blockEnd[bid1]; i++ {
			ra.color[i] = value
			ra.updateTime[i] = time
		}
		ra.rebuild(bid1)
		for bid := bid1 + 1; bid < bid2; bid++ {
			ra.blockColor[bid] = value
			ra.blockUpdateTime[bid] = time
		}
		for i := ra.blockStart[bid2]; i < end; i++ {
			ra.color[i] = value
			ra.updateTime[i] = time
		}
		ra.rebuild(bid2)
	}
}

// 统计 [start, end) 中等于 target 的元素个数.
func (ra *RangeAssignRangeFreq) RangeFreq(start, end int, target int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.color) {
		end = len(ra.color)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.Get(i); cur == target {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.Get(i); cur == target {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		if ra.blockUpdateTime[bid] != -1 {
			if cur := ra.blockColor[bid]; cur == target {
				res += ra.blockEnd[bid] - ra.blockStart[bid]
			}
		} else {
			res += Count(ra.blockSorted[bid], target, 0, len(ra.blockSorted[bid])-1)
		}
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.Get(i); cur == target {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中严格小于 higher 的元素个数.
func (ra *RangeAssignRangeFreq) RangeFreqHigher(start, end int, higher int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.color) {
		end = len(ra.color)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.Get(i); cur < higher {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.Get(i); cur < higher {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		if ra.blockUpdateTime[bid] != -1 {
			if cur := ra.blockColor[bid]; cur < higher {
				res += ra.blockEnd[bid] - ra.blockStart[bid]
			}
		} else {
			res += BisectLeft(ra.blockSorted[bid], higher, 0, len(ra.blockSorted[bid])-1)
		}
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.Get(i); cur < higher {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中小于等于 ceiling 的元素个数.
func (ra *RangeAssignRangeFreq) RangeFreqCeiling(start, end int, ceiling int) int {
	if start < 0 {
		start = 0
	}
	if end > len(ra.color) {
		end = len(ra.color)
	}
	if start >= end {
		return 0
	}
	bid1, bid2 := ra.belong[start], ra.belong[end-1]
	if bid1 == bid2 {
		res := 0
		for i := start; i < end; i++ {
			if cur := ra.Get(i); cur <= ceiling {
				res++
			}
		}
		return res
	}
	res := 0
	for i := start; i < ra.blockEnd[bid1]; i++ {
		if cur := ra.Get(i); cur <= ceiling {
			res++
		}
	}
	for bid := bid1 + 1; bid < bid2; bid++ {
		if ra.blockUpdateTime[bid] != -1 {
			if cur := ra.blockColor[bid]; cur <= ceiling {
				res += ra.blockEnd[bid] - ra.blockStart[bid]
			}
		} else {
			res += BisectRight(ra.blockSorted[bid], ceiling, 0, len(ra.blockSorted[bid])-1)
		}
	}
	for i := ra.blockStart[bid2]; i < end; i++ {
		if cur := ra.Get(i); cur <= ceiling {
			res++
		}
	}
	return res
}

// 统计 [start, end) 中大于等于 floor 的元素个数.
func (ra *RangeAssignRangeFreq) RangeFreqFloor(start, end int, floor int) int {
	return (end - start) - ra.RangeFreqHigher(start, end, floor)
}

// 统计 [start, end) 中严格大于 lower 的元素个数.
func (ra *RangeAssignRangeFreq) RangeFreqLower(start, end int, lower int) int {
	return (end - start) - ra.RangeFreqCeiling(start, end, lower)
}

func (ra *RangeAssignRangeFreq) String() string {
	sb := []string{}
	sb = append(sb, "RangeAssignRangeFreq{")
	for i := 0; i < len(ra.color); i++ {
		sb = append(sb, fmt.Sprintf("%d", ra.Get(i)))
		if i != len(ra.color)-1 {
			sb = append(sb, ",")
		}
	}
	sb = append(sb, "}")
	return strings.Join(sb, "")
}

func (ra *RangeAssignRangeFreq) rebuild(bid int) {
	for i := ra.blockStart[bid]; i < ra.blockEnd[bid]; i++ {
		ra.color[i] = ra.Get(i)
	}
	ra.blockColor[bid] = INF
	ra.blockUpdateTime[bid] = -1
	curSorted := make([]int, ra.blockEnd[bid]-ra.blockStart[bid])
	copy(curSorted, ra.color[ra.blockStart[bid]:ra.blockEnd[bid]])
	sort.Ints(curSorted)
	ra.blockSorted[bid] = curSorted
}

// blockSize := int(math.Sqrt(float64(len(nums))) + 1)
func UseBlock(n int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		tmp := (i + 1) * blockSize
		if tmp > n {
			tmp = n
		}
		blockEnd[i] = tmp
	}
	for i := 0; i < n; i++ {
		belong[i] = i / blockSize
	}

	return struct {
		belong     []int
		blockStart []int
		blockEnd   []int
		blockCount int
	}{belong, blockStart, blockEnd, blockCount}
}

func BisectLeft(nums []int, target int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] >= target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func BisectRight(nums []int, target int, left, right int) int {
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] > target {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}

func Count(nums []int, target int, left, right int) int {
	return BisectRight(nums, target, left, right) - BisectLeft(nums, target, left, right)
}

type _Mocker struct {
	_nums []int
}

func NewMocker(nums []int) *_Mocker {
	return &_Mocker{append(nums[:0:0], nums...)}
}

func (m *_Mocker) Get(index int) int {
	return m._nums[index]
}

func (m *_Mocker) Set(index, value int) {
	m._nums[index] = value
}

func (m *_Mocker) Add(start, end int, delta int) {
	for i := start; i < end; i++ {
		m._nums[i] += delta
	}
}

// assign

func (m *_Mocker) Assign(start, end int, value int) {
	for i := start; i < end; i++ {
		m._nums[i] = value
	}
}

func (m *_Mocker) RangeFreq(start, end int, target int) int {
	res := 0
	for i := start; i < end; i++ {
		if m._nums[i] == target {
			res++
		}
	}
	return res
}

func (m *_Mocker) RangeFreqFloor(start, end int, floor int) int {
	res := 0
	for i := start; i < end; i++ {
		if m._nums[i] >= floor {
			res++
		}
	}
	return res
}

func (m *_Mocker) RangeFreqLower(start, end int, lower int) int {
	res := 0
	for i := start; i < end; i++ {
		if m._nums[i] > lower {
			res++
		}
	}
	return res
}

func (m *_Mocker) RangeFreqCeiling(start, end int, ceiling int) int {
	res := 0
	for i := start; i < end; i++ {
		if m._nums[i] <= ceiling {
			res++
		}
	}
	return res

}

func (m *_Mocker) RangeFreqHigher(start, end int, higher int) int {
	res := 0
	for i := start; i < end; i++ {
		if m._nums[i] < higher {
			res++
		}
	}
	return res
}

func main() {

	testRangeAddRangeFreq := func() {
		N := int(2e4)
		arr := make([]int, N)
		for i := 0; i < N; i++ {
			arr[i] = i
		}
		real := NewRangeAddRangeFreq(arr, int(math.Sqrt(float64(len(arr)))+1))
		mock := NewMocker(arr)
		randint := func(a, b int) int {
			return rand.Intn(b-a+1) + a
		}
		for i := 0; i < 2e4; i++ {
			op := randint(0, 6)
			if op == 0 {
				index := randint(0, N-1)
				value := randint(-10, 10)
				real.Set(index, value)
				mock.Set(index, value)
			} else if op == 1 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				delta := randint(-5, 5)
				real.Add(start, end, delta)
				mock.Add(start, end, delta)
			} else if op == 2 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreq(start, end, floor)
				ans := mock.RangeFreq(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 3 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqFloor(start, end, floor)
				ans := mock.RangeFreqFloor(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 4 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqCeiling(start, end, floor)
				ans := mock.RangeFreqCeiling(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 5 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqHigher(start, end, floor)
				ans := mock.RangeFreqHigher(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 6 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqLower(start, end, floor)
				ans := mock.RangeFreqLower(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			}
		}
	}

	testRangeAssignRangeFreq := func() {
		N := int(2e4)
		arr := make([]int, N)
		for i := 0; i < N; i++ {
			arr[i] = i
		}
		real := NewRangeAssignRangeFreq(arr, int(math.Sqrt(float64(len(arr)))+1))
		mock := NewMocker(arr)
		randint := func(a, b int) int {
			return rand.Intn(b-a+1) + a
		}
		for i := 0; i < 2e4; i++ {
			op := randint(0, 6)
			if op == 0 {
				index := randint(0, N-1)
				value := randint(-10, 10)
				real.Set(index, value)
				mock.Set(index, value)
			} else if op == 1 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				delta := randint(-5, 5)
				real.Assign(start, end, delta)
				mock.Assign(start, end, delta)
			} else if op == 2 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreq(start, end, floor)
				ans := mock.RangeFreq(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 3 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqFloor(start, end, floor)
				ans := mock.RangeFreqFloor(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 4 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqCeiling(start, end, floor)
				ans := mock.RangeFreqCeiling(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 5 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqHigher(start, end, floor)
				ans := mock.RangeFreqHigher(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			} else if op == 6 {
				start := randint(0, N-1)
				end := randint(start, N-1)
				floor := randint(-10, 10)
				res := real.RangeFreqLower(start, end, floor)
				ans := mock.RangeFreqLower(start, end, floor)
				if res != ans {
					fmt.Println(start, end, floor)
					fmt.Println(res, ans)
					panic("")
				}
			}
		}

		fmt.Println("ok 3")
	}

	testTime2 := func() {
		N := int(1e5)
		nums := make([]int, N)
		for i := 0; i < N; i++ {
			nums[i] = i
		}
		ps := NewRangeAddRangeFreq(nums, int(math.Sqrt(float64(len(nums)))+1))
		time1 := time.Now()
		for i := range nums {
			ps.Add(0, N, i)
			ps.RangeFreq(0, len(nums), i)
		}
		time2 := time.Now()
		fmt.Println(time2.Sub(time1))
	}

	testTime3 := func() {
		N := int(1e5)
		nums := make([]int, N)
		for i := 0; i < N; i++ {
			nums[i] = i
		}
		ps := NewRangeAssignRangeFreq(nums, int(0.75*math.Sqrt(float64(len(nums)))+1))
		time1 := time.Now()
		for i := range nums {
			ps.Assign(i, N, i)
			ps.RangeFreq(0, len(nums), i)
		}
		time2 := time.Now()
		fmt.Println(time2.Sub(time1))
	}

	_ = testRangeAddRangeFreq
	_ = testRangeAssignRangeFreq

	testTime2()
	testTime3()

}
