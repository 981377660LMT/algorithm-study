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
    const seg: number[][] = new Array(tmp)
    const ll: number[][] = new Array(tmp)
    const rr: number[][] = new Array(tmp)
    for (let i = 0; i < tmp; i++) {
      seg[i] = []
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
      while (i < len1 && j < len2) {
        if (seg[a][i] < seg[b][j]) {
          seg[k].push(seg[a][i])
          i++
        } else {
          seg[k].push(seg[b][j])
          j++
        }
      }
      seg[k].push(...seg[a].slice(i))
      seg[k].push(...seg[b].slice(j))

      let tail1 = 0
      let tail2 = 0
      for (let i = 0; i < len; i++) {
        while (tail1 < len1 && seg[a][tail1] < seg[k][i]) {
          tail1++
        }
        while (tail2 < len2 && seg[b][tail2] < seg[k][i]) {
          tail2++
        }
        ll[k][i] = tail1
        rr[k][i] = tail2
      }
      ll[k][len] = len1
      rr[k][len] = len2
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
    return (
      this._query(a, b, this._ll[k][lower], this._ll[k][upper], 2 * k + 1, l, (l + r) >> 1) +
      this._query(a, b, this._rr[k][lower], this._rr[k][upper], 2 * k + 2, (l + r) >> 1, r)
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
