package main

import (
	"fmt"
	"math"
	"strconv"
)

func main() {
	EnumeratePalindrome(12, 1000, func(palindrome int) bool {
		fmt.Println(palindrome)
		if palindrome == 121 {
			return true
		}
		return false
	})
}

// 从小到大遍历`[min,max]`闭区间内的回文数.返回 true 可提前终止遍历.
// https://github.com/EndlessCheng/codeforces-go
func EnumeratePalindrome(min, max int, f func(palindrome int) bool) {
	if min > max {
		return
	}

	minLength := len(strconv.Itoa(min))
	startBase := int(math.Pow10((minLength - 1) >> 1))
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
