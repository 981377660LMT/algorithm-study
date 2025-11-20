// 轮询（Polling）的抽象模式通常包含三个核心要素：执行逻辑（Task）、调度策略（Scheduler/Interval）和终止条件（Termination Condition）。
// 在 TypeScript 中，最稳健的轮询抽象通常采用 递归的 setTimeout（而非 setInterval，以防止异步任务执行时间超过间隔导致的任务堆积）配合 Promise 或 Class 来管理状态。

/**
 * 轮询配置选项
 */
interface PollingOptions<T> {
  /** 轮询间隔 (ms) */
  interval: number
  /** 最大尝试次数 (可选，默认无限) */
  maxAttempts?: number
  /** 超时时间 (ms) (可选，默认无限) */
  timeout?: number
  /** 终止条件检查函数：返回 true 则停止轮询 */
  stopCondition?: (data: T) => boolean
  /** 错误处理：返回 true 则继续轮询，false 则终止抛出异常 */
  onError?: (err: any) => boolean
}

/**
 * 轮询器抽象类
 */
export class Poller<T> {
  private timer: NodeJS.Timeout | null = null
  private attempts: number = 0
  private startTime: number = 0
  private isRunning: boolean = false
  private _promise: Promise<T> | null = null
  private _resolve: ((value: T) => void) | null = null
  private _reject: ((reason?: any) => void) | null = null

  constructor(private task: () => Promise<T>, private options: PollingOptions<T>) {}

  /**
   * 开始轮询
   */
  public start(): Promise<T> {
    if (this.isRunning) {
      return this._promise!
    }

    this.isRunning = true
    this.attempts = 0
    this.startTime = Date.now()

    // 创建一个控制整个轮询生命周期的 Promise
    this._promise = new Promise<T>((resolve, reject) => {
      this._resolve = resolve
      this._reject = reject
      this.tick()
    })

    return this._promise
  }

  /**
   * 停止轮询
   */
  public stop(): void {
    this.isRunning = false
    if (this.timer) {
      clearTimeout(this.timer)
      this.timer = null
    }
  }

  /**
   * 单次执行逻辑 (递归核心)
   */
  private async tick() {
    if (!this.isRunning) return

    // 1. 检查最大尝试次数
    if (this.options.maxAttempts && this.attempts >= this.options.maxAttempts) {
      this.fail(new Error('Polling reached max attempts'))
      return
    }

    // 2. 检查超时
    if (this.options.timeout && Date.now() - this.startTime > this.options.timeout) {
      this.fail(new Error('Polling timed out'))
      return
    }

    this.attempts++

    try {
      // 3. 执行任务
      const result = await this.task()

      // 4. 检查是否满足终止条件 (成功)
      if (this.options.stopCondition && this.options.stopCondition(result)) {
        this.succeed(result)
        return
      }

      // 5. 未满足条件，调度下一次
      this.scheduleNext()
    } catch (error) {
      // 6. 错误处理
      const shouldContinue = this.options.onError ? this.options.onError(error) : false
      if (shouldContinue) {
        this.scheduleNext()
      } else {
        this.fail(error)
      }
    }
  }

  private scheduleNext() {
    if (!this.isRunning) return
    this.timer = setTimeout(() => this.tick(), this.options.interval)
  }

  private succeed(value: T) {
    this.stop()
    this._resolve?.(value)
  }

  private fail(error: any) {
    this.stop()
    this._reject?.(error)
  }
}

// 模拟一个异步请求
const mockApi = async () => {
  const status = Math.random() > 0.8 ? 'COMPLETED' : 'PENDING'
  console.log(`Checking status... ${status}`)
  return { status, data: 'some result' }
}

// 创建轮询实例
const poller = new Poller(mockApi, {
  interval: 1000, // 1秒一次
  maxAttempts: 10, // 最多试10次
  timeout: 5000, // 5秒超时
  // 当状态为 COMPLETED 时停止
  stopCondition: res => res.status === 'COMPLETED',
  // 发生错误时打印日志并继续重试
  onError: err => {
    console.error('Retry on error:', err)
    return true
  }
})

// 启动
poller
  .start()
  .then(finalResult => {
    console.log('Polling success:', finalResult)
  })
  .catch(err => {
    console.error('Polling failed:', err)
  })
