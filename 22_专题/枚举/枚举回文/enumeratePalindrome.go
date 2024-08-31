package main

import (
	"math"
	"sort"
	"strconv"
)

// 3272. 统计好整数的数目
// https://leetcode.cn/problems/find-the-count-of-good-integers/description/
func countGoodIntegers(n int, k int) int64 {
	fac := []int{1}
	for i := 1; i <= 15; i++ {
		fac = append(fac, fac[i-1]*i)
	}
	counter := [10]int{}

	calc := func(palindrome int) int {
		for i := range counter {
			counter[i] = 0
		}
		cur := palindrome
		m := 0
		for cur > 0 {
			counter[cur%10]++
			cur /= 10
			m++
		}

		res := (m - counter[0]) * fac[m-1]
		for i := 0; i < 10; i++ {
			res /= fac[counter[i]]
		}
		return res
	}

	res := 0
	visited := map[string]struct{}{}
	EnumeratePalindromeByLen(n, n, func(palindrome int) bool {
		if palindrome%k == 0 {
			digits := []byte(strconv.Itoa(palindrome))
			sort.Slice(digits, func(i, j int) bool { return digits[i] < digits[j] })
			key := string(digits)
			if _, has := visited[key]; !has {
				visited[key] = struct{}{}
				res += calc(palindrome)
			}
		}
		return false
	})

	return int64(res)
}

// 从小到大遍历`[min,max]`闭区间内的回文数.返回 true 可提前终止遍历.
// https://github.com/EndlessCheng/codeforces-go
func EnumeratePalindrome(min, max int, f func(palindrome int) bool) {
	if min > max {
		return
	}

	minLen := len(strconv.Itoa(min))
	startBase := int(math.Pow10((minLen - 1) >> 1))
	for base := startBase; ; base *= 10 {
		// 生成奇数长度回文数，例如 base = 10，生成的范围是 101 ~ 999
		for i := base; i < base*10; i++ {
			x := i
			for t := i / 10; t > 0; t /= 10 {
				x = x*10 + t%10
			}
			if x > max {
				return
			}
			if x >= min {
				if f(x) {
					return
				}
			}
		}

		// 生成偶数长度回文数，例如 base = 10，生成的范围是 1001 ~ 9999
		for i := base; i < base*10; i++ {
			x := i
			for t := i; t > 0; t /= 10 {
				x = x*10 + t%10
			}
			if x > max {
				return
			}
			if x >= min {
				if f(x) {
					return
				}
			}
		}
	}
}

// 从小到大遍历长度在`[minLen,maxLen]`闭区间内的回文数.返回 true 可提前终止遍历.
// maxLen <= 14.
func EnumeratePalindromeByLen(minLen, maxLen int, f func(palindrome int) bool) {
	if minLen > maxLen {
		return
	}
	min, max := int(math.Pow10(minLen-1)), int(math.Pow10(maxLen)-1)
	EnumeratePalindrome(min, max, f)
}
