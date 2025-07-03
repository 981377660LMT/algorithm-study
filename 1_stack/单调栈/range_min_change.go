package main

import (
	"cmp"
	"slices"
)

func main() {
	arr := []int{5, 3, 4, 2, 6}
	defaultVal := 1 << 30
	res := RangeMinChange(arr, defaultVal)
	for i, changes := range res {
		for _, change := range changes {
			println(i, change.lStart, change.lEnd, change.oldMin, change.newMin)
		}
	}
}

type changeInfo[T cmp.Ordered] struct {
	lStart, lEnd   int
	oldMin, newMin T
}

// 维护区间最小值的变化历史。
// 返回：res[i]，表示右端点r=i+1时，所有受影响区间[l,r)的最小值变化记录：(l, r, old_min, new_min)
// 每次右端点推进，所有被当前元素“刷新”最小值的区间都会被记录下来，适用于区间DP、单调栈优化等场景。
func RangeMinChange[T cmp.Ordered](arr []T, defaultVal T) [][]changeInfo[T] {
	n := len(arr)
	res := make([][]changeInfo[T], n)
	stack := make([]changeInfo[T], 0, n)
	for i, v := range arr {
		res[i] = append(res[i], changeInfo[T]{lStart: i, lEnd: i + 1, oldMin: defaultVal, newMin: v})
		ptr := i
		for len(stack) > 0 {
			lc := stack[len(stack)-1]
			if lc.newMin <= v {
				break
			}
			res[i] = append(res[i], changeInfo[T]{lStart: lc.lStart, lEnd: lc.lEnd, oldMin: lc.newMin, newMin: v})
			ptr = lc.lStart
			stack = stack[:len(stack)-1]
		}
		stack = append(stack, changeInfo[T]{lStart: ptr, lEnd: i + 1, newMin: v})
		slices.Reverse(res[i])
	}
	return res
}
