package main

// 允许将数组中的任意字符替换为target字符k次，求target字符的最大连续长度.
func Fix[T comparable](arr []T, target T, k int) int {
	left := 0
	res := 0
	for right := 0; right < len(arr); right++ {
		if arr[right] != target {
			k--
		}
		for k < 0 {
			if arr[left] != target {
				k++
			}
			left++
		}
		res = max(res, right-left+1)
	}
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
