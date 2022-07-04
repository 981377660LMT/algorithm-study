// 给定一棵树的前序遍历 preorder 与中序遍历  inorder。请构造二叉树并返回其根节点。

import { BinaryTree } from '../Tree'

/**
 * @param {number[]} preorder 1 <= preorder.length <= 3000
 * @param {number[]} inorder
 * @return {TreeNode}
 * @description !假设输入的遍历的序列中都不含重复的数字
 * @description 由前序遍历可知preorder数组中第一个数一定是root并弹出，
 * 根据root值在inorder所在位置可将inorder划分为左子树、右子树两部分。
 * 注意inorder 和 postorder一定是长度相等的
 * @summary
 * !优化1: 不用slice 数组 可以直接记录边界的index
 * !优化2: 不用线性查找rootIndex 因为不含重复数字 可以借助哈希表
 * 时间复杂度 O(n)
 */
function buildTree(preorder: number[], inorder: number[]): BinaryTree | null {
  const n = preorder.length
  const indexMap = new Map<number, number>()
  inorder.forEach((value, index) => indexMap.set(value, index))
  return dfs(0, n - 1, 0, n - 1)

  function dfs(pLeft: number, pRight: number, iLeft: number, iRight: number): BinaryTree | null {
    if (pLeft > pRight || iLeft > iRight) return null

    const rootValue = preorder[pLeft]
    const root = new BinaryTree(rootValue)
    const rootIndex = indexMap.get(rootValue)!
    const [leftLen, rightLen] = [rootIndex - iLeft, iRight - rootIndex]

    root.left = dfs(pLeft + 1, pLeft + leftLen, iLeft, rootIndex - 1)
    root.right = dfs(pRight - rightLen + 1, pRight, rootIndex + 1, iRight)
    return root
  }
}

console.log(buildTree([3, 9, 20, 15, 7], [9, 3, 15, 20, 7]))

export { buildTree }
// 时间复杂度：由于每次递归我们的 inorder 和 preorder 的总数都会减 1，因此我们要递归 N 次，故时间复杂度为 $O(N)$，其中 N 为节点个数。
// 空间复杂度：我们使用了递归，也就是借助了额外的栈空间来完成， 由于栈的深度为 N，因此总的空间复杂度为 $O(N)$，其中 N 为节点个数。
