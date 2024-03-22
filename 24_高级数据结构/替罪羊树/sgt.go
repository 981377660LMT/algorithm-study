// ScapeGoatTree (SGT) 替罪羊树
// https://zhuanlan.zhihu.com/p/180545164
// 为了防止二叉搜索树左右不平衡，我们引入平衡树，而其中思路最简单的是替罪羊树（Scapegoat tree）。
// 替罪羊树 是一种依靠重构操作维持平衡的重量平衡树。
// 替罪羊树会在插入、删除操作时，检测途经的节点，若发现失衡，则将以该节点为根的子树重构。
// 当发现某个子树很不平衡时，暴力重构该子树使之平衡(实现时，只重构子树最大的那一个结点(scapegoat))。
// 若左子树或右子树占当前树的比例大于alpha(一般是0.7-0.8) ，则进行重构(很多教程还会判断已删除节点的个数占总大小的比例(0.5)来决定重不重构)
// 重构分为两步操作：先进行一遍中序遍历，把该子树“拉平”，把其中所有数存入一个数组里（BST的性质决定这个数组一定是有序的）；
// 然后，再用这些数据重新建一个平衡的二叉树，放回原位置。
// 一个节点导致树的不平衡，就要导致整棵子树被拍扁，估计这也是“替罪羊”这个名字的由来吧.
//
// https://people.ksp.sk/~kuko/gnarley-trees/Scapegoat.html
// lazy insert: let the tree grow and from time to time, when a subtree gets too imbalanced,
//              rebuild the whole subtree from scratch into a perfectly balanced tree
// lazy delete: just mark the node for deletion; when n/2 nodes are marked,
//              rebuild the whole tree and throw away all deleted nodes
//
// !优点：不受到旋转机制弊端的影响.
// 缺点：无法持久化、无法维护区间信息.
//
// 可以优化的地方：
//  1.使用空节点代替nil，避免特判(优化分支).
//  2.使用内存池，直接分配好内存，避免动态分配内存占用时间.
//  3.使用回收栈，复用结点，避免频繁分配内存.
//  4.使用非递归的方式求排名和第K大.
//  5.使用int32代替结点指针，减少内存占用.
//  6.Init、Clear 生命周期.
//
// https://zhuanlan.zhihu.com/p/21263304
// https://taodaling.github.io/blog/2019/04/19/%E6%9B%BF%E7%BD%AA%E7%BE%8A%E6%A0%91/
// https://riteme.site/blog/2016-4-6/scapegoat.html
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/scapegoat_tree.go
// https://www.bilibili.com/video/BV1sP4y1i7WB/
// https://juejin.cn/post/6844904128150241294 使用替罪羊树实现KD-Tree的增删改查

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
	"sort"
	"strings"
)

func init() {
	debug.SetGCPercent(-1)
}

func main() {
	P6136()
}

// P6136 【模板】普通平衡树（数据加强版）
// https://www.luogu.com.cn/problem/P6136
func P6136() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int32
	fmt.Fscan(in, &n, &q)
	nums := make([]int32, n)
	for i := int32(0); i < n; i++ {
		fmt.Fscan(in, &nums[i])
	}
	sgt := NewScapegoatTree(func(key1, key2 SgtKey) bool { return key1 < key2 }, nums...)

	lastRes := int32(0)
	preRes := int32(0)
	for i := int32(0); i < q; i++ {
		var kind, x int32
		fmt.Fscan(in, &kind, &x)
		x ^= preRes
		switch kind {
		case 1:
			sgt.Add(x)
		case 2:
			sgt.Discard(x)
		case 3:
			preRes = sgt.BisectLeft(x) + 1
			lastRes ^= preRes
		case 4:
			preRes = sgt.At(x - 1)
			lastRes ^= preRes
		case 5:
			preRes, _ = sgt.Prev(x)
			lastRes ^= preRes
		case 6:
			preRes, _ = sgt.Next(x)
			lastRes ^= preRes
		default:
			panic("invalid kind")
		}
	}

	fmt.Fprintln(out, lastRes)
}

func demo() {
	sgt := NewScapegoatTree(func(key1, key2 SgtKey) bool { return key1 < key2 })

	sgt.Add(1)
	sgt.Add(1)
	sgt.Add(2)
	sgt.Add(2)
	sgt.Add(2)
	sgt.Add(3)
	sgt.Add(-3)
	sgt.Add(-3)

	fmt.Println("---")
	fmt.Println(sgt.BisectLeft(1))  // 5
	fmt.Println(sgt.BisectRight(1)) // 5
	fmt.Println(sgt.At(-1))         // 5
	fmt.Println(sgt.At(3))          // 5
	fmt.Println(sgt.At(4))          // 5
	fmt.Println(sgt.Size(), 99)
	fmt.Println(sgt.Pop(0))  // -3
	fmt.Println(sgt.Next(2)) // -3

	fmt.Println()
	sgt.Enumerate(func(key SgtKey) { fmt.Print(key, " ") })
	fmt.Println()
	fmt.Println(sgt)
}

// alpha 的值越小，那么替罪羊树就越容易重构，那么树也就越平衡，查询的效率也就越高，自然修改（加点和删点）的效率也就低了；
// 反之，alpha 的值越大，那么替罪羊树就越不容易重构，那么树也就越不平衡，查询的效率也就越低，自然修改（加点和删点）的效率也就高了。
// 所以，查询多，alpha 就应该小一些；修改多，alpha 就应该大一些。
// alpha = 4/5
const ALPHA_NUM int32 = 4
const ALPHA_DENO int32 = 5

type SgtKey = int32

type SgtNode struct {
	Key         SgtKey
	left, right *SgtNode
	existCount  int32 // 子树有效结点的数量
	allCount    int32 // 子树所有结点的数量
	exist       bool
}

var EMPTY_NODE = &SgtNode{} // 空结点不使用nil，避免特判

type ScapegoatTree struct {
	less func(key1, key2 SgtKey) bool
	root *SgtNode
	// recycles  []*SgtNode // 用于回收被删除的结点(注意: rebuild时才进入回收栈)
	collector []*SgtNode // 用于dfs收集结点
}

func NewScapegoatTree(less func(key1, key2 SgtKey) bool, elements ...SgtKey) *ScapegoatTree {
	res := &ScapegoatTree{less: less, root: EMPTY_NODE}
	if len(elements) > 0 {
		elements = append(elements[:0:0], elements...)
		sort.Slice(elements, func(i, j int) bool { return less(elements[i], elements[j]) })
		nodes := make([]*SgtNode, len(elements))
		for i, key := range elements {
			nodes[i] = res.Alloc(key)
		}
		res.root = res._build(nodes, 0, int32(len(nodes)))
	}
	return res
}

func (t *ScapegoatTree) Alloc(key SgtKey) *SgtNode {
	return &SgtNode{Key: key, left: EMPTY_NODE, right: EMPTY_NODE, existCount: 1, allCount: 1, exist: true}
	// if len(t.recycles) > 0 {
	// 	res := t.recycles[len(t.recycles)-1]
	// 	t.recycles = t.recycles[:len(t.recycles)-1]
	// 	res.Key = key
	// 	res.left, res.right = EMPTY_NODE, EMPTY_NODE
	// 	res.existCount, res.allCount = 1, 1
	// 	res.exist = true
	// 	return res
	// } else {
	// 	return &SgtNode{Key: key, left: EMPTY_NODE, right: EMPTY_NODE, existCount: 1, allCount: 1, exist: true}
	// }
}

func (t *ScapegoatTree) Add(key SgtKey) {
	scapegoat := t._insert(&t.root, key)
	if *scapegoat != EMPTY_NODE {
		t._rebuild(scapegoat)
	}
}

// 删除前需要保证 key 存在.
func (t *ScapegoatTree) Remove(key SgtKey) {
	t._remove(t.root, t.BisectLeft(key)+1)
	if t.root.existCount*2 < t.root.allCount {
		t._rebuild(&t.root)
	}
}

func (t *ScapegoatTree) Discard(key SgtKey) bool {
	ok := t._discard(t.root, t.BisectLeft(key)+1)
	if !ok {
		return false
	}
	if t.root.existCount*2 < t.root.allCount {
		t._rebuild(&t.root)
	}
	return true
}

func (t *ScapegoatTree) Pop(index int32) SgtKey {
	size := t.root.existCount
	if index < 0 {
		index += size
	}
	if index < 0 || index >= size {
		panic(fmt.Sprintf("Pop(%d) not found", index))
	}
	index++
	res := t._remove(t.root, index)
	if t.root.existCount*2 < t.root.allCount {
		t._rebuild(&t.root)
	}
	return res
}

func (t *ScapegoatTree) At(index int32) SgtKey {
	size := t.root.existCount
	if index < 0 {
		index += size
	}
	if index < 0 || index >= size {
		panic(fmt.Sprintf("at(%d) not found", index))
	}
	index++
	cur := t.root
	for cur != EMPTY_NODE {
		if cur.left.existCount+1 == index && cur.exist {
			return cur.Key
		}
		if cur.left.existCount >= index {
			cur = cur.left
		} else {
			index -= cur.left.existCount
			if cur.exist {
				index--
			}
			cur = cur.right
		}
	}
	panic(fmt.Sprintf("at(%d) not found", index))
}

// 严格小于key的数的个数.
func (t *ScapegoatTree) BisectLeft(key SgtKey) int32 {
	cur := t.root
	res := int32(0)
	for cur != EMPTY_NODE {
		if !t.less(cur.Key, key) {
			cur = cur.left
		} else {
			res += cur.left.existCount
			if cur.exist {
				res++
			}
			cur = cur.right
		}
	}
	return res
}

// 小于等于key的数的个数.
func (t *ScapegoatTree) BisectRight(key SgtKey) int32 {
	cur := t.root
	res := int32(0)
	for cur != EMPTY_NODE {
		if t.less(key, cur.Key) {
			cur = cur.left
		} else {
			res += cur.left.existCount
			if cur.exist {
				res++
			}
			cur = cur.right
		}
	}
	return res
}

func (t *ScapegoatTree) Prev(key SgtKey) (res SgtKey, ok bool) {
	less := t.BisectLeft(key)
	if less == 0 {
		return
	}
	res, ok = t.At(less-1), true
	return
}

func (t *ScapegoatTree) Next(key SgtKey) (res SgtKey, ok bool) {
	ngt := t.BisectRight(key)
	if ngt == t.root.existCount {
		return
	}
	res, ok = t.At(ngt), true
	return
}

func (t *ScapegoatTree) Size() int32 {
	return t.root.existCount
}

func (t *ScapegoatTree) Enumerate(f func(key SgtKey)) {
	var dfs func(*SgtNode)
	dfs = func(node *SgtNode) {
		if node == EMPTY_NODE {
			return
		}
		dfs(node.left)
		if node.exist {
			f(node.Key)
		}
		dfs(node.right)
	}
	dfs(t.root)
}

func (t *ScapegoatTree) _collect(cur *SgtNode, nodes *[]*SgtNode) {
	if cur == EMPTY_NODE {
		return
	}
	t._collect(cur.left, nodes)
	if cur.exist {
		*nodes = append(*nodes, cur)
	}
	// if cur.exist {
	// 	*nodes = append(*nodes, cur)
	// } else {
	// 	t.recycles = append(t.recycles, cur)
	// }
	t._collect(cur.right, nodes)
}

func (t *ScapegoatTree) _build(nodes []*SgtNode, left, right int32) *SgtNode {
	if left >= right {
		return EMPTY_NODE
	}
	mid := (left + right) >> 1
	res := nodes[mid]
	res.left = t._build(nodes, left, mid)
	res.right = t._build(nodes, mid+1, right)
	t._pushUp(res)
	return res
}

func (t *ScapegoatTree) _rebuild(nodePtr **SgtNode) {
	t.collector = t.collector[:0]
	t._collect(*nodePtr, &t.collector)
	*nodePtr = t._build(t.collector, 0, int32(len(t.collector)))
}

func (t *ScapegoatTree) _pushUp(node *SgtNode) {
	node.existCount = node.left.existCount + node.right.existCount
	if node.exist {
		node.existCount++
	}
	node.allCount = node.left.allCount + node.right.allCount + 1
}

// 判断是否需要重构.
func (t *ScapegoatTree) _isUnbalanced(node *SgtNode) bool {
	// +5，避免不必要的重构
	threshold := node.allCount*ALPHA_NUM + 5*ALPHA_DENO
	return (node.left.allCount*ALPHA_DENO > threshold) || (node.right.allCount*ALPHA_DENO > threshold)
}

// 返回需要重构的结点.
// 一次修改可能会变更整个搜索路径上的所有子树大小，如果多个子树需要重构，选择最大的那颗。
func (t *ScapegoatTree) _insert(nodePtr **SgtNode, key SgtKey) **SgtNode {
	if *nodePtr == EMPTY_NODE {
		*nodePtr = t.Alloc(key)
		return &EMPTY_NODE
	} else {
		node := *nodePtr
		node.existCount++
		node.allCount++
		if t.less(key, node.Key) {
			res := t._insert(&node.left, key)
			if t._isUnbalanced(node.left) {
				return nodePtr
			} else {
				return res
			}
		} else {
			res := t._insert(&node.right, key)
			if t._isUnbalanced(node.right) {
				return nodePtr
			} else {
				return res
			}
		}
	}
}

func (t *ScapegoatTree) _remove(node *SgtNode, k int32) SgtKey {
	node.existCount--
	offset := node.left.existCount
	if node.exist {
		offset++
	}
	if node.exist && k == offset {
		node.exist = false
		return node.Key
	} else {
		if k <= offset {
			return t._remove(node.left, k)
		} else {
			return t._remove(node.right, k-offset)
		}
	}
}

func (t *ScapegoatTree) _discard(node *SgtNode, k int32) bool {
	if node == EMPTY_NODE {
		return false
	}
	node.existCount--
	offset := node.left.existCount
	if node.exist {
		offset++
	}
	if node.exist && k == offset {
		node.exist = false
		return true
	} else {
		if k <= offset {
			return t._discard(node.left, k)
		} else {
			return t._discard(node.right, k-offset)
		}
	}
}

func (o *SgtNode) String() string {
	return fmt.Sprintf("%v", o.Key)
}

/*
	逆时针旋转 90° 打印这棵树：根节点在最左侧，右子树在上侧，左子树在下侧

效果如下（只打印 key）

Root
│           ┌── 95
│       ┌── 94
│   ┌── 90
│   │   │           ┌── 89
│   │   │       ┌── 88
│   │   │       │   └── 87
│   │   │       │       └── 81
│   │   │   ┌── 74
│   │   └── 66
└── 62

	│           ┌── 59
	│       ┌── 58
	│       │   └── 56
	│       │       └── 47
	│   ┌── 45
	└── 40
	    │       ┌── 37
	    │   ┌── 28
	    └── 25
	        │           ┌── 18
	        │       ┌── 15
	        │   ┌── 11
	        └── 6
	            └── 0
*/
func (o *SgtNode) draw(treeSB, prefixSB *strings.Builder, isTail bool) {
	prefix := prefixSB.String()

	if o.right != EMPTY_NODE {
		newPrefixSB := &strings.Builder{}
		newPrefixSB.WriteString(prefix)
		if isTail {
			newPrefixSB.WriteString("│   ")
		} else {
			newPrefixSB.WriteString("    ")
		}
		o.right.draw(treeSB, newPrefixSB, false)
	}

	treeSB.WriteString(prefix)
	if isTail {
		treeSB.WriteString("└── ")
	} else {
		treeSB.WriteString("┌── ")
	}
	if o.exist {
		treeSB.WriteString(o.String())
	} else {
		treeSB.WriteString("x")
	}
	treeSB.WriteByte('\n')

	if o.left != EMPTY_NODE {
		newPrefixSB := &strings.Builder{}
		newPrefixSB.WriteString(prefix)
		if isTail {
			newPrefixSB.WriteString("    ")
		} else {
			newPrefixSB.WriteString("│   ")
		}
		o.left.draw(treeSB, newPrefixSB, true)
	}
}

func (t *ScapegoatTree) String() string {
	if t.root == EMPTY_NODE {
		return "Empty\n"
	}
	treeSB := &strings.Builder{}
	treeSB.WriteString("Root\n")
	t.root.draw(treeSB, &strings.Builder{}, true)
	return treeSB.String()
}
