//  StringPeriod (字符串周期)

package main

import (
	"fmt"
	"math"
)

func main() {
	nums := []int32{1, 2, 3, 4, 5, 6}
	fmt.Println(MinPeriod(int32(len(nums)), func(i int32) int32 { return nums[i] }))
	fmt.Println(MinRotatePeriod(int32(len(nums)), func(i int32) int32 { return nums[i] }))

	nums = []int32{1, 2, 3, 1, 2, 3}
	fmt.Println(MinPeriod(int32(len(nums)), func(i int32) int32 { return nums[i] }))
	fmt.Println(MinPalindromeRotatePeriod(int32(len(nums)), func(i int32) int32 { return nums[i] }))

	nums = []int32{2, 2, 3, 3, 2, 2}
	fmt.Println(MinPalindromeRotatePeriod(int32(len(nums)), func(i int32) int32 { return nums[i] }))
}

// 查询字符串的最小周期p，满足s[i+p]=s[i] (i+p<n).
func MinPeriod(n int32, f func(i int32) int32) int32 {
	next := getNext(n, f)
	return n - next[n-1]
}

// 查询字符串的最小旋转周期p，满足s[(i+p)%n]=s[i].
// 时间复杂度 O(nlogn).
func MinRotatePeriod(n int32, f func(i int32) int32) int32 {
	check := func(p int32) bool {
		for i, j := int32(0), p; i < n; i, j = i+1, j+1 {
			if j >= n {
				j -= n
			}
			if f(i) != f(j) {
				return false
			}
		}
		return true
	}

	res := n
	enumerateFactors(n, func(p int32) bool {
		if check(p) {
			res = p
			return true
		}
		return false
	})
	return res
}

// 查询字符串的最小回文旋转周期p，满足s[p..n)+s[0..p)是回文.
// 前置条件：s[0..n)是回文.
// 时间复杂度 O(nlogn).
func MinPalindromeRotatePeriod(n int32, f func(i int32) int32) int32 {
	rp := MinRotatePeriod(n, f)
	if rp%2 == 0 {
		rp /= 2
	}
	return rp
}

func isPalindrome(n int32, f func(i int32) int32) bool {
	l, r := int32(0), n-1
	for l < r {
		if f(l) != f(r) {
			return false
		}
		l++
		r--
	}
	return true
}

// `next[i]`表示`[:i+1]`这一段字符串中最长公共前后缀(不含这一段字符串本身,即真前后缀)的长度.
func getNext(n int32, f func(i int32) int32) []int32 {
	next := make([]int32, n)
	j := int32(0)
	for i := int32(1); i < n; i++ {
		for j > 0 && f(i) != f(j) {
			j = next[j-1]
		}
		if f(i) == f(j) {
			j++
		}
		next[i] = j
	}
	return next
}

// 空间复杂度为O(1)的枚举因子.枚举顺序为从小到大.
func enumerateFactors(n int32, f func(factor int32) (shouldBreak bool)) {
	if n <= 0 {
		return
	}
	i := int32(1)
	upper := int32(math.Sqrt(float64(n)))
	for ; i <= upper; i++ {
		if n%i == 0 {
			if f(i) {
				return
			}
		}
	}
	i--
	if i*i == n {
		i--
	}
	for ; i > 0; i-- {
		if n%i == 0 {
			if f(n / i) {
				return
			}
		}
	}
}
