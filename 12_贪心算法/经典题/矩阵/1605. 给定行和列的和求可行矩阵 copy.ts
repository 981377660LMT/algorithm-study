function restoreMatrix(rowSum: number[], colSum: number[]): number[][] {
  const [row, col] = [rowSum.length, colSum.length]
  const res = Array.from<unknown, number[]>({ length: row }, () => Array(col).fill(0))

  for (let r = 0; r < row; r++) {
    if (rowSum[r] === 0) continue
    for (let c = 0; c < col; c++) {
      if (colSum[c] === 0) continue
      const cur = Math.min(rowSum[r], colSum[c])
      res[r][c] = cur
      rowSum[r] -= cur
      colSum[c] -= cur
    }
  }

  return res
}
