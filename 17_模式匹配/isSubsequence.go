package main

import "fmt"

func main() {
	fmt.Println(IsSubsequence([]int{1, 2, 3, 4, 5}, []int{1, 3, 5})) // true
	fmt.Println(MatchSubsequenceString("abcde", "ace"))              // [0 1 1 2 2 3]
}

// 判断`shorter`是否是`longer`的子序列.
// complexity O(n1 + n2)
func IsSubsequenceString(longer, shorter string) bool {
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

// 返回`longer`的每个前缀中的子序列匹配`shorter`的最大长度.
func MatchSubsequenceString(longer, shorter string) []int {
	res := make([]int, len(longer)+1)
	i, j := 0, 0
	for i < len(longer) && j < len(shorter) {
		if longer[i] == shorter[j] {
			j++
		}
		i++
		res[i] = j
	}
	for i++; i < len(res); i++ {
		res[i] = j
	}
	return res
}

// 判断`shorter`是否是`longer`的子序列.
// 如果需要多次匹配，使用`子序列自动机`.
// complexity O(n1 + n2)
func IsSubsequence[S ~[]E, E comparable](longer, shorter S) bool {
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

// 返回`longer`的每个前缀中的子序列匹配`shorter`的最大长度.
func MatchSubsequence[S ~[]E, E comparable](longer, shorter S) []int {
	res := make([]int, len(longer)+1)
	i, j := 0, 0
	for i < len(longer) && j < len(shorter) {
		if longer[i] == shorter[j] {
			j++
		}
		i++
		res[i] = j
	}
	for i++; i < len(res); i++ {
		res[i] = j
	}
	return res
}
