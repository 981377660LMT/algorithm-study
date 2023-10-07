package main

import "math"

// 962. 最大宽度坡
// https://leetcode.cn/problems/maximum-width-ramp/
func maxWidthRamp(nums []int) int {
	res := 0
	Q := NewRightMostLeftMostQuery(nums)
	for i := 0; i < len(nums); i++ {
		rightMostCeiling := Q.RightMostCeiling(i)
		if rightMostCeiling != -1 {
			res = max(res, rightMostCeiling-i)
		}
	}
	return res
}

// 901. 股票价格跨度
// https://leetcode.cn/problems/online-stock-span
type StockSpanner struct {
	Q   *RightMostLeftMostQuery
	ptr int
}

func Constructor() StockSpanner {
	return StockSpanner{
		Q: NewRightMostLeftMostQuery(make([]int, 1e5+10)),
	}
}

func (this *StockSpanner) Next(price int) int {
	pos := this.ptr
	this.ptr++
	this.Q.Set(pos, price)
	leftNearestHigher := this.Q.LeftNearestHigher(pos)
	if leftNearestHigher == -1 {
		return pos + 1
	}
	return pos - leftNearestHigher
}

/**
 * Your StockSpanner object will be instantiated and called as such:
 * obj := Constructor();
 * param_1 := obj.Next(price);
 */

type RightMostLeftMostQuery struct {
	_nums       []int
	_belong     []int
	_blockStart []int
	_blockEnd   []int
	_blockCount int
	_blockMin   []int
	_blockMax   []int
	_blockLazy  []int
}

// 对每个下标，`O(sqrt)`查询 最右侧/最左侧 lower/floor/ceiling/higher 的元素.
func NewRightMostLeftMostQuery(arr []int) *RightMostLeftMostQuery {
	arr = append(arr[:0:0], arr...)
	n := len(arr)
	blockSize := (int(math.Sqrt(float64(n))) + 1)
	blockCount := 1 + (n / blockSize)
	blockStart := make([]int, blockCount)
	blockEnd := make([]int, blockCount)
	belong := make([]int, n)
	for i := 0; i < blockCount; i++ {
		blockStart[i] = i * blockSize
		blockEnd[i] = min((i+1)*blockSize, n)
	}
	for i := 0; i < n; i++ {
		belong[i] = (i / blockSize)
	}
	res := &RightMostLeftMostQuery{
		_nums:       arr,
		_belong:     belong,
		_blockStart: blockStart,
		_blockEnd:   blockEnd,
		_blockCount: blockCount,
		_blockMin:   make([]int, blockCount),
		_blockMax:   make([]int, blockCount),
		_blockLazy:  make([]int, blockCount),
	}
	for bid := 0; bid < blockCount; bid++ {
		res._rebuildBlock(bid)
	}
	return res
}

func (rm *RightMostLeftMostQuery) Set(index, value int) {
	if index < 0 || index >= len(rm._nums) {
		return
	}
	bid := rm._belong[index]
	lazy := rm._blockLazy[bid]
	pre := rm._nums[index] + lazy
	if pre == value {
		return
	}
	rm._nums[index] = value - lazy
	rm._rebuildBlock(bid)
}

func (rm *RightMostLeftMostQuery) AddRange(start, end, delta int) {
	if start < 0 {
		start = 0
	}
	if end > len(rm._nums) {
		end = len(rm._nums)
	}
	if start >= end {
		return
	}
	bid1 := rm._belong[start]
	bid2 := rm._belong[end-1]
	if bid1 == bid2 {
		for i := start; i < end; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid1)
	} else {
		for i := start; i < rm._blockEnd[bid1]; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid1)
		for bid := bid1 + 1; bid < bid2; bid++ {
			rm._blockLazy[bid] += delta
		}
		for i := rm._blockStart[bid2]; i < end; i++ {
			rm._nums[i] += delta
		}
		rm._rebuildBlock(bid2)
	}
}

// 查询`index`右侧最远的下标`j`，使得 `nums[j] < nums[index]`.
// 如果不存在，返回`-1`.
func (rm *RightMostLeftMostQuery) RightMostLower(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightMost(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostFloor(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightMost(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostCeiling(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightMost(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightMostHigher(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightMost(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostLower(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftMost(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostFloor(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftMost(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostCeiling(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftMost(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftMostHigher(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftMost(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestLower(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightNearest(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestFloor(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]

	return rm._queryRightNearest(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestCeiling(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightNearest(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) RightNearestHigher(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryRightNearest(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestLower(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftNearest(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] < cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] < cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestFloor(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftNearest(
		index,
		func(bid int) bool {
			return rm._blockMin[bid]+rm._blockLazy[bid] <= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] <= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestCeiling(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftNearest(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] >= cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] >= cur
		},
	)
}

func (rm *RightMostLeftMostQuery) LeftNearestHigher(index int) int {
	cur := rm._nums[index] + rm._blockLazy[rm._belong[index]]
	return rm._queryLeftNearest(
		index,
		func(bid int) bool {
			return rm._blockMax[bid]+rm._blockLazy[bid] > cur
		},
		func(eid, bid int) bool {
			return rm._nums[eid]+rm._blockLazy[bid] > cur
		},
	)
}

func (rm *RightMostLeftMostQuery) _queryRightMost(
	pos int,
	predicateBlock func(bid int) bool,
	predicateElement func(eid, bid int) bool,
) int {
	bid := rm._belong[pos]
	for i := rm._blockCount - 1; i > bid; i-- {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockEnd[i] - 1; j >= rm._blockStart[i]; j-- {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	for i := rm._blockEnd[bid] - 1; i > pos; i-- {
		if predicateElement(i, bid) {
			return i
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryLeftMost(
	pos int,
	predicateBlock func(bid int) bool,
	predicateElement func(eid, bid int) bool,
) int {
	bid := rm._belong[pos]
	for i := 0; i < bid; i++ {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockStart[i]; j < rm._blockEnd[i]; j++ {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	for i := rm._blockStart[bid]; i < pos; i++ {
		if predicateElement(i, bid) {
			return i
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryRightNearest(
	pos int,
	predicateBlock func(bid int) bool,
	predicateElement func(eid, bid int) bool,
) int {
	bid := rm._belong[pos]
	for i := pos + 1; i < rm._blockEnd[bid]; i++ {
		if predicateElement(i, bid) {
			return i
		}
	}
	for i := bid + 1; i < rm._blockCount; i++ {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockStart[i]; j < rm._blockEnd[i]; j++ {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _queryLeftNearest(
	pos int,
	predicateBlock func(bid int) bool,
	predicateElement func(eid, bid int) bool,
) int {
	bid := rm._belong[pos]
	for i := pos - 1; i >= rm._blockStart[bid]; i-- {
		if predicateElement(i, bid) {
			return i
		}
	}
	for i := bid - 1; i >= 0; i-- {
		if !predicateBlock(i) {
			continue
		}
		for j := rm._blockEnd[i] - 1; j >= rm._blockStart[i]; j-- {
			if predicateElement(j, i) {
				return j
			}
		}
	}
	return -1
}

func (rm *RightMostLeftMostQuery) _rebuildBlock(bid int) {
	rm._blockMin[bid] = math.MaxInt64
	rm._blockMax[bid] = math.MinInt64
	lazy := rm._blockLazy[bid]
	rm._blockLazy[bid] = 0
	for i := rm._blockStart[bid]; i < rm._blockEnd[bid]; i++ {
		rm._nums[i] += lazy
		rm._blockMin[bid] = min(rm._blockMin[bid], rm._nums[i])
		rm._blockMax[bid] = max(rm._blockMax[bid], rm._nums[i])
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
