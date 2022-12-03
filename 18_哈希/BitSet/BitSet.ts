/* eslint-disable generator-star-spacing */

import assert from 'assert'

import { bitCount32, trailingZero32 } from '../../19_数学/acwing专项训练/容斥原理/bitCount'

class BitSet implements Set<number> {
  /** 32 bit per bucket */
  private static readonly _BITS_PER_BUCKET = 1 << 5

  /** bit count of 1 */
  private _size = 0

  /** states */
  private _buckets: Uint32Array

  constructor(nbits: number) {
    this._buckets = new Uint32Array(Math.ceil(nbits / BitSet._BITS_PER_BUCKET))
  }

  add(index: number): this {
    this._ensureCapacity(index + 1)
    const row = Math.floor(index / BitSet._BITS_PER_BUCKET)
    const col = index % BitSet._BITS_PER_BUCKET
    if ((this._buckets[row] & (1 << col)) === 0) {
      this._buckets[row] |= 1 << col
      this._size++
    }
    return this
  }

  clear(): void {
    this._buckets.fill(0)
    this._size = 0
  }

  delete(index: number): boolean {
    this._ensureCapacity(index + 1)
    const row = Math.floor(index / BitSet._BITS_PER_BUCKET)
    const col = index % BitSet._BITS_PER_BUCKET
    if ((this._buckets[row] & (1 << col)) !== 0) {
      this._buckets[row] &= ~(1 << col)
      this._size--
      return true
    }
    return false
  }

  /**
   * @deprecated
   */
  forEach(): void {
    throw new Error(`Method ${this.forEach.name} not implemented.Use ${this.keys.name} instead.`)
  }

  has(index: number): boolean {
    this._ensureCapacity(index + 1)
    const row = Math.floor(index / BitSet._BITS_PER_BUCKET)
    const col = index % BitSet._BITS_PER_BUCKET
    return (this._buckets[row] & (1 << col)) !== 0
  }

  /**
   * @deprecated
   */
  entries(): IterableIterator<[number, number]> {
    throw new Error(`Method ${this.entries.name} not implemented.Use ${this.keys.name} instead.`)
  }

  *keys(): IterableIterator<number> {
    for (let i = 0; i < this._buckets.length; i++) {
      let state = this._buckets[i]
      const offset = i * BitSet._BITS_PER_BUCKET
      while (state > 0) {
        yield offset + trailingZero32(state)
        state &= state - 1
      }
    }
  }

  /**
   * @deprecated
   */
  values(): IterableIterator<number> {
    throw new Error(`Method ${this.values.name} not implemented.Use ${this.keys.name} instead.`)
  }

  and(other: BitSet): void {
    const minLength = Math.min(this._buckets.length, other._buckets.length)
    let newCount = 0
    for (let i = 0; i < minLength; i++) {
      this._buckets[i] &= other._buckets[i]
      newCount += bitCount32(this._buckets[i])
    }

    for (let i = minLength; i < this._buckets.length; i++) {
      this._buckets[i] = 0
    }

    this._size = newCount
  }

  or(other: BitSet): void {
    const maxLength = Math.max(this._buckets.length, other._buckets.length)
    this._ensureCapacity(maxLength * BitSet._BITS_PER_BUCKET)
    let newCount = 0
    for (let i = 0; i < maxLength; i++) {
      if (i < this._buckets.length && i < other._buckets.length) {
        this._buckets[i] |= other._buckets[i]
      } else if (i < this._buckets.length) {
        continue
      } else {
        this._buckets[i] = other._buckets[i]
      }
      newCount += bitCount32(this._buckets[i])
    }

    this._size = newCount
  }

  isSubset(other: BitSet): boolean {
    const maxLength = Math.max(this._buckets.length, other._buckets.length)
    for (let i = 0; i < maxLength; i++) {
      if (i < this._buckets.length && i < other._buckets.length) {
        if ((this._buckets[i] & other._buckets[i]) !== this._buckets[i]) {
          return false
        }
      } else if (i < this._buckets.length) {
        if (this._buckets[i] !== 0) {
          return false
        }
      }
    }

    return true
  }

  toString(): string {
    const keys = [...this.keys()].join(',')
    return `BitSet{${keys}}`
  }

  get size(): number {
    return this._size
  }

  *[Symbol.iterator](): IterableIterator<number> {
    yield* this.keys()
  }

  [Symbol.toStringTag] = 'BitSet'

  private _ensureCapacity(nbits: number): void {
    if (nbits > this._buckets.length * BitSet._BITS_PER_BUCKET) {
      this._resize(nbits)
    }
  }

  private _resize(nbits: number): void {
    let newLength = this._buckets.length << 1
    if (nbits > newLength * BitSet._BITS_PER_BUCKET) {
      newLength = Math.ceil(nbits / BitSet._BITS_PER_BUCKET)
    }

    const newBuckets = new Uint32Array(newLength)
    newBuckets.set(this._buckets)
    this._buckets = newBuckets
  }
}

export { BitSet }

if (require.main === module) {
  const set = new BitSet(100)
  set.add(1).add(2).add(3).add(101)
  // eslint-disable-next-line no-console
  console.log(set.toString())
  assert(set.has(1))
  assert(set.has(2))
  set.delete(2)
  assert(!set.has(2))
  assert(set.size === 3)

  const other = new BitSet(100)
  other.add(101).add(102).add(103)
  set.or(other)
  assert.strictEqual(set.size, 5)

  assert(other.isSubset(set))
  set.and(other)
  assert.strictEqual(set.toString(), other.toString())
  assert(set.isSubset(other))
  assert(other.isSubset(set))
}
