package main

import "sort"

// https://leetcode.cn/problems/number-of-divisible-triplet-sums/description/
// 1 <= nums.length <= 1000
// 1 <= nums[i] <= 109
// 1 <= d <= 109
// !固定左端点，剩下的问题就是两数之和，O(n^2)
func divisibleTripletCount(nums []int, d int) int {
	mods := make([]int, len(nums))
	for i, v := range nums {
		mods[i] = v % d
	}
	S := RangeFreq(mods)
	res := 0
	n := len(nums)
	for i := 0; i < n-2; i++ {
		for j := i + 1; j < n-1; j++ {
			need := (-mods[i] - mods[j]) % d
			if need < 0 {
				need += d
			}
			res += S(j+1, n, need)
		}
	}
	return res
}

type V = int

func RangeFreq(arr []V) func(start, end int, value V) int {
	mp := make(map[V][]int)
	for i, v := range arr {
		mp[v] = append(mp[v], i)
	}

	return func(start, end int, value V) int {
		if start < 0 {
			start = 0
		}
		if end > len(arr) {
			end = len(arr)
		}
		if start >= end {
			return 0
		}

		if indexes, ok := mp[value]; !ok {
			return 0
		} else {
			return sort.SearchInts(indexes, end) - sort.SearchInts(indexes, start)
		}
	}
}
