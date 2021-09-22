// 指定两个节点分别作为起点 start 和终点 end ，
// 请你找出从起点到终点成功概率最大的路径，
// 并返回其成功概率。

import { MinHeap } from '../../../2_queue/minheap'

type Cur = number
type Weight = number
type Edge = [Cur, Weight]

/**
 * @param {number} n
 * @param {number[][]} edges
 * @param {number[]} succProb
 * @param {number} start
 * @param {number} end
 * @return {number}
 */
const maxProbability = function (
  n: number,
  edges: number[][],
  succProb: number[],
  start: number,
  end: number
): number {
  // 0.建图(也可以不建，只是方便获取next)
  const adjList = Array.from<number, [number, number][]>({ length: n }, () => [])
  for (let i = 0; i < edges.length; i++) {
    const [u, v] = edges[i]
    const w = succProb[i]
    adjList[u].push([v, w])
    adjList[v].push([u, w])
  }

  // 1.dist数组
  const dist = Array<number>(n).fill(0)
  dist[start] = 1

  // const visited = new Set<number>()

  // 2.pq优先队列
  const compareFunction = (a: Edge, b: Edge) => b[1] - a[1]
  const priorityQueue = new MinHeap<Edge>(compareFunction)
  priorityQueue.push([start, 1])

  while (priorityQueue.size) {
    // 3.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
    const [cur, maxWeight] = priorityQueue.shift()!
    // if (visited.has(cur)) continue
    if (cur === end) return maxWeight

    // visited.add(cur)

    // 4.利用cur点来更新其相邻节点next与原点的距离
    for (const [next, weight] of adjList[cur]) {
      // if (visited.has(next)) continue
      if (dist[cur] * weight > dist[next]) {
        dist[next] = dist[cur] * weight
        priorityQueue.push([next, dist[next]])
      }
    }
  }

  return 0
}

console.log(
  maxProbability(
    3,
    [
      [0, 1],
      [1, 2],
      [0, 2],
    ],
    [0.5, 0.5, 0.2],
    0,
    2
  )
)
