import { BinaryTree, bt } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {number} low
 * @param {number} high
 * @return {BinaryTree}
 * @description 删除不在[low,high区间的节点]
 * @description 自底向上 后序
 */
const trimBST = function (root: BinaryTree | null, low: number, high: number): BinaryTree | null {
  if (!root) return null
  if (root.val > high) {
    // skip and go left
    return trimBST(root.left, low, high)
  } else if (root.val < low) {
    // skip and go right
    return trimBST(root.right, low, high)
  } else {
    // connect left and right child to the next qualified node
    root.left = trimBST(root.left, low, high)
    root.right = trimBST(root.right, low, high)
    return root
  }
}

console.dir(trimBST(deserializeNode([5, 3, 6, 2, 4, null, 7])!, 2, 4), { depth: null })
