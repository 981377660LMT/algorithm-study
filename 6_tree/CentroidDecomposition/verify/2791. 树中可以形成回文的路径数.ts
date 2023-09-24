/* eslint-disable no-loop-func */
// https://leetcode.cn/problems/count-paths-that-can-form-a-palindrome-in-a-tree/
// 给你一棵 树（即，一个连通、无向且无环的图），根 节点为 0 ，由编号从 0 到 n - 1 的 n 个节点组成。
// 这棵树用一个长度为 n 、下标从 0 开始的数组 parent 表示，其中 parent[i] 为节点 i 的父节点，由于节点 0 为根节点，所以 parent[0] == -1 。
// 另给你一个长度为 n 的字符串 s ，其中 s[i] 是分配给 i 和 parent[i] 之间的边的字符。s[0] 可以忽略。
// !找出满足 u < v ，且从 u 到 v 的路径上分配的字符可以 重新排列 形成 回文 的所有节点对 (u, v) ，并返回节点对的数目。

import { centroidDecomposition } from '../CentroidDecomposition'

function countPalindromePaths(parent: number[], s: string): number {
  const n = parent.length
  const tree: [next: number, cost: number][][] = Array(n)
  for (let i = 0; i < n; i++) tree[i] = []
  for (let cur = 1; cur < n; cur++) {
    const p = parent[cur]
    const cost = 1 << (s.charCodeAt(cur) - 97)
    tree[p].push([cur, cost])
    tree[cur].push([p, cost])
  }

  // 全局状态
  const [centTree, root] = centroidDecomposition(n, tree)
  const removed = new Uint8Array(n)
  let res = 0
  decomposition(root, -1)
  return res

  /**
   * 点分治，对某个点，考虑包含这个点的路径和不包含这个点的路径.
   * 不包含这个点的路径，删除这个点，然后对每个子树递归求解.
   * 包含这个点的路径，可以用dfs收集子树信息，然后根据每个子树内的信息，计算答案.
   */
  function decomposition(cur: number, pre: number): void {
    // 点分树的子树中的答案(不经过重心)
    removed[cur] = 1
    const nexts_ = centTree[cur]
    for (let i = 0; i < nexts_.length; i++) {
      const next = nexts_[i]
      if (!removed[next]) {
        decomposition(next, cur)
      }
    }
    removed[cur] = 0

    // !经过重心的路径
    const counter = new Map<number, number>([[0, 1]])
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      const cost = nexts[i][1]
      if (next === pre || removed[next]) continue

      const sub = new Map<number, number>() // state -> count, 统计子树内(不含cur)
      collect(next, cur, cost, sub)
      sub.forEach((count, state) => {
        res += count * (counter.get(state) || 0)
        for (let j = 0; j < 26; j++) {
          res += count * (counter.get(state ^ (1 << j)) || 0)
        }
      })
      sub.forEach((count, state) => {
        counter.set(state, (counter.get(state) || 0) + count)
      })
    }
  }

  /**
   * 收集子树信息.
   */
  function collect(cur: number, pre: number, state: number, sub: Map<number, number>): void {
    sub.set(state, (sub.get(state) || 0) + 1)
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      const cost = nexts[i][1]
      if (next !== pre && !removed[next]) {
        collect(next, cur, state ^ cost, sub)
      }
    }
  }
}
