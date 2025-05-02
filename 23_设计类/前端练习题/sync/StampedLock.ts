/* eslint-disable brace-style */
/* eslint-disable no-promise-executor-return */
/* eslint-disable no-throw-literal */
/**
 * 表示锁的状态
 */
enum LockState {
  UNLOCKED = 0, // 未锁定
  WRITE_LOCKED, // 写锁定
  READ_LOCKED // 读锁定
}

/**
 * StampedLock - 一种支持乐观读的读写锁
 *
 * 提供三种模式:
 * 1. 写锁 - 排他锁，类似于互斥锁
 * 2. 读锁 - 共享锁，多个读操作可以同时进行
 * 3. 乐观读 - 无锁模式，可能会读取到正在被写入的数据
 *
 * 不同于ReentrantReadWriteLock，StampedLock不可重入但性能更高。
 */
class StampedLock {
  // 内部状态
  private _state: LockState = LockState.UNLOCKED
  private _writeStamp = 0 // 当前写锁的戳记
  private _readCount = 0 // 当前读锁的数量
  private _stamp = 0 // 全局戳记计数器

  // 等待队列
  private _writeWaiters: Array<{
    resolve: (stamp: number) => void
    reject?: (reason?: any) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = []

  private _readWaiters: Array<{
    resolve: (stamp: number) => void
    reject?: (reason?: any) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = []

  /**
   * 获取写锁（排他锁）
   *
   * @param timeout 可选的超时时间(毫秒)
   * @returns 返回一个表示锁持有的戳记（非零值）
   * @throws 如果在超时前无法获取锁，则抛出错误
   */
  async writeLock(timeout?: number): Promise<number> {
    // 如果当前没有锁，直接获取写锁
    if (this._state === LockState.UNLOCKED) {
      return this._acquireWriteLock()
    }

    // 需要等待
    return new Promise<number>((resolve, reject) => {
      const waiter: (typeof this._writeWaiters)[number] = { resolve, reject }
      this._writeWaiters.push(waiter)

      // 设置超时
      if (timeout !== undefined && timeout >= 0) {
        waiter.timeout = setTimeout(() => {
          const index = this._writeWaiters.indexOf(waiter)
          if (index !== -1) {
            this._writeWaiters.splice(index, 1)
            reject(new Error('Write lock acquisition timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 获取读锁（共享锁）
   *
   * @param timeout 可选的超时时间(毫秒)
   * @returns 返回一个表示锁持有的戳记（非零值）
   * @throws 如果在超时前无法获取锁，则抛出错误
   */
  async readLock(timeout?: number): Promise<number> {
    // 如果当前是未锁定或已有其他读锁，且没有写等待者，则可以获取读锁
    if (
      (this._state === LockState.UNLOCKED || this._state === LockState.READ_LOCKED) &&
      this._writeWaiters.length === 0
    ) {
      return this._acquireReadLock()
    }

    // 需要等待
    return new Promise<number>((resolve, reject) => {
      const waiter: (typeof this._readWaiters)[number] = { resolve, reject }
      this._readWaiters.push(waiter)

      // 设置超时
      if (timeout !== undefined && timeout >= 0) {
        waiter.timeout = setTimeout(() => {
          const index = this._readWaiters.indexOf(waiter)
          if (index !== -1) {
            this._readWaiters.splice(index, 1)
            reject(new Error('Read lock acquisition timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 尝试获取写锁，非阻塞
   *
   * @returns 如果成功获取锁，返回戳记；如果失败，返回0
   */
  tryWriteLock(): number {
    if (this._state === LockState.UNLOCKED) {
      return this._acquireWriteLock()
    }
    return 0
  }

  /**
   * 尝试获取读锁，非阻塞
   *
   * @returns 如果成功获取锁，返回戳记；如果失败，返回0
   */
  tryReadLock(): number {
    if (
      (this._state === LockState.UNLOCKED || this._state === LockState.READ_LOCKED) &&
      this._writeWaiters.length === 0
    ) {
      return this._acquireReadLock()
    }
    return 0
  }

  /**
   * 获取乐观读戳记
   * 不会阻塞，即使当前有写锁
   *
   * @returns 当前的戳记值
   */
  tryOptimisticRead(): number {
    // 如果有写锁，返回0，表示无法乐观读取
    if (this._state === LockState.WRITE_LOCKED) {
      return 0
    }
    // 否则返回当前戳记
    return this._stamp
  }

  /**
   * 验证乐观读戳记是否仍然有效
   *
   * @param stamp 之前通过tryOptimisticRead获取的戳记
   * @returns 如果戳记仍然有效（没有写操作发生）返回true；否则返回false
   */
  validate(stamp: number): boolean {
    // 如果戳记为0，始终无效
    if (stamp === 0) return false

    // 如果戳记不等于当前戳记，表示发生了写操作
    return stamp === this._stamp
  }

  /**
   * 释放锁
   *
   * @param stamp 获取锁时返回的戳记
   * @returns 如果戳记有效且成功释放锁，返回true；否则返回false
   */
  unlock(stamp: number): boolean {
    // 检查写锁
    if (this._state === LockState.WRITE_LOCKED && stamp === this._writeStamp) {
      this._releaseWriteLock()
      return true
    }

    // 检查读锁(这里简化处理，不验证特定读锁的戳记)
    if (this._state === LockState.READ_LOCKED && stamp !== 0) {
      this._releaseReadLock()
      return true
    }

    return false
  }

  /**
   * 将乐观读转换为读锁
   *
   * @param stamp 乐观读的戳记
   * @returns 如果转换成功，返回新的读锁戳记；否则返回0
   */
  tryConvertToReadLock(stamp: number): number {
    // 验证乐观读戳记
    if (!this.validate(stamp)) {
      return 0
    }

    // 如果当前没有写锁，可以获取读锁
    if (this._state !== LockState.WRITE_LOCKED) {
      return this._acquireReadLock()
    }

    return 0
  }

  /**
   * 将乐观读转换为写锁
   *
   * @param stamp 乐观读的戳记
   * @returns 如果转换成功，返回新的写锁戳记；否则返回0
   */
  tryConvertToWriteLock(stamp: number): number {
    // 验证乐观读戳记
    if (!this.validate(stamp)) {
      return 0
    }

    // 如果当前是未锁定状态，可以获取写锁
    if (this._state === LockState.UNLOCKED) {
      return this._acquireWriteLock()
    }

    // 如果当前是读锁状态，但只有一个读者(自己)，可以升级为写锁
    if (this._state === LockState.READ_LOCKED && this._readCount === 1) {
      this._state = LockState.WRITE_LOCKED
      this._readCount = 0
      const newStamp = ++this._stamp
      this._writeStamp = newStamp
      return newStamp
    }

    return 0
  }

  /**
   * 将读锁转换为写锁
   *
   * @param stamp 读锁的戳记
   * @returns 如果转换成功，返回新的写锁戳记；否则返回0
   */
  tryConvertReadLockToWriteLock(stamp: number): number {
    // 必须当前持有读锁
    if (this._state !== LockState.READ_LOCKED || stamp === 0) {
      return 0
    }

    // 如果只有一个读者(自己)，可以升级为写锁
    if (this._readCount === 1) {
      this._state = LockState.WRITE_LOCKED
      this._readCount = 0
      const newStamp = ++this._stamp
      this._writeStamp = newStamp
      return newStamp
    }

    return 0
  }

  /**
   * 将写锁降级为读锁
   *
   * @param stamp 写锁的戳记
   * @returns 如果转换成功，返回新的读锁戳记；否则返回0
   */
  tryConvertWriteLockToReadLock(stamp: number): number {
    // 必须当前持有写锁
    if (this._state !== LockState.WRITE_LOCKED || stamp !== this._writeStamp) {
      return 0
    }

    // 降级为读锁
    this._state = LockState.READ_LOCKED
    this._readCount = 1
    this._writeStamp = 0
    return stamp // 可以保持相同的戳记，因为没有其他写操作
  }

  /**
   * 获取当前锁状态的字符串表示
   */
  getStateString(): string {
    switch (this._state) {
      case LockState.UNLOCKED:
        return 'UNLOCKED'
      case LockState.WRITE_LOCKED:
        return 'WRITE_LOCKED'
      case LockState.READ_LOCKED:
        return `READ_LOCKED (count: ${this._readCount})`
      default:
        return 'UNKNOWN'
    }
  }

  /**
   * 使用读锁执行函数
   *
   * @param fn 在获取读锁后执行的函数
   * @returns 函数执行的结果
   */
  async withReadLock<T>(fn: () => Promise<T> | T): Promise<T> {
    const stamp = await this.readLock()
    try {
      return await Promise.resolve(fn())
    } finally {
      this.unlock(stamp)
    }
  }

  /**
   * 使用写锁执行函数
   *
   * @param fn 在获取写锁后执行的函数
   * @returns 函数执行的结果
   */
  async withWriteLock<T>(fn: () => Promise<T> | T): Promise<T> {
    const stamp = await this.writeLock()
    try {
      return await Promise.resolve(fn())
    } finally {
      this.unlock(stamp)
    }
  }

  /**
   * 使用乐观读执行函数，如果检测到冲突则回退到读锁
   *
   * @param fn 在乐观读模式下执行的函数，接收一个validate函数用于检查数据一致性
   * @returns 函数执行的结果
   */
  async withOptimisticRead<T>(fn: (validate: () => boolean) => Promise<T> | T): Promise<T> {
    // 首先尝试乐观读
    const stamp = this.tryOptimisticRead()

    // 如果无法获取乐观读戳记，直接使用读锁
    if (stamp === 0) {
      return this.withReadLock(fn as any)
    }

    try {
      // 创建验证函数
      const validate = () => this.validate(stamp)

      // 执行用户函数
      const result = await Promise.resolve(fn(validate))

      // 如果最终验证失败，使用读锁重试
      if (!validate()) {
        return this.withReadLock(fn as any)
      }

      return result
    } catch (error) {
      // 如果发生错误且是由于验证失败，使用读锁重试
      if (error === 'validation_failed') {
        return this.withReadLock(fn as any)
      }
      throw error
    }
  }

  // 私有辅助方法

  /**
   * 内部方法：获取写锁
   */
  private _acquireWriteLock(): number {
    this._state = LockState.WRITE_LOCKED
    const newStamp = ++this._stamp
    this._writeStamp = newStamp
    return newStamp
  }

  /**
   * 内部方法：获取读锁
   */
  private _acquireReadLock(): number {
    if (this._state !== LockState.READ_LOCKED) {
      this._state = LockState.READ_LOCKED
    }
    this._readCount++
    return this._stamp // 读锁使用当前戳记
  }

  /**
   * 内部方法：释放写锁
   */
  private _releaseWriteLock(): void {
    // 更新戳记，表示写操作完成
    this._stamp++
    this._writeStamp = 0
    this._state = LockState.UNLOCKED

    // 尝试唤醒等待的读者和写者
    this._wakeUpWaiters()
  }

  /**
   * 内部方法：释放读锁
   */
  private _releaseReadLock(): void {
    this._readCount--

    // 如果没有更多读者，更改状态为未锁定
    if (this._readCount === 0) {
      this._state = LockState.UNLOCKED

      // 尝试唤醒等待的写者
      this._wakeUpWaiters()
    }
  }

  /**
   * 内部方法：根据当前状态唤醒等待的线程
   */
  private _wakeUpWaiters(): void {
    if (this._state !== LockState.UNLOCKED) {
      return
    }

    // 优先唤醒写者(防止写饥饿)
    if (this._writeWaiters.length > 0) {
      const waiter = this._writeWaiters.shift()!
      if (waiter.timeout) {
        clearTimeout(waiter.timeout)
      }

      // 获取写锁
      const stamp = this._acquireWriteLock()
      waiter.resolve(stamp)
    }
    // 如果没有写者，唤醒所有读者
    else if (this._readWaiters.length > 0) {
      // 先获取一个读锁的戳记
      const stamp = this._acquireReadLock()

      // 唤醒所有等待的读者
      const waiters = [...this._readWaiters]
      this._readWaiters = []

      // 第一个读者已经通过_acquireReadLock处理了
      const firstWaiter = waiters.shift()!
      if (firstWaiter.timeout) {
        clearTimeout(firstWaiter.timeout)
      }
      firstWaiter.resolve(stamp)

      // 其余读者
      for (const waiter of waiters) {
        if (waiter.timeout) {
          clearTimeout(waiter.timeout)
        }

        // 不需要增加戳记，只需增加读计数
        this._readCount++
        waiter.resolve(stamp)
      }
    }
  }
}

export { StampedLock }

// 使用示例
async function stampedLockExample() {
  const lock = new StampedLock()

  // 共享的数据
  let sharedData = {
    count: 0,
    text: 'Initial value'
  }

  // 1. 写操作示例
  console.log('执行写操作...')
  await lock.withWriteLock(async () => {
    console.log('获取写锁')
    sharedData.count++
    sharedData.text = 'Updated value'
    await new Promise(resolve => setTimeout(resolve, 1000))
    console.log('写操作完成')
  })

  // 2. 读操作示例
  console.log('\n执行多个并发读操作...')
  await Promise.all([
    lock.withReadLock(async () => {
      console.log('读取器1: ', sharedData)
      await new Promise(resolve => setTimeout(resolve, 500))
    }),
    lock.withReadLock(async () => {
      console.log('读取器2: ', sharedData)
      await new Promise(resolve => setTimeout(resolve, 300))
    }),
    lock.withReadLock(async () => {
      console.log('读取器3: ', sharedData)
      await new Promise(resolve => setTimeout(resolve, 200))
    })
  ])

  // 3. 乐观读示例
  console.log('\n执行乐观读...')
  await lock.withOptimisticRead(validate => {
    console.log('乐观读开始，当前数据:', sharedData)

    // 检查数据是否仍然有效
    if (validate()) {
      console.log('乐观读有效，可以安全使用数据')
    } else {
      console.log('乐观读无效，数据可能已更改')
      throw 'validation_failed'
    }
    return sharedData.count
  })

  // 4. 锁升级示例
  console.log('\n演示锁升级...')
  const readStamp = await lock.readLock()
  console.log('获取读锁，当前值:', sharedData.count)

  // 尝试升级到写锁
  const writeStamp = lock.tryConvertReadLockToWriteLock(readStamp)
  if (writeStamp !== 0) {
    console.log('成功将读锁升级为写锁')
    sharedData.count++
    console.log('修改后的值:', sharedData.count)
    lock.unlock(writeStamp)
  } else {
    console.log('无法升级到写锁，释放读锁')
    lock.unlock(readStamp)
  }

  // 5. 锁降级示例
  console.log('\n演示锁降级...')
  const ws = await lock.writeLock()
  console.log('获取写锁')
  sharedData.count *= 2

  // 将写锁降级为读锁
  const rs = lock.tryConvertWriteLockToReadLock(ws)
  if (rs !== 0) {
    console.log('成功将写锁降级为读锁')
    console.log('当前值(只读):', sharedData.count)
    lock.unlock(rs)
  } else {
    console.log('无法降级锁')
    lock.unlock(ws)
  }
}

// 运行示例
stampedLockExample().catch(console.error)
