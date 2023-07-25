/* eslint-disable no-loop-func */
// https://leetcode.cn/problems/number-of-good-paths/
// 一条 好路径需要满足以下条件：

// !开始节点和结束节点的值 相同。
// !开始节点和结束节点中间的所有节点值都小于等于开始节点的值（也就是说开始节点的值应该是路径上所有节点的最大值）。
// 请你返回不同好路径的数目。

import { centroidDecomposition } from '../CentroidDecomposition'

function numberOfGoodPaths(vals: number[], edges: number[][]): number {
  const n = vals.length
  const tree: [next: number, cost: number][][] = Array(n)
  for (let i = 0; i < n; i++) tree[i] = []
  edges.forEach(([u, v]) => {
    tree[u].push([v, 1])
    tree[v].push([u, 1])
  })

  // 全局状态
  const [centTree, root] = centroidDecomposition(n, tree)
  const removed = new Uint8Array(n)
  let res = 0
  decomposition(root, -1)
  return res + n // 一个结点的路径

  /**
   * 点分治，对某个点，考虑包含这个点的路径和不包含这个点的路径.
   * 不包含这个点的路径，删除这个点，然后对每个子树递归求解.
   * 包含这个点的路径，可以用dfs收集子树信息，然后根据每个子树内的信息，计算答案.
   */
  function decomposition(cur: number, pre: number): void {
    // 点分树的子树中的答案(不经过重心)
    removed[cur] = 1
    const _nexts = centTree[cur]
    for (let i = 0; i < _nexts.length; i++) {
      const next = _nexts[i]
      if (!removed[next]) {
        decomposition(next, cur)
      }
    }
    removed[cur] = 0

    const counter = new Map<number, number>([[vals[cur], 1]]) // !经过重心的路径
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next === pre || removed[next]) continue

      const sub = new Map<number, number>() // value -> count
      collect(next, cur, vals[cur], sub)
      sub.forEach((count, value) => {
        res += count * (counter.get(value) || 0)
        counter.set(value, (counter.get(value) || 0) + count)
      })
    }
  }

  /**
   * 收集子树信息.
   */
  function collect(cur: number, pre: number, max: number, sub: Map<number, number>): void {
    const val = vals[cur]
    if (val >= max) {
      sub.set(val, (sub.get(val) || 0) + 1)
      max = val
    }
    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next !== pre && !removed[next]) {
        collect(next, cur, max, sub)
      }
    }
  }
}
