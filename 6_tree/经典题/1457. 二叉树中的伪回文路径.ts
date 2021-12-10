import { BinaryTree } from '../力扣加加/Tree'
import { deserializeNode } from '../力扣加加/构建类/297二叉树的序列化与反序列化'

// 请你返回从根到叶子节点的所有路径中 伪回文 路径的数目。
// 伪回文:路径经过的所有节点值的排列中，存在一个回文序列。(奇数count至多一个)
// 给定二叉树的节点数目在范围 [1, 105] 内
// 1 <= Node.val <= 9   暗示状态压缩

function pseudoPalindromicPaths(root: BinaryTree | null): number {
  let res = 0
  dfs(root, 0)
  return res

  function dfs(root: BinaryTree | null, pathSum: number): void {
    if (!root) return

    if (!root.left && !root.right) {
      pathSum ^= 1 << root.val
      res += binCountOne(pathSum) <= 1 ? 1 : 0
    }

    dfs(root.left, pathSum ^ (1 << root.val))
    dfs(root.right, pathSum ^ (1 << root.val))
  }

  function binCountOne(num: number) {
    let res = 0
    while (num > 0) {
      num &= num - 1
      res++
    }

    return res
  }
}

console.log(pseudoPalindromicPaths(deserializeNode([2, 3, 1, 3, 1, null, 1])))
// 在这些路径中，只有红色和绿色的路径是伪回文路径，因为红色路径 [2,3,3] 存在回文排列 [3,2,3] ，绿色路径 [2,1,1] 存在回文排列 [1,2,1] 。
