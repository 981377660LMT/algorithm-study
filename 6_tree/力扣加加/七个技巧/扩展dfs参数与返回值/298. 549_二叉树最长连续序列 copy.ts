import { BinaryTree } from '../../Tree'
import { deserializeNode } from '../../构建类/297二叉树的序列化与反序列化'
import { treeToGraph } from '../../构建类/treeToGraph'

// 你计算其中 最长连续序列路径 的长度。
// 路径可以是 子-父-子 顺序，并不一定是 父-子 顺序。

function longestConsecutive(root: BinaryTree | null): number {
  if (!root) return 0

  let res = 1

  /**
   *
   * @param root 经过root的最长连续路径 (连续并不是单调 而是1 3 2 3 这种)
   * @returns
   */
  const dfs = (root: BinaryTree | null): [child: number, down: number, up: number] => {
    if (!root) return [Infinity, 0, 0]
    const [pre1, left1, right1] = dfs(root.left)
    const [pre2, left2, right2] = dfs(root.right)
    let [left, right] = [0, 0]

    if (pre1 === root.val - 1) {
      left = Math.max(left, 1 + left1)
    }

    if (pre2 === root.val - 1) {
      left = Math.max(left, 1 + left2)
    }

    if (pre1 === root.val + 1) {
      right = Math.max(right, 1 + right1)
    }

    if (pre2 === root.val + 1) {
      right = Math.max(right, 1 + right2)
    }

    res = Math.max(res, 1 + left + right)
    return [root.val, left, right]
  }

  dfs(root)
  return res
}

export {}
console.log(longestConsecutive(deserializeNode([2, 1, 3])))
