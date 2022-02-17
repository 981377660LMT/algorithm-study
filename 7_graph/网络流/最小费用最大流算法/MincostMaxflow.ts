import { MinHeap } from '../../../8_heap/MinHeap'

class Edge {
  constructor(
    public readonly from: number,
    public readonly to: number,
    public readonly cost: number,
    public cap: number
  ) {}
}

// O(nm(n+m))
function mincostMaxflow(n: number) {
  const adjList = Array.from<unknown, Edge[]>({ length: n }, () => [])
  const dist = Array<number>(n).fill(0)
  const preVertical = Array<number>(n).fill(0)
  const preEdge = Array<number>(n).fill(0)
  const cost = Array<number>(n).fill(0)

  function addEdge(from: number, to: number, cost: number, cap: number): void {
    adjList[from].push(new Edge(adjList[to].length, to, cost, cap))
    adjList[to].push(new Edge(adjList[from].length - 1, from, -cost, 0))
  }

  /**
   *
   * @param from
   * @param to
   * @param flow 最大流大小
   * @returns
   */
  function mincostFlow(from: number, to: number, flow: number): number {
    let res = 0

    while (flow > 0) {
      dist.fill(Infinity)
      dist[from] = 0
      const pq = new MinHeap<[dist: number, cur: number]>((a, b) => a[0] - b[0] || a[1] - b[1])
      pq.heappush([0, from])

      while (pq.size > 0) {
        const [curDis, cur] = pq.heappop()!
        if (dist[cur] < curDis) continue
        for (let i = 0; i < adjList[cur].length; i++) {
          const nextEdge = adjList[cur][i]
          if (
            nextEdge.cap > 0 &&
            dist[nextEdge.to] > dist[cur] + nextEdge.cost + cost[cur] - cost[nextEdge.to]
          ) {
            dist[nextEdge.to] = dist[cur] + nextEdge.cost + cost[cur] - cost[nextEdge.to]
            preVertical[nextEdge.to] = cur
            preEdge[nextEdge.to] = i
            pq.heappush([dist[nextEdge.to], nextEdge.to])
          }
        }
      }

      if (dist[to] === Infinity) return -1
      for (let i = 0; i < n; i++) {
        cost[i] += dist[i]
      }

      let minflow = flow
      for (let i = to; i !== from; i = preVertical[i]) {
        minflow = Math.min(minflow, adjList[preVertical[i]][preEdge[i]].cap)
      }

      flow -= minflow
      res += minflow * cost[to]

      for (let i = to; i !== from; i = preVertical[i]) {
        const edge = adjList[preVertical[i]][preEdge[i]]
        edge.cap -= minflow
        adjList[i][edge.from].cap += minflow
      }
    }

    return res
  }

  return {
    addEdge,
    mincostFlow,
  }
}

export { mincostMaxflow }
