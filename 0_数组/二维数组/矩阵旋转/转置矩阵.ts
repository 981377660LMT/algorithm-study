function transpose(matrix: number[][]): number[][] {
  const m = matrix.length
  const n = matrix[0].length
  const res = Array.from<unknown, number[]>({ length: n }, () => Array(m).fill(0))

  for (let i = 0; i < m; i++) {
    for (let j = 0; j < n; j++) {
      res[j][i] = matrix[i][j]
    }
  }

  return res
}
