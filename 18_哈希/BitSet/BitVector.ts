/* eslint-disable no-param-reassign */
class BitVector {
  private readonly _n: number
  private readonly _block: Uint32Array
  private readonly _sum: Uint32Array

  constructor(n: number) {
    this._n = n
    this._block = new Uint32Array((n + 31) >> 5)
    this._sum = new Uint32Array((n + 31) >> 5)
  }

  add(i: number): void {
    this._block[i >> 5] |= 1 << (i & 31)
  }

  build(): void {
    for (let i = 0; i < this._block.length - 1; i++) {
      this._sum[i + 1] = this._sum[i] + onesCount32(this._block[i])
    }
  }

  has(i: number): boolean {
    return ((this._block[i >> 5] >> (i & 31)) & 1) === 1
  }

  /**
   * 查询 [0, right) 中 digit 的个数.
   */
  countPrefix(digit: 0 | 1, right: number): number {
    if (right < 0) return 0
    if (right > this._n) right = this._n
    const mask = (1 << (right & 31)) - 1
    const res = this._sum[right >> 5] + onesCount32(this._block[right >> 5] & mask)
    return digit === 1 ? res : right - res
  }

  /**
   * 查询 [start, end) 中 digit 的个数.
   */
  count(digit: 0 | 1, start: number, end: number): number {
    return this.countPrefix(digit, end) - this.countPrefix(digit, start)
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
      const mid = (left + right) >> 1
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

  clear(): void {
    this._block.fill(0)
    this._sum.fill(0)
  }

  get size(): number {
    return this._n
  }

  toString(): string {
    const res = []
    for (let i = 0; i < this._n; i++) {
      this.has(i) && res.push(i)
    }
    return `BitVector(${this._n}): [${res.join(', ')}]`
  }
}

function onesCount32(uint32: number): number {
  uint32 -= (uint32 >>> 1) & 0x55555555
  uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
  return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
}

export { BitVector }
