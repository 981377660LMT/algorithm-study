/* eslint-disable @typescript-eslint/no-non-null-assertion */

// https://leetcode.cn/problems/find-closest-lcci/
// https://www.luogu.com.cn/blog/236866/sol-p5397 的无修改版

package main

import "math"

// https://leetcode.cn/problems/find-closest-lcci/description/
func findClosest(words []string, word1 string, word2 string) int {
	finder := NewClosestFinder(words, -1)
	return finder.FindClosest(word1, word2)
}

type V = string

// 给定一个数组，每次查询给定`x`和`y`，查询数组中`x`和`y`的最近距离.
// `O(nsqrt(n))`预处理，`O(sqrt(n))`查询.
type ClosestFinder struct {
	_threshold      int
	_ids            []int
	_mp             [][]int
	_largeResRecord map[int][]int
	_valueToId      map[V]int
}

// threshold: 当元素出现的次数超过这个阈值时，会对此元素进行时空复杂度`O(n)`的预处理.
func NewClosestFinder(arr []V, threshold int) *ClosestFinder {
	if threshold == -1 {
		threshold = 4 * (1 + int(math.Sqrt(float64(len(arr)))))
	}
	res := &ClosestFinder{
		_threshold:      threshold,
		_ids:            make([]int, len(arr)),
		_mp:             [][]int{},
		_largeResRecord: map[int][]int{},
		_valueToId:      map[V]int{},
	}
	for i, v := range arr {
		id := res._getId(v)
		res._ids[i] = id
		if len(res._mp) <= id {
			res._mp = append(res._mp, []int{})
		}
		res._mp[id] = append(res._mp[id], i)
	}
	for id := range res._mp {
		if res._isLarge(id) {
			res._largeResRecord[id] = res._buildLargeRes(id)
		}
	}
	return res
}

// 查询`x`和`y`的最近距离.
// 如果`x`和`y`相同，返回`0`.
// 如果`x`或`y`不存在，返回`-1`.
func (f *ClosestFinder) FindClosest(x, y V) int {
	if x == y {
		return 0
	}
	if _, ok := f._valueToId[x]; !ok {
		return -1
	}
	if _, ok := f._valueToId[y]; !ok {
		return -1
	}
	id1, id2 := f._getId(x), f._getId(y)
	if f._isLarge(id1) {
		return f._largeResRecord[id1][id2]
	}
	if f._isLarge(id2) {
		return f._largeResRecord[id2][id1]
	}
	pos1, pos2 := f._mp[id1], f._mp[id2]
	i, j, res := 0, 0, len(f._ids)
	for i < len(pos1) && j < len(pos2) {
		res = min(res, abs(pos1[i]-pos2[j]))
		if pos1[i] < pos2[j] {
			i++
		} else {
			j++
		}
	}
	return res
}

func (f *ClosestFinder) _getId(value V) int {
	res, ok := f._valueToId[value]
	if ok {
		return res
	}
	id := len(f._valueToId)
	f._valueToId[value] = id
	return id
}

func (f *ClosestFinder) _isLarge(id int) bool {
	return len(f._mp[id]) > f._threshold
}

func (f *ClosestFinder) _buildLargeRes(id int) []int {
	n := len(f._ids)
	res := make([]int, len(f._valueToId))
	for i := range res {
		res[i] = n
	}
	dist := n
	for i := 0; i < n; i++ {
		cur := f._ids[i]
		if cur == id {
			dist = 0
		} else {
			dist++
			res[cur] = min(res[cur], dist)
		}
	}
	dist = n
	for i := n - 1; i >= 0; i-- {
		cur := f._ids[i]
		if cur == id {
			dist = 0
		} else {
			dist++
			res[cur] = min(res[cur], dist)
		}
	}
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

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
