import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297二叉树的序列化与反序列化'

/**
 *
 * @param n
 * @returns
 * @description 满二叉树是一类二叉树，其中每个结点恰好有 0 或 2 个子结点。
 * @description 自顶向下记忆化搜索 分解成左边和右边 然后组合
 * 7 个节点时:左边1+右边5，左边3加右边3 左边5加右边 1
 */
const allPossibleFBT = function (n: number): (BinaryTree | null)[] {
  if (n % 2 == 0) {
    return []
  }

  const memo = new Map<number, (BinaryTree | null)[]>()

  const dfs = (n: number): (BinaryTree | null)[] => {
    if (memo.has(n)) return memo.get(n)!
    if (n === 1) return [new BinaryTree(0)]

    const res: (BinaryTree | null)[] = []
    for (let i = 1; i < n - 1; i += 2) {
      const left = dfs(i)
      const right = dfs(n - i - 1)
      // 左右配对所有可能
      for (const l of left) {
        for (const r of right) {
          const root = new BinaryTree(0)
          root.left = l
          root.right = r
          res.push(root)
        }
      }
    }

    memo.set(n, res)
    return res
  }

  return dfs(n)
}

console.dir(allPossibleFBT(7), { depth: null })

export {}
