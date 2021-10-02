// 你计算其中 最长连续序列路径 的长度。

import { BinaryTree } from '../../Tree'

// 必须从父节点到子节点
function longestConsecutive(root: BinaryTree | null): number {
  if (!root) return 0

  const dfs = (root: BinaryTree | null, parentVal: number, dis: number): number => {
    if (!root) return dis

    if (parentVal + 1 === root.val) dis++
    else dis = 1

    return Math.max(dis, dfs(root.left, root.val, dis), dfs(root.right, root.val, dis))
  }

  return dfs(root, Infinity, 0)
}

export {}
