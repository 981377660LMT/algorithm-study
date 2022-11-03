/* eslint-disable generator-star-spacing */
/* eslint-disable no-param-reassign */

import assert from 'assert'

/**
 * 循环数组实现，慢数组动态扩容
 */
class ArrayDeque<E = number> {
  length: number
  private _head: number
  private _tail: number
  private readonly _capacity: number
  private readonly _data: E[]

  constructor(capacity?: number) {
    this._capacity = capacity ?? 1 << 30
    this._data = [] // 慢数组
    this._head = 0 // 从-1开始向前移动
    this._tail = -1 // 从0开始向后移动
    this.length = 0
  }

  at(index: number): E | undefined {
    if (index < 0) index += this.length
    const pos = (this._head + index + this._capacity) % this._capacity
    return this._data[pos]
  }

  // head前移
  unshift(value: E): boolean {
    if (this._isFull()) return false
    this._head = (this._head - 1 + this._capacity) % this._capacity
    this._data[this._head] = value
    this.length++
    return true
  }

  // tail后移
  push(value: E): boolean {
    if (this._isFull()) return false
    this._tail = (this._tail + 1 + this._capacity) % this._capacity
    this._data[this._tail] = value
    this.length++
    return true
  }

  // head后移
  shift(): E | undefined {
    if (this._isEmpty()) return undefined
    const front = this._front()
    this._head = (this._head + 1 + this._capacity) % this._capacity
    this.length--
    return front
  }

  // tail前移
  pop(): E | undefined {
    if (this._isEmpty()) return undefined
    const rear = this._back()
    this._tail = (this._tail - 1 + this._capacity) % this._capacity
    this.length--
    return rear
  }

  forEach(callback: (value: E, index: number, array: E[]) => void): void {
    let head = this._head
    const times = this.length
    for (let i = 0; i < times; i++) {
      callback(this._data[head], i, this._data)
      head = (head + 1 + this._capacity) % this._capacity
    }
  }

  private _front(): E | undefined {
    return this._isEmpty() ? undefined : this._data[(this._head + this._capacity) % this._capacity]
  }

  private _back(): E | undefined {
    return this._isEmpty() ? undefined : this._data[(this._tail + this._capacity) % this._capacity]
  }

  private _isEmpty(): boolean {
    return this.length === 0
  }

  private _isFull(): boolean {
    return this.length === this._capacity
  }

  *[Symbol.iterator]() {
    let { _head: head } = this
    const times = this.length
    for (let i = 0; i < times; i++) {
      yield this._data[head]
      head = (head + 1 + this._capacity) % this._capacity
    }
  }
}

export { ArrayDeque }

if (require.main === module) {
  const deque = new ArrayDeque<number>()

  deque.push(1)
  deque.unshift(2)

  for (const num of deque) {
    // eslint-disable-next-line no-console
    console.log(num)
  }

  assert.strictEqual(deque.at(0), 2)
  assert.strictEqual(deque.at(1), 1)
  assert.strictEqual(deque.at(2), void 0)
  assert.strictEqual(deque.at(-1), 1)
  assert.strictEqual(deque.at(-2), 2)
  assert.strictEqual(deque.at(-3), void 0)
}
