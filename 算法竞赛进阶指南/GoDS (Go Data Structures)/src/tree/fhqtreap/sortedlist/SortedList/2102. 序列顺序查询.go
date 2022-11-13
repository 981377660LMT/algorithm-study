// https://leetcode.cn/problems/sequentially-ordinal-rank-tracker/

package main

import "strings"

type SORTracker struct {
	sl    *SortedList
	index int
}

type spot struct {
	name  string
	score int
}

func Constructor() SORTracker {
	return SORTracker{
		// !景点评分 越高 ，这个景点越好。如果有两个景点的评分一样，那么 字典序较小 的景点更好。
		sl: NewSortedList(func(a, b interface{}) int {
			sa, sb := a.(spot), b.(spot)
			if sa.score != sb.score {
				return sb.score - sa.score
			}
			return strings.Compare(sa.name, sb.name)
		}, 16),
	}
}

func (this *SORTracker) Add(name string, score int) {
	this.sl.Add(spot{name, score})
}

func (this *SORTracker) Get() string {
	res := this.sl.At(this.index)
	this.index++
	return res.(spot).name
}
