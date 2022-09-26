/* eslint-disable @typescript-eslint/no-unused-vars */
/* eslint-disable no-shadow */

import { BinaryTree } from '../Tree'

// !求该二叉树里节点值之和等于 targetSum 的 路径 的数目。
// !路径方向必须是向下的（只能从父节点到子节点）。
// !O（n^2）的解法
function pathSum(root: BinaryTree | null, target: number): number {
  if (root == null) {
    return 0
  }

  // !以每个root为起点统计
  return dfs(root, target) + pathSum(root.left, target) + pathSum(root.right, target)

  // !以root为起点，统计路径和为sum的路径数
  function dfs(root: BinaryTree | null, sum: number): number {
    if (root == null) {
      return 0
    }

    sum -= root.val
    return (sum === 0 ? 1 : 0) + dfs(root.left, sum) + dfs(root.right, sum)
  }
}

export {}
