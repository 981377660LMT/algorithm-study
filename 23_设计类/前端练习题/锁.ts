/**
 * 异步代码中，有时候需要保证同一时间只有一个任务在执行，这时候就需要锁。
 */
class Lock {
  /** 存储等待解锁的任务的 resolve 函数. */
  private readonly _resolveQueue: ((value: void | PromiseLike<void>) => void)[] = []
  private _locked = false

  /**
   * 请求获取锁。如果锁已被占用（{@link _locked} 为 true），则阻塞请求直到 resolve 被调用。
   */
  async acquire(): Promise<void> {
    if (this._locked) {
      // 阻塞请求直到resolve被调用
      await new Promise(resolve => {
        this._resolveQueue.push(resolve)
      })
    }
    this._locked = true
  }

  /**
   * 释放锁。如果锁未被占用（{@link _locked }为 false），则抛出错误。
   * 否则，将锁的状态设置为解锁（{@link _locked }为 false），并从队列中取出一个等待解锁的任务的 resolve 函数，调用该函数以解锁。
   */
  release(): void {
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
   * 用于执行带有锁的异步任务。
   */
  async execute<T>(asyncFunc: () => Promise<T>): Promise<T> {
    await this.acquire()
    try {
      return await asyncFunc()
    } finally {
      this.release()
    }
  }

  /** 返回锁的状态。 */
  get locked(): boolean {
    return this._locked
  }
}

export {}
