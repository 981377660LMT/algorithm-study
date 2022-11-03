/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 注意堆的索引从0开始，而线段树的索引从1开始
// 堆:root (root<<1)+1 (root<<1)+2
// 线段树: root root<<1 root<<1|1

import assert from 'assert'

type Comparator<T> = (a: T, b: T) => number

class Heap<E> {
  private readonly _heap: E[] = []
  private readonly _comparator: Comparator<E>

  constructor(comparator: Comparator<E>, array?: E[]) {
    this._comparator = comparator
    if (array) {
      this._heap = array
      this._heapify()
    }
  }

  push(value: E): void {
    this._heap.push(value)
    this._pushUp(this._heap.length - 1)
  }

  pop(): E | undefined {
    if (this._heap.length <= 1) return this._heap.pop()
    const returned = this._heap[0]
    const last = this._heap.pop()!
    this._heap[0] = last
    this._pushDown(0)
    return returned
  }

  peek(): E | undefined {
    return this._heap[0]
  }

  get size(): number {
    return this._heap.length
  }

  /**
   * 堆化的复杂度是 `O(n)`
   */
  private _heapify(): void {
    const n = this._heap.length
    for (let i = (n >> 1) - 1; ~i; i--) {
      this._pushDown(i)
    }
  }

  private _pushUp(root: number): void {
    let parent = (root - 1) >> 1
    while (parent >= 0 && this._comparator(this._heap[parent], this._heap[root]) > 0) {
      this._swap(parent, root)
      root = parent
      parent = (parent - 1) >> 1
    }
  }

  private _pushDown(root: number): void {
    // 还有孩子，即不是叶子节点
    while (true) {
      const left = (root << 1) | 1
      if (left >= this._heap.length) break

      const right = left + 1
      let minIndex = root

      if (this._comparator(this._heap[left], this._heap[minIndex]) < 0) {
        minIndex = left
      }

      if (
        right < this._heap.length &&
        this._comparator(this._heap[right], this._heap[minIndex]) < 0
      ) {
        minIndex = right
      }

      if (minIndex === root) return

      this._swap(root, minIndex)
      root = minIndex
    }
  }

  private _swap(i1: number, i2: number): void {
    ;[this._heap[i1], this._heap[i2]] = [this._heap[i2], this._heap[i1]]
  }
}

export { Heap }

if (require.main === module) {
  const heap = new Heap<number>((a, b) => a - b)
  heap.push(1)
  heap.push(8)
  heap.push(3)
  heap.push(5)
  assert.strictEqual(heap.pop(), 1)
  assert.strictEqual(heap.pop(), 3)
  assert.strictEqual(heap.pop(), 5)
  assert.strictEqual(heap.pop(), 8)
  assert.strictEqual(heap.pop(), undefined)
}
