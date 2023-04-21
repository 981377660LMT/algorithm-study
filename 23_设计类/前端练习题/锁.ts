/**
 * 异步代码中，有时候需要保证同一时间只有一个任务在执行，这时候就需要锁。
 */
class Lock {
  private readonly _resolveQueue: ((value?: unknown) => void)[] = []
  private _locked = false

  async acquire() {
    if (this._locked) {
      await new Promise<unknown>(resolve => {
        this._resolveQueue.push(resolve)
      })
    }
    this._locked = true
  }

  release() {
    if (!this._locked) {
      throw new Error('The lock must be acquired before release')
    }
    this._locked = false
    if (this._resolveQueue.length) {
      const resolve = this._resolveQueue.shift()!
      resolve()
    }
  }

  /**
   * Excute async task with lock.
   */
  async execute(asyncFunc: () => Promise<unknown>) {
    await this.acquire()
    try {
      return await asyncFunc()
    } finally {
      this.release()
    }
  }
}

export {}
