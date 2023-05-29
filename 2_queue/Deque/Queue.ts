/* eslint-disable no-param-reassign */

/**
 * 内部使用头尾指针实现的队列.
 * 当入队次数超过`maxEnqueue`时，会使用慢数组, 而不是重构.
 * @warning 由于慢数组的特性, 访问负数索引比正数索引慢很多, 所以不支持`unshift`操作.
 */
class QueueFast<E = number> {
  private _head = 0
  private _tail = 0
  private _data: E[]

  /**
   * @param maxEnqueue 当入队次数超过该值时，会使用慢数组.
   */
  constructor(maxEnqueue: number, iterable: Iterable<E> = []) {
    this._data = Array(maxEnqueue)
    for (const x of iterable) {
      this.push(x)
    }
  }

  push(value: E): this {
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
    return this.length ? this._data[this._head] : undefined
  }

  back(): E | undefined {
    return this.length ? this._data[this._tail - 1] : undefined
  }

  /**
   * 0<=index<len.
   */
  at(index: number): E | undefined {
    return this.length ? this._data[this._head + index] : undefined
  }

  forEach(callbackfn: (value: E, index: number) => void): void {
    let count = 0
    for (let i = this._head; i < this._tail; i++) {
      callbackfn(this._data[i], count++)
    }
  }

  /** 这里不将 {@link _data} 置为空数组. */
  clear(): void {
    this._head = 0
    this._tail = 0
  }

  clone(): QueueFast<E> {
    const res = new QueueFast<E>(0)
    res._head = this._head
    res._tail = this._tail
    res._data = this._data.slice()
    return res
  }

  toString(): string {
    const sb: string[] = []
    this.forEach(x => sb.push(JSON.stringify(x)))
    return `QueueFast{${sb.join(',')}}`
  }

  get length(): number {
    return this._tail - this._head
  }
}

class Queue_<E = number> {
  private _offset = 0
  private _data: E[]
  private readonly _shrinkToFit: boolean

  /**
   * @param shrinkToFit 出队后，如果剩余元素不足原数组长度的一半，则将数组缩小.默认为`false`.
   */
  constructor(iterable: Iterable<E> = [], shrinkToFit = false) {
    this._data = [...iterable]
    this._shrinkToFit = shrinkToFit
  }

  push(value: E): this {
    this._data.push(value)
    return this
  }

  pop(): E | undefined {
    if (!this.length) return undefined
    return this._data.pop()
  }

  shift(): E | undefined {
    if (!this.length) return undefined
    const front = this._data[this._offset++]
    if (this._shrinkToFit && this._offset * 2 > this._data.length) this._shrink()
    return front
  }

  front(): E | undefined {
    return this.length ? this._data[this._offset] : undefined
  }

  back(): E | undefined {
    return this.length ? this._data[this._data.length - 1] : undefined
  }

  /**
   * 0<=index<len.
   */
  at(index: number): E | undefined {
    return this.length ? this._data[this._offset + index] : undefined
  }

  forEach(callbackfn: (value: E, index: number) => void): void {
    let count = 0
    for (let i = this._offset; i < this._data.length; i++) {
      callbackfn(this._data[i], count++)
    }
  }

  toString(): string {
    const sb: string[] = []
    this.forEach(x => sb.push(JSON.stringify(x)))
    return `Queue{${sb.join(',')}}`
  }

  clear(): void {
    this._data = []
    this._offset = 0
  }

  clone(): Queue_<E> {
    return new Queue_(this._data.slice(this._offset), this._shrinkToFit)
  }

  get length(): number {
    return this._data.length - this._offset
  }

  private _shrink(): void {
    this._data = this._data.slice(this._offset)
    this._offset = 0
  }
}

export { QueueFast, Queue_ }

if (require.main === module) {
  const demo = new QueueFast<number>(16)
  demo.push(1)
  demo.push(2)
  console.log(demo.shift())
  console.log(demo.at(-1))
  console.log(demo.at(-2))
  console.log(demo.at(0))
  console.log(demo.at(1))
  demo.forEach((x, i) => console.log(i, x))
  console.log(demo.toString(), demo.clone().toString())

  const n = 1e7
  console.time('FastQueue')
  const fastQueue = new QueueFast<number>(n)
  for (let i = 0; i < n; i++) {
    fastQueue.push(i)
  }
  for (let i = 0; i < n; i++) {
    fastQueue.at(i)
  }
  for (let i = 0; i < n; i++) {
    fastQueue.shift()
  }
  console.timeEnd('FastQueue') // 115.754ms

  console.time('Queue')
  const queue = new Queue_<number>()
  for (let i = 0; i < n; i++) {
    queue.push(i)
  }
  for (let i = 0; i < n; i++) {
    queue.at(i)
  }
  for (let i = 0; i < n; i++) {
    queue.shift()
  }
  console.timeEnd('Queue') // 300ms

  class MyCircularQueue {
    private readonly _queue: QueueFast<number>
    private readonly _capacity: number
    constructor(k: number) {
      this._capacity = k
      this._queue = new QueueFast<number>(k)
    }

    enQueue(value: number): boolean {
      if (this._queue.length === this._capacity) return false
      this._queue.push(value)
      return true
    }

    deQueue(): boolean {
      if (!this._queue.length) return false
      this._queue.shift()
      return true
    }

    Front(): number {
      return this._queue.front() ?? -1
    }

    Rear(): number {
      return this._queue.back() ?? -1
    }

    isEmpty(): boolean {
      return !this._queue.length
    }

    isFull(): boolean {
      return this._queue.length === this._capacity
    }
  }

  /**
   * Your MyCircularQueue object will be instantiated and called as such:
   * var obj = new MyCircularQueue(k)
   * var param_1 = obj.enQueue(value)
   * var param_2 = obj.deQueue()
   * var param_3 = obj.Front()
   * var param_4 = obj.Rear()
   * var param_5 = obj.isEmpty()
   * var param_6 = obj.isFull()
   */
  // ;[
  //   'MyCircularQueue',
  //   'enQueue',
  //   'deQueue',
  //   'Front',
  //   'deQueue',
  //   'Front',
  //   'Rear',
  //   'enQueue',
  //   'isFull',
  //   'deQueue',
  //   'Rear',
  //   'enQueue'
  // ][([3], [7], [], [], [], [], [], [0], [], [], [], [3])]
  const obj = new MyCircularQueue(3)
  console.log(obj.enQueue(7))
  console.log(obj.deQueue())
  console.log(obj.Front())
  console.log(obj.deQueue())
  console.log(obj.Front())
  console.log(obj.Rear())
  console.log(obj.enQueue(0)) // 因该是true
  console.log(obj.isFull())
  console.log(obj.deQueue())
  console.log(obj.Rear())
  console.log(obj.enQueue(3))
}
