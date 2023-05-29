package main

import "fmt"

func main() {
	arr := []interface{}{"a", "a", 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12, 12}
	EnumerateGroup(arr, func(group []interface{}, _, _ int) {
		fmt.Println(group)
	})
}

// 遍历连续相同元素的分组.相当于python中的`itertools.groupby`.
func EnumerateGroup(arr []interface{}, f func(group []interface{}, start, end int)) {
	ptr := 0
	n := len(arr)
	for ptr < n {
		leader := arr[ptr]
		group := []interface{}{leader}
		start := ptr
		ptr++
		for ptr < n && arr[ptr] == leader {
			group = append(group, arr[ptr])
			ptr++
		}
		f(group, start, ptr)
	}
}
