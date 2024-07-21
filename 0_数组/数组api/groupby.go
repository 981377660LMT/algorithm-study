package main

import (
	"fmt"
)

func main() {
	arr := []int32{10, 10, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}

	EnumerateGroup(arr, func(start, end int) {
		fmt.Println(arr[start:end])
	})

	EnumerateGroupByKey(
		arr,
		func(index int) int {
			return index / 3
		},
		func(start, end int) {
			fmt.Println(arr[start:end])
		})

	EnumerateGroupByGroupWhile(
		len(arr),
		func(left, right int) bool {
			return right-left+1 <= 4 // 最多4个元素为一组
		},
		func(start, end int) {
			fmt.Println(arr[start:end])
		},
		false,
	)
}

// 2760. 最长奇偶子数组
// https://leetcode.cn/problems/longest-even-odd-subarray-with-threshold/description/
func longestAlternatingSubarray(nums []int, threshold int) int {
	predicate := func(left, right int) bool {
		if nums[left]&1 != 0 {
			return false
		}
		if right != left && nums[right]&1 == nums[right-1]&1 {
			return false
		}
		if nums[right] > threshold {
			return false
		}
		return true
	}

	res := 0
	EnumerateGroupByGroupWhile(
		len(nums),
		predicate,
		func(start, end int) { res = max(res, end-start) },
		true,
	)
	return res
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

// !分割数组，每段中不同数字的个数不超过k个，求最少段数.
func Solve(nums []int, k int) int {
	nums = append(nums[:0:0], nums...)
	n := len(nums)
	D := NewDictionary[int]()
	for i := 0; i < n; i++ {
		nums[i] = D.Id(nums[i])
	}

	counter := make([]int, D.Size())
	count := 0
	res := 0

	ptr := 0
	for ptr < n {
		counter[nums[ptr]]++
		count = 1
		start := ptr
		ptr++
		for ptr < n {
			v := nums[ptr]
			if counter[v] == 0 {
				count++
			}
			if count > k {
				break
			}
			counter[v]++
			ptr++
		}
		res++
		for i := start; i < ptr; i++ {
			counter[nums[i]]--
		}
	}
	return res
}

// 遍历连续相同元素的分组.相当于python中的`itertools.groupby`.
func EnumerateGroup[E comparable](arr []E, f func(start, end int)) {
	n := len(arr)
	end := 0
	for end < n {
		start := end
		leader := arr[end]
		end++
		for end < n && arr[end] == leader {
			end++
		}
		f(start, end) // [start, end)
	}
}

// 遍历连续key相同元素的分组.
func EnumerateGroupByKey[E any, K comparable](arr []E, key func(index int) K, f func(start, end int)) {
	n := len(arr)
	end := 0
	for end < n {
		start := end
		leader := key(end)
		end++
		for end < n && key(end) == leader {
			end++
		}
		f(start, end) // [start, end)
	}
}

// 遍历分组(分组循环).
//
//	groupWhile: 返回 true 表示 [left, curRight] 内的元素在同一组.
//	skipFalsySingleValueGroup: 是否跳过 groupWhile 结果为 false 且只有一个元素的分组.
func EnumerateGroupByGroupWhile(
	n int,
	groupWhile func(left, curRight int) bool,
	f func(start, end int),
	skipFalsySingleValueGroup bool,
) {
	end := 0
	for end < n {
		start := end
		for end < n && groupWhile(start, end) {
			end++
		}
		isFalsySingleValueGroup := end == start
		if isFalsySingleValueGroup {
			end++
			if skipFalsySingleValueGroup {
				continue
			}
		}
		f(start, end)
	}
}

type Dictionary[V comparable] struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary[V comparable]() *Dictionary[V] {
	return &Dictionary[V]{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary[V]) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary[V]) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary[V]) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary[V]) Size() int {
	return len(d._idToValue)
}
