/**
 * @param {number[][]} grid
 * @return {number}
 * @description 返回矩阵中最短 畅通路径 的长度 即左上角 单元格（即，(0, 0)）到 右下角 单元格（即，(n - 1, n - 1)）的路径
 * @description 相邻是八连通的
 *
 * @summary 求单源最短路径:BFS
 */
const shortestPathBinaryMatrix = (grid: number[][]): number => {
  const n = grid.length
  const visited = new Set<string>()
  const queue: [number, number, number][] = [[0, 0, 1]]
  let res = -1

  const bfs = (i: number, j: number, queue: [number, number, number][]) => {
    visited.add(`${i}#${j}`)
    while (queue.length) {
      const [preI, preJ, preLevel] = queue.shift()!
      const next = [
        [preI + 1, preJ + 1],
        [preI + 1, preJ],
        [preI + 1, preJ - 1],
        [preI, preJ + 1],
        [preI, preJ - 1],
        [preI - 1, preJ + 1],
        [preI - 1, preJ],
        [preI - 1, preJ - 1],
      ]
      for (const [nextI, nextJ] of next) {
        if (
          nextI < n &&
          nextI >= 0 &&
          nextJ < n &&
          nextJ >= 0 &&
          grid[nextI][nextJ] === 0 &&
          !visited.has(`${nextI}#${nextJ}`)
        ) {
          if (nextI === n - 1 && nextJ === n - 1) return (res = preLevel + 1)

          queue.push([nextI, nextJ, preLevel + 1])
          visited.add(`${nextI}#${nextJ}`)
        }
      }
    }
  }

  grid[0][0] === 0 && bfs(0, 0, queue)

  return res
}

console.log(
  shortestPathBinaryMatrix([
    [0, 0, 0],
    [1, 1, 0],
    [1, 1, 1],
  ])
)

export {}
