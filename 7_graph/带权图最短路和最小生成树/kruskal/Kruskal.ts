/* eslint-disable consistent-return */
/* eslint-disable no-param-reassign */

import { UnionFindArray } from '../../../14_并查集/UnionFind'

/**
 * Kruskal算法求无向图最小生成树(森林).
 */
function kruskal(
  n: number,
  edges: [u: number, v: number, w: number][] | number[][]
): {
  /** 最小生成树(森林)的权值之和. */
  mstCost: number
  /** 每条边是否在最小生成树(森林)中. */
  inMst: Uint8Array
  /** 是否是树/连通. */
  isTree: boolean
} {
  const uf = new UnionFindArray(n)
  let count = 0
  let mstCost = 0
  const inMst = new Uint8Array(edges.length)
  let isTree = false

  const order = new Uint32Array(edges.length)
  for (let i = 0; i < edges.length; i++) order[i] = i
  order.sort((a, b) => edges[a][2] - edges[b][2])

  for (let i = 0; i < edges.length; i++) {
    const ei = order[i]
    const { 0: u, 1: v, 2: w } = edges[ei]
    if (uf.union(u, v)) {
      inMst[ei] = 1
      mstCost += w
      count++
      if (count === n - 1) {
        isTree = true
        break
      }
    }
  }

  return { mstCost, inMst, isTree }
}

export { kruskal }

if (require.main === module) {
  const n = 4
  const edges = [
    [0, 1, 1],
    [0, 2, 2],
    [0, 3, 3],
    [1, 2, 4],
    [1, 3, 5],
    [2, 3, 6]
  ]
  console.log(kruskal(n, edges))
  // { mstCost: 6, inMst: Uint8Array [ 1, 1, 1, 0, 0, 0 ], isTree: true }
}
