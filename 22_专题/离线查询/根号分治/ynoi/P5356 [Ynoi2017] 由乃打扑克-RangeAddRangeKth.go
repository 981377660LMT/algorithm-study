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
// 通过分块，统计懒标记，对块内归并来维护块内有序性，然后根据分块进行统计答案即可。
// !2.二分边界取区间最大和最小值
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
	block := UseBlock(nums, int(math.Sqrt(float64(len(nums)))+1))
	belong, blockStart, blockEnd, blockCount := block.belong, block.blockStart, block.blockEnd, block.blockCount

	blockLazy := make([]int, blockCount)
	sortedNums := make([]int, len(nums))
	for bid := 0; bid < blockCount; bid++ {
		start, end := blockStart[bid], blockEnd[bid]
		for i := start; i < end; i++ {
			sortedNums[i] = nums[i]
		}
		sort.Ints(sortedNums[start:end])
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
				min_ = min(min_, blockLazy[i]+sortedNums[blockStart[i]])
				max_ = max(max_, blockLazy[i]+sortedNums[blockEnd[i]-1])
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
				for i := start; i < end; i++ {
					nums[i] += delta
				}
				for i := blockStart[bid1]; i < blockEnd[bid1]; i++ {
					sortedNums[i] = nums[i]
				}
				// TODO: 这里相当于多段有序，可以用merge优化
				sort.Ints(sortedNums[blockStart[bid1]:blockEnd[bid1]])
			} else {
				for i := start; i < blockEnd[bid1]; i++ {
					nums[i] += delta
				}
				for i := blockStart[bid1]; i < blockEnd[bid1]; i++ {
					sortedNums[i] = nums[i]
				}
				sort.Ints(sortedNums[blockStart[bid1]:blockEnd[bid1]])

				for i := bid1 + 1; i < bid2; i++ {
					blockLazy[i] += delta
				}

				for i := blockStart[bid2]; i < end; i++ {
					nums[i] += delta
				}
				for i := blockStart[bid2]; i < blockEnd[bid2]; i++ {
					sortedNums[i] = nums[i]
				}
				sort.Ints(sortedNums[blockStart[bid2]:blockEnd[bid2]])
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
					left, right := blockStart[bid], blockEnd[bid]-1
					ngt := bisectRight(sortedNums, mid-blockLazy[bid], left, right) - left
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
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

func bisectRight(nums []int, target int, left, right int) int {
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

func MergeTwoSortedArray(nums1, nums2 []int) []int {
	n1 := len(nums1)
	if n1 == 0 {
		return nums2
	}
	n2 := len(nums2)
	if n2 == 0 {
		return nums1
	}
	res := make([]int, n1+n2)
	i := 0
	j := 0
	k := 0
	for i < n1 && j < n2 {
		if nums1[i] < nums2[j] {
			res[k] = nums1[i]
			i++
		} else {
			res[k] = nums2[j]
			j++
		}
		k++
	}
	for i < n1 {
		res[k] = nums1[i]
		i++
		k++
	}
	for j < n2 {
		res[k] = nums2[j]
		j++
		k++
	}
	return res
}

func MergeThreeSortedArray(nums1, nums2, nums3 []int) []int {
	n1 := len(nums1)
	if n1 == 0 {
		return MergeTwoSortedArray(nums2, nums3)
	}
	n2 := len(nums2)
	if n2 == 0 {
		return MergeTwoSortedArray(nums1, nums3)
	}
	n3 := len(nums3)
	if n3 == 0 {
		return MergeTwoSortedArray(nums1, nums2)
	}
	res := make([]int, n1+n2+n3)
	i1 := 0
	i2 := 0
	i3 := 0
	k := 0
	for i1 < n1 && i2 < n2 && i3 < n3 {
		if nums1[i1] < nums2[i2] {
			if nums1[i1] < nums3[i3] {
				res[k] = nums1[i1]
				i1++
			} else {
				res[k] = nums3[i3]
				i3++
			}
		} else if nums2[i2] < nums3[i3] {
			res[k] = nums2[i2]
			i2++
		} else {
			res[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i1 < n1 && i2 < n2 {
		if nums1[i1] < nums2[i2] {
			res[k] = nums1[i1]
			i1++
		} else {
			res[k] = nums2[i2]
			i2++
		}
		k++
	}
	for i1 < n1 && i3 < n3 {
		if nums1[i1] < nums3[i3] {
			res[k] = nums1[i1]
			i1++
		} else {
			res[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i2 < n2 && i3 < n3 {
		if nums2[i2] < nums3[i3] {
			res[k] = nums2[i2]
			i2++
		} else {
			res[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i1 < n1 {
		res[k] = nums1[i1]
		i1++
		k++
	}
	for i2 < n2 {
		res[k] = nums2[i2]
		i2++
		k++
	}
	for i3 < n3 {
		res[k] = nums3[i3]
		i3++
		k++
	}
	return res
}
