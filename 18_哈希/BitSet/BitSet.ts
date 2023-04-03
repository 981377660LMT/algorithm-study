/* eslint-disable no-labels */
// 位运算:
// ~~ 或者 | => Math.floor (int)
// ~~ 或者 | + >>>0 => Math.floor (uint)
// & (-1 << k) => 清除最低k位的 1

/**
 * @see {@link https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go}
 */
class BitSet {
  static from(arrayLike: ArrayLike<number | string>): BitSet {
    const bitSet = new BitSet(arrayLike.length)
    for (let i = 0; i < arrayLike.length; i++) {
      if (Number(arrayLike[i]) === 1) {
        bitSet.add(i)
      }
    }
    return bitSet
  }

  // lowbit.bit_length() - 1
  private static _trailingZeros32(uint32: number): number {
    if (uint32 === 0) return 32
    return 31 - Math.clz32(uint32 & -uint32) // bitLength32(uint32 & -uint32) - 1
  }
  private static _onesCount32(uint32: number): number {
    uint32 -= (uint32 >>> 1) & 0x55555555
    uint32 = (uint32 & 0x33333333) + ((uint32 >>> 2) & 0x33333333)
    return (((uint32 + (uint32 >>> 4)) & 0x0f0f0f0f) * 0x01010101) >>> 24
  }
  private static _bitLength32(uint32: number): number {
    return 32 - Math.clz32(uint32)
  }
  private static readonly _W = 32

  private readonly _bits: Uint32Array
  private readonly _n: number

  constructor(n: number) {
    this._bits = new Uint32Array(~~(n / BitSet._W) + 1)
    this._n = n
  }

  add(i: number): void {
    this._bits[~~(i / BitSet._W)] |= 1 << i % BitSet._W
  }

  /**
   * [start, end) 范围内的位置为 1
   */
  addRange(start: number, end: number): void {
    const maskL = ~0 << start % BitSet._W
    const maskR = ~0 << end % BitSet._W
    let i = ~~(start / BitSet._W)
    if (i === ~~(end / BitSet._W)) {
      this._bits[i] |= maskL ^ maskR
      return
    }
    this._bits[i] |= maskL
    for (i++; i < ~~(end / BitSet._W); i++) {
      this._bits[i] = ~0
    }
    this._bits[i] |= ~maskR
  }

  has(i: number): boolean {
    return !!(this._bits[~~(i / BitSet._W)] & (1 << i % BitSet._W))
  }

  discard(i: number): void {
    this._bits[~~(i / BitSet._W)] &= ~(1 << i % BitSet._W)
  }

  /**
   * [start, end) 范围内的位置为 0
   */
  discardRange(start: number, end: number): void {
    const maskL = ~0 << start % BitSet._W
    const maskR = ~0 << end % BitSet._W
    let i = ~~(start / BitSet._W)
    if (i === ~~(end / BitSet._W)) {
      this._bits[i] &= ~maskL | maskR
      return
    }
    this._bits[i] &= ~maskL
    for (i++; i < ~~(end / BitSet._W); i++) {
      this._bits[i] = 0
    }
    this._bits[i] &= maskR
  }

  flip(i: number): void {
    this._bits[~~(i / BitSet._W)] ^= 1 << i % BitSet._W
  }

  /**
   * [start, end) 范围内的位取反
   */
  flipRange(start: number, end: number): void {
    const maskL = ~0 << start % BitSet._W
    const maskR = ~0 << end % BitSet._W
    let i = ~~(start / BitSet._W)
    if (i === ~~(end / BitSet._W)) {
      this._bits[i] ^= maskL ^ maskR
      return
    }
    this._bits[i] ^= maskL
    for (i++; i < ~~(end / BitSet._W); i++) {
      this._bits[i] = ~this._bits[i]
    }
    this._bits[i] ^= ~maskR
  }

  clear(): void {
    this._bits.fill(0)
  }

  /**
   * [start, end) 范围内是否全为 1
   */
  allOne(start: number, end: number): boolean {
    let i = ~~(start / BitSet._W)
    if (i === ~~(end / BitSet._W)) {
      const mask = (~0 << start % BitSet._W) ^ (~0 << end % BitSet._W)
      return (this._bits[i] & mask) === mask
    }
    let mask = ~0 << start % BitSet._W
    if ((this._bits[i] & mask) !== mask) {
      return false
    }
    for (i++; i < ~~(end / BitSet._W); i++) {
      if (~this._bits[i] !== 0) {
        return false
      }
    }
    mask = ~0 << end % BitSet._W
    return ~(this._bits[~~(end / BitSet._W)] | mask) === 0
  }

  /**
   * [start, end) 范围内是否全为 0
   */
  allZero(start: number, end: number): boolean {
    let i = ~~(start / BitSet._W)
    if (i === ~~(end / BitSet._W)) {
      const mask = (~0 << start % BitSet._W) ^ (~0 << end % BitSet._W)
      return (this._bits[i] & mask) === 0
    }
    if (this._bits[i] >> start % BitSet._W !== 0) {
      return false
    }
    for (i++; i < ~~(end / BitSet._W); i++) {
      if (this._bits[i] !== 0) {
        return false
      }
    }
    const mask = ~0 << end % BitSet._W
    return (this._bits[~~(end / BitSet._W)] & ~mask) === 0
  }

  /**
   * 返回第一个 0 的下标，若不存在则返回-1
   * @param position 从哪个位置开始查找
   */
  indexOfZero(position = 0): number {
    if (position === 0) {
      return this._indexOfZero()
    }

    let i = ~~(position / BitSet._W)
    if (i < this._bits.length) {
      let v = this._bits[i]
      if (position % BitSet._W > 0) {
        v |= ~(~0 << position % BitSet._W)
      }
      if (~v !== 0) {
        const res = (i * BitSet._W) | BitSet._trailingZeros32(~v)
        return res < this._n ? res : -1
      }
      for (i++; i < this._bits.length; i++) {
        if (~this._bits[i] !== 0) {
          const res = (i * BitSet._W) | BitSet._trailingZeros32(~this._bits[i])
          return res < this._n ? res : -1
        }
      }
    }

    return -1
  }

  /**
   * 返回第一个 1 的下标，若不存在则返回-1
   * @param position 从哪个位置开始查找
   */
  indexOfOne(position = 0): number {
    if (position === 0) {
      return this._indexOfOne()
    }

    // eslint-disable-next-line space-in-parens
    for (let i = ~~(position / BitSet._W); i < this._bits.length; ) {
      const v = this._bits[i] & (~0 << position % BitSet._W)
      if (v !== 0) {
        return (i * BitSet._W) | BitSet._trailingZeros32(v)
      }
      for (i++; i < this._bits.length; i++) {
        if (this._bits[i] !== 0) {
          return (i * BitSet._W) | BitSet._trailingZeros32(this._bits[i])
        }
      }
    }

    return -1
  }

  /**
   * 返回 [start, end) 范围内 1 的个数
   */
  onesCount(start = 0, end = this._n): number {
    if (start < 0) {
      start = 0
    }
    if (end > this._n) {
      end = this._n
    }
    if (start === 0 && end === this._n) {
      return this._onesCount()
    }

    let pos1 = ~~(start / BitSet._W)
    const pos2 = ~~(end / BitSet._W)
    if (pos1 === pos2) {
      return BitSet._onesCount32(
        this._bits[pos1] & (~0 << start % BitSet._W) & ((1 << end % BitSet._W) - 1)
      )
    }

    let count = 0
    if (start % BitSet._W > 0) {
      count += BitSet._onesCount32(this._bits[pos1] & (~0 << start % BitSet._W))
      pos1++
    }
    for (let i = pos1; i < pos2; i++) {
      count += BitSet._onesCount32(this._bits[i])
    }
    if (end % BitSet._W > 0) {
      count += BitSet._onesCount32(this._bits[pos2] & ((1 << end % BitSet._W) - 1))
    }
    return count
  }

  equals(other: BitSet): boolean {
    if (this._bits.length !== other._bits.length) {
      return false
    }
    for (let i = 0; i < this._bits.length; i++) {
      if (this._bits[i] !== other._bits[i]) {
        return false
      }
    }
    return true
  }

  isSubset(other: BitSet): boolean {
    if (this._bits.length > other._bits.length) {
      return false
    }
    for (let i = 0; i < this._bits.length; i++) {
      if ((this._bits[i] & other._bits[i]) >>> 0 !== this._bits[i]) {
        return false
      }
    }
    return true
  }

  isSuperset(other: BitSet): boolean {
    if (this._bits.length < other._bits.length) {
      return false
    }
    for (let i = 0; i < other._bits.length; i++) {
      if ((this._bits[i] & other._bits[i]) >>> 0 !== other._bits[i]) {
        return false
      }
    }
    return true
  }

  ior(other: BitSet): void {
    for (let i = 0; i < this._bits.length; i++) {
      this._bits[i] |= other._bits[i]
    }
  }

  or(other: BitSet): BitSet {
    const res = new BitSet(this._bits.length)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] | other._bits[i]
    }
    return res
  }

  iand(other: BitSet): void {
    for (let i = 0; i < this._bits.length; i++) {
      this._bits[i] &= other._bits[i]
    }
  }

  and(other: BitSet): BitSet {
    const res = new BitSet(this._bits.length)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] & other._bits[i]
    }
    return res
  }

  copy(): BitSet {
    const res = new BitSet(this._n)
    res._bits.set(this._bits)
    return res
  }

  bitLength(): number {
    return this._lastIndexOfOne() + 1
  }

  toString(): string {
    const sb: string[] = []
    for (let i = 0; i < this._bits.length; i++) {
      // eslint-disable-next-line newline-per-chained-call
      let bits = this._bits[i].toString(2).padStart(BitSet._W, '0').split('').reverse().join('')
      if (i === this._bits.length - 1) {
        bits = bits.slice(0, this._n % BitSet._W)
      }
      sb.push(bits)
    }
    return sb.join('')
  }

  _indexOfZero(): number {
    for (let i = 0; i < this._bits.length; i++) {
      const x = this._bits[i]
      if (~x !== 0) {
        return (i * BitSet._W) | BitSet._trailingZeros32(~x)
      }
    }
    return -1
  }

  _indexOfOne(): number {
    for (let i = 0; i < this._bits.length; i++) {
      const x = this._bits[i]
      if (x !== 0) {
        return (i * BitSet._W) | BitSet._trailingZeros32(x)
      }
    }
    return -1
  }

  _lastIndexOfOne(): number {
    for (let i = this._bits.length - 1; i >= 0; i--) {
      const x = this._bits[i]
      if (x !== 0) {
        return (i * BitSet._W) | (BitSet._bitLength32(x) - 1)
      }
    }
    return -1
  }

  _onesCount(): number {
    let count = 0
    for (let i = 0; i < this._bits.length; i++) {
      count += BitSet._onesCount32(this._bits[i])
    }
    return count
  }
}

export { BitSet }

if (require.main === module) {
  // // onesCountRange test
  // const n = 330
  // const bs = new BitSet(n)
  // const nums = new Uint8Array(n)
  // for (let i = 0; i < 10; i++) {
  //   if (Math.random() > 0.5) {
  //     bs.add(i)
  //     nums[i] = 1
  //   }
  // }
  // // onesCountRange test
  // loop: for (let i = 0; i < n; i++) {
  //   for (let j = i; j < n; j++) {
  //     const count = bs.onesCount(i, j)
  //     const expected = nums.slice(i, j).reduce((a, b) => a + b, 0)
  //     if (count !== expected) {
  //       console.log(i, j, count, expected)
  //       console.log(nums, bs.toString())
  //       break loop
  //     }
  //   }
  // }
  // console.log('ok')
}
