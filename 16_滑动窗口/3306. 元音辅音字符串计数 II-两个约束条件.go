// 3306. 元音辅音字符串计数 II
// https://leetcode.cn/problems/count-of-substrings-containing-every-vowel-and-k-consonants-ii/description/
// 给你一个字符串 word 和一个 非负 整数 k。
// 返回 word 的 子字符串中，每个元音字母（'a'、'e'、'i'、'o'、'u'）至少 出现一次，并且 恰好 包含 k 个辅音字母的子字符串的总数。
//
// !先求出第一个约束条件对应的left，再二分出第二个约束条件对应的maxLeft.

package main

func isVowel(b byte) bool {
	return b == 'a' || b == 'e' || b == 'i' || b == 'o' || b == 'u'
}

func countOfSubstrings(word string, k int) int64 {
	A := AlphaPresum(Str(word))

	// 最多包含k个辅音字母的子字符串个数
	calc := func(word string, k int) int64 {
		res, left, n := 0, 0, len(word)
		consonantCount := 0
		for right := 0; right < n; right++ {
			if !isVowel(word[right]) {
				consonantCount++
			}
			// 包含<=k个辅音字母
			for left <= right && consonantCount > k {
				if !isVowel(word[left]) {
					consonantCount--
				}
				left++
			}

			// 每个元音字母（'a'、'e'、'i'、'o'、'u'）至少出现一次
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

const SIGMA int = 26
const OFFSET int = 97

func AlphaPresum(s Str) func(start, end int, ord int) int {
	preSum := make([][SIGMA]int32, len(s)+1)
	for i := 1; i <= len(s); i++ {
		preSum[i] = preSum[i-1]
		preSum[i][int(s[i-1])-OFFSET]++
	}
	return func(start, end int, ord int) int {
		return int(preSum[end][ord-OFFSET] - preSum[start][ord-OFFSET])
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
