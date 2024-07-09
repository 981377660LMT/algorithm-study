/* eslint-disable @typescript-eslint/no-non-null-assertion */
// https://judge.yosupo.jp/submission/109819
// 懒删除的堆 维护最大最小值
//
// Min-Map Heap (four heaps technique)
//
// Description:
//   A data structure for push/min/max/popmin/popmax.
//
// Algorithm:
//   Maintain two min heaps (minh, minp) an two max heaps (maxh, maxp).
//   The actual elements in the heap is (minh \cup maxh) \setminus (\minp \cup maxp).
//
// Complexity:
//   Amortized O(1) for push and top, and O(log n) for pop.
//

import { Heap } from './Heap'

/**
 * 维护最大和最小值的堆.
 * 注意内部比较元素相等使用的是 `===` 而不是 `deepEqual`.
 * @warning 速度不如`IntervalHeap`.
 */
class MinMaxHeap<T> {
  private _size = 0
  private readonly _minHeapH: Heap<T>
  private readonly _minHeapP: Heap<T>
  private readonly _maxHeapH: Heap<T>
  private readonly _maxHeapP: Heap<T>

  constructor(minHeapLessFunc: (a: T, b: T) => boolean, maxHeapLessFunc: (a: T, b: T) => boolean) {
    this._minHeapH = new Heap({ data: [], less: minHeapLessFunc })
    this._minHeapP = new Heap({ data: [], less: minHeapLessFunc })
    this._maxHeapH = new Heap({ data: [], less: maxHeapLessFunc })
    this._maxHeapP = new Heap({ data: [], less: maxHeapLessFunc })
  }

  push(value: T): void {
    this._minHeapH.push(value)
    this._maxHeapH.push(value)
    this._size++
  }

  popMin(): T | undefined {
    if (!this._size) return undefined
    this._normalize()
    const res = this._minHeapH.pop()!
    this._maxHeapP.push(res)
    this._size--
    return res
  }

  popMax(): T | undefined {
    if (!this._size) return undefined
    this._normalize()
    const res = this._maxHeapH.pop()!
    this._minHeapP.push(res)
    this._size--
    return res
  }

  get min(): T | undefined {
    this._normalize()
    return this._minHeapH.top()
  }

  get max(): T | undefined {
    this._normalize()
    return this._maxHeapH.top()
  }

  get size(): number {
    return this._size
  }

  private _normalize(): void {
    while (this._minHeapP.size && this._minHeapP.top() === this._minHeapH.top()) {
      this._minHeapP.pop()
      this._minHeapH.pop()
    }
    while (this._maxHeapP.size && this._maxHeapP.top() === this._maxHeapH.top()) {
      this._maxHeapP.pop()
      this._maxHeapH.pop()
    }
  }
}

export {}

if (require.main === module) {
  const pq = new MinMaxHeap<number>(
    (a, b) => a < b,
    (a, b) => b < a
  )
  pq.push(3)
  pq.push(1)
  pq.push(2)
  console.log(pq.min)
  console.log(pq.max)
  console.log(pq.popMin())
  console.log(pq.size)
  console.log(pq.popMax())
  console.log(pq.size)
  console.log(pq.popMin())
  console.log(pq.size)
  console.log(pq.popMin())

  // time
  const n = 5e5
  const pq2 = new MinMaxHeap<number>(
    (a, b) => a < b,
    (a, b) => b < a
  )

  console.time('MinMaxHeap')
  for (let i = 0; i < n; i++) {
    pq2.push(i)
    pq2.push(i)
  }
  for (let i = 0; i < n; i++) {
    pq2.popMax()
    pq2.popMin()
  }
  console.timeEnd('MinMaxHeap') // MinMaxHeap: 772.526ms
}
