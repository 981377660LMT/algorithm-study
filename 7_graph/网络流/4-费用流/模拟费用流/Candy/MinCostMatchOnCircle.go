// 环上老鼠进洞问题
// 环上最小费用流
//
// 环上有n个点 0, 1, 2, ... , n - 1
// 点i和点i+1相邻 (0和n-1也相邻)
// 有一些人和一些房子(人和房子数量相等)
// 人可以在相邻点之间移动，每移动一步花费1
// 现在你需要建立人和房子之间的匹配，使得总移动费用最小

package main

import (
	"fmt"
	"sort"
)

func main() {
	n := 5
	people := []int{0, 1, 3}
	houses := []int{1, 2, 4}
	M := NewMinCostMatchOnCircle(n, people, houses)
	fmt.Println(M.GetMinimumCost()) // 9
	for i := 0; i < len(people); i++ {
		fmt.Println(M.GetHouseOf(int32(i)))
	}
}

type MinCostMatchOnCircle struct {
	matching  []int32
	minCost   int
	peopleMap map[int]*Vec32
	houseMap  map[int]*Vec32
	candy     *CandyAssignProblem
	pending   []int32
}

func NewMinCostMatchOnCircle(n int, people, houses []int) *MinCostMatchOnCircle {
	peopleCount := make([]int, len(people))
	houseCount := make([]int, len(houses))
	for i := 0; i < len(people); i++ {
		peopleCount[i] = 1
		houseCount[i] = 1
	}
	return NewMinCostMatchOnCircleWithWeights(n, people, houses, peopleCount, houseCount)
}

func NewMinCostMatchOnCircleWithWeights(n int, people, houses []int, peopleCount, houseCount []int) *MinCostMatchOnCircle {
	if n <= 0 || len(people) != len(houses) {
		panic("invalid input")
	}
	m := int32(len(people))
	candy := NewCandyAssignProblem(n, m*2)
	peopleMap := make(map[int]*Vec32, m)
	houseMap := make(map[int]*Vec32, m)
	matching := make([]int32, m)
	pending := make([]int32, 0, m)
	for i := int32(0); i < m; i++ {
		p := people[i]
		if _, ok := peopleMap[p]; !ok {
			peopleMap[p] = &Vec32{}
		}
		peopleMap[p].Push(i)
		h := houses[i]
		if _, ok := houseMap[h]; !ok {
			houseMap[h] = &Vec32{}
		}
		houseMap[h].Push(i)
	}
	for i := int32(0); i < m; i++ {
		candy.AssignCandyOn(people[i], peopleCount[i])
	}
	for i := int32(0); i < m; i++ {
		candy.RequireCandyOn(houses[i], houseCount[i])
	}
	candy.Solve()
	minCost := candy.MinimumCost()
	res := &MinCostMatchOnCircle{
		matching:  matching,
		minCost:   minCost,
		peopleMap: peopleMap,
		houseMap:  houseMap,
		candy:     candy,
		pending:   pending,
	}
	for i := int32(0); i < candy.candieCnt; i++ {
		last := (i - 1) % candy.candieCnt
		if last < 0 {
			last += candy.candieCnt
		}
		if candy.candies[i].a >= 0 && candy.candies[last].a <= 0 {
			res.buildMatching(i)
		}
	}
	return res
}

func (mcm *MinCostMatchOnCircle) GetMinimumCost() int {
	return mcm.minCost
}

// person 对应的房子.
func (mcm *MinCostMatchOnCircle) GetHouseOf(person int32) int32 {
	return mcm.matching[person]
}

func (mcm *MinCostMatchOnCircle) buildMatching(i int32) {
	people := mcm.peopleMap[mcm.candy.candies[i].location]
	houses := mcm.houseMap[mcm.candy.candies[i].location]
	if people != nil && houses != nil {
		for people.Len() > 0 && houses.Len() > 0 {
			mcm.matching[people.Pop()] = houses.Pop()
		}
	}
	if houses != nil {
		for houses.Len() > 0 && len(mcm.pending) > 0 {
			mcm.matching[mcm.pending[len(mcm.pending)-1]] = houses.Pop()
			mcm.pending = mcm.pending[:len(mcm.pending)-1]
		}
	}
	if mcm.candy.candies[i].a > 0 {
		for len(mcm.pending) < mcm.candy.candies[i].a {
			mcm.pending = append(mcm.pending, people.Pop())
		}
		mcm.buildMatching((i + 1) % mcm.candy.candieCnt)
	}
	last := (i - 1) % mcm.candy.candieCnt
	if last < 0 {
		last += mcm.candy.candieCnt
	}
	if mcm.candy.candies[last].a < 0 {
		for len(mcm.pending) < -mcm.candy.candies[last].a {
			mcm.pending = append(mcm.pending, people.Pop())
		}
		mcm.buildMatching(last)
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

type Vec32 struct{ data []int32 }

func (l *Vec32) Push(v int32) { l.data = append(l.data, v) }
func (l *Vec32) Pop() int32 {
	v := l.data[len(l.data)-1]
	l.data = l.data[:len(l.data)-1]
	return v
}
func (l *Vec32) At(i int32) int32 { return l.data[i] }
func (l *Vec32) Len() int32       { return int32(len(l.data)) }
