/* eslint-disable max-len */

class BitSetBigInt {
  static fromSet(set: Set<number>): BitSetBigInt {
    const res = new BitSetBigInt()
    for (const v of set) {
      res.add(v)
    }
    return res
  }

  private static readonly _V0 = 0n
  private static readonly _V1 = 1n
  private static readonly _V32 = 32n
  private static readonly _MASK32 = (1n << 32n) - 1n

  private _set = BitSetBigInt._V0

  add(value: number): this {
    this._set |= BitSetBigInt._V1 << BigInt(value)
    return this
  }

  has(value: number): boolean {
    return ((this._set >> BigInt(value)) & BitSetBigInt._V1) === BitSetBigInt._V1
  }

  delete(value: number): boolean {
    const mask = BitSetBigInt._V1 << BigInt(value)
    if ((this._set & mask) === BitSetBigInt._V0) {
      return false
    }
    this._set &= ~mask
    return true
  }

  flip(value: number): this {
    this._set ^= BitSetBigInt._V1 << BigInt(value)
    return this
  }

  clear(): void {
    this._set = BitSetBigInt._V0
  }

  lsh(k: number): this {
    this._set <<= BigInt(k)
    return this
  }

  rsh(k: number): this {
    this._set >>= BigInt(k)
    return this
  }

  toSet(): Set<number> {
    const set = new Set<number>()
    let cur = this._set
    let div = 0
    while (cur) {
      const offset = div << 5
      let num = Number(cur & BitSetBigInt._MASK32)
      while (num) {
        const i = 31 - Math.clz32(num & -num)
        set.add(offset | i)
        num ^= 1 << i
      }
      cur >>= BitSetBigInt._V32
      div++
    }
    return set
  }

  forEach(callback: (value: number) => boolean | void): void {
    let cur = this._set
    let div = 0
    while (cur) {
      const offset = div << 5
      let num = Number(cur & BitSetBigInt._MASK32)
      while (num) {
        const i = 31 - Math.clz32(num & -num)
        if (callback(offset | i)) return
        num ^= 1 << i
      }
      cur >>= BitSetBigInt._V32
      div++
    }
  }

  size(): number {
    let res = 0
    let cur = this._set
    while (cur) {
      const num = Number(cur & BitSetBigInt._MASK32)
      res += 32 - Math.clz32(num)
      cur >>= BitSetBigInt._V32
    }
    return res
  }

  equals(other: BitSetBigInt): boolean {
    return this._set === other._set
  }

  isSubsetOf(other: BitSetBigInt): boolean {
    return (this._set & other._set) === this._set
  }

  isSupersetOf(other: BitSetBigInt): boolean {
    return (this._set & other._set) === other._set
  }

  copy(): BitSetBigInt {
    const res = new BitSetBigInt()
    res._set = this._set
    return res
  }

  or(other: BitSetBigInt): BitSetBigInt {
    const res = this.copy()
    res._set |= other._set
    return res
  }

  ior(other: BitSetBigInt): this {
    this._set |= other._set
    return this
  }

  and(other: BitSetBigInt): BitSetBigInt {
    const res = this.copy()
    res._set &= other._set
    return res
  }

  iand(other: BitSetBigInt): this {
    this._set &= other._set
    return this
  }

  xor(other: BitSetBigInt): BitSetBigInt {
    const res = this.copy()
    res._set ^= other._set
    return res
  }

  ixor(other: BitSetBigInt): this {
    this._set ^= other._set
    return this
  }

  not(): BitSetBigInt {
    const res = this.copy()
    res._set = ~res._set
    return res
  }

  inot(): this {
    this._set = ~this._set
    return this
  }
}

// 3181. 执行操作可获得的最大总奖励 II
// leetcode.cn/problems/maximum-total-reward-using-operations-ii/solutions/2805488/typescript-bitsetyou-hua-ke-xing-xing-01-mss1/
// TypeScript bitset优化可行性01背包问题
// BigInt 大数模拟
const maxTotalReward = (rewardValues: number[]): number => {
  rewardValues = [...new Set(rewardValues)].sort((a, b) => a - b)
  const b1 = 1n
  let res = b1
  rewardValues.forEach(v => {
    const bv = BigInt(v)
    const low = ((b1 << bv) - b1) & res
    res |= low << bv
  })
  return res.toString(2).length - 1
}

export { BitSetBigInt }

if (require.main === module) {
  const set = new BitSetBigInt()
  for (let i = 0; i < 500; i++) {
    set.add(i)
  }
  console.log(set.toSet())
  console.time('toSet')
  for (let i = 0; i < 1e5; i++) {
    // set.flip(i)
    set.size()
  }
  console.timeEnd('toSet')
}
