package main

import (
	"bufio"
	"fmt"
	"os"
)

// 空间复杂度`O(n)`的动态树上倍增(LcaOnline).
type CompressedBinaryLiftDynamic struct {
	Depth  []int32
	Parent []int32
	jump   []int32
}

// 不预先给出整棵树,而是动态添加叶子节点,维护树节点的LCA和k级祖先.
func NewCompressedBinaryLiftDynamic(n int32) *CompressedBinaryLiftDynamic {
	res := &CompressedBinaryLiftDynamic{
		Depth:  make([]int32, n),
		Parent: make([]int32, n),
		jump:   make([]int32, n),
	}
	return res
}

func NewCompressedBinaryLiftDynamicWithRoot(n int32, root int32) *CompressedBinaryLiftDynamic {
	res := NewCompressedBinaryLiftDynamic(n)
	res.AddLeaf(root, -1)
	return res
}

// 在树中添加一条从parent到leaf的边，要求parent已经存在于树中或为-1.
// parent为-1时，leaf为根节点.
func (bl *CompressedBinaryLiftDynamic) AddLeaf(leaf int32, parent int32) {
	if parent == -1 {
		bl.Parent[leaf] = -1
		bl.jump[leaf] = leaf
		return
	}
	bl.Depth[leaf] = bl.Depth[parent] + 1
	bl.Parent[leaf] = parent
	if tmp := bl.jump[parent]; bl.Depth[parent]-bl.Depth[tmp] == bl.Depth[tmp]-bl.Depth[bl.jump[tmp]] {
		bl.jump[leaf] = bl.jump[tmp]
	} else {
		bl.jump[leaf] = parent
	}
}

func (bl *CompressedBinaryLiftDynamic) FirstTrue(start int32, predicate func(end int32) bool) int32 {
	for !predicate(start) {
		if predicate(bl.jump[start]) {
			start = bl.Parent[start]
		} else {
			if start == bl.jump[start] {
				return -1
			}
			start = bl.jump[start]
		}
	}
	return start
}

func (bl *CompressedBinaryLiftDynamic) LastTrue(start int32, predicate func(end int32) bool) int32 {
	if !predicate(start) {
		return -1
	}
	for {
		if predicate(bl.jump[start]) {
			if start == bl.jump[start] {
				return start
			}
			start = bl.jump[start]
		} else if predicate(bl.Parent[start]) {
			start = bl.Parent[start]
		} else {
			return start
		}
	}
}

func (bl *CompressedBinaryLiftDynamic) UpToDepth(root int32, toDepth int32) int32 {
	if !(0 <= toDepth && toDepth <= bl.Depth[root]) {
		return -1
	}
	for bl.Depth[root] > toDepth {
		if bl.Depth[bl.jump[root]] < toDepth {
			root = bl.Parent[root]
		} else {
			root = bl.jump[root]
		}
	}
	return root
}

func (bl *CompressedBinaryLiftDynamic) KthAncestor(node, k int32) int32 {
	targetDepth := bl.Depth[node] - k
	return bl.UpToDepth(node, targetDepth)
}

func (bl *CompressedBinaryLiftDynamic) Lca(a, b int32) int32 {
	if bl.Depth[a] > bl.Depth[b] {
		a = bl.KthAncestor(a, bl.Depth[a]-bl.Depth[b])
	} else if bl.Depth[a] < bl.Depth[b] {
		b = bl.KthAncestor(b, bl.Depth[b]-bl.Depth[a])
	}
	for a != b {
		if bl.jump[a] == bl.jump[b] {
			a = bl.Parent[a]
			b = bl.Parent[b]
		} else {
			a = bl.jump[a]
			b = bl.jump[b]
		}
	}
	return a
}

func (bl *CompressedBinaryLiftDynamic) Dist(a, b int32) int32 {
	return bl.Depth[a] + bl.Depth[b] - 2*bl.Depth[bl.Lca(a, b)]
}

func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int32
	fmt.Fscan(in, &q)
	operations := make([][2]int32, q)
	for i := int32(0); i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		switch op {
		case "+":
			var x int32
			fmt.Fscan(in, &x)
			operations[i] = [2]int32{1, x}
		case "-":
			var k int32
			fmt.Fscan(in, &k)
			operations[i] = [2]int32{2, k}
		case "!":
			operations[i] = [2]int32{3, 0}
		case "?":
			operations[i] = [2]int32{4, 0}
		}
	}

	res := Rollbacks(operations)
	for _, v := range res {
		fmt.Fprintln(out, v)
	}
}

// https://www.luogu.com.cn/problem/CF1858E1
// 给定一个初始时为空的数组nums, 需要实现下面四种类型的操作：
// [1, x]: 将x添加到nums尾部
// [2, k]: 将尾部的k个数删除.保证存在k个数.
// [3, 0]: 撤销上一次操作1或2操作
// [4, 0]: 查询当前nums中有多少个不同的数
//
// 1<=q<=1e6,询问次数不超过1e5
//
// !1. 因为要支持撤销，所以需要`保存版本`或者action, 例如immer.js两种方式都会提供.
// 这里保存版本(PersistentStack)比较合适。不同版本构成了一棵树.
// !2. 离线方法是建立一棵版本树，+操作加边，-操作通过倍增上跳到对应节点，
// 操作回退到上个节点，?操作记录当前节点需要记录答案，最后dfs整棵树求解.
// !3. 这个倍增也非常巧妙，不预处理而是是动态的
// !4. 树节点个数最多为n个，因为只有1操作会增加节点.
func Rollbacks(operations [][2]int32) []int32 {
	q := int32(len(operations))
	tree := make([][]int32, q+1) // 0是初始空状态的虚拟根节点
	values := make([]int32, q+1)
	queryGroup := make([][]int32, q+1)
	queryIndex := int32(0)
	history := []int32{0} // 保存节点位置的stack
	nodeId := int32(0)

	lca := NewCompressedBinaryLiftDynamicWithRoot(q+1, 0)
	for _, op := range operations {
		kind, x := op[0], op[1]
		if kind == 1 {
			nodeId++
			values[nodeId] = x
			curNode := history[len(history)-1]
			tree[curNode] = append(tree[curNode], nodeId)
			lca.AddLeaf(nodeId, curNode)
			history = append(history, nodeId)
		} else if kind == 2 {
			cur := history[len(history)-1]
			history = append(history, lca.KthAncestor(cur, x))
		} else if kind == 3 {
			if len(history) <= 1 {
				continue
			}
			history = history[:len(history)-1]
		} else {
			cur := history[len(history)-1]
			queryGroup[cur] = append(queryGroup[cur], queryIndex)
			queryIndex++
		}
	}

	res := make([]int32, queryIndex)
	counter := make(map[int32]int32)
	unique := int32(0)
	var dfs func(int32)
	dfs = func(node int32) {
		curValue := values[node]
		if node != 0 {
			counter[curValue]++
			if counter[curValue] == 1 {
				unique++
			}
		}
		for _, q := range queryGroup[node] {
			res[q] = unique
		}
		for _, next := range tree[node] {
			dfs(next)
		}
		if node != 0 {
			counter[curValue]--
			if counter[curValue] == 0 {
				unique--
			}
		}
	}
	dfs(0)
	return res
}
