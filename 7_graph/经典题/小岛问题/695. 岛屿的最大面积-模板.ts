// 找到给定的二维数组中最大的岛屿面积。
const maxAreaOfIsland = (grid: number[][]): number => {
  let res = 0
  let resCand = 0

  // 1. 确定行列
  const m = grid.length
  const n = grid[0].length

  // 3. 开始dfs
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (grid[i][j] === 1) {
        resCand = 0
        dfs(i, j)
        res = Math.max(res, resCand)
      }
    }
  }

  return res

  // 2. dfs原地标记,同时cand++
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
