/* eslint-disable no-param-reassign */
import { useUnionFindArray } from '../../../14_并查集/useUnionFind'

type KruskalEdge = [from: number, to: number, weight: number]

/**
 * @param n 顶点数
 * @param edges 无向边
 * @returns 最小生成树的权值和，如果不存在则返回 -1
 */
function kruskal(n: number, edges: KruskalEdge[]): number {
  const uf = useUnionFindArray(n)
  let res = 0
  let count = 0

  edges = edges.slice().sort((a, b) => a[2] - b[2])
  for (const [u, v, w] of edges) {
    const [root1, root2] = [uf.find(u), uf.find(v)]
    if (root1 !== root2) {
      uf.union(root1, root2)
      res += w
      count++
      if (count === n - 1) return res
    }
  }

  return -1
}

export { kruskal }
