type Edge = [next: number, remainCapacity: number, rEdge: Edge]

function* iter<T>(iterable: Iterable<T>) {
  yield* iterable
}

/**
 * Dinic算法求最大流
 *
 * 时间复杂度: `O(V^2*E)`
 */
function useDinic(n: number, start: number, end: number) {
  const graph = Array.from<unknown, Edge[]>({ length: n }, () => [])
  let levels: Int32Array
  let edgeIters: IterableIterator<Edge>[]

  /**
   * 添加边 {@link from} -> {@link to}, 容量为 {@link capacity}
   */
  function addEdge(from: number, to: number, capacity: number): void {
    const forward: Edge = [to, capacity, null as any]
    const backward: Edge = [from, 0, forward]
    forward[2] = backward
    graph[from].push(forward)
    graph[to].push(backward)
  }

  /**
   * 针对重边的情况 需要将重边处理后一起添加
   *
   * 添加边 {@link v1} -> {@link v2}, 容量为 {@link w1}
   * 添加边 {@link v2} -> {@link v1}, 容量为 {@link w2}
   */
  function addMultiEdge(v1: number, v2: number, w1: number, w2: number): void {
    const edge1: Edge = [v2, w1, null as any]
    const edge2: Edge = [v1, w2, edge1]
    edge1[2] = edge2
    graph[v1].push(edge1)
    graph[v2].push(edge2)
  }

  function calMaxFlow(): number {
    let flow = 0
    while (bfs()) {
      edgeIters = graph.map(iter) // 当前弧优化
      let delta = Infinity
      while (delta) {
        delta = dfs(start, end, Infinity)
        flow += delta
      }
    }

    return flow
  }

  /**
   * @returns 所有边的残量(剩余的容量)
   */
  function getEdgeRemain(): Map<number, Map<number, number>> {
    const res = new Map<number, Map<number, number>>()
    for (const edges of graph) {
      for (const edge of edges) {
        const [pre, next, remain] = [edge[2][0], edge[0], edge[1]]
        if (!res.has(pre)) {
          res.set(pre, new Map())
        }
        res.get(pre)?.set(next, remain)
      }
    }

    return res
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
      for (const [next, remain] of graph[cur]) {
        if (remain && !visited.has(next)) {
          visited.add(next)
          stack.push(next)
        }
      }
    }

    return visited
  }

  /**
   * 建立分层图
   */
  function bfs(): boolean {
    levels = new Int32Array(n).fill(-1)
    let queue = [start]
    levels[start] = 0

    while (queue.length) {
      const nextQueue: number[] = []
      const len = queue.length
      for (let _ = 0; _ < len; _++) {
        const cur = queue.pop()!
        const dist = levels[cur] + 1
        for (const [next, remain] of graph[cur]) {
          if (levels[next] === -1 && remain > 0) {
            levels[next] = dist
            nextQueue.push(next)
          }
        }
      }

      queue = nextQueue
    }

    return levels[end] !== -1
  }

  /**
   * 寻找增广路
   */
  function dfs(cur: number, target: number, flow: number): number {
    if (cur === target) {
      return flow
    }

    for (const edge of edgeIters[cur]) {
      const [next, remain, rEdge] = edge
      if (remain && levels[cur] < levels[next]) {
        const delta = dfs(next, target, Math.min(flow, remain))
        if (delta) {
          edge[1] -= delta
          rEdge[1] += delta
          return delta
        }
      }
    }

    return 0
  }

  return {
    addEdge,
    addMultiEdge,
    calMaxFlow,
    getEdgeRemain,
    getPath
  }
}

export { useDinic }
