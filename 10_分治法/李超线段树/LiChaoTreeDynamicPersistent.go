package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	demo()
	// LineAddGetMin()
	// SegmentAddGetMin()
}

func demo() {
	tree := NewLiChaoTreeDynamicPersistent(-1e9, 1e9, false, true)
	root1 := tree.NewRoot()
	root2 := tree.AddLine(root1, Line{k: 1, b: 1})
	root3 := tree.AddLine(root2, Line{k: 2, b: 2})
	fmt.Println(tree.Query(root1, 0))
	fmt.Println(tree.Query(root2, 0))
	fmt.Println(tree.Query(root3, 0))

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

	lichao := NewLiChaoTreeDynamicPersistent(-1e9, 1e9, true, false)
	root := lichao.NewRoot()
	for _, op := range operations {
		if op[0] == 0 {
			root = lichao.AddLine(root, Line{k: op[1], b: op[2]})
		} else {
			fmt.Fprintln(out, lichao.Query(root, op[1]).value)
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

	lichao := NewLiChaoTreeDynamicPersistent(-1e9, 1e9, true, false)
	root := lichao.NewRoot()
	for _, op := range operations {
		if op[0] == 0 {
			k, b, startX, endX := op[1], op[2], op[3], op[4]
			root = lichao.AddSegment(root, startX, endX, Line{k: k, b: b})
		} else {
			res := lichao.Query(root, op[1])
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

// 可持久化李超线段树.注意`添加线段`时空间消耗较大.
type LiChaoTreeDynamicPersistent struct {
	start, end int
	minimize   bool
	persistent bool
	lines      []Line
}

func NewLiChaoTreeDynamicPersistent(start, end int, minimize bool, persistent bool) *LiChaoTreeDynamicPersistent {
	end++
	return &LiChaoTreeDynamicPersistent{
		start: start, end: end,
		minimize:   minimize,
		persistent: persistent,
	}
}

func (tree *LiChaoTreeDynamicPersistent) NewRoot() *LiChaoNode {
	return nil
}

// O(logn)
func (tree *LiChaoTreeDynamicPersistent) AddLine(root *LiChaoNode, line Line) *LiChaoNode {
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	if root == nil {
		*&root = &LiChaoNode{lineId: -1}
	}
	return tree._addLine(root, id, tree.start, tree.end)
}

// [start, end)
// O(log^2n)
func (tree *LiChaoTreeDynamicPersistent) AddSegment(root *LiChaoNode, startX, endX int, line Line) *LiChaoNode {
	if startX >= endX {
		return root
	}
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	if root == nil {
		root = &LiChaoNode{lineId: -1}
	}
	return tree._addSegment(root, startX, endX, id, tree.start, tree.end)
}

// O(logn)
func (tree *LiChaoTreeDynamicPersistent) Query(root *LiChaoNode, x int) queryPair {
	if !(tree.start <= x && x < tree.end) {
		panic("x is out of range")
	}
	if root == nil {
		if tree.minimize {
			return queryPair{value: INF, lineId: -1}
		}
		return queryPair{value: -INF, lineId: -1}
	}
	return tree._query(root, x, tree.start, tree.end)
}

func (tree *LiChaoTreeDynamicPersistent) Clear() {
	tree.lines = tree.lines[:0]
}

func (tree *LiChaoTreeDynamicPersistent) Copy(node *LiChaoNode) *LiChaoNode {
	if node == nil || !tree.persistent {
		return node
	}
	return &LiChaoNode{lineId: node.lineId, left: node.left, right: node.right}
}

func (tree *LiChaoTreeDynamicPersistent) _evaluateInner(fid int, x int) T {
	if fid == -1 {
		if tree.minimize {
			return INF
		}
		return -INF
	}
	return Evaluate(tree.lines[fid], x)
}

func (tree *LiChaoTreeDynamicPersistent) _addLine(node *LiChaoNode, fid int, nodeL, nodeR int) *LiChaoNode {
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
		*&node = tree.Copy(node)
		node.lineId = fid
		return node
	}
	if !bl && !br {
		return node
	}
	*&node = tree.Copy(node)
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

func (tree *LiChaoTreeDynamicPersistent) _addSegment(node *LiChaoNode, xl, xr int, fid int, nodeL, nodeR int) *LiChaoNode {
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
		*&node = tree.Copy(node)
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

func (tree *LiChaoTreeDynamicPersistent) _query(node *LiChaoNode, x int, nodeL, nodeR int) queryPair {
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
