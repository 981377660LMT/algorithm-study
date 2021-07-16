// 题目给出的是一个图，类似于大西洋太平洋的题目需要dfs
const numIslands = (grid: string[][]): number => {
  let res = 0

  // 1. 确定行列
  const m = grid.length
  const n = grid[0].length

  // 2. dfs: 碰到为1的陆地就开始深度遍历消除所有陆地
  const dfs = (row: number, column: number) => {
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
        grid[nextRow][nextColumn] === '1'
      ) {
        grid[nextRow][nextColumn] = '0'
        dfs(nextRow, nextColumn)
      }
    })
  }

  // 3. 开始dfs
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (grid[i][j] === '1') {
        res++
        dfs(i, j)
      }
    }
  }

  return res
}

console.log(
  numIslands([
    ['1', '1', '1', '1', '0'],
    ['1', '1', '0', '1', '0'],
    ['1', '1', '0', '0', '0'],
    ['0', '0', '0', '0', '0'],
  ])
)
// 输出：1

console.log(
  numIslands([
    ['1', '1', '0', '0', '0'],
    ['1', '1', '0', '0', '0'],
    ['0', '0', '1', '0', '0'],
    ['0', '0', '0', '1', '1'],
  ])
)
// 输出：3
export {}
