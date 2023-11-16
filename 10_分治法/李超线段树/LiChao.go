package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
	"sort"
)

func main() {
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

	lichao := NewLiChaoTreeCompress(queryX, true)
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

	lichao := NewLiChaoTreeCompress(queryX, true)
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

func demo() {
	tree := NewLiChaoTreeCompress([]int{0, 1, 2, 3, 4, 5, 6, 99999999999}, true)
	tree.AddLine(Line{k: 1, b: 0})
	tree.AddLine(Line{k: 2, b: 0})
	fmt.Println(tree.Query(6))
	tree = NewLiChaoTreeNoCompress(0, 100, true)
	tree.AddSegment(0, 6, Line{k: 1, b: 0})
	tree.AddSegment(0, 100, Line{k: 2, b: 0})
	fmt.Println(tree.Query(6))
}

type T = int

const INF T = 1e18

type Line struct{ k, b T } // y = k * x + b

// Evaluate を書き変えると、totally monotone な関数群にも使える
func Evaluate(line Line, x int) T {
	return line.k*x + line.b
}

type queryPair = struct {
	value  T
	lineId int
}

type LiChaoTree struct {
	n, offset     int
	lower, higher int
	compress      bool
	minimize      bool
	xs            []int
	lines         []Line
	lineIds       []int
}

// 指定查询的 x 值建立李超线段树，采用坐标压缩.
func NewLiChaoTreeCompress(queryX []int, minimize bool) *LiChaoTree {
	set := make(map[int]struct{})
	for _, x := range queryX {
		set[x] = struct{}{}
	}
	unique := make([]int, 0, len(set))
	for x := range set {
		unique = append(unique, x)
	}
	sort.Ints(unique)
	n := len(unique)
	log := 1
	for (1 << log) < n {
		log++
	}
	offset := 1 << log
	lineIds := make([]int, offset<<1)
	for i := range lineIds {
		lineIds[i] = -1
	}
	return &LiChaoTree{
		n: n, offset: offset,
		compress: true, minimize: minimize,
		xs:      unique,
		lineIds: lineIds,
	}
}

// 指定查询的 x 值范围建立李超线段树，不采用坐标压缩.
// higher - lower <= 1e6.
func NewLiChaoTreeNoCompress(lower, higher int, minimize bool) *LiChaoTree {
	n := higher - lower
	log := 1
	for (1 << log) < n {
		log++
	}
	offset := 1 << log
	lineIds := make([]int, offset<<1)
	for i := range lineIds {
		lineIds[i] = -1
	}
	return &LiChaoTree{
		n: n, offset: offset,
		lower: lower, higher: higher,
		compress: false, minimize: minimize,
		lineIds: lineIds,
	}
}

// O(logn)
func (tree *LiChaoTree) AddLine(line Line) {
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	tree._addLineAt(1, id)
}

// [start, end)
// O(log^2n)
func (tree *LiChaoTree) AddSegment(startX, endX int, line Line) {
	if startX >= endX {
		return
	}
	id := len(tree.lines)
	tree.lines = append(tree.lines, line)
	startX = tree._getIndex(startX) + tree.offset
	endX = tree._getIndex(endX) + tree.offset
	for startX < endX {
		if startX&1 == 1 {
			tree._addLineAt(startX, id)
			startX++
		}
		if endX&1 == 1 {
			endX--
			tree._addLineAt(endX, id)
		}
		startX >>= 1
		endX >>= 1
	}
}

// O(logn)
func (tree *LiChaoTree) Query(x int) queryPair {
	x = tree._getIndex(x)
	pos := x + tree.offset
	res := queryPair{lineId: -1}
	if tree.minimize {
		res.value = INF
	} else {
		res.value = -INF
	}

	for pos > 0 {
		if id := tree.lineIds[pos]; id != -1 && id != res.lineId {
			cand := queryPair{value: tree._evaluateInner(id, x), lineId: id}
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
		pos >>= 1
	}
	return res
}

func (tree *LiChaoTree) _addLineAt(i, fid int) {
	upperBit := 31 - bits.LeadingZeros32(uint32(i))
	left := (tree.offset >> upperBit) * (i - (1 << upperBit))
	right := left + (tree.offset >> upperBit)
	minimize := tree.minimize
	for left < right {
		gid := tree.lineIds[i]
		fl := tree._evaluateInner(fid, left)
		fr := tree._evaluateInner(fid, right-1)
		gl := tree._evaluateInner(gid, left)
		gr := tree._evaluateInner(gid, right-1)
		var bl, br bool
		if minimize {
			bl = fl < gl
			br = fr < gr
		} else {
			bl = fl > gl
			br = fr > gr
		}
		if bl && br {
			tree.lineIds[i] = fid
			return
		}
		if !bl && !br {
			return
		}
		mid := (left + right) >> 1
		fm := tree._evaluateInner(fid, mid)
		gm := tree._evaluateInner(gid, mid)
		var bm bool
		if minimize {
			bm = fm < gm
		} else {
			bm = fm > gm
		}
		if bm {
			tree.lineIds[i] = fid
			fid = gid
			if !bl {
				i <<= 1
				right = mid
			} else {
				i = i<<1 | 1
				left = mid
			}
		} else {
			if bl {
				i <<= 1
				right = mid
			} else {
				i = i<<1 | 1
				left = mid
			}
		}
	}
}

func (tree *LiChaoTree) _evaluateInner(fid int, x int) T {
	if fid == -1 {
		if tree.minimize {
			return INF
		}
		return -INF
	}
	var target int
	if tree.compress {
		target = tree.xs[min(x, tree.n-1)]
	} else {
		target = x + tree.lower
	}
	return Evaluate(tree.lines[fid], target)
}

func (tree *LiChaoTree) _getIndex(x int) int {
	if tree.compress {
		return sort.SearchInts(tree.xs, x)
	}
	if x < tree.lower || x > tree.higher {
		panic("x out of range")
	}
	return x - tree.lower
}

func min(a, b int) int {
	if a <= b {
		return a
	}
	return b
}
