/* eslint-disable no-constant-condition */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 注意堆的索引从0开始，而线段树的索引从1开始
// 堆:root (root<<1)+1 (root<<1)+2
// 线段树: root root<<1 root<<1|1
//
// !1.比较器最好用less，而不是compare，因为less更块.
// n=2e7 nlogn
// less: 830ms
// compare: 900ms
// !2.更快的写法是不用class, 直接用函数.
// class: 830ms
// function: 800ms

import assert from 'assert'

class Heap<T> {
  /**
   * 破坏性地合并两个堆，返回合并后的堆.采用启发式合并.
   */
  static mergeDestructively<E>(heap1: Heap<E>, heap2: Heap<E>): Heap<E> {
    if (heap1.size < heap2.size) {
      const tmp = heap1
      heap1 = heap2
      heap2 = tmp
    }
    for (let i = 0; i < heap2.size; i++) heap1.push(heap2._data[i])
    return heap1
  }

  private readonly _data: T[]
  private readonly _less: (a: T, b: T) => boolean

  constructor({ data, less }: { data: T[]; less: (a: T, b: T) => boolean }) {
    this._data = data.slice()
    this._less = less
    if (this._data.length > 1) this._heapify()
  }

  push(value: T): void {
    this._data.push(value)
    this._up(this._data.length - 1)
  }

  pop(): T {
    if (!this._data.length) throw new Error('pop from an empty heap')
    const n = this._data.length - 1
    this._swap(0, n)
    this._down(0, n)
    return this._data.pop()!
  }

  top(): T {
    if (!this._data.length) throw new Error('top from an empty heap')
    return this._data[0]
  }

  /**
   * 弹出并返回堆顶，同时将 value 入堆.
   * `pop` + `push` 的快速版本.
   */
  replace(value: T): T {
    if (!this._data.length) throw new Error('replace from an empty heap')
    const top = this._data[0]
    this._data[0] = value
    this._fix(0)
    return top
  }

  /**
   * 先将 value 入堆，然后弹出并返回堆顶.
   * `push` + `pop` 的快速版本.
   */
  pushPop(value: T): T {
    if (this._data.length && this._less(this._data[0], value)) {
      const tmp = this._data[0]
      this._data[0] = value
      value = tmp
      this._fix(0)
    }
    return value
  }

  clear(): void {
    this._data.length = 0
  }

  get size(): number {
    return this._data.length
  }

  private _heapify(): void {
    const n = this._data.length
    for (let i = (n >>> 1) - 1; ~i; i--) {
      this._down(i, n)
    }
  }

  private _up(j: number): void {
    const { _data, _less } = this
    while (j) {
      const i = (j - 1) >>> 1 // parent
      if (i === j || !_less(_data[j], _data[i])) break
      this._swap(i, j)
      j = i
    }
  }

  private _down(i0: number, n: number): boolean {
    const { _data, _less } = this
    let i = i0
    while (true) {
      const j1 = (i << 1) | 1
      if (j1 >= n || j1 < 0) break
      let j = j1
      const j2 = j1 + 1
      if (j2 < n && _less(_data[j2], _data[j1])) j = j2
      if (!_less(_data[j], _data[i])) break
      this._swap(i, j)
      i = j
    }
    return i > i0
  }

  private _fix(i: number): void {
    if (!this._down(i, this._data.length)) this._up(i)
  }

  private _swap(i: number, j: number): void {
    const tmp = this._data[i]
    this._data[i] = this._data[j]
    this._data[j] = tmp
  }
}

export { Heap }

if (require.main === module) {
  const heap = new Heap({ data: [], less: (a: number, b: number) => a < b })
  heap.push(1)
  heap.push(8)
  heap.push(3)
  heap.push(5)
  assert.strictEqual(heap.pop(), 1)
  assert.strictEqual(heap.pop(), 3)
  assert.strictEqual(heap.pop(), 5)
  assert.strictEqual(heap.pop(), 8)

  const N = 1e7
  console.time('Heap')
  const heap2 = new Heap<number>({ data: [], less: (a, b) => a < b })
  for (let i = 0; i < N; i++) {
    heap2.push(i)
  }
  for (let i = 0; i < N; i++) {
    heap2.pop()
  }
  console.timeEnd('Heap')
}
