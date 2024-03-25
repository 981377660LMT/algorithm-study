// 动态维护lcp的后缀平衡树.

package main

import (
	"bufio"
	"fmt"
	"math"
	"os"
	"runtime/debug"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

const ALPHA_NUM int32 = 4
const ALPHA_DENO int32 = 5

type SbtNode struct {
	Weight             float64
	OffsetToTail       int32
	Lcp                int32
	left, right        *SbtNode
	aliveSize, allSize int32
	key                int32
	alive              bool
	prev, next         *SbtNode
	rangeMinLcp        int32
}

func CreateNILNode() *SbtNode {
	res := &SbtNode{}
	res.Lcp = math.MaxInt32
	res.OffsetToTail = -1
	res.left = res
	res.right = res
	res.key = -1
	res.prev = res
	res.rangeMinLcp = math.MaxInt32
	return res
}

var NIL = CreateNILNode()

func NewSbtNode() *SbtNode {
	res := &SbtNode{
		Lcp:  math.MaxInt32,
		left: NIL, right: NIL,
		prev:        NIL,
		rangeMinLcp: math.MaxInt32,
	}
	return res
}

func NewSbtNodeWithKey(key int32) *SbtNode {
	res := &SbtNode{
		Lcp:  math.MaxInt32,
		left: NIL, right: NIL,
		prev:        NIL,
		key:         key,
		rangeMinLcp: math.MaxInt32,
	}
	res.PushUp()
	return res
}

func (n *SbtNode) PushUp() {
	if n == NIL {
		return
	}
	n.allSize = n.left.allSize + n.right.allSize + 1
	n.aliveSize = n.left.aliveSize + n.right.aliveSize
	n.rangeMinLcp = min32(n.left.rangeMinLcp, n.right.rangeMinLcp)
	if n.alive {
		n.aliveSize++
		n.rangeMinLcp = min32(n.rangeMinLcp, n.Lcp)
	}
}

func (n *SbtNode) PushDown() {}

func (n *SbtNode) CompareTo(o *SbtNode) int8 {
	if n.Weight < o.Weight {
		return -1
	}
	if n.Weight > o.Weight {
		return 1
	}
	return 0
}

func (n *SbtNode) String() string {
	res := strings.Builder{}
	res.WriteString("[")
	remain := 10
	node := n
	for node != nil && remain > 0 {
		res.WriteString(string(node.key))
		res.WriteString(",")
		node = node.next
		if node == n {
			node = nil
		}
		remain--
	}
	if node != nil {
		res.WriteString(",...,")
	}
	res.WriteString("/")
	res.WriteString(string(n.Lcp))
	res.WriteString("]")
	return res.String()
}

type SuffixBalancedTreeLcp struct {
	Root  *SbtNode
	nodes []*SbtNode

	objectHolder **SbtNode
	scapegoat    **SbtNode
	lower, upper float64
	collector    []*SbtNode
}

func NewSuffixBalancedTreeLcp(cap int32) *SuffixBalancedTreeLcp {
	nodes := make([]*SbtNode, 0, max32(cap+1, 16))
	root := NIL
	dummy := NewSbtNodeWithKey(math.MinInt32)
	dummy.next = dummy
	dummy.OffsetToTail = -1
	nodes = append(nodes, dummy)
	return &SuffixBalancedTreeLcp{
		Root:         root,
		nodes:        nodes,
		objectHolder: new(*SbtNode),
	}
}

func (sbt *SuffixBalancedTreeLcp) AddPrefix(x int32) *SbtNode {
	sbt.scapegoat = nil
	sbt._insert(&sbt.Root, x, sbt.nodes[len(sbt.nodes)-1], sbt.objectHolder, 0, 1)
	if sbt.scapegoat != nil {
		sbt._rebuild(sbt.scapegoat, sbt.lower, sbt.upper)
	}

	insertNode := *sbt.objectHolder
	rank := sbt.Rank(insertNode)

	// fix lcp
	var prev, next *SbtNode
	if rank == 1 {
		prev = NIL
	} else {
		prev = sbt._kth(sbt.Root, rank-1)
	}
	if rank == sbt.Root.aliveSize {
		next = NIL
	} else {
		next = sbt._kth(sbt.Root, rank+1)
	}
	sbt._recalcRightLcp(prev, insertNode)
	sbt._recalcRightLcp(insertNode, next)

	sbt.nodes = append(sbt.nodes, insertNode)
	return insertNode
}

func (sbt *SuffixBalancedTreeLcp) RemovePrefix() {
	deleted := sbt.nodes[len(sbt.nodes)-1]
	sbt.nodes = sbt.nodes[:len(sbt.nodes)-1]
	rank := sbt.Rank(deleted)
	var next *SbtNode
	if rank == sbt.Root.aliveSize {
		next = NIL
	} else {
		next = sbt._kth(sbt.Root, rank+1)
	}

	// fix lcp
	if next != NIL {
		nextLcp := min32(next.Lcp, deleted.Lcp)
		next.prev = deleted.prev
		sbt._updateLcp(sbt.Root, next, nextLcp)
	}

	sbt._delete(sbt.Root, deleted)

	// clean or not
	if sbt.Root.aliveSize*2 < sbt.Root.allSize {
		sbt._rebuild(&sbt.Root, 0, 1)
	}
}

func (sbt *SuffixBalancedTreeLcp) Lcp(a, b *SbtNode) int32 {
	if a == b {
		return a.OffsetToTail + 1
	}
	if a.Weight > b.Weight {
		a, b = b, a
	}
	return sbt._rangeLcpExcludeL(sbt.Root, 0, 1, a.Weight, b.Weight)
}

func (sbt *SuffixBalancedTreeLcp) Sa(k int32) *SbtNode {
	k++
	return sbt._kth(sbt.Root, k)
}

func (sbt *SuffixBalancedTreeLcp) Rank(node *SbtNode) int32 {
	return sbt._rank(sbt.Root, node)
}

// <=
func (sbt *SuffixBalancedTreeLcp) Leq(n int32, f func(i int32) int32) int32 {
	return sbt._rankSequence(sbt.Root, n, f)
}

func (sbt *SuffixBalancedTreeLcp) SaAll() []int32 {
	size := sbt.Size()
	res := make([]int32, size)
	ptr := 0
	sbt.Enumerate(func(node *SbtNode) {
		res[ptr] = size - 1 - node.OffsetToTail
		ptr++
	})
	return res
}

func (sbt *SuffixBalancedTreeLcp) Enumerate(f func(node *SbtNode)) {
	var dfs func(node *SbtNode)
	dfs = func(node *SbtNode) {
		if node == NIL {
			return
		}
		node.PushDown()
		dfs(node.left)
		if node.alive {
			f(node)
		}
		dfs(node.right)
	}
	dfs(sbt.Root)
}

func (sbt *SuffixBalancedTreeLcp) Size() int32 {
	return sbt.Root.aliveSize
}

func (sbt *SuffixBalancedTreeLcp) _insert(searchPos **SbtNode, key int32, next *SbtNode, insertNode **SbtNode, L, R float64) {
	if *searchPos == NIL {
		*searchPos = sbt._newNode(key, next, (L+R)/2)
		*insertNode = *searchPos
		return
	}
	(*searchPos).PushDown()
	compareRes := sbt._insertCompare(*searchPos, key, next)
	if compareRes > 0 {
		sbt._insert(&((*searchPos).left), key, next, insertNode, L, (*searchPos).Weight)
	} else {
		sbt._insert(&((*searchPos).right), key, next, insertNode, (*searchPos).Weight, R)
	}
	(*searchPos).PushUp()
	if sbt._isUnbalanced(*searchPos) {
		sbt.scapegoat = searchPos
		sbt.lower = L
		sbt.upper = R
	}
}

func (sbt *SuffixBalancedTreeLcp) _delete(root *SbtNode, node *SbtNode) {
	root.PushDown()
	if root == node {
		root.alive = false
	} else {
		compareRes := root.CompareTo(node)
		if compareRes > 0 {
			sbt._delete(root.left, node)
		} else {
			sbt._delete(root.right, node)
		}
	}
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _updateLcp(root *SbtNode, target *SbtNode, lcp int32) {
	root.PushDown()
	if root == target {
		root.Lcp = lcp
	} else {
		if root.Weight > target.Weight {
			sbt._updateLcp(root.left, target, lcp)
		} else {
			sbt._updateLcp(root.right, target, lcp)
		}
	}
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _rangeLcpExcludeL(root *SbtNode, L, R float64, l, r float64) int32 {
	if root == NIL || R <= l || L > r {
		return math.MaxInt32
	}
	if L > l && R <= r {
		return root.rangeMinLcp
	}
	root.PushDown()
	res := min32(sbt._rangeLcpExcludeL(root.left, L, root.Weight, l, r), sbt._rangeLcpExcludeL(root.right, root.Weight, R, l, r))
	if root.alive && l < root.Weight && root.Weight <= r {
		res = min32(res, root.Lcp)
	}
	return res
}

func (sbt *SuffixBalancedTreeLcp) _considerLcp(a, b *SbtNode) int32 {
	if a.key != b.key {
		return 0
	}
	return 1 + sbt.Lcp(a.next, b.next)
}

func (sbt *SuffixBalancedTreeLcp) _recalcRightLcp(prev, next *SbtNode) {
	if next == NIL {
		return
	}
	next.prev = prev
	lcp := sbt._considerLcp(prev, next)
	sbt._updateLcp(sbt.Root, next, lcp)
}

func (sbt *SuffixBalancedTreeLcp) _kth(root *SbtNode, k int32) (res *SbtNode) {
	if root == NIL {
		return root
	}
	root.PushDown()
	if root.left.aliveSize >= k {
		res = sbt._kth(root.left, k)
	} else {
		leqCount := root.left.aliveSize
		if root.alive {
			leqCount++
		}
		if leqCount >= k {
			res = root
		} else {
			res = sbt._kth(root.right, k-leqCount)
		}
	}
	root.PushUp()
	return
}

func (sbt *SuffixBalancedTreeLcp) _rank(root, node *SbtNode) int32 {
	if root == NIL {
		return 0
	}
	root.PushDown()
	if root == node {
		return root.aliveSize - root.right.aliveSize
	} else {
		compareRes := root.CompareTo(node)
		if compareRes > 0 {
			return sbt._rank(root.left, node)
		} else {
			return root.aliveSize - root.right.aliveSize + sbt._rank(root.right, node)
		}
	}
}

func (sbt *SuffixBalancedTreeLcp) _rankSequence(root *SbtNode, n int32, f func(i int32) int32) int32 {
	if root == NIL {
		return 0
	}
	root.PushDown()
	compareRes := root._sequenceCompare(root, n, f)
	if compareRes > 0 {
		return sbt._rankSequence(root.left, n, f)
	} else {
		return root.aliveSize - root.right.aliveSize + sbt._rankSequence(root.right, n, f)
	}
}

func (sbt *SuffixBalancedTreeLcp) _rebuild(root **SbtNode, L, R float64) {
	sbt.collector = sbt.collector[:0]
	sbt._collect(*root)
	*root = sbt._doRebuild(0, int32(len(sbt.collector))-1, L, R)
}

func (sbt *SuffixBalancedTreeLcp) _doRebuild(l, r int32, L, R float64) *SbtNode {
	if l > r {
		return NIL
	}
	m := (l + r) >> 1
	root := sbt.collector[m]
	root.Weight = (L + R) / 2
	root.left = sbt._doRebuild(l, m-1, L, root.Weight)
	root.right = sbt._doRebuild(m+1, r, root.Weight, R)
	root.PushUp()
	return root
}

func (sbt *SuffixBalancedTreeLcp) _isUnbalanced(node *SbtNode) bool {
	left, right := node.left, node.right
	// +5，避免不必要的重构
	threshold := node.allSize*ALPHA_NUM + 5*ALPHA_DENO
	return (left.allSize*ALPHA_DENO > threshold) || (right.allSize*ALPHA_DENO > threshold)
}

func (sbt *SuffixBalancedTreeLcp) _collect(node *SbtNode) {
	sbt.collector = sbt.collector[:0]
	sbt._doCollect(node)
}

func (sbt *SuffixBalancedTreeLcp) _doCollect(root *SbtNode) {
	if root == NIL {
		return
	}
	root.PushDown()
	sbt._doCollect(root.left)
	sbt.collector = append(sbt.collector, root)
	sbt._doCollect(root.right)
}

func (sbt *SuffixBalancedTreeLcp) _newNode(key int32, next *SbtNode, weight float64) *SbtNode {
	root := NewSbtNode()
	sbt._init(key, root, next, weight)
	return root
}

func (sbt *SuffixBalancedTreeLcp) _init(key int32, root, next *SbtNode, weight float64) {
	root.key = key
	root.Weight = weight
	root.next = next
	root.alive = true
	root.OffsetToTail = next.OffsetToTail + 1
	root.Lcp = math.MaxInt32
	root.prev = NIL
	root.PushUp()
}

func (sbt *SuffixBalancedTreeLcp) _insertCompare(a *SbtNode, key int32, next *SbtNode) int8 {
	if a.key != key {
		if a.key < key {
			return -1
		}
		return 1
	}
	if a.next.Weight < next.Weight {
		return -1
	}
	if a.next.Weight > next.Weight {
		return 1
	}
	return 0

}
func (sbt *SbtNode) _sequenceCompare(root *SbtNode, n int32, f func(i int32) int32) int8 {
	for i := int32(0); i < n; i++ {
		v := f(i)
		if root.key != v {
			if root.key < v {
				return -1
			}
			return 1
		}
		root = root.next
	}
	return 0
}

func (sbt *SuffixBalancedTreeLcp) _nodeCompare(a, b *SbtNode) int8 {
	for a != b {
		if a.key != b.key {
			if a.key < b.key {
				return -1
			}
			return 1
		}
		a = a.next
		b = b.next
	}
	return 0
}

func (sbt *SuffixBalancedTreeLcp) String() string {
	sbt._collect(sbt.Root)
	res := strings.Builder{}
	res.WriteString("{")
	for i, node := range sbt.collector {
		res.WriteString(node.String())
		if i < len(sbt.collector)-1 {
			res.WriteString(",")
		}
	}
	res.WriteString("}")
	return res.String()
}

func min32(a, b int32) int32 {
	if a <= b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a >= b {
		return a
	}
	return b
}

func main() {
	// demo()
	P3809()
}

func demo() {
	s := "banana"
	sbt := NewSuffixBalancedTreeLcp(int32(len(s)))
	suffix := make([]*SbtNode, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		node := sbt.AddPrefix(int32(s[i]))
		suffix = append(suffix, node)
	}

	f := func(start int32) *SbtNode {
		return suffix[int32(len(suffix))-1-start]
	}

	fmt.Println(sbt.Lcp(f(0), f(2)))
	fmt.Println(sbt.Lcp(f(1), f(3)))
	fmt.Println(sbt.Lcp(f(3), f(5)))
}

// P3809 【模板】后缀排序
// https://www.luogu.com.cn/problem/P3809
// 建出后缀平衡树之后，通过中序遍历得到后缀数组。
func P3809() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	sbt := NewSuffixBalancedTreeLcp(int32(len(s)))
	suffix := make([]*SbtNode, 0, len(s))
	for i := len(s) - 1; i >= 0; i-- {
		node := sbt.AddPrefix(int32(s[i]))
		suffix = append(suffix, node)
	}
	f := func(start int32) *SbtNode {
		return suffix[int32(len(suffix))-1-start]
	}
	_ = f

	// 5 3 1 4 2

	// 1 3 5 2 4

	sa := sbt.SaAll()
	for _, v := range sa {
		fmt.Fprint(out, v+1, " ")
	}
}
