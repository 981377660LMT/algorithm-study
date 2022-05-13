import { BinaryTree } from '../Tree'
import { deserializeNode } from '../构建类/297.二叉树的序列化与反序列化'

/**
 * @param {BinaryTree} root
 * @param {number} distance
 * @return {number}
 * @description 如果二叉树中两个 叶 节点之间的 最短路径长度 小于或者等于 distance ，那它们就可以构成一组 好叶子节点对 。
 * @description 其实两个叶子节点的最短路径（距离）可以用其最近的公共祖先来辅助计算
 * 两个叶子节点的最短路径 = 其中一个叶子节点到最近公共祖先的距离 + 另外一个叶子节点到最近公共祖先的距离
 * 自底向上每个root处记录左边的[]与右边的[],向上返回两个[]的连接
 * @summary 计算以每个节点为root的距离左右数组 合并向上
 */
const countPairs = function (root: BinaryTree | null, distance: number): number {
  if (!root) return 0
  let res = 0

  const dfs = (root: BinaryTree | null): number[] => {
    if (!root) return []
    if (!root.left && !root.right) return [0]

    // 如果子树的结果计算出来了，那么父节点只需要把子树的每一项加 1 即可
    const left = dfs(root.left).map(v => v + 1)
    const right = dfs(root.right).map(v => v + 1)
    //  笛卡尔积
    for (const l of left) {
      for (const r of right) {
        if (l + r <= distance) res++
      }
    }

    return [...left, ...right]
  }
  dfs(root)

  return res
}

console.dir(countPairs(deserializeNode([1, 2, 3, 4, 5, 6, 7])!, 3), { depth: null })

export {}
