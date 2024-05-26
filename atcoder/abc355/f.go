package main

import (
	"bufio"
	"fmt"
	stdio "io"
	"math"
	"os"
	"strconv"
)

// from https://atcoder.jp/users/ccppjsrb
var io *Iost

type Iost struct {
	Scanner *bufio.Scanner
	Writer  *bufio.Writer
}

func NewIost(fp stdio.Reader, wfp stdio.Writer) *Iost {
	const BufSize = 2000005
	scanner := bufio.NewScanner(fp)
	scanner.Split(bufio.ScanWords)
	scanner.Buffer(make([]byte, BufSize), BufSize)
	return &Iost{Scanner: scanner, Writer: bufio.NewWriter(wfp)}
}
func (io *Iost) Text() string {
	if !io.Scanner.Scan() {
		panic("scan failed")
	}
	return io.Scanner.Text()
}
func (io *Iost) Atoi(s string) int                 { x, _ := strconv.Atoi(s); return x }
func (io *Iost) Atoi64(s string) int64             { x, _ := strconv.ParseInt(s, 10, 64); return x }
func (io *Iost) Atof64(s string) float64           { x, _ := strconv.ParseFloat(s, 64); return x }
func (io *Iost) NextInt() int                      { return io.Atoi(io.Text()) }
func (io *Iost) NextInt64() int64                  { return io.Atoi64(io.Text()) }
func (io *Iost) NextFloat64() float64              { return io.Atof64(io.Text()) }
func (io *Iost) Print(x ...interface{})            { fmt.Fprint(io.Writer, x...) }
func (io *Iost) Printf(s string, x ...interface{}) { fmt.Fprintf(io.Writer, s, x...) }
func (io *Iost) Println(x ...interface{})          { fmt.Fprintln(io.Writer, x...) }

// F - MST Query
// 頂点に
// 1 から
// N の番号が、辺に
// 1 から
// N−1 の番号が付いた
// N 頂点
// N−1 辺の重み付き無向連結グラフ
// G が与えられます。辺
// i は頂点
// a
// i
// ​
//   と頂点
// b
// i
// ​
//   を結んでおり、その重みは
// c
// i
// ​
//   です。

// Q 個のクエリが与えられるので順に処理してください。
// i 番目のクエリは以下で説明されます。

// 整数
// u
// i
// ​
//
//	,v
//
// i
// ​
//
//	,w
//
// i
// ​
//
//	が与えられる。
//
// G の頂点
// u
// i
// ​
//
//	,v
//
// i
// ​
//
//	の間に重み
//
// w
// i
// ​
//
//	の辺を追加する。その後、
//
// G の最小全域木に含まれる辺の重みの和を出力する。
func main() {
	in := os.Stdin
	out := os.Stdout
	io = NewIost(in, out)
	defer func() {
		io.Writer.Flush()
	}()

	N, Q := io.NextInt(), io.NextInt()
	lct := NewDynamicMST(int32(N))
	for i := 0; i < N-1; i++ {
		a, b, c := io.NextInt(), io.NextInt(), io.NextInt()
		a, b = a-1, b-1
		lct.AddEdge(int32(a), int32(b), c)
	}
	for i := 0; i < Q; i++ {
		u, v, w := io.NextInt(), io.NextInt(), io.NextInt()
		u, v = u-1, v-1
		lct.AddEdge(int32(u), int32(v), w)
		io.Println(lct.GetTotalWeight())
	}
}

const INF int = math.MaxInt64

type DynamicMST struct {
	nodes       []*LCTNode
	totalWeight int
	edgeNum     int32
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
	left, right, father, treeFather *LCTNode
	reverse                         bool
	id                              int32
	a, b                            *LCTNode
	largest                         *LCTNode
	weight                          int
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
