/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-shadow */
// !Dinic算法 数组+边存图 速度快

const INF = 2e15 // !2**53 -1 约为 9e15

/**
 * Dinic算法求最大流
 *
 * 时间复杂度: `O(V^2*E)` , 二分图上为 O(E*√V)
 *
 * @param n 图的顶点个数
 * @param start (虚拟)源点
 * @param end (虚拟)汇点
 */
function useDinic(n: number, start: number, end: number) {
  if (start < 0 || start >= n || end < 0 || end >= n) {
    throw new RangeError(`start: ${start}, end: ${end} out of range [0, ${n - 1}]`)
  }

  const _reGraph = Array.from<unknown, number[]>({ length: n }, () => [])
  const _edges: [next: number, capacity: number][] = []
  const _visitedEdge = new Set<number>()

  /**
   * 添加边 {@link from} -> {@link to}, 容量为 {@link capacity}
   * 注意这种方式会添加重边
   */
  function addEdge(from: number, to: number, capacity: number): void {
    const hash = from * n + to
    _visitedEdge.add(hash)
    _reGraph[from].push(_edges.length)
    _edges.push([to, capacity])
    _reGraph[to].push(_edges.length)
    _edges.push([from, 0])
  }

  /**
   * 如果没有添加过这条边，
   * 则添加边 {@link from} -> {@link to}, 容量为 {@link capacity}
   */
  function addEdgeIfAbsent(from: number, to: number, capacity: number): void {
    const hash = from * n + to
    if (_visitedEdge.has(hash)) return
    _visitedEdge.add(hash)
    _reGraph[from].push(_edges.length)
    _edges.push([to, capacity])
    _reGraph[to].push(_edges.length)
    _edges.push([from, 0])
  }

  function calMaxFlow(): number {
    const levels = new Int32Array(n)
    const curEdges = new Int32Array(n) // 当前弧优化

    let res = 0
    while (bfs(start, end)) {
      curEdges.fill(0)
      res += dfs(start, end, INF)
    }
    return res

    /**
     * 建立分层图
     */
    function bfs(start: number, end: number): boolean {
      let queue = [start]
      levels.fill(-1)
      levels[start] = 0

      while (queue.length) {
        const nextQueue: number[] = []
        const step = queue.length

        for (let _ = 0; _ < step; _++) {
          const cur = queue.pop()!

          // !不要使用 for of 来遍历迭代器循环 速度会变慢
          for (let i = 0; i < _reGraph[cur].length; i++) {
            const ei = _reGraph[cur][i]
            const next = _edges[ei][0] // !不要使用 const [next,capacity] = edges[ei] 解构 速度会变慢
            const capacity = _edges[ei][1]
            if (capacity > 0 && levels[next] === -1) {
              levels[next] = levels[cur] + 1
              if (next === end) return true
              nextQueue.push(next)
            }
          }
        }

        queue = nextQueue
      }

      return false
    }

    /**
     * 寻找增广路
     */
    function dfs(cur: number, end: number, flow: number): number {
      if (cur === end) {
        return flow
      }

      let res = flow
      // 当前弧优化
      for (let ei = curEdges[cur]; ei < _reGraph[cur].length; ei = ++curEdges[cur]) {
        const ej = _reGraph[cur][ei]
        const next = _edges[ej][0]
        const remain = _edges[ej][1]
        if (remain > 0 && levels[next] === levels[cur] + 1) {
          const delta = dfs(next, end, Math.min(res, remain))
          _edges[ej][1] -= delta
          _edges[ej ^ 1][1] += delta
          res -= delta
          if (res === 0) return flow
        }
      }

      return flow - res
    }
  }

  /**
   * @returns 边的残量(剩余的容量)
   * @example
   * ```typescript
   * const maxFlow = useDinic(n, start, end)
   * const query = maxFlow.useQueryRemainOfEdge()
   * console.log(query(0, 1))
   * ```
   */
  function useQueryRemainOfEdge(): (v1: number, v2: number) => number {
    const adjMap = Array.from<number, Map<number, number>>({ length: n }, () => new Map())
    for (let cur = 0; cur < n; cur++) {
      const eis = _reGraph[cur]
      const innerMap = adjMap[cur]

      for (let i = 0; i < eis.length; i++) {
        const ei = eis[i]
        const edge = _edges[ei]
        const next = edge[0]
        const remain = edge[1]

        innerMap.set(next, (innerMap.get(next) || 0) + remain)
      }
    }

    return queryApi

    function queryApi(v1: number, v2: number): number {
      const innerMap = adjMap[v1]
      return innerMap.get(v2) || 0
    }
  }

  /**
   * @returns 最大流经过的点
   */
  function getPath(): Set<number> {
    const visited = new Set<number>()
    const queue = [start]
    while (queue.length) {
      const cur = queue.pop()!
      visited.add(cur)
      for (let i = 0; i < _reGraph[cur].length; i++) {
        const ei = _reGraph[cur][i]
        const edge = _edges[ei]
        const next = edge[0]
        const remain = edge[1]
        if (remain > 0 && !visited.has(next)) {
          visited.add(next)
          queue.push(next)
        }
      }
    }

    return visited
  }

  return {
    addEdge,
    addEdgeIfAbsent,
    calMaxFlow,
    useQueryRemainOfEdge,
    getPath
  }
}

export { useDinic }
