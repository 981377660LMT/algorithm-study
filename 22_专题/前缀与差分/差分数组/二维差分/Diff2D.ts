// 二维差分/二维前缀和

class Diff2D {
  /** 内部矩阵.*/
  readonly matrix: number[]
  private readonly _row: number
  private readonly _col: number
  private readonly _diff: number[]
  private _dirty = false

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

    const newMat = Array(row * col)
    if (typeof arg0 === 'number') {
      for (let i = 0; i < row; i++) {
        const offset = i * col
        for (let j = 0; j < col; j++) {
          newMat[offset + j] = arg2(i, j)
        }
      }
    } else {
      for (let i = 0; i < row; i++) {
        const matRow = arg0[i]
        const offset = i * col
        for (let j = 0; j < col; j++) {
          newMat[offset + j] = matRow[j]
        }
      }
    }

    const diff = Array((row + 2) * (col + 2)).fill(0)

    this.matrix = newMat
    this._row = row
    this._col = col
    this._diff = diff
  }

  /**
   * 区间更新 左上角` (row1, col1)` 到 右下角` (row2, col2) `闭区间所描述的子矩阵的元素
   */
  add(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    const col = this._col + 2
    this._diff[(row1 + 1) * col + col1 + 1] += delta
    this._diff[(row1 + 1) * col + col2 + 2] -= delta
    this._diff[(row2 + 2) * col + col1 + 1] -= delta
    this._diff[(row2 + 2) * col + col2 + 2] += delta
    this._dirty = true
  }

  query(row: number, col: number): number {
    if (this._dirty) this.build()
    return this.matrix[row * this._col + col]
  }

  /**
   * 遍历矩阵，还原对应元素的增量, 并更新矩阵{@link matrix}.
   */
  build(): void {
    if (!this._dirty) return
    this._dirty = false
    const diff = this._diff
    const mat = this.matrix
    for (let i = 0; i < this._row; i++) {
      const offsetD0 = i * (this._col + 2)
      const offsetD1 = (i + 1) * (this._col + 2)
      const offsetM = i * this._col
      for (let j = 0; j < this._col; j++) {
        diff[offsetD1 + j + 1] += diff[offsetD1 + j] + diff[offsetD0 + j + 1] - diff[offsetD0 + j]
        mat[offsetM + j] += diff[offsetD1 + j + 1]
      }
    }
  }
}

class PreSum2D {
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
          preSum[offset1 + c + 1] = arg2(r, c) + preSum[offset0 + c + 1] + preSum[offset1 + c] - preSum[offset0 + c]
        }
      }
    } else {
      for (let r = 0; r < row; r++) {
        const offset0 = r * (col + 1)
        const offset1 = (r + 1) * (col + 1)
        const matRow = arg0[r]
        for (let c = 0; c < col; c++) {
          preSum[offset1 + c + 1] = matRow[c] + preSum[offset0 + c + 1] + preSum[offset1 + c] - preSum[offset0 + c]
        }
      }
    }

    this._preSum = preSum
    this._col = col
  }

  /**
   * 返回 左上角 `(row1, col1)` 、右下角 `(row2, col2)` 闭区间所描述的子矩阵的元素总和 。
   */
  queryRange(row1: number, col1: number, row2: number, col2: number): number {
    const col = this._col + 1
    return (
      this._preSum[(row2 + 1) * col + col2 + 1] -
      this._preSum[(row2 + 1) * col + col1] -
      this._preSum[row1 * col + col2 + 1] +
      this._preSum[row1 * col + col1]
    )
  }
}

export { Diff2D, PreSum2D }

if (require.main === module) {
  const S = new Diff2D([
    [1, 2, 3],
    [4, 5, 6],
    [7, 8, 9]
  ])
  S.add(1, 1, 2, 2, 1)
  console.log(S.query(1, 1))

  const S2 = new Diff2D(3, 3, (r, c) => r * 3 + c + 1)
  S2.add(1, 1, 2, 2, 1)
  console.log(S2.query(1, 1))
}
