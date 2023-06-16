/* eslint-disable no-inner-declarations */

// 分数级联
// 查询区间 [start, end) 中，[floor, ceiling) 范围内的数的个数.

// !适用于负数  O(nlogn) 常数非常大

class SegmentTreeFractionalCascading {
  private static _bisectLeft(array: number[], target: number): number {
    let left = 0
    let right = array.length - 1
    while (left <= right) {
      const mid = (left + right) >> 1
      if (array[mid] < target) {
        left = mid + 1
      } else {
        right = mid - 1
      }
    }
    return left
  }

  private readonly _seg: number[][]
  private readonly _ll: number[][]
  private readonly _rr: number[][]
  private readonly _sz: number

  constructor(array: number[]) {
    const n = array.length
    let sz = 1
    while (sz < n) {
      sz <<= 1
    }

    const tmp = 2 * sz - 1
    const seg: number[][] = Array(tmp)
    const ll: number[][] = Array(tmp)
    const rr: number[][] = Array(tmp)
    for (let i = 0; i < tmp; i++) {
      seg[i] = []
      ll[i] = []
      rr[i] = []
    }

    for (let i = 0; i < n; i++) {
      seg[i + sz - 1].push(array[i])
    }

    for (let k = sz - 2; ~k; k--) {
      const a = 2 * k + 1
      const b = 2 * k + 2
      const len1 = seg[a].length
      const len2 = seg[b].length
      const len = len1 + len2
      ll[k] = Array(len + 1)
      rr[k] = Array(len + 1)

      let i = 0
      let j = 0
      const segA = seg[a]
      const segB = seg[b]
      const segK = seg[k]
      while (i < len1 && j < len2) {
        if (segA[i] < segB[j]) {
          segK.push(segA[i])
          i++
        } else {
          segK.push(segB[j])
          j++
        }
      }

      while (i < len1) {
        segK.push(segA[i])
        i++
      }
      while (j < len2) {
        segK.push(segB[j])
        j++
      }

      let tail1 = 0
      let tail2 = 0
      const llk = ll[k]
      const rrk = rr[k]
      for (let i = 0; i < len; i++) {
        while (tail1 < len1 && segA[tail1] < segK[i]) {
          tail1++
        }
        while (tail2 < len2 && segB[tail2] < segK[i]) {
          tail2++
        }
        llk[i] = tail1
        rrk[i] = tail2
      }
      llk[len] = len1
      rrk[len] = len2
    }

    this._seg = seg
    this._ll = ll
    this._rr = rr
    this._sz = sz
  }

  /**
   * 查询区间 [start, end) 中，[floor, ceiling) 范围内的数的个数.
   */
  query(start: number, end: number, floor: number, ceiling: number): number {
    floor = SegmentTreeFractionalCascading._bisectLeft(this._seg[0], floor)
    ceiling = SegmentTreeFractionalCascading._bisectLeft(this._seg[0], ceiling)
    return this._query(start, end, floor, ceiling, 0, 0, this._sz)
  }

  private _query(
    a: number,
    b: number,
    lower: number,
    upper: number,
    k: number,
    l: number,
    r: number
  ): number {
    if (a >= r || b <= l) {
      return 0
    }
    if (a <= l && r <= b) {
      return upper - lower
    }
    const mid = (l + r) >> 1
    return (
      this._query(a, b, this._ll[k][lower], this._ll[k][upper], 2 * k + 1, l, mid) +
      this._query(a, b, this._rr[k][lower], this._rr[k][upper], 2 * k + 2, mid, r)
    )
  }
}

if (require.main === module) {
  // 9700ms

  function countQuadruplets(nums: number[]) {
    const W = new SegmentTreeFractionalCascading(nums)
    let res = 0
    for (let j = 1; j < nums.length - 2; j++) {
      for (let k = j + 1; k < nums.length - 1; k++) {
        if (nums[k] < nums[j]) {
          const left = W.query(0, j, 0, nums[k])
          const right = W.query(k + 1, nums.length, nums[j] + 1, 1e9)
          res += left * right
        }
      }
    }
    return res
  }
}

export { SegmentTreeFractionalCascading }
