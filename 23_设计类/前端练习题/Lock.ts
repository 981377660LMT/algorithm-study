type VoidPromise = (value: void | PromiseLike<void>) => void

/**
 * 保证同一时刻只有一个异步操作可以执行.
 * 在持有锁的情况下，后续的调用排队等待。
 */
class Lock {
  private _isLocked = false
  private readonly _resolvingQueue: VoidPromise[] = []

  /**
   * @alias acquire
   */
  lock(): Promise<void> {
    if (!this._isLocked) {
      this._isLocked = true
      return Promise.resolve()
    }

    return new Promise<void>(resolve => {
      this._resolvingQueue.push(resolve)
    })
  }

  /**
   * @alias release
   */
  unlock(): void {
    if (this._resolvingQueue.length) {
      const resolve = this._resolvingQueue.shift()!
      resolve()
    } else {
      this._isLocked = false
    }
  }
}

function withLock<T>(task: () => Promise<T>, lock: Lock): Promise<T> {
  return lock.lock().then(() =>
    task().finally(() => {
      lock.unlock()
    })
  )
}

export {}

if (require.main === module) {
  // 使用示例

  {
    const lock = new Lock()

    async function criticalSection() {
      // 尝试获取锁
      await lock.lock()
      try {
        // 执行需要同步的代码
        console.log('Critical section is running')
        // 模拟异步操作
        await new Promise(resolve => setTimeout(resolve, 1000))
      } finally {
        // 释放锁，允许下一个等待的操作执行
        lock.unlock()
      }
    }

    // ！同时启动多个异步操作，但它们将会按顺序执行
    criticalSection()
    criticalSection()
    criticalSection()
  }

  {
    const lock = new Lock()
    const task = () => {
      return new Promise<void>(resolve => {
        setTimeout(() => {
          console.log('task done')
          resolve()
        }, 1000)
      })
    }

    withLock(task, lock)
    withLock(task, lock)
    withLock(task, lock)
  }
}
