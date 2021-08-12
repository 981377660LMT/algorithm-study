// 找到给定的二维数组中最大的岛屿面积。
const maxAreaOfIsland = (grid: number[][]): number => {
  let res = 0

  // 1. 确定行列
  const m = grid.length
  const n = grid[0].length

  // 2. dfs: 碰到为1的陆地就开始深度遍历并标记为已看过(0)
  const dfs = (row: number, column: number, visited: Set<string>) => {
    res = Math.max(res, visited.size)
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
        grid[nextRow][nextColumn] === 1 &&
        !visited.has(`${nextRow}#${nextColumn}`)
      ) {
        visited.add(`${nextRow}#${nextColumn}`)
        dfs(nextRow, nextColumn, visited)
        // 这句很关键 回溯时标记为0 可以避免之后无效的dfs
        grid[nextRow][nextColumn] = 0
      }
    })
  }

  // 3. 开始dfs
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (grid[i][j] === 1) {
        dfs(i, j, new Set([`${i}#${j}`]))
      }
    }
  }

  return res
}

console.log(
  maxAreaOfIsland([
    [1, 1],
    [1, 1],
  ])
)
// 输出：4
console.log(
  maxAreaOfIsland([
    [0, 0, 1, 0, 0, 0, 0, 1, 0, 0, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0],
    [0, 1, 1, 0, 1, 0, 0, 0, 0, 0, 0, 0, 0],
    [0, 1, 0, 0, 1, 1, 0, 0, 1, 0, 1, 0, 0],
    [0, 1, 0, 0, 1, 1, 0, 0, 1, 1, 1, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0],
    [0, 0, 0, 0, 0, 0, 0, 1, 1, 0, 0, 0, 0],
  ])
)
// 6
export {}
