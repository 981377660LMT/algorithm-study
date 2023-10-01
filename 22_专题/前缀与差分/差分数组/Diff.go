package main

import "sort"

// 2251. 花期内花的数目
// https://leetcode.cn/problems/number-of-flowers-in-full-bloom/description/
func fullBloomFlowers(flowers [][]int, people []int) []int {
	diff := NewDiffMap()
	for _, flower := range flowers {
		diff.Add(flower[0], flower[1]+1, 1)
	}
	res := make([]int, len(people))
	for i, p := range people {
		res[i] = diff.Get(p)
	}
	return res
}

type DiffArray struct {
	diff  []int
	dirty bool
}

func NewDiffArray(n int) *DiffArray {
	return &DiffArray{
		diff: make([]int, n+1),
	}
}

func (d *DiffArray) Add(start, end, delta int) {
	if start < 0 {
		start = 0
	}
	if end >= len(d.diff) {
		end = len(d.diff) - 1
	}
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffArray) Build() {
	if d.dirty {
		for i := 1; i < len(d.diff); i++ {
			d.diff[i] += d.diff[i-1]
		}
		d.dirty = false
	}
}

func (d *DiffArray) Get(pos int) int {
	d.Build()
	return d.diff[pos]
}

func (d *DiffArray) GetAll() []int {
	d.Build()
	return d.diff[:len(d.diff)-1]
}

type DiffMap struct {
	diff       map[int]int
	sortedKeys []int
	preSum     []int
	dirty      bool
}

func NewDiffMap() *DiffMap {
	return &DiffMap{
		diff: make(map[int]int),
	}
}

func (d *DiffMap) Add(start, end, delta int) {
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffMap) Build() {
	if d.dirty {
		d.sortedKeys = make([]int, 0, len(d.diff))
		for key := range d.diff {
			d.sortedKeys = append(d.sortedKeys, key)
		}
		sort.Ints(d.sortedKeys)
		d.preSum = make([]int, len(d.sortedKeys)+1)
		for i, key := range d.sortedKeys {
			d.preSum[i+1] = d.preSum[i] + d.diff[key]
		}
		d.dirty = false
	}
}

func (d *DiffMap) Get(pos int) int {
	d.Build()
	return d.preSum[sort.SearchInts(d.sortedKeys, pos+1)]
}
