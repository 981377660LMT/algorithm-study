// 线段树套01Trie的在线做法
// P4585 [FJOI2015]火星商店问题
// https://chenyu-w.github.io./2023/05/01/P4585%20%E7%81%AB%E6%98%9F%E5%95%86%E5%BA%97%E9%97%AE%E9%A2%98/
// 给定 n 个集合，每个集合元素有两个值，一个是价值，一个是存在时间，每个集合初始有一个存在时间无限的物品。
// 每天都有一个 1 操作和若干个 2 操作
// 操作 0 index v : 在编号为index的集合中加入一个物品v。
// 操作 1 left right x day : 在 left 到 right 集合内查询未过期的物品(day天之内)，使 value xor x 最大，输出最大值。

package main

import (
	"bufio"
	"fmt"
	"math/bits"
	"os"
)

const INF int = 1e18

func main() {
	in, out := bufio.NewReader(os.Stdin), bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)
	values := make([]int, n) // 每个商店的不下架物品的价值
	for i := range values {
		fmt.Fscan(in, &values[i])
	}

	operations := make([][4]int, q)
	mutations := [][2]int{} // (id, time)
	queries := [][3]int{}   // (id, startTime, endTime)
	curTime := 0
	for i := 0; i < q; i++ {
		var op int
		fmt.Fscan(in, &op)
		if op == 0 {
			var shop, price int
			fmt.Fscan(in, &shop, &price)
			shop--
			curTime += 1
			operations[i] = [4]int{shop, price, 0, 0}
			mutations = append(mutations, [2]int{i, curTime})
		} else {
			var start, end, x, day int
			fmt.Fscan(in, &start, &end, &x, &day)
			start--
			operations[i] = [4]int{start, end, x, day}
			queries = append(queries, [3]int{i, curTime - day, curTime + 1})
		}
	}

}

type InnerTree = Node

// 线段树套树.
type SegmentTreeDivideInterval struct {
	n               int
	smallN          bool
	offset          int // 线段树中一共offset+n个节点,offset+i对应原来的第i个节点.
	createInnerTree func() InnerTree
	innerTreeList   []InnerTree
	innerTreeMap    map[int]InnerTree
}

// 单点更新，区间查询的树套树.
// n: 第一个维度(一般是序列)的长度.
// createInnerTree: 创建第二个维度(一般是线段树)的树.
// smallN: n较小时，会预先创建好所有的内层树; 否则会用map保存内层树，并在需要的时候创建.
func NewSegmentTree2DPointUpdateRangeQuery(n int, createInnerTree func() InnerTree, smallN bool) *SegmentTreeDivideInterval {
	offset := 1
	for offset < n {
		offset <<= 1
	}
	var innerTreeList []InnerTree
	if smallN {
		innerTreeList = make([]InnerTree, offset+n)
		for i := range innerTreeList {
			innerTreeList[i] = createInnerTree()
		}
	}
	return &SegmentTreeDivideInterval{
		n:               n,
		smallN:          smallN,
		offset:          offset,
		createInnerTree: createInnerTree,
		innerTreeList:   innerTreeList,
		innerTreeMap:    map[int]InnerTree{},
	}
}

func (tree *SegmentTreeDivideInterval) EnumeratePoint(index int, f func(tree InnerTree)) {
	if index < 0 || index >= tree.n {
		return
	}
	index += tree.offset
	for index > 0 {
		f(tree.getTree(index))
		index >>= 1
	}
}

func (tree *SegmentTreeDivideInterval) EnumerateRange(start, end int, f func(tree InnerTree)) {
	if start < 0 {
		start = 0
	}
	if end > tree.n {
		end = tree.n
	}
	if start >= end {
		return
	}

	leftSegments := []InnerTree{}
	rightSegments := []InnerTree{}
	for start, end = start+tree.offset, end+tree.offset; start < end; start, end = start>>1, end>>1 {
		if start&1 == 1 {
			leftSegments = append(leftSegments, tree.getTree(start))
			start++
		}
		if end&1 == 1 {
			end--
			rightSegments = append(rightSegments, tree.getTree(end))
		}
	}

	for i := 0; i < len(leftSegments); i++ {
		f(leftSegments[i])
	}
	for i := len(rightSegments) - 1; i >= 0; i-- {
		f(rightSegments[i])
	}
}

func (tree *SegmentTreeDivideInterval) getTree(segmentId int) InnerTree {
	if tree.smallN {
		return tree.innerTreeList[segmentId]
	} else {
		if v, ok := tree.innerTreeMap[segmentId]; ok {
			return v
		} else {
			newTree := tree.createInnerTree()
			tree.innerTreeMap[segmentId] = newTree
			return newTree
		}
	}
}

type Node struct {
	lastIndex int // 最后一次被更新的时间
	chidlren  [2]*Node
}

type XorTrieSimplePersistent struct {
	bit int
}

func NewXorTrieSimplePersistent(upper int) *XorTrieSimplePersistent {
	return &XorTrieSimplePersistent{bit: bits.Len(uint(upper))}
}

func (trie *XorTrieSimplePersistent) NewRoot() *Node {
	return nil
}

func (trie *XorTrieSimplePersistent) Copy(node *Node) *Node {
	if node == nil {
		return node
	}
	return &Node{
		lastIndex: node.lastIndex,
		chidlren:  node.chidlren,
	}
}

func (trie *XorTrieSimplePersistent) Insert(root *Node, num int, lastIndex int) *Node {
	if root == nil {
		root = &Node{}
	}
	return trie._insert(root, num, trie.bit-1, lastIndex)
}

// 查询num与root中的数异或的最大值以及最大值对应的结点.
// !如果root为nil,返回0.
func (trie *XorTrieSimplePersistent) Query(root *Node, num int, leftIndex int) (maxXor int, node *Node) {
	if root == nil {
		return
	}
	for k := trie.bit - 1; k >= 0; k-- {
		bit := (num >> k) & 1
		if root.chidlren[bit^1] != nil && root.chidlren[bit^1].lastIndex >= leftIndex {
			bit ^= 1
			maxXor |= 1 << k
		}
		root = root.chidlren[bit]
	}
	return maxXor, root
}

func (trie *XorTrieSimplePersistent) _insert(root *Node, num int, depth int, lastIndex int) *Node {
	root = trie.Copy(root)
	root.lastIndex = lastIndex
	if depth < 0 {
		return root
	}
	bit := (num >> depth) & 1
	if root.chidlren[bit] == nil {
		root.chidlren[bit] = &Node{}
	}
	root.chidlren[bit] = trie._insert(root.chidlren[bit], num, depth-1, lastIndex)
	return root
}
