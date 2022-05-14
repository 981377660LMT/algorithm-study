/**
 * @param {number[][]} grid
 * @return {number}
 * @description 遍历 多一个1就+4然后减去上下左右为1的个数
 */
const islandPerimeter = function (grid: number[][]): number {
  let res = 0
  const row = grid.length
  const col = grid[0].length

  for (let i = 0; i < row; i++) {
    for (let j = 0; j < col; j++) {
      const cur = grid[i][j]
      if (cur === 1) {
        res += 4
        ;[
          [i - 1, j],
          [i + 1, j],
          [i, j - 1],
          [i, j + 1],
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
            res--
          }
        })
      }
    }
  }

  return res
}

console.log(
  islandPerimeter([
    [0, 1, 0, 0],
    [1, 1, 1, 0],
    [0, 1, 0, 0],
    [1, 1, 0, 0],
  ])
)
