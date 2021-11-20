import { BinaryTree } from '../../力扣加加/Tree'
import { deserializeNode } from '../../力扣加加/构建类/297二叉树的序列化与反序列化'

// 母题
// 例如，如果路径为 0 -> 1 -> 1 -> 0 -> 1，那么它表示二进制数 01101，也就是 13 。
// 对树上的每一片叶子，我们都要找出从根到该叶子的路径所表示的数字。
// 返回这些数字之和。题目数据保证答案是一个 32 位 整数。

function sumRootToLeaf(root: BinaryTree | null): number {
  let res = 0
  dfs(root, 0)
  return res

  function dfs(root: BinaryTree | null, parent: number): void {
    if (!root) return

    const cur = (parent << 1) + root.val
    if (root.left == undefined && root.right == undefined) res += cur

    dfs(root.left, cur)
    dfs(root.right, cur)
  }
}

console.log(sumRootToLeaf(deserializeNode([1, 0, 1, 0, 1, 0, 1])))
