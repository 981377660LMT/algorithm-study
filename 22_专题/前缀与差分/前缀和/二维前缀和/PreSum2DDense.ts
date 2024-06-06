class PreSum2DDense {
  private readonly _preSum: number[]
  private readonly _col: number

  constructor(mat: ArrayLike<ArrayLike<number>>)
  constructor(row: number, col: number, f: (r: number, c: number) => number)
  constructor(arg0: any, arg1?: any, arg2?: any) {
    let row: number
    let col: number
    if (typeof arg0 === 'number') {
      row = arg0
      col = arg1
    } else {
      row = arg0.length
      col = arg0[0].length
    }
    const preSum = Array<number>((row + 1) * (col + 1)).fill(0)

    if (typeof arg0 === 'number') {
      for (let r = 0; r < row; r++) {
        const offset0 = r * (col + 1)
        const offset1 = (r + 1) * (col + 1)
        for (let c = 0; c < col; c++) {
          preSum[offset1 + c + 1] =
            arg2(r, c) + preSum[offset0 + c + 1] + preSum[offset1 + c] - preSum[offset0 + c]
        }
      }
    } else {
      for (let r = 0; r < row; r++) {
        const offset0 = r * (col + 1)
        const offset1 = (r + 1) * (col + 1)
        const matRow = arg0[r]
        for (let c = 0; c < col; c++) {
          preSum[offset1 + c + 1] =
            matRow[c] + preSum[offset0 + c + 1] + preSum[offset1 + c] - preSum[offset0 + c]
        }
      }
    }

    this._preSum = preSum
    this._col = col
  }

  /**
   * 查询sum(A[x1:x2+1, y1:y2+1])的值(包含边界).
   */
  queryRange(x1: number, x2: number, y1: number, y2: number): number {
    const col = this._col + 1
    return (
      this._preSum[(x2 + 1) * col + y2 + 1] -
      this._preSum[(x2 + 1) * col + y1] -
      this._preSum[x1 * col + y2 + 1] +
      this._preSum[x1 * col + y1]
    )
  }
}

if (require.main === module) {
  const S = new PreSum2DDense([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
  ])
  console.log(S.queryRange(1, 1, 2, 2))

  const S2 = new PreSum2DDense(3, 3, (r, c) => r * 3 + c + 1)
  console.log(S2.queryRange(1, 1, 2, 2))
}

export { PreSum2DDense }
