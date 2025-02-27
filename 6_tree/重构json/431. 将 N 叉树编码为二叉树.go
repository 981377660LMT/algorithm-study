// 431. 将 N 叉树编码为二叉树
// https://leetcode.cn/problems/encode-n-ary-tree-to-binary-tree/?envType=problem-list-v2&envId=design
//
// 步骤：
// 1. 将所有兄弟节点链接在一起，形成一个单向链表。
// 2. 将得到的兄弟节点列表的head与其parent节点相链接
//
// !二叉树节点的左孩子就是节点孩子，二叉树节点的右孩子，就是节点的兄弟。
// !反过来N叉树的第一孩子节点就是其左孩子，其他孩子节点，就是刚刚那个左孩子节点的全部右子树

package main

type Node struct {
	Val      int
	Children []*Node
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

type Codec struct{}

func Constructor() *Codec { return &Codec{} }

func (this *Codec) encode(root *Node) *TreeNode {
	if root == nil {
		return nil
	}
	newRoot := &TreeNode{Val: root.Val}
	// 将 N 叉树节点的第一个子节点编码到二叉树的左侧节点
	if len(root.Children) > 0 {
		newRoot.Left = this.encode(root.Children[0])
	}
	// 对其余兄弟节点进行编码
	sibling := newRoot.Left
	for i := 1; i < len(root.Children); i++ {
		sibling.Right = this.encode(root.Children[i])
		sibling = sibling.Right
	}
	return newRoot
}

func (this *Codec) decode(root *TreeNode) *Node {
	if root == nil {
		return nil
	}
	newRoot := &Node{Val: root.Val}
	// 对所有子节点进行解码
	sibling := root.Left
	for sibling != nil {
		newRoot.Children = append(newRoot.Children, this.decode(sibling))
		sibling = sibling.Right
	}
	return newRoot
}

/**
* Your Codec object will be instantiated and called as such:
* obj := Constructor();
* bst := obj.encode(root);
* ans := obj.decode(bst);
 */
