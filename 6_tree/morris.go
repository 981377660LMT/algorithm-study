package main

import "fmt"

// BinaryTreeNode 定义了一个泛型接口。
// 接口的类型参数 T 表示节点本身的具体类型，
// 必须包含获取左右子节点和设置右子节点的方法，
// 用于 Morris 遍历中的线程指针操作。
// 实现此接口的自定义节点类型可以是任意结构体。
//
// 例如：
//
//	type TreeNode struct {
//	    Val         int
//	    LeftNode    *TreeNode
//	    RightNode   *TreeNode
//	}
//	func (n *TreeNode) Left() *TreeNode      { return n.LeftNode }
//	func (n *TreeNode) Right() *TreeNode     { return n.RightNode }
//	func (n *TreeNode) SetRight(r *TreeNode) { n.RightNode = r }
type BinaryTreeNode[T any] interface {
	Left() T
	Right() T
	SetRight(T)
}

// MorrisInorder 在二叉树上执行中序遍历，
// 无需递归或显式栈，并在遍历时临时创建并移除“线程”指针以恢复树结构。
// T 必须同时实现 BinaryTreeNode[T] 接口和 comparable，以支持 nil 检测。
// visit 回调会被依次调用，参数类型即为遍历到的节点类型 T。
func MorrisInorder[T interface {
	BinaryTreeNode[T]
	comparable
}](root T, visit func(node T)) {
	cur := root
	var nilValue T
	for cur != nilValue {
		left := cur.Left()
		if left == nilValue {
			visit(cur)
			cur = cur.Right()
		} else {
			pred := left
			for {
				right := pred.Right()
				if right == nilValue || right == cur {
					break
				}
				pred = right
			}
			if pred.Right() == nilValue {
				pred.SetRight(cur)
				cur = cur.Left()
			} else {
				pred.SetRight(nilValue)
				visit(cur)
				cur = cur.Right()
			}
		}
	}
}

// TreeNode 定义了二叉搜索树节点
// 实现左右访问与 SetRight 方法
type TreeNode struct {
	Val       int
	LeftNode  *TreeNode
	RightNode *TreeNode
}

func (n *TreeNode) Left() *TreeNode      { return n.LeftNode }
func (n *TreeNode) Right() *TreeNode     { return n.RightNode }
func (n *TreeNode) SetRight(r *TreeNode) { n.RightNode = r }

// recoverTree 在 BST 中恢复正被交换的两个节点值
func recoverTree(root *TreeNode) {
	var first, second, prev *TreeNode
	// 中序遍历，寻找逆序对
	MorrisInorder[*TreeNode](root, func(node *TreeNode) {
		if prev != nil && prev.Val > node.Val {
			if first == nil {
				first = prev
			}
			second = node
		}
		prev = node
	})
	// 交换找到的两个节点的值
	if first != nil && second != nil {
		first.Val, second.Val = second.Val, first.Val
	}
}

// 示例测试
func main() {
	// 构造示例：交换了两个节点
	//   1
	//  / \
	// 3   4
	//  \  /
	//   2
	root := &TreeNode{Val: 1}
	root.LeftNode = &TreeNode{Val: 3}
	root.RightNode = &TreeNode{Val: 4}
	root.LeftNode.RightNode = &TreeNode{Val: 2}

	fmt.Println("恢复前中序：")
	MorrisInorder[*TreeNode](root, func(n *TreeNode) { fmt.Printf("%d ", n.Val) })
	fmt.Println()

	recoverTree(root)

	fmt.Println("恢复后中序：")
	MorrisInorder[*TreeNode](root, func(n *TreeNode) { fmt.Printf("%d ", n.Val) })
	fmt.Println()
}
