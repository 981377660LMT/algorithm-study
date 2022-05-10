import { BinaryTree } from '../6_tree/力扣加加/Tree'
import { deserializeNode } from '../6_tree/力扣加加/构建类/297.二叉树的序列化与反序列化'

// 原来的左子节点变成新的根节点
// 原来的根节点变成新的右子节点
// 原来的右子节点变成新的左子节点
// 关键是记录parent 和之前的右节点
function upsideDownBinaryTree(root: BinaryTree | null): BinaryTree | null {
  let preRight: BinaryTree | null = null
  let parent: BinaryTree | null = null

  while (root) {
    const left = root.left
    const right = root.right

    root.left = preRight
    root.right = parent

    preRight = right
    parent = root
    root = left
  }

  return parent
}

console.log(upsideDownBinaryTree(deserializeNode([1, 2, 3, 4, 5])))
