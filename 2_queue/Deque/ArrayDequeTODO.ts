// golang的deque实现是怎样的
// queue是怎么实现的 参考golang的数据结构

/**
 * 慢数组实现的双端队列,比定长数组+重构更快.
 */
class Queue<E = number> {
  private _head = 0
  private _tail = 0
  private readonly _data: E[]

  constructor(initCapacity = 16, iterable: Iterable<E> = []) {
    this._data = Array(initCapacity)
    for (const x of iterable) {
      this.push(x)
    }
  }

  push(value: E): Queue<E> {
    this._data[this._tail++] = value
    return this
  }

  shift(): E | undefined {
    if (!this.length) return undefined
    const front = this._data[this._head++]
    return front
  }

  at(index: number): E | undefined {
    const n = this.length
    if (index < 0) index += n
    if (index < 0 || index >= n) return undefined
    return this._data[this._head + index]
  }

  get length(): number {
    return this._tail - this._head
  }
}
