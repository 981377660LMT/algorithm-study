package main

// 对每个固定的右端点`right(0<=right<n)`，找到最小的左端点`minLeft`，
// 使得滑动窗口内的元素满足`predicate(minLeft,right)`成立.
// 如果不存在，`minLeft`为`n`.
func GetMinLeft(
	n int32,
	append func(right int32),
	popLeft func(left int32),
	predicate func(left, right int32) bool,
) []int32 {
	minLeft := make([]int32, n)
	left := int32(0)
	for right := int32(0); right < n; right++ {
		append(right)
		for left <= right && !predicate(left, right) {
			popLeft(left)
			left++
		}
		if left > right {
			minLeft[right] = n
		} else {
			minLeft[right] = left
		}
	}
	return minLeft
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

	minLeft := GetMinLeft(n, append, popLeft, predicate)
	res := int32(0)
	for i := int32(0); i < n; i++ {
		res = max32(res, i-minLeft[i]+1)
	}
	return int(res)
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
