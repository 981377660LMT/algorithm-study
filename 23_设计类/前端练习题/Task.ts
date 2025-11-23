// 可调度、可取消、可监控、支持依赖关系的任务运行时系统
// C# 的 Task、Rust 的 Future 以及现代微内核架构中的任务调度思想
// Task、Schedular、Cancellation
//
// - 状态机管理：Pending -> Running -> Completed/Failed/Canceled。
// - 任务调度器 (Scheduler)：控制并发度，支持优先级。
// - 取消机制 (Cancellation)：支持级联取消（父任务取消，子任务自动取消）。
// - 依赖管理：支持 continueWith (链式调用) 和 Task.whenAll。

export {}

enum TaskStatus {
  PENDING = 'PENDING', // 等待调度
  RUNNING = 'RUNNING', // 正在执行
  COMPLETED = 'COMPLETED', // 成功完成
  FAILED = 'FAILED', // 抛出异常
  CANCELED = 'CANCELED' // 被取消
}

class CancellationTokenSource {
  private _isCancellationRequested = false
  private _listeners: Array<() => void> = []

  get token(): CancellationToken {
    return {
      isCancellationRequested: () => this._isCancellationRequested,
      register: (cb: () => void) => this._listeners.push(cb),
      throwIfCancellationRequested: () => {
        if (this._isCancellationRequested) throw new Error('TaskCanceledException')
      }
    }
  }

  cancel() {
    if (this._isCancellationRequested) return
    this._isCancellationRequested = true
    this._listeners.forEach(cb => cb())
    this._listeners = []
  }
}

interface CancellationToken {
  isCancellationRequested(): boolean
  register(callback: () => void): void
  throwIfCancellationRequested(): void
}

class TaskScheduler {
  private _queue: Array<() => void> = []
  private _activeCount = 0
  private _maxConcurrency: number

  constructor(maxConcurrency: number = 4) {
    this._maxConcurrency = maxConcurrency
  }

  // 提交任务到调度器
  schedule(run: () => void) {
    this._queue.push(run)
    this._processQueue()
  }

  private _processQueue() {
    if (this._activeCount >= this._maxConcurrency || this._queue.length === 0) {
      return
    }

    const run = this._queue.shift()
    if (run) {
      this._activeCount++
      // 这里的 run 实际上是 Task 内部封装的执行逻辑
      // 我们不关心它何时结束，只关心它通知我们释放资源
      run()
    }
  }

  // 任务完成（无论成功失败）后调用，释放槽位
  notifyComplete() {
    this._activeCount--
    this._processQueue()
  }
}

const DefaultScheduler = new TaskScheduler(4)

class Task<T> {
  private _status: TaskStatus = TaskStatus.PENDING
  private _result: T | undefined
  private _error: any

  // 真正的执行逻辑
  private _action: (token: CancellationToken) => Promise<T> | T

  // 内部 Promise，用于 await 兼容
  private _completionSource: {
    resolve: (value: T | PromiseLike<T>) => void
    reject: (reason?: any) => void
  }
  private _promise: Promise<T>

  // 取消控制
  private _cancellationTokenSource = new CancellationTokenSource()

  // 调度器引用
  private _scheduler: TaskScheduler

  constructor(
    action: (token: CancellationToken) => Promise<T> | T,
    scheduler: TaskScheduler = DefaultScheduler
  ) {
    this._action = action
    this._scheduler = scheduler

    // 创建一个受控的 Promise
    let resolve!: (value: T | PromiseLike<T>) => void
    let reject!: (reason?: any) => void
    this._promise = new Promise<T>((res, rej) => {
      resolve = res
      reject = rej
    })
    this._completionSource = { resolve, reject }
  }

  // --- 公共 API ---

  public get status() {
    return this._status
  }
  public get result() {
    return this._result
  }
  public get error() {
    return this._error
  }

  /**
   * 启动任务
   * 只有调用 start 后，任务才会进入调度队列
   */
  public start(): this {
    if (this._status !== TaskStatus.PENDING) {
      console.warn('Task already started or finished.')
      return this
    }

    // 提交给调度器
    this._scheduler.schedule(() => this._execute())
    return this
  }

  /**
   * 取消任务
   */
  public cancel() {
    if (this._status === TaskStatus.COMPLETED || this._status === TaskStatus.FAILED) return

    this._status = TaskStatus.CANCELED
    this._cancellationTokenSource.cancel()

    // 如果还在排队没运行，直接拒绝 Promise
    // 如果正在运行，Token 会通知业务逻辑
    this._completionSource.reject(new Error('TaskCanceledException'))
  }

  /**
   * 兼容 await 语法
   */
  public then<TResult1 = T, TResult2 = never>(
    onfulfilled?: ((value: T) => TResult1 | PromiseLike<TResult1>) | null,
    onrejected?: ((reason: any) => TResult2 | PromiseLike<TResult2>) | null
  ): Promise<TResult1 | TResult2> {
    // 如果用户直接 await 一个没 start 的任务，自动 start
    if (this._status === TaskStatus.PENDING) {
      this.start()
    }
    return this._promise.then(onfulfilled, onrejected)
  }

  /**
   * 链式调用：当前任务完成后执行下一个任务
   */
  public continueWith<U>(continuationAction: (prevTask: Task<T>) => U | Promise<U>): Task<U> {
    return new Task<U>(async token => {
      // 等待前置任务结束（无论成功失败）
      try {
        await this._promise
      } catch (e) {
        // 忽略前置错误，传给 continuation 处理
      }

      token.throwIfCancellationRequested()
      return continuationAction(this)
    }, this._scheduler)
  }

  // --- 静态方法 ---

  static run<T>(action: () => Promise<T> | T): Task<T> {
    const task = new Task(action)
    task.start()
    return task
  }

  static delay(ms: number): Task<void> {
    return new Task(
      token =>
        new Promise<void>((resolve, reject) => {
          const timer = setTimeout(() => resolve(), ms)
          token.register(() => {
            clearTimeout(timer)
            reject(new Error('TaskCanceledException'))
          })
        })
    )
  }

  // --- 内部执行逻辑 ---

  private async _execute() {
    // 再次检查取消状态（可能在排队时被取消了）
    if (this._cancellationTokenSource.token.isCancellationRequested()) {
      this._status = TaskStatus.CANCELED
      this._completionSource.reject(new Error('TaskCanceledException'))
      this._scheduler.notifyComplete()
      return
    }

    this._status = TaskStatus.RUNNING

    try {
      // 执行用户逻辑，注入 token
      const result = await this._action(this._cancellationTokenSource.token)

      this._status = TaskStatus.COMPLETED
      this._result = result
      this._completionSource.resolve(result)
    } catch (err: any) {
      if (err.message === 'TaskCanceledException') {
        this._status = TaskStatus.CANCELED
      } else {
        this._status = TaskStatus.FAILED
        this._error = err
      }
      this._completionSource.reject(err)
    } finally {
      // 通知调度器释放资源
      this._scheduler.notifyComplete()
    }
  }
}

{
  // 模拟一个耗时操作
  function downloadFile(
    url: string,
    duration: number
  ): (token: CancellationToken) => Promise<string> {
    return async token => {
      console.log(`[Start] Downloading ${url}...`)

      // 模拟分片下载，每 100ms 检查一次取消状态
      for (let i = 0; i < duration / 100; i++) {
        await new Promise(r => setTimeout(r, 100))
        // 关键：业务逻辑中主动检查取消
        token.throwIfCancellationRequested()
      }

      console.log(`[Done] ${url}`)
      return `Content of ${url}`
    }
  }

  async function main() {
    // 1. 创建一个并发度为 2 的调度器
    const limitedScheduler = new TaskScheduler(2)

    // 2. 创建 3 个任务（超过并发度，Task 3 会排队）
    const t1 = new Task(downloadFile('File A', 1000), limitedScheduler)
    const t2 = new Task(downloadFile('File B', 2000), limitedScheduler)
    const t3 = new Task(downloadFile('File C', 1000), limitedScheduler)

    console.log('--- Scheduling Tasks ---')
    t1.start()
    t2.start()
    t3.start() // 此时 t1, t2 运行，t3 排队

    // 3. 演示取消机制
    setTimeout(() => {
      console.log('!!! Canceling File B !!!')
      t2.cancel()
      // t2 被取消后，调度器会释放槽位，t3 应该立即开始
    }, 500)

    // 4. 演示链式调用 (ContinueWith)
    t1.continueWith(prev => {
      console.log(`File A finished with status: ${prev.status}`)
      return 'Processing A'
    }).start()

    // 等待所有结果（为了演示效果）
    try {
      await Promise.allSettled([t1, t2, t3])
    } catch (e) {}

    console.log('--- Final Status ---')
    console.log('T1:', t1.status)
    console.log('T2:', t2.status) // 应该是 CANCELED
    console.log('T3:', t3.status)
  }

  main()
}
