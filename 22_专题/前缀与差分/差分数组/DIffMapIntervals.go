// !遍历区间的差分字典.

package main

import (
	"fmt"
	"sort"
)

func main() {
	// 示例：创建 DiffMapIntervals 对象并添加区间差分.
	D := NewDiffMapIntervals()
	D.Add(1, 3, 10) // 在区间 [1,3] 加 10
	D.Add(2, 5, -5) // 在区间 [2,5] 加 -5
	D.Add(4, 6, 3)  // 在区间 [4,6] 加 3

	// 枚举从 1 到 7 的累积和区间.
	fmt.Println("累积和区间 (value, start, end):")
	D.Enumerate(1, 7, func(curSum, l, r int) {
		fmt.Printf("(%d, %d, %d)\n", curSum, l, r)
	})
}

// 遍历区间的差分字典.
type DiffMapIntervals struct {
	mp map[int]int
}

// NewDiffMapIntervals 创建一个新的 DiffMapIntervals 对象.
func NewDiffMapIntervals() *DiffMapIntervals {
	return &DiffMapIntervals{
		mp: make(map[int]int),
	}
}

// 闭区间 [left, right] 加上值 x.
func (d *DiffMapIntervals) Add(left, right, x int) {
	if left <= right {
		d.mp[left] += x
		d.mp[right+1] -= x
	}
}

// 枚举从 since 到 until 的累积和区间.
// 对于每个区间，调用回调函数 callback(sum, l, r),
// 表示在区间 [l, r] 内的累积和为 sum.
func (d *DiffMapIntervals) Enumerate(since, until int, f func(sum, l, r int)) {
	curSum := 0
	pre := since
	keys := make([]int, 0, len(d.mp))
	for k := range d.mp {
		keys = append(keys, k)
	}
	sort.Ints(keys)
	for _, t := range keys {
		if t > until {
			break
		}
		if d.mp[t] == 0 {
			continue
		}
		if pre <= t-1 {
			f(curSum, pre, t-1)
		}
		curSum += d.mp[t]
		pre = t
	}
	if pre <= until {
		f(curSum, pre, until)
	}
}
