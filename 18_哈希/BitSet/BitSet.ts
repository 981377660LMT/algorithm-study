/* eslint-disable max-len */
/* eslint-disable no-labels */
// https://github.com/EndlessCheng/codeforces-go/blob/master/copypasta/bits.go
// 位运算:
// ~~ 或者 | => Math.floor (int)
// ~~ 或者 | + >>>0 => Math.floor (uint)
// & (-1 << k) => 清除最低k位的 1

// API:
interface IBitSet {
  readonly size: number
  readonly bits: Uint32Array

  add: (i: number) => void
  has: (i: number) => boolean
  discard: (i: number) => void
  flip: (i: number) => void
  addRange: (start: number, end: number) => void
  discardRange: (start: number, end: number) => void
  flipRange: (start: number, end: number) => void

  lsh: (k: number) => IBitSet
  rsh: (k: number) => IBitSet

  fill: (value: 0 | 1) => void
  clear: () => void

  onesCount: (start?: number, end?: number) => number
  allOne: (start: number, end: number) => boolean
  allZero: (start: number, end: number) => boolean

  indexOfOne: (position?: number) => number
  indexOfZero: (position?: number) => number
  next(index: number): number
  prev(index: number): number

  equals: (other: IBitSet) => boolean
  isSubset: (other: IBitSet) => boolean
  isSuperset: (other: IBitSet) => boolean

  ior: (other: IBitSet) => IBitSet
  iand: (other: IBitSet) => IBitSet
  ixor: (other: IBitSet) => IBitSet
  or: (other: IBitSet) => IBitSet
  and: (other: IBitSet) => IBitSet
  xor: (other: IBitSet) => IBitSet
  iorRange: (start: number, end: number, other: IBitSet) => void
  iandRange: (start: number, end: number, other: IBitSet) => void
  ixorRange: (start: number, end: number, other: IBitSet) => void

  set: (other: IBitSet, offset?: number) => void

  slice: (start?: number, end?: number) => IBitSet
  copy: () => IBitSet
  resize: (size: number) => void

  bitLength: () => number

  forEach: (callback: (value: number) => void) => boolean | void

  toString: () => string
}

/**
 * 位集,用于存储大量的布尔值,可以有效地节省内存.
 * 1e9个元素 => 125MB.
 * 支持 `resize` 操作.
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

  private _n: number
  private _bits: Uint32Array

  /**
   * @param n 左闭右开区间`[0, n)`.
   */
  constructor(n: number, filledValue: 0 | 1 = 0) {
    if (n < 0) throw new RangeError('n must be non-negative')
    this._n = n
    this._bits = filledValue
      ? new Uint32Array((n >>> 5) + 1).fill(~0)
      : new Uint32Array((n >>> 5) + 1)
    if (n) this._bits[this._bits.length - 1] >>>= (this._bits.length << 5) - n
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

  /**
   * 左移k位(<<k).
   * !TODO:需要处理位移超出n位后多出来的1的情况.
   * @warning **不能配合切片使用.必须保证lsh后的位数不超过原位数**.
   */
  lsh(k: number): BitSet {
    if (!k) return this
    const shift = k >>> 5
    const offset = k & 31
    if (shift >= this._bits.length) {
      this._bits.fill(0)
      return this
    }
    if (!offset) {
      this._bits.copyWithin(shift, 0)
    } else {
      for (let i = this._bits.length - 1; i > shift; i--) {
        this._bits[i] =
          (this._bits[i - shift] << offset) | (this._bits[i - shift - 1] >>> (32 - offset))
      }
      this._bits[shift] = this._bits[0] << offset
    }
    this._bits.fill(0, 0, shift)
    return this
  }

  /**
   * 右移k位(>>k).
   */
  rsh(k: number): BitSet {
    if (!k) return this
    const shift = k >>> 5
    const offset = k & 31
    if (shift >= this._bits.length) {
      this._bits.fill(0)
      return this
    }
    const limit = this._bits.length - shift - 1
    if (!offset) {
      this._bits.copyWithin(0, shift)
    } else {
      for (let i = 0; i < limit; i++) {
        this._bits[i] =
          (this._bits[i + shift] >>> offset) | (this._bits[i + shift + 1] << (32 - offset))
      }
      this._bits[limit] = this._bits[this._bits.length - 1] >>> offset
    }
    this._bits.fill(0, limit + 1)
    return this
  }

  clear(): void {
    this._bits.fill(0)
  }

  /**
   * 返回右侧第一个 1 的位置(`包含`当前位置).
   * 如果不存在, 返回 {@link size}.
   */
  next(index: number): number {
    if (index < 0) index = 0
    if (index >= this._n) return this._n
    let mask = index >>> 5
    const buf = this._bits[mask] & (~0 << (index & 31))
    if (buf) return (mask << 5) + BitSet._trailingZeros32(buf)
    for (mask++; mask < this._bits.length; mask++) {
      if (this._bits[mask]) return (mask << 5) + BitSet._trailingZeros32(this._bits[mask])
    }
    return this._n
  }

  /**
   * 返回左侧第一个 1 的位置(`包含`当前位置).
   * 如果不存在, 返回 -1.
   */
  prev(index: number): number {
    if (index >= this._n - 1) index = this._n - 1
    if (index < 0) return -1
    let mask = index >>> 5
    if ((index & 31) < 31) {
      const buf = this._bits[mask] & ((1 << ((index & 31) + 1)) - 1)
      if (buf) return (mask << 5) | (31 - Math.clz32(buf))
      mask--
    }
    for (; ~mask; mask--) {
      if (this._bits[mask]) return (mask << 5) | (31 - Math.clz32(this._bits[mask]))
    }
    return -1
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
        return res < this._n ? res : -1
      }
      for (i++; i < this._bits.length; i++) {
        if (~this._bits[i]) {
          const res = (i << 5) | BitSet._trailingZeros32(~this._bits[i])
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
  onesCount(start = 0, end = this._n): number {
    if (start < 0) {
      start = 0
    }
    if (end > this._n) {
      end = this._n
    }
    if (!start && end === this._n) {
      return this._onesCount()
    }

    let pos1 = start >> 5
    const pos2 = end >> 5
    if (pos1 === pos2) {
      return BitSet._onesCount32(this._bits[pos1] & (~0 << (start & 31)) & ((1 << (end & 31)) - 1))
    }

    let count = 0
    if (start & 31) {
      count += BitSet._onesCount32(this._bits[pos1] & (~0 << (start & 31)))
      pos1++
    }
    for (let i = pos1; i < pos2; i++) {
      count += BitSet._onesCount32(this._bits[i])
    }
    if (end & 31) {
      count += BitSet._onesCount32(this._bits[pos2] & ((1 << (end & 31)) - 1))
    }
    return count
  }

  fill(value: 0 | 1): void {
    this._bits.fill(value ? ~0 : 0)
    this._bits[this._bits.length - 1] >>>= (this._bits.length << 5) - this._n
  }

  slice(start = 0, end = this._n): BitSet {
    if (start < 0) start += this._n
    if (start < 0) start = 0
    if (end < 0) end += this._n
    if (end > this._n) end = this._n
    if (start >= end) return new BitSet(0)
    if (start === 0 && end === this._n) return this.copy()

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
        res._bits[i] ^= this._bits[s + i]
      }
    } else {
      for (let i = 0; i < n; i++) {
        res._bits[i] ^= (this._bits[s + i] >>> hi) ^ (this._bits[s + i + 1] << lo)
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
    const res = new BitSet(this._n)
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
    const res = new BitSet(this._n)
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
    const res = new BitSet(this._n)
    for (let i = 0; i < this._bits.length; i++) {
      res._bits[i] = this._bits[i] ^ other._bits[i]
    }
    return res
  }

  /**
   * 将指定范围内的位与另一个位集进行或运算.
   */
  iorRange(start: number, end: number, other: BitSet): void {
    if (other._n !== end - start) throw new RangeError('length of other must equal to end-start')

    let a = 0
    let b = other._n
    while (start < end && start & 31) {
      this._bits[start >>> 5] |= +other.has(a) << (start & 31)
      a++
      start++
    }
    while (start < end && end & 31) {
      b--
      end--
      this._bits[end >>> 5] |= +other.has(b) << (end & 31)
    }

    const l = start >>> 5
    const r = end >>> 5
    const s = a >>> 5
    const n = r - l
    if (!(a & 31)) {
      for (let i = 0; i < n; i++) {
        this._bits[l + i] |= other._bits[s + i]
      }
    } else {
      const hi = a & 31
      const lo = 32 - hi
      for (let i = 0; i < n; i++) {
        this._bits[l + i] |= (other._bits[s + i] >>> hi) | (other._bits[1 + s + i] << lo)
      }
    }
  }

  /**
   * 将指定范围内的位与另一个位集进行与运算.
   */
  iandRange(start: number, end: number, other: BitSet): void {
    if (other._n !== end - start) throw new RangeError('length of other must equal to end-start')

    let a = 0
    let b = other._n
    while (start < end && start & 31) {
      if (!other.has(a)) this.discard(start)
      a++
      start++
    }
    while (start < end && end & 31) {
      b--
      end--
      if (!other.has(b)) this.discard(end)
    }

    const l = start >>> 5
    const r = end >>> 5
    const s = a >>> 5
    const n = r - l
    if (!(a & 31)) {
      for (let i = 0; i < n; i++) {
        this._bits[l + i] &= other._bits[s + i]
      }
    } else {
      const hi = a & 31
      const lo = 32 - hi
      for (let i = 0; i < n; i++) {
        this._bits[l + i] &= (other._bits[s + i] >>> hi) | (other._bits[1 + s + i] << lo)
      }
    }
  }

  /**
   * 将指定范围内的位与另一个位集进行异或运算.
   */
  ixorRange(start: number, end: number, other: BitSet): void {
    if (other._n !== end - start) throw new RangeError('length of other must equal to end-start')

    let a = 0
    let b = other._n
    while (start < end && start & 31) {
      this._bits[start >>> 5] ^= +other.has(a) << (start & 31)
      a++
      start++
    }
    while (start < end && end & 31) {
      b--
      end--
      this._bits[end >>> 5] ^= +other.has(b) << (end & 31)
    }

    const l = start >>> 5
    const r = end >>> 5
    const s = a >>> 5
    const n = r - l
    if (!(a & 31)) {
      for (let i = 0; i < n; i++) {
        this._bits[l + i] ^= other._bits[s + i]
      }
    } else {
      const hi = a & 31
      const lo = 32 - hi
      for (let i = 0; i < n; i++) {
        this._bits[l + i] ^= (other._bits[s + i] >>> hi) | (other._bits[1 + s + i] << lo)
      }
    }
  }

  /**
   * 类似 {@link Uint8Array.prototype.set}，如果超出赋值范围，抛出异常.
   * @param other 要赋值的位集.
   * @param offset 从哪里开始赋值.
   */
  set(other: BitSet, offset = 0): void {
    let left = offset
    let right = offset + other._n
    if (right > this._n) {
      throw new RangeError(`offset + other._n must be less than or equal to ${this._n}`)
    }
    let a = 0
    let b = other._n
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

  /**
   * @param newSize 拷贝后的位集大小.默认为原位集大小.
   */
  copy(newSize?: number): BitSet {
    if (newSize !== void 0) return this._copyAndResize(newSize)
    const res = new BitSet(this._n)
    res._bits.set(this._bits)
    return res
  }

  bitLength(): number {
    return this._lastIndexOfOne() + 1
  }

  resize(size: number): void {
    const newBits = new Uint32Array((size + 31) >>> 5)
    newBits.set(this._bits.subarray(0, newBits.length))
    const remainingBits = size & 31
    if (remainingBits) {
      const mask = (1 << remainingBits) - 1
      newBits[newBits.length - 1] &= mask
    }
    this._bits = newBits
    this._n = size
  }

  expand(size: number): void {
    if (size <= this._n) return
    this.resize(size)
  }

  toString(): string {
    const sb: string[] = []
    for (let i = 0; i < this._bits.length; i++) {
      // eslint-disable-next-line newline-per-chained-call
      let bits = this._bits[i].toString(2).padStart(32, '0').split('').reverse().join('')
      if (i === this._bits.length - 1) {
        bits = bits.slice(0, this._n - (i << 5))
      }
      sb.push(bits)
    }
    return sb.join('')
  }

  toSet(): Set<number> {
    const set = new Set<number>()
    this.forEach(i => {
      set.add(i)
    })
    return set
  }

  /**
   * 遍历所有 1 的位置.
   */
  forEach(callback: (value: number) => boolean | void): void {
    this._bits.forEach((v, i) => {
      // !注意结束条件是v!==0 而不是v>0
      for (; v; v &= v - 1) {
        const j = (i << 5) | BitSet._trailingZeros32(v)
        if (callback(j)) return
      }
    })
  }

  get size(): number {
    return this._n
  }

  get bits(): Uint32Array {
    return this._bits
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

  private _copyAndResize(size: number): BitSet {
    const res = new BitSet(size)
    res._bits.set(this._bits.subarray(0, res._bits.length))
    const remainingBits = size & 31
    if (remainingBits) {
      const mask = (1 << remainingBits) - 1
      res._bits[res._bits.length - 1] &= mask
    }
    return res
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
  // lsh/rsh demo
  const newA = new BitSet(1000)
  newA.add(20)
  newA.lsh(173)
  console.log(newA.toSet())

  const newB = new BitSet(100)
  newB.add(76)
  newB.add(80)
  newB.add(33)
  newB.add(10)
  console.log(newB.toSet())
  newB.rsh(17)
  console.log(newB.toSet(), newB.has(63))

  // // onesCountRange demo
  // const n = 330
  // const bs = new BitSet(n)
  // const nums = new Uint8Array(n)
  // for (let i = 0; i < 10; i++) {
  //   if (Math.random() > 0.5) {
  //     bs.add(i)
  //     nums[i] = 1
  //   }
  // }
  // // onesCountRange demo
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

  // const set = new BitSet(33)
  // set.fill(1)
  // set.fill(1)
  // console.log(set.has(34))
  // set.discard(1)
  // console.log(set.slice(0, 33).toString())

  // console.log(set.toString())
  // const mm = new BitSet(3)
  // // mm.add(0)

  // console.log(set.toString())

  // console.log('set')
  // //  1001 -> 01
  // const bs1 = BitSet.from('1001')
  // const bs2 = BitSet.from('01')
  // console.log(bs1.toString())
  // bs1.add(1)
  // console.log(bs1.next(1))

  // console.log(bs1.toString())
  // console.log(bs1.prev(3))

  // const bs3 = new BitSet(50)
  // bs3.add(35)
  // console.log(bs3.prev(40))

  // console.time('prev')
  // const big = new BitSet(1e5)
  // for (let i = 0; i < 1e5; i++) {
  //   big.prev(1e5)
  // }
  // console.timeEnd('prev')

  // console.time('expand')
  // const big2 = new BitSet(1e5)
  // for (let i = 1e5; i < 2e5; i++) {
  //   big2.expand(i)
  // }
  // console.timeEnd('expand')

  // console.log('test set')
  // const set1 = new BitSet(100)
  // set1.add(1)
  // const set2 = new BitSet(60)
  // set2.add(20)
  // set2.add(50)
  // set1.set(set2, 18)
  // console.log(set1.toSet())

  // const slice = new BitSet(0)
}
