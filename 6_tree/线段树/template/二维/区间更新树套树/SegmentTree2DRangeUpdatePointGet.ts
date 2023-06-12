/* eslint-disable prefer-destructuring */

// 适用于 ROW*COL<=1e5 稠密的二维矩阵
// !如果ROW/COL很大,需要离散化
// Naive作为内层的树，每次更新是O(col*logrow)的
// !但是行列中有一个不会超过根号，所以可以做一些trick让这里col最坏是sqrt(1e5)的
// !每次更新复杂度为O(sqrt(row*col)*log(row))
// !1000*1000的矩阵更新2e5次：1000 * Math.ceil(Math.log2(1000)) * 2e5 => 2e9次,实际耗时1.1s左右
// 因此对于js只使用遍历+赋值，计算复杂度时可以除以5到20
// https://leetcode.cn/problems/subrectangle-queries/solution/typescript-shi-yong-typedarray-you-hua-d-cdgh/
// 如果可以执行定时任务,可以保存action+倒序查询 ↓
// https://leetcode.cn/problems/subrectangle-queries/solution/bu-bao-li-geng-xin-ju-zhen-yuan-su-de-jie-fa-by-li/
// 一般的二维线段树只能解决单点更新区间查询的情形，
// 这种区间更新单点查询一般是用树套树实现的(线段树维护幺半群，如果是阿贝尔群的话用二维树状数组就可以了)，
// 不过空间占用会很大。因为这里内层的树大小不超过sqrt(row*col)，用朴素的实现和内层线段树的实现差别比较小，
// 而且朴素数组的更新对缓存比较友好，遍历的常数极小，实际运行时间可以做到和log级别的数据结构差不多

const INF = 2e15

interface IRangeUpdatePointGet1D<E, Id> {
  update(start: number, end: number, lazy: Id): void
  get(index: number): E
}

/**
 * 二维区间更新，单点查询的线段树(树套树).
 */
class SegmentTree2DRangeUpdatePointGet<E = number, Id = number> {
  /**
   * 存储内层的"树"结构.
   */
  private readonly _seg: IRangeUpdatePointGet1D<E, Id>[]

  /**
   * 合并两个内层"树"的结果.
   */
  private readonly _mergeRow: (a: E, b: E) => E

  /**
   * 初始化内层"树"的函数.
   */
  private readonly _init1D: () => IRangeUpdatePointGet1D<E, Id>

  /**
   * 当列数超过行数时,需要对矩阵进行旋转,将列数控制在根号以下.
   */
  private readonly _needRotate: boolean

  /**
   * 原始矩阵的行数(未经旋转).
   */
  private readonly _rawRow: number

  private readonly _size: number

  /**
   * @param row 行数.对时间复杂度贡献为`O(log(row))`.
   * @param col 列数.内部树的大小.列数越小,对内部树的时间复杂度要求越低.
   * @param createRangeUpdatePointGet1D 初始化内层"树"的函数.入参为内层"树"的大小.
   * @param mergeRow 合并两个内层"树"的结果.
   */
  constructor(
    row: number,
    col: number,
    createRangeUpdatePointGet1D: (n: number) => IRangeUpdatePointGet1D<E, Id>,
    mergeRow: (a: E, b: E) => E
  ) {
    this._rawRow = row
    this._needRotate = row < col
    if (this._needRotate) {
      row ^= col
      col ^= row
      row ^= col
    }

    let size = 1
    while (size < row) size <<= 1
    this._seg = Array(2 * size - 1)
    this._mergeRow = mergeRow
    this._init1D = () => createRangeUpdatePointGet1D(col)
    this._size = size
  }

  /**
   * 将`[row1,row2)`x`[col1,col2)`的区间值与`lazy`作用.
   */
  update(row1: number, row2: number, col1: number, col2: number, lazy: Id): void {
    if (this._needRotate) {
      const tmp1 = row1
      const tmp2 = row2
      row1 = col1
      row2 = col2
      col1 = this._rawRow - tmp2
      col2 = this._rawRow - tmp1
    }

    this._update(row1, row2, col1, col2, lazy, 0, 0, this._size)
  }

  get(row: number, col: number): E {
    if (this._needRotate) {
      const tmp = row
      row = col
      col = this._rawRow - tmp - 1
    }

    row += this._size - 1
    if (!this._seg[row]) this._seg[row] = this._init1D()
    let res = this._seg[row].get(col)
    while (row > 0) {
      row = (row - 1) >> 1
      if (!this._seg[row]) this._seg[row] = this._init1D()
      res = this._mergeRow(res, this._seg[row].get(col))
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
}

export { SegmentTree2DRangeUpdatePointGet }

if (require.main === module) {
  // https://leetcode.cn/problems/subrectangle-queries/

  /**
   * !区间染色，单点求值的线段树.
   */
  class SubrectangleQueries<V> {
    private _seg2d: SegmentTree2DRangeUpdatePointGet<E, Id>
    private _updateTime = 1
    private readonly _id: Map<V, number> = new Map()
    private readonly _rId: V[] = []

    constructor(rectangle: ArrayLike<ArrayLike<V>>) {
      const row = rectangle.length
      const col = rectangle[0].length
      const seg2d = new SegmentTree2DRangeUpdatePointGet<E, Id>(
        row,
        col,
        n => new NaiveTree(n),
        (a, b) => (a[0] > b[0] ? a : b)
      )
      this._seg2d = seg2d

      for (let i = 0; i < row; ++i) {
        const cache = rectangle[i]
        for (let j = 0; j < col; ++j) {
          this.updateSubrectangle(i, j, i, j, cache[j])
        }
      }
    }

    /**
     * 将左上角为`[row1, col1]`,右下角为`[row2, col2]`的子矩形中的所有元素更新为`newValue`.
     */
    updateSubrectangle(row1: number, col1: number, row2: number, col2: number, newValue: V): void {
      const id = this._getIdByValue(newValue)
      this._seg2d.update(row1, row2 + 1, col1, col2 + 1, [this._updateTime++, id])
    }

    getValue(row: number, col: number): V {
      const id = this._seg2d.get(row, col)![1]
      return this._getValueById(id)
    }

    private _getIdByValue(v: V): number {
      const res = this._id.get(v)
      if (res !== void 0) return res
      const newId = this._id.size
      this._id.set(v, newId)
      this._rId.push(v)
      return newId
    }

    private _getValueById(id: number): V {
      return this._rId[id]
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
   * 也可以不初始化数组,动态开点.
   */
  class NaiveTree {
    private readonly _time: Int32Array
    private readonly _value: Uint32Array

    constructor(n: number) {
      this._time = new Int32Array(n)
      this._value = new Uint32Array(n)
    }

    update(start: number, end: number, lazy: Id): void {
      this._time.fill(lazy[0], start, end)
      this._value.fill(lazy[1], start, end)
    }

    get(index: number): E {
      return [this._time[index], this._value[index]]
    }
  }

  const ROW = 1000
  const COL = 1000
  const tree = new SubrectangleQueries(
    Array(ROW)
      .fill(0)
      .map(() => Array(COL).fill(0))
  )
  console.time('update')
  for (let i = 0; i < 1e6; ++i) {
    tree.updateSubrectangle(0, 0, ROW - 1, COL - 1, i)
    tree.getValue(0, 0)
  }
  console.timeEnd('update') // update: 1.672s
}
