import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297.二叉树的序列化与反序列化'

const isSymmetric = (root: BinaryTree): boolean => {
  if (!root) return true
  const helper = (root1: BinaryTree | null, root2: BinaryTree | null): boolean => {
    if (root1 == null && root2 == null) return true // If both sub trees are empty
    if (root1 == null || root2 == null) return false // If only one of the sub trees are empty
    if (root1.val !== root2.val) return false // If the values dont match up
    return helper(root1.left, root2.right) && helper(root1.right, root2.left)
  }
  return helper(root.left, root.right)
}

console.dir(isSymmetric(deserializeNode([1, 2, 2, 3, 4, 4, 3])!), { depth: null })

export {}
