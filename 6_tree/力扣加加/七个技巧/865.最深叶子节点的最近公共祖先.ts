import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {BinaryTree}
 * @description 自底向上 返回移除了所有不包含 1 的子树的原二叉树(自底向上移除非1的叶子节点即可)
 */
const subtreeWithAllDeepest = function (root: BinaryTree | null): BinaryTree | null {
  if (!root) return root

  const getMaxDepth = (root: BinaryTree | null): number => {
    if (!root) return 0
    return Math.max(getMaxDepth(root.left), getMaxDepth(root.right)) + 1
  }
  const leftDepth = getMaxDepth(root.left)
  const rightDepth = getMaxDepth(root.right)

  if (leftDepth === rightDepth) {
    return root
  } else if (leftDepth < rightDepth) {
    return subtreeWithAllDeepest(root.right)
  } else {
    return subtreeWithAllDeepest(root.left)
  }
}

console.dir(subtreeWithAllDeepest(deserializeNode([3, 5, 1, 6, 2, 0, 8, null, null, 7, 4])!), {
  depth: null,
})

export {}
