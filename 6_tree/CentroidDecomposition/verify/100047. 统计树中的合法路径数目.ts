/* eslint-disable no-loop-func */
// https://leetcode.cn/problems/count-valid-paths-in-a-tree/description/
// 100047. 统计树中的合法路径数目
//
// 给你一棵 n 个节点的无向树，节点编号为 1 到 n 。给你一个整数 n 和一个长度为 n - 1 的二维整数数组 edges ，
// 其中 edges[i] = [ui, vi] 表示节点 ui 和 vi 在树中有一条边。
//
// 请你返回树中的 合法路径数目 。
//
// 如果在节点 a 到节点 b 之间 恰好有一个 节点的编号是质数，那么我们称路径 (a, b) 是 合法的 。
//
// 注意：
//
// 路径 (a, b) 指的是一条从节点 a 开始到节点 b 结束的一个节点序列，序列中的节点 互不相同 ，且相邻节点之间在树上有一条边。
// 路径 (a, b) 和路径 (b, a) 视为 同一条 路径，且只计入答案 一次 。

import { EratosthenesSieve } from '../../../19_数学/因数筛/prime'
import { centroidDecomposition } from '../CentroidDecomposition'

const E = new EratosthenesSieve(1e5 + 10)

function countPaths(n: number, edges: number[][]): number {
  const tree: [next: number, weight: number][][] = Array(n)
  for (let i = 0; i < n; i++) tree[i] = []
  edges.forEach(([u, v]) => {
    tree[u - 1].push([v - 1, 1])
    tree[v - 1].push([u - 1, 1])
  })

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
    const _nexts = centTree[cur]
    for (let i = 0; i < _nexts.length; i++) {
      const next = _nexts[i]
      if (!removed[next]) {
        decomposition(next, cur)
      }
    }
    removed[cur] = 0

    // !经过重心的路径
    const counter = new Map<number, number>()
    const init = +E.isPrime(cur + 1)
    counter.set(init, 1)

    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next === pre || removed[next]) continue

      const sub = new Map<number, number>() // 统计子树内(不含cur)
      collect(next, cur, init, sub)
      if (init === 0) {
        res += (sub.get(0) || 0) * (counter.get(1) || 0)
        res += (sub.get(1) || 0) * (counter.get(0) || 0)
        counter.set(0, (counter.get(0) || 0) + (sub.get(0) || 0))
        counter.set(1, (counter.get(1) || 0) + (sub.get(1) || 0))
      } else {
        res += (sub.get(1) || 0) * (counter.get(1) || 0)
        counter.set(1, (counter.get(1) || 0) + (sub.get(1) || 0))
      }
    }
  }

  /**
   * 收集子树信息.
   */
  function collect(cur: number, pre: number, primeCount: number, sub: Map<number, number>): void {
    if (E.isPrime(cur + 1)) primeCount++
    sub.set(primeCount, (sub.get(primeCount) || 0) + 1)
    if (primeCount > 1) return

    const nexts = tree[cur]
    for (let i = 0; i < nexts.length; i++) {
      const next = nexts[i][0]
      if (next !== pre && !removed[next]) {
        collect(next, cur, primeCount, sub)
      }
    }
  }
}
