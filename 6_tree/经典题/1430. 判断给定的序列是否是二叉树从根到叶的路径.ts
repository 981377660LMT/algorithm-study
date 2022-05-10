import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297.二叉树的序列化与反序列化'

// 从根节点到任意叶节点的任意路径中的节点值所构成的序列为该二叉树的一个 “有效序列” 。
// 检查一个给定的序列是否是给定二叉树的一个 “有效序列” 。

// 总结：判断当前node.val和arr[index]的节点值是否相等，不等直接返回false；
// 相等则进入下一层，递归判断左右子树中，是否存在node.val和下一个元素相等。
// https://leetcode-cn.com/problems/check-if-a-string-is-a-valid-sequence-from-root-to-leaves-path-in-a-binary-tree/comments/377985
function isValidSequence(root: BinaryTree | null, arr: number[]): boolean {
  return dfs(root, arr, 0)
  function dfs(root: BinaryTree | null, arr: number[], index: number): boolean {
    if (!root) return false
    if (root.val !== arr[index]) return false
    if (index === arr.length - 1) return root.left == undefined && root.right == undefined
    return dfs(root.left, arr, index + 1) || dfs(root.right, arr, index + 1)
  }
}

console.log(isValidSequence(deserializeNode([0, 1, 0, 0, 1, 0, null, null, 1, 0, 0]), [0, 1, 0, 1]))
// 输出：true
// 解释：
// 路径 0 -> 1 -> 0 -> 1 是一个“有效序列”（图中的绿色节点）。
// 其他的“有效序列”是：
// 0 -> 1 -> 1 -> 0
// 0 -> 0 -> 0
