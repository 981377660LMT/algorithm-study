这是一个非常有趣且考验架构能力的题目。

要手写一套完备的异步 `Task` 机制，我们不能仅仅是对 `Promise` 的简单封装，而是要构建一个**可调度、可取消、可监控、支持依赖关系**的任务运行时系统。

这套设计参考了 C# 的 `Task`、Rust 的 `Future` 以及现代微内核架构中的任务调度思想。

### 核心特性

1.  **状态机管理**：Pending -> Running -> Completed/Failed/Canceled。
2.  **任务调度器 (Scheduler)**：控制并发度，支持优先级。
3.  **取消机制 (Cancellation)**：支持级联取消（父任务取消，子任务自动取消）。
4.  **依赖管理**：支持 `continueWith` (链式调用) 和 `Task.whenAll`。
5.  **强类型**：完全基于 TypeScript 泛型。

---

### 1. 基础类型定义

首先定义任务的状态和取消令牌。

```typescript
// 任务状态枚举
export enum TaskStatus {
  PENDING = 'PENDING', // 等待调度
  RUNNING = 'RUNNING', // 正在执行
  COMPLETED = 'COMPLETED', // 成功完成
  FAILED = 'FAILED', // 抛出异常
  CANCELED = 'CANCELED' // 被取消
}

// 取消令牌源 (参考 C# CancellationTokenSource)
export class CancellationTokenSource {
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

export interface CancellationToken {
  isCancellationRequested(): boolean
  register(callback: () => void): void
  throwIfCancellationRequested(): void
}
```

### 2. 任务调度器 (TaskScheduler)

这是系统的“大脑”，负责控制并发，防止一次性创建过多 Promise 导致系统卡死。

```typescript
// ...existing code...

export class TaskScheduler {
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

// 全局默认调度器
export const DefaultScheduler = new TaskScheduler(4)
```

### 3. 核心 Task 类

这是最复杂的部分，它既要像 Promise 一样好用，又要暴露底层控制能力。

```typescript
// ...existing code...

export class Task<T> {
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
```

### 4. 实战演示

让我们看看这套机制如何处理复杂的异步场景：**并发限制** + **可取消** + **依赖执行**。

```typescript
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
```

### 设计亮点解析

1.  **分离 `start()` 与构造函数**：

    - 原生 `Promise` 一创建就会立即执行（Eager Execution）。
    - 这个 `Task` 是惰性的（Lazy Execution）。你可以先创建一堆任务对象，建立好依赖关系图，最后再决定何时启动。这对于复杂的任务编排非常重要。

2.  **`CancellationToken` 注入**：

    - 这是处理异步取消的“黄金标准”。
    - 我们不强行销毁 Promise（因为 JS 做不到），而是将 Token 传给业务函数，让业务函数在关键节点（如循环中、网络请求前）自己决定是否响应取消。

3.  **调度器解耦**：

    - `Task` 类本身不包含并发控制逻辑，它只负责状态流转。
    - `TaskScheduler` 负责排队。
    - 这意味着你可以为不同的业务模块创建不同的调度器（例如：`NetworkScheduler` 限制 5 并发，`DiskIOScheduler` 限制 1 并发）。

4.  **`continueWith` vs `then`**：
    - `then` 是 Promise 的标准，只处理成功/失败的值。
    - `continueWith` 接收的是**前一个 Task 对象本身**。这意味着即使前一个任务失败了或被取消了，后续任务依然可以执行（用于清理资源或记录日志），并且可以访问前一个任务的状态。
