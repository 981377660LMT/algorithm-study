/* eslint-disable no-inner-declarations */
/* eslint-disable no-useless-constructor */
// https://leetcode.cn/problems/subrectangle-queries/

// ! 注意题目是否可以这样开空间
class SubrectangleQueries {
  private readonly _tree: SegmentTree2D

  constructor(rectangle: number[][]) {
    const [ROW, COL] = [rectangle.length, rectangle[0].length]
    this._tree = new SegmentTree2D(0, 0, ROW - 1, COL - 1)
    for (let r = 0; r < ROW; r++) {
      for (let c = 0; c < COL; c++) {
        this._tree.update(r, c, r, c, rectangle[r][c])
      }
    }
  }

  updateSubrectangle(
    row1: number,
    col1: number,
    row2: number,
    col2: number,
    newValue: number
  ): void {
    this._tree.update(row1, col1, row2, col2, newValue)
  }

  getValue(row: number, col: number): number {
    return this._tree.query(row, col, row, col)
  }
}

/**
 * Your SubrectangleQueries object will be instantiated and called as such:
 * var obj = new SubrectangleQueries(rectangle)
 * obj.updateSubrectangle(row1,col1,row2,col2,newValue)
 * var param_2 = obj.getValue(row,col)
 */
class SegmentTreeNode2D {
  value = 0

  child1?: SegmentTreeNode2D // 左上
  child2?: SegmentTreeNode2D // 右上
  child3?: SegmentTreeNode2D // 左下
  child4?: SegmentTreeNode2D // 右下

  isLazy = false
  lazyValue = 0
}

class SegmentTree2D {
  private readonly root: SegmentTreeNode2D = new SegmentTreeNode2D()

  /**
   * 0 <= row1 <= row2 < ROW
   * 0 <= col1 <= col2 < COL
   */
  constructor(
    private readonly ROW1: number,
    private readonly COL1: number,
    private readonly ROW2: number,
    private readonly COL2: number
  ) {}

  update(row1: number, col1: number, row2: number, col2: number, target: number): void {
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
      target
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
    target: number
  ): void {
    if (ROW1 <= row1 && row2 <= ROW2 && COL1 <= col1 && col2 <= COL2) {
      root.value = target * (row2 - row1 + 1) * (col2 - col1 + 1)
      root.lazyValue = target
      root.isLazy = true
      return
    }

    const rowMid = Math.floor((row1 + row2) / 2)
    const colMid = Math.floor((col1 + col2) / 2)

    this.pushDown(root, row1, col1, row2, col2, rowMid, colMid)

    if (ROW1 <= rowMid) {
      if (COL1 <= colMid) {
        this._update(root.child1!, ROW1, COL1, ROW2, COL2, row1, col1, rowMid, colMid, target)
      }
      if (colMid < COL2) {
        this._update(root.child2!, ROW1, COL1, ROW2, COL2, row1, colMid + 1, rowMid, col2, target)
      }
    }

    if (rowMid < ROW2) {
      if (COL1 <= colMid) {
        this._update(root.child3!, ROW1, COL1, ROW2, COL2, rowMid + 1, col1, row2, colMid, target)
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
          target
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

    if (root.isLazy) {
      const target = root.lazyValue

      root.child1!.lazyValue = target
      root.child2!.lazyValue = target
      root.child3!.lazyValue = target
      root.child4!.lazyValue = target

      root.child1.value = target * (rowMid - row1 + 1) * (colMid - col1 + 1)
      root.child2.value = target * (rowMid - row1 + 1) * (col2 - colMid)
      root.child3.value = target * (row2 - rowMid) * (colMid - col1 + 1)
      root.child4.value = target * (row2 - rowMid) * (col2 - colMid)

      root.child1.isLazy = true
      root.child2.isLazy = true
      root.child3.isLazy = true
      root.child4.isLazy = true

      root.isLazy = false
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

export {}
