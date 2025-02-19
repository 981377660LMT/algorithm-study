/* eslint-disable no-inner-declarations */
//
// 简单版本 Mutex 实现
//
// 1. Go Mutex 的核心特性:
//    - Lock() 方法获取锁
//    - Unlock() 方法释放锁
//    - 阻塞式等待
//    - 不可重入
//
// 2. 实现计划:
//    - 使用 Promise 实现异步阻塞
//    - 用闭包维护锁状态
//    - 提供 Lock/Unlock 接口
//    - 添加防重入检查
//
// 主要特点:
// - 使用 async/await 实现阻塞等待
// - 维护等待队列实现公平性
// - 提供锁状态检查
// - 支持 try/finally 释放锁

export class Mutex {
  private _locked = false
  private readonly _queue: Array<() => void> = []

  async lock(): Promise<void> {
    if (this._locked) {
      await new Promise<void>(resolve => {
        this._queue.push(resolve)
      })
    }
    this._locked = true
  }

  unlock(): void {
    if (!this._locked) {
      throw new Error('Unlock of unlocked mutex')
    }

    this._locked = false
    const next = this._queue.shift()
    if (next) {
      next()
    }
    console.log(111)
  }

  locked(): boolean {
    return this._locked
  }
}

if (require.main === module) {
  async function example() {
    const mutex = new Mutex()

    // 使用锁保护共享资源
    await mutex.lock()
    try {
      // 临界区代码
      console.log('Critical section')
    } finally {
      mutex.unlock()
    }
  }
  example()
}
