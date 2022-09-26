import { deserializeNode } from '../../../重构json/297.二叉树的序列化与反序列化'
import { BinaryTree } from '../../Tree'

/**
 * @param {BinaryTree} root
 * @return {number}
 * @description 直径等于左右子树高度之和
 */
const diameterOfBinaryTree = function (root: BinaryTree): number {
  if (!root) return 0

  let max = 0
  dfs(root)
  return max

  function dfs(root: BinaryTree | null): number {
    if (!root) {
      return 0
    }

    const left = dfs(root.left)
    const right = dfs(root.right)
    max = Math.max(max, left + right)
    return Math.max(left, right) + 1
  }
}

console.log(diameterOfBinaryTree(deserializeNode([1, 2, 3, 4, 5, null, null])!))

export {}
