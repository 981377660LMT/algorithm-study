// https://www.luogu.com.cn/problem/P5055
// P5055 【模板】可持久化文艺平衡树
// https://www.luogu.com.cn/problem/solution/P5055

package main

import (
	"fmt"
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
	// in := bufio.NewReader(os.Stdin)
	// out := bufio.NewWriter(os.Stdout)
	// defer out.Flush()

	// var n int
	// fmt.Fscan(in, &n)
	// tree := NewPersistentFHQTreap(n * 50) // 没有Build时 根为dummy结点0

	// versions := make([]int, n+1) // 每次操作后的根节点版本号
	// lastRes := 0
	// var version, opt, pos, value, left, right int
	// for i := 1; i <= n; i++ {
	// 	fmt.Fscan(in, &version, &opt)
	// 	versions[i] = versions[version]
	// 	switch opt {
	// 	case 1:
	// 		fmt.Fscan(in, &pos, &value)
	// 		pos ^= lastRes
	// 		value ^= lastRes
	// 		newVersion := tree.Insert(versions[i], pos, value) // !新的根节点的版本编号
	// 		versions[i] = newVersion
	// 	case 2:
	// 		fmt.Fscan(in, &pos)
	// 		pos ^= lastRes
	// 		newVersion := tree.Pop(versions[i], pos-1)
	// 		versions[i] = newVersion
	// 	case 3:
	// 		fmt.Fscan(in, &left, &right)
	// 		left ^= lastRes
	// 		right ^= lastRes
	// 		newVersion := tree.Reverse(versions[i], left, right)
	// 		versions[i] = newVersion
	// 	case 4:
	// 		fmt.Fscan(in, &left, &right)
	// 		left ^= lastRes
	// 		right ^= lastRes
	// 		lastRes = tree.Query(versions[i], left-1, right)
	// 		fmt.Fprintln(out, lastRes)
	// 	}
	// }
	tree := NewPersistentFHQTreap(16)
	id1 := tree.Build([]int{1, 2, 3, 4, 5})
	fmt.Println(tree.InOrder(id1))
	id2 := tree.Reverse(id1, 1, 3)
	fmt.Println(tree.Query(id2, 1, 3))
	fmt.Println(tree.InOrder(id2))
	fmt.Println(tree.Query(id1, 1, 4))
	// id2 := tree.Insert(id1, 2, 6)
	fmt.Println(tree.InOrder(id2))
	fmt.Println(tree.Query(id2, 1, 4))

	// insert有问题？
	id3 := tree.Insert(id2, 2, 3)
	fmt.Println(tree.InOrder(id3))
	fmt.Println(tree.Query(id3, 1, 4))
	id4 := tree.Pop(id3, 2)
	fmt.Println(tree.InOrder(id4))
	fmt.Println(tree.Query(id4, 1, 4))
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
	root  int // 当前(版本)的根节点编号
	seed  uint
	nodes []Node // !不用指针会快很多 优先不用指针 MLE了再换指针
}

var seed uint64 = uint64(time.Now().UnixNano()/2 + 1)

func nextRand() uint64 {
	seed ^= seed << 7
	seed ^= seed >> 9
	return seed & 0xFFFFFFFF
}

// initCapacity 一般是 操作数的50倍
func NewPersistentFHQTreap(initCapacity int) *PersistentFHQTreap {
	treap := &PersistentFHQTreap{
		seed:  uint(time.Now().UnixNano()/2 + 1),
		nodes: make([]Node, 0, initCapacity),
	}
	dummy := &Node{size: 0}
	treap.nodes = append(treap.nodes, *dummy)
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
	res := &Node{
		element: value,
		sum:     value,
		size:    1,
	}
	pt.nodes = append(pt.nodes, *res)
	return len(pt.nodes) - 1 // !返回新节点的编号(当前在nodes中的下标)
}

func (pt *PersistentFHQTreap) copyNode(node int) int {
	nodeCopy := pt.nodes[node] // !赋值浅拷贝结构体
	pt.nodes = append(pt.nodes, nodeCopy)
	return len(pt.nodes) - 1
}

func (pt *PersistentFHQTreap) pushUp(root int) {
	if root == 0 {
		return
	}
	rootRef := &pt.nodes[root] // !注意如果用非指针,记得加上引用
	rootRef.size = 1
	rootRef.sum = rootRef.element
	if rootRef.left != 0 {
		rootRef.size += pt.nodes[rootRef.left].size
		rootRef.sum += pt.nodes[rootRef.left].sum
	}
	if rootRef.right != 0 {
		rootRef.size += pt.nodes[rootRef.right].size
		rootRef.sum += pt.nodes[rootRef.right].sum
	}
}

func (pt *PersistentFHQTreap) pushDown(root int) {
	if root == 0 {
		return
	}

	rootRef := &pt.nodes[root]
	if rootRef.isReversed != 0 {
		rootRef.left, rootRef.right = rootRef.right, rootRef.left
		if rootRef.left != 0 {
			rootRef.left = pt.copyNode(rootRef.left)
			pt.nodes[rootRef.left].isReversed ^= 1
		}
		if rootRef.right != 0 {
			rootRef.right = pt.copyNode(rootRef.right)
			pt.nodes[rootRef.right].isReversed ^= 1
		}
		rootRef.isReversed = 0
	}
}

// !OK
func (pt *PersistentFHQTreap) splitByRank(root, k int, left, right *int) {
	if root == 0 {
		*left, *right = 0, 0
		return
	}
	fmt.Println("splitByRank", root, k, left, right)
	pt.pushDown(root)
	if k <= pt.nodes[pt.nodes[root].left].size {
		*right = pt.copyNode(root)
		pt.splitByRank(pt.nodes[*right].left, k, left, &pt.nodes[*right].left)
		pt.pushUp(*right)
	} else {
		*left = pt.copyNode(root)
		pt.splitByRank(pt.nodes[*left].right, k-pt.nodes[pt.nodes[root].left].size-1, &pt.nodes[*left].right, right)
		pt.pushUp(*left)
	}
}

// !OK
// 返回新版本的根节点编号
// 注意merge中无需再拷贝 因为split总是在merge之前调用
func (t *PersistentFHQTreap) merge(x, y int) int {
	if x == 0 || y == 0 {
		return x | y
	}

	// https://nyaannyaan.github.io/library/rbst/rbst-base.hpp
	// if (int((rng() * (l->cnt + r->cnt)) >> 32) < l->cnt) {
	if int(nextRand()*(uint64(t.nodes[x].size)+uint64(t.nodes[y].size))>>32) < t.nodes[x].size {
		t.pushDown(x)
		t.nodes[x].right = t.merge(t.nodes[x].right, y)
		t.pushUp(x)
		return x
	} else {
		t.pushDown(y)
		t.nodes[y].left = t.merge(x, t.nodes[y].left)
		t.pushUp(y)
		return y
	}
}

// !OK
// 插入元素 返回新版本的根节点编号
func (pt *PersistentFHQTreap) Insert(rootVersion int, index int, value int) int {
	var left, right int
	pt.splitByRank(rootVersion, index, &left, &right)
	newRoot := pt.merge(pt.merge(left, pt.newNode(value)), right)
	pt.root = newRoot
	return newRoot
}

func (pt *PersistentFHQTreap) Pop(rootVersion int, index int) int {
	// index++
	var a, b, c int
	pt.splitByRank(rootVersion, index, &b, &c)
	pt.splitByRank(b, index-1, &a, &b)
	newRoot := pt.merge(a, c)
	pt.root = newRoot
	return newRoot
}

// Reverse [start, stop) in place.
func (pt *PersistentFHQTreap) Reverse(rootVersion int, left, right int) int {
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
	// left++
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

func (pt *PersistentFHQTreap) Size() int {
	return pt.nodes[pt.root].size
}

// Return all elements in index order.
func (pt *PersistentFHQTreap) InOrder(rootVersion int) []int {
	res := make([]int, 0, pt.nodes[rootVersion].size)
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

// https://github.com/EndlessCheng/codeforces-go/blob/f9d97465d8b351af7536b5b6dac30b220ba1b913/copypasta/treap.go#L31
func (pt *PersistentFHQTreap) fastRand() uint {
	pt.seed ^= pt.seed << 13
	pt.seed ^= pt.seed >> 17
	pt.seed ^= pt.seed << 5
	return pt.seed
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
