// 找出并返回满足要求的不同的子数组数，要求子数组中最多 k 个可被 p 整除的元素。

package main

import (
	"fmt"
	"index/suffixarray"
	"reflect"
	"unsafe"
)

func countDistinct(nums []int, k int, p int) (res int) {
	n := len(nums)

	mods := make([]int, n)
	for i := range mods {
		mods[i] = nums[i] % p
	}

	// 1. 先用双指针O(n)的时间计算出所有满足条件的子数组的数量 注意要枚举后缀(固定left 移动right)
	right, countK := 0, 0
	suffixLen := make([]int, n) // 记录每个后缀取到的长度
	for left := 0; left < n; left++ {
		for right < n && countK+boolToInt((mods[right] == 0)) <= k {
			countK += boolToInt((mods[right] == 0))
			right++
		}

		res += right - left
		suffixLen[left] = right - left
		countK -= boolToInt(mods[left] == 0)
	}

	// 2. height数组去重
	// nums to byte array 怎么让后缀数组支持数字?
	byteArray := make([]byte, n)
	for i := 0; i < n; i++ {
		byteArray[i] = byte(nums[i])
	}

	sa, _, height := suffixArray(byteArray)
	// 计算子串重复数量 按后缀排序的顺序枚举后缀 lcp(height)去重
	for i := 0; i < n-1; i++ {
		suffix1, suffix2 := sa[i], sa[i+1]
		subLen1, subLen2 := suffixLen[suffix1], suffixLen[suffix2]
		res -= min(height[i+1], subLen1, subLen2)
	}

	return
}

// https://github.dev/EndlessCheng/codeforces-go/copypasta/strings.go
func suffixArray(s []byte) ([]int32, []int, []int) {
	n := len(s)

	sa := *(*[]int32)(unsafe.Pointer(reflect.ValueOf(suffixarray.New(s)).Elem().FieldByName("sa").Field(0).UnsafeAddr()))

	rank := make([]int, n)
	for i := range rank {
		rank[sa[i]] = i
	}

	height := make([]int, n)
	h := 0
	for i, rk := range rank {
		if h > 0 {
			h--
		}
		if rk > 0 {
			for j := int(sa[rk-1]); i+h < n && j+h < n && s[i+h] == s[j+h]; h++ {
			}
		}
		height[rk] = h
	}

	return sa, rank, height
}

func min(nums ...int) int {
	min := nums[0]
	for _, num := range nums {
		if num < min {
			min = num
		}
	}
	return min
}

func boolToInt(b bool) int {
	if b {
		return 1
	}
	return 0
}

func main() {
	fmt.Println(countDistinct([]int{2, 3, 3, 2, 2}, 2, 2))
	fmt.Println(countDistinct([]int{16, 17, 4, 12, 3}, 4, 1))
	fmt.Println(countDistinct([]int{1, 100, 20, 1}, 1, 4))
}
