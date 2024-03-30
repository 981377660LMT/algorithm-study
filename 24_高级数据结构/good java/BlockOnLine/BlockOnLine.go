package main

type BlockOnLine struct {
	left, right []int32
	n           int32
}

func NewBlockOnLine(n int32) *BlockOnLine {
	b := &BlockOnLine{
		left:  make([]int32, n),
		right: make([]int32, n),
		n:     n,
	}
	b.Init()
	return b
}

func (b *BlockOnLine) Init() {
	n := b.n
	for i := int32(0); i < n; i++ {
		b.left[i] = n
		b.right[i] = -1
	}
}

// Add adds the element at the index position.
// If there are elements on the left or right, the original block is deleted and merged into a new block.
func (b *BlockOnLine) Add(index int32, onAddBlock, onRemoveBlock func(start, end int32)) bool {
	n := b.n
	if !(0 <= index && index < n) {
		return false
	}
	if !(b.left[index] > b.right[index]) {
		return false
	}

	from, to := index, index
	if index > 0 && b.left[index-1] <= b.right[index-1] {
		from = b.left[index-1]
		if onRemoveBlock != nil {
			onRemoveBlock(from, index)
		}
	}
	if index+1 < n && b.left[index+1] <= b.right[index+1] {
		to = b.right[index+1]
		if onRemoveBlock != nil {
			onRemoveBlock(index+1, to+1)
		}
	}
	b.left[from] = from
	b.right[from] = to
	b.left[to] = from
	b.right[to] = to
	if onAddBlock != nil {
		onAddBlock(from, to+1)
	}
	return true
}

// 2382. 删除操作后的最大子段和
// https://leetcode.cn/problems/maximum-segment-sum-after-removals/
func maximumSegmentSum(nums []int, removeQueries []int) []int64 {
	n := int32(len(nums))
	preSum := make([]int64, n+1)
	for i := int32(0); i < n; i++ {
		preSum[i+1] = preSum[i] + int64(nums[i])
	}

	res := make([]int64, len(removeQueries))
	curMax := int64(0)
	B := NewBlockOnLine(n)

	for i := len(removeQueries) - 1; i >= 0; i-- {
		res[i] = int64(curMax)
		B.Add(
			int32(removeQueries[i]),
			func(start, end int32) {
				sum := preSum[end] - preSum[start]
				curMax = max64(curMax, sum)
			},
			nil,
		)
	}
	return res
}

func max64(a, b int64) int64 {
	if a > b {
		return a
	}
	return b
}
