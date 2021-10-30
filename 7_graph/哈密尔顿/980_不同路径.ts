/**
 * @param {number[][]} grid  1 <= grid.length * grid[0].length <= 20
 * @return {number}
 * @description 
 * 1 表示起始方格。且只有一个起始方格。
   2 表示结束方格，且只有一个结束方格。
   0 表示我们可以走过的空方格。
   -1 表示我们无法跨越的障碍。
   每一个无障碍方格都要通过一次，但是一条路径中不能重复通过同一个方格。
   @summary 哈密尔顿路径问题
 */
const uniquePathsIII = (grid: number[][]): number => {
  if (grid.length === 0) return 0
  let res = 0
  const row = grid.length
  const col = grid[0].length
  // 注意这里
  let emptyCount = 0
  let start: number[] = []
  let end: number[] = []

  const next = [
    [-1, 0],
    [0, 1],
    [1, 0],
    [0, -1],
  ]

  for (let i = 0; i < row; i++) {
    for (let j = 0; j < col; j++) {
      if (grid[i][j] !== -1) emptyCount++
      if (grid[i][j] === 1) start = [i, j]
      if (grid[i][j] === 2) end = [i, j]
    }
  }

  const dfs = (x: number, y: number, count: number) => {
    // 1. 回溯终点
    if (grid[x][y] === -1 || grid[x][y] === Infinity) return
    if (x === end[0] && y === end[1]) {
      if (count === emptyCount) res++
      return
    }

    grid[x][y] = Infinity

    for (const [dx, dy] of next) {
      const nextRow = x + dx
      const nextColumn = y + dy
      // 2.在矩阵中
      if (nextRow >= 0 && nextRow < row && nextColumn >= 0 && nextColumn < col) {
        dfs(nextRow, nextColumn, count + 1)
      }
    }

    // 3. 回溯重置
    grid[x][y] = 0
  }

  // 4.每个点开始回溯 count开始为1
  dfs(start[0], start[1], 1)

  return res
}

console.log(
  uniquePathsIII([
    [1, 0, 0, 0],
    [0, 0, 0, 0],
    [0, 0, 2, -1],
  ])
)

export {}
