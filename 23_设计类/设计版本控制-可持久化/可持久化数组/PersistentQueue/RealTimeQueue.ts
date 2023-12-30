/* eslint-disable @typescript-eslint/no-non-null-assertion */
// https://scrapbox.io/data-structures/Realtime_Queue

import assert from 'assert'

import { Stream } from './Stream'
import { PersistentStack } from '../../可持久化栈/PersistentStack'

/**
 * 可持久化队列.
 * @link https://scrapbox.io/data-structures/Realtime_Queue
 */
class RealTimeQueue<T> {
  private static _rotate<T>(f: Stream<T>, b: PersistentStack<T>, t: Stream<T>): Stream<T> {
    return new Stream(() => {
      if (f.empty()) {
        return [b.top()!, t]
      }
      return [f.top()!, this._rotate(f.pop()!, b.pop()!, t.push(b.top()!))]
    })
  }

  private static _makeQueue<T>(f: Stream<T>, b: PersistentStack<T>, s: Stream<T>): RealTimeQueue<T> {
    if (!s.empty()) {
      return new RealTimeQueue(f, b, s.pop()!)
    }
    const tmp = this._rotate(f, b, new Stream())
    return new RealTimeQueue(tmp, new PersistentStack(), tmp)
  }

  private readonly _front = new Stream<T>()
  private readonly _back = new PersistentStack<T>()
  private readonly _schedule = new Stream<T>()

  constructor()
  constructor(f: Stream<T>, b: PersistentStack<T>, s: Stream<T>)
  constructor(f?: Stream<T>, b?: PersistentStack<T>, s?: Stream<T>) {
    const hasArgs = !!arguments.length // 实参数量
    if (hasArgs) {
      this._front = f!
      this._back = b!
      this._schedule = s!
    }
  }

  empty(): boolean {
    return this._front.empty()
  }

  front(): T | undefined {
    return this._front.top()
  }

  push(x: T): RealTimeQueue<T> {
    return RealTimeQueue._makeQueue(this._front!, this._back!.push(x), this._schedule!)
  }

  shift(): RealTimeQueue<T> {
    if (this.empty()) {
      return new RealTimeQueue()
    }
    return RealTimeQueue._makeQueue(this._front!.pop()!, this._back!, this._schedule!)
  }
}

export { RealTimeQueue, RealTimeQueue as PersistentQueue }

if (require.main === module) {
  let queue = new RealTimeQueue<number>()

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
  for (let i = 0; i < 1e7; i++) {
    queue = queue.push(i)
    queue = queue.shift()
    queue.empty()
    queue.front()
  }
  console.timeEnd('queue') // queue: 1.447s
}
