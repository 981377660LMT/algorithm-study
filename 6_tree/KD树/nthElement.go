// !在 KD 树中不如 Sort 快, 谨慎使用

// cpp nthElement

// 排序规则采用默认的升序排序
// void nth_element (RandomAccessIterator first,
// 	RandomAccessIterator nth,
// 	RandomAccessIterator last);
// !排序规则为自定义的 comp 排序规则
// void nth_element (RandomAccessIterator first,
// 	RandomAccessIterator nth,
// 	RandomAccessIterator last,
// 	Compare comp);

package main

import (
	"math/rand"
	"time"
)

// 973. 最接近原点的 K 个点
// https://leetcode.cn/problems/k-closest-points-to-origin/
func kClosest(ps [][]int, k int) [][]int {
	NthElement(ps, k, func(i, j int) bool {
		return ps[i][0]*ps[i][0]+ps[i][1]*ps[i][1] < ps[j][0]*ps[j][0]+ps[j][1]*ps[j][1]
	})
	return ps[:k]
}

type N = []int

// !从数组 nums 中找到第 n 小的元素 a ，并移动到序列中第 n 个下标处。
// !处理后，所有位于 a 之前的元素都不大于 a ，所有位于 a 之后的元素都不小于 a。
//  1 <= nth <= len(nums)
func NthElement(nums []N, nth int, less func(i, j int) bool) {
	nth--
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	for l, r := 0, len(nums)-1; l < r; {
		v := nums[l]
		i, j := l, r+1
		for {
			for i++; i < r && less(i, l); i++ {
			}
			for j--; j > l && less(l, j); j-- {
			}
			if i >= j {
				break
			}
			nums[i], nums[j] = nums[j], nums[i]
		}
		nums[l], nums[j] = nums[j], v
		if j == nth {
			break
		} else if j < nth {
			l = j + 1
		} else {
			r = j - 1
		}
	}
}

// 数组第 k 小 (Quick Select)       kthElement/nthElement/QuickSelect
// 1 <= k <= len(a)
// 代码实现参考算法第四版 p.221
// 算法的平均比较次数为 ~2n+2kln(n/k)+2(n-k)ln(n/(n-k))
// https://en.wikipedia.org/wiki/Quickselect
// https://www.geeksforgeeks.org/quickselect-algorithm/
// 模板题 LC215 https://leetcode-cn.com/problems/kth-largest-element-in-an-array/
//       LC973 https://leetcode-cn.com/problems/k-closest-points-to-origin/submissions/
// 模板题 https://codeforces.com/contest/977/problem/C
// https://github.dev/EndlessCheng/codeforces-go/tree/master/copypasta
func NthElement2(nums []int, nth int) int {
	nth--
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(nums), func(i, j int) { nums[i], nums[j] = nums[j], nums[i] })
	for l, r := 0, len(nums)-1; l < r; { // first, last
		v := nums[l] // 切分元素
		i, j := l, r+1
		for {
			for i++; i < r && nums[i] < v; i++ { // less(i, l)
			}
			for j--; j > l && nums[j] > v; j-- { // less(l, j)
			}
			if i >= j {
				break
			}
			nums[i], nums[j] = nums[j], nums[i]
		}
		nums[l], nums[j] = nums[j], v
		if j == nth {
			break
		} else if j < nth {
			l = j + 1
		} else {
			r = j - 1
		}
	}
	return nums[nth] //  a[:k+1]  a[k:]
}

// 求第 k 小(1 <= k <= len(nums)
//  会改变原数组中元素顺序
func KthMin(nums []int, kth int) int {
	return NthElement2(nums, kth)
}

// 求第 k 大(1 <= k <= len(nums)
//  会改变原数组中元素顺序
func KthMax(nums []int, kth int) int {
	return NthElement2(nums, len(nums)+1-kth)
}
