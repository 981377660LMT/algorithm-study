// !EndlessCheng/codeforces-go/copypasta/leftist_tree.go
// https://nyaannyaan.github.io/library/data-structure/skew-heap.hpp

package main

import (
	"bufio"
	"fmt"
	"os"
)

/*
左偏树(可并堆) leftist tree / leftist heap
代码参考 https://oi-wiki.org/ds/leftist-tree/

模板题
https://www.luogu.com.cn/problem/P3377
https://www.luogu.com.cn/problem/P2713

!非常像并查集
*/

type LeftistTreeNode struct {
	size, value, index  int
	isMin               bool // 是否为小根堆
	left, right, parent *LeftistTreeNode
}

// 一开始有 n 个小根堆，每个堆包含且仅包含一个数。接下来需要支持两种操作：
// !1 x y 合并 x 号堆和 y 号堆 (若第 x 或第 y 个数已经被删除或第 x 和第 y 个数在用一个堆内，则无视此操作)
// !2 x 输出第 x 号堆的堆顶数并删除该数 (若有多个最小数，优先删除先输入的；若第 x 个数已经被删除，则输出 −1 并无视删除操作)
func main() {
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	heapRoots := make([]*LeftistTreeNode, 0, n)
	for i := 0; i < n; i++ {
		var v int
		fmt.Fscan(in, &v)
		heapRoots = append(heapRoots, NewLeftistTreeNode(i, v, true))
	}

	removed := make([]bool, n)
	for i := 0; i < m; i++ {
		var op, x, y int
		fmt.Fscan(in, &op, &x)
		x--
		if op == 1 {
			fmt.Fscan(in, &y)
			y--
			if !removed[x] && !removed[y] {
				heapRoots[x].Push(heapRoots[y]) // !合并
			}
		} else {
			if removed[x] {
				fmt.Fprintln(out, -1)
				continue
			}
			top, _ := heapRoots[x].Pop()
			fmt.Fprintln(out, int(top.value))
			removed[top.index] = true
		}
	}

}

// !EndlessCheng/codeforces-go/copypasta/leftist_tree.go
func NewLeftistTreeNode(index, value int, isMin bool) *LeftistTreeNode {
	o := &LeftistTreeNode{size: 1, value: value, index: index, isMin: isMin}
	o.parent = o
	return o
}

// 注：push 一个节点就相当于 merge 这个节点(所在的组)
func (o *LeftistTreeNode) Push(p *LeftistTreeNode) {
	if o == nil || p == nil {
		return
	}
	o = o.findRoot()
	p = p.findRoot()
	if o == p {
		return
	}
	q := o.merge(p)
	o.parent = q
	p.parent = q
}

// 注：若要复用 top，需要将该节点的 left 和 right 置为 nil，parent 置为自身
func (o *LeftistTreeNode) Pop() (top, newRoot *LeftistTreeNode) {
	o = o.findRoot()
	p := o.left.merge(o.right)
	o.parent = p // 注意这可能会让 parent 指向 nil
	if o.left != nil {
		o.left.parent = p
	}
	if o.right != nil {
		o.right.parent = p
	}
	return o, p
}

// pop and then push back.
func (o *LeftistTreeNode) Peek() int {
	top, newRoot := o.Pop()
	newRoot.Push(top)
	return top.value
}

func (o *LeftistTreeNode) Len() int {
	return o.sizeOrDefault(0)
}

func (o *LeftistTreeNode) findRoot() *LeftistTreeNode {
	if o.parent != o {
		o.parent = o.parent.findRoot() // 路径压缩
	}
	return o.parent
}

func (o *LeftistTreeNode) merge(p *LeftistTreeNode) *LeftistTreeNode {
	if p == nil {
		return o
	}
	if o == nil {
		return p
	}

	if o.isMin {
		if o.value > p.value || o.value == p.value && o.index > p.index { // 大根堆改成 <
			o, p = p, o
		}
	} else {
		if o.value < p.value || o.value == p.value && o.index > p.index {
			o, p = p, o
		}
	}

	o.right = o.right.merge(p)
	if o.left.sizeOrDefault(0) < o.right.sizeOrDefault(0) {
		o.left, o.right = o.right, o.left
	}
	o.size = o.right.sizeOrDefault(0) + 1
	return o
}

func (o *LeftistTreeNode) sizeOrDefault(value int) int {
	if o != nil {
		return o.size
	}
	return value
}
