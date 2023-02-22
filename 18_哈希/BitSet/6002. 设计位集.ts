/* eslint-disable no-param-reassign */

import { BitSet } from './BitSet'

class Bitset {
  private readonly _cap: number
  private readonly _bitSet: BitSet

  // little-endian
  constructor(capacity: number) {
    this._cap = capacity
    this._bitSet = new BitSet(capacity)
  }

  fix(index: number): void {
    this._bitSet.add(index)
  }

  unfix(index: number): void {
    this._bitSet.discard(index)
  }

  flip(): void {
    this._bitSet.flipRange(0, this._cap)
  }

  all(): boolean {
    return this._bitSet.allOne(0, this._cap)
  }

  one(): boolean {
    return !this._bitSet.allZero(0, this._cap)
  }

  count(): number {
    return this._bitSet.onesCount(0, this._cap)
  }

  toString(): string {
    return this._bitSet.toString()
  }
}

export {}
