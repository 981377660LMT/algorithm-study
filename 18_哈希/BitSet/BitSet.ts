/* eslint-disable no-labels */
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go
// 位运算:
// ~~ 或者 | => Math.floor (int)
// ~~ 或者 | + >>>0 => Math.floor (uint)
// & (-1 << k) => 清除最低k位的 1

/**
 * 位集,用于存储大量的布尔值,可以有效地节省内存.
 * 1e9个元素 => 125MB.
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
    if (!uint32) return 32
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

  readonly n: number
  private readonly _bits: Uint32Array

  constructor(n: number, filledValue: 0 | 1 = 0) {
    if (n <= 0) throw new RangeError('n must be positive')
    this.n = n
    this._bits = filledValue
      ? new Uint32Array((n + 31) >>> 5).fill(~0)
      : new Uint32Array((n + 31) >>> 5)
    this._bits[this._bits.length - 1] >>>= (this._bits.length << 5) - n
  }

  add(i: number): void {
    this._bits[i >> 5] |= 1 << (i & 31)
  }

  /**
   * [start, end) 范围内的位置为 1
   */
  addRange(start: number, end: number): void {
    const maskL = ~0 << (start & 31)
    const maskR = ~0 << (end & 31)
    let i = start >> 5
    if (i === end >> 5) {
      this._bits[i] |= maskL ^ maskR
      return
    }
    this._bits[i] |= maskL
    for (i++; i < end >> 5; i++) {
      this._bits[i] = ~0
    }
    this._bits[i] |= ~maskR
  }

  has(i: number): boolean {
    return !!(this._bits[i >> 5] & (1 << (i & 31)))
  }

  discard(i: number): void {
    this._bits[i >> 5] &= ~(1 << (i & 31))
  }

  /**
   * [start, end) 范围内的位置为 0
   */
  discardRange(start: number, end: number): void {
    const maskL = ~0 << (start & 31)
    const maskR = ~0 << (end & 31)
    let i = start >> 5
    if (i === end >> 5) {
      this._bits[i] &= ~maskL | maskR
      return
    }
    this._bits[i] &= ~maskL
    for (i++; i < end >> 5; i++) {
      this._bits[i] = 0
    }
    this._bits[i] &= maskR
  }

  flip(i: number): void {
    this._bits[i >> 5] ^= 1 << (i & 31)
  }

  /**
   * [start, end) 范围内的位取反
   */
  flipRange(start: number, end: number): void {
    const maskL = ~0 << (start & 31)
    const maskR = ~0 << (end & 31)
    let i = start >> 5
    if (i === end >> 5) {
      this._bits[i] ^= maskL ^ maskR
      return
    }
    this._bits[i] ^= maskL
    for (i++; i < end >> 5; i++) {
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
    let i = start >> 5
    if (i === end >> 5) {
      const mask = (~0 << (start & 31)) ^ (~0 << (end & 31))
      return (this._bits[i] & mask) === mask
    }
    let mask = ~0 << (start & 31)
    if ((this._bits[i] & mask) !== mask) {
      return false
    }
    for (i++; i < end >> 5; i++) {
      if (~this._bits[i]) {
        return false
      }
    }
    mask = ~0 << (end & 31)
    return !~(this._bits[end >> 5] | mask)
  }

  /**
   * [start, end) 范围内是否全为 0
   */
  allZero(start: number, end: number): boolean {
    let i = start >> 5
    if (i === end >> 5) {
      const mask = (~0 << (start & 31)) ^ (~0 << (end & 31))
      return !(this._bits[i] & mask)
    }
    if (this._bits[i] >> (start & 31)) {
      return false
    }
    for (i++; i < end >> 5; i++) {
      if (this._bits[i]) {
        return false
      }
    }
    const mask = ~0 << (end & 31)
    return !(this._bits[end >> 5] & ~mask)
  }

  /**
   * 返回第一个 0 的下标，若不存在则返回-1
   * @param position 从哪个位置开始查找
   */
  indexOfZero(position = 0): number {
    if (position === 0) {
      return this._indexOfZero()
    }

    let i = position >> 5
    if (i < this._bits.length) {
      let v = this._bits[i]
      if (position & 31) {
        v |= ~(~0 << (position & 31))
      }
      if (~v) {
        const res = (i << 5) | BitSet._trailingZeros32(~v)
        return res < this.n ? res : -1
      }
      for (i++; i < this._bits.length; i++) {
        if (~this._bits[i]) {
          const res = (i << 5) | BitSet._trailingZeros32(~this._bits[i])
          return res < this.n ? res : -1
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
    if (!position) {
      return this._indexOfOne()
    }

    // eslint-disable-next-line space-in-parens
    for (let i = position >> 5; i < this._bits.length; ) {
      const v = this._bits[i] & (~0 << (position & 31))
      if (v) {
        return (i << 5) | BitSet._trailingZeros32(v)
      }
      for (i++; i < this._bits.length; i++) {
        if (this._bits[i]) {
          return (i << 5) | BitSet._trailingZeros32(this._bits[i])
        }
      }
    }

    return -1
  }

  /**
   * 返回 [start, end) 范围内 1 的个数
   */
  onesCount(start = 0, end = this.n): number {
    if (start < 0) {
      start = 0
    }
    if (end > this.n) {
      end = this.n
    }
    if (!start && end === this.n) {
      return this._onesCount()
    }

    let pos1 = start >> 5
    const pos2 = end >> 5
    if (pos1 === pos2) {
      return BitSet._onesCount32(this._bits[pos1] & (~0 << (start & 31)) & ((1 << (end & 31)) - 1))
    }

    let count = 0
    if ((start & 31) > 0) {
      count += BitSet._onesCount32(this._bits[pos1] & (~0 << (start & 31)))
      pos1++
    }
    for (let i = pos1; i < pos2; i++) {
      count += BitSet._onesCount32(this._bits[i])
    }
    if ((end & 31) > 0) {
      count += BitSet._onesCount32(this._bits[pos2] & ((1 << (end & 31)) - 1))
    }
    return count
  }

  fill(value: 0 | 1): void {
    this._bits.fill(value ? ~0 : 0)
    this._bits[this._bits.length - 1] >>>= (this._bits.length << 5) - this.n
  }

  slice(start = 0, end = this.n): BitSet {
    if (start < 0) start = 0
    if (end > this.n) end = this.n
    const res = new BitSet(end - start)
    const remain = (end - start) & 31
    for (let _ = 0; _ < remain; _++) {
      if (this.has(end - 1)) res.add(end - start - 1)
      end--
    }

    const n = (end - start) >>> 5
    const hi = start & 31
    const lo = 32 - hi
    const s = start >>> 5
    if (!hi) {
      for (let i = 0; i < n; i++) {
        res._bits[i] = this._bits[s + i]
      }
    } else {
      for (let i = 0; i < n; i++) {
        res._bits[i] = (this._bits[s + i] >>> hi) ^ (this._bits[s + i + 1] << lo)
      }
    }
    return res
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

  ior(other: BitSet): BitSet {
    for (let i = 0; i < this._bits.length; i++) {
      this._bits[i] |= other._bits[i]
    }
    return this
  }

  or(other: BitSet): BitSet {
    const res = new BitSet(this.n)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] | other._bits[i]
    }
    return res
  }

  iand(other: BitSet): BitSet {
    for (let i = 0; i < this._bits.length; i++) {
      this._bits[i] &= other._bits[i]
    }
    return this
  }

  and(other: BitSet): BitSet {
    const res = new BitSet(this.n)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] & other._bits[i]
    }
    return res
  }

  ixor(other: BitSet): BitSet {
    for (let i = 0; i < this._bits.length; i++) {
      this._bits[i] ^= other._bits[i]
    }
    return this
  }

  xor(other: BitSet): BitSet {
    const res = new BitSet(this.n)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] ^ other._bits[i]
    }
    return res
  }

  /**
   * 类似 {@link Uint8Array.prototype.set}，区别是超出范围的位会被丢弃，而不是抛出异常.
   * @param other 要赋值的位集.
   * @param offset 从哪里开始赋值.
   */
  set(other: BitSet, offset = 0): void {
    let left = offset
    let right = offset + other.n
    let a = 0
    let b = other.n
    while (left < right && left & 31) {
      if (other.has(a++)) {
        this.add(left++)
      } else {
        this.discard(left++)
      }
    }
    while (left < right && right & 31) {
      if (other.has(--b)) {
        this.add(--right)
      } else {
        this.discard(--right)
      }
    }

    const l = left >>> 5
    const r = right >>> 5
    const s = a >>> 5
    const n = r - l
    if (!(a & 31)) {
      this._bits.set(other._bits.subarray(s, s + n), l)
    } else {
      const hi = a & 31
      const lo = 32 - hi
      for (let i = 0; i < n; i++) {
        this._bits[l + i] = (other._bits[s + i] >>> hi) | (other._bits[1 + s + i] << lo)
      }
    }
  }

  copy(): BitSet {
    const res = new BitSet(this.n)
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
      let bits = this._bits[i].toString(2).padStart(32, '0').split('').reverse().join('')
      if (i === this._bits.length - 1) {
        bits = bits.slice(0, this.n - (i << 5))
      }
      sb.push(bits)
    }
    return sb.join('')
  }

  /**
   * 遍历所有 1 的位置.
   */
  forEach(callback: (value: number) => void): void {
    this._bits.forEach((v, i) => {
      for (; v > 0; v &= v - 1) {
        const j = (i << 5) | BitSet._trailingZeros32(v)
        callback(j)
      }
    })
  }

  private _indexOfZero(): number {
    for (let i = 0; i < this._bits.length; i++) {
      const x = this._bits[i]
      if (~x) {
        return (i << 5) | BitSet._trailingZeros32(~x)
      }
    }
    return -1
  }

  private _indexOfOne(): number {
    for (let i = 0; i < this._bits.length; i++) {
      const x = this._bits[i]
      if (x) {
        return (i << 5) | BitSet._trailingZeros32(x)
      }
    }
    return -1
  }

  private _lastIndexOfOne(): number {
    for (let i = this._bits.length - 1; i >= 0; i--) {
      const x = this._bits[i]
      if (x) {
        return (i << 5) | (BitSet._bitLength32(x) - 1)
      }
    }
    return -1
  }

  private _onesCount(): number {
    let count = 0
    for (let i = 0; i < this._bits.length; i++) {
      count += BitSet._onesCount32(this._bits[i])
    }
    return count
  }

  /**
   * hack.
   * ![start,end) 范围内数与bitset相交的数.end-start<32.
   * @example
   * [start,end) = [10,16)
   * [15,14,13,12,11,10] & Bitset(10,11,13,15) => 101011
   */
  _hasRange(start: number, end: number): number {
    const posL = start >> 5
    const shiftL = start & 31
    const posR = end >> 5
    const shiftR = end & 31
    const maskL = ~(~0 << shiftL)
    const maskR = ~(~0 << shiftR)
    if (posL === posR) return (this._bits[posL] & maskR) >>> shiftL
    return ((this._bits[posL] & ~maskL) >>> shiftL) | ((this._bits[posR] & maskR) << (32 - shiftL))
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

  const set = new BitSet(33)
  set.fill(1)
  set.fill(1)
  console.log(set.has(34))
  set.discard(1)
  console.log(set.slice(0, 33).toString())

  console.log(set.toString())
  const mm = new BitSet(3)
  // mm.add(0)
  set.set(mm, 3)
  console.log(set.toString())

  console.log('set')
  //  1001 -> 01
  const bs1 = BitSet.from('1001')
  const bs2 = BitSet.from('01')
  bs1.set(bs2, 3)
  console.log(bs1.toString())
}
