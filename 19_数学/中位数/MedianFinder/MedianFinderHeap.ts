/* eslint-disable no-else-return */

import { Heap } from '../../../8_heap/Heap'

/**
 * 对顶堆动态维护中位数.
 */
class MedianFinderHeap {
  private readonly _left = new Heap<number>((a, b) => b - a) // 大顶堆
  private readonly _right = new Heap<number>((a, b) => a - b) // 小顶堆

  add(num: number): void {
    if (this._left.size === this._right.size) {
      this._right.push(num)
      this._left.push(this._right.pop()!)
    } else {
      this._left.push(num)
      this._right.push(this._left.pop()!)
    }
  }

  /**
   * @param floor 是否向下取整.默认为true.
   */
  query(floor = true): number | undefined {
    if (this._left.size === 0) return undefined
    if (this._left.size === this._right.size) {
      const sum = this._left.peek()! + this._right.peek()!
      return floor ? Math.floor(sum / 2) : sum / 2
    } else {
      return this._left.peek()!
    }
  }

  get size(): number {
    return this._left.size + this._right.size
  }
}

export { MedianFinderHeap }
