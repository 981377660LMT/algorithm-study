import { BinaryTree } from '../../Tree'
import { deserializeNode } from '../../构建类/297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {number}
 * @description 直径等于左右子树高度之和
 */
const diameterOfBinaryTree = function (root: BinaryTree): number {
  if (!root) return 0
  let max = 0

  const dfs = (root: BinaryTree | null): number => {
    if (!root) return 0
    const left = dfs(root.left)
    const right = dfs(root.right)
    max = Math.max(max, left + right)
    return Math.max(left, right) + 1
  }
  dfs(root)

  return max
}
console.log(diameterOfBinaryTree(deserializeNode([1, 2, 3, 4, 5, null, null])!))

export {}
