/**
 * @param {number[][]} mat
 * @return {number[][]}
 * 请输出一个大小相同的矩阵，其中每一个格子是 mat 中对应位置元素到最近的 0 的距离。
 * @summary
 * 跟之前求海洋到陆地的最大距离的最小值那个做法一样，多源 BFS
 * 「1162.地图分析」
 */
var updateMatrix = function (mat: number[][]): number[][] {
  // 1. 确定行列
  const m = mat.length
  const n = mat[0].length
  const res = Array.from({ length: m }, () => Array(n).fill(Infinity))
  const visited = new Set<number>()

  // 从陆地向海洋多源bfs
  const queue: number[][] = []
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (mat[i][j] === 0) {
        res[i][j] = 0
        queue.push([i, j])
        visited.add(i * n + j)
      }
    }
  }

  while (queue.length) {
    const [row, col] = queue.shift()!

    ;[
      [row - 1, col],
      [row + 1, col],
      [row, col - 1],
      [row, col + 1],
    ].forEach(([nextRow, nextColumn]) => {
      // 1.在矩阵中
      // 2.是陆地
      if (
        nextRow >= 0 &&
        nextRow < m &&
        nextColumn >= 0 &&
        nextColumn < n &&
        !visited.has(nextRow * n + nextColumn)
      ) {
        res[nextRow][nextColumn] = res[row][col] + 1
        queue.push([nextRow, nextColumn])
        visited.add(nextRow * n + nextColumn)
      }
    })
  }

  return res
}

console.log(
  updateMatrix([
    [0, 0, 0],
    [0, 1, 0],
    [0, 0, 0],
  ])
)
// 输出：[[0,0,0],[0,1,0],[0,0,0]]
