/**
 * @description 一个动态大小的位集合。由 BigInt 实现，占用空间很小，但独写性能不如数组。
 * @see https://github.com/harttle/contest.js/blob/master/bitset.ts
 */
class BitSet {
  private readonly cap: number
  private val: bigint

  /**
   *
   * @param val
   * @param cap  默认容量为 Infinity
   */
  constructor(val: bigint | string | number | number[] = 0, cap = Infinity) {
    this.cap = cap
    this.val = 0n

    if (typeof val === 'number' || typeof val === 'bigint') {
      this.val = BigInt(val)
    } else if (typeof val === 'string') {
      for (let i = val.length - 1, bit = 1n; i >= 0; i--, bit <<= 1n) {
        if (val[i] !== '0') this.val |= bit
      }
    } else if (val instanceof Array) {
      for (let i = 0, bit = 1n; i < val.length; i++, bit <<= 1n) {
        if (val[i] !== 0) this.val |= bit
      }
    }

    // -1n 即全1
    const mask = cap === Infinity ? -1n : (1n << BigInt(cap)) - 1n
    this.val &= mask
  }

  get capacity(): number {
    return this.cap
  }

  add(value: number): void {
    this.set(value, 1)
  }

  has(value: number): boolean {
    return this.get(value) === 1
  }

  delete(value: number): void {
    this.set(value, 0)
  }

  /**
   *
   * @returns 返回集合中 1 的位数
   */
  count(): number {
    let count = 0
    for (let bit = 1n; bit <= this.val; bit <<= 1n) if ((bit & this.val) !== 0n) count++
    return count
  }

  toString(): string {
    if (this.cap <= 0) return ''
    if (this.val === 0n) return '0'
    let ans = ''
    let last = null
    for (let val = this.val; val !== 0n && val !== last; last = val, val >>= 1n) {
      ans = String(val & 1n) + ans
    }
    return this.cap === Infinity ? ans : ans.padStart(this.cap, '0')
  }

  shift(len: number): BitSet {
    return new BitSet(this.val << BigInt(len), this.cap)
  }

  unshift(len: number): BitSet {
    return new BitSet(this.val >> BigInt(len), this.cap)
  }

  and(rhs: BitSet | number | bigint | string): BitSet {
    if (!(rhs instanceof BitSet)) rhs = new BitSet(rhs)
    return new BitSet(this.val & rhs.val, Math.max(this.cap, rhs.cap))
  }

  nor(): BitSet {
    return new BitSet(~this.val, this.cap)
  }

  or(rhs: BitSet | number | bigint | string): BitSet {
    if (!(rhs instanceof BitSet)) rhs = new BitSet(rhs)
    return new BitSet(this.val | rhs.val, Math.max(this.cap, rhs.cap))
  }

  xor(rhs: BitSet | number | bigint | string): BitSet {
    if (!(rhs instanceof BitSet)) rhs = new BitSet(rhs)
    return new BitSet(this.val ^ rhs.val, Math.max(this.cap, rhs.cap))
  }

  /**
   *
   * @param index
   * @param value
   * @description 把下标为 i 的位设置位 val
   */
  private set(index: number, value: 1 | 0): void {
    if (value === 1) this.val |= 1n << BigInt(index)
    else this.val &= ~(1n << BigInt(index))
  }

  /**
   *
   * @param index
   * @returns 返回下标为 i 的位的值
   */
  private get(index: number): 1 | 0 {
    return (this.val & (1n << BigInt(index))) !== 0n ? 1 : 0
  }
}

if (require.main === module) {
  const bitSet = new BitSet()
  bitSet.add(2)
  console.log(bitSet.count())
  bitSet.delete(1)
}

export { BitSet }
