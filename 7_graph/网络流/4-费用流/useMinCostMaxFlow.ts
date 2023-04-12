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
  const _edges: Edge[] = []
  const _reGraph = Array(n)
  for (let i = 0; i < n; i++) {
    _reGraph[i] = []
  }

  const _dist = Array<number>(n).fill(INF)
  const _flow = Array<number>(n).fill(0)
  const _pre = new Int16Array(n).fill(-1)

  function addEdge(from: number, to: number, capacity: number, cost: number): void {
    // 原边索引为i 反向边索引为i^1
    _edges.push(new Edge(from, to, capacity, cost, 0))
    _edges.push(new Edge(to, from, 0, -cost, 0))
    const len = _edges.length
    _reGraph[from].push(len - 2)
    _reGraph[to].push(len - 1)
  }

  function work(): [maxFlow: number, minCost: number] {
    let [flow, cost] = [0, 0]
    while (_spfa()) {
      const delta = _flow[end]
      flow += delta
      cost += delta * _dist[end]
      let cur = end
      while (cur !== start) {
        const edgeIndex = _pre[cur]
        _edges[edgeIndex].flow += delta
        _edges[edgeIndex ^ 1].flow -= delta
        cur = _edges[edgeIndex].from
      }
    }

    return [flow, cost]
  }

  return {
    addEdge,
    work
  }

  // spfa沿着最短路寻找增广路径  有负cost的边不能用dijkstra
  function _spfa(): boolean {
    _dist.fill(INF)
    _dist[start] = 0
    const inQueue = new Uint8Array(n)
    inQueue[start] = 1
    let queue = [start]

    _flow.fill(0)
    _flow[start] = INF
    _pre.fill(-1)

    while (queue.length) {
      const nextQueue: number[] = []
      const steps = queue.length
      for (let _ = 0; _ < steps; _++) {
        const cur = queue.pop()!
        inQueue[cur] = 0
        for (const edgeIndex of _reGraph[cur]) {
          const { capacity, cost, flow, to: next } = _edges[edgeIndex]
          if (flow < capacity && _dist[next] > _dist[cur] + cost) {
            _dist[next] = _dist[cur] + cost
            _flow[next] = Math.min(capacity - flow, _flow[cur])
            _pre[next] = edgeIndex
            if (!inQueue[next]) {
              inQueue[next] = 1
              nextQueue.push(next)
            }
          }
        }
      }

      queue = nextQueue
    }

    return _pre[end] !== -1
  }
}

export { useMinCostMaxFlow }
