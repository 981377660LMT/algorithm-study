// 找出一个海洋单元格，这个海洋单元格到离它最近的陆地单元格的距离是最大的。
// 多源bfs:想象存在虚拟节点变为单源bfs
const maxDistance = (grid: number[][]): number => {
  let res = 0

  // 1. 确定行列
  const m = grid.length
  const n = grid[0].length

  const queue: number[][] = []
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      if (grid[i][j] === 1) {
        queue.push([i, j, 0])
      }
    }
  }

  while (queue.length) {
    const [row, column, step] = queue.shift()!
    res = Math.max(res, step)
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
        grid[nextRow][nextColumn] === 0
      ) {
        grid[nextRow][nextColumn] = 1
        queue.push([nextRow, nextColumn, step + 1])
      }
    })
  }

  return res === 0 ? -1 : res
}

console.log(
  maxDistance([
    [1, 0, 1],
    [0, 0, 0],
    [1, 0, 1],
  ])
)
// 输出：2

console.log(
  maxDistance([
    [1, 0, 0],
    [0, 0, 0],
    [0, 0, 0],
  ])
)
// 输出：4
export {}
