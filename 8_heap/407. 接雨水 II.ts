import { MinHeap } from './Heap'

/**
 *
 * @param heightMap
 * 最初，小根堆堆包含地图周边的所有单元格。
 * 检查 4 个邻居，忽略任何已经探索过的（在堆中或已经从堆中弹出）或地图外的
 * 如果新发现的邻居低于当前单元格，则可以在其顶部添加水
 * 将水添加到此相邻单元格（如果相邻单元格高于单元格则为零）并将具有新高度的邻居推入堆中
 * Every cell is pushed and popped from the heap, which is O(mn log(mn))
 */
function trapRainWater(heightMap: number[][]): number {
  if (!heightMap.length || !heightMap[0].length) return 0

  let res = 0
  const m = heightMap.length
  const n = heightMap[0].length
  const pq = new MinHeap<[value: number, row: number, col: number]>((a, b) => a[0] - b[0])
  const visited = new Set<number>()
  // 四个边
  for (let r = 0; r < m; r++) {
    pq.heappush([heightMap[r][0], r, 0])
    visited.add(r * n)
    pq.heappush([heightMap[r][n - 1], r, n - 1])
    visited.add(r * n + n - 1)
  }
  for (let c = 1; c < n - 1; c++) {
    pq.heappush([heightMap[0][c], 0, c])
    visited.add(c)
    pq.heappush([heightMap[m - 1][c], m - 1, c])
    visited.add((m - 1) * n + c)
  }

  // 优先队列bfs
  while (pq.size) {
    const [height, row, col] = pq.heappop()!
    for (const [dr, dc] of [
      [1, 0],
      [-1, 0],
      [0, 1],
      [0, -1]
    ]) {
      const [nextR, nextC] = [row + dr, col + dc]
      const key = nextR * n + nextC
      if (!visited.has(key) && nextR >= 0 && nextR < m && nextC >= 0 && nextC < n) {
        visited.add(key)
        res += Math.max(0, height - heightMap[nextR][nextC])
        pq.heappush([Math.max(height, heightMap[nextR][nextC]), nextR, nextC])
      }
    }
  }

  return res
}

console.log(
  trapRainWater([
    [1, 4, 3, 1, 3, 2],
    [3, 2, 1, 3, 2, 4],
    [2, 3, 3, 2, 3, 1]
  ])
)
// 输出: 4
// 解释: 下雨后，雨水将会被上图蓝色的方块中。总的接雨水量为1+2+1=4。
