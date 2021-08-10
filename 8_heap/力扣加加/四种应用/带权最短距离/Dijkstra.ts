import { PriorityQueue } from '../../../../2_queue/todo优先级队列'

/**
 *
 * @param graph 求带权图邻接表两点间距离，存储[顶点，权值]
 */
const dijkstra = (graph: [number, number][][], start: number, end: number): number => {
  // 堆里的数据都是 [cost, i] 的二元祖，其含义是"从 start 走到 i 的距离是 cost"
  const pq = new PriorityQueue<[number, number]>((a, b) => a[0] - b[0])
  pq.push([0, start])
  const visited = new Set<number>()

  while (pq.length) {
    console.dir(pq, { depth: null })
    const [cost, cur] = pq.shift()!
    if (cur === end) return cost
    if (!visited.has(cur)) {
      visited.add(cur)

      for (const [next, weight] of graph[cur]) {
        if (!visited.has(next)) {
          visited.add(next)
          pq.push([cost + weight, next])
        }
      }
    }
  }

  return -1
}

const graph: [number, number][][] = [
  [[1, 1]],
  [[2, 1]],
  [[4, 1]],
  [
    [0, 1],
    [5, 2],
  ],
  [],
  [[4, 1]],
]

console.dir(dijkstra(graph, 3, 1))
