// !可持久化数组(可以越界访问的实现)
// Persistent Array
// Reference: <https://qiita.com/hotman78/items/9c643feae1de087e6fc5>
//            <https://maspypy.github.io/library/ds/dynamic_array.hpp>
// - 16叉树
// - Fully persistent
// - `Get(root, k)`:  get k-th element
// - `Set(root, k, data)`: make new array whose k-th element is updated to data

package main

import (
	"bufio"
	"fmt"
	"os"
	"runtime/debug"
)

// 单组测试时禁用gc
func init() {
	debug.SetGCPercent(-1)
}

func main() {
	// https://atcoder.jp/contests/abc273/tasks/abc273_e
	// 模拟git
	// 现在给定一个空数组A 和一系列操作,操作有四种

	// 1. ADD x : 在A的末尾添加一个元素x (git commit)
	// 2. DELETE : 删除A的末尾元素 (git revert)
	// 3. SAVE y : 保存当前数组A的状态到y分支 (git checkout -b)
	// 4. LOAD z : 将z分支加载到当前数组A (git merge)
	// 每次操作之后输出当前数组A的最后一个元素,如果数组为空则输出-1

	in := bufio.NewReader(os.Stdin)
	out := bufio.NewWriter(os.Stdout)
	defer out.Flush()

	var q int
	fmt.Fscan(in, &q)

	type pair struct {
		ptr  int
		root *AryNode
	}
	git := make(map[int]pair)

	undefined := -1
	A := NewPersistentArray(undefined)
	root := A.Alloc()
	ptr := 0
	for i := 0; i < q; i++ {
		var op string
		fmt.Fscan(in, &op)
		if op == "ADD" {
			var x int
			fmt.Fscan(in, &x)
			root = A.Set(root, ptr, x)
			ptr++
		} else if op == "DELETE" {
			if ptr > 0 {
				ptr--
			}
		} else if op == "SAVE" {
			var y int
			fmt.Fscan(in, &y)
			git[y] = pair{ptr, root}
		} else if op == "LOAD" {
			var z int
			fmt.Fscan(in, &z)
			ptr, root = git[z].ptr, git[z].root
		}

		last := A.Get(root, ptr-1)
		if last == undefined {
			fmt.Fprint(out, -1, " ")
		} else {
			fmt.Fprint(out, last, " ")
		}
	}
}

type E = int

type AryNode struct {
	data     E
	children [16]*AryNode // !
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
