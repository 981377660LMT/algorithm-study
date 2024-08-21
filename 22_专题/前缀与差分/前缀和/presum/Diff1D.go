// 一维差分

package main

import (
	"fmt"
	"sort"
)

func main() {
	e := func() int { return 0 }
	op := func(a, b int) int { return a + b }
	inv := func(a int) int { return -a }

	diff := NewDiff1D(e, op, inv)
	diff.Init(5, func(i int32) int { return int(i) })
	diff.Add(1, 3, 1)
	diff.Add(2, 4, 2)
	diff.Add(-1, 10, 3)
	for i := int32(0); i < 5; i++ {
		fmt.Println(diff.Get(i))
	}

	// [[1,6],[3,7],[9,12],[4,13]]
	// [2,3,7,11]
	fmt.Println(fullBloomFlowers([][]int{{1, 6}, {3, 7}, {9, 12}, {4, 13}}, []int{2, 3, 7, 11}))
}

// 2251. 花期内花的数目
// https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/
func fullBloomFlowers(flowers [][]int, people []int) []int {
	var allNums []int
	for _, flower := range flowers {
		allNums = append(allNums, flower[0], flower[1])
	}
	for _, p := range people {
		allNums = append(allNums, p)
	}
	_, origin := Discretize(allNums)
	f := func(v int) int { return sort.SearchInts(origin, v) }
	for i := range flowers {
		flowers[i][0] = f(flowers[i][0])
		flowers[i][1] = f(flowers[i][1])
	}
	for i := range people {
		people[i] = f(people[i])
	}

	e, op, inv := func() int { return 0 }, func(a, b int) int { return a + b }, func(a int) int { return -a }
	diff := NewDiff1D(e, op, inv)
	diff.Init(int32(len(origin)), func(i int32) int { return 0 })
	for _, flower := range flowers {
		diff.Add(int32(flower[0]), int32(flower[1]+1), 1)
	}
	res := make([]int, len(people))
	for i, p := range people {
		res[i] = diff.Get(int32(p))
	}
	return res
}

type Diff1D[E any] struct {
	dirty bool
	n     int32
	diff  []E
	data  []E
	e     func() E
	op    func(a, b E) E
	inv   func(a E) E
}

func NewDiff1D[E any](e func() E, op func(a, b E) E, inv func(a E) E) *Diff1D[E] {
	return &Diff1D[E]{e: e, op: op, inv: inv}
}

func (d *Diff1D[E]) Init(n int32, f func(i int32) E) {
	diff := make([]E, n)
	data := make([]E, n)
	for i := int32(0); i < n; i++ {
		diff[i] = d.e()
		data[i] = f(i)
	}
	d.dirty = false
	d.n = n
	d.diff = diff
	d.data = data
}

func (d *Diff1D[E]) Add(start, end int32, v E) {
	if start < 0 {
		start = 0
	}
	if end > d.n {
		end = d.n
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] = d.op(d.diff[start], v)
	if end < d.n {
		d.diff[end] = d.op(d.diff[end], d.inv(v))
	}
}

func (d *Diff1D[E]) Get(index int32) E {
	if d.dirty {
		d.rebuild()
	}
	return d.data[index]
}

func (d *Diff1D[E]) rebuild() {
	if !d.dirty {
		return
	}
	e, op, data, diff := d.e, d.op, d.data, d.diff
	cur := e()
	for i := int32(0); i < d.n; i++ {
		cur = op(cur, diff[i])
		data[i] = op(data[i], cur)
		diff[i] = e()
	}
	d.dirty = false
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
}
