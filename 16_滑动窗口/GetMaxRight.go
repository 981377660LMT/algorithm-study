package main

import "fmt"

// 对每个固定的左端点`left(0<=left<n)`，找到最大的右端点`maxRight`，
// 使得滑动窗口内的元素满足`predicate(left,maxRight)`成立.
// 如果不存在，`maxRight`为-1.
func GetMaxRight(
	n int32,
	append func(right int32),
	popLeft func(left int32),
	predicate func(left, right int32) bool,
) []int32 {
	maxRight := make([]int32, n)
	right := int32(0)
	visitedRight := make([]bool, n)
	visitRight := func(right int32) {
		if visitedRight[right] {
			return
		}
		visitedRight[right] = true
		append(right)
	}

	for left := int32(0); left < n; left++ {
		if right < left {
			right = left
		}
		for right < n {
			visitRight(right)
			if predicate(left, right) {
				right++
			} else {
				break
			}
		}

		if right == n {
			for i := left; i < n; i++ {
				maxRight[i] = n - 1
			}
			break
		}

		if tmp := right - 1; tmp >= left {
			maxRight[left] = tmp
		} else {
			maxRight[left] = -1
		}
		popLeft(left)
	}

	return maxRight
}

// https://leetcode.cn/problems/longest-substring-without-repeating-characters/description/
func lengthOfLongestSubstring(s string) int {
	n := int32(len(s))
	counter := [256]int32{}
	dupCount := int32(0)

	append := func(right int32) {
		ord := s[right] - 'a'
		counter[ord]++
		if counter[ord] == 2 {
			dupCount++
		}
	}

	popLeft := func(left int32) {
		ord := s[left] - 'a'
		counter[ord]--
		if counter[ord] == 1 {
			dupCount--
		}
	}

	predicate := func(left, right int32) bool {
		return dupCount == 0
	}

	maxRight := GetMaxRight(n, append, popLeft, predicate)
	res := int32(0)
	for left, right := range maxRight {
		res = max32(res, right-int32(left)+1)
	}
	return int(res)
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func main() {
	nums := []int32{100, 99, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	curSum := int32(0)
	append := func(right int32) { curSum += nums[right] }
	popLeft := func(left int32) { curSum -= nums[left] }
	predicate := func(left, right int32) bool { return curSum <= 8 }
	maxRight := GetMaxRight(int32(len(nums)), append, popLeft, predicate)
	fmt.Println(maxRight)
}
