// 类型数组+fill的暴力更新
// COL小时,交换ROW和COL让COL更大
// 最坏:320*320,更新满1e5次,耗时: 1.6s
// https://leetcode.cn/submissions/detail/439137174/
// https://leetcode.cn/problems/subrectangle-queries/solution/typescript-shi-yong-typedarray-you-hua-d-cdgh/

/**
 * 二维区间染色,单点查询.
 * 值在0到2^32-1之间.
 */
class RangeAssignPointGet2DNaive {
  private readonly _rawRow: number
  private readonly _rawCol: number
  private readonly _needRotate: boolean
  private readonly _data: Uint32Array

  constructor(grid: ArrayLike<ArrayLike<number>>) {
    const row = grid.length
    const col = grid[0].length
    this._rawRow = row
    this._rawCol = col
    this._needRotate = row > col // 需要遍历的行数不超过sqrt,列上使用fill来加速染色
    this._data = new Uint32Array(row * col)

    if (!this._needRotate) {
      for (let r = 0; r < row; r++) {
        const cache = grid[r]
        const offset = r * col
        for (let c = 0; c < col; c++) {
          this._data[offset + c] = cache[c]
        }
      }
    } else {
      for (let r = 0; r < row; r++) {
        const cache = grid[r]
        for (let c = 0; c < col; c++) {
          const newPos = c * row + (row - r - 1) // (r,c) => (c,row-r-1)
          this._data[newPos] = cache[c]
        }
      }
    }
  }

  /**
   * 将`[row1,row2)`x`[col1,col2)`的区间染色为`value`.
   */
  update(row1: number, row2: number, col1: number, col2: number, value: number): void {
    let curCol = this._rawCol
    if (this._needRotate) {
      curCol = this._rawRow
      const tmp1 = row1
      const tmp2 = row2
      row1 = col1
      row2 = col2
      col1 = this._rawRow - tmp2
      col2 = this._rawRow - tmp1
    }

    for (let r = row1; r < row2; r++) {
      const offset = r * curCol
      this._data.fill(value, offset + col1, offset + col2)
    }
  }

  get(row: number, col: number): number {
    let curCol = this._rawCol
    if (this._needRotate) {
      curCol = this._rawRow
      const tmp = row
      row = col
      col = this._rawRow - tmp - 1
    }

    return this._data[row * curCol + col]
  }
}

export {}

if (require.main === module) {
  class SubrectangleQueries {
    private readonly _manager: RangeAssignPointGet2DNaive

    constructor(rectangle: number[][]) {
      this._manager = new RangeAssignPointGet2DNaive(rectangle)
    }

    /**
     * 将左上角为`[row1, col1]`,右下角为`[row2, col2]`的子矩形中的所有元素更新为`newValue`.
     */
    updateSubrectangle(
      row1: number,
      col1: number,
      row2: number,
      col2: number,
      newValue: number
    ): void {
      this._manager.update(row1, row2 + 1, col1, col2 + 1, newValue)
    }

    getValue(row: number, col: number): number {
      return this._manager.get(row, col)
    }
  }

  /**
   * Your SubrectangleQueries object will be instantiated and called as such:
   * var obj = new SubrectangleQueries(rectangle)
   * obj.updateSubrectangle(row1,col1,row2,col2,newValue)
   * var param_2 = obj.getValue(row,col)
   */

  const N = 1e5
  const sqrt = 1 + (Math.sqrt(N) | 0)
  const ROW = sqrt
  const COL = sqrt
  const tree = new SubrectangleQueries(
    Array(ROW)
      .fill(0)
      .map(() => Array(COL).fill(0))
  )
  console.time('update')
  for (let i = 0; i < 1e5; ++i) {
    tree.updateSubrectangle(0, 0, ROW - 1, COL - 1, i)
    tree.getValue(0, 0)
  }
  console.timeEnd('update') // update: 1.583s
}
