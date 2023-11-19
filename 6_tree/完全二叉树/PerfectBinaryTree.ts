/* eslint-disable no-inner-declarations */

/**
 * 完全二叉树.
 * 根节点编号为1,左子节点编号为2*x,右子节点编号为2*x+1,父结点编号为x>>>1.
 */
class PerfectBinaryTree {
  static depth(u: number): number {
    if (u === 0) return 0
    return 31 - Math.clz32(u)
  }

  /** 完全二叉树中两个节点的最近公共祖先(两个二进制数字的最长公共前缀). */
  static lca(u: number, v: number): number {
    if (u === v) return u
    if (u > v) {
      const tmp = u
      u = v
      v = tmp
    }
    const depth1 = this.depth(u)
    const depth2 = this.depth(v)
    const diff = u ^ (v >>> (depth2 - depth1))
    if (diff === 0) return u
    const len = 32 - Math.clz32(diff)
    return u >>> len
  }

  static dist(u: number, v: number): number {
    return this.depth(u) + this.depth(v) - 2 * this.depth(this.lca(u, v))
  }
}

export { PerfectBinaryTree }

if (require.main === module) {
  // https://leetcode.cn/problems/cycle-length-queries-in-a-tree/description/
  function cycleLengthQueries(n: number, queries: number[][]): number[] {
    return queries.map(([root1, root2]) => PerfectBinaryTree.dist(root1, root2) + 1)
  }
}
