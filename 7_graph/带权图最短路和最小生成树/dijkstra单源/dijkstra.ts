import { Heap } from '../../../8_heap/Heap'

const INF = 2e15
function dijkstra(n: number, graph: [v: number, w: number][][], start: number): number[] {
  const dist = Array(n)
  for (let i = 0; i < n; ++i) dist[i] = INF
  const queue = new Heap<[cur: number, curDist: number]>((a, b) => a[0] - b[0])
  dist[start] = 0
  queue.push([0, start])

  while (queue.size) {
    const [curDist, cur] = queue.pop()!
    if (curDist > dist[cur]) continue
    graph[cur].forEach(([v, w]) => {
      const nextDist = curDist + w
      if (dist[v] > nextDist) {
        dist[v] = nextDist
        queue.push([dist[v], v])
      }
    })
  }

  return dist
}

export { dijkstra }

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  function networkDelayTime(times: number[][], n: number, k: number): number {
    const adjList: [v: number, w: number][][] = Array(n)
    for (let i = 0; i < n; i++) adjList[i] = []
    times.forEach(([u, v, w]) => {
      adjList[u - 1].push([v - 1, w])
    })
    const dist = dijkstra(n, adjList, k - 1)
    return Math.max(...dist) === INF ? -1 : Math.max(...dist)
  }
}
