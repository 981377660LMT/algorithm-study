// 对一个数组A 构建PermutationTree
// 对于一个 1-n 的排列，我们称一个`值域连续`的子数组为段。问一个排列的段的个数
// !即 max(A[i]) - min(A[i]) = right - left (left <= i <= right)

// PermutationTree 是表示区间合并过程的一棵树,
//  叶子节点对应长为1的区间,根节点对应整个区间

// NodeType 为节点类型
//  - LEAF: 叶子节点，表示一个长为1的元素(区间)
//  - JOIN: 从 JOIN 的子结点中选择连续的任意个子节点组成的区间都满足`值域连续`的性质
//          特别地, JOIN_ASC 表示 nums[i]单增, JOIN_DESC 表示nums[i]单减
//  - CUT: 从 CUT 的子结点中选择连续的两个子节点组成的区间都`不`满足`值域连续`的性质

package main

import (
	"fmt"
	"strings"
)

func main() {
	nums := []int{5, 3, 4, 1, 2}
	tree := NewPermutationTree()
	tree.Build(nums)
	fmt.Println(tree.Root)

}

type NodeType int

const (
	JOIN_ASC NodeType = iota
	JOIN_DESC
	LEAF
	CUT
)

type Node struct {
	kind        NodeType
	left, right int // [left, right) => 顶点对应的区间
	minV, maxV  int // [minV, maxV) => 顶点对应的区间中的最小值和最大值
	children    []*Node
}

func (n *Node) String() string {
	mp := map[NodeType]string{
		JOIN_ASC:  "JOIN_ASC",
		JOIN_DESC: "JOIN_DESC",
		LEAF:      "LEAF",
		CUT:       "CUT",
	}
	sb := []string{"Node:\n"}
	sb = append(sb, fmt.Sprintf("kind : %v\n", mp[n.kind]))
	sb = append(sb, fmt.Sprintf("left && right : [%v, %v)\n", n.left, n.right))
	sb = append(sb, fmt.Sprintf("minV && maxV : [%v, %v)\n", n.minV, n.maxV))
	sb = append(sb, fmt.Sprintf("childrenCount : %v", len(n.children)))
	return strings.Join(sb, "  ")
}

func (n *Node) Size() int    { return n.right - n.left }
func (n *Node) IsJoin() bool { return n.kind == JOIN_ASC || n.kind == JOIN_DESC }
func (n *Node) IsLeaf() bool { return n.kind == LEAF }
func (n *Node) IsCut() bool  { return n.kind == CUT }

type PermutationTree struct{ Root *Node }

func NewPermutationTree() *PermutationTree { return &PermutationTree{} }

func (pt *PermutationTree) Build(nums []int) *Node {
	n := len(nums)
	desc, asc := []int{-1}, []int{-1}
	st := []*Node{}

	seg := newRangeAddRangeMin(make([]int, n))
	for i := 0; i < n; i++ {
		for desc[len(desc)-1] != -1 && nums[i] > nums[desc[len(desc)-1]] {
			seg.Update(desc[len(desc)-2]+1, desc[len(desc)-1]+1, nums[i]-nums[desc[len(desc)-1]])
			desc = desc[:len(desc)-1]
		}
		for asc[len(asc)-1] != -1 && nums[i] < nums[asc[len(asc)-1]] {
			seg.Update(asc[len(asc)-2]+1, asc[len(asc)-1]+1, nums[asc[len(asc)-1]]-nums[i])
			asc = asc[:len(asc)-1]
		}
		desc = append(desc, i)
		asc = append(asc, i)

		t := &Node{kind: LEAF, left: i, right: i + 1, minV: nums[i], maxV: nums[i] + 1}
		for {
			kind := CUT
			if len(st) > 0 {
				if st[len(st)-1].maxV == t.minV {
					kind = JOIN_ASC
				} else if t.maxV == st[len(st)-1].minV {
					kind = JOIN_DESC
				}
			}

			if kind != CUT {
				r := st[len(st)-1]
				if kind != r.kind {
					r = &Node{kind: kind, left: r.left, right: r.right, minV: r.minV, maxV: r.maxV, children: []*Node{r}}
				}
				pt.addChild(r, t)
				st = st[:len(st)-1]
				t = r
			} else if seg.Query(0, i+1-t.Size()) == 0 {
				t = &Node{kind: CUT, left: t.left, right: t.right, minV: t.minV, maxV: t.maxV, children: []*Node{t}}
				// do while
				for {
					pt.addChild(t, st[len(st)-1])
					st = st[:len(st)-1]
					if t.maxV-t.minV == t.Size() {
						break
					}
				}
				for i, j := 0, len(t.children)-1; i < j; i, j = i+1, j-1 {
					t.children[i], t.children[j] = t.children[j], t.children[i]
				}
			} else {
				break
			}
		}

		st = append(st, t)
		seg.Update(0, i+1, -1)
	}

	pt.Root = st[0]
	return st[0]
}

func (pt *PermutationTree) addChild(t, c *Node) {
	t.children = append(t.children, c)
	t.left = min(t.left, c.left)
	t.right = max(t.right, c.right)
	t.minV = min(t.minV, c.minV)
	t.maxV = max(t.maxV, c.maxV)
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

const INF int = 1e18

type rangeAddRangeMin struct {
	n, head            int
	rangeMin, rangeAdd []int
}

func newRangeAddRangeMin(nums []int) *rangeAddRangeMin {
	res := &rangeAddRangeMin{}
	res.n = len(nums)
	res.head = 1
	for res.head < res.n {
		res.head <<= 1
	}
	res.rangeMin, res.rangeAdd = make([]int, res.head*2), make([]int, res.head*2)
	for i := range res.rangeMin {
		res.rangeMin[i] = INF
	}
	copy(res.rangeMin[res.head:], nums)
	for i := res.head - 1; i > 0; i-- {
		res.merge(i)
	}
	return res
}

func (r *rangeAddRangeMin) Update(start, end, add int) {
	r.update(start, end, 1, 0, r.head, add)
}
func (r *rangeAddRangeMin) Query(start, end int) int {
	return r.query(start, end, 1, 0, r.head)
}

func (r *rangeAddRangeMin) Get(pos int) int { return r.Query(pos, pos+1) }

func (r *rangeAddRangeMin) query(start, end, pos, left, right int) int {
	if right <= start || end <= left {
		return INF
	}
	if start <= left && right <= end {
		return r.rangeMin[pos] + r.rangeAdd[pos]
	}
	return min(r.query(start, end, pos*2, left, (left+right)/2), r.query(start, end, pos*2+1, (left+right)/2, right)) + r.rangeAdd[pos]
}

func (r *rangeAddRangeMin) update(start, end, pos, left, right, add int) {
	if right <= start || end <= left {
		return
	}
	if start <= left && right <= end {
		r.rangeAdd[pos] += add
		return
	}
	r.update(start, end, pos*2, left, (left+right)/2, add)
	r.update(start, end, pos*2+1, (left+right)/2, right, add)
	r.merge(pos)
}

func (r *rangeAddRangeMin) merge(pos int) {
	r.rangeMin[pos] = min(r.rangeMin[pos*2]+r.rangeAdd[pos*2], r.rangeMin[pos*2+1]+r.rangeAdd[pos*2+1])
}
