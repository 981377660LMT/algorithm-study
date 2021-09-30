import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

/**
 *
 * @param root
 * @param p
 * 给定一棵二叉搜索树root和其中的一个节点 p ，找到该节点在树中的中序后继
 * 如果节点没有中序后继，请返回 null
 */
function inorderSuccessor(root: BinaryTree | null, p: BinaryTree | null): BinaryTree | null {
  if (!root || !p) return null

  let pre: BinaryTree | null = null

  function* inorder(root: BinaryTree | null): Generator<BinaryTree | null> {
    if (!root) return
    yield* inorder(root.left)
    if (pre?.val === p?.val) yield root
    pre = root
    yield* inorder(root.right)
  }

  return inorder(root).next().value || null
}

console.log(inorderSuccessor(deserializeNode([2, 1, 3]), deserializeNode([1])))
console.log(inorderSuccessor(deserializeNode([5, 3, 6, 2, 4, null, null, 1]), deserializeNode([1])))
