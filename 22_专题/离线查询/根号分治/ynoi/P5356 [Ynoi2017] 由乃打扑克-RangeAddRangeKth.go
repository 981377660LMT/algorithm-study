// P5356 [Ynoi2017] 由乃打扑克-RangeAddRangeKth
// https://www.luogu.com.cn/problem/P5356
// https://blog.csdn.net/qq_42101694/article/details/109823342
// RangeAddRangeKth
//
// 区间加区间第 k 小问题。
// 首先是区间第k小，对于这个，我们可以二分答案这个值，
// 然后判断这个值在这个区间的排名，如果小于k，那就变大，如果大于k，那就变小。
// 如何判断排名呢？对于零散块，直接暴力统计。
// 对于整块，我们再做一次二分，而这就需要保证这个块一定要有序。
// 所以我们再建一个数组，这个数组元素和原数组一样，但是在每一个块中都排好序了，二分的时候用这个数组判断即可。
// 然后是区间修改，这个很简单，对于零散块，暴力修改，并对b数组更新+排序，对于整块，直接打lazytag即可。
// 我们在二分答案的时候，由于这题值域特小，所以我们的l和r不用设置inf，设置为这个区间最小，最大值即可。
// 由于b数组排好序的特性，这个最小最大值可以很快地求出。
//
// 优化：
// !1.将角块修改的地方不用直接排序，可以用归并排序讲一些无序和有序的数列段和起来，就将修改的复杂度优化到sqrt(n)
// 通过分块，统计懒标记，对块内归并来维护块内有序性(MergeTrick)，然后根据分块进行统计答案即可。
// !2.二分边界取区间最大和最小值(缩减二分长度)
// !3.查询时，如果当前区间的最大值小于 k那么就可以直接加上区间大小，如果当前区间最小值大于 k那么就不用再去计算它了。
// !4.玄学的块长:当块的大小为b时，修改是O(n/b+b)的，查询是O((n/b)*logb*logA)的。
// !取b=sqrt(n)*logA，单次复杂度就是O(sqrt(n)*logA)了。

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
	block := UseBlock(len(nums), int(math.Sqrt(float64(len(nums)))+1))

	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	blockLazy := make([]int, blockCount)
	blockMergeTrick := make([]*MergeTrick, blockCount)
	for bid := 0; bid < blockCount; bid++ {
		start, end := blockStart[bid], blockEnd[bid]
		blockMergeTrick[bid] = NewMergeTrick(nums[start:end], false)
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

			// 二分答案mid，区间里<=mid的数不超过k个.小段暴力统计，大段二分.
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

type SortedItem = struct{ value, index int }
type MergeTrick struct {
	_nums        []int
	_sortedItems []*SortedItem
	_sortedNums  []int
	_dirty       bool
}

// O(n)区间加, O(n)整体排序.
//
//	shouldCopy: 是否复制nums.
func NewMergeTrick(nums []int, shouldCopy bool) *MergeTrick {
	if shouldCopy {
		nums = append(nums[:0:0], nums...)
	}
	sortedItems := make([]*SortedItem, len(nums))
	for i := range nums {
		sortedItems[i] = &SortedItem{value: nums[i], index: i}
	}
	sort.Slice(sortedItems, func(i, j int) bool { return sortedItems[i].value < sortedItems[j].value })
	sortedNums := make([]int, len(nums))
	for i := range sortedItems {
		sortedNums[i] = sortedItems[i].value
	}
	return &MergeTrick{
		_nums:        nums,
		_sortedItems: sortedItems,
		_sortedNums:  sortedNums,
	}
}

func (mt *MergeTrick) Add(start, end, delta int) {
	mt._dirty = true
	n := len(mt._nums)
	modified := make([]*SortedItem, end-start)
	unmodified := make([]*SortedItem, n-(end-start))
	for i, ptr1, ptr2 := 0, 0, 0; i < n; i++ {
		item := mt._sortedItems[i]
		if index := item.index; index >= start && index < end {
			item.value += delta
			modified[ptr1] = item
			ptr1++
			mt._nums[index] += delta
		} else {
			unmodified[ptr2] = item
			ptr2++
		}
	}

	// 归并
	i1, i2, k := 0, 0, 0
	for i1 < len(modified) && i2 < len(unmodified) {
		if modified[i1].value < unmodified[i2].value {
			mt._sortedItems[k] = modified[i1]
			i1++
		} else {
			mt._sortedItems[k] = unmodified[i2]
			i2++
		}
		k++
	}

	for i1 < len(modified) {
		mt._sortedItems[k] = modified[i1]
		i1++
		k++
	}

	for i2 < len(unmodified) {
		mt._sortedItems[k] = unmodified[i2]
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
		res[i] = mt._sortedItems[i].value
	}
	mt._sortedNums = res
	return res
}
