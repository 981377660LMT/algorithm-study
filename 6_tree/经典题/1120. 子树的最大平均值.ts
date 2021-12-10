import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

// 找出这棵树的 每一棵 子树的 平均值 中的 最大 值。

function maximumAverageSubtree(root: BinaryTree | null): number {
  if (!root) return 0

  let res = -Infinity
  dfs(root)
  return res

  function dfs(root: BinaryTree | null): [sum: number, count: number] {
    if (!root) return [0, 0]

    const [leftSum, leftCount] = dfs(root.left)
    const [rightSum, rightCount] = dfs(root.right)
    const curSum = leftSum + rightSum + root.val
    const curCount = leftCount + rightCount + 1
    res = Math.max(res, curSum / curCount)
    return [curSum, curCount]
  }
}

console.log(maximumAverageSubtree(deserializeNode([5, 6, 1])))
console.log(maximumAverageSubtree(deserializeNode([55, 9, 7])))
