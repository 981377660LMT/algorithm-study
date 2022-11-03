/* eslint-disable generator-star-spacing */

import assert from 'assert'

import { trailingZero32 } from '../../19_数学/acwing专项训练/容斥原理/bitCount'

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

  get size(): number {
    return this._size
  }

  toString(): string {
    const keys = [...this.keys()].join(',')
    return `BitSet{${keys}}`
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
}
