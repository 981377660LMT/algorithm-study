// 给你一个由 0、1 和 2 组成的二维网格，其中：

// 0 代表你可以自由通过和选择建造的空地
// 1 代表你无法通行的建筑物
// 2 代表你无法通行的障碍物
// 通过调研，你希望从它出发能在 最短的距离和 内抵达周边全部的建筑物。请你计算出这个最佳的选址到周边全部建筑物的 最短距离和。
// 请你计算出这个最佳的选址到周边全部建筑物的 最短距离和。

// 类似曼哈顿距离但不是 因为有障碍
// 需要从每个建筑物出发bfs
function shortestDistance(grid: number[][]): number {
  if (grid.length === 0 || grid[0].length === 0) return -1

  const [m, n] = [grid.length, grid[0].length]
  // 有多少个建筑物可以到达
  const reachableCount = Array.from<unknown, number[]>({ length: m }, () => Array(n).fill(0))
  // 可以到达的建筑物们的距离和
  const distanceSum = Array.from<unknown, number[]>({ length: m }, () => Array(n).fill(0))
  // 建筑物1的数量
  let buildingCount = 0

  for (let r = 0; r < m; r++) {
    for (let c = 0; c < n; c++) {
      if (grid[r][c] !== 1) continue
      // 需要从每个建筑物出发bfs
      bfs(r, c, reachableCount, distanceSum)
    }
  }

  let res = Infinity
  for (let r = 0; r < m; r++) {
    for (let c = 0; c < n; c++) {
      if (grid[r][c] !== 0) continue
      // 每个建筑物都可到达的空地
      if (reachableCount[r][c] === buildingCount) res = Math.min(res, distanceSum[r][c])
    }
  }

  return res === Infinity ? -1 : res

  function bfs(
    startRow: number,
    startCol: number,
    reachableCount: number[][],
    distanceSum: number[][]
  ): void {
    const queue: [row: number, col: number, dist: number][] = []
    const visited = Array.from<unknown, boolean[]>({ length: m }, () => Array(n).fill(false))
    buildingCount++
    queue.push([startRow, startCol, 0])

    while (queue.length > 0) {
      const [row, col, dist] = queue.shift()!
      for (const [nextRow, nextCol] of [
        [row + 1, col],
        [row - 1, col],
        [row, col + 1],
        [row, col - 1],
      ]) {
        if (
          nextRow >= 0 &&
          nextRow < m &&
          nextCol >= 0 &&
          nextCol < n &&
          grid[nextRow][nextCol] === 0 &&
          !visited[nextRow][nextCol]
        ) {
          reachableCount[nextRow][nextCol]++
          distanceSum[nextRow][nextCol] += dist + 1
          queue.push([nextRow, nextCol, dist + 1])
          visited[nextRow][nextCol] = true
        }
      }
    }
  }
}

console.log(
  shortestDistance([
    [1, 0, 2, 0, 1],
    [0, 0, 0, 0, 0],
    [0, 0, 1, 0, 0],
  ])
)
export {}
