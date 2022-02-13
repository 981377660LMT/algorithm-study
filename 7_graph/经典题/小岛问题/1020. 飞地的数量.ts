/**
 * @param {number[][]} grid
 * @return {number}
 * 返回网格中无法在任意次数的移动中离开网格边界的陆地单元格的数量。
 */
const numEnclaves = function (grid: number[][]): number {
  let res = 0

  // 1. 确定行列
  const m = grid.length
  const n = grid[0].length

  // 3. 边界开始dfs
  for (let i = 0; i < m; i++) {
    if (grid[i][0] === 1) dfs(i, 0)
    if (grid[i][n - 1] === 1) dfs(i, n - 1)
  }
  for (let j = 0; j < n; j++) {
    if (grid[0][j] === 1) dfs(0, j)
    if (grid[m - 1][j] === 1) dfs(m - 1, j)
  }

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (grid[i][j] === 1) {
        res++
      }
    }
  }

  return res

  // 2. dfs: 碰到为1的陆地就开始深度遍历并标记为已看过(0)
  function dfs(row: number, column: number) {
    grid[row][column] = 0
    ;[
      [row - 1, column],
      [row + 1, column],
      [row, column - 1],
      [row, column + 1],
    ].forEach(([nextRow, nextColumn]) => {
      // 1.在矩阵中
      // 2.是陆地
      if (
        nextRow >= 0 &&
        nextRow < m &&
        nextColumn >= 0 &&
        nextColumn < n &&
        grid[nextRow][nextColumn] === 1
      ) {
        dfs(nextRow, nextColumn)
      }
    })
  }
}

console.log(
  numEnclaves([
    [0, 0, 0, 0],
    [1, 0, 1, 0],
    [0, 1, 1, 0],
    [0, 0, 0, 0],
  ])
)

export {}
