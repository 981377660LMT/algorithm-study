import { BinaryTree } from '../../6_tree/力扣加加/Tree'

// 请你删除 1 条边，使二叉树分裂成两棵子树，且它们子树和的乘积尽可能大。
// 乘积 = 某个节点下所有子节点的和 *（整个树的和 - 某个节点下所有子节点的和）
const MOD = 1e9 + 7
function maxProduct(root: BinaryTree | null): number {
  const total = sum(root)

  let res = 0

  dfs(root)
  return res % MOD

  function sum(root: BinaryTree | null): number {
    if (!root) return 0
    return root.val + sum(root.left) + sum(root.right)
  }

  function dfs(root: BinaryTree | null): number {
    if (!root) return 0
    const left = dfs(root.left)
    const right = dfs(root.right)
    res = Math.max(res, left * (total - left), right * (total - right))
    return root.val + left + right
  }
}

export {}
