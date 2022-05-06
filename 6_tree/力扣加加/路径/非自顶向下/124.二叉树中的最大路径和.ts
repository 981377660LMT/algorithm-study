import { BinaryTree } from '../../Tree'
import { deserializeNode } from '../../构建类/297二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @return {number}
 * @description 自底向上
 * // cur三种路径情况: 1. left+cur 2. right+cur 3. left+cur+right 
   // 其中1,2是要往上探索的。3不能往上。
   // dfs的return值是当前节点[若往上, 即作为子节点]的最大贡献值，是不包含情况3的。
   // 但是3可能是最大路径，因此更新ans时是比较1, 2, 3中最大。
 */
function maxPathSum(root: BinaryTree | null): number {
  if (!root) return 0
  let res = -Infinity

  // 经过root的最大路径长
  const dfs = (root: BinaryTree | null): number => {
    if (!root) return 0
    const leftMax = dfs(root.left)
    const rightMax = dfs(root.right)
    // 和全局变量比较
    res = Math.max(res, root.val + leftMax + rightMax)
    // >0 说明都能使路径变大   类似于求树的高度
    return Math.max(0, Math.max(leftMax, rightMax) + root.val)
  }
  dfs(root)

  return res
}

console.log(maxPathSum(deserializeNode([-10, 9, 20, null, null, 15, 7])!))
