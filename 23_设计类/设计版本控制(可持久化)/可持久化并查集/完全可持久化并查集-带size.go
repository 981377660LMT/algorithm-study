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

func yosupo() {
	// https://judge.yosupo.jp/submission/130167
	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, q int
	fmt.Fscan(in, &n, &q)

	uf := NewPersistentUnionfind()
	versions := make([]*AryNode, 0, q+10)
	versions = append(versions, uf.Alloc())

	for i := 0; i < q; i++ {
		var op, version, u, v int
		fmt.Fscan(in, &op, &version, &u, &v)
		version++
		root := versions[version]
		if op == 0 {
			newRoot, _ := uf.Union(root, u, v)
			root = newRoot
		} else {
			if uf.IsConnected(root, u, v) {
				fmt.Fprintln(out, 1)
			} else {
				fmt.Fprintln(out, 0)
			}
		}
		versions = append(versions, root)
	}
}

func UnionSets() {
	// # https://atcoder.jp/contests/code-thanks-festival-2017-open/tasks/code_thanks_festival_2017_h
	// # 给定n个集合,初始时第i个集合只有一个元素i (i=1,2,...,n)
	// # 之后进行m次合并操作,每次合并ai和bi所在的集合
	// # 如果ai和bi在同一个集合,则无事发生
	// # 给定q个询问,问ai和bi是在第几次操作后第一次连通的,如果不连通则输出-1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var n, m int
	fmt.Fscan(in, &n, &m)
	uf := NewPersistentUnionfind()
	version := make([]*AryNode, 0, m+1)
	version = append(version, uf.Alloc())
	for i := 0; i < m; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		newRoot, _ := uf.Union(version[i], a, b)
		version = append(version, newRoot)
	}

	var q int
	fmt.Fscan(in, &q)
	for i := 0; i < q; i++ {
		var a, b int
		fmt.Fscan(in, &a, &b)
		a, b = a-1, b-1
		if !uf.IsConnected(version[m], a, b) {
			fmt.Fprintln(out, -1)
			continue
		}

		left, right := 0, m
		for left <= right {
			mid := (left + right) / 2
			if uf.IsConnected(version[mid], a, b) {
				right = mid - 1
			} else {
				left = mid + 1
			}
		}
		fmt.Fprintln(out, left)
	}
}

func main() {
	// yosupo()
	UnionSets()
}

type E = int
type AryNode struct {
	data     E
	children [16]*AryNode // !
}

// 完全可持久化并查集,不使用路径压缩.
type PersistentUnionfind struct {
	parents *PersistentArray
}

func NewPersistentUnionfind() *PersistentUnionfind {
	return &PersistentUnionfind{parents: NewPersistentArray(-1)}
}

func (p *PersistentUnionfind) Alloc() *AryNode {
	return p.parents.Alloc()
}

func (p *PersistentUnionfind) Union(root *AryNode, x, y int) (newRoot *AryNode, ok bool) {
	x, y = p.Find(root, x), p.Find(root, y)
	if x == y {
		return root, false
	}
	p1, p2 := p.parents.Get(root, x), p.parents.Get(root, y)
	if -p1 < -p2 {
		x, y = y, x
	}
	newSize := p1 + p2
	root = p.parents.Set(root, x, newSize)
	root = p.parents.Set(root, y, x)
	newRoot, ok = root, true
	return
}

func (p *PersistentUnionfind) Find(root *AryNode, x int) int {
	for {
		p := p.parents.Get(root, x)
		if p < 0 {
			break
		}
		x = p
	}
	return x
}

func (p *PersistentUnionfind) IsConnected(root *AryNode, x, y int) bool {
	return p.Find(root, x) == p.Find(root, y)
}

func (p *PersistentUnionfind) GetSize(root *AryNode, x int) int {
	return -p.parents.Get(root, p.Find(root, x))
}

type PersistentArray struct {
	null E // 当index越界/不存在时返回的值
}

func NewPersistentArray(null E) *PersistentArray {
	return &PersistentArray{null: null}
}

func (o *PersistentArray) Alloc() *AryNode {
	return &AryNode{data: o.null}
}

func (o *PersistentArray) Build(nums []E) *AryNode {
	root := o.Alloc()
	for i := 0; i < len(nums); i++ {
		root = o.setWithoutCopy(root, i, nums[i])
	}
	return root
}

func (o *PersistentArray) Get(root *AryNode, index int) E {
	if root == nil {
		return o.null
	}
	if index == 0 {
		return root.data
	}
	return o.Get(root.children[index&15], (index-1)>>4)
}

func (o *PersistentArray) Set(root *AryNode, index int, data E) *AryNode {
	newNode := o.Alloc()
	if root != nil { // copyNode
		newNode.data = root.data
		for i := 0; i < 16; i++ {
			newNode.children[i] = root.children[i]
		}
	}
	if index == 0 {
		newNode.data = data
		return newNode
	}
	newNode.children[index&15] = o.Set(newNode.children[index&15], (index-1)>>4, data)
	return newNode
}

func (o *PersistentArray) setWithoutCopy(root *AryNode, index int, data E) *AryNode {
	if root == nil {
		root = o.Alloc()
	}
	if index == 0 {
		root.data = data
		return root
	}
	root.children[index&15] = o.setWithoutCopy(root.children[index&15], (index-1)>>4, data)
	return root
}
