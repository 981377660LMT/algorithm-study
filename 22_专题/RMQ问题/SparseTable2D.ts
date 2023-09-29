// 二维RMQ/二维ST表
// https://codeforces.com/blog/entry/45485
// Preprocess : O(n*m*logn*logm)

class SparseTable2D<S> {
  private readonly _e: () => S
  private readonly _op: (a: S, b: S) => S
  private readonly _dp: S[]
  private readonly _row: number
  private readonly _col: number
  private readonly _hash1: number
  private readonly _hash2: number
  private readonly _hash3: number

  constructor(matrix: ArrayLike<ArrayLike<S>>, e: () => S, op: (a: S, b: S) => S) {
    const row = matrix.length
    const col = matrix[0].length
    const rowSize = 32 - Math.clz32(row)
    const colSize = 32 - Math.clz32(col)
    const dp: S[] = Array(rowSize * row * colSize * col)
    for (let i = 0; i < dp.length; i++) dp[i] = e()

    const hash1 = row * colSize * col
    const hash2 = colSize * col
    const hash3 = col

    for (let ir = 0; ir < row; ir++) {
      for (let ic = 0; ic < col; ic++) {
        dp[ir * hash2 + ic] = matrix[ir][ic]
      }
      for (let jc = 1; jc < colSize; jc++) {
        for (let ic = 0; ic + (1 << jc) <= col; ic++) {
          dp[ir * hash2 + jc * hash3 + ic] = op(
            dp[ir * hash2 + (jc - 1) * hash3 + ic],
            dp[ir * hash2 + (jc - 1) * hash3 + ic + (1 << (jc - 1))]
          )
        }
      }
    }

    for (let jr = 1; jr < rowSize; jr++) {
      for (let ir = 0; ir + (1 << jr) <= row; ir++) {
        for (let jc = 0; jc < colSize; jc++) {
          for (let ic = 0; ic + (1 << jc) <= col; ic++) {
            dp[jr * hash1 + ir * hash2 + jc * hash3 + ic] = op(
              dp[(jr - 1) * hash1 + ir * hash2 + jc * hash3 + ic],
              dp[(jr - 1) * hash1 + (ir + (1 << (jr - 1))) * hash2 + jc * hash3 + ic]
            )
          }
        }
      }
    }

    this._e = e
    this._op = op
    this._dp = dp
    this._row = row
    this._col = col
    this._hash1 = hash1
    this._hash2 = hash2
    this._hash3 = hash3
  }

  /**
   * 查询闭区间`[row1,col1,row2,col2]`的贡献值.
   * 0 <= row1 <= row2 < len(matrix).
   * 0 <= col1 <= col2 < len(matrix[0]).
   */
  query(row1: number, col1: number, row2: number, col2: number): S {
    if (row1 < 0) row1 = 0
    if (col1 < 0) col1 = 0
    if (row2 >= this._row) row2 = this._row - 1
    if (col2 >= this._col) col2 = this._col - 1
    if (row1 > row2 || col1 > col2) return this._e()
    const rowK = 32 - Math.clz32(row2 - row1 + 1) - 1
    const colK = 32 - Math.clz32(col2 - col1 + 1) - 1
    const { _hash1, _hash2, _hash3, _dp, _op } = this
    const res1 = _op(
      _dp[rowK * _hash1 + row1 * _hash2 + colK * _hash3 + col1],
      _dp[rowK * _hash1 + row1 * _hash2 + colK * _hash3 + col2 - (1 << colK) + 1]
    )
    const res2 = _op(
      _dp[rowK * _hash1 + (row2 - (1 << rowK) + 1) * _hash2 + colK * _hash3 + col1],
      _dp[
        rowK * _hash1 + (row2 - (1 << rowK) + 1) * _hash2 + colK * _hash3 + col2 - (1 << colK) + 1
      ]
    )
    return _op(res1, res2)
  }
}

class SparseTable2DInt32 {
  private readonly _e: () => number
  private readonly _op: (a: number, b: number) => number
  private readonly _dp: Int32Array
  private readonly _row: number
  private readonly _col: number
  private readonly _hash1: number
  private readonly _hash2: number
  private readonly _hash3: number

  constructor(
    matrix: ArrayLike<ArrayLike<number>>,
    e: () => number,
    op: (a: number, b: number) => number
  ) {
    const row = matrix.length
    const col = matrix[0].length
    const rowSize = 32 - Math.clz32(row)
    const colSize = 32 - Math.clz32(col)
    const dp = new Int32Array(rowSize * row * colSize * col).fill(e())

    const hash1 = row * colSize * col
    const hash2 = colSize * col
    const hash3 = col

    for (let ir = 0; ir < row; ir++) {
      for (let ic = 0; ic < col; ic++) {
        dp[ir * hash2 + ic] = matrix[ir][ic]
      }
      for (let jc = 1; jc < colSize; jc++) {
        for (let ic = 0; ic + (1 << jc) <= col; ic++) {
          dp[ir * hash2 + jc * hash3 + ic] = op(
            dp[ir * hash2 + (jc - 1) * hash3 + ic],
            dp[ir * hash2 + (jc - 1) * hash3 + ic + (1 << (jc - 1))]
          )
        }
      }
    }

    for (let jr = 1; jr < rowSize; jr++) {
      for (let ir = 0; ir + (1 << jr) <= row; ir++) {
        for (let jc = 0; jc < colSize; jc++) {
          for (let ic = 0; ic + (1 << jc) <= col; ic++) {
            dp[jr * hash1 + ir * hash2 + jc * hash3 + ic] = op(
              dp[(jr - 1) * hash1 + ir * hash2 + jc * hash3 + ic],
              dp[(jr - 1) * hash1 + (ir + (1 << (jr - 1))) * hash2 + jc * hash3 + ic]
            )
          }
        }
      }
    }

    this._e = e
    this._op = op
    this._dp = dp
    this._row = row
    this._col = col
    this._hash1 = hash1
    this._hash2 = hash2
    this._hash3 = hash3
  }

  /**
   * 查询闭区间`[row1,col1,row2,col2]`的贡献值.
   * 0 <= row1 <= row2 < len(matrix).
   * 0 <= col1 <= col2 < len(matrix[0]).
   */
  query(row1: number, col1: number, row2: number, col2: number): number {
    if (row1 < 0) row1 = 0
    if (col1 < 0) col1 = 0
    if (row2 >= this._row) row2 = this._row - 1
    if (col2 >= this._col) col2 = this._col - 1
    if (row1 > row2 || col1 > col2) return this._e()
    const rowK = 32 - Math.clz32(row2 - row1 + 1) - 1
    const colK = 32 - Math.clz32(col2 - col1 + 1) - 1
    const { _hash1, _hash2, _hash3, _dp, _op } = this
    const res1 = _op(
      _dp[rowK * _hash1 + row1 * _hash2 + colK * _hash3 + col1],
      _dp[rowK * _hash1 + row1 * _hash2 + colK * _hash3 + col2 - (1 << colK) + 1]
    )
    const res2 = _op(
      _dp[rowK * _hash1 + (row2 - (1 << rowK) + 1) * _hash2 + colK * _hash3 + col1],
      _dp[
        rowK * _hash1 + (row2 - (1 << rowK) + 1) * _hash2 + colK * _hash3 + col2 - (1 << colK) + 1
      ]
    )
    return _op(res1, res2)
  }
}

export { SparseTable2D, SparseTable2DInt32 }

if (require.main === module) {
  const st2d = new SparseTable2DInt32(
    [
      [1, 2, 3],
      [4, 5, 6],
      [7, 8, 9]
    ],
    () => 0,
    Math.max
  )

  console.log(st2d.query(1, 2, 1, 2))
}
