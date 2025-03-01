// # TypeScript 条件变量 (Cond) 实现

// 条件变量是一种同步原语，用于线程等待某个条件满足后再继续执行。在 TypeScript 的异步环境中，我们可以实现类似的机制。

// 这个条件变量实现提供以下功能：

// 1. `wait()` - 释放互斥锁并等待条件变为真
// 2. `signal()` - 唤醒一个等待的线程
// 3. `broadcast()` - 唤醒所有等待的线程
// 4. 支持带谓词的等待，自动处理虚假唤醒
// 5. 支持超时机制

// 这种实现在 JavaScript/TypeScript 的单线程环境中用于协调异步操作，确保它们按照特定条件的顺序执行。

/**
 * 条件变量实现
 * 允许线程等待直到特定条件满足
 */
class Condition {
  private waiters: Array<{
    resolve: () => void
    reject: (error: Error) => void
  }> = []

  /**
   * 创建一个条件变量
   * @param mutex 关联的互斥锁
   */
  // eslint-disable-next-line no-useless-constructor
  constructor(private mutex: Mutex) {}

  /**
   * 等待条件满足
   * 调用前必须已经获取关联的互斥锁
   * @param predicate 可选的谓词函数，如果提供则仅在谓词返回true时才结束等待
   * @param timeout 可选的超时时间(毫秒)
   */
  async wait(predicate?: () => boolean, timeout?: number): Promise<void> {
    if (predicate && predicate()) {
      // 条件已经满足，无需等待
      return
    }

    // 创建等待Promise
    const waitPromise = new Promise<void>((resolve, reject) => {
      const waiter = { resolve, reject }
      this.waiters.push(waiter)

      // 处理超时
      if (timeout !== undefined) {
        setTimeout(() => {
          const index = this.waiters.indexOf(waiter)
          if (index !== -1) {
            this.waiters.splice(index, 1)
            reject(new Error('Condition wait timed out'))
          }
        }, timeout)
      }
    })

    // 释放互斥锁，让其他线程有机会修改条件状态
    this.mutex.unlock()

    try {
      // 等待条件通知
      await waitPromise
    } finally {
      // 重新获取锁
      await this.mutex.lock()

      // 如果提供了谓词，可能需要继续等待
      if (predicate && !predicate()) {
        await this.wait(predicate, timeout)
      }
    }
  }

  /**
   * 通知一个等待的线程
   * 调用前必须已经获取关联的互斥锁
   */
  signal(): void {
    if (this.waiters.length > 0) {
      const waiter = this.waiters.shift()!
      // 在下一个事件循环中解析，确保锁的正确处理
      setTimeout(() => waiter.resolve(), 0)
    }
  }

  /**
   * 通知所有等待的线程
   * 调用前必须已经获取关联的互斥锁
   */
  broadcast(): void {
    const waiters = [...this.waiters]
    this.waiters = []

    // 唤醒所有等待者
    waiters.forEach(waiter => {
      setTimeout(() => waiter.resolve(), 0)
    })
  }
}

/**
 * 互斥锁实现 (供Condition使用)
 */
class Mutex {
  private locked = false
  private waitQueue: Array<{
    resolve: (value: boolean) => void
    reject: (reason: Error) => void
  }> = []

  async lock(timeout?: number): Promise<boolean> {
    if (!this.locked) {
      this.locked = true
      return true
    }

    return new Promise<boolean>((resolve, reject) => {
      const waiter = { resolve, reject }
      this.waitQueue.push(waiter)

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

  unlock(): void {
    if (!this.locked) {
      throw new Error('Mutex is not locked')
    }

    if (this.waitQueue.length > 0) {
      const nextWaiter = this.waitQueue.shift()!
      nextWaiter.resolve(true)
    } else {
      this.locked = false
    }
  }
}

// 使用示例: 经典的生产者-消费者问题
async function example() {
  const mutex = new Mutex()
  const notEmpty = new Condition(mutex)
  const notFull = new Condition(mutex)

  const buffer: number[] = []
  const maxSize = 5

  // 生产者函数
  async function producer() {
    for (let i = 0; i < 10; i++) {
      await mutex.lock()

      // 缓冲区满则等待
      while (buffer.length >= maxSize) {
        console.log('生产者: 缓冲区已满，等待')
        await notFull.wait()
      }

      // 生产一个项目
      const item = i
      buffer.push(item)
      console.log(`生产者: 添加项目 ${item}, 缓冲区大小 ${buffer.length}`)

      // 通知消费者缓冲区不为空
      notEmpty.signal()
      mutex.unlock()

      // 模拟生产时间
      await new Promise(resolve => setTimeout(resolve, 200))
    }
  }

  // 消费者函数
  async function consumer() {
    while (true) {
      await mutex.lock()

      // 缓冲区空则等待
      while (buffer.length === 0) {
        console.log('消费者: 缓冲区为空，等待')
        await notEmpty.wait()
      }

      // 消费一个项目
      const item = buffer.shift()!
      console.log(`消费者: 移除项目 ${item}, 缓冲区大小 ${buffer.length}`)

      // 通知生产者缓冲区不满
      notFull.signal()
      mutex.unlock()

      // 模拟消费时间
      await new Promise(resolve => setTimeout(resolve, 300))

      // 在此例中，消费10个项目后退出
      if (item === 9) break
    }
  }

  // 并行启动生产者和消费者
  await Promise.all([producer(), consumer()])
}

example()

export { Condition }
