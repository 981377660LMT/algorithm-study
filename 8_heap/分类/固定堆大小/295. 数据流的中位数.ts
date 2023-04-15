/* eslint-disable @typescript-eslint/no-non-null-assertion */

import assert from 'assert'
import { Heap } from '../../../Heap'

// 我们的可以维护两个固定堆
// 大顶堆放最小的(n+1)/2个数 小顶堆放最大的n-(n+1)/2个数
// 左侧比右侧多1个，或者相同
class MedianFinder {
  private readonly _left = new Heap<number>((a, b) => b - a) // 大顶堆
  private readonly _right = new Heap<number>((a, b) => a - b) // 小顶堆

  addNum(num: number): void {
    // 往左边加：将右边小的移到左边
    if (this._left.size === this._right.size) {
      this._right.push(num)
      this._left.push(this._right.pop()!)
      // 往右边加：将左边大的移到右边
    } else {
      this._left.push(num)
      this._right.push(this._left.pop()!)
    }
  }

  findMedian(): number {
    if (this._left.size === 0) return 0
    return this._left.size === this._right.size
      ? (this._left.peek()! + this._right.peek()!) / 2
      : this._left.peek()! //  左侧多1位
  }
}

if (require.main === module) {
  const medianFiner = new MedianFinder()
  medianFiner.addNum(6)
  medianFiner.addNum(10)
  medianFiner.addNum(2)
  medianFiner.addNum(6)
  assert.strictEqual(medianFiner.findMedian(), 6)
  medianFiner.addNum(5)
  assert.strictEqual(medianFiner.findMedian(), 6)
  medianFiner.addNum(0)
  assert.strictEqual(medianFiner.findMedian(), 5.5)
}

export {}
