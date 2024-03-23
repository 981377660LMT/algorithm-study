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
// alpha = 3/4
const ALPHA_NUM int32 = 3
const ALPHA_DENO int32 = 4

var EMPTY_NODE = &Node{}

type Char = int32

type Node struct {
	Weight      float64
	left, right int32
	size        int32 // 子树所有结点的数量
}

type SuffixBalancedTree struct {
	data         []Char
	lower, upper float64
	nodes        []*Node
	root         int32
	scapegoat    *int32  // 用于标记需要重构的结点
	collector    []int32 // 用于dfs收集结点
}

func NewSuffixBalancedTree(cap int32) *SuffixBalancedTree {
	data := make([]Char, 0, max32(cap+1, 16))
	nodes := make([]*Node, 0, max32(cap+1, 16))
	data = append(data, 0)
	nodes = append(nodes, EMPTY_NODE)
	return &SuffixBalancedTree{data: data, nodes: nodes}
}

func (t *SuffixBalancedTree) AppendLeft(c Char) {
	t.data = append(t.data, c)
	if len(t.nodes) < len(t.data) {
		t._alloc()
	}
	t.scapegoat = nil
	t._insert(&t.root, 0, 1)
	if t.scapegoat != nil {
		t._rebuild()
	}
}

func (t *SuffixBalancedTree) PopLeft() {
	t._delete(&t.root)
	t.data = t.data[:len(t.data)-1]
}

// 返回下标index的权值，权值越小，字典序越小.
func (t *SuffixBalancedTree) Weight(index int32) float64 {
	index++
	return t.nodes[int32(len(t.data))-index].Weight
}

// 返回下标index的排名(0-based).
func (t *SuffixBalancedTree) Rank(index int32) int32 {
	key := t.Weight(index)
	res, cur := int32(0), t.nodes[t.root]
	for cur != EMPTY_NODE {
		if key < cur.Weight {
			cur = t.nodes[cur.left]
		} else {
			res += t.nodes[cur.left].size + 1
			cur = t.nodes[cur.right]
		}
	}
	return res - 1
}

// 返回字典序小于s的后缀的个数.
func (t *SuffixBalancedTree) RankStr(n int32, f func(i int32) Char) int32 {
	res, cur := int32(0), t.root
	for cur != 0 {
		min_ := min32(cur, n) + 1
		flag := int32(0)
		for i := int32(0); i < min_; i++ {
			v := f(i)
			if v != t.data[cur-i] {
				flag = v - t.data[cur-i]
				break
			}
		}
		if flag <= 0 {
			cur = t.nodes[cur].left
		} else {
			res += t.nodes[t.nodes[cur].left].size + 1
			cur = t.nodes[cur].right
		}
	}
	return res
}

func (t *SuffixBalancedTree) Sa() []int32 {
	sa := make([]int32, t.Size())
	ptr := 0
	t.Enumerate(func(i int32) {
		sa[ptr] = i
		ptr++
	})
	return sa
}

func (t *SuffixBalancedTree) Enumerate(f func(sa int32)) {
	size := t.Size()
	var dfs func(int32)
	dfs = func(x int32) {
		if x == 0 {
			return
		}
		node := t.nodes[x]
		dfs(node.left)
		f(size - x)
		dfs(node.right)
	}
	dfs(t.root)
}

func (t *SuffixBalancedTree) Size() int32 {
	return int32(len(t.data)) - 1
}

func (t *SuffixBalancedTree) _alloc() {
	t.nodes = append(t.nodes, &Node{size: 1})
}

func (t *SuffixBalancedTree) _delete(x *int32) {
	// if *x == 0 {
	// 	return
	// }
	y := int32(len(t.data) - 1)
	nodeY := t.nodes[y]
	nodeX := t.nodes[*x]
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
	nodeX, nodeY := t.nodes[x], t.nodes[y]
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
	node := t.nodes[cur]
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
	node := t.nodes[x]
	node.size = b - a + 1
	node.Weight = (lower + upper) / 2
	node.left = t._build(a, mid-1, lower, node.Weight)
	node.right = t._build(mid+1, b, node.Weight, upper)
	t._pushUp(x)
	return x
}

// 一次修改可能会变更整个搜索路径上的所有子树大小，如果多个子树需要重构，选择最大的那颗。
func (t *SuffixBalancedTree) _insert(x *int32, lower, upper float64) {
	if *x == 0 {
		*x = int32(len(t.data) - 1)
		node := t.nodes[*x]
		node.left, node.right = 0, 0
		node.Weight = (lower + upper) / 2
		node.size = 1
		return
	}
	y := int32(len(t.data) - 1)
	xNode := t.nodes[*x]
	// less
	if t.data[y] < t.data[*x] || (t.data[y] == t.data[*x] && t.nodes[y-1].Weight < t.nodes[*x-1].Weight) {
		t._insert(&xNode.left, lower, xNode.Weight)
	} else {
		t._insert(&xNode.right, xNode.Weight, upper)
	}
	t._pushUp(*x)
	if t._isUnbalanced(*x) {
		t.scapegoat = x
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
	cur := t.nodes[x]
	left, right := t.nodes[cur.left], t.nodes[cur.right]
	cur.size = left.size + right.size + 1
}

func (t *SuffixBalancedTree) _isUnbalanced(x int32) bool {
	// +5，避免不必要的重构
	cur := t.nodes[x]
	left, right := t.nodes[cur.left], t.nodes[cur.right]
	threshold := cur.size*ALPHA_NUM + 5*ALPHA_DENO
	return (left.size*ALPHA_DENO > threshold) || (right.size*ALPHA_DENO > threshold)
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
	// P5346()
	// P5353()
	P6164()
}

func demo() {
	tree := NewSuffixBalancedTree(0)

	for i := 0; i < 100000; i++ {
		tree.AppendLeft('a')
	}
	fmt.Println(tree.Size())
	for i := 0; i < 100000; i++ {
		tree.PopLeft()
	}
	// cbc
	tree.AppendLeft('c')
	tree.AppendLeft('b')
	tree.AppendLeft('a')
	fmt.Println(tree.Size())
	fmt.Println(tree.data)
	fmt.Println(tree.Weight(0))
	fmt.Println(tree.Weight(1), 1)
	fmt.Println(tree.Weight(2), 2)
	fmt.Println(tree.Rank(2))
	fmt.Println(tree.Rank(0))
	fmt.Println(tree.Rank(1))
	fmt.Println(tree.RankStr(6, func(i int32) Char { return 'c' }))

	tree.Enumerate(func(i int32) { fmt.Println(i) })
	fmt.Println(tree.Sa())
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
		tree.AppendLeft(int32(s[i]))
	}

	sa := tree.Sa()
	for _, v := range sa {
		fmt.Fprint(out, v+1, " ")
	}
}

// P5346 【XR-1】柯南家族
// https://www.luogu.com.cn/problem/P5346
func P5346() {

}

// P5353 树上后缀排序
// https://www.luogu.com.cn/problem/P5353
func P5353() {

}

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
func P6164() {}
