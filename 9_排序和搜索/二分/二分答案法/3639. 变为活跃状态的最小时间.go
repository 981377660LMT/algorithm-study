package main

// 3639. 变为活跃状态的最小时间
// https://leetcode.cn/problems/minimum-time-to-activate-string/description/
//
// 给你一个长度为 n 的字符串 s 和一个整数数组 order，其中 order 是范围 [0, n - 1] 内数字的一个 排列 。
// 从时间 t = 0 开始，在每个时间点，将字符串 s 中下标为 order[t] 的字符替换为 '*'。
// 如果 子字符串 包含 至少 一个 '*' ，则认为该子字符串有效。
// 如果字符串中 有效子字符串 的总数大于或等于 k，则称该字符串为 活跃 字符串。
// 返回字符串 s 变为 活跃 状态的最小时间 t。如果无法变为活跃状态，返回 -1。
func minTime(s string, order []int, k int) int {
	n := len(s)
	if n*(n+1)/2 < k {
		return -1
	}

	visited := make([]int, n) // 避免在二分内部反复创建/初始化列表
	check := func(mid int) bool {
		marker := mid + 1
		for i := 0; i <= mid; i++ {
			visited[order[i]] = marker
		}

		count := 0
		pre := -1
		for i, v := range visited {
			if v == marker {
				pre = i
			}
			count += pre + 1
			if count >= k {
				return true
			}
		}
		return false
	}

	left, right := 0, n-1
	for left <= right {
		mid := (left + right) / 2
		if check(mid) {
			right = mid - 1
		} else {
			left = mid + 1
		}
	}
	return left
}
