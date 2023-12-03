/* eslint-disable @typescript-eslint/no-non-null-assertion */

/**
 * 单调队列维护滑动窗口最小值.
 */
class MonoQueue<T> {
  minQueue: SimpleQueue<T> = new SimpleQueue()
  private _len = 0
  private readonly _minQueueCount: SimpleQueue<number> = new SimpleQueue()
  private readonly _compareFn: (a: T, b: T) => number

  constructor(compareFn: (a: T, b: T) => number) {
    this._compareFn = compareFn
  }

  append(value: T): this {
    let count = 1
    while (this.minQueue.length && this._compareFn(this.minQueue.back()!, value) > 0) {
      this.minQueue.pop()
      count += this._minQueueCount.pop()!
    }
    this.minQueue.append(value)
    this._minQueueCount.append(count)
    this._len++
    return this
  }

  popleft(): void {
    if (!this._len) return
    this._minQueueCount.set(0, this._minQueueCount.front()! - 1)
    if (!this._minQueueCount.front()!) {
      this.minQueue.popleft()
      this._minQueueCount.popleft()
    }
    this._len--
  }

  head(): T | undefined {
    return this.minQueue.front()
  }

  toString(): string {
    const res: string[] = []
    for (let i = 0; i < this.minQueue.length; ++i) {
      res.push(JSON.stringify({ value: this.minQueue.at(i), count: this._minQueueCount.at(i) }))
    }
    return `MonoQueue[${res.join(', ')}]`
  }

  get min(): T | undefined {
    return this.minQueue.front()
  }

  get length(): number {
    return this._len
  }
}

class SimpleQueue<T> {
  private readonly _queue: T[] = []
  private _left = 0

  append(value: T): this {
    this._queue.push(value)
    return this
  }

  pop(): T | undefined {
    return this._queue.pop()
  }

  popleft(): T | undefined {
    if (!this.length) return undefined
    return this._queue[this._left++]
  }

  at(index: number): T | undefined {
    return this._queue[index + this._left]
  }

  set(index: number, value: T): void {
    this._queue[index + this._left] = value
  }

  toString(): string {
    const res: string[] = []
    for (let i = 0; i < this.length; ++i) {
      res.push(JSON.stringify(this.at(i)))
    }
    return `[${res.join(', ')}]`
  }

  front(): T | undefined {
    return this._queue[this._left]
  }

  back(): T | undefined {
    return this._queue[this._queue.length - 1]
  }

  get length(): number {
    return this._queue.length
  }
}

export { MonoQueue }

if (require.main === module) {
  const queue = new SimpleQueue<number>()
  queue.append(1)
  console.log(queue.toString())
  const monoQueue = new MonoQueue<number>((a, b) => a - b)
  monoQueue.append(2)
  monoQueue.append(1)
  console.log(monoQueue.toString())
  console.log(monoQueue.minQueue.toString())

  // https://leetcode.cn/problems/minimum-number-of-coins-for-fruits/description/
}
