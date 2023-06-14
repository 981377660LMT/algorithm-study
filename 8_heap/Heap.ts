/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 注意堆的索引从0开始，而线段树的索引从1开始
// 堆:root (root<<1)+1 (root<<1)+2
// 线段树: root root<<1 root<<1|1

import assert from 'assert'

type Comparator<T> = (a: T, b: T) => number

class Heap<E = number> {
  /**
   * 破坏性地合并两个堆，返回合并后的堆.采用启发式合并.
   */
  static mergeDestructively<E>(heap1: Heap<E>, heap2: Heap<E>): Heap<E> {
    if (heap1.size < heap2.size) {
      const tmp = heap1
      heap1 = heap2
      heap2 = tmp
    }
    for (let i = 0; i < heap2.size; i++) heap1.push(heap2._heap[i])
    return heap1
  }

  private _heap: E[]
  private readonly _comparator: Comparator<E>

  constructor()
  constructor(array: E[])
  constructor(comparator: Comparator<E>)
  constructor(array: E[], comparator: Comparator<E>)
  constructor(comparator: Comparator<E>, array: E[])
  constructor(arrayOrComparator1?: E[] | Comparator<E>, arrayOrComparator2?: E[] | Comparator<E>) {
    let defaultArray: E[] = []
    let defaultComparator = (a: E, b: E) => (a as unknown as number) - (b as unknown as number)

    if (arrayOrComparator1) {
      if (Array.isArray(arrayOrComparator1)) {
        defaultArray = arrayOrComparator1
      } else {
        defaultComparator = arrayOrComparator1
      }
    }

    if (arrayOrComparator2) {
      if (Array.isArray(arrayOrComparator2)) {
        defaultArray = arrayOrComparator2
      } else {
        defaultComparator = arrayOrComparator2
      }
    }

    this._comparator = defaultComparator
    this._heap = defaultArray
    if (this._heap.length > 1) this._heapify()
  }

  push(value: E): void {
    this._heap.push(value)
    this._pushUp(this._heap.length - 1)
  }

  pop(): E | undefined {
    if (this._heap.length <= 1) return this._heap.pop()
    const res = this._heap[0]
    this._heap[0] = this._heap.pop()!
    this._pushDown(0)
    return res
  }

  peek(): E | undefined {
    return this._heap[0]
  }

  clear(): void {
    this._heap = []
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
    while (parent >= 0 && this._comparator(this._heap[root], this._heap[parent]) < 0) {
      const tmp = this._heap[root]
      this._heap[root] = this._heap[parent]
      this._heap[parent] = tmp
      root = parent
      parent = (parent - 1) >> 1
    }
  }

  private _pushDown(root: number): void {
    // 还有孩子，即不是叶子节点
    const n = this._heap.length
    for (let left = (root << 1) | 1; left < n; left = (root << 1) | 1) {
      const right = left + 1
      let minIndex = root

      if (this._comparator(this._heap[left], this._heap[minIndex]) < 0) {
        minIndex = left
      }

      if (right < n && this._comparator(this._heap[right], this._heap[minIndex]) < 0) {
        minIndex = right
      }

      if (minIndex === root) return

      const tmp = this._heap[root]
      this._heap[root] = this._heap[minIndex]
      this._heap[minIndex] = tmp

      root = minIndex
    }
  }
}

export { Heap }

if (require.main === module) {
  const heap = new Heap()
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
