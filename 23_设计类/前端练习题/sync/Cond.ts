// 条件变量 `Cond` 建议版本实现
//
// 1. `Cond` 核心特性:
//    - 基于互斥锁(Mutex)实现条件变量
//    - 提供 Wait/Signal/Broadcast 方法
//    - Wait() 会自动释放锁并阻塞等待
//    - Signal() 唤醒一个等待者
//    - Broadcast() 唤醒所有等待者
//
// 2. TypeScript 实现思路:
//    - 使用之前实现的 Mutex 作为基础
//    - 利用 Promise 实现等待队列
//    - 维护等待者列表
//    - 实现三个核心方法
//
// ts 里可以不用lock.

import { Mutex } from './Mutex'

export class Cond {
  readonly lock: Mutex
  private _waiters: Array<() => void> = []

  constructor(mutex: Mutex) {
    this.lock = mutex
  }

  /**
   * 释放锁并挂起,当被唤醒时重新获取锁.
   */
  async wait(): Promise<void> {
    if (!this.lock.locked()) {
      throw new Error('Wait requires lock to be held')
    }

    // runtime_notifyListAdd
    const waiter = new Promise<void>(resolve => {
      this._waiters.push(resolve)
    })
    this.lock.unlock()
    // runtime_notifyListWait
    await waiter
    await this.lock.lock()
  }

  signal(): void {
    if (!this.lock.locked()) {
      throw new Error('Signal requires lock to be held')
    }

    // runtime_notifyListNotifyOne
    const f = this._waiters.shift()
    if (f) {
      f()
    }
  }

  broadcast(): void {
    if (!this.lock.locked()) {
      throw new Error('Broadcast requires lock to be held')
    }

    // runtime_notifyListNotifyAll
    const waiters = this._waiters
    this._waiters = []
    for (const f of waiters) {
      f()
    }
  }
}

if (require.main === module) {
  // eslint-disable-next-line no-inner-declarations
  async function example() {
    const mutex = new Mutex()
    const cond = new Cond(mutex)
    let sharedState = false

    // 等待条件的异步任务
    async function waiter() {
      await mutex.lock()
      try {
        console.log('[Waiter] 获取锁，检查条件')
        while (!sharedState) {
          console.log('[Waiter] 条件未满足，进入等待...')
          // eslint-disable-next-line no-await-in-loop
          await cond.wait() // 释放锁并等待
        }
        console.log('[Waiter] 条件满足，继续执行')
      } finally {
        mutex.unlock()
      }
    }

    // 修改条件并通知的异步任务
    async function notifier() {
      await mutex.lock()
      try {
        console.log('[Notifier] 获取锁，修改条件')
        sharedState = true
        cond.signal() // 通知一个等待者
        // cond.broadcast(); // 通知所有等待者
        console.log('[Notifier] 已发送通知')
      } finally {
        mutex.unlock()
      }
    }

    // 启动两个任务
    waiter()
    setTimeout(notifier, 2000) // 2秒后触发通知
  }

  example()
}
