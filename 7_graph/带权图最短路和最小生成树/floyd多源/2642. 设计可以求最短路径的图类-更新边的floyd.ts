// https://leetcode.cn/problems/design-graph-with-shortest-path-calculator/
// 1 <= n <= 100
// 0 <= edges.length <= n * (n - 1)
// !调用 shortestPath 1e5 次

const INF = 2e15

class FloydDynamic {
  private readonly _n: number
  private readonly _dist: number[]

  constructor(n: number, edges: [u: number, v: number, w: number][]) {
    const dist = Array(n * n)
    for (let i = 0; i < n * n; ++i) dist[i] = INF
    for (let i = 0; i < n; ++i) dist[i * n + i] = 0
    edges.forEach(([u, v, w]) => {
      dist[u * n + v] = Math.min(dist[u * n + v], w)
    })

    for (let k = 0; k < n; ++k) {
      for (let i = 0; i < n; ++i) {
        for (let j = 0; j < n; ++j) {
          const cand = dist[i * n + k] + dist[k * n + j]
          if (dist[i * n + j] > cand) {
            dist[i * n + j] = cand
          }
        }
      }
    }

    this._n = n
    this._dist = dist
  }

  /**
   * 向边集中添加一条有向边.
   * 加边时,枚举每个点对,根据是否经过edge来更新最短路.
   */
  addEdge(edge: [u: number, v: number, w: number]): void {
    const [u, v, w] = edge
    const n = this._n
    for (let i = 0; i < n; ++i) {
      for (let j = 0; j < n; ++j) {
        const cand = this._dist[i * n + u] + w + this._dist[v * n + j]
        if (this._dist[i * n + j] > cand) {
          this._dist[i * n + j] = cand
        }
      }
    }
  }

  /**
   * 求出从`start`到`target`的最短路径长度,如果不存在这样的路径,返回-1.
   */
  shortestPath(start: number, target: number): number {
    const dist = this._dist[start * this._n + target]
    return dist === INF ? -1 : dist
  }
}

const Graph = FloydDynamic
/**
 * Your Graph object will be instantiated and called as such:
 * var obj = new Graph(n, edges)
 * obj.addEdge(edge)
 * var param_2 = obj.shortestPath(node1,node2)
 */

export {}
