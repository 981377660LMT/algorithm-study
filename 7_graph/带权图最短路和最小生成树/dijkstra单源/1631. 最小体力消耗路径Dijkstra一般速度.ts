import { MinHeap } from '../../../2_queue/minheap'

// 单源最短距离问题：权值已知
type CurX = number
type CurY = number
type Weight = number
type Edge = [CurX, CurY, Weight]
function minimumEffortPath(heights: number[][]): number {
  const m = heights.length
  const n = heights[0].length
  const directions = [
    [1, 0],
    [-1, 0],
    [0, 1],
    [0, -1],
  ]

  const dist = Array<number>(m * n).fill(Infinity)
  dist[0] = 0

  const visited = new Set<number>()

  const comparator = (a: Edge, b: Edge) => a[2] - b[2]
  const priorityQueue = new MinHeap<Edge>(comparator)
  priorityQueue.push([0, 0, 0])

  while (priorityQueue.size) {
    // 1.每次都从离原点最近的没更新过的点开始更新(性能瓶颈：可使用优先队列优化成ElogE)
    const [curX, curY, maxWeight] = priorityQueue.shift()!
    if (curX === m - 1 && curY === n - 1) return maxWeight
    const key = curX * n + curY
    if (visited.has(key)) continue

    // 2.加入visited
    visited.add(key)

    // 3.利用cur点来更新其相邻节点next与原点的距离
    for (const [dx, dy] of directions) {
      const nextX = curX + dx
      const nextY = curY + dy
      if (nextX >= 0 && nextY >= 0 && nextX < m && nextY < n && !visited.has(nextX * n + nextY)) {
        priorityQueue.push([
          nextX,
          nextY,
          Math.max(maxWeight, Math.abs(heights[nextX][nextY] - heights[curX][curY])),
        ])
      }
    }
  }

  return -1
}

console.log(
  minimumEffortPath([
    [1, 2, 2],
    [3, 8, 2],
    [5, 3, 5],
  ])
)
