// 覆盖区间的方案数/区间覆盖的方案数
// 给定[0,X]内的n个`闭区间`
// 选择其中若干个,使得这些区间能够覆盖[0,X]内的所有点
// 求方案数
// n<=1e5

// !按终点排序,dp[i][mex]表示前i个区间覆盖,最左端没有被覆盖的点为mex的方案数
// 选区间的问题都可以用数据结构优化
//  把区间按照右端点排序 dp[i][mex] 然后用区间更新

package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

func main() {
	bf := func(intervals [][2]int, x int) int {
		n := len(intervals)
		res := 0
		for s := 1; s < 1<<n; s++ {
			cur := make(map[int]struct{})
			for i := 0; i < n; i++ {
				if s>>i&1 == 1 {
					for v := intervals[i][0]; v <= intervals[i][1]; v++ {
						cur[v] = struct{}{}
					}
				}
			}
			flag := true
			for v := 0; v <= x; v++ {
				if _, ok := cur[v]; !ok {
					flag = false
					break
				}
			}
			if flag {
				res++
			}
		}
		return res
	}
	for i := 0; i < 1; i++ {
		n := rand.Intn(10) + 1
		x := rand.Intn(1e3)
		intervals := make([][2]int, n)
		for i := range intervals {
			s := rand.Intn(1e3)
			e := rand.Intn(1e3)
			if s > e {
				s, e = e, s
			}
			intervals[i] = [2]int{s, e}
		}
		if solve(intervals, x) != bf(intervals, x) {
			fmt.Println(intervals, x)
			fmt.Println(solve(intervals, x), bf(intervals, x))
			break
		}
	}
	fmt.Println("pass")

	for i := 0; i < 5; i++ {
		n := rand.Intn(1e5) + 1
		x := rand.Intn(1e9)
		intervals := make([][2]int, n)
		for i := range intervals {
			s := rand.Intn(1e5)
			e := rand.Intn(1e5)
			if s > e {
				s, e = e, s
			}
			intervals[i] = [2]int{s, e}
		}
		time1 := time.Now()
		solve(intervals, x)
		fmt.Println(time.Since(time1))
	}
}

const MOD int = 1e9 + 7

// 按照区间右端点排序,dp[i][mex]表示前i个区间,未被覆盖的最左端点为mex的方案数.
func solve(intervals [][2]int, x int) int {
	sort.Slice(intervals, func(i, j int) bool { return intervals[i][1] < intervals[j][1] })
	maxEnd := max(x, intervals[len(intervals)-1][1])
	dp := CreateSegmentTree(0, maxEnd+1) // mex: 0-maxEnd+1
	dp.Update(0, 0, Id{mul: 1, add: 1})
	for _, interval := range intervals {
		s, e := interval[0], interval[1]
		if s-1 >= 0 {
			dp.Update(0, s-1, Id{mul: 2, add: 0}) // mex为之前的值,乘以2
		}
		sum_ := dp.Query(s, e+1).sum
		dp.Update(e+1, e+1, Id{mul: 1, add: sum_}) // mex 为 e+1
	}
	return dp.Query(x+1, maxEnd+1).sum
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

// !线段树维护的数据类型 示例: 区间和
type E = struct{ size, sum int }
type Id = struct{ mul, add int }

func e(left, right int) E { return E{size: right - left + 1} }
func id() Id              { return Id{mul: 1, add: 0} }
func op(e1, e2 E) E {
	return E{size: e1.size + e2.size, sum: (e1.sum + e2.sum) % MOD}
}
func mapping(f Id, g E) E {
	return E{size: g.size, sum: (g.sum*f.mul + g.size*f.add) % MOD}
}
func composition(parent, child Id) Id {
	return Id{mul: (parent.mul * child.mul) % MOD, add: (parent.mul*child.add + parent.add) % MOD}
}

//
//
//
// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper int) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           int
	leftChild, rightChild *Node

	data E
	lazy Id
}

// lower<=left<=right<=upper
func (o *Node) Update(left, right int, lazy Id) {
	if left <= o.left && o.right <= right {
		o.propagate(lazy)
		return
	}

	o.pushDown()
	mid := (o.left + o.right) >> 1
	if left <= mid {
		o.leftChild.Update(left, right, lazy)
	}
	if right > mid {
		o.rightChild.Update(left, right, lazy)
	}
	o.pushUp()
}

// lower<=left<=right<=upper
func (o *Node) Query(left, right int) E {
	if left <= o.left && o.right <= right {
		return o.data
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	res := e(left, right)
	if left <= mid {
		res = op(res, o.leftChild.Query(left, right))
	}
	if right > mid {
		res = op(res, o.rightChild.Query(left, right))
	}
	return res
}

func (o *Node) QueryAll() E {
	if o == nil {
		return e(0, -1)
	}
	return o.data
}

// lower<=pos<=upper
func (o *Node) Set(pos int, val E) {
	if o.left == o.right {
		o.data = val
		return
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	if pos <= mid {
		o.leftChild.Set(pos, val)
	} else {
		o.rightChild.Set(pos, val)
	}
	o.pushUp()
}

// lower<=pos<=upper
func (o *Node) Get(pos int) E {
	if o.left == o.right {
		return o.data
	}
	o.pushDown()
	mid := (o.left + o.right) >> 1
	if pos <= mid {
		return o.leftChild.Get(pos)
	}
	return o.rightChild.Get(pos)
}

func newNode(left, right int) *Node {
	return &Node{left: left, right: right, lazy: id(), data: e(left, right)}
}

// op
func (o *Node) pushUp() {
	o.data = op(o.leftChild.QueryAll(), o.rightChild.QueryAll())
}

func (o *Node) pushDown() {
	mid := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = newNode(o.left, mid)
	}
	if o.rightChild == nil {
		o.rightChild = newNode(mid+1, o.right)
	}
	if o.lazy != id() {
		o.leftChild.propagate(o.lazy)
		o.rightChild.propagate(o.lazy)
		o.lazy = id()
	}
}

// mapping + composition
func (o *Node) propagate(lazy Id) {
	o.data = mapping(lazy, o.data)
	o.lazy = composition(lazy, o.lazy)
}
