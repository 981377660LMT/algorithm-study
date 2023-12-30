/* eslint-disable max-len */
/* eslint-disable no-shadow */

// notes:
// !1.注意evalute用 `xx._evaluate =` 重新赋值会比写在构造函数里面快很多;
// !2.由于内部占用了构造函数，因此对外不暴露构造函数，只能通过init方法来创建实例.

import assert from 'assert'

/**
 * 完全可持久化队列.
 */
class PersistentQueue<V> {
  static init<V>(): PersistentQueue<V> {
    return new PersistentQueue()
  }

  private readonly _frontSize: number
  private readonly _rearSize: number
  private readonly _front: _PersistentStack<V> | undefined
  private readonly _rear: _PersistentStack<V> | undefined

  private constructor(frontSize?: number, rearSize?: number, front?: _PersistentStack<V>, rear?: _PersistentStack<V>) {
    this._frontSize = frontSize || 0
    this._rearSize = rearSize || 0
    this._front = front
    this._rear = rear
  }

  empty(): boolean {
    return !this._frontSize
  }

  top(): V | undefined {
    if (!this._front) return undefined
    return this._front.top()
  }

  push(x: V): PersistentQueue<V> {
    return new PersistentQueue(this._frontSize, this._rearSize + 1, this._front, _PersistentStack.push(this._rear, x))._normalize()
  }

  pop(): PersistentQueue<V> {
    return new PersistentQueue(this._frontSize - 1, this._rearSize, this._front && this._front.pop(), this._rear)._normalize()
  }

  length(): number {
    return this._frontSize + this._rearSize
  }

  private _normalize(): PersistentQueue<V> {
    if (this._frontSize >= this._rearSize) {
      return this
    }
    return new PersistentQueue(this._frontSize + this._rearSize, 0, _PersistentStack.concat(this._front, _PersistentStack.reverse(this._rear!)))
  }
}

class _PersistentStack<V> {
  static init<V>(): _PersistentStack<V> {
    return new _PersistentStack()
  }

  static push<V>(x: _PersistentStack<V> | undefined, v: V): _PersistentStack<V> {
    return new _PersistentStack(v, x)
  }

  static concat<V>(x: _PersistentStack<V> | undefined, y: _PersistentStack<V>): _PersistentStack<V> {
    if (!x) {
      return y._evaluate()
    }
    const next = new _PersistentStack<V>()
    next._evaluate = () => _PersistentStack.concat(x.pop()!, y)
    return new _PersistentStack(x._value, next)
  }

  static reverse<V>(head: _PersistentStack<V>): _PersistentStack<V> {
    const res = new _PersistentStack<V>()
    res._evaluate = () => {
      let tmp: _PersistentStack<V> | undefined
      for (let x = head; x; x = x.pop()!) {
        tmp = _PersistentStack.push(tmp, x.top()!)
      }
      return tmp!
    }
    return res
  }

  private _next: _PersistentStack<V> | undefined
  private readonly _value: V | undefined
  private _evaluate: () => _PersistentStack<V>

  private constructor(value?: V, next?: _PersistentStack<V>, evaluate?: () => _PersistentStack<V>) {
    this._next = next
    this._value = value
    this._evaluate = evaluate || (() => this)
  }

  empty(): boolean {
    return !this._next
  }

  top(): V | undefined {
    return this._value
  }

  pop(): _PersistentStack<V> | undefined {
    if (this._next) {
      this._next = this._next._evaluate()
    }
    return this._next
  }
}

export { PersistentQueue }

if (require.main === module) {
  let queue = PersistentQueue.init<number>()

  assert(queue.empty())
  assert(queue.top() === undefined)
  const queue1 = queue.push(1)
  assert(!queue1.empty())
  assert(queue1.top() === 1)

  const queue2 = queue1.push(2)
  assert(!queue2.empty())
  assert(queue2.top() === 1)

  const queue3 = queue2.pop()
  assert(!queue3.empty())
  assert(queue3.top() === 2)

  const queue4 = queue3.pop()
  assert(queue4.empty())
  assert(queue4.top() === undefined)

  console.time('queue')
  for (let i = 0; i < 1e7; i++) {
    queue = queue.push(i)
    queue = queue.pop()
    queue.empty()
    queue.top()
  }
  console.timeEnd('queue') // 325.638ms
}
