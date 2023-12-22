package main

type Sequence = string

// 判断`shorter`是否是`longer`的子序列.
// 如果需要多次匹配，使用`子序列自动机`.
// complexity O(n1 + n2)
func IsSubsequence(longer, shorter Sequence) bool {
	if len(shorter) > len(longer) {
		return false
	}
	if len(shorter) == 0 {
		return true
	}
	n, m := len(longer), len(shorter)
	i, j := 0, 0
	for i < n && j < m {
		if longer[i] == shorter[j] {
			j++
			if j == m {
				return true
			}
		}
		i++
	}
	return false
}
