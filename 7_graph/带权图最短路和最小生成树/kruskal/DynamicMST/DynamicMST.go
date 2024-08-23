// 动态最小生成树.
// https://www.luogu.com.cn/article/xe7hqlw3
// Api:
// 1. NewDynamicMST(n int) *DynamicMST
// 2. (mst *DynamicMST) AddEdge(aId, bId int32, weight int)
// 3. (mst *DynamicMST) GetTotalWeight() int
// 4. (mst *DynamicMST) GetEdgeNum() int32
// 5. (mst *DynamicMST) IsConnected(aId, bId int32) bool
// 6. (mst *DynamicMST) MaxWeightBetween(aId, bId int32) int 从aId到bId的最大边权值

package main

import (
	"fmt"
	"math"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	demo()
}

func demo() {
	mst := NewDynamicMST(5)
	mst.AddEdge(0, 1, -1)
	mst.AddEdge(1, 2, -2)
	mst.AddEdge(2, 3, -3)
	mst.AddEdge(3, 4, -4)
	mst.AddEdge(3, 4, -0)
	mst.AddEdge(0, 4, -0)
	fmt.Println(mst.GetTotalWeight())  // 6
	fmt.Println(mst.GetEdgeNum())      // 4
	fmt.Println(mst.IsConnected(0, 4)) // true
	fmt.Println(mst.MaxWeightBetween(0, 3))
}

// 1697. 检查边长度限制的路径是否存在
// https://leetcode.cn/problems/checking-existence-of-edge-length-limited-paths-ii/description/
func distanceLimitedPathsExist(n int, edgeList [][]int, queries [][]int) []bool {
	mst := NewDynamicMST(int32(n))
	for _, e := range edgeList {
		mst.AddEdge(int32(e[0]), int32(e[1]), e[2])
	}
	res := make([]bool, len(queries))
	for i, q := range queries {
		res[i] = mst.MaxWeightBetween(int32(q[0]), int32(q[1])) < q[2]
	}
	return res
}

const INF int = math.MaxInt64

type DynamicMST struct {
	edgeNum     int32
	totalWeight int
	nodes       []*LCTNode
}

func NewDynamicMST(n int32) *DynamicMST {
	res := &DynamicMST{nodes: make([]*LCTNode, n)}
	for i := int32(0); i < n; i++ {
		res.nodes[i] = NewLCTNode()
		res.nodes[i].id = int32(i)
		res.nodes[i].pushUp()
	}
	for i := int32(1); i < n; i++ {
		node := NewLCTNode()
		node.weight = math.MaxInt64
		node.a = res.nodes[i-1]
		node.b = res.nodes[i]
		node.pushUp()
		join(node.a, node)
		join(node.b, node)
	}
	return res
}

func (mst *DynamicMST) AddEdge(aId, bId int32, weight int) {
	a := mst.nodes[aId]
	b := mst.nodes[bId]
	findRoute(a, b)
	splay(a)
	if a.largest.weight <= weight {
		return
	}
	largest := a.largest
	splay(largest)
	cut(largest.a, largest)
	cut(largest.b, largest)
	if largest.weight < math.MaxInt64 {
		mst.edgeNum--
		mst.totalWeight -= largest.weight
	}
	node := NewLCTNode()
	node.weight = weight
	node.a = a
	node.b = b
	node.pushUp()
	join(node.a, node)
	join(node.b, node)
	mst.edgeNum++
	mst.totalWeight += node.weight
}

func (mst *DynamicMST) GetTotalWeight() int {
	return mst.totalWeight
}

func (mst *DynamicMST) GetEdgeNum() int32 {
	return mst.edgeNum
}

func (mst *DynamicMST) IsConnected(aId, bId int32) bool {
	return mst.MaxWeightBetween(aId, bId) < math.MaxInt64
}

func (mst *DynamicMST) MaxWeightBetween(aId, bId int32) int {
	a := mst.nodes[aId]
	b := mst.nodes[bId]
	findRoute(a, b)
	splay(b)
	return b.largest.weight
}

func createNIL() *LCTNode {
	res := &LCTNode{}
	res.left = res
	res.right = res
	res.father = res
	res.treeFather = res
	res.largest = res
	return res
}

var NIL_NODE = createNIL()

type LCTNode struct {
	reverse                         bool
	id                              int32
	weight                          int
	left, right, father, treeFather *LCTNode
	a, b                            *LCTNode
	largest                         *LCTNode
}

func NewLCTNode() *LCTNode {
	res := &LCTNode{left: NIL_NODE, right: NIL_NODE, father: NIL_NODE, treeFather: NIL_NODE}
	return res
}

func larger(a, b *LCTNode) *LCTNode {
	if a.weight >= b.weight {
		return a
	}
	return b
}

func access(x *LCTNode) {
	last := NIL_NODE
	for x != NIL_NODE {
		splay(x)
		x.right.father = NIL_NODE
		x.right.treeFather = x
		x.setRight(last)
		x.pushUp()
		last = x
		x = x.treeFather
	}
}

func makeRoot(x *LCTNode) {
	access(x)
	splay(x)
	x.reverseFn()
}

func cut(y, x *LCTNode) {
	makeRoot(y)
	access(x)
	splay(y)
	y.right.treeFather = NIL_NODE
	y.right.father = NIL_NODE
	y.setRight(NIL_NODE)
	y.pushUp()
}

func join(y, x *LCTNode) {
	makeRoot(x)
	x.treeFather = y

}

func findRoute(x, y *LCTNode) {
	makeRoot(y)
	access(x)
}

func splay(x *LCTNode) {
	if x == NIL_NODE {
		return
	}
	var y, z *LCTNode
	for y = x.father; y != NIL_NODE; y = x.father {
		if z = y.father; z == NIL_NODE {
			y.pushDown()
			x.pushDown()
			if x == y.left {
				zig(x)
			} else {
				zag(x)
			}
		} else {
			z.pushDown()
			y.pushDown()
			x.pushDown()
			if x == y.left {
				if y == z.left {
					zig(y)
					zig(x)
				} else {
					zig(x)
					zag(x)
				}
			} else {
				if y == z.left {
					zag(x)
					zig(x)
				} else {
					zag(y)
					zag(x)
				}
			}
		}
	}

	x.pushDown()
	x.pushUp()
}

func zig(x *LCTNode) {
	y := x.father
	z := y.father
	b := x.right

	y.setLeft(b)
	x.setRight(y)
	z.changeChild(y, x)

	y.pushUp()
}

func zag(x *LCTNode) {
	y := x.father
	z := y.father
	b := x.left

	y.setRight(b)
	x.setLeft(y)
	z.changeChild(y, x)

	y.pushUp()
}

func findRoot(x *LCTNode) *LCTNode {
	x.pushDown()
	for x.left != NIL_NODE {
		x = x.left
		x.pushDown()
	}
	splay(x)
	return x
}

func (node *LCTNode) pushDown() {
	if node.reverse {
		node.reverse = false

		tmpNode := node.left
		node.left = node.right
		node.right = tmpNode

		node.left.reverseFn()
		node.right.reverseFn()
	}

	node.left.treeFather = node.treeFather
	node.right.treeFather = node.treeFather
}

func (node *LCTNode) reverseFn() {
	node.reverse = !node.reverse
}

func (node *LCTNode) setLeft(x *LCTNode) {
	node.left = x
	x.father = node
}

func (node *LCTNode) setRight(x *LCTNode) {
	node.right = x
	x.father = node
}

func (node *LCTNode) changeChild(y, x *LCTNode) {
	if node.left == y {
		node.setLeft(x)
	} else {
		node.setRight(x)
	}
}

func (node *LCTNode) pushUp() {
	node.largest = larger(node, node.left.largest)
	node.largest = larger(node.largest, node.right.largest)
}
