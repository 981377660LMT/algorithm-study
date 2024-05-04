// https://leetcode.cn/problems/median-of-two-sorted-arrays/

package main

const INF int = int(1e18)

// O(log(m+n)) 求两个有序数组的中位数.
func MedianOfSortedArrays(sortedNums1, sortedNusm2 []int) float64 {
	if len(sortedNums1) > len(sortedNusm2) {
		sortedNums1, sortedNusm2 = sortedNusm2, sortedNums1
	}

	len1, len2 := len(sortedNums1), len(sortedNusm2)
	left, right := 0, len1
	max1, min2 := 0, 0 // 前一部分的最大值和后一部分的最小值

	for left <= right {
		i := (left + right) >> 1
		j := (len1+len2+1)>>1 - i

		a1 := -INF
		if i != 0 {
			a1 = sortedNums1[i-1]
		}
		b1 := INF
		if i != len1 {
			b1 = sortedNums1[i]
		}
		a2 := -INF
		if j != 0 {
			a2 = sortedNusm2[j-1]
		}
		b2 := INF
		if j != len2 {
			b2 = sortedNusm2[j]
		}

		if a1 <= b2 {
			if a1 > a2 {
				max1 = a1
			} else {
				max1 = a2
			}
			if b1 < b2 {
				min2 = b1
			} else {
				min2 = b2
			}
			left = i + 1
		} else {
			right = i - 1
		}
	}

	if (len1+len2)&1 == 1 {
		return float64(max1)
	} else {
		return float64(max1+min2) / 2
	}
}
