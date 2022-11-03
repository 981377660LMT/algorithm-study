/* eslint-disable no-useless-constructor */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
// https://javaguide.cn/cs-basics/data-structure/bloom-filter/

import assert from 'assert'

import { BitSet } from './BitSet/BitSet'

class Hasher {
  private static _genHashCode(str: string): number {
    let hash = 0
    for (const char of str) {
      hash = (hash << 5) - hash + char.codePointAt(0)!
    }
    return hash | 0 // convert to int32
  }

  constructor(private capacity: number, private seed: number) {}

  genHash(value: string): number {
    const hashCode = Hasher._genHashCode(value)
    return Math.abs((this.seed * (hashCode ^ (hashCode >>> 16))) & (this.capacity - 1))
  }
}

interface IBloomFilter {
  add(value: string): this
  has(value: string): boolean
}

class BloomFilter implements IBloomFilter {
  private static readonly _DEFAULT_SIZE: number = 1 << 24
  private static readonly _SEEDS = [3, 13, 46, 71, 91, 134]
  private readonly _bitset = new BitSet(BloomFilter._DEFAULT_SIZE)
  private readonly _hashers: Hasher[] = []

  /**
   * 初始化多个包含 hash 函数的类的数组，每个类中的 hash 函数都不一样
   */
  constructor() {
    for (let i = 0; i < BloomFilter._SEEDS.length; i++) {
      this._hashers.push(new Hasher(BloomFilter._DEFAULT_SIZE, BloomFilter._SEEDS[i]))
    }
  }

  add(value: string): this {
    for (const hasher of this._hashers) {
      this._bitset.add(hasher.genHash(value))
    }
    return this
  }

  has(value: string): boolean {
    for (const hasher of this._hashers) {
      if (!this._bitset.has(hasher.genHash(value))) return false
    }
    return true
  }
}

export { BloomFilter }

if (require.main === module) {
  const bloomFilter = new BloomFilter()
  const value1 = 'cmnx'
  const value2 = 'c'
  const value3 = 'notIn'
  bloomFilter.add(value1).add(value2)
  assert.strictEqual(bloomFilter.has(value1), true)
  assert.strictEqual(bloomFilter.has(value2), true)
  assert.strictEqual(bloomFilter.has(value3), false)
}
