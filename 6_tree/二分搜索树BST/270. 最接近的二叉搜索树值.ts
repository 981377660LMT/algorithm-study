// 思路和心得：

import { BinaryTree } from '../分类/Tree'

// 1.贪心，向target所在的区间走

// 2.做好统计和更新

// 3.发现相同的直接返回就ok
function closestValue(root: BinaryTree | null, target: number): number {
  let res = root!.val
  while (root) {
    if (Math.abs(root.val - target) < Math.abs(res - target)) res = root.val
    if (root.val === target) return root.val
    else if (root.val < target) root = root.right
    else root = root.left
  }
  return res
}
