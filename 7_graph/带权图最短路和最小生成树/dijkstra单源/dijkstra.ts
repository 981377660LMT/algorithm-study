import { Heap } from '../../../8_heap/Heap'
import { HeapUint32 } from '../../../8_heap/SiftHeap'

const INF = 2e15

/**
 * 采用SiftHeap加速的dijkstra算法.求出起点到各点的最短距离.
 * 1e6 => 60ms
 */
function dijkstraSiftHeap1(n: number, graph: [v: number, w: number][][], start: number): number[] {
  const dist = Array(n)
  for (let i = 0; i < n; ++i) dist[i] = INF
  const queue = new HeapUint32(n, (i, j) => dist[i] < dist[j])
  dist[start] = 0
  queue.push(start)

  while (queue.size) {
    const cur = queue.pop()!
    const nexts = graph[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const [next, weight] = nexts[i]
      const cand = dist[cur] + weight
      if (cand < dist[next]) {
        dist[next] = cand
        queue.push(next)
      }
    }
  }

  return dist
}

/**
 * 采用SiftHeap加速的dijkstra算法.求出一条路径.
 */
function dijkstraSiftHeap2(
  n: number,
  adjList: [v: number, w: number][][],
  start: number,
  end: number
): [dist: number, path: number[]] {
  const dist = Array(n)
  const pre = new Int32Array(n)
  for (let i = 0; i < n; ++i) {
    dist[i] = INF
    pre[i] = -1
  }
  dist[start] = 0
  const pq = new HeapUint32(n, (i, j) => dist[i] < dist[j])
  pq.push(start)

  while (pq.size) {
    const cur = pq.pop()!
    const nexts = adjList[cur]
    for (let i = 0; i < nexts.length; ++i) {
      const [next, weight] = nexts[i]
      const cand = dist[cur] + weight
      if (cand < dist[next]) {
        dist[next] = cand
        pq.push(next)
        pre[next] = cur
      }
    }
  }

  if (dist[end] === INF) return [INF, []]

  const path: number[] = []
  let cur = end
  while (~pre[cur]) {
    path.push(cur)
    cur = pre[cur]
  }
  path.push(start)
  for (let i = 0, j = path.length - 1; i < j; ++i, --j) {
    const tmp = path[i]
    path[i] = path[j]
    path[j] = tmp
  }
  return [dist[end], path]
}

/**
 * dijkstra求出起点到各点的最短距离,时间复杂度O((V+E)logV).
 */
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

/**
 * dijkstra求出一条路径.
 */
function dijkstra2(
  n: number,
  adjList: [v: number, w: number][][],
  start: number,
  end: number
): [dist: number, path: number[]] {
  const dist = Array(n)
  const pre = new Int32Array(n)
  for (let i = 0; i < n; ++i) {
    dist[i] = INF
    pre[i] = -1
  }
  dist[start] = 0
  const pq = new Heap<[curDist: number, cur: number]>((a, b) => a[0] - b[0])
  pq.push([0, start])

  while (pq.size) {
    const [curDist, cur] = pq.pop()!
    if (dist[cur] < curDist) continue
    adjList[cur].forEach(([next, weight]) => {
      const cand = dist[cur] + weight
      if (cand < dist[next]) {
        dist[next] = cand
        pq.push([dist[next], next])
        pre[next] = cur
      }
    })
  }

  if (dist[end] === INF) return [INF, []]

  const path: number[] = []
  let cur = end
  while (~pre[cur]) {
    path.push(cur)
    cur = pre[cur]
  }
  path.push(start)
  return [dist[end], path.reverse()]
}

export { dijkstra, dijkstra2, dijkstraSiftHeap1, dijkstraSiftHeap2 }

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

  const n = 1e6
  const edges: [u: number, v: number, w: number][] = []
  for (let i = 0; i < n; ++i) {
    edges.push([i, i + 1, 1])
  }
  const adjList: [v: number, w: number][][] = Array(n + 1)
  for (let i = 0; i <= n; ++i) adjList[i] = []
  edges.forEach(([u, v, w]) => {
    adjList[u].push([v, w])
  })
  console.time('dijkstra')
  dijkstra(n + 1, adjList, 0)
  console.timeEnd('dijkstra')

  console.time('dijkstraSiftHeap')
  dijkstraSiftHeap1(n + 1, adjList, 0)
  console.timeEnd('dijkstraSiftHeap')
}
