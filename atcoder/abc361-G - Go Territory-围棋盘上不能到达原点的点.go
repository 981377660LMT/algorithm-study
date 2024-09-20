package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

const INF int = 1e18

// G - Go Territory (围棋盘上不能到达原点的点，超大的棋盘上围棋围住的点的数量)
// 扫描线+并查集
// https://atcoder.jp/contests/abc361/tasks/abc361_g
//
// 二维平面，有障碍物，可以上下左右走。
// 问有多少个点，不可以走到(-1,-1).
//
// 1.每个点按照x坐标分组，然后按照y坐标排序；
// 2.遍历x，对x相邻的两个空缺的y区间进行合并；
// !3.合并后，以空缺区间为单位，统计哪些区间不与0相连。
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	stones := make([][2]int, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &stones[i][0], &stones[i][1])
	}

	// !所有需要使用到的x坐标(包括间隙和边界)，离散化
	allX := func() []int {
		res := make([]int, 0, n)
		for i := int32(0); i < n; i++ {
			res = append(res, stones[i][0])
		}
		sort.Ints(res)
		toAdd := []int{-INF, INF}
		EnumerateConsecutiveIntervals(
			int32(len(res)), func(i int32) int { return res[i] },
			func(min, max int, isIn bool) {
				if !isIn {
					toAdd = append(toAdd, min)
				}
			},
		)
		return append(toAdd, res...)
	}()
	_, origin := Discretize(allX)
	size := len(origin)
	for i := int32(0); i < n; i++ {
		stones[i][0] = sort.SearchInts(origin, stones[i][0])
	}

	// !将点按照x坐标分组
	xToYs := make([][]int, size)
	for i := int32(0); i < n; i++ {
		x, y := stones[i][0], stones[i][1]
		xToYs[x] = append(xToYs[x], y)
	}
	for _, ys := range xToYs {
		sort.Ints(ys)
	}

	// !扫描线，维护每个x坐标的空缺的的y坐标区间(矩形)
	type interval struct {
		id     int32
		y1, y2 int
	}
	part := int32(0)
	newInterval := func(y1, y2 int) interval {
		id := part
		part++
		return interval{id: id, y1: y1, y2: y2}
	}
	xToIntervals := make(map[int][]interval)
	for x, ys := range xToYs {
		if len(xToYs[x]) == 0 {
			xToIntervals[x] = []interval{newInterval(-INF, INF)}
		} else {
			xToIntervals[x] = append(xToIntervals[x], newInterval(-INF, ys[0]-1))
			EnumerateConsecutiveIntervals(
				int32(len(ys)), func(i int32) int { return ys[i] },
				func(min, max int, isIn bool) {
					if !isIn {
						xToIntervals[x] = append(xToIntervals[x], newInterval(min, max))
					}
				},
			)
			xToIntervals[x] = append(xToIntervals[x], newInterval(ys[len(ys)-1]+1, INF))
		}
	}

	// !合并相邻的空缺区间
	uf := NewUnionFindArraySimple32(part)
	merge := func(inter1, inter2 []interval) {
		EnumerateIntervalsIntersection(
			len(inter1), func(i int) (int, int) { return inter1[i].y1, inter1[i].y2 },
			len(inter2), func(i int) (int, int) { return inter2[i].y1, inter2[i].y2 },
			func(left, right, i, j int) {
				uf.Union(inter1[i].id, inter2[j].id, nil)
			},
		)
	}
	for i := 0; i < len(allX)-1; i++ {
		merge(xToIntervals[i], xToIntervals[i+1])
	}

	res := 0
	for _, inters := range xToIntervals {
		for _, inter := range inters {
			if uf.Find(inter.id) != uf.Find(0) {
				res += inter.y2 - inter.y1 + 1
			}
		}
	}
	fmt.Fprintln(out, res)
}

type UnionFindArraySimple32 struct {
	Part int32
	n    int32
	data []int32
}

func NewUnionFindArraySimple32(n int32) *UnionFindArraySimple32 {
	data := make([]int32, n)
	for i := int32(0); i < n; i++ {
		data[i] = -1
	}
	return &UnionFindArraySimple32{Part: n, n: n, data: data}
}

func (u *UnionFindArraySimple32) Union(key1, key2 int32, beforeMerge func(big, small int32)) bool {
	root1, root2 := u.Find(key1), u.Find(key2)
	if root1 == root2 {
		return false
	}
	if u.data[root1] > u.data[root2] {
		root1, root2 = root2, root1
	}
	if beforeMerge != nil {
		beforeMerge(root1, root2)
	}
	u.data[root1] += u.data[root2]
	u.data[root2] = root1
	u.Part--
	return true
}

func (u *UnionFindArraySimple32) Find(key int32) int32 {
	root := key
	for u.data[root] >= 0 {
		root = u.data[root]
	}
	for key != root {
		key, u.data[key] = u.data[key], root
	}
	return root
}

func (u *UnionFindArraySimple32) Size(key int32) int32 {
	return -u.data[u.Find(key)]
}

// 遍历连续区间.
func EnumerateConsecutiveIntervals(
	n int32, supplier func(i int32) int,
	consumer func(min, max int, isIn bool),
) {
	if n == 0 {
		return
	}
	i := int32(0)
	for i < n {
		start := i
		for i < n-1 && supplier(i)+1 == supplier(i+1) {
			i++
		}
		consumer(supplier(start), supplier(i), true)
		if i+1 < n {
			consumer(supplier(i)+1, supplier(i+1)-1, false)
		}
		i++
	}
}

// 有序区间列表交集(EnumerateIntersection).
func EnumerateIntervalsIntersection(
	n1 int, f1 func(int) (int, int), n2 int, f2 func(int) (int, int),
	f func(left, right, i, j int),
) {
	i, j := 0, 0
	for i < n1 && j < n2 {
		s1, e1 := f1(i)
		s2, e2 := f2(j)
		if (s1 <= e2 && e2 <= e1) || (s2 <= e1 && e1 <= e2) {
			f(max(s1, s2), min(e1, e2), i, j)
		}
		if e1 < e2 {
			i++
		} else {
			j++
		}
	}
}

// 将nums中的元素进行离散化，返回新的数组和对应的原始值.
// origin[newNums[i]] == nums[i]
func Discretize(nums []int) (newNums []int32, origin []int) {
	newNums = make([]int32, len(nums))
	origin = make([]int, 0, len(newNums))
	order := argSort(int32(len(nums)), func(i, j int32) bool { return nums[i] < nums[j] })
	for _, i := range order {
		if len(origin) == 0 || origin[len(origin)-1] != nums[i] {
			origin = append(origin, nums[i])
		}
		newNums[i] = int32(len(origin) - 1)
	}
	origin = origin[:len(origin):len(origin)]
	return
}

func argSort(n int32, less func(i, j int32) bool) []int32 {
	order := make([]int32, n)
	for i := range order {
		order[i] = int32(i)
	}
	sort.Slice(order, func(i, j int) bool { return less(order[i], order[j]) })
	return order
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

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}
