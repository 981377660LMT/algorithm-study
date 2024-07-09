// 懒删除堆/可删除堆
// 堆的删除思路有两种:
// !1. 一种是懒删除，即查询时再实际删除元素；
// 这里的实现是懒删除.
// !2. 另一种是实时删除 index 处的元素
// 调用 heappush 会返回一个 *viPair 指针，记作 p
// 将 p 存于他处（如 slice 或 map），可直接在外部修改 p.v 后调用 fix(p.index)，从而做到修改堆中指定元素
// !调用 remove(p.index) 可以从堆中删除 p.v
// https://github.dev/EndlessCheng/codeforces-go/blob/6be496d4d93d667e718f7f3db5519139a5f17ddf/copypasta/heap.go#L94
// https://cs.opensource.google/go/go/+/refs/tags/go1.19.2:src/container/heap/heap.go
// !需要注意的是, 删除时比较元素使用的是`===`, 而不是deepEqual.

import { Heap } from './Heap'

type LessFunc<E> = (a: E, b: E) => boolean

/**
 * 懒删除堆/可删除堆.
 * 使用两个堆来实现,时间空间复杂度优于使用`Map`来记录元素的删除次数的实现.
 */
class ErasableHeap<E> {
  private readonly _data: Heap<E>
  private readonly _erased: Heap<E>
  private _size: number

  constructor()
  constructor(array: E[])
  constructor(less: LessFunc<E>)
  constructor(array: E[], less: LessFunc<E>)
  constructor(less: LessFunc<E>, array: E[])
  constructor(arrayOrLessFunc1?: E[] | LessFunc<E>, arrayOrLessFunc2?: E[] | LessFunc<E>) {
    let defaultArray: E[] = []
    let defaultLessFunc = (a: E, b: E) => a < b

    if (arrayOrLessFunc1) {
      if (Array.isArray(arrayOrLessFunc1)) {
        defaultArray = arrayOrLessFunc1
      } else {
        defaultLessFunc = arrayOrLessFunc1
      }
    }

    if (arrayOrLessFunc2) {
      if (Array.isArray(arrayOrLessFunc2)) {
        defaultArray = arrayOrLessFunc2
      } else {
        defaultLessFunc = arrayOrLessFunc2
      }
    }

    this._data = new Heap({ data: defaultArray, less: defaultLessFunc })
    this._erased = new Heap({ data: [], less: defaultLessFunc })
    this._size = this._data.size
  }

  push(value: E): void {
    this._data.push(value)
    this._normalize()
    this._size++
  }

  pop(): E | undefined {
    if (!this._size) return undefined
    const value = this._data.pop()
    this._normalize()
    this._size--
    return value
  }

  peek(): E | undefined {
    return this._size ? this._data.top() : undefined
  }

  /**
   * 删除堆中的元素`value`.
   * @warnings 删除前要保证堆中存在该元素，且删除时比较元素使用的是`===`, 而不是`deepEqual`.
   */
  remove(value: E): void {
    this._erased.push(value)
    this._normalize()
    this._size--
  }

  get size(): number {
    return this._size
  }

  private _normalize(): void {
    while (this._data.size && this._erased.size && this._data.top() === this._erased.top()) {
      this._data.pop()
      this._erased.pop()
    }
  }
}

/**
 * 懒删除堆/可删除堆.
 * 使用`Map`来记录元素的删除次数.
 * @deprecated
 */
class _ErasableHeap2<E> {
  private readonly _data: Heap<E>
  private readonly _erased: Map<E, number> = new Map()
  private _size: number

  constructor()
  constructor(array: E[])
  constructor(less: LessFunc<E>)
  constructor(array: E[], less: LessFunc<E>)
  constructor(less: LessFunc<E>, array: E[])
  constructor(arrayOrLessFunc1?: E[] | LessFunc<E>, arrayOrLessFunc2?: E[] | LessFunc<E>) {
    let defaultArray: E[] = []
    let defaultLessFunc = (a: E, b: E) => a < b

    if (arrayOrLessFunc1) {
      if (Array.isArray(arrayOrLessFunc1)) {
        defaultArray = arrayOrLessFunc1
      } else {
        defaultLessFunc = arrayOrLessFunc1
      }
    }

    if (arrayOrLessFunc2) {
      if (Array.isArray(arrayOrLessFunc2)) {
        defaultArray = arrayOrLessFunc2
      } else {
        defaultLessFunc = arrayOrLessFunc2
      }
    }

    this._data = new Heap({ data: defaultArray, less: defaultLessFunc })
    this._size = this._data.size
  }

  push(value: E): void {
    this._data.push(value)
    this._size++
  }

  pop(): E | undefined {
    this._expire()
    const res = this._data.pop()
    if (res !== undefined) this._size--
    return res
  }

  peek(): E | undefined {
    if (!this._size) return undefined
    this._expire()
    return this._data.top()
  }

  /**
   * 删除堆中的元素`value`.
   * @warnings 删除前要保证堆中存在该元素，且删除时比较元素使用的是`===`, 而不是`deepEqual`.
   */
  discard(value: E): void {
    this._erased.set(value, (this._erased.get(value) || 0) + 1)
    this._size--
  }

  get size(): number {
    return this._size
  }

  private _expire(): void {
    while (this._data.size) {
      const top = this._data.top()!
      const erasedCount = this._erased.get(top)
      if (!erasedCount) break
      this._data.pop()
      this._erased.set(top, erasedCount - 1)
    }
  }
}

export { ErasableHeap, ErasableHeap as RemovableHeap }

if (require.main === module) {
  const pq = new ErasableHeap<number>()
  pq.push(4)
  pq.push(1)
  pq.push(2)
  pq.push(3)
  pq.remove(2)

  while (pq.size) {
    console.log(pq.pop())
  }

  const N = 1e7

  {
    // test perf
    console.time('ErasableHeap')
    const pq2 = new ErasableHeap<number>()
    for (let i = 0; i < N; ++i) {
      pq2.push(i)
      pq2.peek()
    }
    for (let i = 0; i < N; ++i) {
      pq2.remove(i)
    }
    console.timeEnd('ErasableHeap') // 1e7: 1.34s
  }

  {
    console.time('ErasableHeap2')
    const pq3 = new _ErasableHeap2<number>()
    for (let i = 0; i < N; ++i) {
      pq3.push(i)
      pq3.peek()
    }
    for (let i = 0; i < N; ++i) {
      pq3.discard(i)
    }
    console.timeEnd('ErasableHeap2') // 1e7: 2.8s
  }
}
