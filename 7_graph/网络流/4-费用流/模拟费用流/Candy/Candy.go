// 环上糖果分配问题
// 在一个环上，有n个点，编号为0,1,...,n-1。
// 每个点上有一些糖果，xi表示点x上有xi个糖果。
// 你可以在相邻的两个点之间传递糖果，每传递一个糖果的代价为1。
// 你需要保证每个点上有yi个糖果，同时x0+x1+...+x_{n-1}=y0+y1+...+y_{n-1}。

package main

import (
	"fmt"
	"sort"
)

func main() {
	cap := NewCandyAssignProblem(10, 20)
	for i := 0; i < 10; i++ {
		cap.RequestOn(i, i, 10-i)
	}
	fmt.Println(cap.Solve()) // 20
	for i := 0; i < 10; i++ {
		fmt.Println(cap.DeliverBetween(i))
	}
}

type CandyAssignProblem struct {
	n            int
	addedCandies []*Candy
	candies      []*Candy
	candieCnt    int32
	minimumCost  int
}

func NewCandyAssignProblem(n int, candyCount int32) *CandyAssignProblem {
	return &CandyAssignProblem{
		n:            n,
		addedCandies: make([]*Candy, 0, candyCount),
	}
}

// 点i上有x个糖果，最后需要y个糖果.
func (cap *CandyAssignProblem) RequestOn(i, x, y int) {
	candy := &Candy{
		location: i,
		x:        x,
		y:        y,
	}
	cap.addedCandies = append(cap.addedCandies, candy)
}

// 点i上有x个糖果.
func (cap *CandyAssignProblem) AssignCandyOn(i, x int) {
	cap.RequestOn(i, x, 0)
}

// 点i上最后需要y个糖果.
func (cap *CandyAssignProblem) RequireCandyOn(i, y int) {
	cap.RequestOn(i, 0, y)
}

func (cap *CandyAssignProblem) Solve() int {
	if len(cap.addedCandies) == 0 {
		cap.addedCandies = append(cap.addedCandies, &Candy{})
	}
	cap.candies = append(cap.candies[:0:0], cap.addedCandies...)
	sort.Slice(cap.candies, func(i, j int) bool {
		return cap.candies[i].location < cap.candies[j].location
	})
	cap.candieCnt = 0
	for i := 1; i < len(cap.candies); i++ {
		if cap.candies[i].location == cap.candies[cap.candieCnt].location {
			cap.candies[cap.candieCnt].x += cap.candies[i].x
			cap.candies[cap.candieCnt].y += cap.candies[i].y
		} else {
			cap.candieCnt++
			cap.candies[cap.candieCnt] = cap.candies[i]
		}
	}
	cap.candieCnt++

	for i := int32(0); i < cap.candieCnt-1; i++ {
		cap.candies[i].w = cap.candies[i+1].location - cap.candies[i].location
	}
	cap.candies[cap.candieCnt-1].w = cap.n + cap.candies[0].location - cap.candies[cap.candieCnt-1].location
	for i := int32(1); i < cap.candieCnt; i++ {
		cap.candies[i].a = cap.candies[i-1].a + cap.candies[i].x - cap.candies[i].y
	}

	sortedByA := make([]*Candy, cap.candieCnt)
	for i := int32(0); i < cap.candieCnt; i++ {
		sortedByA[i] = &Candy{}
	}
	copy(sortedByA, cap.candies)
	sort.Slice(sortedByA, func(i, j int) bool {
		return sortedByA[i].a < sortedByA[j].a
	})
	var prefix int
	half := (cap.n + 1) / 2
	for i := int32(0); i < cap.candieCnt; i++ {
		prefix += sortedByA[i].w
		if prefix >= half {
			cap.candies[0].a = -sortedByA[i].a
			break
		}
	}

	for i := int32(1); i < cap.candieCnt; i++ {
		cap.candies[i].a += cap.candies[0].a
	}
	for i := int32(0); i < cap.candieCnt; i++ {
		cap.minimumCost += abs(cap.candies[i].a) * cap.candies[i].w
	}
	return cap.minimumCost
}

func (cap *CandyAssignProblem) MinimumCost() int {
	return cap.minimumCost
}

// 从i到i+1传递多少糖果，返回值可能为负数。
func (cap *CandyAssignProblem) DeliverBetween(i int) int {
	index := cap.lowerBound(cap.candies, 0, int(cap.candieCnt-1), i)
	if index < 0 {
		index = int(cap.candieCnt - 1)
	}
	return cap.candies[index].a
}

func (cap *CandyAssignProblem) lowerBound(arr []*Candy, l, r int, v int) int {
	for l < r {
		mid := (l + r) >> 1
		if arr[mid].location >= v {
			r = mid
		} else {
			l = mid + 1
		}
	}
	if arr[l].location < v {
		l++
	}
	return l
}

type Candy struct {
	location int
	x, y     int
	a, w     int
}

func abs(a int) int {
	if a < 0 {
		return -a
	}
	return a
}
