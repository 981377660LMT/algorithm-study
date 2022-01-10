// 海量数据处理以及缓存穿透这两个场景让我认识了 布隆过滤器
// https://javaguide.cn/cs-basics/data-structure/bloom-filter/

import assert from 'assert'
import { BitSet } from './BitSet/BitSet'

class SimpleHasher {
  constructor(private capacity: number, private seed: number) {}

  private static genHashCode(str: string) {
    let hash = 0

    for (const char of str) {
      // hash = 31 * hash + char.codePointAt(0)!
      // 与乘以31相同
      hash = (hash << 5) - hash + char.codePointAt(0)!
    }

    return hash | 0 // Convert to 32bit integer
  }

  genHash(value: string): number {
    const hashCode = SimpleHasher.genHashCode(value)
    return Math.abs((this.seed * (hashCode ^ (hashCode >>> 16))) & (this.capacity - 1))
  }
}

interface IBloomFilter {
  add(value: string): void
  has(value: string): boolean
}

class BloomFilter implements IBloomFilter {
  private static DEFAULT_SIZE: number = 2 << 24
  private static SEEDS: number[] = [3, 13, 46, 71, 91, 134]
  private bitset: BitSet = new BitSet(0, BloomFilter.DEFAULT_SIZE)
  private hasherArray: SimpleHasher[] = []

  /**
   * 初始化多个包含 hash 函数的类的数组，每个类中的 hash 函数都不一样
   */
  constructor() {
    for (let i = 0; i < BloomFilter.SEEDS.length; i++) {
      this.hasherArray.push(new SimpleHasher(BloomFilter.DEFAULT_SIZE, BloomFilter.SEEDS[i]))
    }
  }

  add(value: string): void {
    for (const hasher of this.hasherArray) {
      this.bitset.add(hasher.genHash(value))
    }
  }

  has(value: string): boolean {
    for (const hasher of this.hasherArray) {
      if (!this.bitset.has(hasher.genHash(value))) return false
    }

    return true
  }
}

if (require.main === module) {
  const bloomFilter = new BloomFilter()
  const value1 = 'cmnx'
  const value2 = 'c'
  const value3 = 'notIn'
  bloomFilter.add(value1)
  bloomFilter.add(value2)
  assert.strictEqual(bloomFilter.has(value1), true)
  assert.strictEqual(bloomFilter.has(value2), true)
  assert.strictEqual(bloomFilter.has(value3), false)
}
