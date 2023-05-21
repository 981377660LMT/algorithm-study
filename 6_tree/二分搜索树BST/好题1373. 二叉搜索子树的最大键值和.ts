/* eslint-disable no-shadow */
import { BinaryTree } from '../分类/Tree'
import { deserializeNode } from '../重构json/297.二叉树的序列化与反序列化'

type Return = [isValidBST: boolean, min: number, max: number, sum: number]
const INF = 2e15

/**
 * @param {BinaryTree} root
 * @return {number}
 * 给你一棵以 root 为根的 二叉树 ，请你返回 任意 二叉搜索子树的最大键值和。
 */
function maxSumBST(root: BinaryTree): number {
  if (!root) return 0
  let maxSum = 0
  dfs(root)
  return maxSum

  function dfs(root: BinaryTree | null): Return {
    if (!root) return [true, INF, -INF, 0]

    const res: Return = [true, INF, -INF, 0]
    const left = dfs(root.left)
    const right = dfs(root.right)
    if (left[0] && right[0] && root.val > left[2] && root.val < right[1]) {
      res[0] = true // 以 root 为根的二叉树是 BST
      res[1] = Math.min(left[1], root.val) // 计算以 root 为根的这棵 BST 的最小值
      res[2] = Math.max(right[2], root.val) // 计算以 root 为根的这棵 BST 的最大值
      res[3] = left[3] + right[3] + root.val // 计算以 root 为根的这棵 BST 所有节点之和
      maxSum = Math.max(maxSum, res[3])
    } else {
      res[0] = false // 其他的值都没必要计算了，因为用不到
    }

    return res
  }
}

console.log(
  maxSumBST(deserializeNode([1, 4, 3, 2, 4, 2, 5, null, null, null, null, null, null, 4, 6])!)
)
// 输出：20
// 解释：键值为 3 的子树是和最大的二叉搜索树。

// dfs函数在遍历二叉树的同时顺便把之前辅助函数做的事情都做了，
// 避免了在递归函数中调用递归函数，时间复杂度只有 O(N)。
// 如果当前节点要做的事情需要通过左右子树的计算结果推导出来，就要用到后序遍历。
