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

// !分割数组，每段中不同数字的个数不超过k个，求最少段数.
func Solve(nums []int, k int) int {
	nums = append(nums[:0:0], nums...)
	n := len(nums)
	D := NewDictionary()
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
func EnumerateGroupByDivider(arr []Element, isDivider func(elementIndex int, curGroup []Element) bool, f func(group []Element, start, end int)) {
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

type V = int
type Dictionary struct {
	_idToValue []V
	_valueToId map[V]int
}

// A dictionary that maps values to unique ids.
func NewDictionary() *Dictionary {
	return &Dictionary{
		_valueToId: map[V]int{},
	}
}
func (d *Dictionary) Id(value V) int {
	res, ok := d._valueToId[value]
	if ok {
		return res
	}
	id := len(d._idToValue)
	d._idToValue = append(d._idToValue, value)
	d._valueToId[value] = id
	return id
}
func (d *Dictionary) Value(id int) V {
	return d._idToValue[id]
}
func (d *Dictionary) Has(value V) bool {
	_, ok := d._valueToId[value]
	return ok
}
func (d *Dictionary) Size() int {
	return len(d._idToValue)
}
