package main

import (
	"fmt"
	"math/rand"
)

func main() {
	//  test with bf
	nums := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}
	tree := NewRangeAddRangeMin(nums)

	mins := func(nums ...int) int {
		res := nums[0]
		for _, num := range nums {
			if num < res {
				res = num
			}
		}
		return res
	}

	for i := 0; i < 1000; i++ {
		start, end, add := rand.Intn(10), rand.Intn(10), rand.Intn(10)
		if start > end {
			start, end = end, start
		}
		tree.Update(start, end, add)
		for i := start; i < end; i++ {
			nums[i] += add
		}
		for i := 0; i < 10; i++ {
			for j := i + 1; j <= 10; j++ {
				if tree.Query(i, j) != mins(nums[i:j]...) {
					panic("error")
				}
			}
		}
	}

	fmt.Println("test with bf success")
}

const INF int = 1e18

type RangeAddRangeMin struct {
	n, head            int
	rangeMin, rangeAdd []int
}

func NewRangeAddRangeMin(nums []int) *RangeAddRangeMin {
	res := &RangeAddRangeMin{}
	res.n = len(nums)
	res.head = 1
	for res.head < res.n {
		res.head <<= 1
	}
	res.rangeMin, res.rangeAdd = make([]int, res.head*2), make([]int, res.head*2)
	for i := range res.rangeMin {
		res.rangeMin[i] = INF
	}
	copy(res.rangeMin[res.head:], nums)
	for i := res.head - 1; i > 0; i-- {
		res.merge(i)
	}
	return res
}

func (r *RangeAddRangeMin) Update(start, end, add int) {
	r.update(start, end, 1, 0, r.head, add)
}
func (r *RangeAddRangeMin) Query(start, end int) int {
	return r.query(start, end, 1, 0, r.head)
}

func (r *RangeAddRangeMin) Get(pos int) int { return r.Query(pos, pos+1) }

func (r *RangeAddRangeMin) query(start, end, pos, left, right int) int {
	if right <= start || end <= left {
		return INF
	}
	if start <= left && right <= end {
		return r.rangeMin[pos] + r.rangeAdd[pos]
	}
	return min(r.query(start, end, pos*2, left, (left+right)/2), r.query(start, end, pos*2+1, (left+right)/2, right)) + r.rangeAdd[pos]
}

func (r *RangeAddRangeMin) update(start, end, pos, left, right, add int) {
	if right <= start || end <= left {
		return
	}
	if start <= left && right <= end {
		r.rangeAdd[pos] += add
		return
	}
	r.update(start, end, pos*2, left, (left+right)/2, add)
	r.update(start, end, pos*2+1, (left+right)/2, right, add)
	r.merge(pos)
}

func (r *RangeAddRangeMin) merge(pos int) {
	r.rangeMin[pos] = min(r.rangeMin[pos*2]+r.rangeAdd[pos*2], r.rangeMin[pos*2+1]+r.rangeAdd[pos*2+1])
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
