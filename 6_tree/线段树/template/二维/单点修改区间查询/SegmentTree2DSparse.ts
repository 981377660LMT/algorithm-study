/* eslint-disable max-len */
// 二维线段树区间修改单点查询
// Query(lx, rx, ly, ry) 查询区间[lx, rx) * [ly, ry)的值.
// Update(x, y, val) 将点(x, y)的值加上(op) val.
// !每次修改/查询 O(lognlogn)

/**
 * 注意运算需要满足交换律.
 */
interface IOptions<E> {
  xs: number[]
  ys: number[]
  ws?: E[]

  /**
   * 是否离散化x.
   * - 为 true 时对x维度二分离散化,然后用离散化后的值作为下标.
   * - 为 false 时不对x维度二分离散化,而是直接用x的值作为下标(内部给x一个偏移量minX),此时x维度数组长度为最大值减最小值.
   *
   * 默认为 true.
   */
  discretizeX?: boolean

  e: () => E
  op: (a: E, b: E) => E
}

class SegmentTree2DSparse<E> {
  private static _bisectLeft(arr: ArrayLike<number>, x: number, left: number, right: number): number {
    while (left <= right) {
      const mid = (left + right) >> 1
      if (arr[mid] < x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _e: () => E
  private readonly _op: (a: E, b: E) => E
  private readonly _discretizeX: boolean
  private _keyX!: number[]
  private _keyY!: number[]
  private _indptr!: Uint32Array
  private _data!: E[]
  private _minX = 0
  private _n = 0

  constructor(options: IOptions<E>) {
    const { xs, ys, e, op, ws, discretizeX = true } = options
    this._e = e
    this._op = op
    let leaves = ws
    this._discretizeX = discretizeX
    if (!leaves) {
      leaves = Array(xs.length)
      for (let i = 0; i < leaves.length; i++) leaves[i] = e()
    }
    this._build(xs, ys, leaves)
  }

  /**
   * 点 (x,y) 与 val 作用.
   */
  update(x: number, y: number, val: E): void {
    let i = this._xtoi(x) + this._n
    while (i > 0) {
      this._update(i, y, val)
      i >>>= 1
    }
  }

  /**
   * [lx,rx) * [ly,ry) 的和.
   */
  query(lx: number, rx: number, ly: number, ry: number): E {
    if (lx >= rx || ly >= ry) return this._e()
    let left = this._xtoi(lx) + this._n
    let right = this._xtoi(rx) + this._n
    let res = this._e()
    while (left < right) {
      if (left & 1) {
        res = this._op(res, this._prodI(left, ly, ry))
        left++
      }
      if (right & 1) {
        right--
        res = this._op(this._prodI(right, ly, ry), res)
      }
      left >>>= 1
      right >>>= 1
    }
    return res
  }

  private _build(xs: number[], ys: number[], ws: E[]): void {
    if (xs.length !== ys.length || xs.length !== ws.length) {
      throw new Error('Lengths of X, Y, and wt must be equal.')
    }

    if (this._discretizeX) {
      this._keyX = [...new Set(xs)].sort((a, b) => a - b)
      this._n = this._keyX.length
    } else {
      if (xs.length) {
        let min = 0
        let max = 0
        for (let i = 0; i < xs.length; i++) {
          min = Math.min(min, xs[i])
          max = Math.max(max, xs[i])
        }
        this._minX = min
        this._n = max - min + 1
      }
      this._keyX = Array(this._n)
      for (let i = 0; i < this._n; i++) {
        this._keyX[i] = i + this._minX
      }
    }

    const n = this._n
    const keyYRaw: number[][] = Array(n + n)
    const datRaw: E[][] = Array(n + n)
    // !faster than `map`
    for (let i = 0; i < n + n; i++) {
      keyYRaw[i] = []
      datRaw[i] = []
    }
    const order = Array(ys.length)
    for (let i = 0; i < ys.length; i++) order[i] = i
    order.sort((a, b) => ys[a] - ys[b])

    for (let i = 0; i < order.length; i++) {
      const v = order[i]
      let ix = this._xtoi(xs[v]) + n
      const y = ys[v]
      while (ix > 0) {
        const KY = keyYRaw[ix]
        const tmp = datRaw[ix]
        if (!KY.length || KY[KY.length - 1] < y) {
          keyYRaw[ix].push(y)
          tmp.push(ws[v])
        } else {
          tmp[tmp.length - 1] = this._op(tmp[tmp.length - 1], ws[v])
        }
        ix >>>= 1
      }
    }

    this._indptr = new Uint32Array(n + n + 1)
    for (let i = 0; i < n + n; i++) {
      this._indptr[i + 1] = this._indptr[i] + keyYRaw[i].length
    }

    const fullN = this._indptr[n + n]
    this._keyY = Array(fullN).fill(0)
    this._data = Array(2 * fullN)
    for (let i = 0; i < this._data.length; i++) this._data[i] = this._e()

    for (let i = 0; i < n + n; i++) {
      const ptr = this._indptr[i]
      const off = 2 * ptr
      const diff = this._indptr[i + 1] - ptr
      const keyY = keyYRaw[i]
      const dat = datRaw[i]
      for (let j = 0; j < diff; j++) {
        this._keyY[ptr + j] = keyY[j]
        this._data[off + diff + j] = dat[j]
      }
      for (let j = diff - 1; j >= 1; j--) {
        this._data[off + j] = this._op(this._data[off + 2 * j], this._data[off + 2 * j + 1])
      }
    }
  }

  private _xtoi(x: number): number {
    if (this._discretizeX) return SegmentTree2DSparse._bisectLeft(this._keyX, x, 0, this._n - 1)
    const res = x - this._minX
    if (res < 0) return 0
    if (res >= this._n) return this._n
    return res
  }

  private _update(i: number, y: number, val: E): void {
    const lid = this._indptr[i]
    const diff = this._indptr[i + 1] - this._indptr[i]
    let j = SegmentTree2DSparse._bisectLeft(this._keyY, y, lid, lid + diff - 1) - lid
    const offset = 2 * lid
    j += diff
    while (j > 0) {
      this._data[offset + j] = this._op(this._data[offset + j], val)
      j >>= 1
    }
  }

  private _prodI(i: number, ly: number, ry: number): E {
    const lid = this._indptr[i]
    const diff = this._indptr[i + 1] - this._indptr[i]
    let left = SegmentTree2DSparse._bisectLeft(this._keyY, ly, lid, lid + diff - 1) - lid + diff
    let right = SegmentTree2DSparse._bisectLeft(this._keyY, ry, lid, lid + diff - 1) - lid + diff
    const offset = 2 * lid
    let val = this._e()
    while (left < right) {
      if (left & 1) {
        val = this._op(val, this._data[offset + left])
        left++
      }
      if (right & 1) {
        right--
        val = this._op(this._data[offset + right], val)
      }
      left >>= 1
      right >>= 1
    }
    return val
  }
}

if (require.main === module) {
  class NumMatrix {
    private readonly matrix: number[][]
    private readonly ROW: number
    private readonly COL: number
    private readonly tree: SegmentTree2DSparse<number>

    constructor(matrix: number[][]) {
      this.matrix = matrix
      this.ROW = matrix.length
      this.COL = matrix[0].length
      const xs = Array(this.ROW * this.COL)
      const ys = Array(this.ROW * this.COL)
      const ws = Array(this.ROW * this.COL)
      for (let r = 0; r < this.ROW; r++) {
        for (let c = 0; c < this.COL; c++) {
          xs[r * this.COL + c] = r
          ys[r * this.COL + c] = c
          ws[r * this.COL + c] = matrix[r][c]
        }
      }
      this.tree = new SegmentTree2DSparse({
        xs,
        ys,
        ws,
        e: () => 0,
        op: (a, b) => a + b
      })
    }

    update(row: number, col: number, val: number): void {
      const pre = this.matrix[row][col]
      const diff = val - pre
      this.matrix[row][col] = val
      this.tree.update(row, col, diff)
    }

    sumRegion(row1: number, col1: number, row2: number, col2: number): number {
      return this.tree.query(row1, row2 + 1, col1, col2 + 1)
    }
  }
}

export { SegmentTree2DSparse }

if (require.main === module) {
  // test performance
  const ROW = 400
  const COL = 400
  const points: [x: number, y: number][] = []
  for (let i = 0; i < ROW; i++) {
    for (let j = 0; j < COL; j++) {
      points.push([~~(Math.random() * 1e9) - 8e5, ~~(Math.random() * 1e9) - 8e5])
    }
  }

  console.time('1e5')
  const xs = points.map(([x]) => x)
  const ys = points.map(([y]) => y)
  const seg = new SegmentTree2DSparse({
    xs,
    ys,
    e: () => 0,
    op: (a, b) => a + b
  })
  for (let i = 0; i < 500000; i++) {
    seg.update(~~(Math.random() * ROW), ~~(Math.random() * COL), ~~(Math.random() * 1e9) - 8e5)
    seg.query(~~(Math.random() * ROW), ~~(Math.random() * ROW), ~~(Math.random() * COL), ~~(Math.random() * COL))
  }
  console.timeEnd('1e5')
}
