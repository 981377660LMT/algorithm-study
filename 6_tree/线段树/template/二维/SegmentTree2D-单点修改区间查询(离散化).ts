// 二维线段树区间修改单点查询
// Query(lx, rx, ly, ry) 查询区间[lx, rx) * [ly, ry)的值.
// Update(x, y, val) 将点(x, y)的值加上(op) val.
// !每次修改/查询 O(lognlogn)

/**
 * 注意运算需要满足交换律.
 */
interface Options<E> {
  xs: number[]
  ys: number[]
  ws?: E[]
  e: () => E
  op: (a: E, b: E) => E
}

class SegmentTree2DCompress<E> {
  private static _bisectLeft(
    arr: ArrayLike<number>,
    x: number,
    left: number,
    right: number
  ): number {
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
  private _n!: number
  private _keyX!: number[]
  private _keyY!: number[]
  private _indptr!: Uint32Array
  private _data!: E[]

  constructor(options: Options<E>) {
    const { xs, ys, e, op, ws } = options
    this._e = e
    this._op = op
    this._build(
      xs,
      ys,
      ws ??
        Array(xs.length)
          .fill(0)
          .map(() => e())
    )
  }

  private _build(xs: number[], ys: number[], ws: E[]): void {
    if (xs.length !== ys.length || xs.length !== ws.length) {
      throw new Error('Lengths of X, Y, and wt must be equal.')
    }

    this._keyX = [...new Set(xs)].sort((a, b) => a - b)
    this._n = this._keyX.length
    const n = this._n
    const keyYRaw: number[][] = Array(n + n).fill(0)
    const datRaw: E[][] = Array(n + n).fill(0)
    // !faster than `map`
    for (let i = 0; i < n + n; i++) {
      keyYRaw[i] = []
      datRaw[i] = []
    }
    const indices = Array(ys.length).fill(0)
    for (let i = 0; i < ys.length; i++) {
      indices[i] = i
    }
    indices.sort((a, b) => ys[a] - ys[b])

    indices.forEach(i => {
      let ix = this._xtoi(xs[i]) + n
      const y = ys[i]
      while (ix > 0) {
        const KY = keyYRaw[ix]
        const tmp = datRaw[ix]
        if (!KY.length || KY[KY.length - 1] < y) {
          keyYRaw[ix].push(y)
          tmp.push(ws[i])
        } else {
          tmp[tmp.length - 1] = this._op(tmp[tmp.length - 1], ws[i])
        }
        ix >>= 1
      }
    })

    this._indptr = new Uint32Array(n + n + 1)
    for (let i = 0; i < n + n; i++) {
      this._indptr[i + 1] = this._indptr[i] + keyYRaw[i].length
    }

    const fullN = this._indptr[n + n]
    this._keyY = new Array(fullN)
    this._data = Array(2 * fullN).fill(0)
    for (let i = 0; i < this._data.length; i++) {
      this._data[i] = this._e()
    }

    for (let i = 0; i < n + n; i++) {
      const off = 2 * this._indptr[i]
      const diff = this._indptr[i + 1] - this._indptr[i]
      for (let j = 0; j < diff; j++) {
        this._keyY[this._indptr[i] + j] = keyYRaw[i][j]
        this._data[off + diff + j] = datRaw[i][j]
      }
      for (let j = diff - 1; j >= 1; j--) {
        this._data[off + j] = this._op(this._data[off + 2 * j], this._data[off + 2 * j + 1])
      }
    }
  }

  private _xtoi(x: number): number {
    return SegmentTree2DCompress._bisectLeft(this._keyX, x, 0, this._n - 1)
  }

  /**
   * 点 (x,y) 与 val 作用.
   */
  update(x: number, y: number, val: E): void {
    let i = this._xtoi(x) + this._n
    while (i > 0) {
      this._update(i, y, val)
      i >>= 1
    }
  }

  /**
   * [lx,rx) * [ly,ry) 的和.
   */
  query(lx: number, rx: number, ly: number, ry: number): E {
    let L = this._xtoi(lx) + this._n
    let R = this._xtoi(rx) + this._n
    let val = this._e()
    while (L < R) {
      if (L & 1) {
        val = this._op(val, this._prodI(L, ly, ry))
        L++
      }
      if (R & 1) {
        R--
        val = this._op(this._prodI(R, ly, ry), val)
      }
      L >>= 1
      R >>= 1
    }
    return val
  }

  private _update(i: number, y: number, val: E): void {
    const lid = this._indptr[i]
    const diff = this._indptr[i + 1] - this._indptr[i]
    let j = SegmentTree2DCompress._bisectLeft(this._keyY, y, lid, lid + diff - 1) - lid
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
    let left = SegmentTree2DCompress._bisectLeft(this._keyY, ly, lid, lid + diff - 1) - lid + diff
    let right = SegmentTree2DCompress._bisectLeft(this._keyY, ry, lid, lid + diff - 1) - lid + diff
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
    private readonly tree: SegmentTree2DCompress<number>

    constructor(matrix: number[][]) {
      this.matrix = matrix
      this.ROW = matrix.length
      this.COL = matrix[0].length
      const xs = Array(this.ROW * this.COL).fill(0)
      const ys = Array(this.ROW * this.COL).fill(0)
      const ws = Array(this.ROW * this.COL).fill(0)
      for (let r = 0; r < this.ROW; r++) {
        for (let c = 0; c < this.COL; c++) {
          xs[r * this.COL + c] = r
          ys[r * this.COL + c] = c
          ws[r * this.COL + c] = matrix[r][c]
        }
      }
      this.tree = new SegmentTree2DCompress({
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
  const seg = new SegmentTree2DCompress({
    xs,
    ys,
    e: () => 0,
    op: (a, b) => a + b
  })
  for (let i = 0; i < 500000; i++) {
    seg.update(~~(Math.random() * ROW), ~~(Math.random() * COL), ~~(Math.random() * 1e9) - 8e5)
    seg.query(
      ~~(Math.random() * ROW),
      ~~(Math.random() * ROW),
      ~~(Math.random() * COL),
      ~~(Math.random() * COL)
    )
  }
  console.timeEnd('1e5')
}

export { SegmentTree2DCompress }
