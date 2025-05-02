/* eslint-disable no-promise-executor-return */

/**
 * Semaphore 类 - 实现计数信号量
 *
 * 信号量是一种同步原语，用于控制对共享资源的访问数量。
 * 它维护一个内部计数器，表示可用资源数量：
 * - acquire() 方法减少计数器，如果计数器为零则阻塞
 * - release() 方法增加计数器，可能唤醒等待的线程
 */
class Semaphore {
  private _count: number // 当前可用资源数量
  private _waitingQueue: Array<{
    resolve: () => void
    reject?: (reason?: any) => void
  }> // 等待队列

  /**
   * 创建一个新的信号量
   * @param count 初始可用资源数量(默认为1)
   * @throws 如果计数为负数则抛出错误
   */
  constructor(count = 1) {
    if (count < 0) {
      throw new Error('Semaphore count cannot be negative')
    }
    this._count = count
    this._waitingQueue = []
  }

  /**
   * 获取当前信号量的计数值
   */
  get count(): number {
    return this._count
  }

  /**
   * 尝试获取资源，如果资源可用，则减少计数并返回。
   * 如果资源不可用，则等待直到有资源被释放。
   *
   * @param timeout 可选的超时时间(毫秒)，如果指定了超时且在超时前无法获取资源，将抛出错误
   * @returns Promise，在成功获取资源时解析
   * @throws 如果在超时前无法获取资源，则拒绝Promise
   */
  async acquire(timeout?: number): Promise<void> {
    // 如果有可用资源，直接获取
    if (this._count > 0) {
      this._count--
      return Promise.resolve()
    }

    // 否则将请求加入等待队列
    return new Promise<void>((resolve, reject) => {
      const waiter = { resolve }
      this._waitingQueue.push(waiter)

      // 如果指定了超时，设置超时处理
      if (timeout !== undefined && timeout >= 0) {
        const timeoutId = setTimeout(() => {
          // 从等待队列中移除此等待者
          const index = this._waitingQueue.indexOf(waiter)
          if (index !== -1) {
            this._waitingQueue.splice(index, 1)
            reject(new Error('Semaphore acquisition timed out'))
          }
        }, timeout)

        // 覆盖原始的resolve，以便在成功获取时清除超时
        const originalResolve = waiter.resolve
        waiter.resolve = () => {
          clearTimeout(timeoutId)
          originalResolve()
        }
      }
    })
  }

  /**
   * 尝试非阻塞方式获取资源
   *
   * @returns 如果成功获取资源返回true，否则返回false
   */
  tryAcquire(): boolean {
    if (this._count > 0) {
      this._count--
      return true
    }
    return false
  }

  /**
   * 释放一个资源，增加信号量计数
   * 如果有线程在等待，则唤醒队列中的第一个等待线程
   *
   * @param count 要释放的资源数量(默认为1)
   * @throws 如果参数为负数则抛出错误
   */
  release(count = 1): void {
    if (count < 0) {
      throw new Error('Cannot release negative count')
    }

    // 一次只释放一个等待的线程，确保公平性
    for (let i = 0; i < count; i++) {
      if (this._waitingQueue.length > 0) {
        // 有等待的线程，直接唤醒一个
        const waiter = this._waitingQueue.shift()!
        waiter.resolve()
      } else {
        // 没有等待的线程，增加计数
        this._count++
      }
    }
  }

  /**
   * 使用Semaphore执行受限的异步函数
   *
   * @param fn 要执行的函数
   * @returns 函数执行的结果
   */
  async execute<T>(fn: () => Promise<T> | T): Promise<T> {
    await this.acquire()
    try {
      return await Promise.resolve(fn())
    } finally {
      this.release()
    }
  }
}

export { Semaphore }

async function semaphoreExample() {
  const semaphore = new Semaphore(2)

  // 模拟一组需要限制并发的异步任务
  const tasks = Array.from({ length: 5 }, (_, i) => i + 1)

  console.log('开始执行限制并发的任务')

  // 并发执行所有任务，但最多同时执行2个
  await Promise.all(
    tasks.map(async task =>
      semaphore.execute(async () => {
        console.log(`任务 ${task} 开始执行，当前并发: ${2 - semaphore.count}`)
        // 模拟异步操作
        await new Promise(resolve => setTimeout(resolve, 1000))
        console.log(`任务 ${task} 执行完成`)
        return task
      })
    )
  )

  console.log('所有任务执行完毕')
}

// 运行示例
semaphoreExample().catch(console.error)
