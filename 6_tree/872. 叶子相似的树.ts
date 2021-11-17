import { BinaryTree } from './力扣加加/Tree'
import { deserializeNode } from './力扣加加/构建类/297二叉树的序列化与反序列化'

/**
 *
 * @param root1
 * @param root2
 * 一棵二叉树上所有的叶子，这些叶子的值按从左到右的顺序排列形成一个 叶值序列
 * 有两棵二叉树的叶值序列是相同，那么我们就认为它们是 叶相似 的
 */
function leafSimilar(root1: BinaryTree | null, root2: BinaryTree | null): boolean {
  const leaves1 = getLeaves(root1)
  const leaves2 = getLeaves(root2)
  return isSameLeaves(leaves1, leaves2)

  function getLeaves(root: BinaryTree | null): number[] {
    const res: number[] = []
    dfs(root)
    return res

    function dfs(root: BinaryTree | null): void {
      if (!root) return
      if (!root.left && !root.right) res.push(root.val)
      root.left && dfs(root.left)
      root.right && dfs(root.right)
    }
  }

  function isSameLeaves(leaves1: number[], leaves2: number[]): boolean {
    if (leaves1.length !== leaves2.length) return false

    for (let i = 0; i < leaves1.length; i++) {
      if (leaves1[i] !== leaves2[i]) return false
    }

    return true
  }
}
// 层序遍历分别处理两个root即可

console.log(
  leafSimilar(
    deserializeNode([3, 5, 1, 6, 2, 9, 8, null, null, 7, 4]),
    deserializeNode([3, 5, 1, 6, 7, 4, 2, null, null, null, null, null, null, 9, 8])
  )
)
