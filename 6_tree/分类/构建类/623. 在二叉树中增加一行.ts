/* eslint-disable no-param-reassign */
import { BinaryTree } from '../Tree'

type TreeNode = BinaryTree

/**
 * @param root  根节点为第1层，深度为 1
 * @param val
 * @param depth
 *
 * 在其第 d 层追加一行值为 v 的节点。
 * 给定一个深度值 d （正整数），针对深度为 d-1 层的每一非空节点 N，为 N 创建两个值为 v 的左子树和右子树
 * 如果 d 的值为 1，深度 d - 1 不存在，则创建一个新的根节点 v，原先的整棵树将作为 v 的左子树
 */
function addOneRow(root: TreeNode | null, val: number, depth: number): TreeNode | null {
  if (!root) return null
  if (depth === 1) return new TreeNode(val, root, null)
  if (depth === 2) {
    root.left = new TreeNode(val, root.left, null)
    root.right = new TreeNode(val, null, root.right)
  } else {
    root.left && addOneRow(root.left, val, depth - 1)
    root.right && addOneRow(root.right, val, depth - 1)
  }

  return root
}
