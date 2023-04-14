// https://scrapbox.io/data-structures/Physicist's_Queue

import assert from 'assert'
import { PersistentStack } from '../../可持久化栈/PersistentStack'
import { Suspension } from './Suspension'

class PhysicistQueue<T> {
  private readonly _working = new PersistentStack<T>()
  private readonly _front = new Suspension(() => new PersistentStack<T>())
  private readonly _back = new PersistentStack<T>()
  private readonly _fSize: number = 0
  private readonly _bSize: number = 0

  constructor()
  constructor(
    working: PersistentStack<T>,
    front: Suspension<PersistentStack<T>>,
    fSize: number,
    back: PersistentStack<T>,
    bSize: number
  )
  constructor(
    working?: PersistentStack<T>,
    front?: Suspension<PersistentStack<T>>,
    fSize?: number,
    back?: PersistentStack<T>,
    bSize?: number
  ) {
    const hasArgs = !!arguments.length // 实参数量
    if (hasArgs) {
      this._working = working!
      this._front = front!
      this._fSize = fSize!
      this._back = back!
      this._bSize = bSize!
    }
  }

  empty(): boolean {
    return !this._fSize
  }

  front(): T | undefined {
    return this._working.top()
  }

  push(x: T): PhysicistQueue<T> {
    return new PhysicistQueue(
      this._working,
      this._front,
      this._fSize,
      this._back.push(x),
      this._bSize + 1
    )._check()
  }

  shift(): PhysicistQueue<T> {
    if (this.empty()) {
      throw new Error('PhysicistQueue.shift: empty queue')
    }
    return new PhysicistQueue(
      this._working.pop(),
      new Suspension(() => this._front.resolve()!.pop()),
      this._fSize - 1,
      this._back,
      this._bSize
    )._check()
  }

  private _checkR(): PhysicistQueue<T> {
    if (this._fSize >= this._bSize) {
      return this
    }
    const tmp = this._front.resolve()!

    const f = () => {
      let r = this._back.reverse()
      let l = tmp.reverse()
      while (!l.empty()) {
        r = r.push(l.top()!)
        l = l.pop()
      }
      return r
    }

    return new PhysicistQueue(
      tmp,
      new Suspension(f),
      this._fSize + this._bSize,
      new PersistentStack(),
      0
    )
  }

  private _checkW(): PhysicistQueue<T> {
    if (!this._working.empty()) {
      return this
    }
    return new PhysicistQueue(
      this._front.resolve()!,
      this._front,
      this._fSize,
      this._back,
      this._bSize
    )
  }

  private _check(): PhysicistQueue<T> {
    return this._checkR()._checkW()
  }
}

export {}

if (require.main === module) {
  let queue = new PhysicistQueue<number>()

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
  console.timeEnd('queue') // 900ms
}
