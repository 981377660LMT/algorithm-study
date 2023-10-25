package main

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

func MergeKSortedArray(arrays [][]int) []int {
	n := len(arrays)
	if n == 0 {
		return nil
	}
	if n == 1 {
		return arrays[0]
	}
	if n == 2 {
		return MergeTwoSortedArray(arrays[0], arrays[1])
	}

	var merge func(start, end int) []int
	merge = func(start, end int) []int {
		if start >= end {
			return nil
		}
		if end-start == 1 {
			return arrays[start]
		}
		mid := (start + end) >> 1
		return MergeTwoSortedArray(merge(start, mid), merge(mid, end))
	}

	return merge(0, n)
}

func MergeTwoSortedArrayTo(nums1, nums2 []int, target []int) {
	n1 := len(nums1)
	if n1 == 0 {
		copy(target, nums2)
		return
	}
	n2 := len(nums2)
	if n2 == 0 {
		copy(target, nums1)
		return
	}

	i := 0
	j := 0
	k := 0
	for i < n1 && j < n2 {
		if nums1[i] < nums2[j] {
			target[k] = nums1[i]
			i++
		} else {
			target[k] = nums2[j]
			j++
		}
		k++
	}
	for i < n1 {
		target[k] = nums1[i]
		i++
		k++
	}
	for j < n2 {
		target[k] = nums2[j]
		j++
		k++
	}

}

func MergeThreeSortedArrayTo(nums1, nums2, nums3 []int, target []int) {
	n1 := len(nums1)
	if n1 == 0 {
		MergeTwoSortedArrayTo(nums2, nums3, target)
		return
	}
	n2 := len(nums2)
	if n2 == 0 {
		MergeTwoSortedArrayTo(nums1, nums3, target)
		return
	}
	n3 := len(nums3)
	if n3 == 0 {
		MergeTwoSortedArrayTo(nums1, nums2, target)
		return
	}

	i1 := 0
	i2 := 0
	i3 := 0
	k := 0
	for i1 < n1 && i2 < n2 && i3 < n3 {
		if nums1[i1] < nums2[i2] {
			if nums1[i1] < nums3[i3] {
				target[k] = nums1[i1]
				i1++
			} else {
				target[k] = nums3[i3]
				i3++
			}
		} else if nums2[i2] < nums3[i3] {
			target[k] = nums2[i2]
			i2++
		} else {
			target[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i1 < n1 && i2 < n2 {
		if nums1[i1] < nums2[i2] {
			target[k] = nums1[i1]
			i1++
		} else {
			target[k] = nums2[i2]
			i2++
		}
		k++
	}
	for i1 < n1 && i3 < n3 {
		if nums1[i1] < nums3[i3] {
			target[k] = nums1[i1]
			i1++
		} else {
			target[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i2 < n2 && i3 < n3 {
		if nums2[i2] < nums3[i3] {
			target[k] = nums2[i2]
			i2++
		} else {
			target[k] = nums3[i3]
			i3++
		}
		k++
	}
	for i1 < n1 {
		target[k] = nums1[i1]
		i1++
		k++
	}
	for i2 < n2 {
		target[k] = nums2[i2]
		i2++
		k++
	}
	for i3 < n3 {
		target[k] = nums3[i3]
		i3++
		k++
	}

}

// 合并有序数组，保留至多 k 个元素
// https://codeforces.com/problemset/problem/587/C
func MergeTwoSortedArrayWithLimit(a, b []int, k int) []int {
	i, n := 0, len(a)
	j, m := 0, len(b)
	res := make([]int, 0, min(n+m, k))
	for len(res) < k {
		if i == n {
			if len(res)+m-j > k {
				res = append(res, b[j:j+k-len(res)]...)
			} else {
				res = append(res, b[j:]...)
			}
			break
		}
		if j == m {
			if len(res)+n-i > k {
				res = append(res, a[i:i+k-len(res)]...)
			} else {
				res = append(res, a[i:]...)
			}
			break
		}
		if a[i] < b[j] {
			res = append(res, a[i])
			i++
		} else {
			res = append(res, b[j])
			j++
		}
	}
	return res
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
