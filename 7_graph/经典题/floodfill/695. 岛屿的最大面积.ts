// 找到给定的二维数组中最大的岛屿面积。
function maxAreaOfIsland(grid: number[][]): number {
  let res = 0
  let resCand = 0

  const row = grid.length
  const col = grid[0].length

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (grid[r][c] === 1) {
        resCand = 0
        dfs(r, c)
        res = Math.max(res, resCand)
      }
    }
  }

  return res

  //  dfs原地标记
  function dfs(row: number, col: number): void {
    if (grid[row][col] === 1) {
      grid[row][col] = 0
      resCand++
    }

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
        nextRow < row &&
        nextColumn >= 0 &&
        nextColumn < col &&
        grid[nextRow][nextColumn] === 1
      ) {
        dfs(nextRow, nextColumn)
      }
    })
  }
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
