package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	// LineAddGetMin()
	SegmentAddGetMin()
}

// https://judge.yosupo.jp/problem/line_add_get_min
// 0 a b: add line y = ax + b
// 1 x: query min y
func LineAddGetMin() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	operations := make([][3]int, 0, q+n)
	queryX := []int{}
	for i := 0; i < n; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		operations = append(operations, [3]int{0, a, b})
	}
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var a, b int
			fmt.Fscan(in, &a, &b)
			operations = append(operations, [3]int{0, a, b})
		} else {
			var x int
			fmt.Fscan(in, &x)
			operations = append(operations, [3]int{1, x, 0})
			queryX = append(queryX, x)
		}
	}

	lichao := NewLiChaoTreeDynamic(-1e9, 1e9, true)
	for _, op := range operations {
		if op[0] == 0 {
			lichao.AddLine(Line{k: op[1], b: op[2]})
		} else {
			fmt.Fprintln(out, lichao.Query(op[1]).value)
		}
	}
}

// https://judge.yosupo.jp/problem/segment_add_get_min
func SegmentAddGetMin() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	operations := make([][5]int, 0, q+n)
	queryX := []int{}
	for i := 0; i < n; i++ {
		var startX, endX, a, b int
		fmt.Fscan(in, &startX, &endX, &a, &b)
		operations = append(operations, [5]int{0, a, b, startX, endX})
	}
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var startX, endX, a, b int
			fmt.Fscan(in, &startX, &endX, &a, &b)
			operations = append(operations, [5]int{0, a, b, startX, endX})
		} else {
			var x int
			fmt.Fscan(in, &x)
			operations = append(operations, [5]int{1, x, 0, 0, 0})
			queryX = append(queryX, x)
		}
	}

	lichao := NewLiChaoTreeDynamic(-1e9, 1e9, true)
	for _, op := range operations {
		if op[0] == 0 {
			k, b, startX, endX := op[1], op[2], op[3], op[4]
			lichao.AddSegment(startX, endX, Line{k: k, b: b})
		} else {
			res := lichao.Query(op[1])
			if res.lineId == -1 {
				fmt.Fprintln(out, "INFINITY")
			} else {
				fmt.Fprintln(out, res.value)
			}
		}
	}
}

type T = int

const INF T = 1e18

type Line struct{ k, b T } // y = k * x + b

// Evaluate を書き変えると、totally monotone な関数群にも使える
func Evaluate(line Line, x int) T {
	return line.k*x + line.b
}

type LiChaoNode struct {
	lineId      int
	left, right *LiChaoNode
}

type queryPair = struct {
	value  T
	lineId int
}

// 动态开点李超线段树.注意`添加线段`时空间消耗较大.
type LiChaoTreeDynamic struct {
	start, end int
	minimize   bool
	lines      []Line
	root       *LiChaoNode
}

func NewLiChaoTreeDynamic(start, end int, minimize bool) *LiChaoTreeDynamic {
	end++
	return &LiChaoTreeDynamic{
		start: start, end: end,
		minimize: minimize,
		root:     &LiChaoNode{lineId: -1},
	}
}

// O(logn)
func (tree *LiChaoTreeDynamic) AddLine(line Line) {
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	tree._addLine(tree.root, id, tree.start, tree.end)
}

// [start, end)
// O(log^2n)
func (tree *LiChaoTreeDynamic) AddSegment(startX, endX int, line Line) {
	if startX >= endX {
		return
	}
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	tree._addSegment(tree.root, startX, endX, id, tree.start, tree.end)
}

// O(logn)
func (tree *LiChaoTreeDynamic) Query(x int) queryPair {
	if !(tree.start <= x && x < tree.end) {
		panic("x is out of range")
	}
	return tree._query(tree.root, x, tree.start, tree.end)
}

func (tree *LiChaoTreeDynamic) Clear() {
	tree.lines = tree.lines[:0]
}

func (tree *LiChaoTreeDynamic) _evaluateInner(fid int, x int) T {
	if fid == -1 {
		if tree.minimize {
			return INF
		}
		return -INF
	}
	return Evaluate(tree.lines[fid], x)
}

func (tree *LiChaoTreeDynamic) _addLine(node *LiChaoNode, fid int, nodeL, nodeR int) *LiChaoNode {
	gid := node.lineId
	fl := tree._evaluateInner(fid, nodeL)
	fr := tree._evaluateInner(fid, nodeR-1)
	gl := tree._evaluateInner(gid, nodeL)
	gr := tree._evaluateInner(gid, nodeR-1)
	var bl, br bool
	if tree.minimize {
		bl = fl < gl
		br = fr < gr
	} else {
		bl = fl > gl
		br = fr > gr
	}
	if bl && br {
		node.lineId = fid
		return node
	}
	if !bl && !br {
		return node
	}
	nodeM := (nodeL + nodeR) >> 1
	fm := tree._evaluateInner(fid, nodeM)
	gm := tree._evaluateInner(gid, nodeM)
	var bm bool
	if tree.minimize {
		bm = fm < gm
	} else {
		bm = fm > gm
	}
	if bm {
		node.lineId = fid
		if bl {
			if node.right == nil {
				node.right = &LiChaoNode{lineId: -1}
			}
			node.right = tree._addLine(node.right, gid, nodeM, nodeR)
		} else {
			if node.left == nil {
				node.left = &LiChaoNode{lineId: -1}
			}
			node.left = tree._addLine(node.left, gid, nodeL, nodeM)
		}
	} else {
		if !bl {
			if node.right == nil {
				node.right = &LiChaoNode{lineId: -1}
			}
			node.right = tree._addLine(node.right, fid, nodeM, nodeR)
		} else {
			if node.left == nil {
				node.left = &LiChaoNode{lineId: -1}
			}
			node.left = tree._addLine(node.left, fid, nodeL, nodeM)
		}
	}
	return node
}

func (tree *LiChaoTreeDynamic) _addSegment(node *LiChaoNode, xl, xr int, fid int, nodeL, nodeR int) *LiChaoNode {
	if nodeL > xl {
		xl = nodeL
	}
	if nodeR < xr {
		xr = nodeR
	}
	if xl >= xr {
		return node
	}
	if nodeL < xl || xr < nodeR {
		nodeM := (nodeL + nodeR) >> 1
		if node.left == nil {
			node.left = &LiChaoNode{lineId: -1}
		}
		if node.right == nil {
			node.right = &LiChaoNode{lineId: -1}
		}
		node.left = tree._addSegment(node.left, xl, xr, fid, nodeL, nodeM)
		node.right = tree._addSegment(node.right, xl, xr, fid, nodeM, nodeR)
		return node
	}
	return tree._addLine(node, fid, nodeL, nodeR)
}

func (tree *LiChaoTreeDynamic) _query(node *LiChaoNode, x int, nodeL, nodeR int) queryPair {
	fid := node.lineId
	res := queryPair{value: tree._evaluateInner(fid, x), lineId: fid}
	nodeM := (nodeL + nodeR) >> 1
	if x < nodeM && node.left != nil {
		cand := tree._query(node.left, x, nodeL, nodeM)
		if tree.minimize {
			if cand.value < res.value {
				res = cand
			}
		} else {
			if cand.value > res.value {
				res = cand
			}
		}
	}
	if x >= nodeM && node.right != nil {
		cand := tree._query(node.right, x, nodeM, nodeR)
		if tree.minimize {
			if cand.value < res.value {
				res = cand
			}
		} else {
			if cand.value > res.value {
				res = cand
			}
		}
	}
	return res
}

func max(a, b int) int {
	if a >= b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
