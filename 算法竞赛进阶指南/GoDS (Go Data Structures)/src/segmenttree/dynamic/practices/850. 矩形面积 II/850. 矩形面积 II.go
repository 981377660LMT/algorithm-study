// 计算平面中所有 rectangles 所覆盖的 总面积 。任何被两个或多个矩形覆盖的区域应只计算 一次 。
// 返回 总面积 。因为答案可能太大，返回 1e9 + 7 的 模 。
// https://www.luogu.com.cn/problem/P5490 850. 矩形面积 II
// !本题线段树非常特殊
// !1.叶子结点也要pushUp 因为节点值并不全由子节点决定 还与本身有关
//  !2.不用懒标记更新子区间的值 因为区间值由(本身的count)唯一决定

package main

import (
	"fmt"
	"sort"
)

func main() {
	fmt.Println(rectangleArea([][]int{{0, 0, 2, 2}, {1, 0, 2, 3}, {1, 0, 3, 1}}))
}

const MOD int = 1e9 + 7

type Event struct {
	x, y1, y2 int
	kind      uint8 // 0:in 1:out
}

// rectangle[i] = [xi1, yi1, xi2, yi2]
func rectangleArea(rectangles [][]int) int {
	n := len(rectangles)
	events := make([]Event, 0, n*2)
	for i := 0; i < n; i++ {
		x1, y1, x2, y2 := rectangles[i][0], rectangles[i][1], rectangles[i][2], rectangles[i][3]
		in := Event{x: x1, y1: y1, y2: y2, kind: 0}
		out := Event{x: x2, y1: y1, y2: y2, kind: 1}
		events = append(events, in, out)
	}

	// !先进入后退出(此时标记不为负)
	sort.Slice(events, func(i, j int) bool {
		if events[i].x != events[j].x {
			return events[i].x < events[j].x
		}
		return events[i].kind < events[j].kind
	})

	tree := CreateSegmentTree(0, 1e9)
	res := 0
	for i, e := range events {
		if i > 0 {
			res += (e.x - events[i-1].x) * tree.QueryAll().sum
			res %= MOD

		}

		if e.kind == 0 {
			tree.Update(e.y1, e.y2-1, 1) // 注意-1 因为这里是线段表示区间而不是点表示区间
		} else {
			tree.Update(e.y1, e.y2-1, -1)
		}
	}

	return res
}

type Data = struct{ size, sum, count int } // !维护区间覆盖次数以及区间和
type Lazy = int                            // count +1/-1
func e(left, right int) Data               { return Data{size: right - left + 1} }

// !op
func (o *Node) pushUp() {
	if o.data.count == 0 {
		o.data.sum = 0 // !没有标记就从线段树上的两个儿子节点继承
		if o.leftChild != nil {
			o.data.sum += o.leftChild.data.sum
		}
		if o.rightChild != nil {
			o.data.sum += o.rightChild.data.sum
		}
	} else {
		o.data.sum = o.data.size // !如果这个点上有标记， 那么就说明被线段覆盖了， 长度就是 r-l+1
	}
}

func (o *Node) pushDown() {
	mid := (o.left + o.right) >> 1
	if o.leftChild == nil {
		o.leftChild = newNode(o.left, mid)
	}
	if o.rightChild == nil {
		o.rightChild = newNode(mid+1, o.right)
	}
}

// 指定区间上下界建立线段树
func CreateSegmentTree(lower, upper int) *Node {
	root := newNode(lower, upper)
	return root
}

type Node struct {
	left, right           int
	leftChild, rightChild *Node

	data Data
}

func (o *Node) Update(left, right int, lazy Lazy) {
	if left <= o.left && o.right <= right {
		o.data.count += lazy
		o.pushUp() // !注意叶子结点也要pushUp 因为节点值并不全由子节点决定 还与本身有关
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

func (o *Node) QueryAll() Data {
	return o.data
}

func newNode(left, right int) *Node {
	return &Node{left: left, right: right, data: e(left, right)}
}
