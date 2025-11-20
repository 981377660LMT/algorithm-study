Batching（批处理）的抽象模式核心在于**“积攒”**和**“触发”**。它将多个细粒度的请求暂时缓存，达到一定条件（数量阈值或时间阈值）后，合并为一个批量请求进行处理。

在 TypeScript 中，最通用的 Batching 抽象通常包含以下要素：

1.  **Buffer（缓冲区）**：暂存待处理的数据。
2.  **Flush Strategy（刷新策略）**：
    - **Size Limit**：攒够 N 个就发。
    - **Time Limit**：每隔 T 毫秒发一次（无论攒了多少）。
3.  **Executor（执行器）**：实际处理批量数据的函数。

以下是一个通用的 Batching 抽象类实现：

```typescript
/**
 * 批处理配置选项
 */
interface BatcherOptions<T> {
  /** 批次最大数量 (达到此数量立即触发处理) */
  maxSize: number
  /** 最长等待时间 (ms) (达到此时间强制触发处理) */
  maxWait: number
}

/**
 * 批处理抽象类
 */
export class Batcher<T> {
  private buffer: T[] = []
  private timer: NodeJS.Timeout | null = null

  constructor(
    private processBatch: (items: T[]) => Promise<void> | void,
    private options: BatcherOptions<T>
  ) {}

  /**
   * 添加单个任务到队列
   */
  public add(item: T): void {
    this.buffer.push(item)

    // 1. 检查是否达到数量阈值 -> 立即触发
    if (this.buffer.length >= this.options.maxSize) {
      this.flush()
    }
    // 2. 如果是队列中的第一个元素 -> 启动计时器
    else if (!this.timer) {
      this.timer = setTimeout(() => {
        this.flush()
      }, this.options.maxWait)
    }
  }

  /**
   * 强制执行当前批次
   */
  public flush(): void {
    if (this.timer) {
      clearTimeout(this.timer)
      this.timer = null
    }

    if (this.buffer.length === 0) {
      return
    }

    // 取出当前所有数据，并清空缓冲区
    const itemsToProcess = [...this.buffer]
    this.buffer = []

    // 执行批量处理逻辑
    // 注意：这里通常不等待处理完成，以免阻塞后续 add 操作
    // 如果需要错误处理，可以在 processBatch 内部 catch
    try {
      this.processBatch(itemsToProcess)
    } catch (e) {
      console.error('Batch processing failed', e)
    }
  }
}
```

### 使用示例

假设我们要记录日志，为了不频繁调用后端接口，我们使用 Batching 模式：

```typescript
// 模拟批量上传接口
const uploadLogs = async (logs: string[]) => {
  console.log(`[Network] Uploading ${logs.length} logs:`, logs)
  // await api.post('/logs', { logs });
}

// 创建 Batcher 实例
// 策略：攒够 5 条日志发一次，或者每隔 2 秒发一次
const logBatcher = new Batcher<string>(uploadLogs, {
  maxSize: 5,
  maxWait: 2000
})

// 模拟高频调用
logBatcher.add('Log 1') // 启动 2s 计时器
logBatcher.add('Log 2')
logBatcher.add('Log 3')

// ... 1秒后 ...
setTimeout(() => {
  logBatcher.add('Log 4')
  logBatcher.add('Log 5') // 触发 maxSize，立即打印 "Uploading 5 logs..."
}, 1000)

// ... 3秒后 ...
setTimeout(() => {
  logBatcher.add('Log 6') // 启动新计时器
  // 2秒后没有新数据，触发 maxWait，打印 "Uploading 1 logs..."
}, 3000)
```

### 进阶模式：带返回值的 Batching (DataLoader 模式)

上面的模式适用于“发后即忘”（Fire-and-forget）的场景（如埋点、日志）。如果你需要**等待**批量处理的结果（例如 GraphQL 中的 N+1 问题优化），则需要结合 `Promise`。

```typescript
/**
 * 带返回值的批处理任务
 */
interface Task<Input, Output> {
  input: Input
  resolve: (value: Output) => void
  reject: (reason?: any) => void
}

export class PromiseBatcher<Input, Output> {
  private queue: Task<Input, Output>[] = []
  private timer: NodeJS.Timeout | null = null

  constructor(
    // 执行器接收输入数组，必须返回等长的结果数组
    private batchFn: (inputs: Input[]) => Promise<(Output | Error)[]>,
    private options: { maxWait: number; maxSize: number }
  ) {}

  public load(input: Input): Promise<Output> {
    return new Promise((resolve, reject) => {
      this.queue.push({ input, resolve, reject })

      if (this.queue.length >= this.options.maxSize) {
        this.flush()
      } else if (!this.timer) {
        this.timer = setTimeout(() => this.flush(), this.options.maxWait)
      }
    })
  }

  private async flush() {
    if (this.timer) {
      clearTimeout(this.timer)
      this.timer = null
    }
    if (this.queue.length === 0) return

    const currentBatch = [...this.queue]
    this.queue = []

    const inputs = currentBatch.map(task => task.input)

    try {
      // 执行批量请求
      const results = await this.batchFn(inputs)

      // 将结果分发回对应的 Promise
      currentBatch.forEach((task, index) => {
        const result = results[index]
        if (result instanceof Error) {
          task.reject(result)
        } else {
          task.resolve(result)
        }
      })
    } catch (error) {
      // 整体失败
      currentBatch.forEach(task => task.reject(error))
    }
  }
}
```
