/* eslint-disable no-param-reassign */
/* eslint-disable no-shadow */

// https://nyaannyaan.github.io/library/data-structure-2d/2d-segment-tree.hpp
// !二维线段树:单点修改/区间查询 (更新方式为覆盖)

type SegmentTree2D<E> = {
  /**
   * 在 {@link SegmentTree2D.build} 之前调用，设置初始值.
   *
   * 0 <= row < ROW, 0 <= col < COL.
   */
  set: (row: number, col: number, value: E) => void

  /**
   * 如果调用了 {@link SegmentTree2D.set} 初始化，则需要调用此方法构建树.
   */
  build: () => void

  /** 0 <= row < ROW, 0 <= col < COL. */
  get: (row: number, col: number) => E

  /** 0 <= row < ROW, 0 <= col < COL. */
  update: (row: number, col: number, target: E) => void

  /**
   * 查询闭区间 [row1, row2] x [col1, col2] 的区间值.
   *
   * 0 <= row1 <= row2 < ROW.
   * 0 <= col1 <= col2 < COL.
   */
  query: (row1: number, col1: number, row2: number, col2: number) => E
}

function useSegmentTree2D<E = number>(
  row: number,
  col: number,
  e: () => E,
  op: (a: E, b: E) => E
): SegmentTree2D<E> {
  const _row = 1 << (32 - Math.clz32(row - 1))
  const _col = 1 << (32 - Math.clz32(col - 1))
  const _tree = Array.from({ length: (_row * _col) << 2 }, () => e())
  const _id = (r: number, c: number) => ((r * _col) << 1) + c

  function set(row: number, col: number, value: E): void {
    _tree[_id(row + _row, col + _col)] = value
  }

  function build(): void {
    for (let c = _col; c < _col << 1; c++) {
      for (let r = _row - 1; ~r; r--) {
        _tree[_id(r, c)] = op(_tree[_id(r << 1, c)], _tree[_id((r << 1) | 1, c)])
      }
    }
    for (let r = 0; r < _row << 1; r++) {
      for (let c = _col - 1; ~c; c--) {
        _tree[_id(r, c)] = op(_tree[_id(r, c << 1)], _tree[_id(r, (c << 1) | 1)])
      }
    }
  }

  function get(row: number, col: number): E {
    return _tree[_id(row + _row, col + _col)]
  }

  function update(row: number, col: number, target: E): void {
    let r = row + _row
    let c = col + _col
    _tree[_id(r, c)] = target
    for (let i = r >>> 1; i; i >>>= 1) {
      _tree[_id(i, c)] = op(_tree[_id(i << 1, c)], _tree[_id((i << 1) | 1, c)])
    }
    for (; r; r >>>= 1) {
      for (let j = c >>> 1; j; j >>>= 1) {
        _tree[_id(r, j)] = op(_tree[_id(r, j << 1)], _tree[_id(r, (j << 1) | 1)])
      }
    }
  }

  function query(row1: number, col1: number, row2: number, col2: number): E {
    row2++
    col2++
    if (row1 >= row2 || col1 >= col2) return e()
    let res = e()
    row1 += _row
    row2 += _row
    col1 += _col
    col2 += _col
    for (; row1 < row2; row1 >>>= 1, row2 >>>= 1) {
      if (row1 & 1) {
        res = op(res, _query(row1, col1, col2))
        row1++
      }
      if (row2 & 1) {
        row2--
        res = op(res, _query(row2, col1, col2))
      }
    }
    return res
  }

  return {
    set,
    build,
    get,
    update,
    query
  }

  function _query(r: number, c1: number, c2: number): E {
    let res = e()
    for (; c1 < c2; c1 >>>= 1, c2 >>>= 1) {
      if (c1 & 1) {
        res = op(res, _tree[_id(r, c1)])
        c1++
      }

      if (c2 & 1) {
        c2--
        res = op(res, _tree[_id(r, c2)])
      }
    }
    return res
  }
}

export { useSegmentTree2D }
if (require.main === module) {
  // https://leetcode.cn/problems/range-sum-query-2d-mutable/
  class NumMatrix {
    private readonly _ROW: number
    private readonly _COL: number
    private readonly _tree: SegmentTree2D<number>

    constructor(matrix: number[][]) {
      this._ROW = matrix.length
      this._COL = matrix[0].length
      this._tree = useSegmentTree2D(
        this._ROW,
        this._COL,
        () => 0,
        (a, b) => a + b
      )

      for (let r = 0; r < this._ROW; r++) {
        for (let c = 0; c < this._COL; c++) {
          this._tree.set(r, c, matrix[r][c])
        }
      }

      this._tree.build() // !注意如果set了不要忘记 build
    }

    update(row: number, col: number, val: number): void {
      this._tree.update(row, col, val)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this._tree.query(row1, col1, row2, col2)
    }
  }

  const tree = useSegmentTree2D(
    3,
    3,
    () => 0,
    (a, b) => a + b
  )
}
