package main

import "sort"

// https://leetcode.cn/problems/minimum-cost-for-cutting-cake-ii/
func minimumCost(_m int, _n int, horizontalCut []int, verticalCut []int) int64 {
	resCost, _ := MergeTwoArrayWithMinCost(horizontalCut, verticalCut)
	return int64(resCost)
}

type mergeResult struct {
	kind  uint8
	index int32
}

// 合并两个数组，最小化总代价.
// 第i个数的代价为:`(这个数之前"不同"类型的数的个数+1)*这个数`.
// 返回最小代价和合并的结果，结果每一项形如(kind, index).
func MergeTwoArrayWithMinCost(arr1 []int, arr2 []int) (resCost int, resArray []mergeResult) {
	order1 := make([]int32, len(arr1))
	for i := range order1 {
		order1[i] = int32(i)
	}
	sort.Slice(order1, func(i, j int) bool { return arr1[order1[i]] > arr1[order1[j]] })
	order2 := make([]int32, len(arr2))
	for i := range order2 {
		order2[i] = int32(i)
	}
	sort.Slice(order2, func(i, j int) bool { return arr2[order2[i]] > arr2[order2[j]] })
	resArray = make([]mergeResult, 0, len(arr1)+len(arr2))
	count1, count2 := 1, 1
	i, j := 0, 0
	for i < len(arr1) && j < len(arr2) {
		ptr1, ptr2 := order1[i], order2[j]
		if arr1[ptr1] > arr2[ptr2] {
			resCost += arr1[ptr1] * count2
			resArray = append(resArray, mergeResult{kind: 1, index: ptr1})
			i++
			count1++
		} else {
			resCost += arr2[ptr2] * count1
			resArray = append(resArray, mergeResult{kind: 2, index: ptr2})
			j++
			count2++
		}
	}
	for i < len(arr1) {
		ptr1 := order1[i]
		resCost += arr1[ptr1] * count2
		resArray = append(resArray, mergeResult{kind: 1, index: ptr1})
		i++
	}
	for j < len(arr2) {
		ptr2 := order2[j]
		resCost += arr2[ptr2] * count1
		resArray = append(resArray, mergeResult{kind: 2, index: ptr2})
		j++
	}

	return
}
