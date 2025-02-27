// 173. 二叉搜索树迭代器
// https://leetcode.cn/problems/binary-search-tree-iterator/description/?envType=problem-list-v2&envId=design
//
// !Next() 和 hasNext() 的均摊时间复杂度为 O(1)，空间复杂度为 O(h)，其中 h 是树的高度。

package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type BSTIterator struct {
	stack []*TreeNode
	cur   *TreeNode
}

func Constructor(root *TreeNode) BSTIterator {
	return BSTIterator{cur: root}
}

func (this *BSTIterator) Next() int {
	for node := this.cur; node != nil; node = node.Left {
		this.stack = append(this.stack, node)
	}
	this.cur = this.stack[len(this.stack)-1]
	this.stack = this.stack[:len(this.stack)-1]
	res := this.cur.Val
	this.cur = this.cur.Right
	return res
}

func (this *BSTIterator) HasNext() bool {
	return this.cur != nil || len(this.stack) > 0
}

/**
* Your BSTIterator object will be instantiated and called as such:
* obj := Constructor(root);
* param_1 := obj.Next();
* param_2 := obj.HasNext();
 */
