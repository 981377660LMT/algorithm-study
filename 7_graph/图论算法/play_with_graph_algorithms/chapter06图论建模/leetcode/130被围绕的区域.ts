/**
 * @param {character[][]} board
 * @return {void} Do not return anything, modify board in-place instead.
 * @description 任何不在边界上，或不与边界上的 'O' 相连的 'O' 最终都会被填充为 'X'
 * @summary 将非边界的点填成第三种颜色
 */
const solve = (board: string[][]): void => {
  const m = board.length
  const n = board[0].length

  const dfs = (row: number, column: number) => {
    board[row][column] = 'F'
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
        board[nextRow][nextColumn] === 'O'
      ) {
        dfs(nextRow, nextColumn)
      }
    })
  }

  // 思路就是从边缘dfs遍历连接边缘的 把这些点拯救出来
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (board[i][j] === 'O' && ([0, m - 1].includes(i) || [0, n - 1].includes(j))) dfs(i, j)
    }
  }

  // 遍历完后清除
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (board[i][j] === 'O') board[i][j] = 'X'
      if (board[i][j] === 'F') board[i][j] = 'O'
    }
  }

  console.table(board)
}

console.table(
  solve([
    ['X', 'X', 'X', 'X'],
    ['X', 'O', 'O', 'X'],
    ['X', 'X', 'O', 'X'],
    ['X', 'O', 'X', 'X'],
  ])
)

export {}
