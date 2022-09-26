import { BinaryTree } from '../Tree'

/**
 * @param {BinaryTree} root
 * @return {number}
 */
function diameterOfBinaryTree(root: BinaryTree): number {
  if (!root) {
    return 0
  }

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

export {}
