package main

import "fmt"

func main() {
	// "abcccccdddd"
	s := "abcccccdddd"
	nexts := GetNext(s)
	fmt.Println(CountIndexOfAll(s, "cc", 0, nexts))
	fmt.Println(IndexOfAll(s, "cc", 0, nexts))
}

type Str = string

func GetNext(pattern Str) []int {
	next := make([]int, len(pattern))
	j := 0
	for i := 1; i < len(pattern); i++ {
		for j > 0 && pattern[i] != pattern[j] {
			j = next[j-1]
		}
		if pattern[i] == pattern[j] {
			j++
		}
		next[i] = j
	}
	return next
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中匹配次数.
// nexts 数组为nil时, 会调用GetNext(shorter)求nexts数组.
func CountIndexOfAll(longer Str, shorter Str, position int, nexts []int) int {
	if len(shorter) == 0 {
		return 0
	}
	if len(longer) < len(shorter) {
		return 0
	}
	res := 0
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res++
			hitJ = nexts[hitJ-1] // 不允许重叠时 hitJ = 0
		}
	}
	return res
}

// `O(n+m)` 寻找 `shorter` 在 `longer` 中的所有匹配位置.
func IndexOfAll(longer string, shorter string, position int, nexts []int) []int {
	if len(shorter) == 0 {
		return []int{0}
	}
	if len(longer) < len(shorter) {
		return nil
	}
	res := []int{}
	if nexts == nil {
		nexts = GetNext(shorter)
	}
	hitJ := 0
	for i := position; i < len(longer); i++ {
		for hitJ > 0 && longer[i] != shorter[hitJ] {
			hitJ = nexts[hitJ-1]
		}
		if longer[i] == shorter[hitJ] {
			hitJ++
		}
		if hitJ == len(shorter) {
			res = append(res, i-len(shorter)+1)
			hitJ-- // 不允许重叠时 hitJ = 0
		}
	}
	return res
}
