/* eslint-disable @typescript-eslint/no-non-null-assertion */

import { Heap } from '../../../8_heap/Heap'

/**
 * 双向dijkstra求start到end的最短路.
 * 如果不存在，返回-1.
 * @param rg 反向图
 */
function fastDijkstra(
  n: number,
  g: [next: number, weight: number][][],
  rg: [next: number, weight: number][][],
  start: number,
  end: number
): number {
  if (start === end) return 0
  const dist = Array(n)
  const drev = Array(n)
  for (let i = 0; i < n; ++i) {
    dist[i] = -1
    drev[i] = -1
  }
  dist[start] = 0
  drev[end] = 0

  const pq1 = new Heap<[dist: number, vertex: number]>({ data: [], less: (a, b) => a[0] < b[0] })
  const pq2 = new Heap<[dist: number, vertex: number]>({ data: [], less: (a, b) => a[0] < b[0] })
  pq1.push([0, start])
  pq2.push([0, end])

  let res = -1
  while (pq1.size && pq2.size) {
    const d1 = pq1.top()![0]
    const d2 = pq2.top()![0]
    if (res >= 0 && d1 + d2 >= res) break
    if (d1 <= d2) {
      const [d, u] = pq1.pop()!
      if (dist[u] > d) continue
      // eslint-disable-next-line no-loop-func
      g[u].forEach(([v, w]) => {
        const cand = dist[u] + w
        if (dist[v] >= 0 && dist[v] <= cand) return
        dist[v] = cand
        pq1.push([dist[v], v])
        if (drev[v] >= 0) {
          const nu = dist[v] + drev[v]
          if (res < 0 || res > nu) res = nu
        }
      })
    } else {
      const [d, u] = pq2.pop()!
      if (drev[u] > d) continue
      // eslint-disable-next-line no-loop-func
      rg[u].forEach(([v, w]) => {
        const cand = drev[u] + w
        if (drev[v] >= 0 && drev[v] <= cand) return
        drev[v] = cand
        pq2.push([drev[v], v])
        if (dist[v] >= 0) {
          const nu = dist[v] + drev[v]
          if (res < 0 || res > nu) res = nu
        }
      })
    }
  }

  return res
}

export { fastDijkstra }

if (require.main === module) {
  // https://leetcode.cn/problems/design-graph-with-shortest-path-calculator/
  class Graph {
    private readonly _g: [next: number, weight: number][][]
    private readonly _rg: [next: number, weight: number][][]
    constructor(n: number, edges: number[][]) {
      const adjList: [number, number][][] = Array(n)
      for (let i = 0; i < n; i++) adjList[i] = []
      const radjList: [number, number][][] = Array(n)
      for (let i = 0; i < n; i++) radjList[i] = []
      edges.forEach(([u, v, w]) => {
        adjList[u].push([v, w])
        radjList[v].push([u, w])
      })
      this._g = adjList
      this._rg = radjList
    }

    addEdge(edge: number[]): void {
      const [u, v, w] = edge
      this._g[u].push([v, w])
      this._rg[v].push([u, w])
    }

    shortestPath(node1: number, node2: number): number {
      return fastDijkstra(this._g.length, this._g, this._rg, node1, node2)
    }
  }
}
