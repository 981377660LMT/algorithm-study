/**
 * Definition for a binary tree node.
 * class TreeNode {
 *     val: number
 *     left: TreeNode | null
 *     right: TreeNode | null
 *     constructor(val?: number, left?: TreeNode | null, right?: TreeNode | null) {
 *         this.val = (val===undefined ? 0 : val)
 *         this.left = (left===undefined ? null : left)
 *         this.right = (right===undefined ? null : right)
 *     }
 * }
 */

import { useMemoMap } from '../../../5_map/memo'
import { BinaryTree } from '../../../6_tree/力扣加加/Tree'
import { treeToGraph } from '../../../6_tree/力扣加加/构建类/treeToGraph'

// 二叉树转图的应用
function maxPathSum(root: BinaryTree): number {
  const { adjMap, valueById } = treeToGraph(root)
  const dfs = useMemoMap((cur: number, parent: number): number => {
    let res = valueById.get(cur)!
    for (const next of adjMap.get(cur) ?? []) {
      if (next === parent) continue
      res = Math.max(res, dfs(next, cur) + valueById.get(cur)!)
    }
    return res
  })

  let res = root.val
  for (const start of valueById.keys()) {
    res = Math.max(res, dfs(start, -1))
  }
  return res
}
