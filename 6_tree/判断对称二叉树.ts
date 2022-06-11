import { BinaryTree } from './分类/Tree'
import { deserializeNode } from './重构json/297.二叉树的序列化与反序列化'

function isSymmetric(root: BinaryTree): boolean {
  if (!root) return true
  return dfs(root.left, root.right)

  function dfs(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
    if (root1 == null && root2 == null) return true // If both sub trees are empty
    if (root1 == null || root2 == null) return false // If only one of the sub trees are empty
    if (root1.val !== root2.val) return false // If the values dont match up
    return dfs(root1.left, root2.right) && dfs(root1.right, root2.left)
  }
}

console.dir(isSymmetric(deserializeNode([1, 2, 2, 3, 4, 4, 3])!), { depth: null })

export {}
