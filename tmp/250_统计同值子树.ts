import type { BinaryTree } from '../6_tree/力扣加加/Tree'

// 后序遍历
function countUnivalSubtrees(root: BinaryTree | null): number {
  if (!root) return 0
  let res = 0

  dfs(root, 0)
  return res

  function dfs(root: BinaryTree | null, parentValue: number): boolean {
    if (!root) return true

    const left = dfs(root.left, root.val)
    const right = dfs(root.right, root.val)

    if (left && right) {
      res++
      return root.val === parentValue
    }

    return false
  }
}
