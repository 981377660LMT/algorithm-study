function numSpecial(mat: number[][]): number {
  const row = mat.length
  const col = mat[0].length
  const rowCount = Array<number>(row).fill(0)
  const colCount = Array<number>(col).fill(0)

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (mat[r][c] === 1) {
        rowCount[r]++
        colCount[c]++
      }
    }
  }

  let res = 0

  for (let r = 0; r < row; r++) {
    for (let c = 0; c < col; c++) {
      if (check(r, c)) res++
    }
  }

  return res

  function check(row: number, col: number): boolean {
    return mat[row][col] === 1 && rowCount[row] === 1 && colCount[col] === 1
  }
}

console.log(
  numSpecial([
    [1, 0, 0],
    [0, 0, 1],
    [1, 0, 0],
  ])
)
// 特殊位置 定义：如果 mat[i][j] == 1 并且第 i 行和第 j 列中的所有其他元素均为 0（行和列的下标均 从 0 开始 ），则位置 (i, j) 被称为特殊位置。
// 请返回 矩阵 mat 中特殊位置的数目 。

export {}
