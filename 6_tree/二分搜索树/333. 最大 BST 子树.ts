import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

// 给定一个二叉树，找到其中结点最多的二叉搜索树（BST）子树，并返回该子树的大小。
// 时间复杂度为 O(n)
function largestBSTSubtree(root: BinaryTree | null): number {
  if (!root) return 0

  function dfs(root: BinaryTree | null): [min: number, max: number, count: number] {
    // 什么祖先节点都可以接
    if (!root) return [Infinity, -Infinity, 0]

    const [leftMin, leftMax, leftCount] = dfs(root.left)
    const [rightMin, rightMax, rightCount] = dfs(root.right)

    // 前驱后继
    if (root.val > leftMax && root.val < rightMin) {
      return [Math.min(leftMin, root.val), Math.max(rightMax, root.val), 1 + leftCount + rightCount]
    }

    // 什么祖先节点都不能接，不用判断了
    return [-Infinity, Infinity, Math.max(leftCount, rightCount)]
  }

  return dfs(root)[2]
}

console.log(largestBSTSubtree(deserializeNode([10, 5, 15, 1, 8, null, 7])))
