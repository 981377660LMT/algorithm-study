/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * Dinic求最大流 时间复杂度：O(V^2*E)
 *
 * @param n 图的顶点个数
 * @param start (虚拟)源点
 * @param end (虚拟)汇点
 */
function useDinic(n: number, start: number, end: number) {
  if (start < 0 || start >= n || end < 0 || end >= n) {
    throw new RangeError(`start: ${start}, end: ${end} out of range [0, ${n - 1}]`)
  }

  const reGraph = new Map<number, Map<number, number>>() // 残量图

  /**
   * 添加边 {@link from} -> {@link to}, 容量为 {@link capacity}
   *
   * {@link cover} 表示更新容量是否覆盖原有的边
   */
  function addEdge(from: number, to: number, capacity: number, cover = false): void {
    !reGraph.has(from) && reGraph.set(from, new Map())
    let innerMap = reGraph.get(from)!

    if (cover) {
      innerMap.set(to, capacity) // 覆盖
    } else {
      innerMap.set(to, (innerMap.get(to) ?? 0) + capacity) // 增加
    }

    !reGraph.has(to) && reGraph.set(to, new Map())
    innerMap = reGraph.get(to)!
    if (!innerMap.has(from)) innerMap.set(from, 0) // 防止自环边影响
  }

  function calMaxFlow(): number {
    let res = 0
    const depth = new Int32Array(n)
    let iterMap: Map<number, IterableIterator<number>> // 当前弧优化

    while (bfs()) {
      iterMap = makeCurEdge(reGraph)
      while (true) {
        const delta = dfs(start, end, Infinity)
        if (delta === 0) break
        res += delta
      }
    }

    return res

    function bfs(): boolean {
      depth.fill(-1)
      depth[start] = 0

      let queue: number[] = [start]

      while (queue.length) {
        const nextQueue: number[] = []
        const steps = queue.length
        for (let _ = 0; _ < steps; _++) {
          const cur = queue.pop()!
          const nextDist = depth[cur] + 1
          for (const [next, remain] of reGraph.get(cur) ?? []) {
            if (depth[next] === -1 && remain > 0) {
              depth[next] = nextDist
              nextQueue.push(next)
            }
          }
        }

        queue = nextQueue
      }

      return depth[end] !== -1
    }

    /**
     * @param cur 当前点
     * @param minFlow 路径上的最小流量
     * @returns 增广路径上的最小流量
     */
    function dfs(cur: number, end: number, minFlow: number): number {
      if (cur === end) return minFlow

      const innerMap1 = reGraph.get(cur)!
      for (const next of iterMap.get(cur) ?? []) {
        const innerMap2 = reGraph.get(next)!

        const remain = innerMap1.get(next)!
        if (remain && depth[cur] < depth[next]) {
          const nextFlow = dfs(next, end, Math.min(minFlow, remain))
          if (nextFlow) {
            innerMap1.set(next, innerMap1.get(next)! - nextFlow)
            innerMap2.set(cur, innerMap2.get(cur)! + nextFlow)
            return nextFlow
          }
        }
      }

      return 0
    }

    function makeCurEdge(
      reGraph: Map<number, Map<number, number>>
    ): Map<number, IterableIterator<number>> {
      const res = new Map<number, IterableIterator<number>>()
      for (const key of reGraph.keys()) res.set(key, reGraph.get(key)!.keys())
      return res
    }
  }

  /**
   * @returns 最大流经过的点
   */
  function getPath(): Set<number> {
    const visited = new Set<number>()
    const stack = [start]
    while (stack.length) {
      const cur = stack.pop()!
      visited.add(cur)
      for (const [next, remain] of reGraph.get(cur) ?? []) {
        if (remain && !visited.has(next)) {
          visited.add(next)
          stack.push(next)
        }
      }
    }

    return visited
  }

  /**
   * @returns 边的残量(剩余的容量)
   */
  function getRemainOfEdge(v1: number, v2: number): number {
    const innerMap = reGraph.get(v1)
    return innerMap?.get(v2) ?? 0
  }

  return {
    addEdge,
    calMaxFlow,
    getPath,
    getRemainOfEdge
  }
}

export { useDinic }
