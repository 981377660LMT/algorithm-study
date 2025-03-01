// 这个实现提供了可重入锁的基本功能：
// - 同一个上下文可以多次获取锁
// - 每次获取增加计数，每次释放减少计数
// - 只有当计数归零时才真正释放锁
// - 支持等待队列和可选的超时机制
//
// 在实际项目中可能需要根据具体需求进行调整，特别是关于如何标识锁的当前所有者。

/**
 * 可重入锁实现
 * 允许同一个持有者多次获取锁而不会引起死锁
 */
class ReentrantLock {
  private owner: symbol | null = null
  private count = 0
  private waitQueue: Array<{
    resolve: () => void
    reject: (error: Error) => void
  }> = []

  /**
   * 尝试获取锁，如果锁已被当前持有者获取，则递增计数
   * 如果锁被其他持有者获取，则等待
   * @param timeout 可选的超时时间(毫秒)
   * @returns Promise，解析为true表示获取锁成功
   */
  async acquire(timeout?: number): Promise<boolean> {
    const currentOwner = this.getCurrentOwner()

    // 如果当前调用者已持有锁，增加计数
    if (this.owner === currentOwner) {
      this.count++
      return true
    }

    // 如果锁被其他人持有，等待
    if (this.owner !== null) {
      return new Promise<boolean>((resolve, reject) => {
        const waitPromise = { resolve: () => resolve(this.doAcquire(currentOwner)), reject }

        this.waitQueue.push(waitPromise)

        // 处理超时
        if (timeout !== undefined) {
          setTimeout(() => {
            // 从等待队列中移除
            const index = this.waitQueue.indexOf(waitPromise)
            if (index !== -1) {
              this.waitQueue.splice(index, 1)
              reject(new Error('Lock acquisition timed out'))
            }
          }, timeout)
        }
      })
    }

    // 锁未被持有，直接获取
    return this.doAcquire(currentOwner)
  }

  /**
   * 释放锁，减少计数，如果计数为0，则完全释放锁
   * @returns 是否成功释放
   */
  release(): boolean {
    const currentOwner = this.getCurrentOwner()

    if (this.owner !== currentOwner) {
      throw new Error('Cannot release a lock that is not owned by the current context')
    }

    this.count--

    if (this.count === 0) {
      this.owner = null

      // 如果等待队列中有等待者，唤醒第一个
      if (this.waitQueue.length > 0) {
        const next = this.waitQueue.shift()!
        next.resolve()
      }
    }

    return true
  }

  /**
   * 获取当前锁的状态信息
   */
  getStatus(): { locked: boolean; count: number; queueLength: number } {
    return {
      locked: this.owner !== null,
      count: this.count,
      queueLength: this.waitQueue.length
    }
  }

  /**
   * 内部方法：实际获取锁
   */
  private doAcquire(owner: symbol): boolean {
    this.owner = owner
    this.count = 1
    return true
  }

  /**
   * 获取当前执行上下文的标识符
   * 在真实场景中，这可能需要更复杂的机制来唯一标识调用者
   */
  private getCurrentOwner(): symbol {
    // 如果当前执行没有关联的所有者，创建一个新的
    if (!this._currentOwner) {
      this._currentOwner = Symbol('lock-owner')
    }
    return this._currentOwner
  }

  // 存储当前上下文的所有者标识
  private _currentOwner: symbol | null = null
}

// 使用示例
async function example() {
  const lock = new ReentrantLock()

  // 首次获取锁
  await lock.acquire()
  console.log('锁已获取', lock.getStatus())

  // 再次获取相同的锁（重入）
  await lock.acquire()
  console.log('锁再次获取', lock.getStatus())

  // 释放一次
  lock.release()
  console.log('释放一次', lock.getStatus())

  // 再次释放，完全释放锁
  lock.release()
  console.log('完全释放', lock.getStatus())
}
example()

export { ReentrantLock }
