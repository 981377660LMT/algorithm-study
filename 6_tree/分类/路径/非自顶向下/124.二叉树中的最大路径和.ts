/* eslint-disable no-shadow */
import { deserializeNode } from '../../../重构json/297.二叉树的序列化与反序列化'
import { BinaryTree } from '../../Tree'

const INF = 2e15

function maxPathSum(root: BinaryTree | null): number {
  if (!root) return 0

  let res = -INF
  dfs(root)
  return res

  // 经过root的最大路径长
  function dfs(root: BinaryTree | null): number {
    if (!root) {
      return 0
    }

    const leftMax = dfs(root.left)
    const rightMax = dfs(root.right)
    res = Math.max(res, root.val + leftMax + rightMax) // !经过当前节点的最长路径
    return Math.max(0, Math.max(leftMax, rightMax) + root.val) // !子树往上的贡献值
  }
}

console.log(maxPathSum(deserializeNode([-10, 9, 20, null, null, 15, 7])!))
