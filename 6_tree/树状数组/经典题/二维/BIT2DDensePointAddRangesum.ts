class BIT2DDensePointAddRangeSum {
  private readonly _row: number
  private readonly _col: number
  private readonly _data: Float64Array

  constructor(row: number, col: number) {
    this._row = row
    this._col = col
    this._data = new Float64Array(row * col)
  }

  build(f: (r: number, c: number) => number): void {
    const ROW = this._row
    const COL = this._col
    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        this._data[COL * r + c] = f(r, c)
      }
    }
    for (let r = 1; r <= ROW; r++) {
      for (let c = 1; c <= COL; c++) {
        const nc = c + (c & -c)
        if (nc <= COL) {
          this._data[this._id(r, nc)] += this._data[this._id(r, c)]
        }
      }
    }
    for (let r = 1; r <= ROW; r++) {
      for (let c = 1; c <= COL; c++) {
        const nr = r + (r & -r)
        if (nr <= ROW) {
          this._data[this._id(nr, c)] += this._data[this._id(r, c)]
        }
      }
    }
  }

  /**
   * 点 (r,c) 的值加上 delta.
   */
  add(r: number, c: number, delta: number): void {
    r++
    for (let i = r; i <= this._row; i += i & -i) {
      this._add(i, c, delta)
    }
  }

  /**
   * [0,r) * [0,c).
   */
  queryPrefix(r: number, c: number): number {
    if (r > this._row) r = this._row
    if (c > this._col) c = this._col
    let res = 0
    for (let i = r; i > 0; i &= i - 1) {
      res += this._queryPrefix(i, c)
    }
    return res
  }

  /**
   * [r1,r2) * [c1,c2).
   */
  queryRange(r1: number, r2: number, c1: number, c2: number): number {
    if (r1 < 0) r1 = 0
    if (c1 < 0) c1 = 0
    if (r2 > this._row) r2 = this._row
    if (c2 > this._col) c2 = this._col
    if (r1 >= r2 || c1 >= c2) return 0
    let pos = 0
    let neg = 0
    for (; r1 < r2; r2 &= r2 - 1) {
      pos += this._queryRange(r2, c1, c2)
    }
    for (; r2 < r1; r1 &= r1 - 1) {
      neg += this._queryRange(r1, c1, c2)
    }
    return pos - neg
  }

  toString(): string {
    const res: string[] = []
    for (let r = 0; r < this._row; r++) {
      const row: number[] = []
      for (let c = 0; c < this._col; c++) {
        row.push(this.queryRange(r, r + 1, c, c + 1))
      }
      res.push(row.join(','))
    }
    return res.join('\n')
  }

  private _id(r: number, c: number): number {
    return (r - 1) * this._col + c - 1
  }

  private _add(r: number, c: number, delta: number): void {
    c++
    for (let i = c; i <= this._col; i += i & -i) {
      this._data[this._id(r, i)] += delta
    }
  }

  private _queryPrefix(r: number, c: number): number {
    let pos = 0
    for (let i = c; i > 0; i &= i - 1) {
      pos += this._data[this._id(r, i)]
    }
    return pos
  }

  private _queryRange(r: number, c1: number, c2: number): number {
    let pos = 0
    let neg = 0
    for (; c1 < c2; c2 &= c2 - 1) {
      pos += this._data[this._id(r, c2)]
    }
    for (; c2 < c1; c1 &= c1 - 1) {
      neg += this._data[this._id(r, c1)]
    }
    return pos - neg
  }
}

export { BIT2DDensePointAddRangeSum }

if (require.main === module) {
  const bit = new BIT2DDensePointAddRangeSum(10, 10)
  bit.add(1, 1, 1)
  bit.add(2, 2, 1)
  console.log(bit.toString())
  console.log(bit.queryPrefix(3, 3))
  console.log(bit.queryRange(1, 3, 1, 200))

  // https://leetcode.cn/problems/range-sum-query-2d-mutable/
  class NumMatrix {
    private readonly _bit: BIT2DDensePointAddRangeSum
    private readonly _matrix: number[][]

    constructor(matrix: number[][]) {
      this._matrix = matrix
      this._bit = new BIT2DDensePointAddRangeSum(matrix.length, matrix[0].length)
      this._bit.build((r, c) => matrix[r][c])
    }

    update(row: number, col: number, val: number): void {
      const pre = this._matrix[row][col]
      const diff = val - pre
      this._matrix[row][col] = val
      this._bit.add(row, col, diff)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._bit.queryRange(row1, row2 + 1, col1, col2 + 1)
    }
  }
}
