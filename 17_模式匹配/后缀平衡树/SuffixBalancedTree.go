// 后缀平衡树(动态后缀数组)
// 陈立杰 <<重量平衡树和后缀平衡树在信息学奥赛中的应用>> 2013论文集
// https://oi-wiki.org/string/suffix-bst/
// https://www.cnblogs.com/cjfdf/p/10322533.html
// https://www.luogu.com.cn/article/0e8jjpqk
// https://www.luogu.com.cn/article/wq2yn9lr
// https://www.luogu.com/article/1l23lxcd
// https://blog.csdn.net/YxuanwKeith/article/details/52741250
//
// 后缀平衡树是一种"动态维护后缀排序"的数据结构。
// 它支持在串S的开头添加/删除一个字符，复杂度不依赖于字符集的大小.

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

func init() {
	debug.SetGCPercent(-1)
}

// alpha 的值越小，那么替罪羊树就越容易重构，那么树也就越平衡，查询的效率也就越高，自然修改（加点和删点）的效率也就低了；
// 反之，alpha 的值越大，那么替罪羊树就越不容易重构，那么树也就越不平衡，查询的效率也就越低，自然修改（加点和删点）的效率也就高了。
// 所以，查询多，alpha 就应该小一些；修改多，alpha 就应该大一些。
const ALPHA_NUM int32 = 4
const ALPHA_DENO int32 = 5

var EMPTY_NODE = &Node{}

type Char = byte

type Node struct {
	Weight      float64
	Len         int32 // 代表的后缀的长度
	Id          int32
	left, right int32
	size        int32 // 子树所有结点的数量
}

type SuffixBalancedTree struct {
	Ords         []Char
	Nodes        []*Node
	lower, upper float64
	root         int32

	scapegoat *int32  // 临时保存，用于标记需要重构的结点
	collector []int32 // 临时保存，用于dfs收集结点
}

func NewSuffixBalancedTree(cap int32) *SuffixBalancedTree {
	ords := make([]Char, 0, max32(cap+1, 16))
	nodes := make([]*Node, 0, max32(cap+1, 16))
	ords = append(ords, 0)
	nodes = append(nodes, EMPTY_NODE)
	return &SuffixBalancedTree{Ords: ords, Nodes: nodes}
}

func (t *SuffixBalancedTree) AppendLeft(char Char, id int32) {
	t.Ords = append(t.Ords, char)
	if len(t.Nodes) < len(t.Ords) {
		t._alloc()
	}
	t.scapegoat = nil
	t._insert(&t.root, 0, 1, id)
	if t.scapegoat != nil {
		t._rebuild()
	}
}

// 不使用标记删除.
func (t *SuffixBalancedTree) PopLeft() {
	t._delete(&t.root)
	t.Ords = t.Ords[:len(t.Ords)-1]
}

// 用于树上后缀排序.
// !hack: 自定义比较函数，用于插入时比较两个后缀的大小.
// curPos: 当前插入的字符的位置(1-indexed).
// searchPos: 当前遍历到的结点代表的后缀的位置(1-indexed).
func (t *SuffixBalancedTree) Add(char Char, id int32, less func(tree *SuffixBalancedTree, curPos, searchPos int32) bool) {
	t.Ords = append(t.Ords, char)
	if len(t.Nodes) < len(t.Ords) {
		t._alloc()
	}
	t.scapegoat = nil
	t._insertWithLess(&t.root, 0, 1, id, less)
	if t.scapegoat != nil {
		t._rebuild()
	}
}

// 返回后缀s[index:]的权值，权值越小，字典序越小.
func (t *SuffixBalancedTree) Weight(index int32) float64 {
	index++
	return t.Nodes[int32(len(t.Ords))-index].Weight
}

func (t *SuffixBalancedTree) GetNode(index int32) *Node {
	index++
	return t.Nodes[int32(len(t.Ords))-index]
}

// 求排名第k个的后缀是谁(sa).
func (t *SuffixBalancedTree) Kth(k int32) int32 {
	k++
	cur := t.Nodes[t.root]
	for cur != EMPTY_NODE {
		leftNode := t.Nodes[cur.left]
		if leftNode.size+1 == k {
			return t.Size() - cur.Len
		}
		if leftNode.size >= k {
			cur = leftNode
		} else {
			k -= leftNode.size + 1
			cur = t.Nodes[cur.right]
		}
	}
	panic("unreachable")
}

// 返回后缀s[index:]的排名(0-based).
func (t *SuffixBalancedTree) Rank(index int32) int32 {
	key := t.Weight(index)
	res, cur := int32(0), t.Nodes[t.root]
	for cur != EMPTY_NODE {
		if key < cur.Weight {
			cur = t.Nodes[cur.left]
		} else {
			res += t.Nodes[cur.left].size + 1
			cur = t.Nodes[cur.right]
		}
	}
	return res - 1
}

// 返回字典序严格小于的后缀的个数.
// 如果要求s出现的次数，可以用前驱后继相减得到.
func (t *SuffixBalancedTree) BisectLeftString(n int32, f func(i int32) Char) int32 {
	res, cur := int32(0), t.root
	for cur != 0 {
		compareRes := t._compareString(n, f, cur)
		if compareRes <= 0 {
			cur = t.Nodes[cur].left
		} else {
			res += t.Nodes[t.Nodes[cur].left].size + 1
			cur = t.Nodes[cur].right
		}
	}
	return res
}

// 返回字典序小于等于s的后缀的个数.
// 如果要求s出现的次数，可以用上下界相减得到, 见 P6164.
func (t *SuffixBalancedTree) BisectRightString(n int32, f func(i int32) Char) int32 {
	res, cur := int32(0), t.root
	for cur != 0 {
		compareRes := t._compareString(n, f, cur)
		if compareRes < 0 {
			cur = t.Nodes[cur].left
		} else {
			res += t.Nodes[t.Nodes[cur].left].size + 1
			cur = t.Nodes[cur].right
		}
	}
	return res
}

func (t *SuffixBalancedTree) Sa() []int32 {
	size := t.Size()
	ptr := 0
	sa := make([]int32, size)
	var dfs func(int32)
	dfs = func(x int32) {
		if x == 0 {
			return
		}
		node := t.Nodes[x]
		dfs(node.left)
		sa[ptr] = size - x
		ptr++
		dfs(node.right)
	}
	dfs(t.root)
	return sa
}

func (t *SuffixBalancedTree) Enumerate(f func(node *Node)) {
	var dfs func(int32)
	dfs = func(x int32) {
		if x == 0 {
			return
		}
		node := t.Nodes[x]
		dfs(node.left)
		f(node)
		dfs(node.right)
	}
	dfs(t.root)
}

func (t *SuffixBalancedTree) Size() int32 {
	return int32(len(t.Ords)) - 1
}

func (t *SuffixBalancedTree) _alloc() {
	t.Nodes = append(t.Nodes, &Node{})
}

func (t *SuffixBalancedTree) _delete(x *int32) {
	y := int32(len(t.Ords) - 1)
	nodeY := t.Nodes[y]
	nodeX := t.Nodes[*x]
	nodeX.size--
	if *x == y {
		*x = t._merge(nodeX.left, nodeX.right)
	} else if nodeY.Weight < nodeX.Weight {
		t._delete(&nodeX.left)
	} else {
		t._delete(&nodeX.right)
	}
}

func (t *SuffixBalancedTree) _merge(x, y int32) int32 {
	if x == 0 || y == 0 {
		return x | y
	}
	nodeX, nodeY := t.Nodes[x], t.Nodes[y]
	if nodeX.size > nodeY.size {
		nodeX.right = t._merge(nodeX.right, y)
		t._pushUp(x)
		return x
	} else {
		nodeY.left = t._merge(x, nodeY.left)
		t._pushUp(y)
		return y
	}
}

func (t *SuffixBalancedTree) _collect(cur int32) {
	if cur == 0 {
		return
	}
	node := t.Nodes[cur]
	t._collect(node.left)
	t.collector = append(t.collector, cur)
	t._collect(node.right)
}

func (t *SuffixBalancedTree) _build(a, b int32, lower, upper float64) int32 {
	if a > b {
		return 0
	}
	mid := (a + b) >> 1
	x := t.collector[mid]
	node := t.Nodes[x]
	node.Weight = (lower + upper) / 2
	node.size = b - a + 1
	node.left = t._build(a, mid-1, lower, node.Weight)
	node.right = t._build(mid+1, b, node.Weight, upper)
	t._pushUp(x)
	return x
}

// 一次修改可能会变更整个搜索路径上的所有子树大小，如果多个子树需要重构，选择最大的那颗。
func (t *SuffixBalancedTree) _insert(searchPos *int32, lower, upper float64, id int32) {
	if *searchPos == 0 {
		*searchPos = int32(len(t.Ords) - 1)
		// init
		node := t.Nodes[*searchPos]
		node.Weight = (lower + upper) / 2
		node.Len = t.Size()
		node.Id = id
		node.left, node.right = 0, 0
		node.size = 1
		return
	}
	curPos := int32(len(t.Ords) - 1)
	xNode := t.Nodes[*searchPos]
	// 默认比较函数，如果当前字符相等则去掉一个字符比较.
	if t.Ords[curPos] < t.Ords[*searchPos] || (t.Ords[curPos] == t.Ords[*searchPos] && t.Nodes[curPos-1].Weight < t.Nodes[*searchPos-1].Weight) {
		t._insert(&xNode.left, lower, xNode.Weight, id)
	} else {
		t._insert(&xNode.right, xNode.Weight, upper, id)
	}
	t._pushUp(*searchPos)
	if t._isUnbalanced(*searchPos) {
		t.scapegoat = searchPos
		t.lower = lower
		t.upper = upper
	}
}

func (t *SuffixBalancedTree) _insertWithLess(searchPos *int32, lower, upper float64, id int32, less func(tree *SuffixBalancedTree, curPos, searchPos int32) bool) {
	if *searchPos == 0 {
		*searchPos = int32(len(t.Ords) - 1)
		// init
		node := t.Nodes[*searchPos]
		node.Weight = (lower + upper) / 2
		node.Len = t.Size()
		node.Id = id
		node.left, node.right = 0, 0
		node.size = 1
		return
	}
	curPos := int32(len(t.Ords) - 1)
	xNode := t.Nodes[*searchPos]
	if less(t, curPos, *searchPos) {
		t._insertWithLess(&xNode.left, lower, xNode.Weight, id, less)
	} else {
		t._insertWithLess(&xNode.right, xNode.Weight, upper, id, less)
	}
	t._pushUp(*searchPos)
	if t._isUnbalanced(*searchPos) {
		t.scapegoat = searchPos
		t.lower = lower
		t.upper = upper
	}
}

func (t *SuffixBalancedTree) _rebuild() {
	t.collector = t.collector[:0]
	t._collect(*t.scapegoat)
	*t.scapegoat = t._build(0, int32(len(t.collector))-1, t.lower, t.upper)
}

func (t *SuffixBalancedTree) _pushUp(x int32) {
	cur := t.Nodes[x]
	left, right := t.Nodes[cur.left], t.Nodes[cur.right]
	cur.size = left.size + right.size + 1
}

func (t *SuffixBalancedTree) _isUnbalanced(x int32) bool {
	// +5，避免不必要的重构
	cur := t.Nodes[x]
	left, right := t.Nodes[cur.left], t.Nodes[cur.right]
	threshold := cur.size*ALPHA_NUM + 5*ALPHA_DENO
	return (left.size*ALPHA_DENO > threshold) || (right.size*ALPHA_DENO > threshold)
}

func (t *SuffixBalancedTree) _compareString(n1 int32, f1 func(i int32) Char, n2 int32) int8 {
	for i := int32(0); i < n1; i++ {
		if i >= n2 {
			return 1
		}
		v1, v2 := f1(i), t.Ords[n2-i]
		if v1 != v2 {
			if v1 < v2 {
				return -1
			} else {
				return 1
			}
		}
	}
	if n1 < n2 {
		return -1
	} else {
		return 0
	}
}

func min32(a, b int32) int32 {
	if a < b {
		return a
	}
	return b
}

func max32(a, b int32) int32 {
	if a > b {
		return a
	}
	return b
}

func main() {
	// demo()
	// P3809()
	// P5353()
	P6164()
}

func demo() {
	tree := NewSuffixBalancedTree(0)

	for i := 0; i < 100000; i++ {
		tree.AppendLeft('a', -1)
	}
	fmt.Println(tree.Size())
	for i := 0; i < 100000; i++ {
		tree.PopLeft()
	}

	// aaa
	tree.AppendLeft('a', -1)
	tree.AppendLeft('a', -1)
	tree.AppendLeft('a', -1)
	fmt.Println(tree.Size())
	fmt.Println(tree.Ords)
	fmt.Println(tree.Weight(0))
	fmt.Println(tree.Weight(1), 1)
	fmt.Println(tree.Weight(2), 2)
	fmt.Println(tree.Rank(2))
	fmt.Println(tree.Rank(0))
	fmt.Println(tree.Rank(1))
	fmt.Println(tree.BisectRightString(6, func(i int32) Char { return 'c' }), 999)

	tree.Enumerate(func(node *Node) { fmt.Println(node.Len) })
	fmt.Println(tree.Sa())
	fmt.Println(tree.Kth(2))

	s := "a"
	bytes := []byte(s)
	bytes = append(bytes, 255)
	res1 := tree.BisectLeftString(int32(len(bytes)), func(i int32) Char { return bytes[i] })
	bytes = bytes[:len(bytes)-1]
	bytes[len(bytes)-1]--
	res2 := tree.BisectLeftString(int32(len(bytes)), func(i int32) Char { return bytes[i] })
	fmt.Println(res1-res2, res1, res2, 987)
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
	tree := NewSuffixBalancedTree(int32(len(s)))
	for i := len(s) - 1; i >= 0; i-- {
		tree.AppendLeft(Char(s[i]), -1)
	}

	sa := tree.Sa()
	for _, v := range sa {
		fmt.Fprint(out, v+1, " ")
	}
}

// P5353 树上后缀排序(树上SA)
// https://www.luogu.com.cn/problem/P5353
// 树上的每个字符串为根节点到每个结点路径组成的字符串.
// 对这些字符串按照字典序排序.
// 如果两个节点所代表的字符串完全相同，
// 它们的大小由它们父亲排名的大小决定，即谁的父亲排名大谁就更大；
// 如果仍相同，则由它们编号的大小决定，即谁的编号大谁就更大。
//
// !给定一棵以 1 为根包含 n 个节点的树，保证对于 2∼n 的每个节点，其父亲的编号均小于自己的编号。
// 输出一行 n 个正整数，第 i 个正整数表示代表排名第 i 的字符串的节点编号。
//
// !因为父结点编号均小于子结点编号，所以可以按照编号从小到大的顺序插入，保证父结点在子结点之前插入.
func P5353() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	parents := make([]int32, n+1) // 1-indexed
	for i := int32(2); i < n+1; i++ {
		fmt.Fscan(in, &parents[i])
	}
	var values string // 每个结点上的字符
	fmt.Fscan(in, &values)

	sbt := NewSuffixBalancedTree(n)

	// 如果两个节点所代表的字符串完全相同，
	// 它们的大小由它们父亲排名的大小决定，即谁的父亲排名大谁就更大；
	// 如果仍相同，则由它们编号的大小决定，即谁的编号大谁就更大。
	for id := int32(0); id < n; id++ {
		sbt.Add(
			values[id], id,
			func(tree *SuffixBalancedTree, curPos, searchPos int32) bool {
				ords, nodes := tree.Ords, tree.Nodes
				if ords[curPos] != ords[searchPos] {
					return ords[curPos] < ords[searchPos]
				}
				parentWeight1, parentWeight2 := nodes[parents[curPos]].Weight, nodes[parents[searchPos]].Weight
				if parentWeight1 != parentWeight2 {
					return parentWeight1 < parentWeight2
				}
				return id < nodes[searchPos].Id // 注意curPos对应node的Id还未赋值，需要使用传入的id.
			},
		)
	}

	sbt.Enumerate(func(node *Node) {
		fmt.Fprint(out, node.Id+1, " ")
	})
}

// P5346 【XR-1】柯南家族 (树上后缀排序+二维数点)
// https://www.luogu.com.cn/problem/P5346
// 树上后缀排序+离散化+dfs序转化为二维区间第k小问题, waveletMatrix求解.
func P5346() {}

// P6164 【模板】后缀平衡树
// https://www.luogu.com.cn/problem/P6164
// 给定初始字符串 s 和 q 个操作：
// 1.在当前字符串的后面插入若干个字符。
// 2.在当前字符串的后面删除若干个字符。
// 3.询问字符串 t 作为连续子串在当前字符串中出现了几次？
//
// t 的出现次数等于以 t 为前缀的后缀数量，
// 而以 t 为前缀的后缀数量等于其后继的排名减去其前驱的排名.
// 在 t 后面加入一个极大的字符，就可以构造出 t 的一个后继。
// 将 t 的最后一个字符减小 1，就可以构造出 t 的一个前驱。
// 现在要查询某一个串 t 在后缀平衡树中排名，由于不能保证 t 在后缀平衡树中出现过，所以每次只能暴力比较字符串大小。
func P6164() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()
	// reader, _ := os.Open("P6164_1.in")
	// in := bufio.NewReader(reader)
	// // out := bufio.NewWriter(os.Stdout)
	// writer, _ := os.Create("P6164_1.out")
	// out := bufio.NewWriter(writer)

	var q int32
	fmt.Fscan(in, &q)
	var s string
	fmt.Fscan(in, &s)
	sbt := NewSuffixBalancedTree(q + int32(len(s)))
	for _, c := range s {
		sbt.AppendLeft(Char(c), -1)
	}

	decode := func(bytes []byte, mask int32) {
		m := int32(len(bytes))
		for i := int32(0); i < m; i++ {
			mask = (mask*131 + i) % m
			bytes[i], bytes[mask] = bytes[mask], bytes[i]
		}
	}

	add := func(bytes []byte) {
		for _, c := range bytes {
			sbt.AppendLeft(Char(c), -1)
		}
	}

	delete := func(k int32) {
		for i := int32(0); i < k; i++ {
			sbt.PopLeft()
		}
	}

	// Count.
	query := func(bytes []byte) int32 {
		n := int32(len(bytes))
		// !翻转字符串.
		tmp := make([]byte, len(bytes)+1)
		tmp[0] = 255
		copy(tmp[1:], bytes)
		res1 := sbt.BisectLeftString(n+1, func(i int32) Char { return tmp[n+1-1-i] })

		res2 := sbt.BisectLeftString(n, func(i int32) Char { return bytes[n-1-i] })
		return res1 - res2
	}

	preRes := int32(0)
	for i := int32(0); i < q; i++ {
		var kind string
		fmt.Fscan(in, &kind)
		if kind == "ADD" {
			var str string
			fmt.Fscan(in, &str)
			bytes := []byte(str)
			decode(bytes, preRes)
			add(bytes)
		} else if kind == "DEL" {
			var k int32
			fmt.Fscan(in, &k)
			delete(k)
		} else {
			var str string
			fmt.Fscan(in, &str)
			bytes := []byte(str)
			decode(bytes, preRes)
			res := query(bytes)
			fmt.Fprintln(out, res)
			preRes = preRes ^ res
		}
	}
}
