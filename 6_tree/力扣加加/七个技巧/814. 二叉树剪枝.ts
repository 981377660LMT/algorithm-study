import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 * @description 自底向上 返回移除了所有不包含 1 的子树的原二叉树(自底向上移除非1的叶子节点即可)
 */
const pruneTree = function (root: BinaryTree): BinaryTree | null {
  root.left && (root.left = pruneTree(root.left))
  root.right && (root.right = pruneTree(root.right))
  if (root.val === 0 && !root.left && !root.right) {
    return null
  } else {
    return root
  }
}

console.dir(pruneTree(deserializeNode([1, 0, 1, 0, 0, 0, 1])!), { depth: null })

export {}
