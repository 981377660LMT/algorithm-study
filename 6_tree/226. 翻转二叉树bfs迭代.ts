import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 */
var invertTree = function (root: BinaryTree): BinaryTree {
  if (!root) return root
  const queue: BinaryTree[] = [root]
  while (queue.length) {
    const head = queue.shift()!
    ;[head.left, head.right] = [head.right, head.left]
    head.left && queue.push(head.left)
    head.right && queue.push(head.right)
  }
  return root
}

console.dir(invertTree(deserializeNode([4, 2, 7, 1, 3, 6, 9])!))
export {}
