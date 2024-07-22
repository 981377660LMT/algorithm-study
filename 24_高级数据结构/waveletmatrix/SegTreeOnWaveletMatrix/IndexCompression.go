package main

import (
	"fmt"
	"sort"
)

func main() {
	{
		arr := []int{1449, 12, 12, 987}
		newArr, index := indexCompressionDistinctSmall(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
	}

	{
		arr := []int{1449, 12, 12, 987}
		newArr, index := indexCompressionSameSmall(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
	}

	{
		arr := []int{1449, 12, 12, 987, 1e9 + 7}
		newArr, index := indexCompressionDistinctLarge(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
		fmt.Println(index(1e9 + 7))
	}

	{
		arr := []int{1449, 12, 12, 987, 1e9 + 7}
		newArr, index := indexCompressionSameLarge(arr)
		fmt.Println(newArr)
		fmt.Println(index(1000000))
		fmt.Println(index(1e9 + 7))
	}
}

// IndexCompression 用于对数组进行压缩.
//
//	 allowSame: 相同大小的元素编号是否能相同.
//		true -> [2,3,2] -> [0,1,0]
//		false -> [2,3,2] -> [0,2,0]
//	 smallRange: 数据极差较小(不超过1e7).
func IndexCompression(arr []int, allowSame bool, smallRange bool) (compressedArr []int32, index func(int) int32) {
	if allowSame {
		if smallRange {
			return indexCompressionSameSmall(arr)
		} else {
			return indexCompressionSameLarge(arr)
		}
	} else {
		if smallRange {
			return indexCompressionDistinctSmall(arr)
		} else {
			return indexCompressionDistinctLarge(arr)
		}
	}
}

func indexCompressionSameLarge(arr []int) (compressedArr []int32, index func(int) int32) {
	var data []int
	order := argSort(arr)
	compressedArr = make([]int32, len(arr))
	for _, v := range order {
		if len(data) == 0 || data[len(data)-1] != arr[v] {
			data = append(data, arr[v])
		}
		compressedArr[v] = int32(len(data) - 1)
	}
	data = data[:len(data):len(data)]
	index = func(x int) int32 { return int32(sort.SearchInts(data, x)) }
	return
}

func indexCompressionDistinctLarge(arr []int) (compressedArr []int32, index func(int) int32) {
	var data []int
	order := argSort(arr)
	compressedArr = make([]int32, len(arr))
	for _, v := range order {
		compressedArr[v] = int32(len(data))
		data = append(data, arr[v])
	}
	data = data[:len(data):len(data)]
	index = func(x int) int32 { return int32(sort.SearchInts(data, x)) }
	return
}

func indexCompressionSameSmall(arr []int) (compressedArr []int32, index func(int) int32) {
	var min_, max_ int
	var data []int32
	compressedArr = make([]int32, len(arr))
	for i, v := range arr {
		compressedArr[i] = int32(v)
	}
	min32, max32 := int32(0), int32(-1)
	if len(compressedArr) > 0 {
		for _, x := range compressedArr {
			if x < min32 {
				min32 = x
			}
			if x > max32 {
				max32 = x
			}
		}
	}
	data = make([]int32, max32-min32+2)
	for _, x := range compressedArr {
		data[x-min32+1] = 1
	}
	for i := 0; i < len(data)-1; i++ {
		data[i+1] += data[i]
	}
	for i, v := range compressedArr {
		compressedArr[i] = data[v-min32]
	}
	min_, max_ = int(min32), int(max32)
	index = func(x int) int32 { return data[clamp(x-min_, 0, max_-min_+1)] }
	return
}

func indexCompressionDistinctSmall(arr []int) (compressedArr []int32, index func(int) int32) {
	var min_, max_ int
	var data []int32
	compressedArr = make([]int32, len(arr))
	for i, v := range arr {
		compressedArr[i] = int32(v)
	}
	min32, max32 := int32(0), int32(-1)
	if len(compressedArr) > 0 {
		for _, x := range compressedArr {
			if x < min32 {
				min32 = x
			}
			if x > max32 {
				max32 = x
			}
		}
	}
	data = make([]int32, max32-min32+2)
	for _, x := range compressedArr {
		data[x-min32+1]++
	}
	for i := 0; i < len(data)-1; i++ {
		data[i+1] += data[i]
	}
	for i, x := range compressedArr {
		compressedArr[i] = data[x-min32]
		data[x-min32]++
	}
	copy(data[1:], data)
	data[0] = 0
	min_, max_ = int(min32), int(max32)
	index = func(x int) int32 { return data[clamp(x-min_, 0, max_-min_+1)] }
	return
}

func clamp(x, min, max int) int {
	if x < min {
		return min
	}
	if x > max {
		return max
	}
	return x
}

func argSort(nums []int) []int32 {
	order := make([]int32, len(nums))
	for i := int32(0); i < int32(len(order)); i++ {
		order[i] = i
	}
	sort.Slice(order, func(i, j int) bool { return nums[order[i]] < nums[order[j]] })
	return order
}
