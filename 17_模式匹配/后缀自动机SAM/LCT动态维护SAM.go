package main

import (
	"bufio"
	"fmt"
	"os"
)

func main() {
	P5212()
}

// P5212 SubString
// https://www.luogu.com.cn/problem/P5212
// 给定一个字符串 init，要求支持两个操作：
// 在当前字符串的后面插入一个字符串。
// 询问字符串 s 在当前字符串中出现了几次。
// 强制在线。
//
// 动态子串出现次数 => 定位到结点，答案是一个结点子树size之和，lct 维护子树和即可。
func P5212() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	normalize := func(chars []int32, mask int32) {
		m := int32(len(chars))
		for i := int32(0); i < m; i++ {
			mask = (mask*131 + int32(i)) % m
			chars[i], chars[mask] = chars[mask], chars[i]
		}
	}

	sam := NewSuffixAutomaton()
	lct := NewLinkCutTreeSubTreeSum()
	nodes := []*treeNode{}

	onCreate := func(pos int32) {
		for int32(len(nodes)) <= pos {
			nodes = append(nodes, lct.Alloc(0))
		}
	}
	onLink := func(child, parent int32) {
		lct.LinkEdge(nodes[child], nodes[parent])
	}
	onCut := func(child, parent int32) {
		lct.CutEdge(nodes[child], nodes[parent])
	}

	update := func(bytes []int32) {
		for _, b := range bytes {
			lastPos := sam.Add(b, onCreate, onLink, onCut)
			lct.Set(nodes[lastPos], 1)
		}
	}
	query := func(bytes []int32) int32 {
		pos := int32(0)
		for _, b := range bytes {
			pos = sam.Nodes[pos].Next[b-OFFSET]
			if pos == -1 {
				return 0
			}
		}
		node := nodes[pos]
		lct.Evert(nodes[0])
		return lct.QuerySubTree(node)
	}

	var q int32
	fmt.Fscan(in, &q)
	var init string
	fmt.Fscan(in, &init)

	update([]int32(init))

	mask := int32(0)
	for i := int32(0); i < q; i++ {
		var op, s string
		fmt.Fscan(in, &op, &s)
		bytes := []int32(s)
		normalize(bytes, mask)
		if op == "QUERY" {
			res := query(bytes)
			fmt.Fprintln(out, res)
			mask ^= res
		} else {
			update(bytes)
		}
	}
}

const INF int32 = 1e9 + 10

const SIGMA int32 = 2    // 字符集大小
const OFFSET int32 = 'A' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
	End    int32        // 最长的字符在原串的下标, 实点下标为非负数, 虚点下标为负数
}

type SuffixAutomaton struct {
	Nodes   []*Node
	LastPos int32 // 当前插入的字符对应的节点(实点，原串的一个前缀)
	n       int32 // 当前字符串长度
}

func NewSuffixAutomaton() *SuffixAutomaton {
	res := &SuffixAutomaton{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0, -1))
	return res
}

// 每次插入会增加一个实点，可能增加一个虚点.
// 返回当前前缀对应的节点编号(lastPos).
func (sam *SuffixAutomaton) Add(
	char int32,
	onCreate func(pos int32), onLink func(child, parent int32), onCut func(child, parent int32),
) int32 {
	c := char - OFFSET
	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[sam.LastPos].MaxLen+1, sam.Nodes[sam.LastPos].End+1))
	onCreate(newNode)
	p := sam.LastPos
	for p != -1 && sam.Nodes[p].Next[c] == -1 {
		sam.Nodes[p].Next[c] = newNode
		p = sam.Nodes[p].Link
	}
	q := int32(0)
	if p != -1 {
		q = sam.Nodes[p].Next[c]
	}
	if p == -1 || sam.Nodes[p].MaxLen+1 == sam.Nodes[q].MaxLen {
		sam.Nodes[newNode].Link = q
	} else {
		// 不够用，需要新增一个虚点
		newQ := int32(len(sam.Nodes))
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1, -abs32(sam.Nodes[q].End)))
		sam.Nodes[len(sam.Nodes)-1].Next = sam.Nodes[q].Next
		onCreate(newQ)
		onLink(newQ, sam.Nodes[newQ].Link)
		onCut(q, sam.Nodes[q].Link)
		onLink(q, newQ)
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[c] == q {
			sam.Nodes[p].Next[c] = newQ
			p = sam.Nodes[p].Link
		}
	}
	sam.n++
	sam.LastPos = newNode
	onLink(newNode, sam.Nodes[newNode].Link)
	return sam.LastPos
}

func (sam *SuffixAutomaton) Size() int32 {
	return int32(len(sam.Nodes))
}

func (sam *SuffixAutomaton) newNode(link, maxLen, end int32) *Node {
	res := &Node{Link: link, MaxLen: maxLen, End: end}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}

type E = int32 // 子树和

// 维护子树和的 LCT.
type LinkCutTreeSubTreeSum struct{}

// check: 删除、添加边时是否检查边的存在.
func NewLinkCutTreeSubTreeSum() *LinkCutTreeSubTreeSum {
	return &LinkCutTreeSubTreeSum{}
}

func (lct *LinkCutTreeSubTreeSum) Build(n int32, f func(i int32) E) []*treeNode {
	nodes := make([]*treeNode, n)
	for i := int32(0); i < n; i++ {
		nodes[i] = lct.Alloc(f(i))
	}
	return nodes
}

func (lct *LinkCutTreeSubTreeSum) Alloc(key E) *treeNode {
	res := newTreeNode(key)
	lct.update(res)
	return res
}

// 将 t 转为根节点.
func (lct *LinkCutTreeSubTreeSum) Evert(t *treeNode) {
	lct.expose(t)
	lct.toggle(t)
	lct.push(t)
}

// 需要保证 u,v 之间没有边.
func (lct *LinkCutTreeSubTreeSum) LinkEdge(child, parent *treeNode) {
	lct.Evert(child)
	lct.expose(parent)
	child.p = parent
	parent.r = child
	lct.update(parent)
}

// 需要保证u,v之间有边.
func (lct *LinkCutTreeSubTreeSum) CutEdge(u, v *treeNode) {
	lct.Evert(u)
	lct.expose(v)
	parent := v.l
	v.l = nil
	lct.update(v)
	parent.p = nil
}

func (lct *LinkCutTreeSubTreeSum) Lca(u, v *treeNode) *treeNode {
	if !lct.IsConnected(u, v) {
		return nil
	}
	lct.expose(u)
	return lct.expose(v)
}

func (lct *LinkCutTreeSubTreeSum) KthAncestor(x *treeNode, k int32) *treeNode {
	lct.expose(x)
	for x != nil {
		lct.push(x)
		if x.r != nil && x.r.cnt > k {
			x = x.r
		} else {
			if x.r != nil {
				k -= x.r.cnt
			}
			if k == 0 {
				return x
			}
			k--
			x = x.l
		}
	}
	return nil
}

func (lct *LinkCutTreeSubTreeSum) GetParent(t *treeNode) *treeNode {
	lct.expose(t)
	p := t.l
	if p == nil {
		return nil
	}
	for {
		lct.push(p)
		if p.r == nil {
			return p
		}
		p = p.r
	}
}

func (lct *LinkCutTreeSubTreeSum) Jump(from, to *treeNode, k int32) *treeNode {
	lct.Evert(to)
	lct.expose(from)
	for {
		lct.push(from)
		rs := int32(0)
		if from.r != nil {
			rs = from.r.cnt
		}
		if k < rs {
			from = from.r
			continue
		}
		if k == rs {
			break
		}
		k -= rs + 1
		from = from.l
	}
	lct.splay(from)
	return from
}

// !查询前注意要调用 Evert 选定根节点(换根).
func (lct *LinkCutTreeSubTreeSum) QuerySubTree(t *treeNode) E {
	lct.expose(t)
	return t.key + t.sub
}

func (lct *LinkCutTreeSubTreeSum) Set(t *treeNode, key E) *treeNode {
	lct.expose(t)
	t.key = key
	lct.update(t)
	return t
}

func (lct *LinkCutTreeSubTreeSum) Get(t *treeNode) E {
	return t.key
}

func (lct *LinkCutTreeSubTreeSum) IsConnected(u, v *treeNode) bool {
	return u == v || lct.GetRoot(u) == lct.GetRoot(v)
}

func (lct *LinkCutTreeSubTreeSum) GetRoot(t *treeNode) *treeNode {
	lct.expose(t)
	for t.l != nil {
		lct.push(t)
		t = t.l
	}
	return t
}

func (lct *LinkCutTreeSubTreeSum) expose(t *treeNode) *treeNode {
	var rp *treeNode
	for cur := t; cur != nil; cur = cur.p {
		lct.splay(cur)
		if cur.r != nil {
			cur.Add(cur.r)
		}
		cur.r = rp
		if cur.r != nil {
			cur.Erase(cur.r)
		}
		lct.update(cur)
		rp = cur
	}
	lct.splay(t)
	return rp
}

func (lct *LinkCutTreeSubTreeSum) update(t *treeNode) {
	t.cnt = 1
	if t.l != nil {
		t.cnt += t.l.cnt
	}
	if t.r != nil {
		t.cnt += t.r.cnt
	}
	t.Merge(t.l, t.r)
}

func (lct *LinkCutTreeSubTreeSum) rotr(t *treeNode) {
	x := t.p
	y := x.p
	x.l = t.r
	if t.r != nil {
		t.r.p = x
	}
	t.r = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeSubTreeSum) rotl(t *treeNode) {
	x := t.p
	y := x.p
	x.r = t.l
	if t.l != nil {
		t.l.p = x
	}
	t.l = x
	x.p = t
	lct.update(x)
	lct.update(t)
	t.p = y
	if y != nil {
		if y.l == x {
			y.l = t
		}
		if y.r == x {
			y.r = t
		}
		lct.update(y)
	}
}

func (lct *LinkCutTreeSubTreeSum) toggle(t *treeNode) {
	t.l, t.r = t.r, t.l
	t.rev = !t.rev
}

func (lct *LinkCutTreeSubTreeSum) push(t *treeNode) {
	if t.rev {
		if t.l != nil {
			lct.toggle(t.l)
		}
		if t.r != nil {
			lct.toggle(t.r)
		}
		t.rev = false
	}
}

func (lct *LinkCutTreeSubTreeSum) splay(t *treeNode) {
	lct.push(t)
	for !t.IsRoot() {
		q := t.p
		if q.IsRoot() {
			lct.push(q)
			lct.push(t)
			if q.l == t {
				lct.rotr(t)
			} else {
				lct.rotl(t)
			}
		} else {
			r := q.p
			lct.push(r)
			lct.push(q)
			lct.push(t)
			if r.l == q {
				if q.l == t {
					lct.rotr(q)
					lct.rotr(t)
				} else {
					lct.rotl(t)
					lct.rotr(t)
				}
			} else {
				if q.r == t {
					lct.rotl(q)
					lct.rotl(t)
				} else {
					lct.rotr(t)
					lct.rotl(t)
				}
			}
		}
	}
}

type treeNode struct {
	l, r, p       *treeNode
	key, sum, sub E
	cnt           int32
	rev           bool
}

func newTreeNode(key E) *treeNode {
	res := &treeNode{key: key, sum: key, cnt: 1}
	return res
}

func (n *treeNode) IsRoot() bool {
	return n.p == nil || (n.p.l != n && n.p.r != n)
}

func (n *treeNode) Add(other *treeNode)   { n.sub += other.sum }
func (n *treeNode) Erase(other *treeNode) { n.sub -= other.sum }
func (n *treeNode) Merge(n1, n2 *treeNode) {
	var tmp1, tmp2 E
	if n1 != nil {
		tmp1 = n1.sum
	}
	if n2 != nil {
		tmp2 = n2.sum
	}
	n.sum = tmp1 + n.key + n.sub + tmp2
}

func abs32(x int32) int32 {
	if x < 0 {
		return -x
	}
	return x
}
