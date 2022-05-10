import { deserializeNode } from './297.二叉树的序列化与反序列化'

class BinaryTree {
  val: number
  left: BinaryTree | null
  right: BinaryTree | null
  constructor(val: number) {
    this.val = val
    this.left = null
    this.right = null
  }
}
// 修改树
/**
 * @param {BinaryTree} root1
 * @param {BinaryTree} root2
 * @return {BinaryTree}
 */
// const mergeTrees = function (root1: BinaryTree, root2: BinaryTree): BinaryTree {
//   if (!root1) return root2
//   if (!root2) return root1
//   const preOrder = (root1: BinaryTree | null, root2: BinaryTree | null): BinaryTree | null => {
//     if (!root1) return root2
//     if (!root2) return root1
//     root1.val += root2.val
//     root1.left = preOrder(root1.left, root2.left)
//     root1.right = preOrder(root1.right, root2.right)
//     return root1
//   }
//   return preOrder(root1, root2)!
// }
// 不修改树
/**
 * @param {BinaryTree} root1
 * @param {BinaryTree} root2
 * @return {BinaryTree}
 */
const mergeTrees = function (
  root1: BinaryTree | null,
  root2: BinaryTree | null
): BinaryTree | null {
  if (!root1) return root2
  if (!root2) return root1
  const root = new BinaryTree(root1.val + root2.val)
  root.left = mergeTrees(root1.left, root2.left)
  root.right = mergeTrees(root1.right, root2.right)
  return root
}

console.dir(mergeTrees(deserializeNode([1, 2, 3])!, deserializeNode([2, 3, 4])!), { depth: null })
