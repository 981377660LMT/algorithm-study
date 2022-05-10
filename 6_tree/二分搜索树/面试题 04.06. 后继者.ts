import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {BinaryTree} p
 * @return {BinaryTree}
 * 找出二叉搜索树中指定节点的“下一个”节点（也即中序后继）。
 */
var inorderSuccessor = function (root: BinaryTree, p: BinaryTree): BinaryTree {
  let pre: BinaryTree | null = null

  function* inorder(root: BinaryTree | null): Generator<BinaryTree> {
    if (!root) return
    yield* inorder(root.left)
    if (pre && pre.val === p.val) yield root
    pre = root
    yield* inorder(root.right)
  }

  return inorder(root).next().value ?? null
}

console.log(inorderSuccessor(deserializeNode([5, 3, 6, 2, 4, null, null, 1])!, new BinaryTree(6)))
console.log(inorderSuccessor(deserializeNode([2, 1, 3])!, new BinaryTree(1)))
