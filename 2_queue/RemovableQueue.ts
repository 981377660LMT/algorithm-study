/* eslint-disable max-len */

interface IQueue<V> {
  readonly length: number
  push(value: V): void
  pop(): V | undefined
  top(): V | undefined
}

class RemovableQueue<V> {
  private readonly _queue: IQueue<V>
  private readonly _removedQueue: IQueue<V>
  private readonly _equals?: (v1: V, v2: V) => boolean

  constructor(createQueue: () => IQueue<V>, equals?: (v1: V, v2: V) => boolean) {
    this._queue = createQueue()
    this._removedQueue = createQueue()
    this._equals = equals
  }

  push(value: V): void {
    this._queue.push(value)
  }

  pop(): V | undefined {
    this._refresh()
    return this._queue.pop()
  }

  top(): V | undefined {
    this._refresh()
    return this._queue.top()
  }

  remove(value: V): void {
    this._removedQueue.push(value)
  }

  empty(): boolean {
    return this._removedQueue.length === 0
  }

  get length(): number {
    return this._queue.length - this._removedQueue.length
  }

  private _refresh(): void {
    if (!this._equals) {
      while (this._removedQueue.length && this._removedQueue.top() === this._queue.top()) {
        this._removedQueue.pop()
        this._queue.pop()
      }
    } else {
      while (this._removedQueue.length && this._equals(this._removedQueue.top()!, this._queue.top()!)) {
        this._removedQueue.pop()
        this._queue.pop()
      }
    }
  }
}

export {}
