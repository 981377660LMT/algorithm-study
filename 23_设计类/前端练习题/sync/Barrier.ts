/**
 * Barrier 类 - 实现同步屏障
 *
 * 屏障是一种同步原语，它允许多个任务在某一点等待，
 * 直到所有任务都到达该点后才能同时继续执行。
 *
 * 每次使用后，屏障可以自动重置，以便重复使用。
 */
class Barrier {
  private _parties: number // 需要等待的任务数量
  private _waitCount = 0 // 当前等待的任务数量
  private _generation = 0 // 当前代数，用于处理重用
  private _broken = false // 屏障是否已损坏
  private _waiters: Array<{
    resolve: (index: number) => void
    reject: (reason?: any) => void
    generation: number
  }> = [] // 等待的任务队列

  private _onReset?: () => void // 屏障重置时的回调函数

  /**
   * 创建一个新的屏障
   *
   * @param parties 需要到达屏障的任务数量
   * @param action 可选的回调函数，在所有任务到达时执行一次
   * @throws 如果parties小于1，则抛出错误
   */
  constructor(parties: number, action?: () => void) {
    if (parties < 1) {
      throw new Error('Barrier parties must be at least 1')
    }
    this._parties = parties
    this._onReset = action
  }

  /**
   * 获取需要到达屏障的任务数量
   */
  get parties(): number {
    return this._parties
  }

  /**
   * 获取当前在屏障处等待的任务数量
   */
  get numberOfWaitingParties(): number {
    return this._waitCount
  }

  /**
   * 获取屏障是否已损坏
   */
  get isBroken(): boolean {
    return this._broken
  }

  /**
   * 等待其他任务到达屏障
   *
   * 当所有任务都到达屏障时，屏障会被打开，所有等待的任务会被释放，
   * 屏障会被重置以备下一次使用。
   *
   * @param timeout 可选的超时时间(毫秒)，超时后将抛出错误
   * @returns 返回当前线程在此屏障上的到达索引，范围是0到parties-1
   * @throws 如果屏障已损坏，则抛出错误
   * @throws 如果等待超时，则抛出错误
   */
  async wait(timeout?: number): Promise<number> {
    if (this._broken) {
      throw new Error('Barrier is broken')
    }

    const myGeneration = this._generation
    const myIndex = this._waitCount++

    // 如果这是最后一个到达的任务
    if (myIndex === this._parties - 1) {
      // 重置屏障
      return this.reset()
    }

    // 等待其他任务到达
    return new Promise<number>((resolve, reject) => {
      const waiter = { resolve, reject, generation: myGeneration }
      this._waiters.push(waiter)

      // 设置超时
      if (timeout !== undefined && timeout >= 0) {
        setTimeout(() => {
          const index = this._waiters.indexOf(waiter)
          if (index >= 0 && waiter.generation === myGeneration) {
            this._waiters.splice(index, 1)
            this.breakBarrier()
            reject(new Error('Barrier wait timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 重置屏障为初始状态
   *
   * 此方法通常由最后一个到达的任务自动调用，但也可以手动调用来重置屏障。
   *
   * @returns 当前线程的到达索引，如果是由最后一个到达的任务调用，则返回parties-1
   */
  private reset(): number {
    // 获取当前等待任务的数量（这将是最后一个到达任务的索引）
    const index = this._waitCount - 1

    // 增加代数，使旧的等待者失效
    this._generation++

    // 重置计数器和破坏标志
    this._waitCount = 0
    this._broken = false

    // 执行重置回调（如果有）
    try {
      if (this._onReset) {
        this._onReset()
      }
    } catch (error) {
      // 如果回调抛出错误，将屏障标记为损坏
      this._broken = true
    }

    // 唤醒所有等待的任务
    const currentWaiters = this._waiters.filter(w => w.generation === this._generation - 1)
    this._waiters = this._waiters.filter(w => w.generation !== this._generation - 1)

    for (const waiter of currentWaiters) {
      if (this._broken) {
        waiter.reject(new Error('Barrier is broken'))
      } else {
        waiter.resolve(index)
      }
    }

    // 如果回调导致屏障被破坏，抛出错误
    if (this._broken) {
      throw new Error('Barrier is broken due to exception in barrier action')
    }

    return index
  }

  /**
   * 将屏障置于损坏状态
   *
   * 当屏障处于损坏状态时，所有等待的任务将收到错误，
   * 后续对wait的调用也将失败，直到屏障被重置。
   */
  breakBarrier(): void {
    if (this._broken) {
      return // 已经是损坏状态
    }

    this._broken = true
    this._generation++ // 使当前等待者失效

    // 拒绝所有等待的任务
    const currentWaiters = [...this._waiters]
    this._waiters = []

    for (const waiter of currentWaiters) {
      waiter.reject(new Error('Barrier is broken'))
    }
  }

  /**
   * 重置屏障状态，允许它再次使用
   *
   * 即使屏障曾经被破坏，也可以通过此方法重新启用它。
   */
  resetManually(): void {
    this._broken = false
    this._waitCount = 0
    this._generation++
    this._waiters = this._waiters.filter(w => w.generation === this._generation)
  }
}

export { Barrier }

// 使用示例
async function barrierExample() {
  // 创建一个需要3个任务到达的屏障，当所有任务到达时执行回调
  const barrier = new Barrier(3, () => {
    console.log('所有任务都已到达屏障！继续执行...')
  })

  // 模拟多个异步任务
  async function task(id: number, delay: number): Promise<void> {
    console.log(`任务 ${id} 开始执行`)

    // 模拟工作
    // eslint-disable-next-line no-promise-executor-return
    await new Promise(resolve => setTimeout(resolve, delay))

    console.log(`任务 ${id} 到达屏障，等待其他任务...`)

    // 等待屏障
    const index = await barrier.wait()

    console.log(`任务 ${id} 通过屏障，到达索引: ${index}`)
  }

  // 并行启动3个任务，每个任务有不同的延迟
  await Promise.all([task(1, 1000), task(2, 2000), task(3, 3000)])

  console.log('所有任务完成')

  // 演示屏障重用
  console.log('\n重用屏障...\n')

  await Promise.all([task(4, 500), task(5, 1000), task(6, 1500)])

  console.log('第二轮任务完成')
}

// 运行示例
barrierExample().catch(console.error)
