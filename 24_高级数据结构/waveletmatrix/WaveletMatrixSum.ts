import assert from 'assert'
/* eslint-disable no-inner-declarations */
/* eslint-disable semi-style */
/* eslint-disable no-param-reassign */

// 维护区间贡献的 Wavelet Matrix
// !注意查询区间贡献时, 异或无效

// CountRange(start, end, a, b, xor) - 区间 [start, end) 中值在 [a, b) 之间的数的个数和这些数的和.
// CountPrefix(start, end, x, xor) - 区间 [start, end) 中值在 [0, x) 之间的数的个数和这些数的和.

// Kth(start, end, k, xor) - 区间 [start, end) 中第 k 小的数(0-indexed) 和前 k 小的数的和(不包括这个数).

// Floor(start, end, x, xor) - 区间 [start, end) 中值小于等于 x 的最大值
// Ceiling(start, end, x, xor) - 区间 [start, end) 中值大于等于 x 的最小值

// MaxRightValue(start, end, xor, check) -
//   返回使得 check(prefixSum) 为 true 的最大value, 其中prefixSum为[0,val)内的数的和.
// MaxRightCount(start, end, xor, check) -
//   返回使得 check(prefixSum) 为 true 的区间前缀个数的最大值.

/**
 * 常数大约是线段树的级别.
 */
class WaveletMatrixSum {
  private readonly _log: number
  private readonly _mid: Uint32Array
  private readonly _bv: _BV[]
  private readonly _preSum: number[][]

  /**
   * @param nums 0 <= nums[i] < 2**31 (2e9+)
   * 如果不在这个范围内，需要将所有数(松)离散化.
   * @param log 如果要支持异或,则需要按照异或的值来决定值域.
   * 设为`-1`时表示不使用异或.
   */
  constructor(nums: Uint32Array, log = -1) {
    nums = nums.slice()
    if (log === -1) {
      log = 32 - Math.clz32(Math.max(...nums)) // bitLength
    }

    const n = nums.length
    const mid = new Uint32Array(log)
    const bv = Array(log)
    for (let i = 0; i < log; i++) {
      bv[i] = new _BV(n) // 1-based
    }
    const preSum = Array(log + 1).fill(0)
    for (let i = 0; i <= log; i++) {
      preSum[i] = Array(n + 1).fill(0)
    }
    let a0 = new Uint32Array(n)
    const a1 = new Uint32Array(n)
    for (let d = log - 1; d >= -1; d--) {
      let p0 = 0
      let p1 = 0
      for (let i = 0; i < n; i++) {
        preSum[d + 1][i + 1] = preSum[d + 1][i] + nums[i]
      }
      if (d === -1) {
        break
      }

      for (let i = 0; i < n; i++) {
        const f = (nums[i] >> d) & 1
        if (f) {
          bv[d].add(i)
          a1[p1] = nums[i]
          p1++
        } else {
          a0[p0] = nums[i]
          p0++
        }
      }

      mid[d] = p0
      bv[d].build()
      const tmp = nums
      nums = a0
      a0 = tmp
      for (let i = 0; i < p1; i++) {
        nums[p0 + i] = a1[i]
      }
    }

    this._log = log
    this._mid = mid
    this._bv = bv
    this._preSum = preSum
  }

  /**
   * [start, end) 中范围在 [floor, higher) 之间的 `数的个数, 这些数的和`.
   */
  countRange(
    start: number,
    end: number,
    floor: number,
    higher: number,
    xor = 0
  ): [res: number, sum: number] {
    const [c1, s1] = this.countPrefix(start, end, floor, xor)
    const [c2, s2] = this.countPrefix(start, end, higher, xor)
    return [c2 - c1, s2 - s1]
  }

  /**
   * [start, end) 中范围在 [0, higher) 之间的 `数的个数, 这些数的和`.
   */
  countPrefix(start: number, end: number, higher: number, xor = 0): [res: number, sum: number] {
    if (higher >= 1 << this._log) {
      return [end - start, this._get(this._log, start, end)]
    }
    let count = 0
    let sum = 0
    for (let d = this._log - 1; ~d; d--) {
      const add = (higher >>> d) & 1
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)

      if (add) {
        count += kf
        if (f) {
          sum += this._get(d, start + this._mid[d] - l0, end + this._mid[d] - r0)
          start = l0
          end = r0
        } else {
          sum += this._get(d, l0, r0)
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

    return [count, sum]
  }

  /**
   * [start, end) 中的 `第k小的元素, 前k个元素(不包括第k小的元素)的和`.
   * k 从 `0` 开始.
   * 如果 k < 0, 返回 [null, 0]; 如果 k >= right-left, 返回 [null, 区间和].
   */
  kth(
    start: number,
    end: number,
    k: number,
    xor = 0
  ): [res: number, sum: number] | [res: null, sum: number] {
    if (k < 0) return [null, 0]
    if (end - start <= k) return [null, this._get(this._log, start, end)]
    let res = 0
    let sum = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f * (end - start - (r0 - l0)) + (f ^ 1) * (r0 - l0)
      if (k < kf) {
        if (f) {
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        } else {
          start = l0
          end = r0
        }
      } else {
        k -= kf
        res |= 1 << d
        if (f) {
          sum += this._get(d, start + this._mid[d] - l0, end + this._mid[d] - r0)
          start = l0
          end = r0
        } else {
          sum += this._get(d, l0, r0)
          start += this._mid[d] - l0
          end += this._mid[d] - r0
        }
      }
    }

    if (k) {
      sum += this._get(0, start, start + k)
    }
    return [res, sum]
  }

  /**
   * 返回使得 check(prefixSum) 为 true 的值域上的最大值 val.
   * !(即区间内小于 val 的数的和 prefixSum 满足 check函数, 找到这样的最大的 val)
   * 如果整个区间都满足, 则返回`Infinity`.
   * @example
   * val = 5 => 即区间内值域在 [0,5) 中的数的和满足 check 函数.
   */
  maxRightValue(
    start: number,
    end: number,
    predicate: (preSum: number) => boolean,
    xor = 0
  ): number {
    if (predicate(this._get(this._log, start, end))) return Infinity

    let res = 0
    let sum = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const loSum = f
        ? this._get(d, start + this._mid[d] - l0, end + this._mid[d] - r0)
        : this._get(d, l0, r0)
      if (predicate(sum + loSum)) {
        sum += loSum
        res |= 1 << d
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

    return res
  }

  /**
   * 返回使得 check(prefixSum) 为 true 的区间前缀个数的最大值.
   * @example
   * count = 4 => 即区间内的数排序后, 前4个数的和满足 check 函数.
   */
  maxRightCount(
    start: number,
    end: number,
    predicate: (preSum: number) => boolean,
    xor = 0
  ): number {
    if (predicate(this._get(this._log, start, end))) return end - start

    let res = 0
    let sum = 0
    for (let d = this._log - 1; ~d; d--) {
      const f = (xor >>> d) & 1
      const l0 = this._bv[d].rank(start, 0)
      const r0 = this._bv[d].rank(end, 0)
      const kf = f ? end - start - (r0 - l0) : r0 - l0
      const loSum = f
        ? this._get(d, start + this._mid[d] - l0, end + this._mid[d] - r0)
        : this._get(d, l0, r0)

      if (predicate(sum + loSum)) {
        sum += loSum
        res += kf
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

    res += WaveletMatrixSum._binarySearch(
      k => predicate(sum + this._get(0, start, start + k)),
      0,
      end - start
    )
    return res
  }

  floor(start: number, end: number, value: number, xor = 0): number | null {
    const less = this.countPrefix(start, end, value, xor)[0]
    return less ? this.kth(start, end, less - 1, xor)[0] : null
  }

  ceiling(start: number, end: number, value: number, xor = 0): number | null {
    const less = this.countPrefix(start, end, value, xor)[0]
    return less === end - start ? null : this.kth(start, end, less, xor)[0]
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
  const wm = new WaveletMatrixSum(new Uint32Array([3, 1, 2, 4, 5, 6, 7, 8, 9, 10]))
  assert.deepStrictEqual(wm.countRange(0, 10, 3, 7, 0), [4, 18])
  assert.deepStrictEqual(wm.kth(0, 10, 3, 0), [4, 6])
  assert.strictEqual(
    wm.maxRightValue(0, 10, preSum => preSum < 11, 0),
    5
  ) // 5 即值域在 [0,5) 中的数的和小于 11
  assert.strictEqual(
    wm.maxRightCount(0, 10, preSum => preSum < 11, 0),
    4
  ) // 4 即排序后前 4 个数的和小于 11
  assert.deepStrictEqual(wm.ceiling(0, 10, 3, 0), 3)
  assert.deepStrictEqual(wm.floor(0, 10, 3, 0), 2)
  console.log('pass')
}

export { WaveletMatrixSum }
