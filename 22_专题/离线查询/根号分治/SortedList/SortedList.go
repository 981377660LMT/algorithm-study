package main

import "sort"

func main() {

}

const _LOAD int = 1 << 9

type S [2]int

// 使用分块+树状数组维护的有序序列.
type SortedList struct {
	less          func(a, b S) bool
	size          int
	blocks        [][]S
	blockLens     []int
	mins          []S
	tree          []int
	shouldRebuild bool
}

func NewSortedList(less func(a, b S) bool, elements ...S) *SortedList {
	res := &SortedList{less: less}
	sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
	n := len(elements)
	lists := [][]S{}
	for i := 0; i < n; i += _LOAD {
		lists = append(lists, elements[i:min(i+_LOAD, n)])
	}
	listLens := make([]int, len(lists))
	mins := make([]S, len(lists))
	for i, cur := range lists {
		listLens[i] = len(cur)
		mins[i] = cur[0]
	}
	res.size = n
	res.blocks = lists
	res.blockLens = listLens
	res.mins = mins
	res.shouldRebuild = true
	return res
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
