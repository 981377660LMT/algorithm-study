/* eslint-disable no-param-reassign */
// !这个版本比较快

import { strictEqual } from 'assert'

class WaveletMatrix {
  private readonly _n: number
  private readonly _maxLog: number
  private readonly _mat: _BV[]
  private readonly _zs: Uint32Array
  private readonly _buff1: Uint32Array
  private readonly _buff2: Uint32Array

  /**
   * @param nums 0 <= nums[i] < 2**31 (2e9+)
   * 如果不在这个范围内，需要将所有数(松)离散化.
   */
  constructor(nums: Uint32Array) {
    nums = nums.slice()
    let max_ = 0
    for (let i = 0; i < nums.length; i++) {
      max_ = Math.max(max_, nums[i])
    }
    const n = nums.length
    const maxLog = 32 - Math.clz32(max_) + 1 // bit_len + 1
    const mat = Array(maxLog)
    for (let i = 0; i < maxLog; i++) mat[i] = new _BV(n + 1)
    const zs = new Uint32Array(maxLog)
    const buff1 = new Uint32Array(maxLog)
    const buff2 = new Uint32Array(maxLog)
    let ls = new Uint32Array(n)
    const rs = new Uint32Array(n)
    for (let dep = 0; dep < maxLog; dep++) {
      let p = 0
      let q = 0
      for (let i = 0; i < n; i++) {
        const k = (nums[i] >>> (maxLog - dep - 1)) & 1
        if (k) {
          rs[q] = nums[i]
          mat[dep].add(i)
          q++
        } else {
          ls[p] = nums[i]
          p++
        }
      }

      zs[dep] = p
      mat[dep].build()
      ls = nums
      nums.set(rs.subarray(0, q), p)
    }

    this._n = n
    this._maxLog = maxLog
    this._mat = mat
    this._zs = zs
    this._buff1 = buff1
    this._buff2 = buff2
  }

  /**
   * [start, end) 中值为 value 的数的个数.
   */
  count(start: number, end: number, value: number): number {
    return this._count(value, end) - this._count(value, start)
  }

  /**
   * [start, end) 中值在 [lower, upper) 内的数的个数.
   */
  countRange(start: number, end: number, floor: number, higher: number): number {
    return this._freqDfs(0, start, end, 0, floor, higher)
  }

  /**
   * [start, end) 中第 k(0-indexed) 小的数.
   */
  kth(start: number, end: number, k: number): number | null {
    return this.kthMax(start, end, end - start - k - 1)
  }

  /**
   * [start, end) 中第 k(0-indexed) 大的数.
   */
  kthMax(start: number, end: number, k: number): number | null {
    if (k < 0 || k >= end - start) {
      return null
    }
    let res = 0
    for (let dep = 0; dep < this._maxLog; dep++) {
      const p = this._mat[dep].countPrefix(1, start)
      const q = this._mat[dep].countPrefix(1, end)
      if (k < q - p) {
        start = this._zs[dep] + p
        end = this._zs[dep] + q
        res |= 1 << (this._maxLog - dep - 1)
      } else {
        k -= q - p
        start -= p
        end -= q
      }
    }
    return res
  }

  /**
   * [start, end) 中值小于 value 的最大值.
   */
  lower(start: number, end: number, value: number): number | null {
    const k = this._lt(start, end, value)
    return k ? this.kth(start, end, k - 1) : null
  }

  /**
   * [start, end) 中值大于 value 的最小值.
   */
  higher(start: number, end: number, value: number): number | null {
    const k = this._le(start, end, value)
    return k === end - start ? null : this.kth(start, end, k)
  }

  /**
   * [start, end) 中值不超过 value 的最大值.
   */
  floor(start: number, end: number, value: number): number | null {
    return this.count(start, end, value) ? value : this.lower(start, end, value)
  }

  /**
   * [start, end) 中值不小于 value 的最小值.
   */
  ceiling(start: number, end: number, value: number): number | null {
    return this.count(start, end, value) ? value : this.higher(start, end, value)
  }

  /**
   * 第k(0-indexed)个value的下标.不存在则返回-1.
   */
  index(value: number, k: number): number | -1 {
    this._count(value, this._n) // mutates buff1 and buff2
    for (let dep = this._maxLog - 1; ~dep; dep--) {
      const bit = ((value >>> (this._maxLog - dep - 1)) & 1) as 0 | 1
      k = this._mat[dep].kthWithStart(bit, k, this._buff1[dep])
      if (k < 0 || k >= this._buff2[dep]) {
        return -1
      }
      k -= this._buff1[dep]
    }
    return k
  }

  /**
   * 第k(0-indexed)个value的下标.不存在则返回-1.从start开始计数.
   */
  indexWithStart(value: number, k: number, start: number): number | -1 {
    return this.index(value, k + this._count(value, start))
  }

  private _count(value: number, end: number): number {
    let left = 0
    let right = end
    for (let dep = 0; dep < this._maxLog; dep++) {
      this._buff1[dep] = left
      this._buff2[dep] = right
      const bit = ((value >>> (this._maxLog - dep - 1)) & 1) as 0 | 1
      left = this._mat[dep].countPrefix(bit, left) + this._zs[dep] * bit
      right = this._mat[dep].countPrefix(bit, right) + this._zs[dep] * bit
    }
    return right - left
  }

  private _freqDfs(
    d: number,
    left: number,
    right: number,
    val: number,
    a: number,
    b: number
  ): number {
    if (left === right) return 0
    if (d === this._maxLog) return a <= val && val < b ? right - left : 0
    const nv = (1 << (this._maxLog - d - 1)) | val
    const nnv = ((1 << (this._maxLog - d - 1)) - 1) | nv
    if (nnv < a || b <= val) return 0
    if (a <= val && nnv < b) return right - left
    const lc = this._mat[d].countPrefix(1, left)
    const rc = this._mat[d].countPrefix(1, right)
    return (
      this._freqDfs(d + 1, left - lc, right - rc, val, a, b) +
      this._freqDfs(d + 1, lc + this._zs[d], rc + this._zs[d], nv, a, b)
    )
  }

  private _ll(left: number, right: number, val: number): [number, number] {
    let res = 0
    for (let dep = 0; dep < this._maxLog; dep++) {
      this._buff1[dep] = left
      this._buff2[dep] = right
      const bit = ((val >>> (this._maxLog - dep - 1)) & 1) as 0 | 1
      if (bit) {
        res +=
          right - left + this._mat[dep].countPrefix(1, left) - this._mat[dep].countPrefix(1, right)
      }
      left = this._mat[dep].countPrefix(bit, left) + this._zs[dep] * bit
      right = this._mat[dep].countPrefix(bit, right) + this._zs[dep] * bit
    }
    return [res, right - left]
  }

  private _lt(left: number, right: number, val: number): number {
    return this._ll(left, right, val)[0]
  }

  private _le(left: number, right: number, val: number): number {
    const res = this._ll(left, right, val)
    return res[0] + res[1]
  }
}

class _BV {
  private readonly _n: number
  private readonly _block: Uint32Array
  private readonly _sum: Uint32Array

  constructor(n: number) {
    this._n = n
    const size = 1 + (n >>> 5)
    this._block = new Uint32Array(size)
    this._sum = new Uint32Array(size)
  }

  add(i: number): void {
    this._block[i >>> 5] |= 1 << (i & 31)
  }

  build(): void {
    for (let i = 0; i < this._block.length - 1; i++) {
      this._sum[i + 1] = this._sum[i] + onesCount32(this._block[i])
    }
  }

  /**
   * 查询 [0, right) 中 digit 的个数.
   */
  countPrefix(digit: 0 | 1, right: number): number {
    const mask = (1 << (right & 31)) - 1
    const res = this._sum[right >>> 5] + onesCount32(this._block[right >>> 5] & mask)
    return digit * res + (digit ^ 1) * (right - res) // digit === 1 ? res : right - res
  }

  /**
   * 查询第 k 个 digit 的位置. k 从 0 开始.
   */
  kth(digit: 0 | 1, k: number): number {
    if (k < 0 || this.countPrefix(digit, this._n) <= k) {
      return -1
    }

    let left = 0
    let right = this._n
    while (right - left > 1) {
      const mid = (left + right) >>> 1
      if (this.countPrefix(digit, mid) >= k + 1) {
        right = mid
      } else {
        left = mid
      }
    }
    return right - 1
  }

  kthWithStart(digit: 0 | 1, k: number, start: number): number {
    return this.kth(digit, k + this.countPrefix(digit, start))
  }
}

function onesCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}

if (require.main === module) {
  const M = new WaveletMatrix(new Uint32Array([1, 2, 3, 1, 5, 6, 7, 8, 9, 10]))

  strictEqual(M.count(0, 1, 1), 1)
  strictEqual(M.countRange(0, 10, 1, 5), 4)
  strictEqual(M.index(1, 1), 3)
  strictEqual(M.indexWithStart(1, 0, 2), 3)
  strictEqual(M.kth(0, 10, 2), 2)
  strictEqual(M.kthMax(0, 10, 2), 8)
  strictEqual(M.lower(0, 3, 2), 1)
  strictEqual(M.floor(0, 3, 2), 2)
  strictEqual(M.higher(0, 10, 1), 2)
  strictEqual(M.ceiling(0, 10, 1), 1)

  // eslint-disable-next-line no-inner-declarations
  function countQuadruplets(nums: number[]) {
    const W = new WaveletMatrix(new Uint32Array(nums))
    let res = 0
    for (let j = 1; j < nums.length - 2; j++) {
      for (let k = j + 1; k < nums.length - 1; k++) {
        if (nums[k] < nums[j]) {
          const left = W.countRange(0, j, 0, nums[k])
          const right = W.countRange(k + 1, nums.length, nums[j] + 1, 1e9)
          res += left * right
        }
      }
    }
    return res
  }
}

export {}
