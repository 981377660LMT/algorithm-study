// TODO 有问题

// https://www.luogu.com.cn/problem/P5055
// P5055 【模板】可持久化文艺平衡树
// https://www.luogu.com.cn/problem/solution/P5055

package main

import (
	"bufio"
	"fmt"
	"os"
	"time"
)

// 您需要写一种数据结构，来维护一个序列，
// 其中需要提供以下操作（对于各个以往的历史版本）：

// 1.在第 p 个数后插入数 x 。
// 2.删除第 p 个数。
// 3.翻转区间 [l,r]，例如原序列是 {5,4,3,2,1}，翻转区间 [2,4]后，
// 结果是 {5,2,3,4,1}。
// 4.查询区间 [l,r] 中所有数的和。
// 和原本平衡树不同的一点是，每一次的任何操作都是基于某一个历史版本，
// 同时生成一个新的版本（操作 4 即保持原版本无变化），
// 新版本即编号为此次操作的序号。

// 本题强制在线。
// 第一行包含一个整数n，表示操作的总数。
// 接下来n行，每行前两个整数vi, opt，vi表示基于的过去版本号(0<vi<i) , opt,表示操作的序号(1≤opti≤ 4)。
// 若opti= 1，则接下来两个整数pi,xi，表示操作为在第pi个数后插入数x。
// 若opti= 2，则接下来一个整数pi，表示操作为删除第pi个数。
// 若opti= 3，则接下来两个整数li,ri，表示操作为翻转区间[li, ri]。
// 若opt,= 4，则接下来两个整数li,ri，表示操作为查询区间[li, ri]的和。
// 强制在线规则:
// 令当前操作之前的最后一次4操作的答案为lastRes(如果之前没有4操作，则lastRes = 0)。
// 则此次操作的pi, xi或li, ri均按位异或上lastRes即可得到真实的pi,xi或li, ri。

// !注：如果用的是指针写法，必要时禁止 GC，能加速不少
// !但是如果是多组数据，就不要禁止 GC 了，否则MLE
// func init() { debug.SetGCPercent(-1) }

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n int
	fmt.Fscan(in, &n)
	tree := NewPersistentFHQTreap(n * 60) // !开60倍空间避免扩容

	versions := make([]int, n+1) // 每次操作后的根节点版本号
	lastRes := 0
	var version, opt, pos, value, left, right int
	for i := 1; i <= n; i++ {
		fmt.Fscan(in, &version, &opt)
		versions[i] = versions[version]
		switch opt {
		case 1:
			fmt.Fscan(in, &pos, &value)
			pos ^= lastRes
			value ^= lastRes
			newVersion := tree.Insert(versions[i], pos, value) // 新的根节点的版本编号
			versions[i] = newVersion
		case 2:
			fmt.Fscan(in, &pos)
			pos ^= lastRes
			newVersion := tree.Pop(versions[i], pos-1)
			versions[i] = newVersion
		case 3:
			fmt.Fscan(in, &left, &right)
			left ^= lastRes
			right ^= lastRes
			newVersion := tree.Reverse(versions[i], left-1, right)
			versions[i] = newVersion
		case 4:
			fmt.Fscan(in, &left, &right)
			left ^= lastRes
			right ^= lastRes
			lastRes = tree.Query(versions[i], left-1, right)
			fmt.Fprintln(out, lastRes)
		}
	}
}

type Node struct {
	// !Raw value
	element int

	// !Data and lazy tag maintained by segment tree
	sum     int
	lazyAdd int

	left, right int
	size        int
	isReversed  uint8
}

type PersistentFHQTreap struct {
	root      int // 当前(版本)的根节点编号
	seed      uint64
	nodeCount int
	nodes     []Node // !不用指针会快很多(但是拷贝会内存占用更高)
}

// initCapacity 一般是 操作数的60倍
func NewPersistentFHQTreap(initCapacity int) *PersistentFHQTreap {
	treap := &PersistentFHQTreap{
		seed:  uint64(time.Now().UnixNano()/2 + 1),
		nodes: make([]Node, max(initCapacity, 128)),
	}
	return treap
}

// 返回build后的根节点版本编号
func (t *PersistentFHQTreap) Build(nums []int) int {
	t.root = t.build(1, len(nums), nums)
	return t.root
}

func (pt *PersistentFHQTreap) build(left, right int, nums []int) int {
	if left > right {
		return 0
	}
	mid := (left + right) >> 1
	newNode := pt.newNode(nums[mid-1])
	pt.nodes[newNode].left = pt.build(left, mid-1, nums)
	pt.nodes[newNode].right = pt.build(mid+1, right, nums)
	pt.pushUp(newNode)
	return newNode
}

func (pt *PersistentFHQTreap) newNode(value int) int {
	pt.growBy(1)
	pt.nodeCount++
	pt.nodes[pt.nodeCount].element = value
	pt.nodes[pt.nodeCount].sum = value
	pt.nodes[pt.nodeCount].size = 1
	return pt.nodeCount
}

func (pt *PersistentFHQTreap) copyNode(node int) int {
	if node == 0 {
		return 0
	}
	pt.growBy(1)
	pt.nodeCount++
	pt.nodes[pt.nodeCount] = pt.nodes[node]
	return pt.nodeCount
}

// OK
func (pt *PersistentFHQTreap) pushUp(root int) {
	if root == 0 {
		return
	}

	pt.nodes[root].size = 1
	pt.nodes[root].sum = pt.nodes[root].element
	if pt.nodes[root].left != 0 {
		pt.nodes[root].size += pt.nodes[pt.nodes[root].left].size
		pt.nodes[root].sum += pt.nodes[pt.nodes[root].left].sum
	}
	if pt.nodes[root].right != 0 {
		pt.nodes[root].size += pt.nodes[pt.nodes[root].right].size
		pt.nodes[root].sum += pt.nodes[pt.nodes[root].right].sum
	}
}

// !OK
func (pt *PersistentFHQTreap) pushDown(root int) {
	if root == 0 {
		return
	}

	if pt.nodes[root].isReversed != 0 {
		pt.nodes[root].left, pt.nodes[root].right = pt.nodes[root].right, pt.nodes[root].left
		if pt.nodes[root].left != 0 {
			pt.nodes[root].left = pt.copyNode(pt.nodes[root].left)
			pt.nodes[pt.nodes[root].left].isReversed ^= 1
		}
		if pt.nodes[root].right != 0 {
			pt.nodes[root].right = pt.copyNode(pt.nodes[root].right)
			pt.nodes[pt.nodes[root].right].isReversed ^= 1
		}
		pt.nodes[root].isReversed = 0
	}
}

func (pt *PersistentFHQTreap) splitByRank(root, k int, x, y *int) {
	if root == 0 {
		*x, *y = 0, 0
		return
	}

	pt.pushDown(root)
	if k <= pt.nodes[pt.nodes[root].left].size {
		*y = pt.copyNode(root)
		pt.splitByRank(pt.nodes[root].left, k, x, &pt.nodes[*y].left)
		pt.pushUp(*y)
	} else {
		*x = pt.copyNode(root)
		pt.splitByRank(pt.nodes[root].right, k-pt.nodes[pt.nodes[root].left].size-1, &pt.nodes[*x].right, y)
		pt.pushUp(*x)
	}
}

// 返回新版本的根节点编号
// 注意merge中无需再拷贝 因为split总是在merge之前调用
func (pt *PersistentFHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x + y
	}

	if int(pt.nextRand()*(uint64(pt.nodes[x].size)+uint64(pt.nodes[y].size))>>32) < pt.nodes[x].size {
		pt.pushDown(x)
		pt.nodes[x].right = pt.merge(pt.nodes[x].right, y)
		pt.pushUp(x)
		return x
	} else {
		pt.pushDown(y)
		pt.nodes[y].left = pt.merge(x, pt.nodes[y].left)
		pt.pushUp(y)
		return y
	}
}

// 插入元素 返回新版本的根节点编号
func (pt *PersistentFHQTreap) Insert(rootVersion int, index int, value int) int {
	var left, right int
	pt.splitByRank(rootVersion, index, &left, &right)
	newRoot := pt.merge(pt.merge(left, pt.newNode(value)), right)
	pt.root = newRoot
	return newRoot
}

func (pt *PersistentFHQTreap) Pop(rootVersion int, index int) int {
	index++
	var a, b, c int
	pt.splitByRank(rootVersion, index, &b, &c)
	pt.splitByRank(b, index-1, &a, &b)
	newRoot := pt.merge(a, c)
	pt.root = newRoot
	return newRoot
}

// Reverse [start, stop) in place.
func (pt *PersistentFHQTreap) Reverse(rootVersion int, left, right int) int {
	left++
	var a, b, c int
	pt.splitByRank(rootVersion, right, &a, &c)
	pt.splitByRank(a, left-1, &a, &b)
	if b != 0 {
		pt.nodes[b].isReversed ^= 1
	}
	newRoot := pt.merge(pt.merge(a, b), c)
	pt.root = newRoot
	return newRoot
}

// Query [start, stop) (defaults to range sum).
//  0 <= start <= stop <= n
func (pt *PersistentFHQTreap) Query(root int, left, right int) int {
	left++
	var a, b, c int
	pt.splitByRank(root, right, &a, &c)
	pt.splitByRank(a, left-1, &a, &b)
	res := 0
	if b != 0 {
		res = pt.nodes[b].sum
	}
	pt.root = pt.merge(a, pt.merge(b, c))
	return res
}

// Return all elements in index order.
func (pt *PersistentFHQTreap) InOrder(rootVersion int) []int {
	res := make([]int, 0, pt.Size(rootVersion))
	pt.inOrder(rootVersion, &res)
	return res
}

func (pt *PersistentFHQTreap) inOrder(root int, res *[]int) {
	if root == 0 {
		return
	}
	pt.pushDown(root)
	pt.inOrder(pt.nodes[root].left, res)
	*res = append(*res, pt.nodes[root].element)
	pt.inOrder(pt.nodes[root].right, res)
}

func (pt *PersistentFHQTreap) Size(rootVersion int) int {
	return pt.nodes[rootVersion].size
}

func (pt *PersistentFHQTreap) nextRand() uint64 {
	pt.seed ^= pt.seed << 7
	pt.seed ^= pt.seed >> 9
	return pt.seed & 0xFFFFFFFF
}

func (pt *PersistentFHQTreap) growBy(n int) {
	currentCapacity := cap(pt.nodes)
	if pt.nodeCount+n >= currentCapacity {
		newCapacity := int(2.0 * float32(currentCapacity+n))
		pt.resize(newCapacity)
	}
}

func (pt *PersistentFHQTreap) resize(cap int) {
	newNodes := make([]Node, cap, cap)
	copy(newNodes, pt.nodes)
	pt.nodes = newNodes
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
