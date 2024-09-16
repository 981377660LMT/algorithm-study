package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	abc363_c()
}

// C - Avoid K Palindrome 2
// https://atcoder.jp/contests/abc363/tasks/abc363_c
// 2<=k<=N<=10
// 给定一个字符串，将字母排序，问有多少种情况，其不存在一个长度为k的回文串。
func abc363_c() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int
	fmt.Fscan(in, &n, &k)
	var s string
	fmt.Fscan(in, &s)

	bs := []byte(s)
	sort.Slice(bs, func(i, j int) bool { return bs[i] < bs[j] })
	res := 0

	isPalindrome := func(bs []byte) bool {
		for i := 0; i < len(bs)/2; i++ {
			if bs[i] != bs[len(bs)-1-i] {
				return false
			}
		}
		return true
	}
	check := func(bs []byte) bool {
		for i := 0; i <= len(bs)-k; i++ {
			if isPalindrome(bs[i : i+k]) {
				return false
			}
		}
		return true
	}

	for {
		if check(bs) {
			res++
		}
		if !nextPermutation(bs) {
			break
		}
	}

	fmt.Println(res)
}

// 原地返回下一个字典序的排列.
// 不包含重复排列.
func nextPermutation[T byte | int32 | int | string](nums []T) bool {
	left, right := len(nums)-1, len(nums)-1
	for left > 0 && nums[left-1] >= nums[left] {
		left--
	}
	if left == 0 {
		return false
	}
	last := left - 1
	for nums[right] <= nums[last] {
		right--
	}
	nums[last], nums[right] = nums[right], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return true
}

// 原地返回上一个字典序的排列.
// 不包含重复排列.
func prePermutation[T byte | int32 | int | string](nums []T) bool {
	left, right := len(nums)-1, len(nums)-1
	for left > 0 && nums[left-1] <= nums[left] {
		left--
	}
	if left == 0 {
		return false
	}
	last := left - 1
	for nums[right] >= nums[last] {
		right--
	}
	nums[last], nums[right] = nums[right], nums[last]
	for i, j := last+1, len(nums)-1; i < j; i, j = i+1, j-1 {
		nums[i], nums[j] = nums[j], nums[i]
	}
	return true
}
