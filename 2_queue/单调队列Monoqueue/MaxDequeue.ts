/* eslint-disable @typescript-eslint/no-explicit-any */
/* eslint-disable no-console */
/* eslint-disable no-param-reassign */
/* eslint-disable @typescript-eslint/no-non-null-assertion */
// 用两个栈模拟一个双端队列，注意到每个栈的前缀最小值是可以 O(1) 修改和询问的
// 当某个栈空的时候将另一个栈分一半给这个已经空的栈（即暴力重构）

interface IDeque<E = unknown> {
  append(item: E): void
  appendLeft(item: E): void
  pop(): E | undefined
  popLeft(): E | undefined
  at(index: number): E | undefined
  readonly size: number
}

interface IStack<E = unknown> {
  append(item: E): void
  pop(): E | undefined
  at(index: number): E | undefined
  readonly size: number
}

type CompareFunc<E> = (a: E, b: E) => number

class MaxDeque<E extends { value: V }, V = E['value']> implements IDeque<E> {
  private readonly _left: _MaxStack<E, V>
  private readonly _right: _MaxStack<E, V>
  private readonly _compareValue: CompareFunc<V>

  constructor(compareValue: CompareFunc<V> = (a: any, b: any) => a - b) {
    this._left = new _MaxStack(compareValue)
    this._right = new _MaxStack(compareValue)
    this._compareValue = compareValue
  }

  append(item: E): void {
    this._right.append(item)
  }

  appendLeft(item: E): void {
    this._left.append(item)
  }

  pop(): E | undefined {
    if (this._right.size > 0) {
      return this._right.pop()
    }
    const tmp: E[] = []
    const n = this._left.size
    for (let i = 0; i < n; i++) {
      tmp.push(this._left.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._left.append(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._right.append(tmp[i])
    }
    return this._right.pop()
  }

  popLeft(): E | undefined {
    if (this._left.size > 0) {
      return this._left.pop()
    }
    const tmp: E[] = []
    const n = this._right.size
    for (let i = 0; i < n; i++) {
      tmp.push(this._right.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._right.append(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._left.append(tmp[i])
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
    const leftSize = this._left.size
    if (index < leftSize) {
      return this._left.at(leftSize - index - 1)
    }
    return this._right.at(index - leftSize)
  }

  toString(): string {
    const sb: string[] = []
    const n = this.size
    for (let i = 0; i < n; i++) {
      sb.push(JSON.stringify(this.at(i)))
    }
    return `MaxDeque(${sb.join(', ')})`
  }

  get size(): number {
    return this._left.size + this._right.size
  }

  get max(): V {
    if (this._left.size === 0) {
      return this._right.max
    }
    if (this._right.size === 0) {
      return this._left.max
    }
    if (this._compareValue(this._left.max, this._right.max) >= 0) {
      return this._left.max
    }
    return this._right.max
  }
}

class _MaxStack<E extends { value: V }, V = E['value']> implements IStack<E> {
  private readonly _stack: E[] = []
  private readonly _maxes: V[] = []
  private readonly _compareValue: CompareFunc<V>

  constructor(compareValue: CompareFunc<V> = (a: any, b: any) => a - b) {
    this._compareValue = compareValue
  }

  append(item: E): void {
    this._stack.push(item)
    if (
      this._maxes.length === 0 ||
      this._compareValue(item.value, this._maxes[this._maxes.length - 1]) >= 0
    ) {
      this._maxes.push(item.value)
    }
  }

  pop(): E | undefined {
    const res = this._stack.pop()
    if (res && res.value === this._maxes[this._maxes.length - 1]) {
      this._maxes.pop()
    }
    return res
  }

  at(index: number): E | undefined {
    return this._stack[index]
  }

  get size(): number {
    return this._stack.length
  }

  get max(): V {
    return this._maxes[this._maxes.length - 1]
  }
}

class MinDeque<E extends { value: V }, V = E['value']> implements IDeque<E> {
  private readonly _left: _MinStack<E, V>
  private readonly _right: _MinStack<E, V>
  private readonly _compareValue: CompareFunc<V>

  constructor(compareValue: CompareFunc<V> = (a: any, b: any) => a - b) {
    this._left = new _MinStack(compareValue)
    this._right = new _MinStack(compareValue)
    this._compareValue = compareValue
  }

  append(item: E): void {
    this._right.append(item)
  }

  appendLeft(item: E): void {
    this._left.append(item)
  }

  pop(): E | undefined {
    if (this._right.size > 0) {
      return this._right.pop()
    }
    const tmp: E[] = []
    const n = this._left.size
    for (let i = 0; i < n; i++) {
      tmp.push(this._left.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._left.append(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._right.append(tmp[i])
    }
    return this._right.pop()
  }

  popLeft(): E | undefined {
    if (this._left.size > 0) {
      return this._left.pop()
    }
    const tmp: E[] = []
    const n = this._right.size
    for (let i = 0; i < n; i++) {
      tmp.push(this._right.pop()!)
    }
    const half = n >> 1
    for (let i = half - 1; i >= 0; i--) {
      this._right.append(tmp[i])
    }
    for (let i = half; i < n; i++) {
      this._left.append(tmp[i])
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
    const leftSize = this._left.size
    if (index < leftSize) {
      return this._left.at(leftSize - index - 1)
    }
    return this._right.at(index - leftSize)
  }

  toString(): string {
    const sb: string[] = []
    const n = this.size
    for (let i = 0; i < n; i++) {
      sb.push(JSON.stringify(this.at(i)))
    }
    return `MinDeque(${sb.join(', ')})`
  }

  get size(): number {
    return this._left.size + this._right.size
  }

  get min(): V {
    if (this._left.size === 0) {
      return this._right.min
    }
    if (this._right.size === 0) {
      return this._left.min
    }
    if (this._compareValue(this._left.min, this._right.min) <= 0) {
      return this._left.min
    }
    return this._right.min
  }
}

class _MinStack<E extends { value: V }, V = E['value']> implements IStack<E> {
  private readonly _stack: E[] = []
  private readonly _mins: V[] = []
  private readonly _compareValue: CompareFunc<V>

  constructor(compareValue: CompareFunc<V> = (a: any, b: any) => a - b) {
    this._compareValue = compareValue
  }

  append(item: E): void {
    this._stack.push(item)
    if (
      this._mins.length === 0 ||
      this._compareValue(item.value, this._mins[this._mins.length - 1]) <= 0
    ) {
      this._mins.push(item.value)
    }
  }

  pop(): E | undefined {
    const res = this._stack.pop()
    if (res && res.value === this._mins[this._mins.length - 1]) {
      this._mins.pop()
    }
    return res
  }

  at(index: number): E | undefined {
    return this._stack[index]
  }

  get size(): number {
    return this._stack.length
  }

  get min(): V {
    return this._mins[this._mins.length - 1]
  }
}

if (require.main === module) {
  const maxDeque = new MaxDeque()
  maxDeque.append({ value: 1 })
  maxDeque.append({ value: 3 })
  maxDeque.append({ value: 2 })
  maxDeque.append({ value: 4 })
  maxDeque.append({ value: 5 })
  console.log(maxDeque.toString())
  const res = maxDeque.popLeft()
  console.log(res)
  console.log(maxDeque.toString())
  maxDeque.appendLeft({ value: 6 })
  console.log(maxDeque.toString())
  const res2 = maxDeque.popLeft()
  console.log(res2)
  console.log(maxDeque.toString())
  console.log(maxDeque.max)
  console.log(maxDeque.at(0))
  console.log(maxDeque.size)

  class MyCircularDeque {
    private readonly _queue: IDeque<{ value: number }> = new MinDeque()
    private readonly _k: number
    constructor(k: number) {
      this._k = k
    }

    insertFront(value: number): boolean {
      if (this._queue.size === this._k) {
        return false
      }
      this._queue.appendLeft({ value })
      return true
    }

    insertLast(value: number): boolean {
      if (this._queue.size === this._k) {
        return false
      }
      this._queue.append({ value })
      return true
    }

    deleteFront(): boolean {
      if (this._queue.size === 0) {
        return false
      }
      this._queue.popLeft()
      return true
    }

    deleteLast(): boolean {
      if (this._queue.size === 0) {
        return false
      }
      this._queue.pop()
      return true
    }

    getFront(): number {
      const res = this._queue.at(0)
      return res ? res.value : -1
    }

    getRear(): number {
      const res = this._queue.at(this._queue.size - 1)
      return res ? res.value : -1
    }

    isEmpty(): boolean {
      return this._queue.size === 0
    }

    isFull(): boolean {
      return this._queue.size === this._k
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

  const myMinDeque = new MaxDeque<{ value: string }>((a, b) => a.localeCompare(b))
  myMinDeque.append({ value: 'a' })
  myMinDeque.append({ value: 'c' })
}

export { MaxDeque, MinDeque }
