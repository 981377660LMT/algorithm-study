package main

import "sort"

type RangeFreqQuery struct {
	mp map[int][]int
}

func Constructor(arr []int) RangeFreqQuery {
	mp := map[int][]int{}
	for i, v := range arr {
		mp[v] = append(mp[v], i)
	}
	return RangeFreqQuery{mp: mp}
}

// 查询[left,right]区间内等于value的元素个数.
func (this *RangeFreqQuery) Query(left int, right int, value int) int {
	pos := this.mp[value]
	return sort.SearchInts(pos, right+1) - sort.SearchInts(pos, left)
}

/**
 * Your RangeFreqQuery object will be instantiated and called as such:
 * obj := Constructor(arr);
 * param_1 := obj.Query(left,right,value);
 */
