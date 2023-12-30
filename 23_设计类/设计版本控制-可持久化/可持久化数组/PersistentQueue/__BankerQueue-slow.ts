// #include "data_structure/stream.cpp"
// #include "other/int_alias.cpp"

import assert from 'assert'
import { Stream } from './Stream'

class BankerQueue<T> {
  private readonly _front = new Stream<T>()
  private readonly _back = new Stream<T>()
  private readonly _fSize: number = 0
  private readonly _bSize: number = 0

  constructor()
  constructor(f: Stream<T>, b: Stream<T>, fSize: number, bSize: number)
  constructor(f?: Stream<T>, b?: Stream<T>, fSize?: number, bSize?: number) {
    const hasArgs = !!arguments.length // 实参数量
    if (hasArgs) {
      this._front = f!
      this._back = b!
      this._fSize = fSize!
      this._bSize = bSize!
    }
  }

  empty(): boolean {
    return !this._fSize
  }

  front(): T | undefined {
    if (this.empty()) {
      return undefined
    }
    return this._front.top()
  }

  push(x: T): BankerQueue<T> {
    return new BankerQueue(
      this._front,
      this._back.push(x),
      this._fSize,
      this._bSize + 1
    )._normalize()
  }

  shift(): BankerQueue<T> {
    if (this.empty()) {
      return new BankerQueue()
    }
    return new BankerQueue(
      this._front.pop()!,
      this._back,
      this._fSize - 1,
      this._bSize
    )._normalize()
  }

  private _normalize(): BankerQueue<T> {
    if (this._fSize >= this._bSize) {
      return this
    }
    return new BankerQueue(
      Stream.concat(this._front, this._back.reverse()),
      new Stream<T>(),
      this._fSize + this._bSize,
      0
    )
  }
}

export {}

if (require.main === module) {
  let queue = new BankerQueue<number>()

  assert(queue.empty())
  assert(queue.front() === undefined)
  const queue1 = queue.push(1)
  assert(!queue1.empty())
  assert(queue1.front() === 1)

  const queue2 = queue1.push(2)
  assert(!queue2.empty())
  assert(queue2.front() === 1)

  const queue3 = queue2.shift()
  assert(!queue3.empty())
  assert(queue3.front() === 2)

  const queue4 = queue3.shift()
  assert(queue4.empty())
  assert(queue4.front() === undefined)

  console.time('queue')
  for (let i = 0; i < 1e5; i++) {
    queue = queue.push(i)
    queue = queue.shift()
    queue.empty()
    queue.front()
  }
  console.timeEnd('queue') // 450ms
}
