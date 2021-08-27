const rotate = (matrix: number[][]): void => {
  matrix.reverse()
  matrix = Array.from({ length: matrix.length }, (_, i) =>
    //   第i行
    Array(matrix.length)
      .fill(i)
      .map((_, row) => matrix[row][i])
  )

  console.table(matrix)
}

rotate([
  [1, 2, 3],
  [4, 5, 6],
  [7, 8, 9],
])

export {}
