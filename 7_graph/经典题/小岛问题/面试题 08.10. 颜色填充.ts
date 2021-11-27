function floodFill(image: number[][], sr: number, sc: number, newColor: number): number[][] {
  // 1. 确定行列
  const m = image.length
  const n = image[0].length
  const dfs = (row: number, column: number, startColor: number) => {
    image[row][column] = Infinity // 做标记
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
        image[nextRow][nextColumn] === startColor
      ) {
        dfs(nextRow, nextColumn, startColor)
      }
    })
  }

  // 3. 开始dfs
  dfs(sr, sc, image[sr][sc])

  // 重置标记
  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      image[i][j] === Infinity && (image[i][j] = newColor)
    }
  }

  return image
}

console.log(
  floodFill(
    [
      [0, 0, 0],
      [0, 1, 1],
    ],
    1,
    1,
    1
  )
)
