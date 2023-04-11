// https://noshi91.github.io/Library/data_structure/range_mode_query.cpp

package main

import (
	"sort"
)

// https://leetcode.cn/problems/online-majority-element-in-subarray/
type MajorityChecker struct {
	RMQ *RangeModeQuery
}

func Constructor(arr []int) MajorityChecker {
	return MajorityChecker{NewRangeModeQuery(arr)}
}

func (this *MajorityChecker) Query(left int, right int, threshold int) int {
	mode, freq := this.RMQ.Query(left, right+1)
	if freq >= threshold {
		return mode
	}
	return -1
}

// 在线查询区间众数(出现次数最多的数和出现次数).
type RangeModeQuery struct {
	value, rank []int   // 值和id(0-based) -> bNode
	mode, freq  [][]int // 值和出现次数 -> sNode

	qs [][]int
	t  int

	sorted []int
}

// O(nsqrt(n))构建.
func NewRangeModeQuery(nums []int) *RangeModeQuery {
	n := len(nums)
	sorted, rank_ := sortedSet(nums) // 离散化
	value := make([]int, n)
	for i, e := range nums {
		value[i] = rank_[e]
	}

	res := &RangeModeQuery{}
	t := 1
	for t*t < n {
		t++
	}

	rank := make([]int, n)
	qs := make([][]int, n)
	for i := 0; i < n; i++ {
		e := value[i]
		rank[i] = len(qs[e])
		qs[e] = append(qs[e], i)
	}

	bc := n/t + 1
	mode, freq := make([][]int, bc), make([][]int, bc)
	for i := 0; i < bc; i++ {
		mode[i] = make([]int, bc)
		freq[i] = make([]int, bc)
	}

	for f := 0; f*t <= n; f++ {
		freq_ := make([]int, n)
		curMode, curFreq := 0, 0
		for l := f + 1; l*t <= n; l++ {
			for i := (l - 1) * t; i != l*t; i++ {
				e := value[i]
				freq_[e]++
				if freq_[e] > curFreq {
					curMode, curFreq = e, freq_[e]
				}
			}
			mode[f][l] = curMode
			freq[f][l] = curFreq
		}
	}

	res.value = value
	res.rank = rank
	res.mode = mode
	res.freq = freq
	res.qs = qs
	res.t = t
	res.sorted = sorted

	return res
}

// O(sqrt(n))查询区间 [start, end) 中出现次数最多的数mode, 以及出现的次数freq.
//  0 <= start < end <= len(nums)
func (rmq *RangeModeQuery) Query(start, end int) (mode, freq int) {
	if start >= end {
		panic("start>=end")
	}
	if start < 0 {
		start = 0
	}
	if end > len(rmq.value) {
		end = len(rmq.value)
	}

	T := rmq.t
	bf := start/T + 1
	bl := end / T
	if bf >= bl {
		resMode, resFreq := 0, 0
		for i := start; i < end; i++ {
			rank, value := rmq.rank[i], rmq.value[i]
			v := rmq.qs[value]
			lenV := len(v)
			for {
				idx := rank + resFreq
				if idx >= lenV || v[idx] >= end {
					break
				}
				resMode = value
				resFreq++
			}
		}
		return rmq.sorted[resMode], resFreq
	}

	resMode, resFreq := rmq.mode[bf][bl], rmq.freq[bf][bl]
	for i := start; i < bf*T; i++ {
		rank, value := rmq.rank[i], rmq.value[i]
		v := rmq.qs[value]
		lenV := len(v)
		for {
			idx := rank + resFreq
			if idx >= lenV || v[idx] >= end {
				break
			}
			resMode = value
			resFreq++
		}
	}

	for i := bl * T; i < end; i++ {
		rank, value := rmq.rank[i], rmq.value[i]
		v := rmq.qs[value]
		lenV := len(v)
		for {
			idx := rank - resFreq
			if idx < 0 || idx >= lenV || v[idx] < start {
				break
			}
			resMode = value
			resFreq++
		}
	}

	return rmq.sorted[resMode], resFreq
}

func sortedSet(xs []int) (sorted []int, rank map[int]int) {
	set := make(map[int]struct{}, len(xs))
	for _, v := range xs {
		set[v] = struct{}{}
	}
	sorted = make([]int, 0, len(set))
	for k := range set {
		sorted = append(sorted, k)
	}
	sort.Ints(sorted)
	rank = make(map[int]int, len(sorted))
	for i, v := range sorted {
		rank[v] = i
	}
	return
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
