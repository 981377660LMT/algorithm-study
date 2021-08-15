// 给定一棵树的前序遍历 preorder 与中序遍历  inorder。请构造二叉树并返回其根节点。

import { BinaryTree } from '../Tree'

/**
 * @param {number[]} preorder
 * @param {number[]} inorder
 * @return {TreeNode}
 * @description 假设输入的遍历的序列中都不含重复的数字(找根要找对)
 * @description 由前序遍历可知preorder数组中第一个数一定是root并弹出，
 * 根据root值在inorder所在位置可将inorder划分为左子树、右子树两部分。
 * 注意inorder 和 postorder一定是长度相等的
 */
const buildTree = (preorder: number[], inorder: number[]): BinaryTree | null => {
  // # 实际上inorder 和 postorder一定是长度相等的
  if (!preorder.length || !inorder.length) return null

  // 先找根
  const rootValue = preorder[0]
  const root = new BinaryTree(rootValue)
  const rootIndex = inorder.indexOf(rootValue)
  // 去除根，然后包括进左子树的个数
  root.left = buildTree(preorder.slice(1, rootIndex + 1), inorder.slice(0, rootIndex))
  root.right = buildTree(preorder.slice(rootIndex + 1), inorder.slice(rootIndex + 1))

  return root
}

console.log(buildTree([3, 9, 20, 15, 7], [9, 3, 15, 20, 7]))

export { buildTree }
// 时间复杂度：由于每次递归我们的 inorder 和 preorder 的总数都会减 1，因此我们要递归 N 次，故时间复杂度为 $O(N)$，其中 N 为节点个数。
// 空间复杂度：我们使用了递归，也就是借助了额外的栈空间来完成， 由于栈的深度为 N，因此总的空间复杂度为 $O(N)$，其中 N 为节点个数。
