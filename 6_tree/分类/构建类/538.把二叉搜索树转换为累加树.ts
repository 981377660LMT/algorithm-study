import { BinaryTree } from '../Tree'
import { deserializeNode } from './297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 * 从下到上，从右到左的顺序累加
 * 累加树（Greater Sum Tree）
 * BST的中序遍历就是从小到大,那么反过来就是从大到小,然后累加就好了.(右中左)
 */
const bstToGst = function (root: BinaryTree): BinaryTree {
  let sum = 0
  const reversedInOrder = (root: BinaryTree | null) => {
    if (!root) return
    reversedInOrder(root.right)
    const tmp = root.val
    root.val += sum
    sum += tmp
    reversedInOrder(root.left)
  }
  reversedInOrder(root)
  return root
}

console.dir(bstToGst(deserializeNode([4, 1, 6, 0, 2, 5, 7])!), { depth: null })
