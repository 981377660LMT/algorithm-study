package main

import "fmt"

// word = "ieaouqqieaouqq", k = 1
func main() {
	// fmt.Println(countOfSubstrings("ieaouqqieaouqq", 1)) // 3
	// word = "aeiou", k = 0
	fmt.Println(countOfSubstrings("aeiou", 0)) // 1
}

func isVowel(b byte) bool {
	return b == 'a' || b == 'e' || b == 'i' || b == 'o' || b == 'u'
}

func countOfSubstrings(word string, k int) int64 {
	A := AlphaPresum(Str(word), 26, 'a')

	calc := func(word string, k int) int64 {
		res, left, n := 0, 0, len(word)
		consonantCount := 0
		for right := 0; right < n; right++ {
			if !isVowel(word[right]) {
				consonantCount++
			}
			for left <= right && consonantCount > k {
				if !isVowel(word[left]) {
					consonantCount--
				}
				left++
			}

			maxLeft := MinLeft(right+1, func(left int) bool {
				if A(left, right+1, 'a') == 0 {
					return true
				}
				if A(left, right+1, 'e') == 0 {
					return true
				}
				if A(left, right+1, 'i') == 0 {
					return true
				}
				if A(left, right+1, 'o') == 0 {
					return true
				}
				if A(left, right+1, 'u') == 0 {
					return true
				}
				return false
			}, left)

			res += maxLeft - left + 1
		}

		return int64(res)
	}

	return calc(word, k) - calc(word, k-1)
}

type Str = []byte

// 换成数组
func AlphaPresum(s Str, sigma int, offset int) func(start, end int, ord int) int {
	preSum := make([][26]int32, len(s)+1)
	for i := 1; i <= len(s); i++ {
		preSum[i] = preSum[i-1]
		preSum[i][int(s[i-1])-offset]++
	}
	return func(start, end int, ord int) int {
		return int(preSum[end][ord-offset] - preSum[start][ord-offset])
	}
}

// 返回最大的 right 使得 [left,right) 内的值满足 check.
// !注意check内的right不包含,使用时需要right-1.
// right<=upper.
func MaxRight(left int, check func(right int) bool, upper int) int {
	ok, ng := left, upper+1
	for ok+1 < ng {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}

// 返回最小的 left 使得 [left,right) 内的值满足 check.
// left>=lower.
func MinLeft(right int, check func(left int) bool, lower int) int {
	ok, ng := right, lower-1
	for ng+1 < ok {
		mid := (ok + ng) >> 1
		if check(mid) {
			ok = mid
		} else {
			ng = mid
		}
	}
	return ok
}
