/* eslint-disable @typescript-eslint/no-non-null-assertion */
/* eslint-disable generator-star-spacing */
/* eslint-disable no-param-reassign */

// class Deque
// !参考golang的双端队列实现,两个queue头尾拼接实现
// !left尾部 -> left头部 <-> right头部 <- right尾部

import assert from 'assert'

import { Queue_, QueueFast } from './Queue'

/**
 * 内部的队列使用头尾指针实现.
 * 当入队次数超过`maxEnqueue`时，会使用慢数组.
 * @deprecated
 */
class ArrayDequeFast<E = number> {
  private _left: QueueFast<E>
  private _right: QueueFast<E>

  /**
   * @param maxEnqueue 当入队次数超过该值时，会使用慢数组.
   */
  constructor(maxEnqueue: number, iterable: Iterable<E> = []) {
    this._left = new QueueFast(maxEnqueue)
    this._right = new QueueFast(maxEnqueue)
    for (const x of iterable) {
      this.push(x)
    }
  }

  push(value: E): this {
    this._right.push(value)
    return this
  }

  unshift(value: E): this {
    this._left.push(value)
    return this
  }

  pop(): E | undefined {
    return this._right.length ? this._right.pop() : this._left.shift()
  }

  shift(): E | undefined {
    return this._left.length ? this._left.pop() : this._right.shift()
  }

  forEach(callback: (value: E, index: number) => void): void {
    const ll = this._left.length
    for (let i = 0; i < ll; i++) callback(this._left.at(ll - 1 - i)!, i)
    const rl = this._right.length
    for (let i = 0; i < rl; i++) callback(this._right.at(i)!, ll + i)
  }

  /**
   * 0<=index<len.
   */
  at(index: number): E | undefined {
    const ll = this._left.length
    if (index < ll) return this._left.at(ll - 1 - index)
    return this._right.at(index - ll)
  }

  front(): E | undefined {
    return this._left.length ? this._left.back() : this._right.front()
  }

  back(): E | undefined {
    return this._right.length ? this._right.back() : this._left.front()
  }

  clear(): void {
    this._left.clear()
    this._right.clear()
  }

  clone(): ArrayDequeFast<E> {
    const res = new ArrayDequeFast<E>(0)
    res._left = this._left.clone()
    res._right = this._right.clone()
    return res
  }

  toString(): string {
    const sb: string[] = []
    this.forEach(x => sb.push(JSON.stringify(x)))
    return `ArrayDequeFast{${sb.join(',')}}`
  }

  empty(): boolean {
    return !this._left.length && !this._right.length
  }

  get length(): number {
    return this._left.length + this._right.length
  }
}

class ArrayDeque<E = number> {
  private _left: Queue_<E>
  private _right: Queue_<E>
  private _shrinkToFit: boolean

  /**
   * @param shrinkToFit 出队后，如果剩余元素不足原数组长度的一半，则将数组缩小.默认为`false`.
   */
  constructor(iterable: Iterable<E> = [], shrinkToFit = false) {
    this._left = new Queue_([], shrinkToFit)
    this._right = new Queue_([], shrinkToFit)
    for (const x of iterable) {
      this.push(x)
    }
    this._shrinkToFit = shrinkToFit
  }

  push(value: E): this {
    this._right.push(value)
    return this
  }

  unshift(value: E): this {
    this._left.push(value)
    return this
  }

  pop(): E | undefined {
    return this._right.length ? this._right.pop() : this._left.shift()
  }

  shift(): E | undefined {
    return this._left.length ? this._left.pop() : this._right.shift()
  }

  forEach(callback: (value: E, index: number) => void): void {
    const ll = this._left.length
    for (let i = 0; i < ll; i++) callback(this._left.at(ll - 1 - i)!, i)
    const rl = this._right.length
    for (let i = 0; i < rl; i++) callback(this._right.at(i)!, ll + i)
  }

  /**
   * 0<=index<len.
   */
  at(index: number): E | undefined {
    const ll = this._left.length
    if (index < ll) return this._left.at(ll - 1 - index)
    return this._right.at(index - ll)
  }

  front(): E | undefined {
    return this._left.length ? this._left.back() : this._right.front()
  }

  back(): E | undefined {
    return this._right.length ? this._right.back() : this._left.front()
  }

  clear(): void {
    this._left.clear()
    this._right.clear()
  }

  clone(): ArrayDeque<E> {
    const res = new ArrayDeque<E>([], this._shrinkToFit)
    res._left = this._left.clone()
    res._right = this._right.clone()
    return res
  }

  toString(): string {
    const sb: string[] = []
    this.forEach(x => sb.push(JSON.stringify(x)))
    return `ArrayDeque{${sb.join(',')}}`
  }

  empty(): boolean {
    return !this._left.length && !this._right.length
  }

  get length(): number {
    return this._left.length + this._right.length
  }
}

export { ArrayDequeFast, ArrayDeque }

if (require.main === module) {
  const deque = new ArrayDequeFast<number>(16)

  deque.push(1)
  deque.unshift(2)
  deque.unshift(3)
  deque.push(4)
  console.log(deque.toString(), deque.at(0), deque.at(1), deque.at(2))

  assert.strictEqual(deque.at(0), 3)
  assert.strictEqual(deque.at(1), 2)
  assert.strictEqual(deque.at(2), 1)
  assert.strictEqual(deque.back(), 4)
  assert.strictEqual(deque.front(), 3)

  // 641. 设计循环双端队列
  class MyCircularDeque {
    private readonly _capacity: number
    private readonly _queue: ArrayDequeFast<number>
    constructor(k: number) {
      this._capacity = k
      this._queue = new ArrayDequeFast(k)
    }

    insertFront(value: number): boolean {
      if (this._queue.length === this._capacity) return false
      this._queue.unshift(value)
      return true
    }

    insertLast(value: number): boolean {
      if (this._queue.length === this._capacity) return false
      this._queue.push(value)
      return true
    }

    deleteFront(): boolean {
      if (this._queue.length === 0) return false
      this._queue.shift()
      return true
    }

    deleteLast(): boolean {
      if (this._queue.length === 0) return false
      this._queue.pop()
      return true
    }

    getFront(): number {
      return this._queue.front() ?? -1
    }

    getRear(): number {
      return this._queue.back() ?? -1
    }

    isEmpty(): boolean {
      return this._queue.length === 0
    }

    isFull(): boolean {
      return this._queue.length === this._capacity
    }
  }

  /**
   * Your MyCircularDeque object will be instantiated and called as such:
   * var obj = new MyCircularDeque(k)
   * var param_1 = obj.insertFront(value)
   * var param_2 = obj.insertLast(value)
   * var param_3 = obj.deleteFront()
   * var param_4 = obj.deleteLast()
   * var param_5 = obj.getFront()
   * var param_6 = obj.getRear()
   * var param_7 = obj.isEmpty()
   * var param_8 = obj.isFull()
   */

  const n = 1e7
  console.time('FastDeque')
  const fdq = new ArrayDequeFast<number>(n)
  for (let i = 0; i < n; i++) {
    fdq.push(i)
  }
  for (let i = 0; i < n; i++) {
    fdq.shift()
  }
  for (let i = 0; i < n; i++) {
    fdq.unshift(i)
  }
  for (let i = 0; i < n; i++) {
    fdq.pop()
  }
  for (let i = 0; i < n; i++) {
    fdq.at(i)
  }
  console.timeEnd('FastDeque') // 235.25ms

  console.time('Deque')
  const dq = new ArrayDeque<number>()
  for (let i = 0; i < n; i++) {
    dq.push(i)
  }
  for (let i = 0; i < n; i++) {
    dq.shift()
  }
  for (let i = 0; i < n; i++) {
    dq.unshift(i)
  }
  for (let i = 0; i < n; i++) {
    dq.pop()
  }
  for (let i = 0; i < n; i++) {
    dq.at(i)
  }
  console.timeEnd('Deque') // 620.207ms
}
