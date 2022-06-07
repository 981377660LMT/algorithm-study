// 308. 二维区域和检索 - 可变
// https://leetcode.cn/problems/range-sum-query-2d-mutable/

class SegmentTreeNode2D {
  value = 0

  child1?: SegmentTreeNode2D // 左上
  child2?: SegmentTreeNode2D // 右上
  child3?: SegmentTreeNode2D // 左下
  child4?: SegmentTreeNode2D // 右下

  // isLazy = false
  lazyValue = 0 // !叠加可以不需要isLazy标记 直接通过lazyValue来判断
}

class SegmentTree2D {
  private readonly root: SegmentTreeNode2D = new SegmentTreeNode2D()

  constructor(
    private readonly ROW1: number,
    private readonly COL1: number,
    private readonly ROW2: number,
    private readonly COL2: number
  ) {}

  update(row1: number, col1: number, row2: number, col2: number, delta: number): void {
    this.checkRange(row1, col1, row2, col2)
    this._update(
      this.root,
      row1,
      col1,
      row2,
      col2,
      this.ROW1,
      this.COL1,
      this.ROW2,
      this.COL2,
      delta
    )
  }

  query(row1: number, col1: number, row2: number, col2: number): number {
    this.checkRange(row1, col1, row2, col2)
    return this._query(
      this.root,
      row1,
      col1,
      row2,
      col2,
      this.ROW1,
      this.COL1,
      this.ROW2,
      this.COL2
    )
  }

  queryAll(): number {
    return this.root.value
  }

  private _query(
    root: SegmentTreeNode2D,
    ROW1: number,
    COL1: number,
    ROW2: number,
    COL2: number,
    row1: number,
    col1: number,
    row2: number,
    col2: number
  ): number {
    if (ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2) {
      return root.value
    }

    const rowMid = Math.floor((row1 + row2) / 2)
    const colMid = Math.floor((col1 + col2) / 2)
    this.pushDown(root, row1, col1, row2, col2, rowMid, colMid)

    let res = 0

    if (ROW1 <= rowMid) {
      if (COL1 <= colMid) {
        res += this._query(root.child1!, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid)
      }
      if (colMid < COL2) {
        res += this._query(root.child2!, ROW1, COL1, ROW2, COL2, row1, colMid + 1, rowMid, col2)
      }
    }

    if (rowMid < ROW2) {
      if (COL1 <= colMid) {
        res += this._query(root.child3!, ROW1, COL1, ROW2, COL2, rowMid + 1, col1, row2, colMid)
      }
      if (colMid < COL2) {
        res += this._query(root.child4!, ROW1, COL1, ROW2, COL2, rowMid + 1, colMid + 1, row2, col2)
      }
    }

    return res
  }

  private _update(
    root: SegmentTreeNode2D,
    ROW1: number,
    COL1: number,
    ROW2: number,
    COL2: number,
    row1: number,
    col1: number,
    row2: number,
    col2: number,
    delta: number
  ): void {
    if (ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2) {
      root.value += delta * (row2 - row1 + 1) * (col2 - col1 + 1)
      root.lazyValue += delta
      return
    }

    const rowMid = Math.floor((row1 + row2) / 2)
    const colMid = Math.floor((col1 + col2) / 2)

    this.pushDown(root, row1, col1, row2, col2, rowMid, colMid)

    if (ROW1 <= rowMid) {
      if (COL1 <= colMid) {
        this._update(root.child1!, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid, delta)
      }
      if (colMid < COL2) {
        this._update(root.child2!, ROW1, COL1, ROW2, COL2, row1, colMid + 1, rowMid, col2, delta)
      }
    }

    if (rowMid < ROW2) {
      if (COL1 <= colMid) {
        this._update(root.child3!, ROW1, COL1, ROW2, COL2, rowMid + 1, col1, row2, colMid, delta)
      }
      if (colMid < COL2) {
        this._update(
          root.child4!,
          ROW1,
          COL1,
          ROW2,
          COL2,
          rowMid + 1,
          colMid + 1,
          row2,
          col2,
          delta
        )
      }
    }

    this.pushUp(root)
  }

  private pushUp(root: SegmentTreeNode2D): void {
    root.value = root.child1!.value + root.child2!.value + root.child3!.value + root.child4!.value
  }

  private pushDown(
    root: SegmentTreeNode2D,
    row1: number,
    col1: number,
    row2: number,
    col2: number,
    rowMid: number,
    colMid: number
  ): void {
    !root.child1 && (root.child1 = new SegmentTreeNode2D())
    !root.child2 && (root.child2 = new SegmentTreeNode2D())
    !root.child3 && (root.child3 = new SegmentTreeNode2D())
    !root.child4 && (root.child4 = new SegmentTreeNode2D())

    if (root.lazyValue) {
      const delta = root.lazyValue

      root.child1.lazyValue += delta
      root.child2.lazyValue += delta
      root.child3.lazyValue += delta
      root.child4.lazyValue += delta

      root.child1.value += delta * (rowMid - row1 + 1) * (colMid - col1 + 1)
      root.child2.value += delta * (rowMid - row1 + 1) * (col2 - colMid)
      root.child3.value += delta * (row2 - rowMid) * (colMid - col1 + 1)
      root.child4.value += delta * (row2 - rowMid) * (col2 - colMid)

      root.lazyValue = 0
    }
  }

  private checkRange(row1: number, col1: number, row2: number, col2: number): void {
    if (
      this.ROW1 <= row1 &&
      row1 <= row2 &&
      row2 <= this.ROW2 &&
      this.COL1 <= col1 &&
      col1 <= col2 &&
      col2 <= this.COL2
    ) {
      return
    }

    throw new RangeError(
      `[${row1}, ${col1}, ${row2}, ${col2}] out of range: [${this.ROW1}, ${this.COL1}, ${this.ROW2}, ${this.COL2}]`
    )
  }
}

if (require.main === module) {
  class NumMatrix {
    private readonly matrix: number[][]
    private readonly ROW: number
    private readonly COL: number
    private readonly tree: SegmentTree2D

    constructor(matrix: number[][]) {
      this.matrix = matrix
      this.ROW = matrix.length
      this.COL = matrix[0].length
      this.tree = new SegmentTree2D(0, 0, this.ROW, this.COL)
      for (let r = 0; r < this.ROW; r++) {
        for (let c = 0; c < this.COL; c++) {
          this.tree.update(r, c, r, c, matrix[r][c])
        }
      }
    }

    update(row: number, col: number, val: number): void {
      const delta = val - this.matrix[row][col]
      this.matrix[row][col] = val
      this.tree.update(row, col, row, col, delta)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this.tree.query(row1, col1, row2, col2)
    }
  }
}

export {}
