package main

import "fmt"

func main() {
	arr := []Element{10, 10, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	EnumerateGroup(arr, func(group []Element, _, _ int) {
		fmt.Println(group)
	})

	EnumerateGroupByDivider(
		arr,
		func(index int, curGroup []Element) bool {
			return len(curGroup) == 3 // 最多3个元素为一组
		},
		func(group []Element, _, _ int) {
			fmt.Println(group)
		})
}

type Element = int

// 遍历连续相同元素的分组.相当于python中的`itertools.groupby`.
func EnumerateGroup(arr []Element, f func(group []Element, start, end int)) {
	ptr := 0
	n := len(arr)
	for ptr < n {
		leader := arr[ptr]
		group := []Element{leader}
		start := ptr
		ptr++
		for ptr < n && arr[ptr] == leader {
			group = append(group, arr[ptr])
			ptr++
		}
		f(group, start, ptr)
	}
}

// 遍历连续key相同元素的分组.
func EnumerateGroupByKey(arr []Element, key func(index int) interface{}, f func(group []Element, start, end int)) {
	ptr := 0
	n := len(arr)
	for ptr < n {
		leader := key(ptr)
		group := []Element{arr[ptr]}
		start := ptr
		ptr++
		for ptr < n && key(ptr) == leader {
			group = append(group, arr[ptr])
			ptr++
		}
		f(group, start, ptr)
	}
}

// 遍历分组(分组循环).
//  isDivider: 判断当前元素是否为分组的分界点.如果返回true,则以当前元素为分界点,新建下一个分组.
func EnumerateGroupByDivider(arr []Element, isDivider func(index int, curGroup []Element) bool, f func(group []Element, start, end int)) {
	ptr := 0
	n := len(arr)
	for ptr < n {
		leader := arr[ptr]
		group := []Element{leader}
		start := ptr
		ptr++
		for ptr < n && !isDivider(ptr, group) {
			group = append(group, arr[ptr])
			ptr++
		}
		f(group, start, ptr)
	}
}
