// 广义后缀自动机(General Suffix Automaton，GSA)
// https://www.yhzq-blog.cc/%e5%b9%bf%e4%b9%89%e5%90%8e%e7%bc%80%e8%87%aa%e5%8a%a8%e6%9c%ba%e6%80%bb%e7%bb%93/
// https://zhuanlan.zhihu.com/p/34838533
// https://oi-wiki.org/string/general-sam/
// https://www.luogu.com.cn/article/pm10t1pc
// https://www.luogu.com/article/w967d5rp
//
// note:
//
// 0. 一个能接受多个串所有子串的自动机。
// 1. 构建方式：
//   - 伪广义后缀自动机:
//     !如果给出的是多个字符串而不是一个trie，则可以使用.
//     对每个串，重复在同一个 SAM 上进行建立.
//     !每次建完一个串以后就把lastPos 指针移到root上面，接着建下一个串。
//     注意"插入字符串时需要看一下当前准备插入的位置是否已经有结点了".
//     如果有的话我们只需要在其基础上额外判断一下拆分 SAM 结点的情况；否则的话就和普通的 SAM 插入一模一样了
//   - Trie树上的广义后缀自动机：建立在 Trie 树上的 SAM 称为广义 SAM
//
// !2. 自动机和广义后缀自动机中"用于构建"该自动机的所有串的所有前缀节点的树链的并的长度和是 O(L*sqrt(L)) 的。
//
//	!对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.标记次数之和不超过O(Lsqrt(L)).
//	记count[u]为结点u的子树中的endPos的原串个数，则sum(count)的数量级为O(Lsqrt(L))，L为所有串长之和.
//	证明(利用根号分治)：
//	则若一个串长S>SQRT(L)，这样的串显然不超SQRT(L)个，而由于广义 SAM 上的节点数量级线性所以这里的总贡献数量级为O(LSQRT(L))。
//	而对于串长不超过SQRT(L)的串，贡献数量级为O(nSQRT(L))。
//	https://blog.csdn.net/ylsoi/article/details/94476894
//	https://ddosvoid.github.io/2020/10/18/%E6%B5%85%E8%B0%88%E6%A0%B9%E5%8F%B7%E7%AE%97%E6%B3%95/
//	喵星球上的点名 https://www.luogu.com.cn/problem/P2336
//	Sevenk Love Oimaster https://www.luogu.com.cn/problem/SP8093
//  !可以被线段树合并优化成O(nlogn).
//
// !3. 广义 SAM 出现子串查询：
//	  对于 n 个串的广义后缀自动机，求出每个点对应的字符串是哪些原串的子串。
//		和线段树合并维护 Endpos 集合基本一致，将每个后缀对应的点附上对应串的标记，
//    然后在树结构上 DFS 进行线段树合并即可得到每个串的出现位置。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const SIGMA int32 = 26   // 字符集大小
const OFFSET int32 = 'a' // 字符集的起始字符

type Node struct {
	Next   [SIGMA]int32 // SAM 转移边
	Link   int32        // 后缀链接
	MaxLen int32        // 当前节点对应的最长子串的长度
}

type SuffixAutomatonGeneral struct {
	Nodes []*Node
	n     int32 // 当前字符串长度

	doubling *DoublingSimple
}

func NewSuffixAutomatonGeneral() *SuffixAutomatonGeneral {
	res := &SuffixAutomatonGeneral{}
	res.Nodes = append(res.Nodes, res.newNode(-1, 0))
	return res
}

// !需要在插入新串之前将lastPos置为0.
// eg:
//
//	sam := NewSuffixAutomatonGeneral()
//	for _,word := range words {
//	  lastPos = 0
//	  for _,c := range word {
//	    lastPos = sam.Add(lastPos,c)
//	  }
//	}
//
// 返回当前前缀对应的节点编号(lastPos).
func (sam *SuffixAutomatonGeneral) Add(lastPos int32, char int32) int32 {
	c := char - OFFSET
	sam.n++

	// 判断当前转移结点是否存在.
	if tmp := sam.Nodes[lastPos].Next[c]; tmp != -1 {
		lastNode, nextNode := sam.Nodes[lastPos], sam.Nodes[tmp]
		if lastNode.MaxLen+1 == nextNode.MaxLen {
			return tmp
		} else {
			newQ := int32(len(sam.Nodes))
			sam.Nodes = append(sam.Nodes, sam.newNode(nextNode.Link, lastNode.MaxLen+1))
			sam.Nodes[newQ].Next = nextNode.Next
			sam.Nodes[tmp].Link = newQ
			for lastPos != -1 && sam.Nodes[lastPos].Next[c] == tmp {
				sam.Nodes[lastPos].Next[c] = newQ
				lastPos = sam.Nodes[lastPos].Link
			}
			return newQ
		}
	}

	newNode := int32(len(sam.Nodes))
	// 新增一个实点以表示当前最长串
	sam.Nodes = append(sam.Nodes, sam.newNode(-1, sam.Nodes[lastPos].MaxLen+1))
	p := lastPos
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
		sam.Nodes = append(sam.Nodes, sam.newNode(sam.Nodes[q].Link, sam.Nodes[p].MaxLen+1))
		sam.Nodes[len(sam.Nodes)-1].Next = sam.Nodes[q].Next
		sam.Nodes[q].Link = newQ
		sam.Nodes[newNode].Link = newQ
		for p != -1 && sam.Nodes[p].Next[c] == q {
			sam.Nodes[p].Next[c] = newQ
			p = sam.Nodes[p].Link
		}
	}
	return newNode
}

func (sam *SuffixAutomatonGeneral) AddString(s string, prefixEndFn func(i, pos int32)) (lastPos int32) {
	lastPos = 0
	for i, c := range s {
		lastPos = sam.Add(lastPos, c)
		if prefixEndFn != nil {
			prefixEndFn(int32(i), lastPos)
		}
	}
	return
}

func (sam *SuffixAutomatonGeneral) Size() int32 {
	return int32(len(sam.Nodes))
}

// 后缀链接树.也叫 parent tree.
func (sam *SuffixAutomatonGeneral) BuildTree() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(1); v < n; v++ {
		p := sam.Nodes[v].Link
		graph[p] = append(graph[p], v)
	}
	return graph
}

func (sam *SuffixAutomatonGeneral) BuildDAG() [][]int32 {
	n := int32(len(sam.Nodes))
	graph := make([][]int32, n)
	for v := int32(0); v < n; v++ {
		for _, to := range sam.Nodes[v].Next {
			if to != -1 {
				graph[v] = append(graph[v], to)
			}
		}
	}
	return graph
}

// 将结点按照长度进行计数排序，返回后缀链接树的dfs顺序.
// 注意：后缀链接树上父亲的MaxLen值一定小于儿子，但不能认为编号小的节点MaxLen值也小.
// 常数比建图 + dfs 小.
func (sam *SuffixAutomatonGeneral) GetDfsOrder() []int32 {
	nodes, size, n := sam.Nodes, sam.Size(), sam.n
	counter := make([]int32, n+1)
	for i := int32(0); i < size; i++ {
		counter[nodes[i].MaxLen]++
	}
	for i := int32(1); i <= n; i++ {
		counter[i] += counter[i-1]
	}
	order := make([]int32, size)
	for i := size - 1; i >= 0; i-- {
		v := nodes[i].MaxLen
		counter[v]--
		order[counter[v]] = i
	}
	return order
}

// 对每个模式串，返回其在sam上各个endPos的大小.
// dfsOrder: 后缀链接树的dfs顺序.
// isPrefix: 判断pos是否是模式串的前缀.
func (sam *SuffixAutomatonGeneral) GetEndPosSize(dfsOrder []int32, isPrefix func(pos int32) bool) []int32 {
	size := sam.Size()
	endPosSize := make([]int32, size)
	for i := size - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		if isPrefix(cur) { // 实点
			endPosSize[cur]++
		}
		pre := sam.Nodes[cur].Link
		endPosSize[pre] += endPosSize[cur]
	}
	return endPosSize
}

// 给定子串的起始位置和结束位置，返回子串在fail树上的位置.
// 快速定位子串, 可以与其它字符串算法配合使用.
// 倍增往上跳到 MaxLen>=end-start 的最后一个节点.
// start: 子串起始位置, end: 子串结束位置, endPosOfEnd: 子串结束位置在fail树上的位置.
func (sam *SuffixAutomatonGeneral) LocateSubstring(start, end int32, endPosOfEnd int32) (pos int32) {
	target := end - start
	_, pos = sam.Doubling().MaxStep(endPosOfEnd, func(p int32) bool { return sam.Nodes[p].MaxLen >= target })
	return
}

func (sam *SuffixAutomatonGeneral) DistinctSubstringAt(pos int32) int32 {
	if pos == 0 {
		return 0
	}
	return sam.Nodes[pos].MaxLen - sam.Nodes[sam.Nodes[pos].Link].MaxLen
}

// 本质不同的子串个数.
func (sam *SuffixAutomatonGeneral) DistinctSubstring() int {
	res := 0
	for i := 1; i < len(sam.Nodes); i++ {
		res += int(sam.DistinctSubstringAt(int32(i)))
	}
	return res
}

// 获取pattern在sam上的位置.
func (sam *SuffixAutomatonGeneral) GetPos(pattern string) (pos int32, ok bool) {
	pos = 0
	for _, c := range pattern {
		pos = sam.Nodes[pos].Next[c-OFFSET]
		if pos == -1 {
			return -1, false
		}
	}
	return pos, true
}

// 后缀连接树上倍增.
func (sam *SuffixAutomatonGeneral) Doubling() *DoublingSimple {
	if sam.doubling == nil {
		size := sam.Size()
		doubling := NewDoubling(size, int(size))
		for i := int32(1); i < size; i++ {
			doubling.Add(i, sam.Nodes[i].Link)
		}
		doubling.Build()
		sam.doubling = doubling
	}
	return sam.doubling
}

func (sam *SuffixAutomatonGeneral) newNode(link, maxLen int32) *Node {
	res := &Node{Link: link, MaxLen: maxLen}
	for i := int32(0); i < SIGMA; i++ {
		res.Next[i] = -1
	}
	return res
}

type E = int32

func e() E { return 0 }
func op(a, b E) E {
	return a + b
}
func merge(a, b E) E { // 合并两个不同的树的结点的函数
	// return min32(1, a+b)
	return a + b
}

type SegNode struct {
	value                 E
	leftChild, rightChild *SegNode
}

func (n *SegNode) String() string {
	return fmt.Sprintf("%v", n.value)
}

type SegmentTreeMerger struct {
	left, right int32
}

// 指定闭区间[left,right]建立Merger.
func NewSegmentTreeMerger(left, right int32) *SegmentTreeMerger {
	return &SegmentTreeMerger{left: left, right: right}
}

// NewRoot().
func (sm *SegmentTreeMerger) Alloc() *SegNode {
	return &SegNode{value: e()}
}

// 权值线段树求第 k 小.
// 调用前需保证 1 <= k <= node.value.
func (sm *SegmentTreeMerger) Kth(node *SegNode, k int32, getCount func(node *SegNode) int32) (res int32, ok bool) {
	if k < 1 || k > getCount(node) {
		return
	}
	return sm._kth(k, node, sm.left, sm.right, getCount), true
}

func (sm *SegmentTreeMerger) Get(node *SegNode, index int32) E {
	return sm._get(node, index, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) Set(node *SegNode, index int32, value E) {
	sm._set(node, index, value, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) Query(node *SegNode, left, right int32) E {
	return sm._query(node, left, right, sm.left, sm.right)
}

func (sm *SegmentTreeMerger) QueryAll(node *SegNode) E {
	return sm._eval(node)
}

func (sm *SegmentTreeMerger) Update(node *SegNode, index int32, value E) {
	sm._update(node, index, value, sm.left, sm.right)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeMerger) Merge(a, b *SegNode) *SegNode {
	return sm._merge(a, b, sm.left, sm.right)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeMerger) MergeDestructively(a, b *SegNode) *SegNode {
	return sm._mergeDestructively(a, b, sm.left, sm.right)
}

// 线段树分裂，将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分.
func (sm *SegmentTreeMerger) Split(node *SegNode, left, right int32) (this, other *SegNode) {
	this, other = sm._split(node, nil, left, right, sm.left, sm.right)
	return
}

func (sm *SegmentTreeMerger) _kth(k int32, node *SegNode, left, right int32, getCount func(*SegNode) int32) int32 {
	if left == right {
		return left
	}
	mid := (left + right) >> 1
	leftCount := int32(0)
	if node.leftChild != nil {
		leftCount = getCount(node.leftChild)
	}
	if leftCount >= k {
		return sm._kth(k, node.leftChild, left, mid, getCount)
	} else {
		return sm._kth(k-leftCount, node.rightChild, mid+1, right, getCount)
	}
}

func (sm *SegmentTreeMerger) _get(node *SegNode, index int32, left, right int32) E {
	if node == nil {
		return e()
	}
	if left == right {
		return node.value
	}
	mid := (left + right) >> 1
	if index <= mid {
		return sm._get(node.leftChild, index, left, mid)
	} else {
		return sm._get(node.rightChild, index, mid+1, right)
	}
}
func (sm *SegmentTreeMerger) _query(node *SegNode, L, R int32, left, right int32) E {
	if node == nil {
		return e()
	}
	if L <= left && right <= R {
		return node.value
	}
	mid := (left + right) >> 1
	if R <= mid {
		return sm._query(node.leftChild, L, R, left, mid)
	}
	if L > mid {
		return sm._query(node.rightChild, L, R, mid+1, right)
	}
	return op(sm._query(node.leftChild, L, R, left, mid), sm._query(node.rightChild, L, R, mid+1, right))
}

func (sm *SegmentTreeMerger) _set(node *SegNode, index int32, value E, left, right int32) {
	if left == right {
		node.value = value
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.leftChild, index, value, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._set(node.rightChild, index, value, mid+1, right)
	}
	node.value = op(sm._eval(node.leftChild), sm._eval(node.rightChild))
}

func (sm *SegmentTreeMerger) _update(node *SegNode, index int32, value E, left, right int32) {
	if left == right {
		node.value = op(node.value, value)
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._update(node.leftChild, index, value, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._update(node.rightChild, index, value, mid+1, right)
	}
	node.value = op(sm._eval(node.leftChild), sm._eval(node.rightChild))
}

func (sm *SegmentTreeMerger) _merge(a, b *SegNode, left, right int32) *SegNode {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := sm.Alloc()
	if left == right {
		newNode.value = merge(a.value, b.value)
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	newNode.value = op(sm._eval(newNode.leftChild), sm._eval(newNode.rightChild))
	return newNode
}

func (sm *SegmentTreeMerger) _mergeDestructively(a, b *SegNode, left, right int32) *SegNode {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.value = merge(a.value, b.value)
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	a.value = op(sm._eval(a.leftChild), sm._eval(a.rightChild))
	return a
}

func (sm *SegmentTreeMerger) _split(a, b *SegNode, L, R int32, left, right int32) (*SegNode, *SegNode) {
	if a == nil || L > right || R < left {
		return a, nil
	}
	if L <= left && right <= R {
		return nil, a
	}
	if b == nil {
		b = sm.Alloc()
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	a.value = op(sm._eval(a.leftChild), sm._eval(a.rightChild))
	b.value = op(sm._eval(b.leftChild), sm._eval(b.rightChild))
	return a, b
}

func (sm *SegmentTreeMerger) _eval(node *SegNode) E {
	if node == nil {
		return e()
	}
	return node.value
}

// !注意不存在的情况(maxCount=0).
type MaxCount = int32
type SegNode2 struct {
	MaxCount              MaxCount // 出现次数最多的权值出现的次数.
	MaxIndex              int32    // 出现次数最多的权值.如果有多个，取最小的.
	leftChild, rightChild *SegNode2
}

func (n *SegNode2) String() string {
	return fmt.Sprintf("%v", n.MaxCount)
}

type SegmentTreeOnRangeWithIndex struct {
	min, max int32
}

// 指定闭区间[min,max]建立权值线段树.
func NewSegmentTreeOnRangeWithIndex(min, max int32) *SegmentTreeOnRangeWithIndex {
	return &SegmentTreeOnRangeWithIndex{min: min, max: max}
}

// NewRoot().
func (sm *SegmentTreeOnRangeWithIndex) Alloc() *SegNode2 {
	return &SegNode2{}
}

func (sm *SegmentTreeOnRangeWithIndex) Get(node *SegNode2, index int32) MaxCount {
	return sm._get(node, index, sm.min, sm.max)
}

func (sm *SegmentTreeOnRangeWithIndex) Set(node *SegNode2, index int32, count MaxCount) {
	sm._set(node, index, count, sm.min, sm.max)
}

func (sm *SegmentTreeOnRangeWithIndex) Query(node *SegNode2, left, right int32) (maxCount MaxCount, maxIndex int32) {
	return sm._query(node, left, right, sm.min, sm.max)
}

func (sm *SegmentTreeOnRangeWithIndex) QueryAll(node *SegNode2) (maxCount MaxCount, maxIndex int32) {
	if node == nil {
		return
	}
	return node.MaxCount, node.MaxIndex
}

func (sm *SegmentTreeOnRangeWithIndex) Add(node *SegNode2, index int32, count MaxCount) {
	sm._update(node, index, count, sm.min, sm.max)
}

// 用一个新的节点存合并的结果，会生成重合节点数量的新节点.
func (sm *SegmentTreeOnRangeWithIndex) Merge(a, b *SegNode2) *SegNode2 {
	return sm._merge(a, b, sm.min, sm.max)
}

// 把第二棵树直接合并到第一棵树上，比较省空间，缺点是会丢失合并前树的信息.
func (sm *SegmentTreeOnRangeWithIndex) MergeDestructively(a, b *SegNode2) *SegNode2 {
	return sm._mergeDestructively(a, b, sm.min, sm.max)
}

// 线段树分裂，将区间 [left,right] 从原树分离到 other 上, this 为原树的剩余部分.
func (sm *SegmentTreeOnRangeWithIndex) Split(node *SegNode2, left, right int32) (this, other *SegNode2) {
	this, other = sm._split(node, nil, left, right, sm.min, sm.max)
	return
}

func (sm *SegmentTreeOnRangeWithIndex) _get(node *SegNode2, index int32, left, right int32) MaxCount {
	if node == nil {
		return 0
	}
	if left == right {
		return node.MaxCount
	}
	mid := (left + right) >> 1
	if index <= mid {
		return sm._get(node.leftChild, index, left, mid)
	} else {
		return sm._get(node.rightChild, index, mid+1, right)
	}
}

func (sm *SegmentTreeOnRangeWithIndex) _query(node *SegNode2, L, R int32, left, right int32) (maxCount MaxCount, maxIndex int32) {
	if node == nil {
		return
	}
	if L <= left && right <= R {
		return node.MaxCount, node.MaxIndex
	}
	mid := (left + right) >> 1
	if R <= mid {
		return sm._query(node.leftChild, L, R, left, mid)
	}
	if L > mid {
		return sm._query(node.rightChild, L, R, mid+1, right)
	}
	c1, i1 := sm._query(node.leftChild, L, R, left, mid)
	c2, i2 := sm._query(node.rightChild, L, R, mid+1, right)
	if c1 > c2 {
		return c1, i1
	} else if c1 < c2 {
		return c2, i2
	}
	if i1 <= i2 {
		return c1, i1
	} else {
		return c2, i2
	}
}

func (sm *SegmentTreeOnRangeWithIndex) _set(node *SegNode2, index int32, count MaxCount, left, right int32) {
	if left == right {
		node.MaxCount = count
		node.MaxIndex = left
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._set(node.leftChild, index, count, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._set(node.rightChild, index, count, mid+1, right)
	}
	sm._pushUp(node)
}

func (sm *SegmentTreeOnRangeWithIndex) _update(node *SegNode2, index int32, count MaxCount, left, right int32) {
	if left == right {
		node.MaxCount += count
		node.MaxIndex = left
		return
	}
	mid := (left + right) >> 1
	if index <= mid {
		if node.leftChild == nil {
			node.leftChild = sm.Alloc()
		}
		sm._update(node.leftChild, index, count, left, mid)
	} else {
		if node.rightChild == nil {
			node.rightChild = sm.Alloc()
		}
		sm._update(node.rightChild, index, count, mid+1, right)
	}
	sm._pushUp(node)
}

func (sm *SegmentTreeOnRangeWithIndex) _merge(a, b *SegNode2, left, right int32) *SegNode2 {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	newNode := sm.Alloc()
	if left == right {
		newNode.MaxCount = a.MaxCount + b.MaxCount
		newNode.MaxIndex = left
		return newNode
	}
	mid := (left + right) >> 1
	newNode.leftChild = sm._merge(a.leftChild, b.leftChild, left, mid)
	newNode.rightChild = sm._merge(a.rightChild, b.rightChild, mid+1, right)
	sm._pushUp(newNode)
	return newNode
}

func (sm *SegmentTreeOnRangeWithIndex) _mergeDestructively(a, b *SegNode2, left, right int32) *SegNode2 {
	if a == nil || b == nil {
		if a == nil {
			return b
		}
		return a
	}
	if left == right {
		a.MaxCount += b.MaxCount
		a.MaxIndex = left
		return a
	}
	mid := (left + right) >> 1
	a.leftChild = sm._mergeDestructively(a.leftChild, b.leftChild, left, mid)
	a.rightChild = sm._mergeDestructively(a.rightChild, b.rightChild, mid+1, right)
	sm._pushUp(a)
	return a
}

func (sm *SegmentTreeOnRangeWithIndex) _split(a, b *SegNode2, L, R int32, left, right int32) (*SegNode2, *SegNode2) {
	if a == nil || L > right || R < left {
		return a, nil
	}
	if L <= left && right <= R {
		return nil, a
	}
	if b == nil {
		b = sm.Alloc()
	}
	mid := (left + right) >> 1
	a.leftChild, b.leftChild = sm._split(a.leftChild, b.leftChild, L, R, left, mid)
	a.rightChild, b.rightChild = sm._split(a.rightChild, b.rightChild, L, R, mid+1, right)
	sm._pushUp(a)
	sm._pushUp(b)
	return a, b
}

func (sm *SegmentTreeOnRangeWithIndex) _evelCount(node *SegNode2) MaxCount {
	if node == nil {
		return 0
	}
	return node.MaxCount
}

func (sm *SegmentTreeOnRangeWithIndex) _pushUp(node *SegNode2) {
	left, right := node.leftChild, node.rightChild
	b1, b2 := left == nil, right == nil
	if b1 || b2 {
		if b1 && b2 {
			return
		}
		if b1 {
			node.MaxCount = right.MaxCount
			node.MaxIndex = right.MaxIndex
		} else {
			node.MaxCount = left.MaxCount
			node.MaxIndex = left.MaxIndex
		}
	} else {
		if left.MaxCount > right.MaxCount {
			node.MaxCount = left.MaxCount
			node.MaxIndex = left.MaxIndex
		} else if left.MaxCount < right.MaxCount {
			node.MaxCount = right.MaxCount
			node.MaxIndex = right.MaxIndex
		} else {
			if left.MaxIndex <= right.MaxIndex {
				node.MaxCount = left.MaxCount
				node.MaxIndex = left.MaxIndex
			} else {
				node.MaxCount = right.MaxCount
				node.MaxIndex = right.MaxIndex
			}
		}
	}
}

type DoublingSimple struct {
	n   int32
	log int32
	to  []int32
}

func NewDoubling(n int32, maxStep int) *DoublingSimple {
	res := &DoublingSimple{n: n, log: int32(bits.Len(uint(maxStep)))}
	size := n * res.log
	res.to = make([]int32, size)
	for i := int32(0); i < size; i++ {
		res.to[i] = -1
	}
	return res
}

func (d *DoublingSimple) Add(from, to int32) {
	d.to[from] = to
}

func (d *DoublingSimple) Build() {
	n := d.n
	for k := int32(0); k < d.log-1; k++ {
		for v := int32(0); v < n; v++ {
			w := d.to[k*n+v]
			next := (k+1)*n + v
			if w == -1 {
				d.to[next] = -1
				continue
			}
			d.to[next] = d.to[k*n+w]
		}
	}
}

// 求从 `from` 状态开始转移 `step` 次的最终状态的编号。
// 不存在时返回 -1。
func (d *DoublingSimple) Jump(from int32, step int) (to int32) {
	to = from
	for k := int32(0); k < d.log; k++ {
		if to == -1 {
			return
		}
		if step&(1<<k) != 0 {
			to = d.to[k*d.n+to]
		}
	}
	return
}

// 求从 `from` 状态开始转移 `step` 次，满足 `check` 为 `true` 的最大的 `step` 以及最终状态的编号。
func (d *DoublingSimple) MaxStep(from int32, check func(next int32) bool) (step int, to int32) {
	for k := d.log - 1; k >= 0; k-- {
		tmp := d.to[k*d.n+from]
		if tmp == -1 {
			continue
		}
		if check(tmp) {
			step |= 1 << k
			from = tmp
		}
	}
	to = from
	return
}

type Bitset []uint

func NewBitset(n int32) Bitset    { return make(Bitset, n>>6+1) }
func (b Bitset) Set(p int32)      { b[p>>6] |= 1 << (p & 63) }
func (b Bitset) Has(p int32) bool { return b[p>>6]&(1<<(p&63)) != 0 }
func (b Bitset) Reset(p int32)    { b[p>>6] &^= 1 << (p & 63) }
func (b Bitset) Flip(p int32)     { b[p>>6] ^= 1 << (p & 63) }

func abs32(a int32) int32 {
	if a < 0 {
		return -a
	}
	return a
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
	// P3181()
	// P4081()
	// P6139()

	// bzoj3926()
	// SP8093()

	// CF204E()
	// CF204E线段树合并()
	// CF316G3()
	// CF427D()
	// CF452E()
	CF547E()
	// CF666E()
}

// P3181 [HAOI2016] 找相同字符 (分别维护不同串的 size)
// https://www.luogu.com.cn/problem/P3181
// !求两个字符串的相同子串数量。
//
// 输入SAM,求出每个串对应的endPosSize后，按照dfs序逆序dp.
// 每个endPos的贡献为size1*size2*count.
func P3181() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s, t string
	fmt.Fscan(in, &s, &t)

	sam := NewSuffixAutomatonGeneral()
	maxSize := int32(2 * (len(s) + len(t)))
	isPrefix1 := NewBitset(maxSize)
	sam.AddString(s, func(_, pos int32) { isPrefix1.Set(pos) })
	isPrefix2 := NewBitset(maxSize)
	sam.AddString(t, func(_, pos int32) { isPrefix2.Set(pos) })

	size := sam.Size()
	dfsOrder := sam.GetDfsOrder()
	endPosSize1 := sam.GetEndPosSize(dfsOrder, isPrefix1.Has)
	endPosSize2 := sam.GetEndPosSize(dfsOrder, isPrefix2.Has)
	res := 0
	for i := size - 1; i >= 1; i-- {
		node := sam.Nodes[i]
		count := int(node.MaxLen - sam.Nodes[node.Link].MaxLen)
		size1, size2 := int(endPosSize1[i]), int(endPosSize2[i])
		res += size1 * size2 * count
	}
	fmt.Fprintln(out, res)
}

// P4081 [USACO17DEC] Standing Out from the Herd P
// https://www.luogu.com.cn/problem/P4081
//
// 给定n个模式串.对每个模式串，求出本质不同的子串个数，且子串不在其他模式串中出现.
// 每个串在自动机上跑一下，然后把经过的点以及他们的parent树上的祖先全部标记上当前串的编号.
// 如果一个点被标记了两次，那么这个点所代表的子串必然不是本质相同的，就不能算
func P4081() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int32
	fmt.Fscan(in, &n)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	sam := NewSuffixAutomatonGeneral()
	for _, v := range words {
		sam.AddString(v, nil)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belong := make([][]int32, size)
	visitedTime := make([]int32, size)
	for i := int32(0); i < size; i++ {
		visitedTime[i] = -1
	}

	// 对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.
	// 标记次数之和不超过O(Lsqrt(L)).
	markChain := func(sid int32, pos int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			if len(belong[pos]) >= 2 {
				break
			}
			belong[pos] = append(belong[pos], sid)
			pos = nodes[pos].Link
		}
	}

	// 标记所有文本串的子串.
	for i, w := range words {
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos)
		}
	}

	res := make([]int32, n)
	for i := int32(1); i < size; i++ {
		if len(belong[i]) == 1 {
			res[belong[i][0]] += sam.DistinctSubstringAt(i)
		}
	}

	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// P6139 【模板】广义后缀自动机（广义 SAM）
// https://www.luogu.com.cn/problem/P6139
// 求多个字符串的本质不同子串个数.
// 同时需要输出广义后缀自动机的点数.
func P6139() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	sam := NewSuffixAutomatonGeneral()
	for i := 0; i < n; i++ {
		var s string
		fmt.Fscan(in, &s)
		sam.AddString(s, nil)
	}
	fmt.Fprintln(out, sam.DistinctSubstring())
	fmt.Fprintln(out, sam.Size())
}

// bzoj 3926 [Zjoi2015]诸神眷顾的幻想乡 (树上本质不同路径数)
// 给出一颗叶子结点不超过 20 个的无根树，每个节点上都有一个不超过 10 的数字.
// 求树上本质不同的路径个数（两条路径相同定义为：其路径上所有节点上的数字依次相连组成的字符串相同）。
//
// 如果只建立一个SAM，会出现路径不是一条链(被 LCA 折断)的情况，不方便统计.
// !由于叶子节点仅有20个，因此从每个叶子节点开始，整棵树都会形成一个字典树。将这 棵 Trie 树拼在一起求 GSAM 即可。
func bzoj3926() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	colors := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &colors[i])
	}
	tree := make([][]int32, n)
	deg := make([]int32, n)
	for i := int32(0); i < n-1; i++ {
		var u, v int32
		fmt.Fscan(in, &u, &v)
		u--
		v--
		tree[u] = append(tree[u], v)
		tree[v] = append(tree[v], u)
		deg[u]++
		deg[v]++
	}

	sam := NewSuffixAutomatonGeneral()
	var dfs func(cur, pre, pos int32)
	dfs = func(cur, pre, pos int32) {
		pos = sam.Add(pos, colors[cur])
		for _, next := range tree[cur] {
			if next != pre {
				dfs(next, cur, pos)
			}
		}
	}
	for i := int32(0); i < n; i++ {
		if deg[i] == 1 {
			dfs(i, -1, 0)
		}
	}

	fmt.Fprintln(out, sam.DistinctSubstring())
}

// JZPGYZ - Sevenk Love Oimaster
// https://www.luogu.com.cn/problem/SP8093
// 给定 n 个模板串，以及 q 个查询串
// 依次查询每一个查询串是多少个模板串的子串
//
// 对原串建一个广义SAM，然后把每一个原串放到SAM上跑一跑，记录一下每一个状态属于多少个原串，用belongCount表示。
// 这样的话查询串直接在SAM上跑，如果失配输出0，否则直接输出记录在上面的belongCount就好了。
func SP8093() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	texts := make([]string, n)
	for i := 0; i < n; i++ {
		fmt.Fscan(in, &texts[i])
	}
	patterns := make([]string, q)
	for i := 0; i < q; i++ {
		fmt.Fscan(in, &patterns[i])
	}

	sam := NewSuffixAutomatonGeneral()
	for _, v := range texts {
		sam.AddString(v, nil)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belongCount := make([]int32, size) // 每个状态属于多少个原串
	visitedTime := make([]int32, size)
	for i := range visitedTime {
		visitedTime[i] = -1
	}

	// 对文本串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.
	// 标记次数之和不超过O(Lsqrt(L)).
	markChain := func(sid int32, pos int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			belongCount[pos]++
			pos = nodes[pos].Link
		}
	}

	// 查询模式串是多少个文本串的子串.
	query := func(pattern string) int32 {
		pos := int32(0)
		for _, c := range pattern {
			pos = nodes[pos].Next[c-OFFSET]
			if pos == -1 {
				return 0
			}
		}
		return belongCount[pos]
	}

	// 标记所有文本串的子串.
	for i, w := range texts {
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos)
		}
	}

	for _, w := range patterns {
		fmt.Fprintln(out, query(w))
	}
}

// Little Elephant and Strings (广义SAM根号技巧 + 记忆化dfs预处理)
// https://www.luogu.com.cn/problem/CF204E
// 给定n个字符串，对每个字符串，求出它有多少个子串属于至少k个字符串.
//
// - 标记时，使用 visitedTime 保证一个字符串的多个子串对同一个endPos只标记一次，复杂度O(Lsqrt(L)).
// - 查询时，先预处理出合法的endPosCount个数，再一边转移一边查询.
func CF204E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
	}

	sam := NewSuffixAutomatonGeneral()
	for _, v := range words {
		sam.AddString(v, nil)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belongCount := make([]int32, size) // 每个状态属于多少个原串
	visitedTime := make([]int32, size)

	// 对字符串t[i]的每个前缀，在后缀链接树上向上跳标记每个endPos，表示该endPos包含了t[i]的子串.
	// 标记次数之和不超过O(Lsqrt(L)).
	markChain := func(sid int32, pos int32) {
		for pos > 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			belongCount[pos]++
			pos = nodes[pos].Link
		}
	}

	for i := int32(0); i < size; i++ {
		visitedTime[i] = -1
	}
	// 标记所有文本串的子串.
	for i, w := range words {
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos)
		}
	}

	// 查询每个字符串有多少个子串属于至少k个字符串.
	// 注意不能暴力上跳，需要预处理从每个节点一直上跳获得的总合法子串个数.
	memo := make([]int, size)
	visited := make([]bool, size)
	var dfs func(int32)
	dfs = func(pos int32) {
		if pos == 0 || visited[pos] {
			return
		}
		visited[pos] = true
		link := nodes[pos].Link
		dfs(link)
		memo[pos] += memo[link]
	}
	for i := int32(1); i < size; i++ {
		if belongCount[i] >= k {
			memo[i] = int(sam.DistinctSubstringAt(i))
		}
	}
	for i := int32(1); i < size; i++ {
		dfs(i)
	}

	res := make([]int, n)
	for i, w := range words {
		count := 0
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			count += memo[pos]
		}
		res[i] = count
	}
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// 将CF204E的markChain函数换成线段树合并处理.
// 线段树合并更慢一些.
func CF204E线段树合并() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, k int32
	fmt.Fscan(in, &n, &k)
	words := make([]string, n)
	allLen := int32(0)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
		allLen += int32(len(words[i]))
	}

	sam := NewSuffixAutomatonGeneral()
	// func e() E { return 0 }
	// func op(a, b E) E {
	// 	return a + b
	// }
	// func merge(a, b E) E { // 合并两个不同的树的结点的函数
	// 	return min32(1, a+b)
	// }
	seg := NewSegmentTreeMerger(0, n-1) // 维护每个endPos出现的字符串下标.
	nodes := make([]*SegNode, allLen*2)
	for i := range nodes {
		nodes[i] = seg.Alloc()
	}
	for wi, v := range words {
		pos := int32(0)
		for _, c := range v {
			pos = sam.Add(pos, c)
			seg.Set(nodes[pos], int32(wi), 1)
		}
	}

	size := sam.Size()
	tree := sam.BuildTree()
	isEndPosOk := make([]bool, size) // 状态属于>=k个原串则为ok.
	var mergeEndPos func(int32)
	mergeEndPos = func(cur int32) {
		for _, next := range tree[cur] {
			mergeEndPos(next)
			nodes[cur] = seg.MergeDestructively(nodes[cur], nodes[next]) // 线段树合并
		}
		isEndPosOk[cur] = seg.QueryAll(nodes[cur]) >= k
	}
	mergeEndPos(0)

	// 查询每个字符串有多少个子串属于至少k个字符串.
	// 注意不能暴力上跳，需要预处理从每个节点一直上跳获得的总合法子串个数.
	memo := make([]int, size)
	visited := make([]bool, size)
	var dfs func(int32)
	dfs = func(pos int32) {
		if pos == 0 || visited[pos] {
			return
		}
		visited[pos] = true
		link := sam.Nodes[pos].Link
		dfs(link)
		memo[pos] += memo[link]
	}
	for i := int32(1); i < size; i++ {
		if isEndPosOk[i] {
			memo[i] = int(sam.DistinctSubstringAt(i))
		}
	}
	for i := int32(1); i < size; i++ {
		dfs(i)
	}

	res := make([]int, n)
	for i, w := range words {
		count := 0
		pos := int32(0)
		for _, c := range w {
			pos = sam.Nodes[pos].Next[c-OFFSET]
			count += memo[pos]
		}
		res[i] = count
	}
	for _, v := range res {
		fmt.Fprint(out, v, " ")
	}
}

// Good Substrings
// https://www.luogu.com.cn/problem/CF316G3
// 给定一个文本串s和m个限制，问有多少个本质不同子串s'满足所有限制：
// 每个限制形如(模式串w,left,right)：在w中，s'的出现次数在[left,right]之间.
// m<=10,len(w)<=5e4
//
// 看到本质不同子串，文本串和模式串全部丢到广义 SAM 上。
// 求出endPos中每种串的出现次数即可.
func CF316G3() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var m int32
	fmt.Fscan(in, &m)
	type limit struct {
		w           string
		left, right int32
	}
	limits := make([]limit, m)
	for i := int32(0); i < m; i++ {
		var w string
		var left, right int32
		fmt.Fscan(in, &w, &left, &right)
		limits[i] = limit{w, left, right}
	}

	N := int32(len(s))
	for _, v := range limits {
		N += int32(len(v.w))
	}
	isPrefix := make([]Bitset, m+1)
	for i := range isPrefix {
		isPrefix[i] = NewBitset(2 * N)
	}
	sam := NewSuffixAutomatonGeneral()
	insert := func(s string, id int32) {
		pos := int32(0)
		for _, c := range s {
			pos = sam.Add(pos, c)
			isPrefix[id].Set(pos)
		}
	}

	insert(s, 0)
	for i, v := range limits {
		insert(v.w, int32(i+1))
	}

	dfsOrder := sam.GetDfsOrder()
	endPosSize := make([][]int32, m+1)
	for i := int32(0); i <= m; i++ {
		endPosSize[i] = sam.GetEndPosSize(dfsOrder, isPrefix[i].Has)
	}

	res := 0
	check := func(pos int32) bool {
		for i, v := range limits {
			if endPosSize[i+1][pos] < v.left || endPosSize[i+1][pos] > v.right {
				return false
			}
		}
		return true
	}
	for i := int32(1); i < sam.Size(); i++ {
		if endPosSize[0][i] > 0 && check(i) {
			res += int(sam.DistinctSubstringAt(i))
		}
	}
	fmt.Fprintln(out, res)
}

// Match & Catch
// https://www.luogu.com.cn/problem/CF427D
// 求两个串的最短公共唯一子串
func CF427D() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const INF int32 = 1e9 + 10

	var s, t string
	fmt.Fscan(in, &s, &t)

	sam := NewSuffixAutomatonGeneral()
	maxSize := int32(2 * (len(s) + len(t)))
	isPrefix1 := NewBitset(maxSize)
	sam.AddString(s, func(_, pos int32) { isPrefix1.Set(pos) })
	isPrefix2 := NewBitset(maxSize)
	sam.AddString(t, func(_, pos int32) { isPrefix2.Set(pos) })

	size := sam.Size()
	dfsOrder := sam.GetDfsOrder()
	endPosSize1 := sam.GetEndPosSize(dfsOrder, isPrefix1.Has)
	endPosSize2 := sam.GetEndPosSize(dfsOrder, isPrefix2.Has)

	res := INF
	for i := int32(1); i < size; i++ {
		if endPosSize1[i] == 1 && endPosSize2[i] == 1 {
			node := sam.Nodes[i]
			minLen := sam.Nodes[node.Link].MaxLen + 1
			res = min32(res, minLen)
		}
	}
	if res == INF {
		fmt.Fprintln(out, -1)
	} else {
		fmt.Fprintln(out, res)
	}
}

// Three strings
// https://www.luogu.com.cn/problem/CF452E
//
// 给定三个字符串s1,s2,s3，记最短串长度为L.
// 对每个1<=len<=L，
// 求三元组(i,j,k)的个数模1e9+7，满足s1[i:i+len-1]=s2[j:j+len-1]=s3[k:k+len-1].
//
// 先建一棵广义SAM，求出每个点可以到达的A,B,C的字串的个数，
// 然后这个点贡献的值就是三个串分别的个数乘起来，发现一个点会对[minLen,maxLen]的答案都有贡献。
// 差分求解即可.
func CF452E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	const MOD int = 1e9 + 7

	var s1, s2, s3 string
	fmt.Fscan(in, &s1, &s2, &s3)

	L := min32(min32(int32(len(s1)), int32(len(s2))), int32(len(s3)))

	maxSize := int32(2 * (len(s1) + len(s2) + len(s3)))
	isPrefix := [3]Bitset{NewBitset(maxSize), NewBitset(maxSize), NewBitset(maxSize)}
	sam := NewSuffixAutomatonGeneral()
	for i, s := range []string{s1, s2, s3} {
		pos := int32(0)
		for _, c := range s {
			pos = sam.Add(pos, c)
			isPrefix[i].Set(pos)
		}
	}
	size := sam.Size()
	dfsOrder := sam.GetDfsOrder()
	endPosSize := [3][]int32{nil, nil, nil}
	for i := range endPosSize {
		endPosSize[i] = sam.GetEndPosSize(dfsOrder, isPrefix[i].Has)
	}

	diff := make([]int, size+1)
	for i := int32(1); i < size; i++ {
		mul := int(endPosSize[0][i]) * int(endPosSize[1][i]) % MOD * int(endPosSize[2][i]) % MOD
		link := sam.Nodes[i].Link
		minLen, maxLen := sam.Nodes[link].MaxLen+1, sam.Nodes[i].MaxLen
		diff[minLen] = (diff[minLen] + mul) % MOD
		diff[maxLen+1] = (diff[maxLen+1] - mul + MOD) % MOD
	}
	for i := int32(1); i <= L; i++ {
		diff[i] = (diff[i] + diff[i-1]) % MOD
	}
	for i := int32(1); i <= L; i++ {
		fmt.Fprint(out, diff[i], " ")
	}
}

// Mike and Friends
// https://www.luogu.com.cn/problem/CF547E
// 给定n个字符串words和q个查询.
// 每次查询words[index]在words[start:end)中出现的次数.
func CF547E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	allLen := int32(0)
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
		allLen += int32(len(words[i]))
	}

	sam := NewSuffixAutomatonGeneral()
	seg := NewSegmentTreeMerger(0, n-1)
	nodes := make([]*SegNode, 2*allLen)
	for i := range nodes {
		nodes[i] = seg.Alloc()
	}

	endPos := make([]int32, n) // 每个串的pos
	for wi, w := range words {
		last := sam.AddString(w, func(_, pos int32) { seg.Set(nodes[pos], int32(wi), 1) })
		endPos[wi] = last
	}
	dfsOrder := sam.GetDfsOrder()
	for i := sam.Size() - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		link := sam.Nodes[cur].Link
		nodes[link] = seg.Merge(nodes[link], nodes[cur])
	}

	query := func(start, end, index int32) int32 {
		pos := endPos[index]
		return seg.Query(nodes[pos], start, end-1)
	}

	for i := int32(0); i < q; i++ {
		var start, end, index int32
		fmt.Fscan(in, &start, &end, &index)
		start--
		index--
		fmt.Fprintln(out, query(start, end, index))
	}
}

// Forensic Examination [CF666E] (线段树合并维护 endPosSize)
// https://www.luogu.com.cn/problem/CF666E
// 给定一个字符串s和n个字符串words[i].
// 处理q个查询，每次查询s[start1:end1)在words[start2:end2)的哪个串里出现次数最多.
// 如果出现次数相同，返回最小的下标，并输出出现次数.
//
// 1.所有串建立广义SAM.
// 2.对每个查询，倍增定位s[start1:end1)在fail树上的位置.
// 3.利用线段树合并求出每个endPos的(最大出现次数，最小下标)。
func CF666E() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var s string
	fmt.Fscan(in, &s)
	var n int32
	fmt.Fscan(in, &n)

	allLen := int32(len(s))
	words := make([]string, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &words[i])
		allLen += int32(len(words[i]))
	}

	sam := NewSuffixAutomatonGeneral()
	prefixEnd := make([]int32, len(s))
	sam.AddString(s, func(i, pos int32) { prefixEnd[i] = pos })

	seg := NewSegmentTreeOnRangeWithIndex(0, n-1)
	nodes := make([]*SegNode2, 2*allLen)
	for i := range nodes {
		nodes[i] = seg.Alloc()
	}
	for wi, w := range words {
		sam.AddString(w, func(_, pos int32) { seg.Set(nodes[pos], int32(wi), 1) })
	}
	dfsOrder := sam.GetDfsOrder()
	for i := sam.Size() - 1; i >= 1; i-- {
		cur := dfsOrder[i]
		link := sam.Nodes[cur].Link
		nodes[link] = seg.Merge(nodes[link], nodes[cur]) // 注意不能用MergeDestructively
	}

	// 每次查询s[start1:end1)在words[start2:end2)的哪个串里出现次数最多.
	// 如果出现次数相同，返回最小的下标，并输出出现次数.
	query := func(start1, end1 int32, start2, end2 int32) (maxCount int32, minIndex int32) {
		pos := sam.LocateSubstring(start1, end1, prefixEnd[end1-1])
		maxCount, minIndex = seg.Query(nodes[pos], start2, end2-1)
		return
	}

	var q int32
	fmt.Fscan(in, &q)
	for i := int32(0); i < q; i++ {
		var a, b, c, d int32
		fmt.Fscan(in, &a, &b, &c, &d)
		a--
		c--
		maxCount, minIndex := query(c, d, a, b)

		// !注意不存在的情况
		if maxCount == 0 {
			fmt.Fprintln(out, a+1, 0)
		} else {
			fmt.Fprintln(out, minIndex+1, maxCount)
		}
	}
}

// 1408. 数组中的字符串匹配
// https://leetcode.cn/problems/string-matching-in-an-array/description/
// 对每个单词，如果它是另一个单词的子串，就把它加入答案中.
func stringMatching(words []string) (res []string) {
	sam := NewSuffixAutomatonGeneral()
	endPos := make([]int32, len(words)) // 每个串的pos
	for i, w := range words {
		endPos[i] = sam.AddString(w, nil)
	}

	size := sam.Size()
	nodes := sam.Nodes
	belongCount := make([]int32, size) // 每个状态属于多少个原串
	visitedTime := make([]int32, size)
	for i := range visitedTime {
		visitedTime[i] = -1
	}

	markChain := func(sid int32, pos int32) {
		for pos >= 0 && visitedTime[pos] != sid {
			visitedTime[pos] = sid
			belongCount[pos]++
			pos = nodes[pos].Link
		}
	}
	for i, w := range words {
		pos := int32(0)
		for _, c := range w {
			pos = nodes[pos].Next[c-OFFSET]
			markChain(int32(i), pos)
		}
	}

	for i, pos := range endPos {
		if belongCount[pos] >= 2 {
			res = append(res, words[i])
		}
	}
	return
}
