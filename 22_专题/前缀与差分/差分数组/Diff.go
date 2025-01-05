package main

import (
	"fmt"
	"slices"
	"sort"
)

// 2251. 花期内花的数目
// https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/
func fullBloomFlowers(flowers [][]int, people []int) []int {
	diff := NewDiffMap()
	for _, flower := range flowers {
		diff.AddRange(flower[0], flower[1]+1, 1)
	}
	res := make([]int, len(people))
	for i, p := range people {
		res[i] = diff.Get(p)
	}
	return res
}

func main() {
	{

		D := NewDiffArray(10)
		D.AddRange(1, 3, 1)
		D.AddRange(2, 10, 1)
		fmt.Println(D.GetAll())
	}

	{
		D := NewDiffMap()
		D.AddRange(1, 3, 1)
		D.AddRange(2, 10, 1)
		for i := 0; i < 10; i++ {
			fmt.Println(i, D.Get(i))
		}
	}
}

type DiffArray struct {
	n     int
	diff  []int
	dirty bool
}

// [0, n)
func NewDiffArray(n int) *DiffArray {
	return &DiffArray{
		n:    n,
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) AddRange(start, end, delta int) {
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
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Get(pos int) int {
	if pos < 0 || pos >= d.n {
		return 0
	}
	d.build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.build()
	return d.diff[:len(d.diff)-1]
}

func (d *DiffArray) build() {
	if !d.dirty {
		return
	}
	d.dirty = false
	for i := 1; i < len(d.diff); i++ {
		d.diff[i] += d.diff[i-1]
	}
}

// DiffArraySparse.
type DiffMap struct {
	diff       map[int]int
	sortedKeys []int
	preSum     []int
	dirty      bool
}

func NewDiffMap() *DiffMap {
	return &DiffMap{
		diff:  make(map[int]int),
		dirty: true,
	}
}

func (d *DiffMap) AddRange(start, end, delta int) {
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffMap) Build() {
	if !d.dirty {
		return
	}
	d.dirty = false

	d.sortedKeys = make([]int, 0, len(d.diff))
	for key := range d.diff {
		d.sortedKeys = append(d.sortedKeys, key)
	}
	slices.Sort(d.sortedKeys)
	d.preSum = make([]int, len(d.sortedKeys)+1)
	for i, key := range d.sortedKeys {
		d.preSum[i+1] = d.preSum[i] + d.diff[key]
	}
}

func (d *DiffMap) Get(pos int) int {
	d.Build()
	return d.preSum[sort.SearchInts(d.sortedKeys, pos+1)]
}
