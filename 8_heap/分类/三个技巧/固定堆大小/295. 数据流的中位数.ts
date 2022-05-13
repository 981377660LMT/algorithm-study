import assert from 'assert'
import { PriorityQueue } from '../../../../2_queue/优先级队列'

// 我们的可以维护两个固定堆
// 大顶堆放最小的~~(n+1)/2个数 小顶堆放最大的n-~~(n+1)/2个数
// 左侧比右侧多1个，或者相同
class MedianFinder {
  private left: PriorityQueue<number> // 大顶堆
  private right: PriorityQueue<number> // 小顶堆

  constructor() {
    this.left = new PriorityQueue<number>((a, b) => b - a)
    this.right = new PriorityQueue<number>((a, b) => a - b)
  }

  addNum(num: number): void {
    // 往左边加：将右边小的移到左边
    if (this.left.length === this.right.length) {
      this.right.push(num)
      this.left.push(this.right.shift()!)
      // 往右边加：将左边大的移到右边
    } else {
      this.left.push(num)
      this.right.push(this.left.shift()!)
    }
  }

  findMedian(): number {
    return this.left.length === this.right.length
      ? (this.left.peek() + this.right.peek()) / 2
      : this.left.peek() //  左侧多1位
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
  console.dir(medianFiner, { depth: null })
  assert.strictEqual(medianFiner.findMedian(), 6)
  medianFiner.addNum(0)
  console.dir(medianFiner, { depth: null })
  assert.strictEqual(medianFiner.findMedian(), 5.5)
}

export {}
