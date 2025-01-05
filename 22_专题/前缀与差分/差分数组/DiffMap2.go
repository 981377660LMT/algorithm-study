package main

import (
	"fmt"
	"math/rand"
	"slices"
	"sort"
)

const INF int = 1e18

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

// 3413. 收集连续 K 个袋子可以获得的最多硬币数量
// https://leetcode.cn/problems/maximum-coins-from-k-consecutive-bags/description/
func maximumCoins(coins [][]int, k int) int64 {
	seg := NewDiffMap()
	for _, coin := range coins {
		l, r, w := coin[0], coin[1], coin[2]
		seg.AddRange(l, r+1, w)
	}

	res := 0
	for _, coin := range coins {
		l, r := coin[0], coin[1]
		res = max(res, seg.GetRange(l, l+k))
		res = max(res, seg.GetRange(r-k+1, r+1))
	}
	return int64(res)
}

// 支持区间和查询的差分Map.
type DiffMap2 struct {
	diff  map[int]int
	pos   []int
	sum0  []int
	sum1  []int
	dirty bool
}

func NewDiffMap() *DiffMap2 {
	return &DiffMap2{diff: make(map[int]int)}
}

// [start, end) += delta.
func (d *DiffMap2) AddRange(start, end, delta int) {
	if start >= end {
		return
	}
	d.dirty = true
	d.diff[start] += delta
	d.diff[end] -= delta
}

func (d *DiffMap2) Get(pos int) int {
	d.build()
	i := sort.SearchInts(d.pos, pos+1)
	if i == 0 {
		return 0
	}
	return d.sum0[i-1]
}

func (d *DiffMap2) GetRange(start, end int) int {
	if start >= end {
		return 0
	}
	d.build()
	return d.presum(end) - d.presum(start)
}

func (d *DiffMap2) build() {
	if !d.dirty {
		return
	}
	d.dirty = false

	pos := make([]int, 0, len(d.diff))
	for p := range d.diff {
		pos = append(pos, p)
	}
	if len(pos) == 0 {
		return
	}

	slices.Sort(pos)
	d.pos = pos
	d.sum0 = make([]int, len(pos))
	d.sum1 = make([]int, len(pos))

	pre := pos[0]
	s0 := d.diff[pre]
	s1 := 0
	d.sum0[0] = s0
	for i := 1; i < len(pos); i++ {
		cur := pos[i]

		s1 += (cur - pre) * s0
		d.sum1[i] = s1

		s0 += d.diff[cur]
		d.sum0[i] = s0

		pre = cur
	}
}

func (d *DiffMap2) presum(v int) int {
	if len(d.pos) == 0 {
		return 0
	}
	if v <= d.pos[0] {
		return 0
	}
	if v >= d.pos[len(d.pos)-1] {
		return d.sum1[len(d.pos)-1]
	}
	i := sort.SearchInts(d.pos, v+1)
	if i == 0 {
		return 0
	}
	res := d.sum1[i-1]
	width := v - d.pos[i-1]
	res += width * d.sum0[i-1]
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

func check() {
	type MockedMap map[int]int

	addRange := func(m MockedMap, start, end, delta int) {
		for i := start; i < end; i++ {
			m[i] += delta
		}
	}
	_ = addRange

	get := func(m MockedMap, pos int) int {
		return m[pos]
	}
	_ = get

	getRange := func(m MockedMap, start, end int) int {
		res := 0
		for i := start; i < end; i++ {
			res += m[i]
		}
		return res
	}
	_ = getRange

	run := func() {
		mp1 := make(MockedMap)
		mp2 := NewDiffMap()

		for i := 0; i < 1000; i++ {
			start := rand.Intn(100)
			end := rand.Intn(100)
			delta := rand.Intn(100)
			addRange(mp1, start, end, delta)
			mp2.AddRange(start, end, delta)
		}

		for i := 0; i < 1000; i++ {
			pos := rand.Intn(100)
			res1 := get(mp1, pos)
			res2 := mp2.Get(pos)
			if res1 != res2 {
				fmt.Println("error", pos, res1, res2)
				panic("")
			}
		}

		for i := 0; i < 1000; i++ {
			start := rand.Intn(100)
			end := rand.Intn(100)
			if getRange(mp1, start, end) != mp2.GetRange(start, end) {
				fmt.Println("error", start, end, getRange(mp1, start, end), mp2.GetRange(start, end))
				panic("")
			}
		}
	}

	for i := 0; i < 10000; i++ {
		run()
	}
	fmt.Println("done")
}
