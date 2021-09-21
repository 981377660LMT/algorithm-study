import { MinHeap } from './minheap'

type Value = number
type Row = number
type Col = number
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
  const rows = heightMap.length
  const cols = heightMap[0].length
  const pq = new MinHeap<[Value, Row, Col]>((a, b) => a[0] - b[0])
  const visited = new Set<string>()
  // 四个边
  for (let r = 0; r < rows; r++) {
    pq.push([heightMap[r][0], r, 0])
    visited.add(`${r}#${0}`)
    pq.push([heightMap[r][cols - 1], r, cols - 1])
    visited.add(`${r}#${cols - 1}`)
  }
  for (let c = 1; c < cols - 1; c++) {
    pq.push([heightMap[0][c], 0, c])
    visited.add(`${0}#${c}`)
    pq.push([heightMap[rows - 1][c], rows - 1, c])
    visited.add(`${rows - 1}#${c}`)
  }

  // 优先队列bfs
  while (pq.size) {
    const [h, r, c] = pq.shift()!
    for (const [dr, dc] of [
      [1, 0],
      [-1, 0],
      [0, 1],
      [0, -1],
    ]) {
      const [nextR, nextC] = [r + dr, c + dc]
      console.log(h, r, c)
      const key = `${nextR}#${nextC}`
      if (!visited.has(key) && nextR >= 0 && nextR < rows && nextC >= 0 && nextC < cols) {
        visited.add(key)
        res += Math.max(0, h - heightMap[nextR][nextC])
        pq.push([Math.max(h, heightMap[nextR][nextC]), nextR, nextC])
      }
    }
  }

  return res
}

console.log(
  trapRainWater([
    [1, 4, 3, 1, 3, 2],
    [3, 2, 1, 3, 2, 4],
    [2, 3, 3, 2, 3, 1],
  ])
)
// 输出: 4
// 解释: 下雨后，雨水将会被上图蓝色的方块中。总的接雨水量为1+2+1=4。
