/* eslint-disable generator-star-spacing */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */

// 用两个栈模拟一个双端队列
// 当某个栈空的时候将另一个栈分一半给这个已经空的栈（即暴力重构）

class StackDeque<E> {
  static from<T>(iterable: Iterable<T>): StackDeque<T> {
    const deque = new StackDeque<T>()
    for (const item of iterable) {
      deque.append(item)
    }
    return deque
  }

  private _left: E[] = []
  private _right: E[] = []

  append(item: E): void {
    this._right.push(item)
  }

  appendLeft(item: E): void {
    this._left.push(item)
  }

  pop(): E | undefined {
    if (this._right.length > 0) {
      return this._right.pop()
    }
    const tmp: E[] = []
    const n = this._left.length
    for (let i = 0; i < n; i++) {
      tmp.push(this._left.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._left.push(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._right.push(tmp[i])
    }
    return this._right.pop()
  }

  popLeft(): E | undefined {
    if (this._left.length > 0) {
      return this._left.pop()
    }
    const tmp: E[] = []
    const n = this._right.length
    for (let i = 0; i < n; i++) {
      tmp.push(this._right.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._right.push(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._left.push(tmp[i])
    }
    return this._left.pop()
  }

  at(index: number): E | undefined {
    const n = this.size
    if (index < 0) {
      index += n
    }
    if (index < 0 || index >= n) {
      return undefined
    }
    const leftSize = this._left.length
    if (index < leftSize) {
      return this._left[leftSize - index - 1]
    }
    return this._right[index - leftSize]
  }

  reverse(): void {
    ;[this._left, this._right] = [this._right, this._left]
  }

  toString(): string {
    const items = [...this].map(v => JSON.stringify(v))
    return `StackDeque(${items.join(', ')})`
  }

  get size(): number {
    return this._left.length + this._right.length
  }

  *entries(): IterableIterator<[number, E]> {
    const leftSize = this._left.length
    const rightSize = this._right.length
    let count = 0
    for (let i = leftSize - 1; i >= 0; i--) {
      yield [count, this._left[i]]
      count++
    }
    for (let i = 0; i < rightSize; i++) {
      yield [count, this._right[i]]
      count++
    }
  }

  *[Symbol.iterator](): Iterator<E> {
    const leftSize = this._left.length
    const rightSize = this._right.length
    for (let i = leftSize - 1; i >= 0; i--) {
      yield this._left[i]
    }
    for (let i = 0; i < rightSize; i++) {
      yield this._right[i]
    }
  }
}

if (require.main === module) {
  const deque = new StackDeque<number>()
  deque.append(1)
  deque.append(2)
  deque.appendLeft(3)
  console.log(deque.toString())
  deque.reverse()
}

export { StackDeque }
