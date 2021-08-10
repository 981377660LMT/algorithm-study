import assert from 'assert'
import { PriorityQueue } from '../../../../2_queue/todo优先级队列'

// 我们的可以维护两个固定堆
// 大顶堆放最小的~~(n+1)/2个数 小顶堆放最大的n-~~(n+1)/2个数
class MedianFiner {
  private small: PriorityQueue<number>
  private big: PriorityQueue<number>
  private length: number

  constructor() {
    this.small = new PriorityQueue<number>((a, b) => b - a)
    this.big = new PriorityQueue<number>((a, b) => a - b)
    this.length = 0
  }

  addNum(num: number): this {
    if (this.small.length === 0 || num < this.small.peek()) {
      this.small.push(num)
    } else {
      this.big.push(num)
    }

    // 两个堆的大小最多相差 1
    if (this.small.length < this.big.length) {
      this.small.push(this.big.shift()!)
    } else if (this.small.length > this.big.length + 1) {
      this.big.push(this.small.shift()!)
    }
    this.length++

    return this
  }

  findMedian(): number {
    if (this.length % 2 === 0) {
      return (this.small.peek() + this.big.peek()) / 2
    } else {
      return this.small.peek()
    }
  }
}

if (require.main === module) {
  const medianFiner = new MedianFiner()
  medianFiner.addNum(1).addNum(2).addNum(3)
  assert.strictEqual(medianFiner.findMedian(), 2)
  medianFiner.addNum(4)
  assert.strictEqual(medianFiner.findMedian(), 2.5)
}

export {}
