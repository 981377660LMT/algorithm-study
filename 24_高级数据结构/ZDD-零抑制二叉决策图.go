//
// `Zero-suppressed binary decision diagram` with family algebra operations
// !零抑制二叉决策图(ZDD/ZBDD)
//
// Description:
//   ZDD maintains a set of sets of integers. Possible operations are
//     - or(A, B): union of A and B                 { 1, 12, 123 } or { 1, 23 } = { 1, 12, 123, 23 }
//     - and(A, B): intersection of A and B         { 1, 12, 123 } and { 12 } = { 12 }
//     - difference(A, B): difference of A and B    { 1, 12, 123 } difference { 12 } = { 1, 123 }
//     - count(A) : the number of sets in A
//		 - var(v): a set containing v                 var(1) = { 1 }
//
//     - mul(A, B): Kronecker product of A and B    { 1, 12, 123 } mul { 34 } = { 134, 1234 }
//     - div(A, B): quotient of A and B,            { 1, 12, 123 } div { 12 } = { {}, 3 }
//     - mod(A, B): reminder of A and B             { 1, 12, 123 } mod { 12 } = { 1 }
//
// Complexity:
//   O( #essentially disjoint subsets ).
//
// References:
//   S. Minato (1993):
//     Zero-suppressed BDDs for set manipulation in combinatorial problems.
//     Proceedings of the 30st annual Design Automation Conference, pp. 272-277.
//   S. Minato (1994):
//     Calculation of unate cube set algebra using zero-suppressed BDDs.
//     Proceedings of the 31st annual Design Automation Conference, pp. 420-424.
//

package main

import (
	"fmt"
	"sort"
	"strconv"
	"strings"
)

func main() {
	zdd := NewZDD()
	set1 := zdd.NewSetWith([]int{1, 12, 123})
	set2 := zdd.NewSetWith([]int{2})
	fmt.Println(zdd.ToSet(zdd.Or(set1, set2)))         // 并
	fmt.Println(zdd.ToSet(zdd.And(set1, set2)))        // 交
	fmt.Println(zdd.ToSet(zdd.Difference(set1, set2))) // 差
	fmt.Println(zdd.Count(set1))

	fmt.Println(zdd.ToSet(zdd.Mul(set1, set2)))
	fmt.Println(zdd.ToSet(zdd.Div(set1, set2)))
	fmt.Println(zdd.ToSet(zdd.Mod(set1, set2)))

}

// !ZDD 压缩集合.
type ZDD struct {
	nodes []ZNode
}

func NewZDD() *ZDD {
	return &ZDD{nodes: []ZNode{{-1, 0, 0}, {-1, 0, 0}}}
}

func (zdd *ZDD) NewSet() int {
	return 0 // !0: 空集 1: 只包含空集的集合
}

func (zdd *ZDD) NewSetWith(nums []int) int {
	res := zdd.NewSet()
	for _, v := range nums {
		res = zdd.Or(res, zdd.Var(v))
	}
	return res
}

func (zdd *ZDD) Var(v int) int {
	return zdd._getNode(v, 0, 1)
}

func (zdd *ZDD) Or(s1, s2 int) int {
	h := make(map[[2]int]int)
	if s1 > s2 {
		s1, s2 = s2, s1
	}
	if s1 == 0 || s1 == s2 {
		return s2
	}
	pair := [2]int{s1, s2}
	if i, ok := h[pair]; ok {
		return i
	}
	xNode, yNode := zdd.nodes[s1], zdd.nodes[s2]
	if xNode[0] > yNode[0] {
		res := zdd._getNode(xNode[0], zdd.Or(xNode[1], s2), xNode[2])
		h[pair] = res
		return res
	}
	if xNode[0] < yNode[0] {
		res := zdd._getNode(yNode[0], zdd.Or(yNode[1], s1), yNode[2])
		h[pair] = res
		return res
	}
	res := zdd._getNode(xNode[0], zdd.Or(xNode[1], yNode[1]), zdd.Or(xNode[2], yNode[2]))
	h[pair] = res
	return res
}

func (zdd *ZDD) And(s1, s2 int) int {
	h := make(map[[2]int]int)
	if s1 > s2 {
		s1, s2 = s2, s1
	}
	if s1 == 0 || s1 == s2 {
		return s1
	}
	pair := [2]int{s1, s2}
	if i, ok := h[pair]; ok {
		return i
	}
	xNode, yNode := zdd.nodes[s1], zdd.nodes[s2]
	if xNode[0] > yNode[0] {
		res := zdd.And(xNode[1], s2)
		h[pair] = res
		return res
	}
	if xNode[0] < yNode[0] {
		res := zdd.And(yNode[1], s1)
		h[pair] = res
		return res
	}
	res := zdd._getNode(xNode[0], zdd.And(xNode[1], yNode[1]), zdd.And(xNode[2], yNode[2]))
	h[pair] = res
	return res
}

func (zdd *ZDD) Difference(s1, s2 int) int {
	h := make(map[[2]int]int)
	if s1 == 0 || s2 == 0 {
		return s1
	}
	if s1 == s2 {
		return 0
	}
	pair := [2]int{s1, s2}
	if i, ok := h[pair]; ok {
		return i
	}
	xNode, yNode := zdd.nodes[s1], zdd.nodes[s2]
	if xNode[0] > yNode[0] {
		res := zdd._getNode(xNode[0], zdd.Difference(xNode[1], s2), xNode[2])
		h[pair] = res
		return res
	}
	if xNode[0] < yNode[0] {
		res := zdd.Difference(s1, yNode[1])
		h[pair] = res
		return res
	}
	res := zdd._getNode(xNode[0], zdd.Difference(xNode[1], yNode[1]), zdd.Difference(xNode[2], yNode[2]))
	h[pair] = res
	return res
}

// 统计集合中元素个数.
func (zdd *ZDD) Count(s int) int {
	h := make(map[int]int)
	if s <= 1 {
		return s
	}
	if i, ok := h[s]; ok {
		return i
	}
	xNode := zdd.nodes[s]
	res := zdd.Count(xNode[1]) + zdd.Count(xNode[2])
	h[s] = res
	return res
}

func (zdd *ZDD) ToSet(s int) string {
	res := []int{}
	path := []int{}
	var dfs func(x int)
	dfs = func(x int) {
		if x == 1 {
			res = append(res, path...)
		} else if x != 0 {
			path = append(path, zdd.nodes[x][0])
			dfs(zdd.nodes[x][2])
			path = path[:len(path)-1]
			dfs(zdd.nodes[x][1])
		}
	}
	dfs(s)
	sort.Ints(res)
	sb := strings.Builder{}
	sb.WriteString(fmt.Sprintf("ZDDSet(%d){", len(res)))
	for i, v := range res {
		if i > 0 {
			sb.WriteByte(',')
		}
		sb.WriteString(strconv.Itoa(v))
	}
	sb.WriteByte('}')
	return sb.String()
}

func (zdd *ZDD) _getNode(v, lo, hi int) int {
	mp := make(map[ZNode]int)
	if hi == 0 {
		return lo
	}
	t := ZNode{v, lo, hi}
	if i, ok := mp[t]; ok {
		return i
	}
	zdd.nodes = append(zdd.nodes, t)
	newV := len(zdd.nodes) - 1
	mp[t] = newV
	return newV
}

type ZNode [3]int // v, lo, hi (集合中的元素, 最小值, 最大值)

func ZLess(x, y ZNode) bool {
	if v1, v2 := x[0], y[0]; v1 != v2 {
		return v1 < v2
	}
	if lo1, lo2 := x[1], y[1]; lo1 != lo2 {
		return lo1 < lo2
	}
	return x[2] < y[2]
}

func (z *ZDD) Mul(s1, s2 int) int {
	h := make(map[[2]int]int)
	if s1 > s2 {
		s1, s2 = s2, s1
	}
	if s1 == 0 || s2 == 1 {
		return s1
	}
	if s2 == 0 || s1 == 1 {
		return s2
	}
	pair := [2]int{s1, s2}
	if i, ok := h[pair]; ok {
		return i
	}
	xn, yn := z.nodes[s1], z.nodes[s2]
	if xn[0] > yn[0] {
		res := z._getNode(xn[0], z.Mul(xn[1], s2), z.Mul(xn[2], s2))
		h[pair] = res
		return res
	}
	if xn[0] < yn[0] {
		res := z._getNode(yn[0], z.Mul(yn[1], s1), z.Mul(yn[2], s1))
		h[pair] = res
		return res
	}
	res := z._getNode(xn[0], z.Mul(xn[1], yn[1]), z.Or(z.Or(z.Mul(xn[2], yn[2]), z.Mul(xn[2], yn[1])), z.Mul(xn[1], yn[2])))
	h[pair] = res
	return res
}

func (zdd *ZDD) Div(s1, s2 int) int {
	h := make(map[[2]int]int)
	if s2 == 1 {
		return s1
	}
	if s1 <= 1 {
		return 0
	}
	if s1 == s2 {
		return 1
	}
	xNode, yNode := zdd.nodes[s1], zdd.nodes[s2]
	if xNode[0] < yNode[0] { // node[y].v does not occur in x
		return 0
	}
	pair := [2]int{s1, s2}
	if i, ok := h[pair]; ok {
		return i
	}
	if xNode[0] != yNode[0] {
		res := zdd._getNode(xNode[0], zdd.Div(xNode[1], s2), zdd.Div(xNode[2], s2))
		h[pair] = res
		return res
	}
	z := zdd.Div(xNode[2], yNode[2])
	if z != 0 && yNode[1] != 0 {
		res := zdd.And(z, zdd.Div(xNode[1], yNode[1]))
		h[pair] = res
		return res
	}
	h[pair] = z
	return z
}

func (zdd *ZDD) Mod(s1, s2 int) int {
	return zdd.Difference(s1, zdd.Mul(s2, zdd.Div(s1, s2)))
}
