/* eslint-disable no-inner-declarations */
/* eslint-disable prefer-destructuring */
// https://kopricky.github.io/code/SegmentTrees/merge_segtree.html
// 静态二维矩形区间计数(RectangleSum-MergeSegtree)
// O(nlogn)空间构建 O(logn^2)查询
// 虽然Fractional Cascading (分数级联)理论可以O(logn)查询 但是更慢

/**
 * @alias MergeSegtree
 */
class RectangleSum {
  private static _lowerBound(arr: ArrayLike<number>, x: number): number {
    let left = 0
    let right = arr.length - 1
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

  private _n!: number
  private _xs!: number[]
  private _ys!: number[][]
  private _sum!: number[][]
  private _rawX: number[] = []
  private _rawY: number[] = []
  private _rawWeight: number[] = []
  private _hasBuilt = false

  addPoint(x: number, y: number, weight: number): void {
    this._rawX.push(x)
    this._rawY.push(y)
    this._rawWeight.push(weight)
  }

  /**
   * [x1, x2) * [y1, y2) 的矩形区域内的权值和.
   */
  query(x1: number, y1: number, x2: number, y2: number): number {
    if (!this._hasBuilt) {
      this._build()
      this._hasBuilt = true
    }
    const lxid = RectangleSum._lowerBound(this._xs, x1)
    const rxid = RectangleSum._lowerBound(this._xs, x2)
    if (lxid >= rxid) return 0
    return this._query(lxid, rxid, y1, y2)
  }

  private _build() {
    if (this._hasBuilt) return
    const rx = this._rawX
    const ry = this._rawY

    const weights = this._rawWeight
    const sz = rx.length
    const xs = Array<number>(sz)
    let n = 1
    while (n < sz) n <<= 1

    const sorted: [x: number, i: number][] = Array(sz)
    for (let i = 0; i < sz; i++) {
      sorted[i] = [rx[i], i]
    }
    sorted.sort((a, b) => a[0] - b[0])

    const ys: number[][] = Array(2 * n)
    const sum: number[][] = Array(2 * n)
    const hoge: [y: number, i: number][][] = Array(2 * n)

    for (let i = 0; i < sz; i++) {
      const a = sorted[i][0]
      const b = sorted[i][1]
      xs[i] = a
      ys[i + n] = [ry[b]]
      hoge[i + n] = [[ry[b], b]]
      sum[i + n] = [0, weights[b]]
    }
    for (let i = sz + n; i < 2 * n; i++) {
      hoge[i] = []
    }

    for (let i = n - 1; i >= 1; i--) {
      const child1 = hoge[i << 1]
      const child2 = hoge[(i << 1) | 1]
      hoge[i] = Array(child1.length + child2.length)
      const cur = hoge[i]
      let p = 0
      let q = 0
      let r = 0
      while (p < child1.length && q < child2.length) {
        if (child1[p][0] < child2[q][0]) {
          cur[r++] = child1[p++]
        } else {
          cur[r++] = child2[q++]
        }
      }
      while (p < child1.length) {
        cur[r++] = child1[p++]
      }
      while (q < child2.length) {
        cur[r++] = child2[q++]
      }

      ys[i] = Array(cur.length)
      sum[i] = Array(cur.length + 1)
      sum[i][0] = 0
      for (let j = 0; j < cur.length; j++) {
        ys[i][j] = cur[j][0]
        sum[i][j + 1] = sum[i][j] + weights[cur[j][1]]
      }
    }

    this._n = n
    this._xs = xs
    this._ys = ys
    this._sum = sum
    this._rawX = []
    this._rawY = []
    this._rawWeight = []

    this._hasBuilt = true
  }

  private _query(a: number, b: number, ly: number, ry: number): number {
    let res1 = 0
    let res2 = 0
    a += this._n
    b += this._n
    while (a ^ b) {
      if (a & 1) {
        const c = RectangleSum._lowerBound(this._ys[a], ly)
        const d = RectangleSum._lowerBound(this._ys[a], ry)
        res1 += this._sum[a][d] - this._sum[a][c]
        a++
      }
      if (b & 1) {
        b--
        const c = RectangleSum._lowerBound(this._ys[b], ly)
        const d = RectangleSum._lowerBound(this._ys[b], ry)
        res2 += this._sum[b][d] - this._sum[b][c]
      }
      a >>= 1
      b >>= 1
    }
    return res1 + res2
  }
}

export { RectangleSum }
if (require.main === module) {
  // test 1e5

  console.time('query')
  const rs = new RectangleSum()
  for (let i = 0; i < 1e5; i++) {
    rs.addPoint(i, i, 1)
  }
  for (let i = 0; i < 1e5; i++) {
    console.log(rs.query(0, 0, 1e9, 1e9))
  }
  console.timeEnd('query') // 350ms

  // https://leetcode.cn/problems/count-number-of-rectangles-containing-each-point/
  // 2250. 统计包含每个点的矩形数目
  // 请你返回一个整数数组 count ，长度为 points.length，其中 count[j]是 包含 第 j 个点的矩形数目。
  const INF = 2e15
  function countRectangles(rectangles: number[][], points: number[][]): number[] {
    const R = new RectangleSum()
    rectangles.forEach(([x, y]) => {
      R.addPoint(x, y, 1)
    })
    return points.map(([x, y]) => R.query(x, y, INF, INF))
  }
}
