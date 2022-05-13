import { BinaryTree } from '../Tree'
import { buildTree } from './105. 从前序与中序遍历序列构造二叉树'

/**
 * @param {number[]} preorder
 * @return {BinaryTree}
 * @description 相当于告诉了中序遍历，等价于前序+中序唯一确定
 */
const bstFromPreorder = (preorder: number[]): BinaryTree | null => {
  const inorder = preorder.slice().sort((a, b) => a - b)
  return buildTree(preorder, inorder)
}

console.dir(bstFromPreorder([8, 5, 1, 7, 10, 12]), { depth: null })
