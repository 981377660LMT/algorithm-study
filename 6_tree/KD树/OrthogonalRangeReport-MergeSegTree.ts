// https://kopricky.github.io/code/SegmentTrees/orthogonal_range_report.html

// 正交矩形报告(Orthogonal Range Report)
// 统计矩形内的点, 每次查询 O(logN^2 + 点数)

/**
 * 正交范围报告(Orthogonal Range Report).
 */
class OrthogonalRangeReport {
  private static _lowerBound<E, T>(
    arr: ArrayLike<E>,
    x: T,
    key: (e: E) => T = e => e as unknown as T
  ): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (key(arr[mid]) < x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private static _upperBound<E, T>(
    arr: ArrayLike<E>,
    x: T,
    key: (e: E) => T = e => e as unknown as T
  ): number {
    let left = 0
    let right = arr.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (key(arr[mid]) <= x) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _n: number
  private readonly _xs: number[]
  private readonly _ys: [y: number, i: number][][]

  constructor(points: [x: number, y: number][] | number[][]) {
    const sz = points.length
    let n = 1
    while (n < sz) n <<= 1
    const sorted: [x: number, y: number, i: number][] = Array(sz)
    for (let i = 0; i < sz; i++) {
      sorted[i] = [points[i][0], points[i][1], i]
    }
    sorted.sort((a, b) => a[0] - b[0])
    const xs: number[] = Array(sz)
    const ys: [y: number, i: number][][] = Array(2 * n - 1)

    for (let i = 0; i < sz; i++) {
      // eslint-disable-next-line prefer-destructuring
      xs[i] = sorted[i][0]
      ys[i + n - 1] = [[sorted[i][1], sorted[i][2]]]
    }
    for (let i = sz; i < n; i++) {
      ys[i + n - 1] = []
    }

    for (let i = n - 2; ~i; i--) {
      const nums1 = ys[2 * i + 1]
      const nums2 = ys[2 * i + 2]
      ys[i] = Array(nums1.length + nums2.length)
      let p = 0
      let p1 = 0
      let p2 = 0
      while (p1 < nums1.length && p2 < nums2.length) {
        if (nums1[p1][0] < nums2[p2][0]) {
          ys[i][p++] = nums1[p1++]
        } else {
          ys[i][p++] = nums2[p2++]
        }
      }
      while (p1 < nums1.length) {
        ys[i][p++] = nums1[p1++]
      }
      while (p2 < nums2.length) {
        ys[i][p++] = nums2[p2++]
      }
    }

    this._n = n
    this._xs = xs
    this._ys = ys
  }

  /**
   * 返回[x1,x2)×[y1,y2)的矩形内的点的索引.
   */
  query(x1: number, x2: number, y1: number, y2: number): number[] {
    const lxid = OrthogonalRangeReport._lowerBound(this._xs, x1)
    const rxid = OrthogonalRangeReport._lowerBound(this._xs, x2)
    if (lxid >= rxid) return []
    const report: number[] = []
    this._query(lxid, rxid, y1, y2, report, 0, 0, this._n)
    return report
  }

  private _query(
    lxid: number,
    rxid: number,
    ly: number,
    ry: number,
    report: number[],
    k: number,
    l: number,
    r: number
  ): void {
    if (r <= lxid || rxid <= l) return
    if (lxid <= l && r <= rxid) {
      const ysk = this._ys[k]
      const st = OrthogonalRangeReport._lowerBound(ysk, ly, e => e[0])
      const ed = OrthogonalRangeReport._upperBound(ysk, ry, e => e[0])
      for (let i = st; i < ed; ++i) {
        report.push(ysk[i][1])
      }
    } else {
      this._query(lxid, rxid, ly, ry, report, 2 * k + 1, l, (l + r) >>> 1)
      this._query(lxid, rxid, ly, ry, report, 2 * k + 2, (l + r) >>> 1, r)
    }
  }
}

export { OrthogonalRangeReport }

if (require.main === module) {
  const n = 1e5
  const points = Array(n)
  for (let i = 0; i < n; i++) {
    points[i] = [i, i]
  }
  console.time('ort')
  const ort = new OrthogonalRangeReport(points)
  for (let i = 0; i < 100; i++) {
    ort.query(0, n, 0, n)
  }
  console.timeEnd('ort')

  const points2 = [
    [1, 1],
    [2, 2],
    [3, 3],
    [4, 4],
    [5, 5]
  ]
  const ort2 = new OrthogonalRangeReport(points2)
  console.log(ort2.query(2, 3, 1, 3))
}
