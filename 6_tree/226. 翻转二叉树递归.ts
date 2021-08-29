import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 */
var invertTree = function (root: BinaryTree): BinaryTree {
  const helper = (root: BinaryTree | null) => {
    if (!root) return
    helper(root.left)
    helper(root.right)
    ;[root.left, root.right] = [root.right, root.left]
  }
  helper(root)
  return root
}

console.dir(invertTree(deserializeNode([4, 2, 7, 1, 3, 6, 9])!))
export {}
