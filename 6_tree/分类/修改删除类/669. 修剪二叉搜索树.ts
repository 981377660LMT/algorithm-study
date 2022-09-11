/* eslint-disable no-param-reassign */

import { deserializeNode } from '../../重构json/297.二叉树的序列化与反序列化'
import { BinaryTree } from '../Tree'

/**
 * @param {BinaryTree} root
 * @param {number} low
 * @param {number} high
 * @return {BinaryTree}
 * @description 删除不在[low,high区间的节点]
 * @description 自底向上 后序
 */
function trimBST(root: BinaryTree | null, low: number, high: number): BinaryTree | null {
  if (!root) {
    return null
  }

  if (root.val > high) {
    return trimBST(root.left, low, high)
  }

  if (root.val < low) {
    return trimBST(root.right, low, high)
  }

  root.left = trimBST(root.left, low, high)
  root.right = trimBST(root.right, low, high)
  return root
}

console.dir(trimBST(deserializeNode([5, 3, 6, 2, 4, null, 7])!, 2, 4), { depth: null })
