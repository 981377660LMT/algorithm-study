/* eslint-disable no-param-reassign */
import { useUnionFindMap } from '../../../14_并查集/useUnionFind'

type KruskalEdge<V = unknown> = [u: V, v: V, weight: number]

/**
 * @param n 顶点数
 * @param edges 无向边
 */
function kruskal<V = unknown>(n: number, edges: KruskalEdge<V>[]): number {
  const uf = useUnionFindMap<V>()
  let res = 0
  let hit = 0

  edges = edges.slice().sort((a, b) => a[2] - b[2])
  for (const [u, v, w] of edges) {
    const [root1, root2] = [uf.find(u), uf.find(v)]
    if (root1 !== root2) {
      res += w
      hit++
      if (hit === n - 1) break
      uf.union(root1, root2)
    }
  }

  return hit === n - 1 ? res : Infinity
}

export { kruskal, KruskalEdge }
