/* eslint-disable no-inner-declarations */
/* eslint-disable semi-style */
/* eslint-disable no-param-reassign */

// CountRange: 返回区间 [left, right) 中 范围在 [a, b) 中的 元素的个数.
// Kth: 返回区间 [left, right) 中的 第k小的元素
// KthValueAndSum: 返回区间 [left, right) 中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
// Sum: 返回区间 [left, right) 中第[k1, k2)个元素的 op 的结果
// SumAll: 返回区间 [left, right) 的 op 的结果
// Median: 返回区间 [left, right) 的中位数

// CountRangeSegments: 返回所有区间中 范围在 [a, b) 中的 元素的个数.
// KthSegments: 返回所有区间中的 第k小的元素
// KthValueAndSumSegments: 返回所有区间中的 (第k小的元素, 前k个元素(不包括第k小的元素) 的 op 的结果)
// SumSegments: 返回所有区间中所有数的 op 的结果
// SumAllSegments: 返回所有区间的 op 的结果
// MedianSegments: 返回所有区间的中位数

// MaxRight: 返回使得 check(count,prefixSum) 为 true 的最大 (count, prefixSum) 对.

import assert from 'assert'

/**
 * 常数大约是线段树的级别.
 */
class WaveletMatrixSegments {
  private readonly _log: number
  private readonly _mid: Uint32Array
  private readonly _bv: _BV[]
  private readonly _preSum: number[][]

  /**
   * @param nums 0 <= nums[i] < 2**31 (2e9+)
   * 如果不在这个范围内，需要将所有数(松)离散化.
   * @param log 如果要支持异或,则需要按照异或的值来决定值域.
   * 设为`-1`时表示不使用异或.
   * @param sumData 如果要支持区间和,则需要传入每个位置的权值.
   */
  constructor(nums: Uint32Array, log = -1, sumData: number[] = []) {
    nums = nums.slice()
    sumData = sumData.slice()
    if (log === -1) {
      log = 32 - Math.clz32(Math.max(...nums, 1)) // bitLength
    }

    const n = nums.length
    const mid = new Uint32Array(log)
    const bv = Array(log)
    for (let i = 0; i < log; i++) {
      bv[i] = new _BV(n) // 1-based
    }

    const needMakeSum = sumData.length > 0
    let preSum: number[][] = []
    if (needMakeSum) {
      preSum = Array(log + 1).fill(0)
      for (let i = 0; i <= log; i++) {
        preSum[i] = Array(n + 1).fill(0)
      }
    }

    let a0 = new Uint32Array(n)
    const a1 = new Uint32Array(n)
    let s0 = Array(n).fill(0)
    const s1 = Array(n).fill(0)
    for (let d = log - 1; d >= -1; d--) {
      let p0 = 0
      let p1 = 0
      if (needMakeSum) {
        for (let i = 0; i < n; i++) {
          preSum[d + 1][i + 1] = preSum[d + 1][i] + sumData[i]
        }
      }
      if (d === -1) {
        break
      }

      for (let i = 0; i < n; i++) {
        const f = (nums[i] >> d) & 1
        if (f) {
          if (needMakeSum) {
            s1[p1] = sumData[i]
          }
          bv[d].add(i)
          a1[p1] = nums[i]
          p1++
        } else {
          if (needMakeSum) {
            s0[p0] = sumData[i]
          }
          a0[p0] = nums[i]
          p0++
        }
      }

      mid[d] = p0
      bv[d].build()
      const tmp1 = nums
      nums = a0
      a0 = tmp1
      const tmp2 = sumData
      sumData = s0
      s0 = tmp2
      for (let i = 0; i < p1; i++) {
        nums[p0 + i] = a1[i]
        if (needMakeSum) {
          sumData[p0 + i] = s1[i]
        }
      }
    }

    this._log = log
    this._mid = mid
    this._bv = bv
    this._preSum = preSum
  }

  /**
   * 返回区间 [start, end) 中 范围在 [floor, higher) 中的 元素的个数.
   */
  countRange(start: number, end: number, floor: number, higher: number, xor = 0): number {
    return this._prefixCount(start, end, higher, xor) - this._prefixCount(start, end, floor, xor)
  }

  countRangeSegments(segments: [number, number][], floor: number, higher: number, xor = 0): number {
    let res = 0
    segments.forEach(([start, end]) => {
      res += this.countRange(start, end, floor, higher, xor)
    })
    return res
  }

  /**
   * [start, end) 中的 `第k小的元素, 前k个元素(不包括第k小的元素)的和`.
   * k 从 `0` 开始.
   * 如果 k < 0, 返回 [null, 0]; 如果 k >= right-left, 返回 [null, 区间和].
   */
  kthValueAndSum(
    start: number,
    end: number,
    k: number,
    xor = 0
  ): [res: number, preSum: number] | [res: null, preSum: number] {
    if (k < 0) return [null, 0]
    if (end - start <= k) return [null, this._get(this._log, start, end)]
    let res = 0
    let sum = 0
    let count = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      if (count + kf > k) {
        if (f) {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        } else {
          start = l0
          end = r0
        }
      } else {
        const s = f
          ? this._get(d, this._mid[d] - l0 + start, this._mid[d] - r0 + end)
          : this._get(d, l0, r0)
        count += kf
        res |= 1 << d
        sum += s
        if (f) {
          start = l0
          end = r0
        } else {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        }
      }
    }

    sum += this._get(0, start, start + k - count)
    return [res, sum]
  }

  kthValueAndSumSegments(
    segments: [number, number][],
    k: number,
    xor = 0
  ): [res: number, preSum: number] | [res: null, preSum: number] {
    if (k < 0) return [null, 0]
    let totalLen = 0
    segments.forEach(([start, end]) => {
      totalLen += end - start
    })
    if (totalLen <= k) return [null, this.sumAllSegments(segments)]
    let res = 0
    let sum = 0
    let count = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      let kf = 0
      segments.forEach(([start, end]) => {
        const l0 = this._bv[d].rank(start, 0)
        const r0 = this._bv[d].rank(end, 0)
        kf += f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      })

      if (count + kf > k) {
        for (let i = 0; i < segments.length; i++) {
          const seg = segments[i]
          const L = seg[0]
          const R = seg[1]
          const l0 = this._bv[d].rank(L, 0)
          const r0 = this._bv[d].rank(R, 0)
          if (f) {
            seg[0] += this._mid[d] - l0
            seg[1] += this._mid[d] - r0
          } else {
            seg[0] = l0
            seg[1] = r0
          }
        }
      } else {
        count += kf
        res |= 1 << d
        for (let i = 0; i < segments.length; i++) {
          const seg = segments[i]
          const L = seg[0]
          const R = seg[1]
          const l0 = this._bv[d].rank(L, 0)
          const r0 = this._bv[d].rank(R, 0)
          const s = f
            ? this._get(d, this._mid[d] - l0 + L, this._mid[d] - r0 + R)
            : this._get(d, l0, r0)
          sum += s
          if (f) {
            seg[0] = l0
            seg[1] = r0
          } else {
            seg[0] += this._mid[d] - l0
            seg[1] += this._mid[d] - r0
          }
        }
      }
    }

    segments.forEach(([start, end]) => {
      const t = Math.min(end - start, k - count)
      sum += this._get(0, start, start + t)
      count += t
    })

    return [res, sum]
  }

  kth(start: number, end: number, k: number, xor = 0): number {
    if (k < 0 || k >= end - start) return -1
    let res = 0
    let count = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const c = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      if (count + c > k) {
        if (f) {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        } else {
          start = l0
          end = r0
        }
      } else {
        count += c
        res |= 1 << d
        if (f) {
          start = l0
          end = r0
        } else {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        }
      }
    }
    return res
  }

  kthSegments(segments: [number, number][], k: number, xor = 0): number {
    if (k < 0) return -1
    let totalLen = 0
    segments.forEach(([start, end]) => {
      totalLen += end - start
    })
    if (k >= totalLen) return -1
    let count = 0
    let res = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      let c = 0
      for (let i = 0; i < segments.length; i++) {
        const seg = segments[i]
        const L = seg[0]
        const R = seg[1]
        const l0 = this._bv[d].rank(L, 0)
        const r0 = this._bv[d].rank(R, 0)
        c += f * (R - L - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      }
      if (count + c > k) {
        for (let i = 0; i < segments.length; i++) {
          const seg = segments[i]
          const L = seg[0]
          const R = seg[1]
          const l0 = this._bv[d].rank(L, 0)
          const r0 = this._bv[d].rank(R, 0)
          if (f) {
            seg[0] += this._mid[d] - l0
            seg[1] += this._mid[d] - r0
          } else {
            seg[0] = l0
            seg[1] = r0
          }
        }
      } else {
        count += c
        res |= 1 << d
        for (let i = 0; i < segments.length; i++) {
          const seg = segments[i]
          const L = seg[0]
          const R = seg[1]
          const l0 = this._bv[d].rank(L, 0)
          const r0 = this._bv[d].rank(R, 0)
          if (f) {
            seg[0] = l0
            seg[1] = r0
          } else {
            seg[0] += this._mid[d] - l0
            seg[1] += this._mid[d] - r0
          }
        }
      }
    }

    return res
  }

  /**
   * 区间中位数.
   * @param upper true表示上中位数, false表示下中位数.
   */
  median(upper: boolean, start: number, end: number, xor = 0): number {
    const n = end - start
    const k = upper ? n >>> 1 : (n - 1) >>> 1
    return this.kth(start, end, k, xor)
  }

  medianSegments(upper: boolean, segments: [number, number][], xor = 0): number {
    let n = 0
    segments.forEach(([start, end]) => {
      n += end - start
    })
    const k = upper ? n >>> 1 : (n - 1) >>> 1
    return this.kthSegments(segments, k, xor)
  }

  sum(start: number, end: number, k1: number, k2: number, xor = 0): number {
    return this._prefixSum(start, end, k2, xor) - this._prefixSum(start, end, k1, xor)
  }

  sumSegments(segments: [number, number][], k1: number, k2: number, xor = 0): number {
    return this._prefixSumSegments(segments, k2, xor) - this._prefixSumSegments(segments, k1, xor)
  }

  sumAll(start: number, end: number): number {
    return this._get(this._log, start, end)
  }

  sumAllSegments(segments: [number, number][]): number {
    let res = 0
    segments.forEach(([start, end]) => {
      res += this._get(this._log, start, end)
    })
    return res
  }

  /**
   * 返回使得 check(count,prefixSum) 为 true 的最大 (count, prefixSum) 对.
   * !(即区间内小于 val 的数的个数count和 和 prefixSum 满足 check函数, 找到这样的最大的 (count, prefixSum).
   */
  maxRight(
    start: number,
    end: number,
    predicate: (count: number, preSum: number) => boolean,
    xor = 0
  ): [count: number, preSum: number] {
    const tmp = this._get(this._log, start, end)
    if (predicate(end - start, tmp)) return [end - start, tmp]

    let res = 0
    let count = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      const loSum = f
        ? this._get(d, start + this._mid[d] - l0, end + this._mid[d] - r0)
        : this._get(d, l0, r0)

      const tmp = res + loSum
      if (predicate(count + kf, tmp)) {
        count += kf
        res = tmp
        if (f) {
          start = l0
          end = r0
        } else {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        }
      } else if (f) {
        start += this._mid[d] - l0
        end += this._mid[d] - r0
      } else {
        start = l0
        end = r0
      }
    }

    const k = WaveletMatrixSegments._binarySearch(
      k => predicate(count + k, res + this._get(0, start, start + k)),
      0,
      end - start
    )
    count += k
    res += this._get(0, start, start + k)
    return [count, res]
  }

  private _prefixCount(start: number, end: number, higher: number, xor = 0): number {
    if (higher <= 0) return 0
    if (higher >= 1 << this._log) return end - start
    let count = 0
    for (let d = this._log - 1; ~d; d--) {
      const add = (higher >>> d) & 1
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)

      if (add) {
        count += kf
        if (f) {
          start = l0
          end = r0
        } else {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        }
      } else if (f) {
        start += this._mid[d] - l0
        end += this._mid[d] - r0
      } else {
        start = l0
        end = r0
      }
    }

    return count
  }

  private _prefixSum(start: number, end: number, k: number, xor = 0): number {
    return this.kthValueAndSum(start, end, k, xor)[1]
  }

  private _prefixSumSegments(segments: [number, number][], k: number, xor = 0): number {
    return this.kthValueAndSumSegments(segments, k, xor)[1]
  }

  private static _binarySearch(f: (e: number) => boolean, ok: number, ng: number): number {
    while (Math.abs(ok - ng) > 1) {
      const x = (ok + ng) >>> 1
      if (f(x)) {
        ok = x
      } else {
        ng = x
      }
    }
    return ok
  }

  private _get(d: number, l: number, r: number): number {
    return this._preSum[d][r] - this._preSum[d][l]
  }
}

class _BV {
  private readonly _data: [count: number, sum: number][] = []

  constructor(n: number) {
    this._data = Array(((n + 63) >> 5) + 1)
    for (let i = 0; i < this._data.length; i++) {
      this._data[i] = [0, 0]
    }
  }

  add(i: number): void {
    this._data[i >> 5][0] |= 1 << (i & 31)
  }

  build(): void {
    for (let i = 0; i < this._data.length - 1; i++) {
      this._data[i + 1][1] = this._data[i][1] + onesCount32(this._data[i][0])
    }
  }

  // [0, k) 内 1 的个数.
  rank(k: number, digit: 0 | 1): number {
    const [a, b] = this._data[k >> 5]
    const res = b + onesCount32(a & ((1 << (k & 31)) - 1))
    return digit * res + (digit ^ 1) * (k - res) // digit ? res : k - res
  }
}

function onesCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}

if (require.main === module) {
  const nums = [1, 2, 3, 4, 5, 6, 7, 8, 9, 10]
  const wm = new WaveletMatrixSegments(new Uint32Array(nums), -1, nums)

  assert.strictEqual(wm.countRange(0, 10, 3, 7, 0), 4)
  assert.strictEqual(wm.kth(0, 10, 9, 0), 10)
  assert.deepStrictEqual(wm.kthValueAndSum(0, 10, 10, 0), [null, 55])
  assert.strictEqual(wm.sum(0, 10, 1, 3, 0), 5)
  assert.strictEqual(wm.sumAll(0, 10), 55)
  assert.strictEqual(wm.median(false, 0, 5, 0), 3)

  assert.strictEqual(
    wm.countRangeSegments(
      [
        [0, 1],
        [5, 10]
      ],
      3,
      90,
      0
    ),
    5
  )

  assert.strictEqual(
    wm.kthSegments(
      [
        [0, 1],
        [5, 10]
      ],
      3,
      0
    ),
    8
  )

  assert.deepStrictEqual(
    wm.kthValueAndSumSegments(
      [
        [0, 1],
        [5, 10]
      ],
      30,
      0
    ),
    [null, 41]
  )

  assert.strictEqual(
    wm.sumSegments(
      [
        [0, 1],
        [5, 10]
      ],
      1,
      3,
      0
    ),
    13
  )

  assert.strictEqual(
    wm.sumAllSegments([
      [0, 1],
      [5, 10]
    ]),
    41
  )

  assert.strictEqual(wm.medianSegments(true, [[0, 5]], 0), 3)

  const [count, sum] = wm.maxRight(0, 10, (count, sum) => count <= 3)
  assert.strictEqual(count, 3)
  assert.strictEqual(sum, 6)
  console.log('Pass')

  console.time('WaveletMatrixSegments')
  const getBig = (len = 1e5): number[] => {
    const big = Array(len)
    for (let i = 0; i < len; i++) {
      big[i] = ~~(Math.random() * 2e9) - 1e9
    }
    return big
  }
  const big = getBig()
  const wmBig = new WaveletMatrixSegments(new Uint32Array(big), -1, big)
  for (let i = 0; i < 1e5; i++) {
    wmBig.kth(0, big.length, 0, 1e9)
    wmBig.kthValueAndSum(0, big.length, 0, 10)
    wmBig.kthValueAndSum(0, big.length, 0, 1e9)
    wmBig.sumAll(0, big.length)
  }
  console.timeEnd('WaveletMatrixSegments')
}

export { WaveletMatrixSegments }
