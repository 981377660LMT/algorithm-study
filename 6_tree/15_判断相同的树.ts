import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297.二叉树的序列化与反序列化'

const isSameTree = (root1: BinaryTree | null, root2: BinaryTree | null): boolean => {
  if (root1 == null && root2 == null) return true
  if (root1 == null || root2 == null) return false // If only one of the sub trees are empty

  return (
    root1.val === root2.val &&
    isSameTree(root1.left, root2.left) &&
    isSameTree(root1.right, root2.right)
  )
}

console.dir(
  isSameTree(deserializeNode([1, 2, 2, 3, 4, 4, 3]), deserializeNode([1, 2, 2, 3, 4, 4, 3])),
  { depth: null }
)

export {}
