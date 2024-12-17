/* eslint-disable no-param-reassign */

/**
 * 固定大小的环形缓冲区, 可用于双端队列、固定大小滑动窗口等场景.
 */
class RingBuffer<T> {
  private readonly _capacity: number
  private readonly _values: T[]
  private _start = 0
  private _end = 0
  private _size = 0

  constructor(capacity: number) {
    if (capacity <= 0) throw new Error('capacity must be greater than 0')
    this._capacity = capacity
    this._values = Array(capacity)
  }

  append(value: T): void {
    if (this.full()) this.popleft()
    this._values[this._end] = value
    this._end++
    if (this._end >= this._capacity) this._end = 0
    this._size++
  }

  appendLeft(value: T): void {
    if (this.full()) this.pop()
    this._start--
    if (this._start < 0) this._start = this._capacity - 1
    this._values[this._start] = value
    this._size++
  }

  pop(): T | undefined {
    if (this.empty()) return undefined
    this._end--
    if (this._end < 0) this._end = this._capacity - 1
    const value = this._values[this._end]
    this._size--
    return value
  }

  popleft(): T | undefined {
    if (this.empty()) return undefined
    const value = this._values[this._start]
    this._start++
    if (this._start >= this._capacity) this._start = 0
    this._size--
    return value
  }

  head(): T | undefined {
    if (this.empty()) return undefined
    return this._values[this._start]
  }

  tail(): T | undefined {
    if (this.empty()) return undefined
    const index = this._end - 1
    return this._values[index < 0 ? this._capacity - 1 : index]
  }

  at(index: number): T | undefined {
    if (index < 0) index += this._size
    if (index < 0 || index >= this._size) return undefined
    index += this._start
    if (index >= this._capacity) index -= this._capacity
    return this._values[index]
  }

  clear() {
    this._start = 0
    this._end = 0
    this._size = 0
  }

  empty(): boolean {
    return this._size === 0
  }

  full(): boolean {
    return this._size === this._capacity
  }

  size(): number {
    return this._size
  }

  forEach(consumer: (value: T) => void): void {
    let ptr = this._start
    for (let i = 0; i < this._size; i++) {
      consumer(this._values[ptr])
      ptr++
      if (ptr >= this._capacity) ptr = 0
    }
  }

  toString(): string {
    const sb: string[] = []
    this.forEach(x => sb.push(JSON.stringify(x)))
    return `RingBuffer{${sb.join(',')}}`
  }
}

export { RingBuffer }

if (require.main === module) {
  class MyCircularDeque {
    private readonly _k: number
    private readonly _buffer: RingBuffer<number>

    constructor(k: number) {
      this._k = k
      this._buffer = new RingBuffer(k)
    }

    insertFront(value: number): boolean {
      if (this.isFull()) return false
      this._buffer.appendLeft(value)
      return true
    }

    insertLast(value: number): boolean {
      if (this.isFull()) return false
      this._buffer.append(value)
      return true
    }

    deleteFront(): boolean {
      if (this.isEmpty()) return false
      this._buffer.popleft()
      return true
    }

    deleteLast(): boolean {
      if (this.isEmpty()) return false
      this._buffer.pop()
      return true
    }

    getFront(): number {
      return this._buffer.head() ?? -1
    }

    getRear(): number {
      return this._buffer.tail() ?? -1
    }

    isEmpty(): boolean {
      return this._buffer.empty()
    }

    isFull(): boolean {
      return this._buffer.full()
    }
  }
}
