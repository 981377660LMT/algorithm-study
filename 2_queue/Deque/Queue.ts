import { Queue as Q1 } from 'datastructures-js'
import { Queue } from 'typescript-collections'

/* eslint-disable no-param-reassign */

/**
 * 内部使用头尾指针实现的双端队列.
 * 当入队次数超过`maxEnqueue`时，会使用慢数组,而不是重构.
 */
class FastQueue<E = number> {
  private _head = 0
  private _tail = 0
  private readonly _data: E[]

  /**
   * @param maxEnqueue 当入队次数超过该值时，会使用慢数组.
   */
  constructor(maxEnqueue = 16, iterable: Iterable<E> = []) {
    this._data = Array(maxEnqueue)
    this._head = 0
    this._tail = 0
    for (const x of iterable) {
      this.push(x)
    }
  }

  push(value: E): FastQueue<E> {
    this._data[this._tail++] = value
    return this
  }

  pop(): E | undefined {
    if (!this.length) return undefined
    const last = this._data[--this._tail]
    return last
  }

  shift(): E | undefined {
    if (!this.length) return undefined
    const front = this._data[this._head++]
    return front
  }

  front(): E | undefined {
    return this._data[this._head]
  }

  back(): E | undefined {
    return this._data[this._tail - 1]
  }

  at(index: number): E | undefined {
    const n = this.length
    if (index < 0) index += n
    if (index < 0 || index >= n) return undefined
    return this._data[this._head + index]
  }

  forEach(callbackfn: (value: E, index: number) => void): void {
    let count = 0
    for (let i = this._head; i < this._tail; i++) {
      callbackfn(this._data[i], count++)
    }
  }

  toString(): string {
    return this._data.slice(this._head, this._tail).toString()
  }

  get length(): number {
    return this._tail - this._head
  }
}

/**
 * 当元素个数超过一定值时，会重构数组.
 */
class Queue<E = number> {}

if (require.main === module) {
  // 每次删除头部元素并不是真的移除，而是标记其已经被移除
  // 即用数组+左右两个指针 替代stack、queue
  const queue = new FastQueue<number>(1e6)
  queue.push(1)
  queue.push(2)
  console.log(queue.shift())
  console.log(queue.at(-1))
  console.log(queue.at(-2))
  console.log(queue.at(0))
  console.log(queue.at(1))
  queue.forEach((x, i) => console.log(i, x))
  console.log(queue.toString())

  console.time('SimpleQueue')
  const n = 1e7
  for (let i = 0; i < n; i++) {
    queue.push(i)
  }
  for (let i = 0; i < n; i++) {
    queue.at(i)
  }
  for (let i = 0; i < n; i++) {
    queue.shift()
  }
  console.timeEnd('SimpleQueue') // 150ms
}

export { FastQueue }
