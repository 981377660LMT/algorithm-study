import { BinaryTree } from '../Tree'

/**
 *
 * @param root
 * @param subRoot
 * 检验 root 中是否包含和 subRoot 具有相同结构和节点值的子树
 */
function isSubtree(root: BinaryTree | null, subRoot: BinaryTree | null): boolean {
  if (!root && !subRoot) return true
  if (!root || !subRoot) return false
  return (
    isSameTree(root, subRoot) || isSubtree(root.left, subRoot) || isSubtree(root.right, subRoot)
  )
}

// 快速比较两个json
function isSameTree(t1: BinaryTree | null, t2: BinaryTree | null): boolean {
  // 递归的终点
  if (!t1 && !t2) return true
  if (!t1 || !t2) return false
  return t1.val === t2.val && isSameTree(t1.left, t2.left) && isSameTree(t1.right, t2.right)
}
