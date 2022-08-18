/* eslint-disable no-shadow */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable no-param-reassign */
/* eslint-disable no-constant-condition */
/* eslint-disable no-useless-constructor */

const INF = 2e15

class Edge {
  constructor(
    public readonly from: number,
    public readonly to: number,
    public readonly capacity: number,
    public readonly cost: number,
    public flow: number
  ) {}
}

/**
 * @param n 顶点个数
 * @param start (虚拟)源点
 * @param end (虚拟)汇点
 */
function useMinCostMaxFlow(n: number, start: number, end: number) {
  n += 10
  const edges: Edge[] = []
  const reGraph: number[][] = Array.from({ length: n }, () => []) // 残量图存储的是边的下标

  function addEdge(from: number, to: number, capacity: number, cost: number): void {
    // 原边索引为i 反向边索引为i^1
    edges.push(new Edge(from, to, capacity, cost, 0))
    edges.push(new Edge(to, from, 0, -cost, 0))
    const len = edges.length
    reGraph[from].push(len - 2)
    reGraph[to].push(len - 1)
  }

  function work(): [maxFlow: number, minCost: number] {
    const dist = Array<number>(n).fill(INF)
    let [flow, cost] = [0, 0]
    while (true) {
      const delta = spfa()
      if (delta === 0) break
      flow += delta
      cost += delta * dist[end]
    }

    return [flow, cost]

    // spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra
    function spfa(): number {
      dist.fill(INF)
      dist[start] = 0
      const inQueue = new Uint8Array(n)
      let queue = [start]

      const inFlow = Array<number>(n).fill(0)
      inFlow[start] = INF
      const pre = new Int32Array(n).fill(-1)

      while (queue.length) {
        const nextQueue: number[] = []
        const steps = queue.length
        for (let _ = 0; _ < steps; _++) {
          const cur = queue.pop()!
          inQueue[cur] = 0
          for (const edgeIndex of reGraph[cur]) {
            const { capacity, cost, flow, to: next } = edges[edgeIndex]
            if (flow < capacity && dist[next] > dist[cur] + cost) {
              dist[next] = dist[cur] + cost
              inFlow[next] = Math.min(capacity - flow, inFlow[cur])
              pre[next] = edgeIndex
              if (!inQueue[next]) {
                inQueue[next] = 1
                nextQueue.push(next)
              }
            }
          }
        }

        queue = nextQueue
      }

      const resDelta = inFlow[end]
      if (resDelta > 0) {
        let cur = end
        while (cur !== start) {
          const edgeIndex = pre[cur]
          edges[edgeIndex].flow += resDelta
          edges[edgeIndex ^ 1].flow -= resDelta
          cur = edges[edgeIndex].from
        }
      }

      return resDelta
    }
  }

  return {
    addEdge,
    work
  }
}

export { useMinCostMaxFlow }
