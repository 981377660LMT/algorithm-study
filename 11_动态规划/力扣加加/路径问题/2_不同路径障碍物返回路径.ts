// 返回一条可行的路径
function pathWithObstacles(obstacleGrid: number[][]): number[][] {
  // 1. 确定行列
  const m = obstacleGrid.length
  const n = obstacleGrid[0].length
  if (obstacleGrid[0][0] === 1 || obstacleGrid[m - 1][n - 1] === 1) return []
  const directions = [
    [1, 0],
    [0, 1],
  ]

  function* dfs(row: number, column: number, path: number[][]): Generator<number[][]> {
    if (row === m - 1 && column === n - 1) yield path

    for (const [dx, dy] of directions) {
      const nextRow = row + dx
      const nextColumn = column + dy
      // 1.在矩阵中
      // 2.是陆地
      if (
        nextRow >= 0 &&
        nextRow < m &&
        nextColumn >= 0 &&
        nextColumn < n &&
        obstacleGrid[nextRow][nextColumn] === 0
      ) {
        obstacleGrid[nextRow][nextColumn] = 1
        path.push([nextRow, nextColumn])
        yield* dfs(nextRow, nextColumn, path)
        path.pop()
      }
    }
  }

  // 3. 开始dfs
  return dfs(0, 0, [[0, 0]]).next().value || []
}

console.log(
  pathWithObstacles([
    [0, 0, 0, 0],
    [0, 1, 0, 0],
    [0, 0, 0, 0],
    [0, 0, 1, 0],
    [0, 0, 0, 0],
  ])
)

export {}
