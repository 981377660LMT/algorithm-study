// 这个实现提供了互斥锁的基本功能：
// - 确保在同一时刻只有一个调用者可以持有锁
// - 支持等待队列，先到先得地分配锁
// - 提供尝试获取锁的非阻塞方法
// - 支持超时机制
// - 提供便捷的 withLock 方法来自动管理锁的释放
//
// 在 JavaScript/TypeScript 单线程环境中，这种互斥锁主要用于协调异步操作，确保对共享资源的访问是序列化的。

/**
 * 互斥锁实现.
 * 确保在任何时刻只有一个执行上下文可以持有锁.
 */
export class Mutex2 {
  private locked = false
  private waitQueue: Array<{
    resolve: (value: boolean) => void
    reject: (reason: Error) => void
  }> = []

  /**
   * 获取锁，如果锁已被获取，则等待直到锁可用
   * @param timeout 可选的超时时间(毫秒)
   * @returns Promise，解析为true表示获取锁成功
   */
  async lock(timeout?: number): Promise<boolean> {
    // 如果锁空闲，立即获取
    if (!this.locked) {
      this.locked = true
      return true
    }

    // 否则加入等待队列
    return new Promise<boolean>((resolve, reject) => {
      const waiter = { resolve, reject }
      this.waitQueue.push(waiter)

      // 处理超时
      if (timeout !== undefined) {
        setTimeout(() => {
          const index = this.waitQueue.indexOf(waiter)
          if (index !== -1) {
            this.waitQueue.splice(index, 1)
            reject(new Error('Lock acquisition timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 尝试获取锁，如果无法立即获取则返回false
   * @returns 是否成功获取锁
   */
  tryLock(): boolean {
    if (!this.locked) {
      this.locked = true
      return true
    }
    return false
  }

  /**
   * 释放锁
   * @throws Error 如果锁未被获取
   */
  unlock(): void {
    if (!this.locked) {
      throw new Error('Mutex is not locked')
    }

    if (this.waitQueue.length > 0) {
      // 唤醒队列中的下一个等待者
      const nextWaiter = this.waitQueue.shift()!
      nextWaiter.resolve(true)
    } else {
      this.locked = false
    }
  }

  /**
   * 使用锁执行一个函数，并在函数执行完毕后自动释放锁
   * @param fn 需要在锁内执行的异步函数
   * @param timeout 可选的超时时间
   * @returns 函数的返回值Promise
   */
  async withLock<T>(fn: () => Promise<T> | T, timeout?: number): Promise<T> {
    const acquired = await this.lock(timeout)
    if (!acquired) {
      throw new Error('Failed to acquire lock')
    }

    try {
      return await Promise.resolve(fn())
    } finally {
      this.unlock()
    }
  }

  /**
   * 获取锁的状态信息
   */
  getStatus(): { locked: boolean; queueLength: number } {
    return {
      locked: this.locked,
      queueLength: this.waitQueue.length
    }
  }
}

// 使用示例
async function example() {
  const mutex = new Mutex2()

  // 获取锁
  await mutex.lock()
  console.log('锁已获取', mutex.getStatus())

  // 尝试获取锁（将失败，因为锁已被持有）
  const acquired = mutex.tryLock()
  console.log('尝试获取锁:', acquired)

  // 释放锁
  mutex.unlock()
  console.log('锁已释放', mutex.getStatus())

  // 使用 withLock 自动管理锁
  await mutex.withLock(async () => {
    console.log('在锁保护中执行操作')
    // eslint-disable-next-line no-promise-executor-return
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('操作完成，锁将被自动释放')
  })

  console.log('最终状态', mutex.getStatus())
}
example()
