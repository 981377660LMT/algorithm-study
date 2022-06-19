import { deserializeNode } from '../../重构json/297.二叉树的序列化与反序列化'
import { BinaryTree } from '../Tree'

// 假如通过节点 node 的每种可能的 “根-叶” 路径上值的总和全都小于给定的 limit，则该节点被称之为「不足节点」，需要被删除。
// 请你删除所有不足节点，并返回生成的二叉树的根。
function sufficientSubset(root: BinaryTree | null, limit: number): BinaryTree | null {
  const dummy = new BinaryTree(0, root)
  dfs(dummy, 0)
  return dummy.left

  function dfs(root: BinaryTree | null, curSum: number): number {
    if (!root) return -Infinity
    if (!root.left && !root.right) return curSum + root.val
    const left = dfs(root.left, curSum + root.val)
    const right = dfs(root.right, curSum + root.val)
    if (left < limit) root.left = null
    if (right < limit) root.right = null
    return Math.max(left, right)
  }
}

console.log(
  sufficientSubset(deserializeNode([1, 2, 3, 4, -99, -99, 7, 8, 9, -99, -99, 12, 13, -99, 14]), 1)
)
