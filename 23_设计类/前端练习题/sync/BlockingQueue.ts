/* eslint-disable consistent-return */
/* eslint-disable no-promise-executor-return */
/* eslint-disable no-await-in-loop */

/**
 * BlockingQueue<T> - 线程安全的阻塞队列
 *
 * 当队列为空时，take()操作会阻塞直到有元素可用
 * 当队列已满时，put()操作会阻塞直到有空间可用
 */
class BlockingQueue<T> {
  private _queue: T[] = [] // 内部队列存储
  private _capacity: number // 队列最大容量
  private _putters: Array<{
    resolve: () => void
    element: T
    timeout?: ReturnType<typeof setTimeout>
  }> = [] // 等待添加元素的线程
  private _takers: Array<{
    resolve: (value: T) => void
    timeout?: ReturnType<typeof setTimeout>
  }> = [] // 等待获取元素的线程

  /**
   * 创建一个新的阻塞队列
   *
   * @param capacity 队列的最大容量(默认为Infinity，表示无界队列)
   * @throws 如果容量小于等于0，则抛出错误
   */
  constructor(capacity = Infinity) {
    if (capacity <= 0) {
      throw new Error('Queue capacity must be greater than 0')
    }
    this._capacity = capacity
  }

  /**
   * 获取队列中当前的元素数量
   */
  get size(): number {
    return this._queue.length
  }

  /**
   * 获取队列的最大容量
   */
  get capacity(): number {
    return this._capacity
  }

  /**
   * 检查队列是否为空
   */
  get isEmpty(): boolean {
    return this._queue.length === 0
  }

  /**
   * 检查队列是否已满
   */
  get isFull(): boolean {
    return this._queue.length >= this._capacity
  }

  /**
   * 尝试将元素添加到队列尾部
   * 如果队列已满，则阻塞等待
   *
   * @param element 要添加的元素
   * @param timeout 可选的超时时间(毫秒)，超时后抛出错误
   * @throws 如果在超时前无法添加元素，则抛出错误
   */
  async put(element: T, timeout?: number): Promise<void> {
    // 如果有等待的消费者，直接将元素传递给第一个消费者
    if (this._takers.length > 0) {
      const taker = this._takers.shift()!
      // 清除taker的超时定时器(如果有)
      if (taker.timeout) {
        clearTimeout(taker.timeout)
      }
      // 传递元素给消费者
      taker.resolve(element)
      return
    }

    // 如果队列未满，直接添加元素
    if (this._queue.length < this._capacity) {
      this._queue.push(element)
      return
    }

    // 队列已满，需要等待
    return new Promise<void>((resolve, reject) => {
      const putter: (typeof this._putters)[number] = { resolve, element }
      this._putters.push(putter)

      // 如果设置了超时
      if (timeout !== undefined && timeout >= 0) {
        putter.timeout = setTimeout(() => {
          // 从等待队列中移除
          const index = this._putters.indexOf(putter)
          if (index !== -1) {
            this._putters.splice(index, 1)
            reject(new Error('Put operation timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 尝试立即将元素添加到队列
   * 如果队列已满，则返回false而不阻塞
   *
   * @param element 要添加的元素
   * @returns 如果成功添加则返回true，否则返回false
   */
  offer(element: T): boolean {
    // 如果有等待的消费者，直接将元素传递给它
    if (this._takers.length > 0) {
      const taker = this._takers.shift()!
      if (taker.timeout) {
        clearTimeout(taker.timeout)
      }
      taker.resolve(element)
      return true
    }

    // 如果队列未满，则添加元素
    if (this._queue.length < this._capacity) {
      this._queue.push(element)
      return true
    }

    // 队列已满，返回false
    return false
  }

  /**
   * 从队列头部获取并移除元素
   * 如果队列为空，则阻塞等待
   *
   * @param timeout 可选的超时时间(毫秒)，超时后抛出错误
   * @returns 队列头部的元素
   * @throws 如果在超时前无法获取元素，则抛出错误
   */
  async take(timeout?: number): Promise<T> {
    // 如果队列非空，直接返回头部元素
    if (this._queue.length > 0) {
      const element = this._queue.shift()!

      // 如果有等待的生产者，让第一个生产者添加元素
      if (this._putters.length > 0) {
        const putter = this._putters.shift()!
        // 清除putter的超时定时器(如果有)
        if (putter.timeout) {
          clearTimeout(putter.timeout)
        }
        // 添加元素到队列
        this._queue.push(putter.element)
        // 通知生产者已添加成功
        putter.resolve()
      }

      return element
    }

    // 队列为空，需要等待
    return new Promise<T>((resolve, reject) => {
      const taker: (typeof this._takers)[number] = { resolve }
      this._takers.push(taker)

      // 如果设置了超时
      if (timeout !== undefined && timeout >= 0) {
        taker.timeout = setTimeout(() => {
          // 从等待队列中移除
          const index = this._takers.indexOf(taker)
          if (index !== -1) {
            this._takers.splice(index, 1)
            reject(new Error('Take operation timed out'))
          }
        }, timeout)
      }
    })
  }

  /**
   * 尝试立即从队列头部获取并移除元素
   * 如果队列为空，则返回undefined而不阻塞
   *
   * @returns 队列头部的元素，如果队列为空则返回undefined
   */
  poll(): T | undefined {
    // 如果队列为空，返回undefined
    if (this._queue.length === 0) {
      return undefined
    }

    // 获取头部元素
    const element = this._queue.shift()!

    // 如果有等待的生产者，让第一个生产者添加元素
    if (this._putters.length > 0) {
      const putter = this._putters.shift()!
      if (putter.timeout) {
        clearTimeout(putter.timeout)
      }
      this._queue.push(putter.element)
      putter.resolve()
    }

    return element
  }

  /**
   * 查看队列头部的元素，但不移除
   *
   * @returns 队列头部的元素，如果队列为空则返回undefined
   */
  peek(): T | undefined {
    return this._queue.length > 0 ? this._queue[0] : undefined
  }

  /**
   * 清空队列，并唤醒所有等待中的生产者
   * 注意：这会使等待中的消费者继续等待
   */
  clear(): void {
    // 清空队列
    this._queue = []

    // 唤醒所有等待的生产者，告诉他们可以添加元素了
    // (因为队列现在是空的)
    const putters = [...this._putters]
    this._putters = []

    for (const putter of putters) {
      if (putter.timeout) {
        clearTimeout(putter.timeout)
      }

      // 只有在队列未满的情况下才添加元素
      if (this._queue.length < this._capacity) {
        this._queue.push(putter.element)
        putter.resolve()
      } else {
        // 重新添加到等待队列
        this._putters.push(putter)
      }
    }
  }

  /**
   * 将指定元素插入此队列的尾部，如果队列已满则等待
   * 此方法等同于put()，为了与Java API兼容
   */
  async add(element: T, timeout?: number): Promise<void> {
    return this.put(element, timeout)
  }

  /**
   * 检索并移除此队列的头，如果队列为空则等待
   * 此方法等同于take()，为了与Java API兼容
   */
  async remove(timeout?: number): Promise<T> {
    return this.take(timeout)
  }

  /**
   * 以数组形式返回队列中的所有元素
   * 注意：这不会修改队列
   */
  toArray(): T[] {
    return [...this._queue]
  }

  /**
   * 创建一个带有初始元素的阻塞队列
   *
   * @param elements 初始元素数组
   * @param capacity 队列容量(默认为数组长度)
   */
  static of<T>(elements: T[], capacity?: number): BlockingQueue<T> {
    const actualCapacity = capacity !== undefined ? capacity : Math.max(elements.length, 1)
    const queue = new BlockingQueue<T>(actualCapacity)

    // 添加初始元素(不会超过容量)
    for (let i = 0; i < Math.min(elements.length, actualCapacity); i++) {
      queue._queue.push(elements[i])
    }

    return queue
  }
}

export { BlockingQueue }

// 使用示例
async function blockingQueueExample() {
  console.log('创建有界阻塞队列(容量为2)')
  const queue = new BlockingQueue<number>(2)

  // 生产者函数
  async function producer() {
    for (let i = 1; i <= 5; i++) {
      console.log(`生产者: 尝试添加元素 ${i}`)
      await queue.put(i)
      console.log(`生产者: 成功添加元素 ${i}, 队列大小: ${queue.size}`)
      await new Promise(resolve => setTimeout(resolve, 500)) // 生产间隔
    }
  }

  // 消费者函数
  async function consumer() {
    // 先等待一段时间，让队列有机会填满
    await new Promise(resolve => setTimeout(resolve, 1000))

    for (let i = 1; i <= 5; i++) {
      console.log('消费者: 尝试获取元素')
      const item = await queue.take()
      console.log(`消费者: 获取到元素 ${item}, 队列大小: ${queue.size}`)
      await new Promise(resolve => setTimeout(resolve, 1000)) // 消费间隔
    }
  }

  // 并行运行生产者和消费者
  await Promise.all([producer(), consumer()])

  console.log('\n演示超时功能:')
  const timeoutQueue = new BlockingQueue<number>(1)

  // 填满队列
  await timeoutQueue.put(999)
  console.log('队列已满')

  try {
    // 尝试添加元素，但设置1秒超时
    console.log('尝试添加元素，设置1秒超时')
    await timeoutQueue.put(1000, 1000)
  } catch (error) {
    console.log(`超时错误: ${(error as Error).message}`)
  }

  console.log('\n演示非阻塞操作:')
  const q = new BlockingQueue<number>(2)

  // 使用offer尝试添加元素
  console.log(`offer(1) 结果: ${q.offer(1)}`)
  console.log(`offer(2) 结果: ${q.offer(2)}`)
  console.log(`offer(3) 结果: ${q.offer(3)}`) // 失败，队列已满

  // 使用poll尝试获取元素
  console.log(`poll() 结果: ${q.poll()}`) // 返回 1
  console.log(`poll() 结果: ${q.poll()}`) // 返回 2
  console.log(`poll() 结果: ${q.poll()}`) // 返回 undefined，队列已空
}

// 运行示例
blockingQueueExample().catch(console.error)
