/* eslint-disable no-promise-executor-return */

/**
 * CountDownLatch - 一种同步辅助类
 *
 * CountDownLatch 允许一个或多个线程等待，直到在其他线程中执行的一组操作完成。
 * 这个类设计为一次性使用 - 计数无法被重置。如果需要重置的版本，
 * 可以考虑使用 CyclicBarrier。
 */
class CountDownLatch {
  private _count: number // 剩余计数
  private _waiters: Array<{
    resolve: () => void
    reject?: (reason?: any) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = [] // 等待计数器归零的线程

  /**
   * 创建一个新的 CountDownLatch，初始计数为指定值
   *
   * @param count 初始计数，必须大于或等于0
   * @throws 如果计数小于0，则抛出错误
   */
  constructor(count: number) {
    if (count < 0) {
      throw new Error('Count must be non-negative')
    }
    this._count = count
  }

  /**
   * 获取当前计数
   */
  get count(): number {
    return this._count
  }

  /**
   * 递减锁存器的计数，如果计数达到零，则释放所有等待的线程
   *
   * 如果当前计数大于零，则递减计数。
   * 如果递减后计数为零，则释放所有等待的线程。
   * 如果当前计数已经为零，则此方法不起作用。
   */
  countDown(): void {
    if (this._count <= 0) {
      return
    }

    this._count--

    // 如果计数已达到零，唤醒所有等待的线程
    if (this._count === 0) {
      const waiters = [...this._waiters]
      this._waiters = []

      for (const waiter of waiters) {
        if (waiter.timeout) {
          clearTimeout(waiter.timeout)
        }
        waiter.resolve()
      }
    }
  }

  /**
   * 一次将计数减少指定的值
   *
   * @param delta 要减少的数量，必须大于0且不超过当前计数
   * @throws 如果delta小于1或大于当前计数，则抛出错误
   */
  countDownMultiple(delta: number): void {
    if (delta <= 0) {
      throw new Error('Delta must be positive')
    }

    if (delta > this._count) {
      throw new Error('Cannot reduce count below zero')
    }

    const newCount = this._count - delta
    this._count = newCount

    // 如果计数已达到零，唤醒所有等待的线程
    if (newCount === 0) {
      const waiters = [...this._waiters]
      this._waiters = []

      for (const waiter of waiters) {
        if (waiter.timeout) {
          clearTimeout(waiter.timeout)
        }
        waiter.resolve()
      }
    }
  }

  /**
   * 使当前线程等待，直到锁存器的计数递减到零，或发生超时
   *
   * 如果当前计数为零，则此方法立即返回。
   *
   * @param timeout 可选的超时时间(毫秒)，如果指定，则在超时后会抛出错误
   * @returns 返回一个Promise，当计数达到零时解析，或在超时时拒绝
   */
  async await(timeout?: number): Promise<void> {
    // 如果计数已经为零，立即返回
    if (this._count === 0) {
      return Promise.resolve()
    }

    // 创建一个Promise来等待计数器归零
    return new Promise<void>((resolve, reject) => {
      const waiter: (typeof this._waiters)[number] = { resolve, reject }
      this._waiters.push(waiter)

      // 如果设置了超时
      if (timeout !== undefined && timeout >= 0) {
        waiter.timeout = setTimeout(() => {
          // 从等待队列中移除
          const index = this._waiters.indexOf(waiter)
          if (index !== -1) {
            this._waiters.splice(index, 1)
            reject(new Error('CountDownLatch await timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 尝试等待计数器归零，最多等待指定的时间
   *
   * @param timeout 等待的最长时间(毫秒)
   * @returns 如果计数达到零则返回true，如果等待超时则返回false
   */
  async awaitWithBoolean(timeout: number): Promise<boolean> {
    // 如果计数已经为零，立即返回
    if (this._count === 0) {
      return true
    }

    try {
      await this.await(timeout)
      return true
    } catch (e) {
      return false
    }
  }

  /**
   * 一次性将计数减少到零，并释放所有等待的线程
   */
  countDownAll(): void {
    if (this._count > 0) {
      this._count = 0

      const waiters = [...this._waiters]
      this._waiters = []

      for (const waiter of waiters) {
        if (waiter.timeout) {
          clearTimeout(waiter.timeout)
        }
        waiter.resolve()
      }
    }
  }
}

export { CountDownLatch }

// 使用示例
async function countDownLatchExample() {
  console.log('CountDownLatch 基本用法示例:')

  // 创建一个计数为3的CountDownLatch
  const latch = new CountDownLatch(3)
  console.log(`初始计数: ${latch.count}`)

  // 工作线程函数
  async function worker(id: number, delay: number) {
    console.log(`工作线程 ${id} 开始执行`)
    await new Promise(resolve => setTimeout(resolve, delay))
    console.log(`工作线程 ${id} 完成，调用countDown()`)
    latch.countDown()
    console.log(`当前计数: ${latch.count}`)
  }

  // 等待线程函数
  async function waiter() {
    console.log('等待线程: 等待所有工作线程完成')
    await latch.await()
    console.log('等待线程: 所有工作线程已完成，继续执行')
  }

  // 启动工作线程和等待线程
  const promises = [worker(1, 1000), worker(2, 2000), worker(3, 3000), waiter()]

  await Promise.all(promises)

  // 超时示例
  console.log('\nCountDownLatch 超时示例:')
  const timeoutLatch = new CountDownLatch(2)

  async function timeoutExample() {
    try {
      console.log('尝试等待带有1秒超时')
      await timeoutLatch.await(1000)
      console.log('成功等待!')
    } catch (error) {
      console.log(`等待超时: ${(error as Error).message}`)
    }

    // 使用布尔返回值的版本
    const result = await timeoutLatch.awaitWithBoolean(1000)
    console.log(`使用awaitWithBoolean的结果: ${result ? '成功' : '超时'}`)
  }

  // 只减少一个计数，使超时发生
  timeoutLatch.countDown()
  console.log(`减少一个计数，当前计数: ${timeoutLatch.count}`)

  await timeoutExample()

  // countDownMultiple示例
  console.log('\ncountDownMultiple示例:')
  const multiLatch = new CountDownLatch(5)
  console.log('创建计数为5的锁存器')

  // 减少多个计数
  multiLatch.countDownMultiple(3)
  console.log(`减少3个计数，当前计数: ${multiLatch.count}`)

  // 使用countDownAll
  console.log('\ncountDownAll示例:')
  console.log('调用countDownAll()')
  multiLatch.countDownAll()
  console.log(`当前计数: ${multiLatch.count}`)

  // 验证锁存器已释放
  const isReleased = await multiLatch.awaitWithBoolean(0)
  console.log(`锁存器已释放: ${isReleased}`)
}

// 运行示例
countDownLatchExample().catch(console.error)
