/* eslint-disable no-promise-executor-return */
/* eslint-disable brace-style */
/**
 * ReadWriteLock - 读写锁实现
 *
 * 允许多个读操作同时进行，但写操作必须独占。
 * 提供公平模式选项，可以避免读操作长期占用导致写操作饥饿。
 */
class ReadWriteLock {
  private _readCount = 0 // 当前读取器数量
  private _writeRequested = false // 是否有写请求等待中
  private _writing = false // 是否有写操作进行中
  private _fairMode: boolean // 是否使用公平模式

  // 等待队列
  private _readWaiters: Array<{
    resolve: () => void
    reject?: (reason?: any) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = []

  private _writeWaiters: Array<{
    resolve: () => void
    reject?: (reason?: any) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = []

  /**
   * 创建一个新的读写锁
   *
   * @param fairMode 是否使用公平模式。在公平模式下，如果有写操作在等待，
   *                新的读操作将等待写操作完成，防止写操作饥饿。默认为true。
   */
  constructor(fairMode = true) {
    this._fairMode = fairMode
  }

  /**
   * 获取读锁
   *
   * 在非公平模式下，只要没有写操作进行中，就可以立即获取。
   * 在公平模式下，如果有写操作等待，则需要等待写操作完成。
   *
   * @param timeout 可选的超时时间(毫秒)
   * @returns 返回一个Promise，当获取读锁成功时解析
   * @throws 如果在超时前无法获取读锁，则抛出错误
   */
  async readLock(timeout?: number): Promise<void> {
    // 如果没有写操作进行中且（没有等待的写操作或非公平模式），直接获取读锁
    if (!this._writing && (!this._fairMode || !this._writeRequested)) {
      this._readCount++
      return Promise.resolve()
    }

    // 需要等待
    return new Promise<void>((resolve, reject) => {
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
   * 释放读锁
   *
   * 减少读取器计数，如果没有更多读取器且有写操作等待，则唤醒一个写操作。
   */
  readUnlock(): void {
    if (this._readCount <= 0) {
      throw new Error('No read lock to release')
    }

    this._readCount--

    // 如果没有更多读取器且有写操作在等待，唤醒一个写操作
    if (this._readCount === 0 && this._writeWaiters.length > 0) {
      this._wakeUpWriter()
    }
  }

  /**
   * 获取写锁
   *
   * 只有当没有读操作和其他写操作时，才能获取写锁。
   *
   * @param timeout 可选的超时时间(毫秒)
   * @returns 返回一个Promise，当获取写锁成功时解析
   * @throws 如果在超时前无法获取写锁，则抛出错误
   */
  async writeLock(timeout?: number): Promise<void> {
    // 标记有写请求
    this._writeRequested = true

    // 如果没有读取器和写操作，直接获取写锁
    if (this._readCount === 0 && !this._writing) {
      this._writing = true
      this._writeRequested = false
      return Promise.resolve()
    }

    // 需要等待
    return new Promise<void>((resolve, reject) => {
      const waiter: (typeof this._writeWaiters)[number] = { resolve, reject }
      this._writeWaiters.push(waiter)

      // 设置超时
      if (timeout !== undefined && timeout >= 0) {
        waiter.timeout = setTimeout(() => {
          const index = this._writeWaiters.indexOf(waiter)
          if (index !== -1) {
            this._writeWaiters.splice(index, 1)

            // 如果没有更多写等待者，取消写请求标记
            if (this._writeWaiters.length === 0) {
              this._writeRequested = false
            }

            reject(new Error('Write lock acquisition timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 释放写锁
   *
   * 标记写操作完成，并根据等待队列唤醒读操作或下一个写操作。
   */
  writeUnlock(): void {
    if (!this._writing) {
      throw new Error('No write lock to release')
    }

    this._writing = false

    // 如果有写操作等待且在公平模式下，优先唤醒写操作
    if (this._writeWaiters.length > 0 && this._fairMode) {
      this._wakeUpWriter()
    }
    // 否则唤醒所有等待的读操作
    else if (this._readWaiters.length > 0) {
      this._wakeUpReaders()
    }
    // 如果没有读操作等待且有写操作等待，唤醒一个写操作
    else if (this._writeWaiters.length > 0) {
      this._wakeUpWriter()
    }
    // 如果没有任何等待操作，清除写请求标记
    else {
      this._writeRequested = false
    }
  }

  /**
   * 尝试获取读锁，非阻塞
   *
   * @returns 如果成功获取读锁返回true，否则返回false
   */
  tryReadLock(): boolean {
    if (!this._writing && (!this._fairMode || !this._writeRequested)) {
      this._readCount++
      return true
    }
    return false
  }

  /**
   * 尝试获取写锁，非阻塞
   *
   * @returns 如果成功获取写锁返回true，否则返回false
   */
  tryWriteLock(): boolean {
    if (this._readCount === 0 && !this._writing) {
      this._writing = true
      this._writeRequested = false
      return true
    }
    return false
  }

  /**
   * 获取当前状态信息
   *
   * @returns 包含当前锁状态的对象
   */
  getStatus(): {
    readCount: number
    writing: boolean
    writeRequested: boolean
    readWaiters: number
    writeWaiters: number
    fairMode: boolean
  } {
    return {
      readCount: this._readCount,
      writing: this._writing,
      writeRequested: this._writeRequested,
      readWaiters: this._readWaiters.length,
      writeWaiters: this._writeWaiters.length,
      fairMode: this._fairMode
    }
  }

  /**
   * 使用读锁执行函数
   *
   * @param fn 在获取读锁后执行的函数
   * @returns 函数执行的结果
   */
  async withReadLock<T>(fn: () => Promise<T> | T): Promise<T> {
    await this.readLock()
    try {
      return await Promise.resolve(fn())
    } finally {
      this.readUnlock()
    }
  }

  /**
   * 使用写锁执行函数
   *
   * @param fn 在获取写锁后执行的函数
   * @returns 函数执行的结果
   */
  async withWriteLock<T>(fn: () => Promise<T> | T): Promise<T> {
    await this.writeLock()
    try {
      return await Promise.resolve(fn())
    } finally {
      this.writeUnlock()
    }
  }

  /**
   * 切换公平模式
   *
   * @param fairMode 是否启用公平模式
   */
  setFairMode(fairMode: boolean): void {
    this._fairMode = fairMode
  }

  // 私有辅助方法

  /**
   * 唤醒一个等待的写操作
   */
  private _wakeUpWriter(): void {
    if (this._writeWaiters.length > 0) {
      const waiter = this._writeWaiters.shift()!
      if (waiter.timeout) {
        clearTimeout(waiter.timeout)
      }
      this._writing = true

      // 如果没有更多写等待者，清除写请求标记
      if (this._writeWaiters.length === 0) {
        this._writeRequested = false
      }

      waiter.resolve()
    }
  }

  /**
   * 唤醒所有等待的读操作
   */
  private _wakeUpReaders(): void {
    if (this._readWaiters.length > 0) {
      // 唤醒所有等待的读操作
      const waiters = [...this._readWaiters]
      this._readWaiters = []

      // 增加读取器计数
      this._readCount += waiters.length

      // 清除所有超时并解析Promise
      for (const waiter of waiters) {
        if (waiter.timeout) {
          clearTimeout(waiter.timeout)
        }
        waiter.resolve()
      }
    }
  }
}

export { ReadWriteLock }

// 使用示例
async function readWriteLockExample() {
  // 创建读写锁（启用公平模式）
  const rwLock = new ReadWriteLock(true)

  // 共享数据
  const sharedData = {
    value: 0,
    lastUpdatedBy: ''
  }

  // 1. 多个并发读操作
  console.log('执行多个并发读操作...')
  await Promise.all([
    rwLock.withReadLock(async () => {
      console.log('读取器 1: 值 =', sharedData.value)
      await new Promise(resolve => setTimeout(resolve, 100))
      console.log('读取器 1: 完成')
    }),
    rwLock.withReadLock(async () => {
      console.log('读取器 2: 值 =', sharedData.value)
      await new Promise(resolve => setTimeout(resolve, 200))
      console.log('读取器 2: 完成')
    }),
    rwLock.withReadLock(async () => {
      console.log('读取器 3: 值 =', sharedData.value)
      await new Promise(resolve => setTimeout(resolve, 50))
      console.log('读取器 3: 完成')
    })
  ])

  // 2. 写操作必须独占
  console.log('\n执行写操作(独占)...')
  await rwLock.withWriteLock(async () => {
    console.log('写入器: 开始更新')
    sharedData.value++
    sharedData.lastUpdatedBy = 'Writer 1'
    await new Promise(resolve => setTimeout(resolve, 300))
    console.log('写入器: 完成更新，新值 =', sharedData.value)
  })

  // 3. 演示写操作后的读操作
  console.log('\n执行写操作后的读操作...')
  await rwLock.withReadLock(() => {
    console.log('更新后读取: 值 =', sharedData.value, '更新者:', sharedData.lastUpdatedBy)
  })

  // 4. 演示公平模式（写操作不会被新的读操作饿死）
  console.log('\n演示公平模式...')

  // 启动一个长时间运行的读操作
  const longReadPromise = rwLock.withReadLock(async () => {
    console.log('长时间读取: 开始')
    await new Promise(resolve => setTimeout(resolve, 500))
    console.log('长时间读取: 完成')
  })

  // 等待100ms后请求写锁
  await new Promise(resolve => setTimeout(resolve, 100))
  const writePromise = rwLock.writeLock()
  console.log('写入器: 请求写锁（需要等待所有读操作完成）')

  // 尝试开始另一个读操作（在公平模式下，应该等待写操作）
  await new Promise(resolve => setTimeout(resolve, 100))
  const newReadPromise = rwLock.readLock()
  console.log('新读取器: 请求读锁（在公平模式下，会等待写操作完成）')

  // 等待长读操作完成
  await longReadPromise

  // 等待写操作获得锁并完成
  await writePromise
  console.log('写入器: 获得写锁')
  sharedData.value++
  sharedData.lastUpdatedBy = 'Writer 2'
  console.log('写入器: 完成更新，新值 =', sharedData.value)
  rwLock.writeUnlock()
  console.log('写入器: 释放写锁')

  // 等待新读操作完成
  await newReadPromise
  console.log('新读取器: 获得读锁，值 =', sharedData.value)
  rwLock.readUnlock()
  console.log('新读取器: 释放读锁')

  // 5. 演示非阻塞尝试
  console.log('\n演示非阻塞尝试...')

  // 获取读锁
  if (rwLock.tryReadLock()) {
    console.log('成功获取读锁（非阻塞）')
    rwLock.readUnlock()
  } else {
    console.log('无法获取读锁（非阻塞）')
  }

  // 获取写锁
  await rwLock.writeLock()
  console.log('获取写锁')

  // 尝试获取读锁（应该失败，因为写锁正在进行中）
  if (rwLock.tryReadLock()) {
    console.log('成功获取读锁（非阻塞）')
    rwLock.readUnlock()
  } else {
    console.log('无法获取读锁（非阻塞），因为写锁在进行中')
  }

  // 尝试获取写锁（应该失败，因为已持有写锁）
  if (rwLock.tryWriteLock()) {
    console.log('成功获取另一个写锁（非阻塞）')
    rwLock.writeUnlock()
  } else {
    console.log('无法获取另一个写锁（非阻塞），因为已持有写锁')
  }

  // 释放写锁
  rwLock.writeUnlock()
  console.log('释放写锁')

  // 打印最终状态
  console.log('\n最终锁状态:', rwLock.getStatus())
}

// 运行示例
readWriteLockExample().catch(console.error)
