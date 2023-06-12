/* eslint-disable prefer-destructuring */

// 适用于 ROW*COL<=1e5 稠密的二维矩阵
// !如果ROW/COL很大,需要离散化
// Naive作为内层的树，每次更新是O(col*logrow)的
// !但是行列中有一个不会超过根号，所以可以做一些trick让这里col最坏是sqrt(1e5)的
// !每次更新复杂度为O(sqrt(row*col)*log(row))
// !1000*1000的矩阵更新2e5次：1000 * Math.ceil(Math.log2(1000)) * 2e5 => 2e9次,实际耗时1.1s左右
// 因此对于js只使用遍历+赋值，计算复杂度时可以除以5到20

const INF = 2e15

interface IRangeUpdatePointGet1D<E, Id> {
  update(start: number, end: number, lazy: Id): void
  get(index: number): E
}

/**
 * 二维区间更新，单点查询的线段树(树套树).
 */
class SegTree2DRangeUpdatePointGet<E = number, Id = number> {
  /**
   * 存储内层的"树"结构.
   */
  private readonly _seg: IRangeUpdatePointGet1D<E, Id>[]

  /**
   * 合并两个内层"树"的结果.
   */
  private readonly _mergeE: (a: E, b: E) => E

  /**
   * 初始化内层"树"的函数.
   */
  private readonly _init1D: () => IRangeUpdatePointGet1D<E, Id>

  private readonly _size: number
  private readonly _needTranspose: boolean

  /**
   * @param row 行数.对时间复杂度贡献为`O(log(row))`.
   * @param col 列数.内部树的大小为.列数越小,对内部树的时间复杂度要求越低.
   * @param mergeE 合并两个内层"树"的结果.
   * @param createRangeUpdatePointGet1D 初始化内层"树"的函数.入参为内层"树"的大小.
   */
  constructor(
    row: number,
    col: number,
    mergeE: (a: E, b: E) => E,
    createRangeUpdatePointGet1D: (n: number) => IRangeUpdatePointGet1D<E, Id>
  ) {
    let size = 1
    while (size < row) size <<= 1
    this._size = size
    this._seg = Array(2 * size - 1)
    this._mergeE = mergeE
    this._init1D = () => createRangeUpdatePointGet1D(col)
    this._needTranspose = row < col
  }

  update(row1: number, row2: number, col1: number, col2: number, lazy: Id): void {
    this._update(row1, row2, col1, col2, lazy, 0, 0, this._size)
  }

  get(row: number, col: number): E {
    row += this._size - 1
    if (!this._seg[row]) this._seg[row] = this._init1D()
    let res = this._seg[row].get(col)
    while (row > 0) {
      row = (row - 1) >> 1
      if (!this._seg[row]) this._seg[row] = this._init1D()
      res = this._mergeE(res, this._seg[row].get(col))
    }
    return res
  }

  private _update(
    R: number,
    C: number,
    start: number,
    end: number,
    lazy: Id,
    pos: number,
    r: number,
    c: number
  ): void {
    if (c <= R || C <= r) return
    if (R <= r && c <= C) {
      if (!this._seg[pos]) this._seg[pos] = this._init1D()
      this._seg[pos].update(start, end, lazy)
    } else {
      const mid = (r + c) >>> 1
      this._update(R, C, start, end, lazy, 2 * pos + 1, r, mid)
      this._update(R, C, start, end, lazy, 2 * pos + 2, mid, c)
    }
  }

  private _normalize(
    row1: number,
    col1: number,
    row2: number,
    col2: number
  ): [number, number, number, number] {}
}

export { SegTree2DRangeUpdatePointGet }

if (require.main === module) {
  // https://leetcode.cn/problems/subrectangle-queries/

  /**
   * !区间染色，单点求值.
   */
  class SubrectangleQueries {
    private _seg2d: SegTree2DRangeUpdatePointGet<E, Id>
    private _updateTime = 2

    constructor(rectangle: number[][]) {
      const ROW = rectangle.length
      const COL = rectangle[0].length
      const seg2d = new SegTree2DRangeUpdatePointGet<E, Id>(
        ROW,
        COL,
        (a, b) => {
          if (!a || !b) return a || b
          return a[0] > b[0] ? a : b
        },
        n => new NaiveTree(n)
      )

      for (let i = 0; i < ROW; ++i) {
        for (let j = 0; j < COL; ++j) {
          seg2d.update(i, i + 1, j, j + 1, [1, rectangle[i][j]])
        }
      }

      this._seg2d = seg2d
    }

    updateSubrectangle(
      row1: number,
      col1: number,
      row2: number,
      col2: number,
      newValue: number
    ): void {
      this._seg2d.update(row1, row2 + 1, col1, col2 + 1, [this._updateTime++, newValue])
    }

    getValue(row: number, col: number): number {
      return this._seg2d.get(row, col)![1]
    }
  }

  /**
   * Your SubrectangleQueries object will be instantiated and called as such:
   * var obj = new SubrectangleQueries(rectangle)
   * obj.updateSubrectangle(row1,col1,row2,col2,newValue)
   * var param_2 = obj.getValue(row,col)
   */

  type E = [time: number, value: number]
  type Id = E

  /**
   * 内层"树"的实现.
   * 这里把Id拆成两个类型数组存，节省空间.
   */
  class NaiveTree {
    private readonly _nums: E[]

    constructor(n: number) {
      // this._nums = Array(n)
      // for (let i = 0; i < n; ++i) this._nums[i] = undefined // 1s(手动填充undefined)

      // this._nums = Array(n).fill(0) // 1.78s(填充0)

      // this._nums = Array(n).fill(undefined) // 1.8s(不填充)

      this._nums = Array(n)
      for (let i = 0; i < n; ++i) this._nums[i] = [] // 320ms(手动填充[])
    }

    update(start: number, end: number, lazy: Id): void {
      for (let i = start; i < end; ++i) this._nums[i] = lazy
    }

    get(index: number): E {
      return this._nums[index]
    }
  }

  //   SubrectangleQueries subrectangleQueries = new SubrectangleQueries([[1,1,1],[2,2,2],[3,3,3]]);
  // subrectangleQueries.getValue(0, 0); // 返回 1
  // subrectangleQueries.updateSubrectangle(0, 0, 2, 2, 100);
  // subrectangleQueries.getValue(0, 0); // 返回 100
  // subrectangleQueries.getValue(2, 2); // 返回 100
  // subrectangleQueries.updateSubrectangle(1, 1, 2, 2, 20);
  // subrectangleQueries.getValue(2, 2); // 返回 20

  // 来源：力扣（LeetCode）
  // 链接：https://leetcode.cn/problems/subrectangle-queries
  // 著作权归领扣网络所有。商业转载请联系官方授权，非商业转载请注明出处。
  const ROW = 600
  const COL = 600
  const tree = new SubrectangleQueries(
    Array(ROW)
      .fill(0)
      .map(() => Array(COL).fill(0))
  )
  console.time('update')
  for (let i = 0; i < 1e5; ++i) {
    tree.updateSubrectangle(0, 0, ROW - 1, COL - 1, i)
  }
  console.timeEnd('update') // update: 228.916ms
}
