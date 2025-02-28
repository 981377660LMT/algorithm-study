// 1586. 二叉搜索树迭代器 II
// https://leetcode.cn/problems/binary-search-tree-iterator-ii/solutions/2387842/er-cha-sou-suo-shu-die-dai-qi-ii-by-leet-3o49/

package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type BSTIterator struct {
	pre     *TreeNode
	stack   []*TreeNode
	arr     []int
	pointer int
}

func Constructor(root *TreeNode) BSTIterator {
	return BSTIterator{
		pre:     root,
		pointer: -1,
	}
}

func (this *BSTIterator) HasNext() bool {
	return len(this.stack) > 0 || this.pre != nil || this.pointer < len(this.arr)-1
}

func (this *BSTIterator) Next() int {
	this.pointer++
	if this.pointer == len(this.arr) {
		for this.pre != nil {
			this.stack = append(this.stack, this.pre)
			this.pre = this.pre.Left
		}
		cur := this.stack[len(this.stack)-1]
		this.stack = this.stack[:len(this.stack)-1]
		this.pre = cur.Right
		this.arr = append(this.arr, cur.Val)
	}
	return this.arr[this.pointer]
}

func (this *BSTIterator) HasPrev() bool {
	return this.pointer > 0
}

func (this *BSTIterator) Prev() int {
	this.pointer--
	return this.arr[max(0, this.pointer)]
}

/**
* Your BSTIterator object will be instantiated and called as such:
* obj := Constructor(root);
* param_1 := obj.HasNext();
* param_2 := obj.Next();
* param_3 := obj.HasPrev();
* param_4 := obj.Prev();
 */
