import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {number} target
 * @return {BinaryTree}
 */
const removeLeafNodes = function (root: BinaryTree, target: number): BinaryTree | null {
  root.left && (root.left = removeLeafNodes(root.left, target))
  root.right && (root.right = removeLeafNodes(root.right, target))
  if (root.val === target && !root.left && !root.right) {
    return null
  } else {
    return root
  }
}

console.dir(removeLeafNodes(deserializeNode([1, 2, 3, 2, null, 2, 4])!, 2), { depth: null })

export {}
