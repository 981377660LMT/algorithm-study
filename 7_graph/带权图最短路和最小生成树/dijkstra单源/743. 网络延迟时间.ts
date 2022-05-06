import { PriorityQueue } from '../../../2_queue/优先级队列'

type Cur = number
type Weight = number
type Edge = [Cur, Weight]

/**
 * @param {number[][]} times
 * @param {number} n  有 n 个网络节点，标记为 1 到 n。
 * @param {number} k
 * @return {number}
 * @description O((V+E)logV)
 */
const networkDelayTime = function (times: number[][], n: number, k: number): number {
  const adjList = Array.from<number, [number, number][]>({ length: n + 1 }, () => [])
  times.forEach(([u, v, w]) => adjList[u].push([v, w]))

  const dist = Array<number>(n + 1).fill(Infinity)
  dist[k] = 0

  const visited = new Set<number>()

  const compareFunction = (a: Edge, b: Edge) => a[1] - b[1]
  const priorityQueue = new PriorityQueue<Edge>(compareFunction)
  priorityQueue.push([k, 0])

  while (priorityQueue.length) {
    // 1.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
    const [cur, _] = priorityQueue.shift()!
    if (visited.has(cur)) continue

    // 2.加入visited
    visited.add(cur)

    // 3.利用cur点来更新其相邻节点next与原点的距离
    for (const [next, weight] of adjList[cur]) {
      if (visited.has(next)) continue
      if (dist[cur] + weight < dist[next]) {
        dist[next] = dist[cur] + weight
        priorityQueue.push([next, dist[next]])
      }
    }
  }

  let res = -1
  for (let i = 1; i <= n; i++) {
    res = Math.max(res, dist[i])
  }

  return res === Infinity ? -1 : res
}

export {}

console.log(
  networkDelayTime(
    [
      [2, 1, 1],
      [2, 3, 1],
      [3, 4, 1],
    ],
    4,
    2
  )
)
