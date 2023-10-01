package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"sort"
)

const INF int = 1e18

// 0 start end delta
// 1 start end k (k从0开始)
func RangeAddRangeKth(nums []int, operations [][4]int) []int {
	nums = append(nums[:0:0], nums...)
	block := UseBlock(nums, int(math.Sqrt(float64(len(nums)))+1))
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	blockLazy := make([]int, blockCount)
	blockMergeTrick := make([]*MergeTrick, blockCount)
	for bid := 0; bid < blockCount; bid++ {
		start, end := blockStart[bid], blockEnd[bid]
		blockMergeTrick[bid] = NewMergeTrick(nums[start:end])
	}

	getMinAndMax := func(start, end int, bid1, bid2 int) (min_, max_ int) {
		min_, max_ = INF, -INF

		if bid1 == bid2 {
			for i := start; i < end; i++ {
				cur := nums[i] + blockLazy[bid1]
				min_ = min(min_, cur)
				max_ = max(max_, cur)
			}
		} else {
			for i := start; i < blockEnd[bid1]; i++ {
				cur := nums[i] + blockLazy[bid1]
				min_ = min(min_, cur)
				max_ = max(max_, cur)
			}
			for i := bid1 + 1; i < bid2; i++ {
				sortedNums := blockMergeTrick[i].GetSortedNums()
				min_ = min(min_, blockLazy[i]+sortedNums[0])
				max_ = max(max_, blockLazy[i]+sortedNums[len(sortedNums)-1])
			}
			for i := blockStart[bid2]; i < end; i++ {
				cur := nums[i] + blockLazy[bid2]
				min_ = min(min_, cur)
				max_ = max(max_, cur)
			}
		}

		return
	}

	res := []int{}
	for _, op := range operations {
		kind := op[0]
		if kind == 0 {
			start, end, delta := op[1], op[2], op[3]
			bid1, bid2 := belong[start], belong[end-1]
			if bid1 == bid2 {
				blockMergeTrick[bid1].Add(start-blockStart[bid1], end-blockStart[bid1], delta)
			} else {
				blockMergeTrick[bid1].Add(start-blockStart[bid1], blockEnd[bid1]-blockStart[bid1], delta)
				for i := bid1 + 1; i < bid2; i++ {
					blockLazy[i] += delta
				}
				blockMergeTrick[bid2].Add(0, end-blockStart[bid2], delta)
			}
		} else {
			start, end, k := op[1], op[2], op[3]
			if k < 0 || k > end-start-1 {
				res = append(res, -1)
				continue
			}
			bid1, bid2 := belong[start], belong[end-1]

			// 二分答案mid，区间里<=mid的数不超过k个
			check := func(mid int) bool {
				res := 0
				if bid1 == bid2 {
					for i := start; i < end; i++ {
						cur := nums[i] + blockLazy[bid1]
						if cur <= mid {
							res++
							if res > k {
								return false
							}
						}
					}
					return res <= k
				}

				for bid := bid1 + 1; bid < bid2; bid++ {
					sortedNums := blockMergeTrick[bid].GetSortedNums()
					if sortedNums[0]+blockLazy[bid] > mid {
						continue
					}
					if sortedNums[len(sortedNums)-1]+blockLazy[bid] <= mid {
						res += len(sortedNums)
						if res > k {
							return false
						}
						continue
					}
					ngt := bisectRight(sortedNums, mid-blockLazy[bid])
					res += ngt
					if res > k {
						return false
					}
				}
				for i := start; i < blockEnd[bid1]; i++ {
					cur := nums[i] + blockLazy[bid1]
					if cur <= mid {
						res++
						if res > k {
							return false
						}
					}
				}
				for i := blockStart[bid2]; i < end; i++ {
					cur := nums[i] + blockLazy[bid2]
					if cur <= mid {
						res++
						if res > k {
							return false
						}
					}
				}
				return res <= k
			}

			left, right := getMinAndMax(start, end, bid1, bid2)
			if k == 0 {
				res = append(res, left)
				continue
			}
			if k == end-start-1 {
				res = append(res, right)
				continue
			}

			for left <= right {
				mid := (left + right) >> 1
				if check(mid) {
					left = mid + 1
				} else {
					right = mid - 1
				}
			}
			res = append(res, left)

		}
	}

	return res
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	nums := make([]int, n)
	for i := range nums {
		fmt.Fscan(in, &nums[i])
	}

	operations := make([][4]int, q)
	for i := range operations {
		var op int
		fmt.Fscan(in, &op)
		if op == 1 {
			var start, end, k int
			fmt.Fscan(in, &start, &end, &k)
			start--
			k--
			operations[i] = [4]int{1, start, end, k}
		} else {
			var start, end, delta int
			fmt.Fscan(in, &start, &end, &delta)
			start--
			operations[i] = [4]int{0, start, end, delta}
		}
	}

	res := RangeAddRangeKth(nums, operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// blockSize = int(math.Sqrt(float64(len(nums)))+1)
func UseBlock(nums []int, blockSize int) struct {
	belong     []int // 下标所属的块.
	blockStart []int // 每个块的起始下标(包含).
	blockEnd   []int // 每个块的结束下标(不包含).
	blockCount int   // 块的数量.
} {
	n := len(nums)

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

func max(a, b int) int {
	if a > b {
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

func bisectRight(nums []int, target int) int {
	left, right := 0, len(nums)-1
	for left <= right {
		mid := (left + right) >> 1
		if nums[mid] <= target {
			left = mid + 1
		} else {
			right = mid - 1
		}
	}
	return left
}

type MergeTrick struct {
	_nums        []int
	_originIndex []int
	_sortedNums  []int
	_dirty       bool
}

func NewMergeTrick(nums []int) *MergeTrick {
	nums = append(nums[:0:0], nums...)
	order := make([]int, len(nums))
	for i := range order {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool {
		return nums[order[i]] < nums[order[j]]
	})
	originIndex := make([]int, len(nums))
	for i := range originIndex {
		originIndex[order[i]] = i
	}
	sortedNums := append(nums[:0:0], nums...)
	sort.Ints(sortedNums)
	return &MergeTrick{_nums: nums, _originIndex: originIndex, _sortedNums: sortedNums}
}

func (mt *MergeTrick) Add(start, end, delta int) {
	mt._dirty = true
	n := len(mt._nums)
	modified := make([]int, end-start)
	unmodified := make([]int, n-(end-start))
	for i, ptr1, ptr2 := 0, 0, 0; i < n; i++ {
		index := mt._originIndex[i]
		if index >= start && index < end {
			modified[ptr1] = index
			mt._nums[index] += delta
			ptr1++
		} else {
			unmodified[ptr2] = index
			ptr2++
		}
	}

	// 归并
	i1, i2, k := 0, 0, 0
	for i1 < len(modified) && i2 < len(unmodified) {
		if mt._nums[modified[i1]] < mt._nums[unmodified[i2]] {
			mt._originIndex[k] = modified[i1]
			i1++
		} else {
			mt._originIndex[k] = unmodified[i2]
			i2++
		}
		k++
	}

	for i1 < len(modified) {
		mt._originIndex[k] = modified[i1]
		i1++
		k++
	}

	for i2 < len(unmodified) {
		mt._originIndex[k] = unmodified[i2]
		i2++
		k++
	}
}

// 返回原始数组.
func (mt *MergeTrick) GetNums() []int {
	return mt._nums
}

// 返回排序后的数组.
func (mt *MergeTrick) GetSortedNums() []int {
	if !mt._dirty {
		return mt._sortedNums
	}
	mt._dirty = false
	res := make([]int, len(mt._nums))
	for i := range res {
		res[i] = mt._nums[mt._originIndex[i]]
	}
	mt._sortedNums = res
	return res
}
