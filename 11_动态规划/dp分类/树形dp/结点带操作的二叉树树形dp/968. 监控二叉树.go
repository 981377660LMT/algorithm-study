// 监控二叉树 放照相
// 给定一个二叉树，我们在树的节点上安装摄像头。
// 节点上的每个摄影头都可以监视其父对象、自身及其直接子对象。
// 计算监控树的所有节点所需的最小摄像头数量。
//
// n<=1000

package main

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

const INF int = 1e18

func minCameraCover(root *TreeNode) int {
	var f func(*TreeNode) (int, int, int)
	f = func(node *TreeNode) (choose int, byParent int, byChildren int) {
		if node == nil {
			return INF, 0, 0 // 空节点不能安装摄像头，也无需被监控到
		}
		leftChoose, leftByParent, leftByChildren := f(node.Left)
		rightChoose, rightByParent, rightByChildren := f(node.Right)
		choose = min(leftChoose, leftByParent) + min(rightChoose, rightByParent) + 1
		byParent = min(leftChoose, leftByChildren) + min(rightChoose, rightByChildren)
		byChildren = min(leftChoose+rightByChildren, min(leftByChildren+rightChoose, leftChoose+rightChoose))
		return
	}

	choose, _, byChildren := f(root) // 根节点没有父节点
	return min(choose, byChildren)
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}
